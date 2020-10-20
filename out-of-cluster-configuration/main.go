package main

import (
	"context"
	"encoding/base64"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func main() {
	restConfig, err := buildRestConfig("Please enter kubeconfig string of base64")
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
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

func buildRestConfig(base64KubeConfig string) (resetConfig *rest.Config, err error) {
	kubeConfig, err := base64.StdEncoding.DecodeString(base64KubeConfig)
	if err != nil {
		return nil, err
	}

	conf, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (config *clientcmdapi.Config, e error) {
		return clientcmd.Load(kubeConfig)
	})

	if err != nil {
		return nil, err
	}
	return conf, nil
}

