package util

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/asdine/storm/v3"
)

// ConnectStorm Create database connection
func ConnectStorm() *storm.DB {
	db, err := storm.Open(GetEnvStr("DB_FILE", "geek-life.db"))
	FatalIfError(err, "Could not connect Embedded Database File")

	return db
}

// UnixToTime create time.Time from string timestamp
func UnixToTime(timestamp string) time.Time {
	parts := strings.Split(timestamp, ".")
	i, err := strconv.ParseInt(parts[0], 10, 64)
	if LogIfError(err, "Could not parse timestamp : "+timestamp+" (using current time instead)") {
		return time.Unix(i, 0)
	}

	return time.Now()
}

// LogIfError logs the error and returns true on Error. think as IfError
func LogIfError(err error, msgOrPattern string, args ...interface{}) bool {
	if err != nil {
		message := fmt.Sprintf(msgOrPattern, args...)
		log.Printf("%s: %w\n", message, err)

		return true
	}

	return false
}

// FatalIfError logs the error and Exit program on Error
func FatalIfError(err error, msgOrPattern string, args ...interface{}) {
	message := fmt.Sprintf(msgOrPattern, args...)

	if LogIfError(err, message) {
		log.Fatal("FATAL ERROR: Exiting program! - ", message, "\n")
	}
}
