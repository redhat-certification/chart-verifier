package profiles

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestProfile(t *testing.T) {

	testProfile := Profile{}

	testProfile.Apiversion = "v1"
	testProfile.Kind = "verifier-profile"
	testProfile.Name = "profile-1.0.0"

	testProfile.Annotations = []Annotation{DigestAnnotation, OCPVersionAnnotation, LastCertifiedTimestampAnnotation}

	testProfile.Checks = []*ProfileCheck{
		{Name: fmt.Sprintf("%s/%s", "v1.0", checks.HasReadmeName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", "v1.0", checks.IsHelmV3Name), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", "v1.0", checks.ContainsTestName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", "v1.0", checks.ContainsValuesName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", "v1.0", checks.ContainsValuesSchemaName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", "v1.0", checks.HasKubeversionName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", "v1.0", checks.NotContainsCRDsName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", "v1.0", checks.HelmLintName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", "v1.0", checks.NotContainCsiObjectsName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", "v1.0", checks.ImagesAreCertifiedName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", "v1.0", checks.ChartTestingName), Type: checks.MandatoryCheckType},
	}

	t.Run("Profile read from disk should match test profile", func(t *testing.T) {
		diskProfile := GetProfile()
		assert.True(t, cmp.Equal(diskProfile, &testProfile), "profiles do not match")

	})

}

func TestProfileFilter(t *testing.T) {

	defaultRegistry := checks.NewRegistry()

	defaultRegistry.Add(checks.HasReadmeName, "v1.0", checks.HasReadme)
	defaultRegistry.Add(checks.IsHelmV3Name, "v1.0", checks.IsHelmV3)
	defaultRegistry.Add(checks.ContainsTestName, "v1.0", checks.ContainsTest)
	defaultRegistry.Add(checks.ContainsValuesName, "v1.0", checks.ContainsValues)

	defaultRegistry.Add(checks.HasReadmeName, "v1.1", checks.HasReadme)
	defaultRegistry.Add(checks.IsHelmV3Name, "v1.1", checks.IsHelmV3)
	defaultRegistry.Add(checks.ContainsTestName, "v1.1", checks.ContainsTest)
	defaultRegistry.Add(checks.ContainsValuesName, "v1.1", checks.ContainsValues)
	defaultRegistry.Add(checks.ContainsValuesSchemaName, "v1.1", checks.ContainsValuesSchema)
	defaultRegistry.Add(checks.HasKubeversionName, "v1.1", checks.HasKubeVersion)
	defaultRegistry.Add(checks.NotContainsCRDsName, "v1.1", checks.NotContainCRDs)
	defaultRegistry.Add(checks.HelmLintName, "v1.1", checks.HelmLint)
	defaultRegistry.Add(checks.NotContainsCRDsName, "v1.1", checks.NotContainCSIObjects)
	defaultRegistry.Add(checks.ImagesAreCertifiedName, "v1.1", checks.ImagesAreCertified)
	defaultRegistry.Add(checks.ChartTestingName, "v1.1", checks.ChartTesting)

	defaultRegistry.Add("BadHasReadme", "v1.0", checks.HasReadme)
	defaultRegistry.Add("BadIsHelmV3Name", "v1.0", checks.IsHelmV3)
	defaultRegistry.Add("BadContainsTestName", "v1.o", checks.ContainsTest)

	expectedChecks := make(map[checks.CheckName]checks.Check)
	expectedChecks[checks.HasReadmeName] = checks.Check{CheckId: checks.CheckId{Name: checks.HasReadmeName, Version: "v1.0"}, Type: checks.MandatoryCheckType, Func: checks.HasReadme}
	expectedChecks[checks.IsHelmV3Name] = checks.Check{CheckId: checks.CheckId{Name: checks.IsHelmV3Name, Version: "v1.0"}, Type: checks.MandatoryCheckType, Func: checks.IsHelmV3}
	expectedChecks[checks.ContainsTestName] = checks.Check{CheckId: checks.CheckId{Name: checks.ContainsTestName, Version: "v1.0"}, Type: checks.MandatoryCheckType, Func: checks.ContainsTest}
	expectedChecks[checks.ContainsValuesName] = checks.Check{CheckId: checks.CheckId{Name: checks.ContainsValuesName, Version: "v1.0"}, Type: checks.MandatoryCheckType, Func: checks.ContainsValues}

	t.Run("Checks filtered using profile subset", func(t *testing.T) {
		filteredChecks := GetProfile().FilterChecks(defaultRegistry.AllChecks())
		compareCheckMaps(t, expectedChecks, filteredChecks)
	})

	defaultRegistry.Add(checks.ContainsValuesSchemaName, "v1.0", checks.ContainsValuesSchema)
	defaultRegistry.Add(checks.HasKubeversionName, "v1.0", checks.HasKubeVersion)
	defaultRegistry.Add(checks.NotContainsCRDsName, "v1.0", checks.NotContainCRDs)
	defaultRegistry.Add(checks.HelmLintName, "v1.0", checks.HelmLint)
	defaultRegistry.Add(checks.NotContainCsiObjectsName, "v1.0", checks.NotContainCSIObjects)
	defaultRegistry.Add(checks.ImagesAreCertifiedName, "v1.0", checks.ImagesAreCertified)
	defaultRegistry.Add(checks.ChartTestingName, "v1.0", checks.ChartTesting)

	expectedChecks[checks.ContainsValuesSchemaName] = checks.Check{CheckId: checks.CheckId{Name: checks.ContainsValuesSchemaName, Version: "v1.0"}, Type: checks.MandatoryCheckType, Func: checks.ContainsValuesSchema}
	expectedChecks[checks.HasKubeversionName] = checks.Check{CheckId: checks.CheckId{Name: checks.HasKubeversionName, Version: "v1.0"}, Type: checks.MandatoryCheckType, Func: checks.HasKubeVersion}
	expectedChecks[checks.NotContainsCRDsName] = checks.Check{CheckId: checks.CheckId{Name: checks.NotContainsCRDsName, Version: "v1.0"}, Type: checks.MandatoryCheckType, Func: checks.NotContainCRDs}
	expectedChecks[checks.HelmLintName] = checks.Check{CheckId: checks.CheckId{Name: checks.HelmLintName, Version: "v1.0"}, Type: checks.MandatoryCheckType, Func: checks.HelmLint}
	expectedChecks[checks.NotContainCsiObjectsName] = checks.Check{CheckId: checks.CheckId{Name: checks.NotContainCsiObjectsName, Version: "v1.0"}, Type: checks.MandatoryCheckType, Func: checks.NotContainCSIObjects}
	expectedChecks[checks.ImagesAreCertifiedName] = checks.Check{CheckId: checks.CheckId{Name: checks.ImagesAreCertifiedName, Version: "v1.0"}, Type: checks.MandatoryCheckType, Func: checks.ImagesAreCertified}
	expectedChecks[checks.ChartTestingName] = checks.Check{CheckId: checks.CheckId{Name: checks.ChartTestingName, Version: "v1.0"}, Type: checks.MandatoryCheckType, Func: checks.ChartTesting}

	t.Run("Checks filtered using profile - full set", func(t *testing.T) {
		filteredChecks := GetProfile().FilterChecks(defaultRegistry.AllChecks())
		compareCheckMaps(t, expectedChecks, filteredChecks)
	})

}

