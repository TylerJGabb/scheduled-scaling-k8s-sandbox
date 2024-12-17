package models

import "encoding/json"

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

func (s *ScheduleConfig) FromJson(data []byte) error {
	return json.Unmarshal(data, s)
}
