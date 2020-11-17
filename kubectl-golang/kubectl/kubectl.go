package kubectl

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"io"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/penglongli/kubernetes-demo/kubectl-golang/k8s"
)

var (
	decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
)

func Apply(ctx context.Context, namespace string, data []byte) (result []string, err error) {
	chanMes, chanErr := readYaml(data)
	for {
		select {
			case mes, ok := <- chanMes: {
				if !ok {
					return result, nil
				}

				// Prepare discovery mapper
				mapper, err := k8s.GetDiscoveryMapper()
				if err != nil {
					return result, err
				}

				// Decode YAML manifest into unstructured.Unstructured
				// Get GVK
				obj := &unstructured.Unstructured{}
				_, gvk, err := decUnstructured.Decode(mes, nil, obj)
				if err != nil {
					return result, err
				}

				// Find GVR
				mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
				if err != nil {
					return result, err
				}

				// Prepare dynamic client
				dynamicClient, err := k8s.GetDynamicClient()
				if err != nil {
					return result, err
				}

				// Obtain REST interface for the GVR
				var dr dynamic.ResourceInterface
				if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
					// namespaced resources should specify the namespace
					dr = dynamicClient.Resource(mapping.Resource).Namespace(obj.GetNamespace())
				} else {
					// for cluster-wide resources
					dr = dynamicClient.Resource(mapping.Resource)
				}

				_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
					FieldManager: "kubectl-golang",
				})
				if err != nil {
					result = append(result, err.Error())
				} else {
					result = append(result, obj.GetName() + " patched.")
				}
			}
			case err, ok := <-chanErr:
				if !ok {
					return result, nil
				}
				if err == nil {
					continue
				}
				result = append(result, err.Error())
		}
	}
}

func readYaml(data []byte) (<-chan []byte, <-chan error) {
	var (
		chanErr        = make(chan error)
		chanBytes        = make(chan []byte)
		multidocReader = utilyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(data)))
	)

	go func() {
		defer close(chanErr)
		defer close(chanBytes)

		buf, err := multidocReader.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			chanErr <- errors.Wrap(err, "failed to read yaml data")
			return
		}
		chanBytes <- buf
	}()
	return chanBytes, chanErr
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
