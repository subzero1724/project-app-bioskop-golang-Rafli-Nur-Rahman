package logger

import (
	"cinema-booking-system/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new Zap logger
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	// Parse log level
	level := zapcore.InfoLevel
	switch cfg.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	// Create logger config
	var zapConfig zap.Config

	if cfg.Log.Encoding == "json" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	zapConfig.Level = zap.NewAtomicLevelAt(level)
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Build logger
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
