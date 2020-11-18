package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/penglongli/kubernetes-demo/kubectl-golang/kubectl"
)

var (
	applyYAML = `
apiVersion: v1
kind: Namespace
metadata:
  name: test-008
  labels:
    app: test-008
spec:
  finalizers:
  - kubernetes

---
apiVersion: v1
kind: Namespace
metadata:
  name: test-apply1
  labels:
    app: test-apply1
spec:
  finalizers:
  - kubernetes
---
apiVersion: v1
kind: Namespace
metadata:
  name: test-apply2
  labels:
    app: test-apply2
spec:
  finalizers:
  - kubernetes
`
)

func main() {
	// APPLY
	result, err := kubectl.Apply(context.Background(), []byte(applyYAML))
	// DELETE
	// result, err := kubectl.Delete(context.Background(), []byte(applyYAML))
	fmt.Println(err)
	fmt.Println(strings.Join(result, "\n"))
}
