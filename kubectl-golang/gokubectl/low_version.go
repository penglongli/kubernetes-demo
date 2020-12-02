package gokubectl
import (
	"context"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	scalev1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/certificates/v1"
	coordinatev1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	schedulev1 "k8s.io/api/scheduling/v1"
	storagev1 "k8s.io/api/storage/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8syaml "sigs.k8s.io/yaml"

	"github.com/penglongli/kubernetes-demo/kubectl-golang/k8s"
)

var (
	kindHandleFunc = map[string]func(context.Context, *k8s.KubeClient, []byte) error{
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
		return nil
	}

	// Handle kind resource
	f, ok := kindHandleFunc[strings.ToLower(typeMeta.Kind)]
	if !ok {
		return errors.Errorf("Unknown kind: %s", typeMeta.Kind)
	}
	return f(low.ctx, low.kubeClient, data)
}

func certificatesigningrequest(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
   	if err != nil {
   		return err
   	}

   	// Decode obj
   	obj := new(v1.CertificateSigningRequest)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CertificatesV1().CertificateSigningRequests().Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CertificatesV1().CertificateSigningRequests().Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CertificatesV1().CertificateSigningRequests().Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func clusterrolebinding(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(rbacv1.ClusterRoleBinding)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.RbacV1().ClusterRoleBindings().Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.RbacV1().ClusterRoleBindings().Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.RbacV1().ClusterRoleBindings().Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func clusterrole(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(rbacv1.ClusterRole)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.RbacV1().ClusterRoles().Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.RbacV1().ClusterRoles().Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.RbacV1().ClusterRoles().Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func configmap(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.ConfigMap)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().ConfigMaps(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().ConfigMaps(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().ConfigMaps(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func cronjob(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(batchv1beta1.CronJob)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.BatchV1beta1().CronJobs(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.BatchV1beta1().CronJobs(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.BatchV1beta1().CronJobs(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func csidriver(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(storagev1.CSIDriver)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.StorageV1().CSIDrivers().Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.StorageV1().CSIDrivers().Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.StorageV1().CSIDrivers().Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func csinode(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(storagev1.CSINode)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.StorageV1().CSINodes().Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.StorageV1().CSINodes().Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.StorageV1().CSINodes().Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func daemonset(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(appsv1.DaemonSet)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.AppsV1().DaemonSets(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.AppsV1().DaemonSets(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.AppsV1().DaemonSets(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func deployment(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(appsv1.Deployment)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.AppsV1().Deployments(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.AppsV1().Deployments(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.AppsV1().Deployments(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func endpoint(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.Endpoints)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().Endpoints(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().Endpoints(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().Endpoints(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func event(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.Event)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().Events(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().Events(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().Events(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func horizontalpodautoscaler(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(scalev1.HorizontalPodAutoscaler)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.AutoscalingV1().HorizontalPodAutoscalers(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.AutoscalingV1().HorizontalPodAutoscalers(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.AutoscalingV1().HorizontalPodAutoscalers(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func ingress(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(networkv1.Ingress)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.NetworkingV1().Ingresses(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.NetworkingV1().Ingresses(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.NetworkingV1().Ingresses(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func job(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(batchv1.Job)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.BatchV1().Jobs(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.BatchV1().Jobs(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.BatchV1().Jobs(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func lease(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(coordinatev1.Lease)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoordinationV1().Leases(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoordinationV1().Leases(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoordinationV1().Leases(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func limitrange(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.LimitRange)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().LimitRanges(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().LimitRanges(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().LimitRanges(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func mutatingwebhookconfiguration(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(admissionv1.MutatingWebhookConfiguration)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.AdmissionregistrationV1().MutatingWebhookConfigurations().
		Get(ctx, obj.GetName(), metav1.GetOptions{
			ResourceVersion: "0",
		})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.AdmissionregistrationV1().MutatingWebhookConfigurations().Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func namespace(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.Namespace)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().Namespaces().Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().Namespaces().Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().Namespaces().Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func networkpolicy(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(networkv1.NetworkPolicy)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.NetworkingV1().NetworkPolicies(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
			ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.NetworkingV1().NetworkPolicies(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.NetworkingV1().NetworkPolicies(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func node(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.Node)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().Nodes().Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().Nodes().Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().Nodes().Update(ctx, obj, metav1.UpdateOptions{})
	return err
}





func persistentvolumeclaim(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.PersistentVolumeClaim)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().PersistentVolumeClaims(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().PersistentVolumeClaims(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().PersistentVolumeClaims(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func persistentvolume(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.PersistentVolume)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().PersistentVolumes().Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().PersistentVolumes().Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().PersistentVolumes().Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func pod(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.Pod)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().Pods(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().Pods(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().Pods(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func priorityclass(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(schedulev1.PriorityClass)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.SchedulingV1().PriorityClasses().Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.SchedulingV1().PriorityClasses().Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.SchedulingV1().PriorityClasses().Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func replicaset(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(appsv1.ReplicaSet)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.AppsV1().ReplicaSets(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.AppsV1().ReplicaSets(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.AppsV1().ReplicaSets(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func replicationcontroller(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.ReplicationController)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().ReplicationControllers(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().ReplicationControllers(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().ReplicationControllers(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func resourcequota(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.ResourceQuota)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().ResourceQuotas(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().ResourceQuotas(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().ResourceQuotas(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func rolebinding(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(rbacv1.RoleBinding)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.RbacV1().RoleBindings(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.RbacV1().RoleBindings(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.RbacV1().RoleBindings(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func role(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(rbacv1.Role)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.RbacV1().Roles(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.RbacV1().Roles(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.RbacV1().Roles(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func secret(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.Secret)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().Secrets(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().Secrets(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().Secrets(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func serviceaccount(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.ServiceAccount)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().ServiceAccounts(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().ServiceAccounts(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().ServiceAccounts(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func service(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(corev1.Service)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.CoreV1().Services(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.CoreV1().Services(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.CoreV1().Services(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func statefulset(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(appsv1.StatefulSet)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.AppsV1().StatefulSets(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.AppsV1().StatefulSets(obj.GetNamespace()).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.AppsV1().StatefulSets(obj.GetNamespace()).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func validatingwebhookconfiguration(ctx context.Context, kubeClient *k8s.KubeClient, data []byte) error {
	clientSet, err := kubeClient.GetClientSet()
	if err != nil {
		return err
	}

	// Decode obj
	obj := new(admissionv1.ValidatingWebhookConfiguration)
	if err = yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	_, err = clientSet.AdmissionregistrationV1().ValidatingWebhookConfigurations().Get(ctx, obj.GetName(), metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		if !k8sErrors.IsNotFound(err) {
			return err
		}
		// Create if not exist
		_, err = clientSet.AdmissionregistrationV1().ValidatingWebhookConfigurations().Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	_, err = clientSet.AdmissionregistrationV1().ValidatingWebhookConfigurations().Update(ctx, obj, metav1.UpdateOptions{})
	return err
}
