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

var testsData []testData

func TestOCVersions(t *testing.T) {

	testsData = append(testsData, testData{getVersionOut: output120, OCVersion: "4.7"})
	testsData = append(testsData, testData{getVersionOut: output121, OCVersion: "4.8"})
	testsData = append(testsData, testData{getVersionOut: output122, OCVersion: "4.9"})

	for _, testdata := range testsData {
		clientset := fake.NewSimpleClientset()
		clientset.Discovery().(*discoveryfake.FakeDiscovery).FakedServerVersion = &version.Info{
			Major: testdata.getVersionOut.Major,
			Minor: testdata.getVersionOut.Minor,
		}
		oc := Oc{clientset: clientset}
		version, err := oc.GetOcVersion()
		if err != nil {
			t.Error(err)
		}
		if version != testdata.OCVersion {
			t.Error(fmt.Sprintf("Version mismatch expected: %s, got: %s", testdata.OCVersion, version))
		}
	}
}
