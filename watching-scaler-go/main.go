package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"watxhing-scaler-go/k8sclient"
	"watxhing-scaler-go/models"
	"watxhing-scaler-go/scaler"

	_ "time/tzdata"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	ns := os.Getenv("NAMESPACE")
	if ns == "" {
		panic("NAMESPACE environment variable is required")
	}
	deployName := os.Getenv("DEPLOYMENT")
	if deployName == "" {
		panic("DEPLOYMENT environment variable is required")
	}
	schedules := os.Getenv("SCHEDULES")
	if schedules == "" {
		panic("SCHEDULES environment variable is required")
	}
	scheduleConfig := models.SchedulesConfig{}
	if err := json.Unmarshal([]byte(schedules), &scheduleConfig); err != nil {
		panic(fmt.Errorf("Error parsing schedules: %s", err))
	}

	if err := scheduleConfig.Validate(); err != nil {
		panic(fmt.Errorf("Error validating schedules: %s", err))
	}

	fmt.Printf("Found %d schedules\n", len(scheduleConfig.Schedules))

	clientset, err := getClientset()
	if err != nil {
		panic(fmt.Errorf("Error creating k8s clientset: %s", err))
	}

	client := k8sclient.New(clientset)
	s := scaler.NewScaler(client, scheduleConfig.Schedules, ns, deployName)

	for {
		s.ApplyScheduledScalings(scheduleConfig.Schedules)
		time.Sleep(30 * time.Second)
	}
}
