package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// Initialize the logger with a JSON formatter
func InitLogger() {
	Logger = logrus.New()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Log to a file instead of stdout
	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Logger.SetOutput(file)
	} else {
		Logger.Info("Failed to log to file, using default stderr")
	}
}

// LogError logs an error message
func LogError(source, message string, err error) {
	logrus.WithFields(logrus.Fields{
		"source": source,
		"error":  err,
	}).Error(message)
}

// LogWarn logs a warning message
func LogWarn(source, message string) {
	Logger.WithField("source", source).Warn(message)
}

// LogInfo logs an informational message
func LogInfo(source, message string) {
	// logrus.Info(message)
	Logger.WithField("source", source).Info(message)
}
