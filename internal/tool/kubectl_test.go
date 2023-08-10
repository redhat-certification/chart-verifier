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

var output126 = version.Info{
	Major: "1",
	Minor: "26",
}

var latestVersion = output126

var testsData []testData

func TestOCVersions(t *testing.T) {
	testsData = append(testsData, testData{getVersionOut: output120, OCVersion: "4.7"})
	testsData = append(testsData, testData{getVersionOut: output121, OCVersion: "4.8"})
	testsData = append(testsData, testData{getVersionOut: output122, OCVersion: "4.9"})
	testsData = append(testsData, testData{getVersionOut: output123, OCVersion: "4.10"})
	testsData = append(testsData, testData{getVersionOut: output124, OCVersion: "4.11"})
	testsData = append(testsData, testData{getVersionOut: output125, OCVersion: "4.12"})
	testsData = append(testsData, testData{getVersionOut: output126, OCVersion: "4.13"})

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

var testDaemonSets = []v1.DaemonSet{
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test0"},
		Status:     v1.DaemonSetStatus{NumberUnavailable: 1},
	},
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test1"},
		Status:     v1.DaemonSetStatus{NumberUnavailable: 2},
	},
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test2"},
		Status:     v1.DaemonSetStatus{NumberUnavailable: 3},
	},
}

var testStatefulSets = []v1.StatefulSet{
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test0"},
		Status:     v1.StatefulSetStatus{Replicas: 1, AvailableReplicas: 0},
	},
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test1"},
		Status:     v1.StatefulSetStatus{Replicas: 2, AvailableReplicas: 0},
	},
	{
		ObjectMeta: metav1.ObjectMeta{Name: "test0"},
		Status:     v1.StatefulSetStatus{Replicas: 3, AvailableReplicas: 0},
	},
}

var DeploymentList1 []v1.Deployment
var DaemonSetList1 []v1.DaemonSet
var StatefulSetList1 []v1.StatefulSet

func TestGetDeploymentList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientset := fake.NewSimpleClientset()
	_, err := clientset.AppsV1().Deployments("").Create(ctx, &testDeployments[0], metav1.CreateOptions{})

	kubectl := Kubectl{clientset: clientset}
	deps, err := getDeploymentsList(kubectl, ctx, "", "")

	require.NoError(t, err)
	require.Equal(t, 1, len(deps))
	require.Equal(t, testDeployments[0], deps[0])

}

func TestGetDaemonSetList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientset := fake.NewSimpleClientset()
	_, err := clientset.AppsV1().DaemonSets("").Create(ctx, &testDaemonSets[0], metav1.CreateOptions{})

	kubectl := Kubectl{clientset: clientset}
	daemons, err := getDaemonSetsList(kubectl, ctx, "", "")

	require.NoError(t, err)
	require.Equal(t, 1, len(daemons))
	require.Equal(t, testDaemonSets[0], daemons[0])

}

func TestGetStatefulSetList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientset := fake.NewSimpleClientset()
	_, err := clientset.AppsV1().StatefulSets("").Create(ctx, &testStatefulSets[0], metav1.CreateOptions{})

	kubectl := Kubectl{clientset: clientset}
	sts, err := getStatefulSetsList(kubectl, ctx, "", "")

	require.NoError(t, err)
	require.Equal(t, 1, len(sts))
	require.Equal(t, testStatefulSets[0], sts[0])

}

func TestWaitForWorkloadResources(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	listDeployments = deploymentTestListGood
	listDaemonSets = daemonSetTestListGood
	listStatefulSets = statefulSetTestListGood

	DeploymentList1 = make([]v1.Deployment, len(testDeployments))
	DaemonSetList1 = make([]v1.DaemonSet, len(testDaemonSets))
	StatefulSetList1 = make([]v1.StatefulSet, len(testStatefulSets))

	copy(DeploymentList1, testDeployments)
	copy(DaemonSetList1, testDaemonSets)
	copy(StatefulSetList1, testStatefulSets)

	k := new(Kubectl)
	err := k.WaitForWorkloadResources(ctx, "testNameSpace", "selector")
	require.NoError(t, err)
}

func TestBadToGoodWaitForDeployments(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	listDeployments = deploymentTestListBadToGood
	DeploymentList1 = make([]v1.Deployment, len(testDeployments))
	copy(DeploymentList1, testDeployments)

	listDaemonSets = daemonSetTestListGood
	DaemonSetList1 = make([]v1.DaemonSet, len(testDaemonSets))
	copy(DaemonSetList1, testDaemonSets)

	listStatefulSets = statefulSetTestListGood
	StatefulSetList1 = make([]v1.StatefulSet, len(testStatefulSets))
	copy(StatefulSetList1, testStatefulSets)

	k := new(Kubectl)
	err := k.WaitForWorkloadResources(ctx, "testNameSpace", "selector")
	require.NoError(t, err)
}

