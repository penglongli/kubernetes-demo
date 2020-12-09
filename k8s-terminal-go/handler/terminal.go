package handler

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"

	"github.com/penglongli/kubernetes-demo/k8s-terminal-go/k8s"
)

var (
	END_OF_TRANSMISSION = "\u0004"

	validShell = []string{"bash", "sh", "powershell", "cmd"}
)

type TerminalMessage struct {
	Op, Data   string
	Code       int
	Rows, Cols uint16
}

type TerminalSession struct {
	SockSession sockjs.Session
	SizeChan    chan *remotecommand.TerminalSize

	Namespace string
	Pod       string
	Container string

	// Connection will be close after timeout.
	Timeout        time.Duration
	RefreshTimeout chan struct{}
	NeedClose      bool
}

// Read will read the input of front-end by sockjs
func (session *TerminalSession) Read(p []byte) (int, error) {
	if session.NeedClose {
		_ = session.SockSession.Close(128, "connection need to be closed.")
		return copy(p, END_OF_TRANSMISSION), errors.Errorf("connection need to be closed.")
	} else {
		session.RefreshTimeout <- struct{}{}
	}

	m, err := session.SockSession.Recv()
	if err != nil {
		klog.Error(err)
		return copy(p, END_OF_TRANSMISSION), err
	}

	msg := new(TerminalMessage)
	if err := json.Unmarshal([]byte(m), &msg); err != nil {
		return copy(p, END_OF_TRANSMISSION), err
	}

	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		// Disable resize, refer to Next() function
		return 0, nil
	default:
		return copy(p, END_OF_TRANSMISSION), errors.Errorf("unknown message type: %s", msg.Op)
	}
}

// Write will send bytes(p) to front-end by sockjs
func (session *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(&TerminalMessage{
		Op:   "stdout",
		Data: string(p),
	})
	if err != nil {
		klog.Error(err)
		return 0, err
	}

	if err = session.SockSession.Send(string(msg)); err != nil {
		klog.Error(err)
		return 0, err
	}
	return len(p), nil
}

// Next returns the new terminal size after the terminal has been resized. It returns nil when
// monitoring has been stopped.
// NOTE: Only first init terminal size, and then stop the size monitor.
func (session *TerminalSession) Next() *remotecommand.TerminalSize {
	if v, ok := <-session.SizeChan; ok {
		defer close(session.SizeChan)
		return v
	} else {
		return nil
	}
}

func (session *TerminalSession) HandleTimeout() {
	tf := time.After(session.Timeout)
	select {
	case <-tf:
		klog.Errorf("Expire timeout, connection need to be closed.")
		session.NeedClose = true
		session.Read([]byte(""))
		return
	case <-session.RefreshTimeout:
		return
	}
}

func (session *TerminalSession) CheckShellInPod() ([]string, error) {
	clientSet, err := k8s.GetClientSet()
	if err != nil {
		return nil, err
	}

	req := clientSet.CoreV1().RESTClient().Post().Resource("pods").Namespace(session.Namespace).Name(session.Pod).
		SubResource("exec").
		VersionedParams(
			&v1.PodExecOptions{
				Container: session.Container,
				Command:   []string{"ls", "/bin"},
				Stdin:     false,
				Stdout:    true,
				Stderr:    true,
				TTY:       false,
			}, scheme.ParameterCodec)

	restConfig, err := k8s.GetRestConfig()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	str := stdout.String()
	for _, shell := range validShell {
		if strings.Contains(str, shell) {
			return []string{shell}, nil
		}
	}
	return nil, errors.Errorf("no such available shell")
}

func (session *TerminalSession) Exec(cmd []string) error {
	clientSet, err := k8s.GetClientSet()
	if err != nil {
		return err
	}

	req := clientSet.CoreV1().RESTClient().Post().Resource("pods").
		Namespace(session.Namespace).
		Name(session.Pod).
		SubResource("exec").
		VersionedParams(
			&v1.PodExecOptions{
				Container: session.Container,
				Command:   cmd,
				Stdin:     true,
				Stdout:    true,
				Stderr:    true,
				TTY:       true,
			}, scheme.ParameterCodec)

	restConfig, err := k8s.GetRestConfig()
	if err != nil {
		klog.Error(err)
		return err
	}
	exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
	if err != nil {
		klog.Error(err)
		return err
	}

	// Handle websocket timeout
	go func(session *TerminalSession) {
		for {
			session.HandleTimeout()
			if session.NeedClose {
				break
			}
		}
	}(session)

	session.SizeChan <- &remotecommand.TerminalSize{Width: 150, Height: 50}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             session,
		Stdout:            session,
		Stderr:            session,
		TerminalSizeQueue: session,
		Tty:               true,
	})
	if err != nil {
		klog.Error(err)
		return err
	}
	return nil
}

func (session *TerminalSession) Logging() error {
	return nil
}
