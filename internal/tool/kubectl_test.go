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
	"k8s.io/client-go/kubernetes/fake"
)

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

var (
	DeploymentList1  []v1.Deployment
	DaemonSetList1   []v1.DaemonSet
	StatefulSetList1 []v1.StatefulSet
)

func TestGetDeploymentList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientset := fake.NewSimpleClientset()
	_, err := clientset.AppsV1().Deployments("").Create(ctx, &testDeployments[0], metav1.CreateOptions{})
	require.NoError(t, err)

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
	require.NoError(t, err)

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
	require.NoError(t, err)

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

	listDaemonSets = daemonSetTestListGood
	listStatefulSets = statefulSetTestListGood

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
	require.Contains(t, err.Error(), "error getting Deployment from namespace")
	require.Contains(t, err.Error(), "pretend error getting deployment list")
}

func deploymentTestListGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.Deployment, error) {
	// Mock listing Deployments
	// Return contents of testDeployments and decrease
	// the nuber of unavailable pods on each call
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
	// Mock listing DaemonSets
	// Return contents of testDaemonSets and decrease
	// the nuber of unavailable pods on each call
	for index := 0; index < len(testDaemonSets); index++ {
		if DaemonSetList1[index].Status.NumberUnavailable > 0 {
			DaemonSetList1[index].Status.NumberUnavailable--
			fmt.Printf("NumberUnavailable set to %d for daemonset %s\n", DaemonSetList1[index].Status.NumberUnavailable, DaemonSetList1[index].Name)
		}
	}
	return DaemonSetList1, nil
}

func statefulSetTestListGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.StatefulSet, error) {
	// Mock listing StatefulSets
	// Return contents of testStatefulSets and increase
	// the nuber of available pods on each call
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
	// Mock listing Deployments with error
	fmt.Println("deploymentTestListBad called")
	return nil, errors.New("pretend error getting deployment list")
}

var (
	deploymentErrorSent  = false
	daemonSetErrorSent   = false
	statefulSetErrorSent = false
)

func deploymentTestListBadToGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.Deployment, error) {
	// Mock listing Deployments
	// Return error on first call and pass to deploymentTestListGood
	// on subsequent calls
	if !deploymentErrorSent {
		fmt.Println("deploymentTestListBadToGood bad path")
		deploymentErrorSent = true
		return nil, errors.New("pretend error getting deployment list")
	}
	fmt.Println("deploymentTestListBadToGood good path")
	return deploymentTestListGood(k, context, namespace, selector)
}

func daemonSetTestListBadToGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.DaemonSet, error) {
	// Mock listing DaemonSets
	// Return error on first call and pass to daemonSetTestListGood
	// on subsequent calls
	if !daemonSetErrorSent {
		fmt.Println("daemonSetTestListBadToGood bad path")
		daemonSetErrorSent = true
		return nil, errors.New("pretend error getting daemonSet list")
	}
	fmt.Println("deploymentTestListBadToGood good path")
	return daemonSetTestListGood(k, context, namespace, selector)
}

func statefulSetTestListBadToGood(k Kubectl, context context.Context, namespace string, selector string) ([]v1.StatefulSet, error) {
	// Mock listing StatefulSets
	// Return error on first call and pass to statefulSetTestListGood
	// on subsequent calls
	if !statefulSetErrorSent {
		fmt.Println("statefulSetTestListBadToGood bad path")
		statefulSetErrorSent = true
		return nil, errors.New("pretend error getting statefulSet list")
	}
	fmt.Println("statefulSetTestListBadToGood good path")
	return statefulSetTestListGood(k, context, namespace, selector)
}
