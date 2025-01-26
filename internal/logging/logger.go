package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

// Init initializes the logger
func CreateLogger(logLevel zapcore.Level) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(logLevel)
	zlog, _ := config.Build()
	logger = zlog.Sugar()
}

// Get returns a new zap logger
func Get() *zap.SugaredLogger {
	return logger
}
