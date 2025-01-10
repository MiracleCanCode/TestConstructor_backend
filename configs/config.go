package configs

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DB     string
	PORT   string
	SECRET string
}

func Load(log *zap.Logger) *Config {

	if err := godotenv.Load(".env.local"); err != nil {
		log.Error("Failed get env file, err:" + err.Error())
	}
	db := os.Getenv("DB")
	port := os.Getenv("PORT")

	return &Config{
		DB:     db,
		PORT:   port,
		SECRET: "SUPERSECRETKEYFORBESTAPPINTHEWORLD",
	}
}
