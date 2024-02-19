package logging

import "go.uber.org/zap"

// Get returns a new zap logger
func Get() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}
