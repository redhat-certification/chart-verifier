package tool

import (
	"context"
	"fmt"

	"helm.sh/helm/v3/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
)

// Based on https://access.redhat.com/solutions/4870701
var kubeOpenShiftVersionMap map[string]string = map[string]string{
	"1.22": "4.9",
	"1.21": "4.8",
	"1.20": "4.7",
	"1.19": "4.6",
	"1.18": "4.5",
	"1.17": "4.4",
	"1.16": "4.3",
	"1.14": "4.2",
	"1.13": "4.1",
}

type Kubectl struct {
	clientset kubernetes.Interface
}

func NewKubectl(settings *cli.EnvSettings) (*Kubectl, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if len(settings.KubeConfig) > 0 {
		loadingRules = &clientcmd.ClientConfigLoadingRules{ExplicitPath: settings.KubeConfig}
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		&clientcmd.ConfigOverrides{})
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	config.APIPath = "/api"
	config.GroupVersion = &schema.GroupVersion{Group: "core", Version: "v1"}
	config.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	kubectl := new(Kubectl)
	kubectl.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubectl, nil
}

func (k Kubectl) WaitForDeployments(namespace string, selector string) error {
	deployments, err := k.clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return err
	}

	for _, deployment := range deployments.Items {
		// Just after rollout, pods from the previous deployment revision may still be in a
		// terminating state.
		unavailable := deployment.Status.UnavailableReplicas
		if unavailable != 0 {
			return fmt.Errorf("%d replicas unavailable", unavailable)
		}
	}

	return nil
}

func (k Kubectl) DeleteNamespace(namespace string) error {
	if err := k.clientset.CoreV1().Namespaces().Delete(context.TODO(), namespace, *metav1.NewDeleteOptions(0)); err != nil {
		return err
	}
	return nil
}

func (k Kubectl) GetOcVersion() (string, error) {
	version, err := k.clientset.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}

	// Relying on Kubernetes version can be replaced after fixing this issue:
	// https://bugzilla.redhat.com/show_bug.cgi?id=1850656
	kubeVersion := fmt.Sprintf("%s.%s", version.Major, version.Minor)
	osVersion, ok := kubeOpenShiftVersionMap[kubeVersion]
	if !ok {
		return "", fmt.Errorf("internal error: %q not found in Kubernetes-OpenShift version map", kubeVersion)
	}

	return osVersion, nil
}

func GetKubeOpenShiftVersionMap() map[string]string {
	return kubeOpenShiftVersionMap
}
