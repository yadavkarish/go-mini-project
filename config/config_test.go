package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDBConnectionString(t *testing.T) {
	// Test Case 1: DB_CONNECTION_STRING is set with a valid value
	t.Run("Valid DB_CONNECTION_STRING", func(t *testing.T) {
		// Set the environment variable
		os.Setenv("DB_CONNECTION_STRING", "valid_connection_string")

		// Call the function
		connStr := GetDBConnectionString()

		// Assert the value is correct
		assert.Equal(t, "valid_connection_string", connStr)

		// Clean up
		os.Unsetenv("DB_CONNECTION_STRING")
	})

	// Test Case 2: DB_CONNECTION_STRING is not set
	t.Run("DB_CONNECTION_STRING not set", func(t *testing.T) {
		// Unset the environment variable to simulate it not being set
		os.Unsetenv("DB_CONNECTION_STRING")

		// Call the function
		connStr := GetDBConnectionString()

		// Assert it returns an empty string
		assert.Equal(t, "", connStr)
	})

	// Test Case 3: DB_CONNECTION_STRING is set with an empty value
	t.Run("DB_CONNECTION_STRING is empty", func(t *testing.T) {
		// Set the environment variable to an empty string
		os.Setenv("DB_CONNECTION_STRING", "")

		// Call the function
		connStr := GetDBConnectionString()

		// Assert the value is an empty string
		assert.Equal(t, "", connStr)

		// Clean up
		os.Unsetenv("DB_CONNECTION_STRING")
	})
}
