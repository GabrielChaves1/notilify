package logger

import (
	"GabrielChaves1/notilify/internal/domain/types"
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Environment types.Environment
}

func NewLogger(config Config) *logrus.Entry {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	
	standardFields := logrus.Fields{
		"app":     "notilify",
		"version": "1.0.0",
		"environment": config.Environment,
	}

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})

	if config.Environment == types.Development {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.WithFields(standardFields).Info("Logger initialized")
	return logger.WithFields(standardFields)
}