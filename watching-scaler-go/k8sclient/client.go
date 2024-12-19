package k8sclient

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client interface {
	ScaleDeployment(namespace string, deploymentName string, replicas int) error
}

type ClientImpl struct {
	clientset kubernetes.Interface
}

func New(clientset kubernetes.Interface) Client {
	return &ClientImpl{clientset: clientset}
}

func (c *ClientImpl) ScaleDeployment(namespace string, deploymentName string, replicas int) error {
	deployments := c.clientset.AppsV1().Deployments(namespace)
	deploy, err := deployments.Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	i32Replicas := int32(replicas)
	if *deploy.Spec.Replicas == i32Replicas {
		fmt.Printf("Deployment %s already has %d replicas\n", deploy.Name, replicas)
		return nil
	}

	deploy.Spec.Replicas = &i32Replicas
	if deploy.Annotations == nil {
		deploy.Annotations = make(map[string]string)
	}
	deploy.Annotations["lastScaledAt"] = time.Now().Format(time.RFC3339)
	_, err = deployments.Update(context.Background(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Deployment %s scaled to %d replicas\n", deploy.Name, replicas)
	return nil
}
