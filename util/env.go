package util

import (
	"log"
	"os"
	"strconv"
	// "github.com/subosito/gotenv"
)

// func init() {
// 	gotenv.Load()
// }

// GetEnvInt finds an ENV variable and converts to int, otherwise return default value
func GetEnvInt(key string, defaultVal int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		} else {
			log.Fatal(err)
		}

		return 0
	}

	return defaultVal
}

// GetEnvStr finds an ENV variable, otherwise return default value
func GetEnvStr(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defaultVal
}
