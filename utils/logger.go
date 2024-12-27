package utils

import (
	"github.com/sirupsen/logrus"
)

// Initialize the logger with a JSON formatter
func InitLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

// LogError logs an error message
func LogError(message string, err error) {
	logrus.WithFields(logrus.Fields{
		"error": err,
	}).Error(message)
}

// LogWarn logs a warning message
func LogWarn(message string) {
	logrus.Warn(message)
}

// LogInfo logs an informational message
func LogInfo(message string) {
	logrus.Info(message)
}
