package configs

import (
	"os"

	"github.com/MiracleCanCode/zaperr"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DB     string
	PORT   string
	SECRET string
}

func LoadConfig(log *zap.Logger, handleErrors *zaperr.Zaperr) *Config {

	handleErrors.LogPanicError(godotenv.Load(".env.local"), "Failed get env file")

	db := os.Getenv("DB")
	port := os.Getenv("PORT")
	secret := os.Getenv("SECRET")

	return &Config{
		DB:     db,
		PORT:   port,
		SECRET: secret,
	}
}
