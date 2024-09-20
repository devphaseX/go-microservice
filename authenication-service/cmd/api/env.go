package main

import "github.com/spf13/viper"

type AppEnvConfig struct {
	Addr            string `mapstructure:"LISTEN_ADDR"`
	SymmetricKey    string `mapstructure:"SYMMETRIC_KEY"`
	DbMaxRetryCount int    `mapstructure:"DB_MAX_RETRY_COUNT"`
}

func LoanEnv(path string) (envConfig *AppEnvConfig, err error) {
	vp := viper.New()

	vp.AddConfigPath(path)
	vp.SetConfigName("app")
	vp.SetConfigType("env")

	vp.AutomaticEnv()

	if err = vp.ReadInConfig(); err != nil {
		return nil, err
	}

	err = vp.Unmarshal(&envConfig)
	return
}
