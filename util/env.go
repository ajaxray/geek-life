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
	var err error
	intVal := defaultVal

	if v, ok := os.LookupEnv(key); ok {
		if intVal, err = strconv.Atoi(v); err != nil {
			log.Fatal(err)
		}
	}

	return intVal
}

// GetEnvStr finds an ENV variable, otherwise return default value
func GetEnvStr(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defaultVal
}
