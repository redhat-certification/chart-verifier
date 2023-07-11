package tool

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/cli"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/utils"
)

//go:embed kubeOpenShiftVersionMap.yaml
var content embed.FS

// Based on https://access.redhat.com/solutions/4870701
var (
	kubeOpenShiftVersionMap map[string]string
	listDeployments         = getDeploymentsList
	latestKubeVersion       *semver.Version
)

type versionMap struct {
	Versions []*versionMapping `yaml:"versions"`
}

type versionMapping struct {
	KubeVersion string `yaml:"kube-version"`
	OcpVersion  string `yaml:"ocp-version"`
}

type deploymentNotReady struct {
	Name        string
	Unavailable int32
}

func init() {
	kubeOpenShiftVersionMap = make(map[string]string)

	yamlFile, err := content.ReadFile("kubeOpenShiftVersionMap.yaml")
	if err != nil {
		utils.LogError(fmt.Sprintf("Error reading content of kubeOpenShiftVersionMap.yaml: %v", err))
		return
	}

	versions := versionMap{}
	err = yaml.Unmarshal(yamlFile, &versions)
	if err != nil {
		utils.LogError(fmt.Sprintf("Error reading content of kubeOpenShiftVersionMap.yaml: %v", err))
		return
	}

	latestKubeVersion, _ = semver.NewVersion("0.0")
	for _, versionMap := range versions.Versions {
		currentVersion, _ := semver.NewVersion(versionMap.KubeVersion)
		if currentVersion.GreaterThan(latestKubeVersion) {
			latestKubeVersion = currentVersion
		}
		kubeOpenShiftVersionMap[versionMap.KubeVersion] = versionMap.OcpVersion
	}
}

type Kubectl struct {
	clientset kubernetes.Interface
}

func NewKubectl(kubeConfig clientcmd.ClientConfig) (*Kubectl, error) {
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

func (k Kubectl) WaitForDeployments(context context.Context, namespace string, selector string) error {
	deadline, _ := context.Deadline()
	unavailableDeployments := []deploymentNotReady{{Name: "none", Unavailable: 1}}
	getDeploymentsError := ""

	utils.LogInfo(fmt.Sprintf("Start wait for deployments. --timeout time left: %s ", time.Until(deadline).String()))

	for deadline.After(time.Now()) && len(unavailableDeployments) > 0 {
		unavailableDeployments = []deploymentNotReady{}
		deployments, err := listDeployments(k, context, namespace, selector)
		if err != nil {
			unavailableDeployments = []deploymentNotReady{{Name: "none", Unavailable: 1}}
			getDeploymentsError = fmt.Sprintf("error getting deployments from namespace %s : %v", namespace, err)
			utils.LogWarning(getDeploymentsError)
			time.Sleep(time.Second)
		} else {
			getDeploymentsError = ""
			for _, deployment := range deployments {
				// Just after rollout, pods from the previous deployment revision may still be in a
				// terminating state.
				if deployment.Status.UnavailableReplicas > 0 {
					unavailableDeployments = append(unavailableDeployments, deploymentNotReady{Name: deployment.Name, Unavailable: deployment.Status.UnavailableReplicas})
				}
			}
			if len(unavailableDeployments) > 0 {
				utils.LogInfo(fmt.Sprintf("Wait for %d deployments:", len(unavailableDeployments)))
				for _, unavailableDeployment := range unavailableDeployments {
					utils.LogInfo(fmt.Sprintf("    - %s with %d unavailable replicas", unavailableDeployment.Name, unavailableDeployment.Unavailable))
				}
				time.Sleep(time.Second)
			} else {
				utils.LogInfo(fmt.Sprintf("Finish wait for deployments, --timeout time left %s", time.Until(deadline).String()))
			}
		}
	}

	if len(getDeploymentsError) > 0 {
		errorMsg := fmt.Sprintf("Time out retrying after %s", getDeploymentsError)
		utils.LogError(errorMsg)
		return errors.New(errorMsg)
	}
	if len(unavailableDeployments) > 0 {
		errorMsg := "error unavailable deployments, timeout has expired, please consider increasing the timeout using the chart-verifier --timeout flag"
		utils.LogError(errorMsg)
		return errors.New(errorMsg)
	}

	return nil
}

func (k Kubectl) DeleteNamespace(context context.Context, namespace string) error {
	if err := k.clientset.CoreV1().Namespaces().Delete(context, namespace, *metav1.NewDeleteOptions(0)); err != nil {
		return err
	}
	return nil
}

func (k Kubectl) GetServerVersion() (*version.Info, error) {
	version, err := k.clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, err
	}
	return version, err
}

func GetKubeOpenShiftVersionMap() map[string]string {
	return kubeOpenShiftVersionMap
}

func GetClientConfig(envSettings *cli.EnvSettings) clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if len(envSettings.KubeConfig) > 0 {
		loadingRules = &clientcmd.ClientConfigLoadingRules{ExplicitPath: envSettings.KubeConfig}
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		&clientcmd.ConfigOverrides{CurrentContext: envSettings.KubeContext})
}

func getDeploymentsList(k Kubectl, context context.Context, namespace string, selector string) ([]v1.Deployment, error) {
	list, err := k.clientset.AppsV1().Deployments(namespace).List(context, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	return list.Items, err
}

func GetLatestKubeVersion() string {
	return latestKubeVersion.String()
}
