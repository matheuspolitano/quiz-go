package config

import (
	"github.com/spf13/viper"
)

// Config holds all configuration values for our application.
type Config struct {
	API_URL string `mapstructure:"API_URL"`
}

// LoadConfig reads configuration from:
// 1) A file named "app.env" in the specified path
// 2) Environment variables (override file)
// And unmarshals into the Config struct.
func LoadConfig(path string) (Config, error) {
	var cfg Config
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
