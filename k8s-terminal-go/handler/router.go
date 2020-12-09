package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"

	"github.com/penglongli/kubernetes-demo/k8s-terminal-go/result"
)

type terminalRequest struct {
	Namespace string `json:"namespace" binding:"required"`
	Pod       string `json:"pod" binding:"required"`
	Container string `json:"container" binding:"required"`
}

func Router(r *gin.Engine) {
	// Kubernetes resource
	r.GET("/k8s/namespaces", GetNamespaces)
	r.GET("/k8s/namespaces/:namespace/pods", GetPods)

	// Terminal
	r.GET("/terminal/exec/*path", execHandler)
	r.GET("/terminal/logging/*path", loggingHandler)

	// Page and static files
	r.StaticFS("/page/terminal", http.Dir("k8s-terminal-go/webapp/html/"))
	r.StaticFS("/static", http.Dir("k8s-terminal-go/webapp/static/"))
}

func execHandler(ctx *gin.Context) {
	sockHandler := sockjs.NewHandler("/terminal/exec", sockjs.DefaultOptions, func(session sockjs.Session) {
		go func(c *gin.Context) {
			defer func() {
				if recoverErr := recover(); recoverErr != nil {
					klog.Error(recoverErr)
				}
			}()

			request := &terminalRequest{
				Namespace: ctx.Query("namespace"),
				Pod:       ctx.Query("pod"),
				Container: ctx.Query("container"),
			}
			klog.Infof("Exec received request: %#v", request)
			// TODO CheckNamespace CheckPod

			terminalSession := &TerminalSession{
				SockSession:    session,
				SizeChan:       make(chan *remotecommand.TerminalSize, 1),
				Namespace:      request.Namespace,
				Pod:            request.Pod,
				Container:      request.Container,
				Timeout:        10 * time.Minute,
				RefreshTimeout: make(chan struct{}, 1),
			}
			// Get available shell in pod
			shell, err := terminalSession.CheckShellInPod()
			if err != nil {
				result.Failed(ctx, result.ERROR, err.Error())
				return
			}
			// Exec
			if err = terminalSession.Exec(shell); err != nil {
				_ = terminalSession.SockSession.Close(126, err.Error())
			}
			return
		}(ctx)
	})
	sockHandler.ServeHTTP(ctx.Writer, ctx.Request)
}

func loggingHandler(ctx *gin.Context) {

}
