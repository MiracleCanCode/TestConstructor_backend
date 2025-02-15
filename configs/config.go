package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DB           string
	PORT         string
	SECRET       string
	CLIENT_URL   string
	REDIS_HOST   string
	RABBITMQ_URL string
}

func Load(log *zap.Logger) (*Config, error) {
	envFile := ".env.production"

	if err := godotenv.Load(envFile); err != nil {
		if os.IsNotExist(err) {
			log.Warn("Env file not found", zap.String("file", envFile))
		} else {
			log.Error("Failed to load env file", zap.String("file", envFile), zap.Error(err))
			return nil, fmt.Errorf("Load: failed to load %s: %w", envFile, err)
		}
	}

	db, ok := os.LookupEnv("DB")
	if !ok {
		log.Error("DB env variable not set")
		return nil, fmt.Errorf("Load: DB env variable not set")
	}

	rabbit, ok := os.LookupEnv("RABBITMQ_URL")
	if !ok {
		log.Error("RABBITMQ_URL env variable not set")
		return nil, fmt.Errorf("Load: RABBITMQ_URL env variable not set")
	}

	clientUrl, ok := os.LookupEnv("CLIENT_URL")
	if !ok {
		log.Error("CLIENT_URL env variable not set")
		return nil, fmt.Errorf("Load: CLIENT_URL env variable not set")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Error("PORT env variable not set")
		return nil, fmt.Errorf("Load: PORT env variable not set")
	}

	secret, ok := os.LookupEnv("SECRET")
	if !ok {
		log.Warn("SECRET env variable not set. Using default secret!")
		secret = "SUPERSECRETKEYFORBESTAPPINTHEWORLD"
	}

	redis, ok := os.LookupEnv("REDIS_HOST")
	if !ok {
		log.Error("REDIS_HOST env variable not set")
		return nil, fmt.Errorf("Load: REDIS_HOST env variable not set")
	}

	return &Config{
		DB:           db,
		PORT:         port,
		SECRET:       secret,
		CLIENT_URL:   clientUrl,
		REDIS_HOST:   redis,
		RABBITMQ_URL: rabbit,
	}, nil
}
