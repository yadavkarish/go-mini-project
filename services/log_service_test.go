package services

import (
	"bytes"
	// "csv-microservice/services"
	"csv-microservice/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Mock OsOpen function for testing
var OsOpen = func(name string) (*os.File, error) {
	return os.Open(name)
}

// Test LogEntry functionality
func TestLogEntry(t *testing.T) {
	buffer := new(bytes.Buffer)
	utils.Logger = logrus.New()
	utils.Logger.SetOutput(buffer)

	// Reset logs for the test using ResetLogs()
	ResetLogs()

	// Call the LogEntry function
	LogEntry("info", "testSource", "Test message")

	// Assertions
	assert.Len(t, GetLogsSlice(), 1) // Verify that a log entry was added
	assert.Contains(t, buffer.String(), "Test message")
	assert.Contains(t, buffer.String(), "testSource")
}

// TestGetLogs_Success tests GetLogs API success scenario.
func TestGetLogs_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/logs", GetLogs)

	// Mock log file
	logContent := `time="2023-01-01T00:00:00Z" level=info msg="Test log 1" source=testSource1
time="2023-01-02T00:00:00Z" level=warn msg="Test log 2" source=testSource2`
	tempFile, err := os.CreateTemp("", "logs.log")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	_, _ = tempFile.WriteString(logContent)
	tempFile.Close()

	// Set up test
	os.Rename(tempFile.Name(), "logs.log")
	defer os.Remove("logs.log")

	req, _ := http.NewRequest("GET", "/logs?start=2023-01-01&end=2023-01-02&level=info", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test log 1")
	assert.NotContains(t, w.Body.String(), "Test log 2")
}

// TestGetLogs_FileNotFound tests the scenario where the log file is missing.
func TestGetLogs_FileNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/logs", GetLogs)

	req, _ := http.NewRequest("GET", "/logs", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to read logs")
}

func TestFilterLogs(t *testing.T) {
	// Initialize the logger
	utils.InitLogger()

	// Clear any pre-existing logs
	ResetLogs()

	// Set up test logs with specific timestamps
	testLogs := []logrus.Entry{
		{
			Time:    time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC), // Within range
			Level:   logrus.InfoLevel,
			Message: "Test log 1",
		},
		{
			Time:    time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC), // Within range
			Level:   logrus.WarnLevel,
			Message: "Test log 2",
		},
		{
			Time:    time.Date(2022, 12, 31, 12, 0, 0, 0, time.UTC), // Outside range
			Level:   logrus.ErrorLevel,
			Message: "Test log 3",
		},
	}

	// Add logs to the global slice
	Log = append(Log, testLogs...)

	// Filter logs with specific criteria
	filteredLogs := FilterLogs("2023-01-01", "2023-12-31", "warning")
	// fmt.Println(filteredLogs)
	// Validate filtered results
	assert.Len(t, filteredLogs, 1)
	assert.Equal(t, "Test log 2", filteredLogs[0].Message)
}

// TestParseDate tests the ParseDate function with valid and invalid inputs.
func TestParseDate(t *testing.T) {
	validDate := ParseDate("2023-01-01")
	assert.Equal(t, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), validDate)

	invalidDate := ParseDate("invalid-date")
	assert.Equal(t, time.Time{}, invalidDate)
}

// TestParseLogLevel tests the ParseLogLevel function with various inputs.
func TestParseLogLevel(t *testing.T) {
	assert.Equal(t, logrus.InfoLevel, ParseLogLevel("info"))
	assert.Equal(t, logrus.WarnLevel, ParseLogLevel("warn"))
	assert.Equal(t, logrus.InfoLevel, ParseLogLevel("invalid"))
}

// TestParseLogLine_Valid tests parseLogLine with valid log lines.
func TestParseLogLine_Valid(t *testing.T) {
	logLine := `time="2023-01-01T00:00:00Z" level=info msg="Test log 1" source=testSource1`
	entry, err := ParseLogLine(logLine)
	assert.NoError(t, err)
	assert.Equal(t, "Test log 1", entry.Message)
	assert.Equal(t, logrus.InfoLevel, entry.Level)
	assert.Equal(t, "testSource1", entry.Data["source"])
}

// TestParseLogLine_Invalid tests parseLogLine with invalid log lines.
func TestParseLogLine_Invalid(t *testing.T) {
	logLine := `invalid log line format`
	_, err := ParseLogLine(logLine)
	assert.Error(t, err)
}

// TestGetLogs_ScannerError tests GetLogs for scanner errors.
// func TestGetLogs_ScannerError(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	ResetLogs()
// 	router := gin.Default()
// 	router.GET("/logs", GetLogs)

// 	// Mock a file with invalid lines to trigger scanner errors
// 	logContent := `time="invalid" level=info msg="Test log" source=source`
// 	tempFile, err := os.CreateTemp("", "logs.log")
// 	assert.NoError(t, err)
// 	defer os.Remove(tempFile.Name())

// 	_, _ = tempFile.WriteString(logContent)
// 	tempFile.Close()

// 	// Rename to match expected log file
// 	os.Rename(tempFile.Name(), "logs.log")
// 	defer os.Remove("logs.log")

// 	// Make a GET request to fetch logs
// 	req, _ := http.NewRequest("GET", "/logs", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// Validate response
// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
// 	assert.Contains(t, w.Body.String(), "Failed to parse logs")
// }

func TestGetLogs_ScannerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/logs", GetLogs)

	// Reset state
	ResetLogs()

	// Write invalid content directly to logs.log
	logFile, err := os.Create("logs.log")
	assert.NoError(t, err)
	defer os.Remove("logs.log") // Cleanup after test

	logContent := `time="invalid_time" level=info msg="Test log" source=source`
	_, writeErr := logFile.WriteString(logContent)
	assert.NoError(t, writeErr)
	logFile.Close()

	// Make a GET request to fetch logs
	req, _ := http.NewRequest("GET", "/logs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validate response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to parse logs")

	// Reset state
	ResetLogs()
}
