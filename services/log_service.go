package services

import (
	"bufio"
	"csv-microservice/utils"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Log []logrus.Entry

// Initialize logger
// func InitLogger() {
// 	log.SetFormatter(&logrus.JSONFormatter{})
// 	log.SetLevel(logrus.InfoLevel)
// }

func ResetLogs() {
	Log = []logrus.Entry{}
}

// GetLogsSlice returns the current logs (used for testing purposes).
func GetLogsSlice() []logrus.Entry {
	return Log
}

// Log an entry
func LogEntry(level, source, message string) {
	entry := logrus.Entry{
		Time:  time.Now(),
		Level: ParseLogLevel(level),
		// Data:    logrus.Fields{"source": source},
		Message: message,
		Data: map[string]interface{}{
			"source": source,
		},
	}
	// logLevel, _ := logrus.ParseLevel(level)
	// entry.Level = logLevel
	// log.WithFields(entry.Data).Log(logLevel, message)
	Log = append(Log, entry)
	utils.Logger.WithFields(entry.Data).Log(entry.Level, message)
}

// API Endpoint to fetch logs
func GetLogs(c *gin.Context) {
	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")
	level := c.DefaultQuery("level", "info")

	logFile, err := os.Open("logs.log")
	if err != nil {

		// If the file doesn't exist, return a specific error message
		if os.IsNotExist(err) {
			utils.LogError("GetLogs", "Log file not found", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Log file not found",
			})
			return
		}

		// If it's any other error, log it as a failure to open the file
		utils.LogError("GetLogs", "Failed to open log file", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to read logs",
		})
		return
	}
	defer logFile.Close()

	Log = []logrus.Entry{}
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		line := scanner.Text()
		entry, parseErr := ParseLogLine(line)
		if parseErr != nil {
			utils.LogError("GetLogs", "Failed to parse log line", parseErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to parse logs",
			})
			return
		}
		Log = append(Log, entry)
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
	filteredLogs := FilterLogs(start, end, level)
	c.JSON(http.StatusOK, gin.H{
		"logs": filteredLogs,
	})
}

// Filter logs based on date range and log level
func FilterLogs(start, end, level string) []logrus.Entry {
	fmt.Println("Filtering logs")
	// fmt.Printf("Start: %s, End: %s, Level: %s\n", start, end, level)
	var filtered []logrus.Entry
	for _, log := range Log {
		// fmt.Println("Log Entry Time:", log.Time, "Level:", log.Level)
		if level != "" && log.Level.String() != level {
			continue
		}

		if start != "" && log.Time.Before(ParseDate(start)) {
			continue
		}

		if end != "" && log.Time.After(ParseDate(end)) {
			continue
		}

		filtered = append(filtered, log)
	}
	// fmt.Println("f: ", filtered)
	return filtered
}

// Helper function to parse date strings
func ParseDate(dateStr string) time.Time {
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
	}
	return parsed
}

// Parse log level safely
func ParseLogLevel(level string) logrus.Level {
	parsed, err := logrus.ParseLevel(level)
	if err != nil {
		return logrus.InfoLevel // Default to INFO
	}
	return parsed
}

func ParseLogLine(line string) (logrus.Entry, error) {
	var entry logrus.Entry

	if !strings.Contains(line, "level=") || !strings.Contains(line, "msg=") {
		fmt.Printf("Skipping line (invalid format): %s\n", line)
		return logrus.Entry{}, fmt.Errorf("invalid log line format: %s", line)
	}
	// Parse the log line manually
	var timeStr, level, msg, source string
	_, err := fmt.Sscanf(line, `time=%q level=%s msg=%q source=%s`, &timeStr, &level, &msg, &source)
	if err != nil {
		return logrus.Entry{}, fmt.Errorf("failed to parse log line: %v", err)
	}

	// Example: Parsing logic using a regex or JSON unmarshalling.
	// Parse time
	parsedTime, timeErr := time.Parse(time.RFC3339, timeStr)
	if timeErr != nil {
		return logrus.Entry{}, fmt.Errorf("failed to parse time: %v", timeErr)
	}

	// Parse log level
	parsedLevel, levelErr := logrus.ParseLevel(level)
	if levelErr != nil {
		return logrus.Entry{}, fmt.Errorf("failed to parse log level: %v", levelErr)
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
