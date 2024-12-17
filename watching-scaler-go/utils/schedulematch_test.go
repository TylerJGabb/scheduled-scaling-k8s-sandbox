package utils_test

import (
	"testing"
	"time"
	"watxhing-scaler-go/models"
	"watxhing-scaler-go/utils"
)

func TestIsTimeInSchedule_Happy(t *testing.T) {
	t1 := time.Date(2021, 1, 1, 12, 30, 0, 0, utils.NYC)
	schedule := models.ScheduleConfig{
		Name:            "test",
		StartTime:       "12:00",
		DurationMinutes: 60,
		Replicas:        1,
		Days:            []int{int(t1.Weekday())},
	}
	if !utils.IsTimeInSchedule(t1, schedule) {
		t.Error("Expected true, got false")
	}
	t2, _ := time.Parse("15:04", "11:30")
	if utils.IsTimeInSchedule(t2, schedule) {
		t.Error("Expected false, got true")
	}
	t3, _ := time.Parse("15:04", "13:00")
	if utils.IsTimeInSchedule(t3, schedule) {
		t.Error("Expected false, got true")
	}
	t4, _ := time.Parse("15:04", "13:01")
	if utils.IsTimeInSchedule(t4, schedule) {
		t.Error("Expected false, got true")
	}
}

func TestIsTimeInSchedule_SpansMidnight(t *testing.T) {
	t1 := time.Date(2021, 1, 1, 1, 30, 0, 0, utils.NYC)
	schedule := models.ScheduleConfig{
		Name:            "test",
		StartTime:       "23:00",
		DurationMinutes: 240,
		Replicas:        1,
		Days:            []int{(int(t1.Weekday()) - 1) % 7},
	}
	if !utils.IsTimeInSchedule(t1, schedule) {
		t.Error("Expected true, got false")
	}
}
