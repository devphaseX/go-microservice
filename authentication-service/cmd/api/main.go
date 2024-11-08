package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	db "github.com/devphaseX/go-microservice/authenication-service/db/sqlc"
	// This import is required for the file:// source driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Config struct {
	tokenMaker TokenMaker
	env        *AppEnvConfig
	store      db.Store
	hash       *Argon2idHash
}

func newConfig(store db.Store, env *AppEnvConfig) (*Config, error) {
	pasetoMaker, err := NewPasetoMaker(env.SymmetricKey)
	if err != nil {
		return nil, err
	}

	hash := DefaultArgonHash()
	return &Config{
		store:      store,
		tokenMaker: pasetoMaker,
		env:        env,
		hash:       hash,
	}, nil
}

func main() {
	appEnvConfig, err := LoanEnv(".")

	if err != nil {
		log.Fatal(err)
	}

	if err := createDatabase(appEnvConfig); err != nil {
		log.Panic("failed to create db", err)
	}

	if err := runMigrations(appEnvConfig); err != nil {
		log.Panic(err)
	}

	dbConn := connect(appEnvConfig.DbSource, appEnvConfig.DbMaxRetryCount)
	config, err := newConfig(db.NewStore(dbConn), appEnvConfig)

	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", appEnvConfig.Addr),
		Handler: config.routes(),
	}

	fmt.Printf("server listening on port: %s", appEnvConfig.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(DbSource string) (*pgxpool.Pool, error) {
	pgConfig, err := pgxpool.ParseConfig(DbSource)
	if err != nil {
		return nil, errors.New("Unable to parse connection string")
	}

	// Set some reasonable pool limits
	pgConfig.MaxConns = 20
	pgConfig.MinConns = 2
	pgConfig.MaxConnLifetime = time.Hour
	pgConfig.MaxConnIdleTime = 30 * time.Minute

	// Set some reasonable timeouts
	pgConfig.ConnConfig.ConnectTimeout = 5 * time.Second
	pgConfig.ConnConfig.RuntimeParams["statement_timeout"] = "30000" // 30 seconds

	conn, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	return conn, err
}

func connect(DbSource string, dbMaxRetryCount int) *pgxpool.Pool {
	dbFailRetryAttemptCount := 0

	for {
		conn, err := openDB(DbSource)

		if err != nil {
			log.Printf("cannot connect to DB: %s", err)
		} else {
			log.Println("connection to DB sucessful")
			return conn
		}

		if dbFailRetryAttemptCount < dbMaxRetryCount {
			time.Sleep(time.Second * 2)
			continue
		}

		return nil
	}
}

func createDatabase(config *AppEnvConfig) error {
	// Connect to postgres database to create new database
	// Modify connection string to connect to 'postgres' database instead
	connURL, err := url.Parse(config.DbSource)
	if err != nil {
		return fmt.Errorf("error parsing DB source URL: %v", err)
	}

	// Override the database name to connect to postgres
	connURL.Path = "/postgres"
	db, err := sql.Open("postgres", connURL.String())

	if err != nil {
		return fmt.Errorf("error connecting to postgres: %v", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", config.DbName)
	fmt.Print(query)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if database exists: %v", err)
	}

	if !exists {
		// Create database if it doesn't exist
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DbName))
		if err != nil {
			return fmt.Errorf("error creating database: %v", err)
		}
		log.Printf("Database %s created successfully", config.DbName)
	} else {
		log.Printf("Database %s already exists", config.DbName)
	}

	return nil
}

func runMigrations(config *AppEnvConfig) error {
	// Connect to the newly created database
	db, err := sql.Open("postgres", config.DbSource)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error creating migration driver: %v", err)
	}

	// Check if migrations directory exists
	if _, err := os.Stat(config.MigrationsPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", config.MigrationsPath)
	}

	// Initialize migrations
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", config.MigrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %v", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
