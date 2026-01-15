package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Env struct {
	DbName     string `env:"DB_NAME"`
	DbUsername string `env:"DB_USERNAME"`
	DbPassword string `env:"DB_PASSWORD"`
	DbPort     int    `env:"DB_PORT"`
	DbHost     string `env:"DB_HOST"`
	IpAddress  string `env:"IP_ADDRESS"`
	ApiPort    int    `env:"API_PORT"`

	// JWT
	JWTSecret             string `env:"JWT_SECRET"`
	JWTExpireHours        int    `env:"JWT_EXPIRE_HOURS"`
	JWTRefreshExpireHours int    `env:"JWT_REFRESH_EXPIRE_HOURS"`

	Environment string `env:"ENVIRONMENT"`
}

type Config struct {
	Env    Env
	Client *pgxpool.Pool
}

var config Config

func GetConfig() *Config {
	config.Env = *getEnv()

	return &config
}

func getEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found, using system environment variables")
	}

	var cfg Env
	err = env.Parse(&cfg)
	if err != nil {
		slog.Error("Failed to parse env: %v", "error", err.Error())
		panic(err)
	}

	return &cfg
}
