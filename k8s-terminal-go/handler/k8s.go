package handler

import (
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/penglongli/kubernetes-demo/k8s-terminal-go/k8s"
	"github.com/penglongli/kubernetes-demo/k8s-terminal-go/result"
)

type Pod struct {
	Name       string   `json:"name"`
	Namespace  string   `json:"namespace"`
	Containers []string `json:"containers"`
}

func GetNamespaces(ctx *gin.Context) {
	clientSet, err := k8s.GetClientSet()
	if err != nil {
		result.Failed(ctx, result.ERROR, err.Error())
		return
	}

	list, err := clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		result.Failed(ctx, result.ERROR, err.Error())
		return
	}
	var namespaces []string
	for i := range list.Items {
		namespaces = append(namespaces, list.Items[i].Name)
	}
	result.Success(ctx, namespaces)
}

func GetPods(ctx *gin.Context) {
	namespace := ctx.Param("namespace")
	if namespace == "" {
		result.Failed(ctx, result.ERROR, "namespace cannot be null")
		return
	}

	clientSet, err := k8s.GetClientSet()
	if err != nil {
		result.Failed(ctx, result.ERROR, err.Error())
		return
	}

	list, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		result.Failed(ctx, result.ERROR, err.Error())
		return
	}
	var pods []*Pod
	for i := range list.Items {
		pod := &Pod{
			Name:      list.Items[i].Name,
			Namespace: list.Items[i].Namespace,
		}

		containers := list.Items[i].Spec.Containers
		for j := range containers {
			pod.Containers = append(pod.Containers, containers[j].Name)
		}
		pods = append(pods, pod)
	}
	result.Success(ctx, pods)
}
