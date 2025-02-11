package logger

import (
	"sync"

	log "github.com/MiracleCanCode/example_configuration_logger"
	"go.uber.org/zap"
)

var (
	instance *zap.Logger
	once     sync.Once
)

func GetInstance() *zap.Logger {
	once.Do(func() {
		instance = log.Logger(log.DefaultLoggerConfig())
		defer func() {
			if err := instance.Sync(); err != nil {
				instance.Error("Failed to sync logger", zap.Error(err))
			}
		}()
	})

	return instance
}
