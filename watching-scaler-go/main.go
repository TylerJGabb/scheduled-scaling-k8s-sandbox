package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"watxhing-scaler-go/models"
	"watxhing-scaler-go/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "time/tzdata"
)

func getClientset() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if os.Getenv("LOCAL") == "true" {
		fmt.Println("Using out of cluster config")
		kubeConfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Using in cluster config")
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}

func main() {

	clientset, err := getClientset()
	if err != nil {
		panic(err)
	}

	scheduleConfig := models.SchedulesConfig{}
	err = json.Unmarshal([]byte(os.Getenv("SCHEDULES")), &scheduleConfig)
	if err != nil {
		panic(fmt.Errorf("Error parsing schedules: %s", err))
	}
	if err := scheduleConfig.Validate(); err != nil {
		panic(fmt.Errorf("Error validating schedules: %s", err))
	}

	fmt.Printf("Found %d schedules\n", len(scheduleConfig.Schedules))

	for {
		current := time.Now()

		fmt.Println("-----------------------------------")
		fmt.Printf("Current time: %v\n", current.In(utils.TIMEZONE))

		for _, schedule := range scheduleConfig.Schedules {
			fmt.Printf("Processing schedule: %s\n", schedule.Name)
			// start and end are in the format HH:MM

			if utils.IsTimeInSchedule(current, schedule) {
				fmt.Printf("Matched schedule %s, scaling to %d replicas\n", schedule.Name, schedule.Replicas)
				deployments := clientset.AppsV1().Deployments(os.Getenv("NAMESPACE"))
				deploy, err := deployments.Get(context.Background(), os.Getenv("DEPLOYMENT"), metav1.GetOptions{})
				if err != nil {
					panic(err)
				}
				i32Replicas := int32(schedule.Replicas)
				if *deploy.Spec.Replicas == i32Replicas {
					fmt.Printf("Deployment %s already has %d replicas\n", deploy.Name, schedule.Replicas)
					break
				}

				deploy.Spec.Replicas = &i32Replicas
				_, err = deployments.Update(context.Background(), deploy, metav1.UpdateOptions{})
				if err != nil {
					panic(err)
				}
				fmt.Printf("Scaled %s successfully to %d replicas\n", deploy.Name, schedule.Replicas)
				break
			}
		}
		time.Sleep(30 * time.Second)
	}
}
