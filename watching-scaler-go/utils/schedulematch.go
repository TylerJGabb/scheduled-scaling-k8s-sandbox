package utils

import (
	"fmt"
	"time"
	"watxhing-scaler-go/models"
)

var NYC, _ = time.LoadLocation("America/New_York")

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
	if arrayContains(schedule.Days, weekday) {
		start, _ := time.Parse("15:04", schedule.StartTime)
		start = time.Date(
			t.Year(),
			t.Month(),
			t.Day(),
			start.Hour(),
			start.Minute(),
			t.Second(),
			t.Nanosecond(),
			NYC,
		)
		end := start.Add(time.Minute * time.Duration(schedule.DurationMinutes))
		fmt.Printf("Start:   %v\n", start)
		fmt.Printf("End:     %v\n", end)
		fmt.Printf("Current: %v\n", t)

		if t.After(start) && t.Before(end) {
			return true
		}
	}
	if arrayContains(schedule.Days, (weekday-1)%7) {
		start, _ := time.Parse("15:04", schedule.StartTime)
		start = time.Date(
			t.Year(),
			t.Month(),
			t.Day(),
			start.Hour(),
			start.Minute(),
			t.Second(),
			t.Nanosecond(),
			NYC,
		)
		start = start.Add(-time.Hour * 24)
		end := start.Add(time.Minute * time.Duration(schedule.DurationMinutes))
		fmt.Printf("Start:   %v\n", start)
		fmt.Printf("End:     %v\n", end)
		fmt.Printf("Current: %v\n", t)

		if t.After(start) && t.Before(end) {
			return true
		}
	}
	return false
}
