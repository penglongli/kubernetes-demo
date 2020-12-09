package main

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/penglongli/kubernetes-demo/k8s-leader-election/utils"
)

func main()  {
	// configMap()
	// endpoint()
	lease()
}

func execute(ctx context.Context) {
	fmt.Println("I'm executing...")
	time.Sleep(1000 * time.Second)
}

func configMap() {
	clientSet, err := getClientSet()
	if err != nil {
		panic(err)
	}

	leaderElection := utils.NewLeaderElectionWithConfigMap(clientSet, "i-am-a-cm-lock", execute)
	if err := leaderElection.Run(); err != nil {
		panic(err)
	}
}

func endpoint() {
	clientSet, err := getClientSet()
	if err != nil {
		panic(err)
	}

	leaderElection := utils.NewLeaderElectionWithEndpoint(clientSet, "i-am-a-cm-lock", execute)
	if err := leaderElection.Run(); err != nil {
		panic(err)
	}
}

func lease() {
	clientSet, err := getClientSet()
	if err != nil {
		panic(err)
	}

	leaderElection := utils.NewLeaderElectionWithLease(clientSet, "i-am-a-cm-lock", execute)
	if err := leaderElection.Run(); err != nil {
		panic(err)
	}
}

func getClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}
