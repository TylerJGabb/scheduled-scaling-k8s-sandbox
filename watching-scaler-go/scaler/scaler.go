package scaler

import (
	"fmt"
	"time"
	"watxhing-scaler-go/k8sclient"
	"watxhing-scaler-go/models"
)

type Scaler struct {
	client         k8sclient.Client
	namespace      string
	deploymentName string
}

func NewScaler(
	client k8sclient.Client,
	schedules []models.ScheduleConfig,
	namespace string,
	deploymentName string,
) *Scaler {
	return &Scaler{
		client:         client,
		namespace:      namespace,
		deploymentName: deploymentName,
	}
}

func (s *Scaler) ApplyScheduledScalings(
	t time.Time,
	schedules []models.ScheduleConfig) {
	for _, schedule := range schedules {
		if schedule.IsActive(t) {
			err := s.client.ScaleDeployment(
				s.namespace,
				s.deploymentName,
				schedule.Replicas,
			)
			if err != nil {
				fmt.Printf("Error scaling deployment: %v\n", err)
			}
			return // stop after the first matched schedule
		}
	}
}
