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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ScheduleConfig struct {
	Name      string `json:"name"`
	StartHour int    `json:"startHour"`
	EndHour   int    `json:"endHour"`
	Replicas  int    `json:"replicas"`
	Days      []int  `json:"days"`
}

type SchedulesConfig struct {
	Schedules []ScheduleConfig `json:"schedules"`
}

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

func arrayContains(arr []int, elem int) bool {
	for _, item := range arr {
		if item == elem {
			return true
		}
	}
	return false
}

func main() {

	clientset, err := getClientset()
	if err != nil {
		panic(err)
	}

	scheduleConfig := SchedulesConfig{}
	err = json.Unmarshal([]byte(os.Getenv("SCHEDULES")), &scheduleConfig)
	if err != nil {
		panic(err)
	}

	currentTime := time.Now()
	currentHour := currentTime.Hour()
	currentDay := int(currentTime.Weekday())

	for _, schedule := range scheduleConfig.Schedules {
		fmt.Printf("Processing schedule: %s\n", schedule.Name)
		if arrayContains(schedule.Days, currentDay) && currentHour >= schedule.StartHour && currentHour < schedule.EndHour {
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
			fmt.Printf("Scaled %s successfully to %d replicas", deploy.Name, schedule.Replicas)
			break
		}
	}
}
