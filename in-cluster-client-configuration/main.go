package main

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// list namespaces
	nsList, err := clientset.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		panic(err)
	}
	for _, ns := range nsList.Items {
		fmt.Println(ns.Name)
	}
}
