package tool

import (
	"fmt"
	"testing"

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

var testsData []testData

func TestOCVersions(t *testing.T) {

	testsData = append(testsData, testData{getVersionOut: output120, OCVersion: "4.7"})
	testsData = append(testsData, testData{getVersionOut: output121, OCVersion: "4.8"})
	testsData = append(testsData, testData{getVersionOut: output122, OCVersion: "4.9"})
	testsData = append(testsData, testData{getVersionOut: output123, OCVersion: "4.10"})
	testsData = append(testsData, testData{getVersionOut: output124, OCVersion: "4.11"})

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
			t.Error(fmt.Sprintf("server version mismatch, expected: %+v, got: %+v", testdata.getVersionOut, serverVersion))
		}
		kubeVersion := fmt.Sprintf("%s.%s", serverVersion.Major, serverVersion.Minor)
		ocVersion := GetKubeOpenShiftVersionMap()[kubeVersion]
		if ocVersion != testdata.OCVersion {
			t.Error(fmt.Sprintf("version mismatch, expected: %s, got: %s", testdata.OCVersion, ocVersion))
		}
	}
}
