package tool

import (
	"fmt"

	"helm.sh/helm/v3/pkg/cli"
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

type Oc struct {
	clientset kubernetes.Interface
}

func NewOc(setttings *cli.EnvSettings) (*Oc, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if len(setttings.KubeConfig) > 0 {
		loadingRules = &clientcmd.ClientConfigLoadingRules{ExplicitPath: setttings.KubeConfig}
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
	oc := new(Oc)
	oc.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return oc, nil
}

func (o Oc) GetOcVersion() (string, error) {
	version, err := o.clientset.Discovery().ServerVersion()
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
