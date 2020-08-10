package utils

import "os"

// Get value of an env variable and return a default if it doesn't
// exist
func GetEnv(key, defaultVal string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultVal
	}
	return value
}
