package models

import (
	"encoding/json"
	"fmt"
)

type ScheduleConfig struct {
	Name            string `json:"name"`
	StartTime       string `json:"startTime"`
	DurationMinutes int    `json:"durationMinutes"`
	Replicas        int    `json:"replicas"`
	Days            []int  `json:"days"`
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
	if s.Replicas <= 0 || s.Replicas > 10 {
		return fmt.Errorf("`replicas` must be present and between 1 and 10 for schedule %s", s.Name)
	}
	if len(s.Days) == 0 {
		return fmt.Errorf("`days` must be present for schedule %s", s.Name)
	}
	for _, day := range s.Days {
		if day < 0 || day > 6 {
			return fmt.Errorf("each day in `days` must be between 0 and 6 for schedule %s", s.Name)
		}
	}
	return nil
}
