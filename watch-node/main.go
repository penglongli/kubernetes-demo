package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	yearSeconds int64 = 365 * 24 * 60 * 60
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

	nodeCfg := NewNodeWatchConfig(clientset)
	go nodeCfg.Watch()

	ch := make(chan struct{})
	fmt.Println(<-ch)
}

func NewNodeWatchConfig(clientSet *kubernetes.Clientset) *NodeWatchConfig {
	return &NodeWatchConfig{
		stopChan: make(chan struct{}),
		clientSet: clientSet,
	}
}

type NodeWatchConfig struct {
	stopChan chan struct{}
	clientSet *kubernetes.Clientset
}

func (config *NodeWatchConfig) Watch()  {
	nodeWatch, err := getNodeWatchChannel(config.clientSet)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	fmt.Println("[Info] watch nodes goroutine started.")
	nodeChan := nodeWatch.ResultChan()

L:
	for {
		select {
		case event, ok := <-nodeChan: {
			if !ok {
				bs, _ := json.Marshal(event)
				fmt.Printf("[Warn] Watch node is closed, event: %s\n", string(bs))

				fmt.Println("[Warn] Start reconnect k8s.")
				watchChan, err := getNodeWatchChannel(config.clientSet)
				if err != nil {
					fmt.Printf("[Error] Reconnect failed, err: %s\n", err.Error())
					break
				}

				nodeWatch = watchChan
				nodeChan = watchChan.ResultChan()
				continue
			}

			if event.Type == watch.Error {
				bs, _ := json.Marshal(event)
				fmt.Printf("[Error] received Error event, event: %s\n", string(bs))
				continue
			}

			if event.Type == watch.Added || event.Type == watch.Deleted {
				bs, err := json.Marshal(event.Object)
				if err != nil {
					fmt.Printf("[Error] marshal failed, event: %s, err: %s, \n", string(bs), err.Error())
					continue
				}

				node := new(v1.Node)
				if err = json.Unmarshal(bs, node); err != nil {
					fmt.Printf("[Error] unmarshal failed, event: %s, err: %s\n", string(bs), err.Error())
					continue
				}

				// TODO
				// 此处拿到了 Node 节点，根据   增加、删除   做授权、接触授权操作

			}
		}
		case <-config.stopChan:
			fmt.Println("[Warn] node watch closed.")
			break L
		}
	}

	// 正常关闭
	if nodeWatch != nil {
		 nodeWatch.Stop()
	}
}

func (config *NodeWatchConfig) Stop()  {
	config.stopChan <- struct{}{}
}

func getNodeWatchChannel(clientSet *kubernetes.Clientset) (nodeWatch watch.Interface, err error) {
	nodeWatch, err = clientSet.CoreV1().Nodes().Watch(context.Background(), metav1.ListOptions{
		Watch: true,
		ResourceVersion: "0",
		TimeoutSeconds: &yearSeconds,
	})
	if err == nil {
		return nodeWatch, nil
	}

	// 尝试重连
	fmt.Printf("[Error] connect k8s error: %s\n", err.Error())
	for i := 0; i < 10; i++ {
		nodeWatch, err = clientSet.CoreV1().Nodes().Watch(context.Background(), metav1.ListOptions{
			Watch: true,
			ResourceVersion: "0",
			TimeoutSeconds: &yearSeconds,
		})
		if err != nil {
			fmt.Printf("[Error] connect k8s error: %s\n", err.Error())
			time.Sleep(5 * time.Duration(time.Second))
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return nodeWatch, nil
}
