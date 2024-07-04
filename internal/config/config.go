package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     uint   `mapstructure:"DB_PORT"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	HttpServer  string `mapstructure:"HTTP_SERVER"`
	Salt        string `mapstructure:"APP_SALT"`
	SigningKey  string `mapstructure:"APP_SIGNING_KEY"`
	GRPCAddress string `mapstructure:"GRPC_ADDRESS"`
}

func NewConfig() (*Config, error) {
	var config Config

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("Error reading env file", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	logrus.Debug("Parsed config values")
	logrus.Debugln(config)

	return &config, nil
}
