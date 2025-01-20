package utils

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	// Temporary file for testing
	tempLogFile := "logs.log"
	defer os.Remove(tempLogFile)

	// Call InitLogger
	InitLogger()

	// Check Logger instance is initialized
	assert.NotNil(t, Logger, "Logger should be initialized")

	// Ensure the logger writes to a file
	_, err := os.Stat(tempLogFile)
	assert.NoError(t, err, "Log file should be created")
}

func TestLogError(t *testing.T) {
	// Redirect output for testing
	var buf bytes.Buffer
	InitLogger()
	Logger.SetOutput(&buf)

	// Log an error
	testSource := "TestSource"
	testMessage := "This is a test error"
	testError := errors.New("test error")
	LogError(testSource, testMessage, testError)

	// Check output
	logOutput := buf.String()
	assert.Contains(t, logOutput, "error", "Log should contain 'error'")
	assert.Contains(t, logOutput, testSource, "Log should contain the source")
	assert.Contains(t, logOutput, testMessage, "Log should contain the message")
	assert.Contains(t, logOutput, testError.Error(), "Log should contain the error")
}

func TestLogWarn(t *testing.T) {
	// Redirect output for testing
	var buf bytes.Buffer
	InitLogger()
	Logger.SetOutput(&buf)

	// Log a warning
	testSource := "TestSource"
	testMessage := "This is a test warning"
	LogWarn(testSource, testMessage)

	// Check output
	logOutput := buf.String()
	assert.Contains(t, logOutput, "warning", "Log should contain 'warning'")
	assert.Contains(t, logOutput, testSource, "Log should contain the source")
	assert.Contains(t, logOutput, testMessage, "Log should contain the message")
}

func TestLogInfo(t *testing.T) {
	// Redirect output for testing
	var buf bytes.Buffer
	InitLogger()
	Logger.SetOutput(&buf)

	// Log an info message
	testSource := "TestSource"
	testMessage := "This is a test info"
	LogInfo(testSource, testMessage)

	// Check output
	logOutput := buf.String()
	assert.Contains(t, logOutput, "info", "Log should contain 'info'")
	assert.Contains(t, logOutput, testSource, "Log should contain the source")
	assert.Contains(t, logOutput, testMessage, "Log should contain the message")
}
