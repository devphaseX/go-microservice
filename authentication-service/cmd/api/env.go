package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type AppEnvConfig struct {
	Addr             string        `mapstructure:"LISTEN_ADDR"`
	SymmetricKey     string        `mapstructure:"SYMMETRIC_KEY"`
	DbMaxRetryCount  int           `mapstructure:"DB_MAX_RETRY_COUNT"`
	DbSource         string        `mapstructure:"DB_SOURCE"`
	MigrationsPath   string        `mapstructure:"MIGRATIONS_PATH"`
	DbName           string        `mapstructue:"DBNAME"`
	AccessTokenTime  time.Duration `mapstructure:"ACCESS_TOKEN_TIME"`
	RefreshTokenTime time.Duration `mapstructure:"REFRESH_TOKEN_TIME"`
}

func LoanEnv(path string) (envConfig *AppEnvConfig, err error) {
	vp := viper.New()

	vp.AddConfigPath(path)
	vp.SetConfigName("app")
	vp.SetConfigType("env")

	vp.AutomaticEnv()

	fmt.Printf("Looking for config in: %s\n", path)
	if err = vp.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return nil, err
	}
	fmt.Printf("Successfully read config file\n")

	err = vp.Unmarshal(&envConfig)
	return
}
