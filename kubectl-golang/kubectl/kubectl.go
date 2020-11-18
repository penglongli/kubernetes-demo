package kubectl

import (
	"bufio"
	"bytes"
	"context"
	"io"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"

	"github.com/penglongli/kubernetes-demo/kubectl-golang/k8s"
)

var (
	decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
)

func Apply(ctx context.Context, data []byte) (result []string, err error) {
	chanMes, chanErr := readYaml(data)
	for {
		select {
			case mes, ok := <- chanMes: {
				if !ok {
					return result, nil
				}

				// Get obj and dr
				obj, dr, err := buildDynamicResourceClient(mes)
				if err != nil {
					return result, err
				}

				// Create or Update
				_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, mes, metav1.PatchOptions{
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

func Delete(ctx context.Context, data []byte) (result []string, err error) {
	chanMes, chanErr := readYaml(data)
	for {
		select {
		case mes, ok := <- chanMes: {
			if !ok {
				return result, nil
			}

			// Get obj and dr
			obj, dr, err := buildDynamicResourceClient(mes)
			if err != nil {
				return result, err
			}

			// Delete
			deletePolicy := metav1.DeletePropagationBackground
			err = dr.Delete(ctx, obj.GetName(), metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
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

func buildDynamicResourceClient(data []byte) (obj *unstructured.Unstructured, dr dynamic.ResourceInterface, err error) {
	// Decode YAML manifest into unstructured.Unstructured
	// Get GVK
	obj = &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(data, nil, obj)
	if err != nil {
		return obj, dr, errors.Wrap(err, "received err while decoding yaml.")
	}

	// Prepare discovery mapper
	mapper, err := k8s.GetDiscoveryMapper()
	if err != nil {
		return obj, dr, errors.Wrap(err, "prepare discovery mapper failed")
	}

	// Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return obj, dr, errors.Wrap(err, "mapping kind with version failed")
	}

	// Prepare dynamic client
	dynamicClient, err := k8s.GetDynamicClient()
	if err != nil {
		return obj, dr, errors.Wrap(err, "prepare dynamic client failed.")
	}

	// Obtain REST interface for the GVR
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dynamicClient.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dynamicClient.Resource(mapping.Resource)
	}
	return obj, dr, nil
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

		for {
			buf, err := multidocReader.Read()
			if err != nil {
				if err == io.EOF {
					return
				}
				chanErr <- errors.Wrap(err, "failed to read yaml data")
				return
			}
			chanBytes <- buf
		}
	}()
	return chanBytes, chanErr
}
