package scaler_test

import (
	"fmt"
	"testing"
	"time"
	"watxhing-scaler-go/models"
	"watxhing-scaler-go/scaler"
	"watxhing-scaler-go/utils"
)

// Here is a fictitious scenario
// we have two schedules in the list, on is on weekdays and the other is on weekends
// on weekdays we have 6 replicas and on weekends we have 2 replicas

var testSchedules = []models.ScheduleConfig{
	{
		Name:      "business-hours",
		Days:      []int{1, 2, 3, 4, 5},
		StartTime: "05:00",
		EndTime:   "21:00",
		Replicas:  1,
	},
	{
		Name:      "mon-thru-thurs-off-hours",
		Days:      []int{1, 2, 3, 4},
		StartTime: "21:00",
		EndTime:   "05:00",
		Replicas:  2,
	},
	{
		Name:      "friday-off-hours",
		Days:      []int{5},
		StartTime: "21:00",
		EndTime:   "00:00",
		Replicas:  3,
	},
	{
		Name:      "saturday-hours",
		Days:      []int{6},
		StartTime: "00:00",
		EndTime:   "00:00",
		Replicas:  4,
	},
	{
		Name:      "sunday-hours",
		Days:      []int{0},
		StartTime: "00:00",
		EndTime:   "05:00",
		Replicas:  5,
	},
}

type TestClient struct {
	throw    error
	replicas int
}

func (tc *TestClient) ScaleDeployment(namespace string, deploymentName string, replicas int) error {
	if tc.throw != nil {
		return tc.throw
	}
	tc.replicas = replicas
	return nil
}

func Test_ApplyScheduledScalings_AllSchedulesApplied(t *testing.T) {
	mondayBusiness := time.Date(2021, 1, 4, 8, 0, 0, 0, utils.TIMEZONE)
	mondayOff := time.Date(2021, 1, 4, 22, 0, 0, 0, utils.TIMEZONE)
	fridayOff := time.Date(2021, 1, 8, 22, 0, 0, 0, utils.TIMEZONE)
	saturday := time.Date(2021, 1, 9, 12, 0, 0, 0, utils.TIMEZONE)
	sunday := time.Date(2021, 1, 10, 2, 0, 0, 0, utils.TIMEZONE)

	tc := &TestClient{
		replicas: 0,
	}

	s := scaler.NewScaler(tc, "test", "test")

	s.ApplyScheduledScalings(mondayBusiness, testSchedules)
	if tc.replicas != 1 {
		t.Errorf("Expected 1, got %d", tc.replicas)
	}

	s.ApplyScheduledScalings(mondayOff, testSchedules)
	if tc.replicas != 2 {
		t.Errorf("Expected 2, got %d", tc.replicas)
	}

	s.ApplyScheduledScalings(fridayOff, testSchedules)
	if tc.replicas != 3 {
		t.Errorf("Expected 3, got %d", tc.replicas)
	}

	s.ApplyScheduledScalings(saturday, testSchedules)
	if tc.replicas != 4 {
		t.Errorf("Expected 4, got %d", tc.replicas)
	}

	s.ApplyScheduledScalings(sunday, testSchedules)
	if tc.replicas != 5 {
		t.Errorf("Expected 5, got %d", tc.replicas)
	}
}

func Test_ApplyScheduledScalings_IfNoSchedulesMatchNothingIsDone(t *testing.T) {

	tc := &TestClient{
		replicas: 0,
	}

	monday := time.Date(2021, 1, 4, 8, 0, 0, 0, utils.TIMEZONE)
	s := scaler.NewScaler(tc, "test", "test")

	s.ApplyScheduledScalings(monday, []models.ScheduleConfig{})
	if tc.replicas != 0 {
		t.Errorf("Expected 0, got %d", tc.replicas)
	}
}

func Test_ApplyScheduledScalings_IfClientHasErrorThatErrorIsPassedOut(t *testing.T) {
	expectedError := fmt.Errorf("test error")
	tc := &TestClient{
		throw: expectedError,
	}

	monday := time.Date(2021, 1, 4, 8, 0, 0, 0, utils.TIMEZONE)
	s := scaler.NewScaler(tc, "test", "test")

	err := s.ApplyScheduledScalings(monday, testSchedules)
	if err != expectedError {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}
}
