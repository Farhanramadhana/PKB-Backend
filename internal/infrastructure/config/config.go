package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AppEnv   string
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret   string
	TTLHours int
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("JWT_TTL_HOURS", 24)
	viper.SetDefault("APP_ENV", "development")

	return &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			URL: viper.GetString("DATABASE_URL"),
		},
		JWT: JWTConfig{
			Secret:   viper.GetString("JWT_SECRET"),
			TTLHours: viper.GetInt("JWT_TTL_HOURS"),
		},
		AppEnv: viper.GetString("APP_ENV"),
	}
}
