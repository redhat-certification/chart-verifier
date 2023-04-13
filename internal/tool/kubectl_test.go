package tool

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	discoveryfake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes/fake"
)

type testData struct {
	getVersionOut version.Info
	OCVersion     string
}

var output120 = version.Info{
	Major: "1",
	Minor: "20",
}

var output121 = version.Info{
	Major: "1",
	Minor: "21",
}

var output122 = version.Info{
	Major: "1",
	Minor: "22",
}

var output123 = version.Info{
	Major: "1",
	Minor: "23",
}

var output124 = version.Info{
	Major: "1",
	Minor: "24",
}

var output125 = version.Info{
	Major: "1",
	Minor: "25",
}

var latestVersion = output125

var testsData []testData

func TestOCVersions(t *testing.T) {
	testsData = append(testsData, testData{getVersionOut: output120, OCVersion: "4.7"})
	testsData = append(testsData, testData{getVersionOut: output121, OCVersion: "4.8"})
	testsData = append(testsData, testData{getVersionOut: output122, OCVersion: "4.9"})
	testsData = append(testsData, testData{getVersionOut: output123, OCVersion: "4.10"})
	testsData = append(testsData, testData{getVersionOut: output124, OCVersion: "4.11"})
	testsData = append(testsData, testData{getVersionOut: output125, OCVersion: "4.12"})

	for _, testdata := range testsData {
		clientset := fake.NewSimpleClientset()
		clientset.Discovery().(*discoveryfake.FakeDiscovery).FakedServerVersion = &version.Info{
			Major: testdata.getVersionOut.Major,
			Minor: testdata.getVersionOut.Minor,
		}
		kubectl := Kubectl{clientset: clientset}
		serverVersion, err := kubectl.GetServerVersion()
		if err != nil {
			t.Error(err)
		}
		if serverVersion.Major != testdata.getVersionOut.Major || serverVersion.Minor != testdata.getVersionOut.Minor {
			t.Errorf("server version mismatch, expected: %+v, got: %+v", testdata.getVersionOut, serverVersion)
		}
		kubeVersion := fmt.Sprintf("%s.%s", serverVersion.Major, serverVersion.Minor)
		ocVersion := GetKubeOpenShiftVersionMap()[kubeVersion]
		if ocVersion != testdata.OCVersion {
			t.Errorf("version mismatch, expected: %s, got: %s", testdata.OCVersion, ocVersion)
		}
	}

	latestKV := GetLatestKubeVersion()
	expectedLatestKV := fmt.Sprintf("%s.%s.0", latestVersion.Major, latestVersion.Minor)
	if latestKV != expectedLatestKV {
		t.Errorf("latest kubversion mismatch, expected: %s, got: %s", expectedLatestKV, latestKV)
	}
}

var testDeployments = []v1.Deployment{
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test0"},
		Status:     v1.DeploymentStatus{UnavailableReplicas: 1},
	},
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test1"},
		Status:     v1.DeploymentStatus{UnavailableReplicas: 2},
	},
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test2"},
		Status:     v1.DeploymentStatus{UnavailableReplicas: 3},
	},
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test3"},
		Status:     v1.DeploymentStatus{UnavailableReplicas: 4},
	},
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test4"},
		Status:     v1.DeploymentStatus{UnavailableReplicas: 5},
	},
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test5"},
		Status:     v1.DeploymentStatus{UnavailableReplicas: 6},
	},
}

var DeploymentList1 []v1.Deployment

func TestWaitForDeployments(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	listDeployments = deploymentTestListGood
	DeploymentList1 = make([]v1.Deployment, len(testDeployments))
	copy(DeploymentList1, testDeployments)

	k := new(Kubectl)
	err := k.WaitForDeployments(ctx, "testNameSpace", "selector")
	require.NoError(t, err)
}

func TestBadToGoodWaitForDeployments(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	listDeployments = deploymentTestListBadToGood
	DeploymentList1 = make([]v1.Deployment, len(testDeployments))
	copy(DeploymentList1, testDeployments)

	k := new(Kubectl)
	err := k.WaitForDeployments(ctx, "testNameSpace", "selector")
	require.NoError(t, err)
}

func TestTimeExpirationWaitingForDeployments(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	listDeployments = deploymentTestListGood
	DeploymentList1 = make([]v1.Deployment, len(testDeployments))
	copy(DeploymentList1, testDeployments)

	k := new(Kubectl)
	err := k.WaitForDeployments(ctx, "testNameSpace", "selector")
	require.Error(t, err)
	require.Contains(t, err.Error(), "error unavailable deployments, timeout has expired,")
}

func TestTimeExpirationGetDeploymentsFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	listDeployments = deploymentTestListBad
	DeploymentList1 = make([]v1.Deployment, len(testDeployments))
	copy(DeploymentList1, testDeployments)

	k := new(Kubectl)
	err := k.WaitForDeployments(ctx, "testNameSpace", "selector")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Time out retrying after")
	require.Contains(t, err.Error(), "error getting deployments from namespace")
	require.Contains(t, err.Error(), "pretend error getting deployment list")
}

func deploymentTestListGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.Deployment, error) {
	fmt.Println("deploymentTestListGood called")
	for index := 0; index < len(DeploymentList1); index++ {
		if DeploymentList1[index].Status.UnavailableReplicas > 0 {
			DeploymentList1[index].Status.UnavailableReplicas--
			fmt.Printf("UnavailableReplicas set to %d for deployment %s\n", DeploymentList1[index].Status.UnavailableReplicas, DeploymentList1[index].Name)
		}
	}
	return DeploymentList1, nil
}

func deploymentTestListBad(k Kubectl, context context.Context, namespace string, selector string) ([]v1.Deployment, error) {
	fmt.Println("deploymentTestListBad called")
	return nil, errors.New("pretend error getting deployment list")
}

var errorSent = false

func deploymentTestListBadToGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.Deployment, error) {
	if !errorSent {
		fmt.Println("deploymentTestListBadToGoodToBad bad path")
		errorSent = true
		return nil, errors.New("pretend error getting deployment list")
	}
	fmt.Println("deploymentTestListBadToGoodToBad good path")
	return deploymentTestListGood(k, context, namespace, selector)
}
