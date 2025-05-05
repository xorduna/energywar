package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
	URI string `mapstructure:"uri"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

func LoadConfig() (*Config, error) {
	// Set default configuration
	viper.SetDefault("database.uri", "game.db")
	viper.SetDefault("server.port", ":8080")

	// Enable reading from environment variables
	//viper.SetEnvPrefix("GS")
	viper.AutomaticEnv()

	// Bind environment variables
	viper.BindEnv("database.uri", "DB_URI")
	viper.BindEnv("server.port", "SERVER_PORT")

	// Read config file if exists (not an error if it doesn't)
	_ = viper.ReadInConfig()

	// Create config struct
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
