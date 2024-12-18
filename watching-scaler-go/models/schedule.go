package models

import (
	"encoding/json"
	"fmt"
	"time"
	"watxhing-scaler-go/utils"
)

type ScheduleConfig struct {
	Name            string `json:"name"`
	StartTime       string `json:"startTime"`
	DurationMinutes int    `json:"durationMinutes"`
	Replicas        int    `json:"replicas"`
	Days            []int  `json:"days"`
}

func (schedule ScheduleConfig) IsActive(t time.Time) bool {
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
		utils.TIMEZONE,
	)
	end := start.Add(time.Minute * time.Duration(schedule.DurationMinutes))
	if utils.ArrayContains(schedule.Days, weekday) && t.After(start) && t.Before(end) {
		return true
	}

	previousWeekday := (weekday - 1) % 7
	start = start.Add(-time.Hour * 24)
	end = start.Add(time.Minute * time.Duration(schedule.DurationMinutes))
	return utils.ArrayContains(schedule.Days, previousWeekday) && t.After(start) && t.Before(end)
}

type SchedulesConfig struct {
	Schedules []ScheduleConfig `json:"schedules"`
}

func (s *SchedulesConfig) Validate() error {
	if len(s.Schedules) == 0 {
		return fmt.Errorf("`schedules` must be present, and contain at least one schedule")
	}
	for _, schedule := range s.Schedules {
		if err := schedule.Validate(); err != nil {
			return err
		}
	}
	return nil
}
func (s *ScheduleConfig) FromJson(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *ScheduleConfig) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("`name` is missing from a schedule")
	}
	if s.StartTime == "" {
		return fmt.Errorf("`startTime` is missing from schedule %s", s.Name)
	}
	if s.DurationMinutes <= 0 {
		return fmt.Errorf("`durationMinutes` must be present and greater than 0 for schedule %s", s.Name)
	}
	if s.Replicas < 0 || s.Replicas > 10 {
		return fmt.Errorf("`replicas` must be present and within [0, 10] for schedule %s", s.Name)
	}
	if len(s.Days) == 0 {
		return fmt.Errorf("`days` must be present for schedule %s", s.Name)
	}
	for _, day := range s.Days {
		if day < 0 || day > 6 {
			return fmt.Errorf("each day in `days` must be within [0, 6] for schedule %s", s.Name)
		}
	}
	return nil
}
