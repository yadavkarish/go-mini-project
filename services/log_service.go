package services

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logs []logrus.Entry

// Initialize logger
func InitLogger() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

// Log an entry
func LogEntry(level, message string) {
	entry := logrus.Entry{
		Time:    time.Now(),
		Message: message,
	}
	logLevel, _ := logrus.ParseLevel(level)
	entry.Level = logLevel
	log.WithFields(entry.Data).Log(logLevel, message)
	logs = append(logs, entry)
}

// API Endpoint to fetch logs
func GetLogs(c *gin.Context) {
	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")
	level := c.DefaultQuery("level", "info")

	// Filter logs by date range and level
	filteredLogs := filterLogs(start, end, level)
	c.JSON(http.StatusOK, gin.H{
		"logs": filteredLogs,
	})
}

// Filter logs based on date range and log level
func filterLogs(start, end, level string) []logrus.Entry {
	var filtered []logrus.Entry
	for _, log := range logs {
		if level != "" && log.Level.String() != level {
			continue
		}

		if start != "" && log.Time.Before(parseDate(start)) {
			continue
		}

		if end != "" && log.Time.After(parseDate(end)) {
			continue
		}

		filtered = append(filtered, log)
	}
	return filtered
}

// Helper function to parse date strings
func parseDate(dateStr string) time.Time {
	parsed, _ := time.Parse("2006-01-02", dateStr)
	return parsed
}
