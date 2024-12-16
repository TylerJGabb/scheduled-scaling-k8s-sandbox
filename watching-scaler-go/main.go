package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
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
	if err != nil {
		panic(err)
	}

	// list all deployments in the current namespace
	deployments, err := clientset.AppsV1().Deployments(os.Getenv("NAMESPACE")).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range deployments.Items {
		fmt.Printf("Deployment %s - %d replicas\n", d.Name, *d.Spec.Replicas)
	}

}
