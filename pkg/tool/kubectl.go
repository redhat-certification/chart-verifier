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
