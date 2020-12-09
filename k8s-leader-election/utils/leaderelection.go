package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

const (
	defaultLeaseDuration = 4 * time.Second
	defaultRenewDeadline = 2 * time.Second
	defaultRetryPeriod   = 1 * time.Second
)

type leaderElection struct {
	// runFunc will execute when acquire the leader.
	runFunc func(ctx context.Context)

	// namespace is to store the lock resource
	namespace string
	// lockName defined the lock resource name between all the member
	lockName string
	// identity is the unique identity of the currently running member
	// default the pod hostname
	identity string
	// resourceLock defines the lock type.
	// k8s.io/client-go/tools/leaderelection/resourcelock/interface.go
	resourceLock string

	leaseDuration time.Duration
	renewDeadline time.Duration
	retryPeriod   time.Duration

	clientset kubernetes.Interface
}

func NewLeaderElectionWithLease(clientset kubernetes.Interface, lockName string, runFunc func(ctx context.Context)) *leaderElection {
	return &leaderElection{
		runFunc:       runFunc,
		lockName:      lockName,
		resourceLock:  resourcelock.LeasesResourceLock,
		leaseDuration: defaultLeaseDuration,
		renewDeadline: defaultRenewDeadline,
		retryPeriod:   defaultRetryPeriod,
		clientset:     clientset,
	}
}

func NewLeaderElectionWithConfigMap(clientset kubernetes.Interface, lockName string, runFunc func(ctx context.Context)) *leaderElection {
	return &leaderElection{
		runFunc:       runFunc,
		lockName:      lockName,
		resourceLock:  resourcelock.ConfigMapsResourceLock,
		leaseDuration: defaultLeaseDuration,
		renewDeadline: defaultRenewDeadline,
		retryPeriod:   defaultRetryPeriod,
		clientset:     clientset,
	}
}

func NewLeaderElectionWithEndpoint(clientset kubernetes.Interface, lockName string, runFunc func(ctx context.Context)) *leaderElection {
	return &leaderElection{
		runFunc:       runFunc,
		lockName:      lockName,
		resourceLock:  resourcelock.EndpointsResourceLock,
		leaseDuration: defaultLeaseDuration,
		renewDeadline: defaultRenewDeadline,
		retryPeriod:   defaultRetryPeriod,
		clientset:     clientset,
	}
}

func (l *leaderElection) Run() error {
	if l.identity == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}
		l.identity = hostname
	}
	if l.namespace == "" {
		l.namespace = inClusterNamespace()
	}

	broadcaster := record.NewBroadcaster()
	broadcaster.StartRecordingToSink(&corev1.EventSinkImpl{Interface: l.clientset.CoreV1().Events(l.namespace)})
	eventRecorder := broadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: fmt.Sprintf("%s/%s", l.lockName, string(l.identity))})

	rlConfig := resourcelock.ResourceLockConfig{
		Identity:      l.identity,
		EventRecorder: eventRecorder,
	}

	lock, err := resourcelock.New(l.resourceLock, l.namespace, l.lockName, l.clientset.CoreV1(), l.clientset.CoordinationV1(), rlConfig)
	if err != nil {
		return err
	}

	leaderConfig := leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: l.leaseDuration,
		RenewDeadline: l.renewDeadline,
		RetryPeriod:   l.retryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				fmt.Println("became leader, starting")
				l.runFunc(ctx)
			},
			OnStoppedLeading: func() {
				fmt.Println("stopped leading")
			},
			OnNewLeader: func(identity string) {
				fmt.Printf("new leader detected, current leader: %s\n", identity)
			},
		},
	}

	leaderelection.RunOrDie(context.TODO(), leaderConfig)
	return nil
}

func inClusterNamespace() string {
	if ns := os.Getenv("POD_NAMESPACE"); ns != "" {
		return ns
	}

	if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
			return ns
		}
	}
	return "default"
}
