package config

import (
	"github.com/spf13/viper"
)

// Config type to manage main application setup
type Config struct {
	ApiPort         string `mapstructure:"API_PORT"`
	ApiTimeShutdown int    `mapstructure:"API_TIME_SHUTDOWN"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {

			return
		}
		err = nil
	}
	err = viper.Unmarshal(&config)
	return
}
