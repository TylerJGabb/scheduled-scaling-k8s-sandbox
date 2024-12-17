package utils

import (
	"os"
	"time"
	"watxhing-scaler-go/models"
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

func arrayContains(arr []int, elem int) bool {
	for _, item := range arr {
		if item == elem {
			return true
		}
	}
	return false
}

func IsTimeInSchedule(
	t time.Time,
	schedule models.ScheduleConfig,
) bool {
	weekday := int(t.Weekday())
	start, _ := time.Parse("15:04", schedule.StartTime)
	start = time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		start.Hour(),
		start.Minute(),
		t.Second(),
		t.Nanosecond(),
		TIMEZONE,
	)
	end := start.Add(time.Minute * time.Duration(schedule.DurationMinutes))
	if arrayContains(schedule.Days, weekday) && t.After(start) && t.Before(end) {
		return true
	}

	previousWeekday := (weekday - 1) % 7
	start = start.Add(-time.Hour * 24)
	end = start.Add(time.Minute * time.Duration(schedule.DurationMinutes))
	return arrayContains(schedule.Days, previousWeekday) && t.After(start) && t.Before(end)
}
