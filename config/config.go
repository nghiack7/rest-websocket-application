package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig() (*viper.Viper, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Set default values for durations
	viper.SetDefault("server.timeout", "10s")
	viper.SetDefault("server.read_timeout", "5s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.shutdown_timeout", "30s")
	viper.SetDefault("auth.jwt_expiration", "24h")
	viper.SetDefault("database.conn_max_lifetime", "5m")

	// Process environment variable substitutions with defaults
	// This handles ${VAR:default} syntax in the config file
	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)
		if len(value) > 4 && value[0:2] == "${" && value[len(value)-1:] == "}" {
			// Extract variable name and default value
			varAndDefault := value[2 : len(value)-1]
			parts := strings.SplitN(varAndDefault, ":", 2)

			envValue := os.Getenv(parts[0])
			if envValue != "" {
				// Use environment variable if set
				viper.Set(key, envValue)
			} else if len(parts) > 1 {
				// Use default value if provided and env var not set
				viper.Set(key, parts[1])
			}
		}
	}
	return viper.GetViper(), nil
}
