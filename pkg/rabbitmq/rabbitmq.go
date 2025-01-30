package rabbitmq

import (
	logger "github.com/MiracleCanCode/example_configuration_logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/server/configs"
	"go.uber.org/zap"
)

type Rabbitmq struct {
	Conn *amqp.Connection
}

func New() *Rabbitmq {
	log := logger.Logger(logger.DefaultLoggerConfig())
	cfg, err := configs.Load(log)
	if err != nil {
		log.Error("Failed to load config", zap.Error(err))
		return nil
	}
	conn, err := amqp.Dial(cfg.RABBITMQ_URL)
	if err != nil {
		log.Error("Failed to connect to RabbitMQ: %s", zap.Error(err))
	}

	return &Rabbitmq{
		Conn: conn,
	}
}
