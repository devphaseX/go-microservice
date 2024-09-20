package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	lnAddr          = "80"
	symmetricKey    = ""
	dbMaxRetryCount = 10
)

type Config struct {
	tokenMaker TokenMaker
	env        *AppEnvConfig
}

func newConfig(env *AppEnvConfig) (*Config, error) {
	pasetoMaker, err := NewPasetoMaker(env.SymmetricKey)

	if err != nil {
		return nil, err
	}

	return &Config{
		tokenMaker: pasetoMaker,
		env:        env,
	}, nil
}

func main() {
	appEnvConfig, err := LoanEnv("../../")

	if err != nil {
		log.Fatal(err)
	}

	config, err := newConfig(appEnvConfig)

	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    appEnvConfig.Addr,
		Handler: config.routes(),
	}

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

func connect(DbSource string) *pgxpool.Pool {
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
