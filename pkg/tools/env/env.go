package env

import "os"

// GetEnvOrFallback it will get value of environment variable but if is not initialized will
// return fallback value (default value)
func GetEnvOrFallback(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
