package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment         string `mapstructure:"ENVIRONMENT"`
	ServerAddress       string `mapstructure:"SERVER_ADDRESS"`
	GeoDBAddress        string `mapstructure:"GEODB_ADDRESS"`
	GeoDBAPIKey         string `mapstructure:"GEODB_API_KEY"`
	GeoDBAPIHost        string `mapstructure:"GEODB_API_HOST"`
	AerisWeatherAddress string `mapstructure:"AERISWEATHER_API_ADDRESS"`
	AerisWeatherAPIKey  string `mapstructure:"AERISWEATHER_API_KEY"`
	AerisWeatherAPIHost string `mapstructure:"AERISWEATHER_API_HOST"`
	GinMode             string `mapstructure:"GIN_MODE"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
