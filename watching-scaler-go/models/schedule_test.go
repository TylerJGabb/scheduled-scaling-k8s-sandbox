package models_test

import (
	"testing"
	"time"
	"watxhing-scaler-go/models"
	"watxhing-scaler-go/utils"
)

func Test_IsActive_Happy(t *testing.T) {
	t1 := time.Date(2021, 1, 1, 12, 30, 0, 0, utils.TIMEZONE)
	schedule := models.ScheduleConfig{
		Name:            "test",
		StartTime:       "12:00",
		DurationMinutes: 60,
		Replicas:        1,
		Days:            []int{int(t1.Weekday())},
	}
	if !schedule.IsActive(t1) {
		t.Error("Expected true, got false")
	}
	t2, _ := time.Parse("15:04", "11:30")
	if schedule.IsActive(t2) {
		t.Error("Expected false, got true")
	}
	t3, _ := time.Parse("15:04", "13:00")
	if schedule.IsActive(t3) {
		t.Error("Expected false, got true")
	}
	t4, _ := time.Parse("15:04", "13:01")
	if schedule.IsActive(t4) {
		t.Error("Expected false, got true")
	}
}

func Test_IsActive_SpansMidnight(t *testing.T) {
	t1 := time.Date(2021, 1, 1, 1, 30, 0, 0, utils.TIMEZONE)
	schedule := models.ScheduleConfig{
		Name:            "test",
		StartTime:       "23:00",
		DurationMinutes: 240,
		Replicas:        1,
		Days:            []int{(int(t1.Weekday()) - 1) % 7},
	}
	if !schedule.IsActive(t1) {
		t.Error("Expected true, got false")
	}
	t2 := t1.Add(-time.Hour * 2)
	if !schedule.IsActive(t2) {
		t.Error("Expected true, got false")
	}
	t3 := t1.Add(-time.Hour * 4)
	if schedule.IsActive(t3) {
		t.Error("Expected false, got true")
	}
	t4 := t1.Add(time.Hour * 3)
	if schedule.IsActive(t4) {
		t.Error("Expected false, got true")
	}
}
