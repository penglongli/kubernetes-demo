package main

import (
	"github.com/gin-gonic/gin"

	"github.com/penglongli/kubernetes-demo/k8s-terminal-go/handler"
)

func main() {
	r := gin.Default()
	handler.Router(r)
	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
