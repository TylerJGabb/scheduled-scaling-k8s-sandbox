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
	namespace string,
	deploymentName string,
) *Scaler {
	return &Scaler{
		client:         client,
		namespace:      namespace,
		deploymentName: deploymentName,
	}
}

func (s *Scaler) ApplyScheduledScalings(t time.Time, schedules []models.ScheduleConfig) error {
	for _, schedule := range schedules {
		if schedule.IsActive(t) {
			fmt.Printf("Applying schedule '%s' with target replicas '%d'\n", schedule.Name, schedule.Replicas)
			if err := s.client.ScaleDeployment(s.namespace, s.deploymentName, schedule.Replicas); err != nil {
				return err
			}
			return nil // stop after the first matched schedule
		}
	}
	return nil
}
