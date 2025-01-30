package configs

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type RabbitmqConfig struct {
	Host     string
	Port     string
	User     string
	Queue    string
	Password string
}

func NewRabbitmqConfig(log *zap.Logger) *RabbitmqConfig {
	var envFile string

	if PRODACTION {
		envFile = ".env.prodaction"
	} else {
		envFile = ".env.local"
	}

	if err := godotenv.Load(envFile); err != nil {
		if os.IsNotExist(err) {
			log.Warn("Env file not found", zap.String("file", envFile))
		} else {
			log.Error("Failed to load env file", zap.String("file", envFile), zap.Error(err))
			return nil
		}
	}

	host, ok := os.LookupEnv("RABBITMQ_HOST")
	if !ok {
		log.Error("RABBITMQ_HOST env variable not set")
		return nil
	}

	password, ok := os.LookupEnv("RABBITMQ_PASSWORD")
	if !ok {
		log.Error("RABBITMQ_PASSWORD env variable not set")
		return nil
	}

	user, ok := os.LookupEnv("RABBITMQ_USER")
	if !ok {
		log.Error("RABBITMQ_USER env variable not set")
		return nil
	}

	queue, ok := os.LookupEnv("RABBITMQ_QUEUE")
	if !ok {
		log.Error("RABBITMQ_QUEUE env variable not set")
		return nil
	}

	port, ok := os.LookupEnv("RABBITMQ_PORT")
	if !ok {
		log.Error("RABBITMQ_PORT env variable not set")
		return nil
	}

	return &RabbitmqConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Queue:    queue,
		Password: password,
	}
}
