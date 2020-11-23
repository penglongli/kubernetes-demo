package k8s

import (
	"encoding/base64"
	"strconv"

	"github.com/pkg/errors"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type KubeClient struct {
	Base64KubeConfig string
}

func (kube *KubeClient) GetClientSet() (*kubernetes.Clientset, error) {
	restConfig, err := kube.buildRestConfig()
	if err != nil {
		return nil, errors.Wrap(err, "build restConfig failed")
	}

	return kubernetes.NewForConfig(restConfig)
}

func (kube *KubeClient) GetDynamicClient() (dynamic.Interface, error) {
	restConfig, err := kube.buildRestConfig()
	if err != nil {
		return nil, errors.Wrap(err, "build restConfig failed")
	}

	return dynamic.NewForConfig(restConfig)
}

func (kube *KubeClient) GetDiscoveryMapper() (*restmapper.DeferredDiscoveryRESTMapper, error) {
	restConfig, err := kube.buildRestConfig()
	if err != nil {
		return nil, errors.Wrap(err, "build restConfig failed")
	}

	// Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "new dc failed")
	}

	discoveryMapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	return discoveryMapper, nil
}

func (kube *KubeClient) CompareVersion() (bool, error) {
	client, err := kube.GetClientSet()
	if err != nil {
		return false, err
	}
	serverInfo, err := client.ServerVersion()
	if err != nil {
		return false, err
	}

	minor, err := strconv.Atoi(serverInfo.Minor)
	if err != nil {
		return false, err
	}
	// version < 1.16
	if serverInfo.Major == "1" && minor < 16 {
		return true, nil
	}
	return false, nil
}

func (kube *KubeClient) buildRestConfig() (resetConfig *rest.Config, err error) {
	kubeConfig, err := base64.StdEncoding.DecodeString(kube.Base64KubeConfig)
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