func compareCheckMaps(t *testing.T, expectedChecks, filteredChecks map[checks.CheckName]checks.Check) {

	assert.Equal(t, len(expectedChecks), len(filteredChecks), fmt.Sprintf("Expected map length : %d does not match returned mao length : %d", len(expectedChecks), len(filteredChecks)))
	for k, v := range filteredChecks {
		_, ok := expectedChecks[k]
		if !ok {
			assert.True(t, ok, "Entry not found in expected: %s", k)
		} else {
			assert.Equal(t, v.CheckId.Name, expectedChecks[k].CheckId.Name, fmt.Sprintf("%s: Map names do not match! got:%s, expect:%s ", k, v.CheckId.Name, expectedChecks[k].CheckId.Name))
			assert.Equal(t, v.CheckId.Version, expectedChecks[k].CheckId.Version, fmt.Sprintf("%s: Map versions do not match! got:%s, expect:%s", k, v.CheckId.Version, expectedChecks[k].CheckId.Version))
			assert.Equal(t, v.Type, expectedChecks[k].Type, fmt.Sprintf("%s: Map types do not match! got:%s, expect:%s", k, v.Type, expectedChecks[k].Type))
			runFunc := filepath.Base(runtime.FuncForPC(reflect.ValueOf(v.Func).Pointer()).Name())
			expectFunc := filepath.Base(runtime.FuncForPC(reflect.ValueOf(expectedChecks[k].Func).Pointer()).Name())
			assert.Equal(t, runFunc, expectFunc, fmt.Sprintf("%s: Map funcs do not match! got:%v, expect:%v", k, runFunc, expectFunc))
		}
	}
}
