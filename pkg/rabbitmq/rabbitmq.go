package rabbitmq

import (
	"fmt"

	logger "github.com/MiracleCanCode/example_configuration_logger"
	"github.com/server/configs"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Rabbitmq struct {
	conn *amqp.Connection
}

func New() *Rabbitmq {
	log := logger.Logger(logger.DefaultLoggerConfig())
	cfg := configs.NewRabbitmqConfig(log)
	rabbitmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		log.Error("Failed to connect to RabbitMQ: %s", zap.Error(err))
	}

	return &Rabbitmq{
		conn: conn,
	}
}
