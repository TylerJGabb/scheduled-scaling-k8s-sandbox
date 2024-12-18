package k8sclient_test

import (
	"context"
	"testing"
	"time"
	"watxhing-scaler-go/k8sclient"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/fake"
)

func Test_ScaleDeployment_NeedsScaling(t *testing.T) {
	var replicas int32 = 1

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
	}

	clientset := fake.NewSimpleClientset(deploy)

	client := k8sclient.New(clientset)

	err := client.ScaleDeployment("test", "test", 2)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	deploy, err = clientset.AppsV1().Deployments("test").Get(context.Background(), deploy.Name, metav1.GetOptions{})
	if err != nil {
		t.Errorf("Error getting deployment: %v", err)
	}

	if *deploy.Spec.Replicas != 2 {
		t.Errorf("Expected 2, got %d", *deploy.Spec.Replicas)
	}

	scaledAt, err := time.Parse(time.RFC3339, deploy.Annotations["lastScaledAt"])
	if err != nil {
		t.Errorf("Error parsing time: %v", err)
	}

	if time.Now().Sub(scaledAt) > 5*time.Second {
		t.Errorf("Expected scaledAt to be within the last 5 seconds, got %v", scaledAt)
	}

}

func Test_ScaleDeployment_AlreadyScaled(t *testing.T) {
	var replicas int32 = 2

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
	}

	clientset := fake.NewSimpleClientset(deploy)

	client := k8sclient.New(clientset)

	err := client.ScaleDeployment("test", "test", 2)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	deploy, err = clientset.AppsV1().Deployments("test").Get(context.Background(), deploy.Name, metav1.GetOptions{})
	if err != nil {
		t.Errorf("Error getting deployment: %v", err)
	}

	if *deploy.Spec.Replicas != 2 {
		t.Errorf("Expected 2, got %d", *deploy.Spec.Replicas)
	}

	if deploy.Annotations["lastScaledAt"] != "" {
		t.Error("Expected lastScaledAt annotation to be empty")
	}

}
