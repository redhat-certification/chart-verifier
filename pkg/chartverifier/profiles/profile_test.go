package profiles

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	NoVersion           string     = ""
	configVersion00     string     = "v0.0"
	configVersion10     string     = "v1.0"
	configVersion11     string     = "v1.1"
	configVersion12     string     = "v1.2"
	checkVersion10      string     = CheckVersion10
	checkVersion11      string     = "v1.1"
	NoVendorType        VendorType = ""
	PartnerVendorType   VendorType = "partner"
	RedhatVendorType    VendorType = "redhat"
	CommunityVendorType VendorType = "community"
)

func TestProfile(t *testing.T) {

	testProfile := getDefaultProfile("test")
	testProfile.Name = "profile-partner-1.1"
	config := viper.New()
	config.Set(VendorTypeConfigName, string(PartnerVendorType))

	t.Run("Profile read from disk should match test profile", func(t *testing.T) {
		diskProfile := New(config)
		if !cmp.Equal(diskProfile, testProfile) {
			assert.Equal(t, testProfile.Name, diskProfile.Name, "Name mismatch")
			assert.Equal(t, testProfile.Vendor, diskProfile.Vendor, "Vendor mismatch")
			assert.Equal(t, testProfile.Version, diskProfile.Version, "Version mismatch")
			assert.Equal(t, len(testProfile.Annotations), len(diskProfile.Annotations), "Annotations number mismatch")
			for _, testAnnotation := range testProfile.Annotations {
				found := false
				for _, diskAnnotation := range diskProfile.Annotations {
					if testAnnotation == diskAnnotation {
						found = true
						break
					}
				}
				assert.True(t, found, fmt.Sprintf("Annotation not found : %s", testAnnotation))
			}
			assert.Equal(t, len(testProfile.Checks), len(diskProfile.Checks), "Checks number mismatch")
			for _, testCheck := range testProfile.Checks {
				found := false
				for _, diskCheck := range diskProfile.Checks {
					if strings.Compare(testCheck.Name, diskCheck.Name) == 0 {
						if testCheck.Type == diskCheck.Type {
							found = true
							break
						}
					}
				}
				assert.True(t, found, fmt.Sprintf("Check not matched : %s : %s", testCheck.Name, testCheck.Type))
			}
			assert.True(t, cmp.Equal(diskProfile, testProfile), "profiles do not match")
		}
	})

}

func TestGetProfiles(t *testing.T) {

	getAndCheckProfile(t, PartnerVendorType, PartnerVendorType, configVersion11, configVersion11)
	getAndCheckProfile(t, RedhatVendorType, RedhatVendorType, configVersion11, configVersion11)
	getAndCheckProfile(t, CommunityVendorType, CommunityVendorType, configVersion11, configVersion11)
	getAndCheckProfile(t, NoVendorType, PartnerVendorType, configVersion11, configVersion11)
	getAndCheckProfile(t, RedhatVendorType, RedhatVendorType, NoVersion, configVersion11)
	getAndCheckProfile(t, NoVendorType, PartnerVendorType, NoVersion, configVersion11)
	getAndCheckProfile(t, PartnerVendorType, PartnerVendorType, configVersion12, configVersion11)
	getAndCheckProfile(t, PartnerVendorType, PartnerVendorType, configVersion00, configVersion11)
	getAndCheckProfile(t, RedhatVendorType, RedhatVendorType, configVersion12, configVersion11)
	getAndCheckProfile(t, RedhatVendorType, RedhatVendorType, configVersion00, configVersion11)
	getAndCheckProfile(t, CommunityVendorType, CommunityVendorType, configVersion00, configVersion11)
	getAndCheckProfile(t, CommunityVendorType, CommunityVendorType, configVersion12, configVersion11)
}

