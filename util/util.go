package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/mitchellh/go-homedir"
)

// ConnectStorm Create database connection
func ConnectStorm(dbFilePath string) *storm.DB {
	var dbPath string

	if dbFilePath != "" {
		info, err := os.Stat(dbFilePath)
		if err == nil && info.IsDir() {
			fmt.Println("Mentioned DB path is a directory. Please specify a file or ignore to create automatically in home directory.")
			os.Exit(1)
		}

		dbPath = dbFilePath
	} else {
		dbPath = GetEnvStr("DB_FILE", "")
	}

	var err error
	if dbPath == "" {
		// Try in home dir
		dbPath, err = homedir.Expand("~/.geek-life/default.db")

		// If home dir is not detected, try in system tmp dir
		if err != nil {
			f, _ := ioutil.TempFile("geek-life", "default.db")
			dbPath = f.Name()
		}
	}

	CreateDirIfNotExist(path.Dir(dbPath))

	db, openErr := storm.Open(dbPath)
	FatalIfError(openErr, "Could not connect Embedded Database File")

	return db
}

// CreateDirIfNotExist creates a directory if not found
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
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
