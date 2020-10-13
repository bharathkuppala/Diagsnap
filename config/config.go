package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// AppConfig ...
type AppConfig struct {
	Server ServerConfig
	Minio  MinioConfig
}

// ServerConfig ...
type ServerConfig struct {
	Port         string
	TimeoutRead  time.Duration
	TimeoutWrite time.Duration
}

// MinioConfig ...
type MinioConfig struct {
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioObject    string
}

// LoadEnv ...
func LoadEnv() *AppConfig {
	viper.SetConfigName("config")

	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	var appConfig AppConfig

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return nil
	}

	err := viper.Unmarshal(&appConfig)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		return nil
	}

	fmt.Println(appConfig)

	return &appConfig
}
