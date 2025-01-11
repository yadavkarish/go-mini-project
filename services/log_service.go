package services

import (
	"bufio"
	"csv-microservice/utils"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logs []logrus.Entry

// Initialize logger
// func InitLogger() {
// 	log.SetFormatter(&logrus.JSONFormatter{})
// 	log.SetLevel(logrus.InfoLevel)
// }

// Log an entry
func LogEntry(level, source, message string) {
	entry := logrus.Entry{
		Time:  time.Now(),
		Level: parseLogLevel(level),
		// Data:    logrus.Fields{"source": source},
		Message: message,
		Data: map[string]interface{}{
			"source": source,
		},
	}
	// logLevel, _ := logrus.ParseLevel(level)
	// entry.Level = logLevel
	// log.WithFields(entry.Data).Log(logLevel, message)
	logs = append(logs, entry)
	utils.Logger.WithFields(entry.Data).Log(entry.Level, message)
}

// API Endpoint to fetch logs
func GetLogs(c *gin.Context) {
	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")
	level := c.DefaultQuery("level", "info")

	logFile, err := os.Open("logs.log")
	if err != nil {
		utils.LogError("GetLogs", "Failed to open log file", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to read logs",
		})
		return
	}
	defer logFile.Close()

	// var logs []logrus.Entry
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		line := scanner.Text()
		entry, parseErr := parseLogLine(line)
		fmt.Println(parseErr)
		if parseErr == nil {
			logs = append(logs, entry)
			fmt.Println(entry)
		}
	}

	if err := scanner.Err(); err != nil {
		utils.LogError("GetLogs", "Error reading log file", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to parse logs",
		})
		return
	}

	// Filter logs by date range and level
	filteredLogs := filterLogs(start, end, level)
	c.JSON(http.StatusOK, gin.H{
		"logs": filteredLogs,
	})
}

// Filter logs based on date range and log level
func filterLogs(start, end, level string) []logrus.Entry {
	fmt.Println("Filtering logs")
	var filtered []logrus.Entry
	for _, log := range logs {
		fmt.Println(log.Time)
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

// Parse log level safely
func parseLogLevel(level string) logrus.Level {
	parsed, err := logrus.ParseLevel(level)
	if err != nil {
		return logrus.InfoLevel // Default to INFO
	}
	return parsed
}

// func parseLogLine(line string) (logrus.Entry, error) {
// 	var entry logrus.Entry
// 	err := json.Unmarshal([]byte(line), &entry)
// 	fmt.Println("Error: ", err)
// 	if err != nil {
// 		return logrus.Entry{}, err
// 	}
// 	return entry, nil
// }

func parseLogLine(line string) (logrus.Entry, error) {
	var entry logrus.Entry

	// Parse the log line manually
	var timeStr, level, msg, source string
	_, err := fmt.Sscanf(line, `time=%q level=%s msg=%q source=%s`, &timeStr, &level, &msg, &source)
	if err != nil {
		return logrus.Entry{}, err
	}

	// Parse time
	parsedTime, timeErr := time.Parse(time.RFC3339, timeStr)
	if timeErr != nil {
		return logrus.Entry{}, timeErr
	}

	// Parse log level
	parsedLevel, levelErr := logrus.ParseLevel(level)
	if levelErr != nil {
		return logrus.Entry{}, levelErr
	}

	// Create logrus.Entry
	entry = logrus.Entry{
		Time:    parsedTime,
		Level:   parsedLevel,
		Message: msg,
		Data: map[string]interface{}{
			"source": source,
		},
	}
	return entry, nil
}
