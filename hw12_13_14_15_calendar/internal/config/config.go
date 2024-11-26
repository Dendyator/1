package config

import (
	"log"

	"github.com/spf13/viper" //nolint:depguard
)

type Config struct {
	Server   ServerConfig
	Logger   LoggerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type LoggerConfig struct {
	Level string
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

func LoadConfig(configPath string) Config {
	var config Config
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	return config
}