func getAndCheckProfile(t *testing.T, configVendorType, expectVendorType VendorType, configVersion, expectVersion string) {

	config := viper.New()
	if len(configVendorType) > 0 {
		config.Set(VendorTypeConfigName, string(configVendorType))
	}
	if len(configVersion) > 0 {
		config.Set(VersionConfigName, configVersion)
	}

	t.Run(fmt.Sprintf("Request : VendorType config %s expect %s : Version config %s expect %s ", configVendorType, expectVendorType, configVersion, expectVersion), func(t *testing.T) {
		profile := New(config)
		assert.Equal(t, expectVendorType, profile.Vendor, "VendorType did not match")
		assert.Equal(t, expectVersion, profile.Version, "Version did not match")
		profile = Get()
		assert.Equal(t, expectVendorType, profile.Vendor, "VendorType did not match")
		assert.Equal(t, expectVersion, profile.Version, "Version did not match")

	})
}
func TestProfileFilter(t *testing.T) {

	defaultRegistry := checks.NewRegistry()

	defaultRegistry.Add(checks.HasReadmeName, checkVersion10, checks.HasReadme)
	defaultRegistry.Add(checks.IsHelmV3Name, checkVersion10, checks.IsHelmV3)
	defaultRegistry.Add(checks.ContainsTestName, checkVersion10, checks.ContainsTest)
	defaultRegistry.Add(checks.ContainsValuesName, checkVersion10, checks.ContainsValues)

	defaultRegistry.Add(checks.HasReadmeName, checkVersion11, checks.HasReadme)
	defaultRegistry.Add(checks.IsHelmV3Name, checkVersion11, checks.IsHelmV3)
	defaultRegistry.Add(checks.ContainsTestName, checkVersion11, checks.ContainsTest)
	defaultRegistry.Add(checks.ContainsValuesName, checkVersion11, checks.ContainsValues)
	defaultRegistry.Add(checks.ContainsValuesSchemaName, checkVersion11, checks.ContainsValuesSchema)
	defaultRegistry.Add(checks.HasKubeversionName, checkVersion11, checks.HasKubeVersion_V1_1)
	defaultRegistry.Add(checks.NotContainsCRDsName, checkVersion11, checks.NotContainCRDs)
	defaultRegistry.Add(checks.HelmLintName, checkVersion11, checks.HelmLint)
	defaultRegistry.Add(checks.NotContainsCRDsName, checkVersion11, checks.NotContainCSIObjects)
	defaultRegistry.Add(checks.ImagesAreCertifiedName, checkVersion11, checks.ImagesAreCertified)
	defaultRegistry.Add(checks.ChartTestingName, checkVersion11, checks.ChartTesting)

	defaultRegistry.Add("BadHasReadme", checkVersion10, checks.HasReadme)
	defaultRegistry.Add("BadIsHelmV3Name", checkVersion10, checks.IsHelmV3)
	defaultRegistry.Add("BadContainsTestName", "v1.o", checks.ContainsTest)

	expectedChecks := FilteredRegistry{}
	expectedChecks[checks.HasReadmeName] = checks.Check{CheckId: checks.CheckId{Name: checks.HasReadmeName, Version: checkVersion10}, Type: checks.MandatoryCheckType, Func: checks.HasReadme}
	expectedChecks[checks.IsHelmV3Name] = checks.Check{CheckId: checks.CheckId{Name: checks.IsHelmV3Name, Version: checkVersion10}, Type: checks.MandatoryCheckType, Func: checks.IsHelmV3}
	expectedChecks[checks.ContainsTestName] = checks.Check{CheckId: checks.CheckId{Name: checks.ContainsTestName, Version: checkVersion10}, Type: checks.MandatoryCheckType, Func: checks.ContainsTest}
	expectedChecks[checks.ContainsValuesName] = checks.Check{CheckId: checks.CheckId{Name: checks.ContainsValuesName, Version: checkVersion10}, Type: checks.MandatoryCheckType, Func: checks.ContainsValues}
	expectedChecks[checks.HasKubeversionName] = checks.Check{CheckId: checks.CheckId{Name: checks.HasKubeversionName, Version: checkVersion11}, Type: checks.MandatoryCheckType, Func: checks.HasKubeVersion_V1_1}

	config := viper.New()
	t.Run("Checks filtered using profile subset", func(t *testing.T) {
		filteredChecks := New(config).FilterChecks(defaultRegistry.AllChecks())
		CompareCheckMaps(t, expectedChecks, filteredChecks)
	})

	defaultRegistry.Add(checks.ContainsValuesSchemaName, checkVersion10, checks.ContainsValuesSchema)
	defaultRegistry.Add(checks.HasKubeversionName, checkVersion10, checks.HasKubeVersion)
	defaultRegistry.Add(checks.NotContainsCRDsName, checkVersion10, checks.NotContainCRDs)
	defaultRegistry.Add(checks.HelmLintName, checkVersion10, checks.HelmLint)
	defaultRegistry.Add(checks.NotContainCsiObjectsName, checkVersion10, checks.NotContainCSIObjects)
	defaultRegistry.Add(checks.ImagesAreCertifiedName, checkVersion10, checks.ImagesAreCertified)
	defaultRegistry.Add(checks.ChartTestingName, checkVersion10, checks.ChartTesting)
	defaultRegistry.Add(checks.RequiredAnnotationsPresentName, checkVersion10, checks.RequiredAnnotationsPresent)

	expectedChecks[checks.ContainsValuesSchemaName] = checks.Check{CheckId: checks.CheckId{Name: checks.ContainsValuesSchemaName, Version: checkVersion10}, Type: checks.MandatoryCheckType, Func: checks.ContainsValuesSchema}
	expectedChecks[checks.NotContainsCRDsName] = checks.Check{CheckId: checks.CheckId{Name: checks.NotContainsCRDsName, Version: checkVersion10}, Type: checks.MandatoryCheckType, Func: checks.NotContainCRDs}
	expectedChecks[checks.HelmLintName] = checks.Check{CheckId: checks.CheckId{Name: checks.HelmLintName, Version: checkVersion10}, Type: checks.MandatoryCheckType, Func: checks.HelmLint}
	expectedChecks[checks.NotContainCsiObjectsName] = checks.Check{CheckId: checks.CheckId{Name: checks.NotContainCsiObjectsName, Version: checkVersion10}, Type: checks.MandatoryCheckType, Func: checks.NotContainCSIObjects}
	expectedChecks[checks.ImagesAreCertifiedName] = checks.Check{CheckId: checks.CheckId{Name: checks.ImagesAreCertifiedName, Version: checkVersion10}, Type: checks.MandatoryCheckType, Func: checks.ImagesAreCertified}
	expectedChecks[checks.ChartTestingName] = checks.Check{CheckId: checks.CheckId{Name: checks.ChartTestingName, Version: checkVersion10}, Type: checks.MandatoryCheckType, Func: checks.ChartTesting}
	expectedChecks[checks.RequiredAnnotationsPresentName] = checks.Check{CheckId: checks.CheckId{Name: checks.RequiredAnnotationsPresentName, Version: checkVersion10}, Type: checks.MandatoryCheckType, Func: checks.RequiredAnnotationsPresent}

	t.Run("Checks filtered using profile - full set", func(t *testing.T) {
		filteredChecks := New(config).FilterChecks(defaultRegistry.AllChecks())
		CompareCheckMaps(t, expectedChecks, filteredChecks)
	})

}

func CompareCheckMaps(t *testing.T, expectedChecks, filteredChecks FilteredRegistry) {

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
