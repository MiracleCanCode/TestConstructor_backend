package queuemanager

import (
	"errors"
	"github.com/rabbitmq/amqp091-go"
	"github.com/server/pkg/rabbitmq"
	"go.uber.org/zap"
)

type QueueManagerInterface interface {
	PublishMessage(body []byte, name string, contentType string) error
	ConsumeMessages(name string, handler func(message []byte) error) error
}

type QueueManager struct {
	rabbitmq *rabbitmq.Rabbitmq
	logger   *zap.Logger
}

func New(logger *zap.Logger) *QueueManager {
	return &QueueManager{
		rabbitmq: rabbitmq.New(),
		logger:   logger,
	}
}

func (s *QueueManager) PublishMessage(body []byte, name string, contentType string) error {
	rabbit := s.rabbitmq.Conn
	if rabbit == nil {
		s.logger.Error("RabbitMQ connection is nil")
		return errors.New("rabbitmq connection is nil")
	}

	channRabbit, err := rabbit.Channel()
	if err != nil {
		s.logger.Error("Failed to create rabbitmq channel", zap.Error(err))
		return err
	}
	defer channRabbit.Close()

	_, err = channRabbit.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.logger.Error("Failed to declare queue", zap.Error(err))
		return err
	}

	err = channRabbit.Publish(
		"",
		name,
		false,
		false,
		amqp091.Publishing{
			ContentType: contentType,
			Body:        body,
		},
	)
	if err != nil {
		s.logger.Error("Failed to publish message", zap.Error(err))
		return err
	}

	s.logger.Info("Message published successfully", zap.String("queue", name))
	return nil
}

func (s *QueueManager) ConsumeMessages(name string, handler func(message []byte) error) error {
	rabbit := s.rabbitmq.Conn
	if rabbit == nil {
		s.logger.Error("RabbitMQ connection is nil")
		return errors.New("rabbitmq connection is nil")
	}

	channRabbit, err := rabbit.Channel()
	if err != nil {
		s.logger.Error("Failed to create rabbitmq channel", zap.Error(err))
		return err
	}
	defer channRabbit.Close()

	_, err = channRabbit.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.logger.Error("Failed to declare queue", zap.Error(err))
		return err
	}

	msgs, err := channRabbit.Consume(
		name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.logger.Error("Failed to register a consumer", zap.Error(err))
		return err
	}

	go func() {
		for msg := range msgs {
			s.logger.Info("Received message", zap.String("queue", name), zap.ByteString("message", msg.Body))

			if err := handler(msg.Body); err != nil {
				s.logger.Error("Error processing message", zap.Error(err))
			}
		}
	}()

	s.logger.Info("Started consuming messages", zap.String("queue", name))
	return nil
}
