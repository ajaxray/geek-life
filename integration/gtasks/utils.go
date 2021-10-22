package gtasks

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return tok
}

func gTaskToUnixTime(rfc3339Time string) int64 {
	goTime, _ := time.Parse(time.RFC3339, rfc3339Time)
	return goTime.Unix()
}

func unixToGtaskTime(unixTime int64) string {
	return time.Unix(unixTime, 0).Format(time.RFC3339)
}
