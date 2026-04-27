package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var Log *zap.Logger
var err error
var logConfig zap.Config

func InitLogger(env string, appName string) {
	switch env {
	case "production":
		logConfig = zap.NewProductionConfig()
		logConfig.EncoderConfig.TimeKey = "timestamp"
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	case "staging":
		logConfig = zap.NewProductionConfig()
		logConfig.EncoderConfig.TimeKey = "timestamp"
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	case "development", "local", "testing":
		logConfig = zap.NewDevelopmentConfig()
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	default:
		logConfig = zap.NewDevelopmentConfig()
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logConfig.InitialFields = map[string]interface{}{
		"app": appName,
		"env": env,
	}

	Log, err = logConfig.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

func Sync() {
	if Log != nil {
		if err := Log.Sync(); err != nil {
			log.Println("failed to sync logger:", err)
		}
	}
}
