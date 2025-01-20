package configs

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DB         string
	PORT       string
	SECRET     string
	PRODACTION bool
}

func Load(log *zap.Logger) *Config {

	if err := godotenv.Load(".env.local"); err != nil {
		log.Error("Failed get env file", zap.Error(err))
	}
	db := os.Getenv("DB")
	port := os.Getenv("PORT")

	return &Config{
		DB:         db,
		PORT:       port,
		SECRET:     "SUPERSECRETKEYFORBESTAPPINTHEWORLD",
		PRODACTION: false,
	}
}
