package utils

import (
	"os"
	"time"
)

func getTimezone() *time.Location {
	tz := os.Getenv("IANA_TIMEZONE")
	if tz == "" {
		tz = "America/New_York"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}
	return loc
}

var TIMEZONE = getTimezone()

func ArrayContains(arr []int, elem int) bool {
	for _, item := range arr {
		if item == elem {
			return true
		}
	}
	return false
}
