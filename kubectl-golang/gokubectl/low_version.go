package gokubectl

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	k8syaml "sigs.k8s.io/yaml"

	"github.com/penglongli/kubernetes-demo/kubectl-golang/k8s"
)

var (
	kindHandleFunc = map[string]func(context.Context, *k8s.KubeClient, *unstructured.Unstructured, []byte) error{
		"certificatesigningrequest":      certificatesigningrequest,
		"clusterrolebinding":             clusterrolebinding,
		"clusterrole":                    clusterrole,
		"configmap":                      configmap,
		"cronjob":                        cronjob,
		"csidriver":                      csidriver,
		"csinode":                        csinode,
		"daemonset":                      daemonset,
		"deployment":                     deployment,
		"endpoint":                       endpoint,
		"event":                          event,
		"horizontalpodautoscaler":        horizontalpodautoscaler,
		"ingress":                        ingress,
		"job":                            job,
		"lease":                          lease,
		"limitrange":                     limitrange,
		"mutatingwebhookconfiguration":   mutatingwebhookconfiguration,
		"namespace":                      namespace,
		"networkpolicy":                  networkpolicy,
		"node":                           node,
		"persistentvolumeclaim":          persistentvolumeclaim,
		"persistentvolume":               persistentvolume,
		"pod":                            pod,
		"priorityclass":                  priorityclass,
		"replicaset":                     replicaset,
		"replicationcontroller":          replicationcontroller,
		"resourcequota":                  resourcequota,
		"rolebinding":                    rolebinding,
		"role":                           role,
		"secret":                         secret,
		"serviceaccount":                 serviceaccount,
		"service":                        service,
		"statefulset":                    statefulset,
		"validatingwebhookconfiguration": validatingwebhookconfiguration,
	}
)

type lowVersion struct {
	ctx        context.Context
	kubeClient *k8s.KubeClient
}

func (low *lowVersion) apply(data []byte) error {
	var typeMeta runtime.TypeMeta
	if err := k8syaml.Unmarshal(data, &typeMeta); err != nil {
		return errors.Wrapf(err, "Decode yaml failed. ")
	}
	if typeMeta.Kind == "" {
		return errors.Errorf("Type kind is empty.")
	}

	// Decode to unstructured object
	obj := &unstructured.Unstructured{}
	_, _, err := decUnstructured.Decode(data, nil, obj)
	if err != nil {
		return errors.Wrapf(err, "Unmarshal yaml failed.")
	}

	// Handle kind resource
	f, ok := kindHandleFunc[strings.ToLower(obj.GetKind())]
	if !ok {
		return errors.Errorf("Unknown kind: %s", obj.GetKind())
	}
	return f(low.ctx, low.kubeClient, obj, data)
}

func certificatesigningrequest(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CertificatesV1().CertificateSigningRequests().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func clusterrolebinding(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.RbacV1().ClusterRoleBindings().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func clusterrole(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.RbacV1().ClusterRoles().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func configmap(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().ConfigMaps(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func cronjob(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.BatchV1beta1().CronJobs(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func csidriver(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.StorageV1().CSIDrivers().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func csinode(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.StorageV1().CSINodes().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func daemonset(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.AppsV1().DaemonSets(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func deployment(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.AppsV1().Deployments(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func endpoint(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().Endpoints(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func event(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().Events(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func horizontalpodautoscaler(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.AutoscalingV1().HorizontalPodAutoscalers(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func ingress(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.NetworkingV1().IngressClasses().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func job(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.BatchV1().Jobs(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func lease(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoordinationV1().Leases(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func limitrange(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().LimitRanges(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func mutatingwebhookconfiguration(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.AdmissionregistrationV1().MutatingWebhookConfigurations().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func namespace(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().Namespaces().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}
func networkpolicy(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.NetworkingV1().NetworkPolicies(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func node(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().Nodes().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func persistentvolumeclaim(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().PersistentVolumeClaims(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func persistentvolume(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().PersistentVolumes().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func pod(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().Pods(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func priorityclass(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.SchedulingV1().PriorityClasses().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func replicaset(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.AppsV1().ReplicaSets(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func replicationcontroller(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().ReplicationControllers(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func resourcequota(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().ResourceQuotas(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func rolebinding(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.RbacV1().RoleBindings(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func role(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.RbacV1().Roles(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}
func secret(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().Secrets(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}
func serviceaccount(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().ServiceAccounts(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}
func service(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.CoreV1().Services(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}

func statefulset(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.AppsV1().StatefulSets(obj.GetNamespace()).Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}
func validatingwebhookconfiguration(ctx context.Context, kubeClient *k8s.KubeClient, obj *unstructured.Unstructured, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
    if err != nil {
        return err
    }

	_, err = clientSet.AdmissionregistrationV1().ValidatingWebhookConfigurations().Patch(ctx, obj.GetName(), types.ApplyPatchType,
		data, metav1.PatchOptions{FieldManager: "kubectl-golang"})
	return err
}