func TestBadToGoodWaitForDaemonSets(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	listDeployments = deploymentTestListGood
	DeploymentList1 = make([]v1.Deployment, len(testDeployments))
	copy(DeploymentList1, testDeployments)

	listDaemonSets = daemonSetTestListBadToGood
	DaemonSetList1 = make([]v1.DaemonSet, len(testDaemonSets))
	copy(DaemonSetList1, testDaemonSets)

	listStatefulSets = statefulSetTestListGood
	StatefulSetList1 = make([]v1.StatefulSet, len(testStatefulSets))
	copy(StatefulSetList1, testStatefulSets)

	k := new(Kubectl)
	err := k.WaitForWorkloadResources(ctx, "testNameSpace", "selector")
	require.NoError(t, err)
}

func TestBadToGoodWaitForStatefulSets(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	listDeployments = deploymentTestListGood
	DeploymentList1 = make([]v1.Deployment, len(testDeployments))
	copy(DeploymentList1, testDeployments)

	listDaemonSets = daemonSetTestListGood
	DaemonSetList1 = make([]v1.DaemonSet, len(testDaemonSets))
	copy(DaemonSetList1, testDaemonSets)

	listStatefulSets = statefulSetTestListBadToGood
	StatefulSetList1 = make([]v1.StatefulSet, len(testStatefulSets))
	copy(StatefulSetList1, testStatefulSets)

	k := new(Kubectl)
	err := k.WaitForWorkloadResources(ctx, "testNameSpace", "selector")
	require.NoError(t, err)

}

func TestTimeExpirationWaitingForWorkloadResources(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	listDeployments = deploymentTestListGood
	DeploymentList1 = make([]v1.Deployment, len(testDeployments))
	copy(DeploymentList1, testDeployments)

	listDaemonSets = daemonSetTestListEmpty
	listStatefulSets = statefulSetTestListEmpty

	k := new(Kubectl)
	err := k.WaitForWorkloadResources(ctx, "testNameSpace", "selector")
	require.Error(t, err)
	require.Contains(t, err.Error(), "error unavailable workload resources, timeout has expired,")
}

func TestTimeExpirationGetDeploymentsFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	listDeployments = deploymentTestListBad
	DeploymentList1 = make([]v1.Deployment, len(testDeployments))
	copy(DeploymentList1, testDeployments)

	k := new(Kubectl)
	err := k.WaitForWorkloadResources(ctx, "testNameSpace", "selector")
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

func daemonSetTestListGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.DaemonSet, error) {
	fmt.Println("daemonSetTestListGood called")
	for index := 0; index < len(testDaemonSets); index++ {
		if DaemonSetList1[index].Status.NumberUnavailable > 0 {
			DaemonSetList1[index].Status.NumberUnavailable--
			fmt.Printf("NumberUnavailable set to %d for daemonset %s\n", DaemonSetList1[index].Status.NumberUnavailable, DaemonSetList1[index].Name)
		}
	}
	return DaemonSetList1, nil
}
func statefulSetTestListGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.StatefulSet, error) {
	fmt.Println("statefulSetSetTestListGood called")
	for index := 0; index < len(testDaemonSets); index++ {
		unavailableReplicas := StatefulSetList1[index].Status.Replicas - StatefulSetList1[index].Status.AvailableReplicas
		if unavailableReplicas > 0 {
			StatefulSetList1[index].Status.AvailableReplicas++
			unavailableReplicas = StatefulSetList1[index].Status.Replicas - StatefulSetList1[index].Status.AvailableReplicas

			fmt.Printf("Unavailable Replicas set to %d for statefulset %s\n", unavailableReplicas, StatefulSetList1[index].Name)
		}
	}
	return StatefulSetList1, nil
}

func deploymentTestListBad(k Kubectl, context context.Context, namespace string, selector string) ([]v1.Deployment, error) {
	fmt.Println("deploymentTestListBad called")
	return nil, errors.New("pretend error getting deployment list")
}

var deploymentErrorSent = false
var daemonSetErrorSent = false
var statefulSetErrorSent = false

func deploymentTestListBadToGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.Deployment, error) {
	if !deploymentErrorSent {
		fmt.Println("deploymentTestListBadToGood bad path")
		deploymentErrorSent = true
		return nil, errors.New("pretend error getting deployment list")
	}
	fmt.Println("deploymentTestListBadToGood good path")
	return deploymentTestListGood(k, context, namespace, selector)
}

func daemonSetTestListBadToGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.DaemonSet, error) {
	if !daemonSetErrorSent {
		fmt.Println("daemonSetTestListBadToGood bad path")
		daemonSetErrorSent = true
		return nil, errors.New("pretend error getting daemonSet list")
	}
	fmt.Println("deploymentTestListBadToGood good path")
	return daemonSetTestListGood(k, context, namespace, selector)
}

func statefulSetTestListBadToGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.StatefulSet, error) {
	if !statefulSetErrorSent {
		fmt.Println("statefulSetTestListBadToGood bad path")
		statefulSetErrorSent = true
		return nil, errors.New("pretend error getting statefulSet list")
	}
	fmt.Println("statefulSetTestListBadToGood good path")
	return statefulSetTestListGood(k, context, namespace, selector)
}

func daemonSetTestListEmpty(k Kubectl, context context.Context, namespace string, selector string) ([]v1.DaemonSet, error) {
	return []v1.DaemonSet{}, nil
}

func statefulSetTestListEmpty(k Kubectl, context context.Context, namespace string, selector string) ([]v1.StatefulSet, error) {
	return []v1.StatefulSet{}, nil
}
