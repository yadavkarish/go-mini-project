package config

import "os"

func GetDBConnectionString() string {
	return os.Getenv("DB_CONNECTION_STRING")
}
