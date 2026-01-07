package profiles

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"github.com/redhat-certification/chart-verifier/internal/chartverifier/checks"
	apiChecks "github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
)

const (
	NoVersion       string = ""
	configVersion00 string = "v0.0"
	// TODO(komish) identify where this was intended to be used, or remove.
	configVersion10     string     = "v1.0"
	configVersion11     string     = "v1.1"
	configVersion12     string     = "v1.2"
	configVersion13     string     = "v1.3"
	checkVersion10      string     = CheckVersion10
	checkVersion11      string     = "v1.1"
	NoVendorType        VendorType = ""
	PartnerVendorType   VendorType = "partner"
	RedhatVendorType    VendorType = "redhat"
	CommunityVendorType VendorType = "community"
)

func TestProfile(t *testing.T) {
	testProfile := getDefaultProfile("test")
	testProfile.Name = "profile-partner-1.3"
	config := make(map[string]interface{})
	config[VendorTypeConfigName] = PartnerVendorType

	t.Run("Profile read from disk should match test profile", func(t *testing.T) {
		diskProfile := New(config)
		if !cmp.Equal(diskProfile, testProfile) {
			assert.Equal(t, testProfile.Name, diskProfile.Name, "Name mismatch")
			assert.Equal(t, testProfile.Vendor, diskProfile.Vendor, "Vendor mismatch")
			assert.Equal(t, testProfile.Version, diskProfile.Version, "Version mismatch")
			assert.Equal(t, len(testProfile.Annotations), len(diskProfile.Annotations), "Annotations number mismatch")
			for _, testAnnotation := range testProfile.Annotations {
				found := slices.Contains(diskProfile.Annotations, testAnnotation)
				assert.True(t, found, fmt.Sprintf("Annotation not found : %s", testAnnotation))
			}
			assert.Equal(t, len(testProfile.Checks), len(diskProfile.Checks), "Checks number mismatch")
			for _, testCheck := range testProfile.Checks {
				found := slices.Contains(diskProfile.Checks, testCheck)
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
	getAndCheckProfile(t, RedhatVendorType, RedhatVendorType, NoVersion, configVersion13)
	getAndCheckProfile(t, NoVendorType, PartnerVendorType, NoVersion, configVersion13)
	getAndCheckProfile(t, PartnerVendorType, PartnerVendorType, configVersion13, configVersion13)
	getAndCheckProfile(t, PartnerVendorType, PartnerVendorType, configVersion00, configVersion13)
	getAndCheckProfile(t, RedhatVendorType, RedhatVendorType, configVersion13, configVersion13)
	getAndCheckProfile(t, RedhatVendorType, RedhatVendorType, configVersion00, configVersion13)
	getAndCheckProfile(t, CommunityVendorType, CommunityVendorType, configVersion00, configVersion13)
	getAndCheckProfile(t, CommunityVendorType, CommunityVendorType, configVersion13, configVersion13)
}

func getAndCheckProfile(t *testing.T, configVendorType, expectVendorType VendorType, configVersion, expectVersion string) {
	config := make(map[string]interface{})

	if len(configVendorType) > 0 {
		config[VendorTypeConfigName] = configVendorType
	}
	if len(configVersion) > 0 {
		config[VersionConfigName] = configVersion
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

	defaultRegistry.Add(apiChecks.HasReadme, checkVersion10, checks.HasReadme)
	defaultRegistry.Add(apiChecks.IsHelmV3, checkVersion10, checks.IsHelmV3)
	defaultRegistry.Add(apiChecks.ContainsTest, checkVersion10, checks.ContainsTest)
	defaultRegistry.Add(apiChecks.ContainsValues, checkVersion10, checks.ContainsValues)

	defaultRegistry.Add(apiChecks.HasReadme, checkVersion11, checks.HasReadme)
	defaultRegistry.Add(apiChecks.IsHelmV3, checkVersion11, checks.IsHelmV3)
	defaultRegistry.Add(apiChecks.ContainsTest, checkVersion11, checks.ContainsTest)
	defaultRegistry.Add(apiChecks.ContainsValues, checkVersion11, checks.ContainsValues)
	defaultRegistry.Add(apiChecks.ContainsValuesSchema, checkVersion11, checks.ContainsValuesSchema)
	defaultRegistry.Add(apiChecks.HasKubeVersion, checkVersion11, checks.HasKubeVersion_V1_1)
	defaultRegistry.Add(apiChecks.NotContainsCRDs, checkVersion11, checks.NotContainCRDs)
	defaultRegistry.Add(apiChecks.HelmLint, checkVersion11, checks.HelmLint)
	defaultRegistry.Add(apiChecks.NotContainsCRDs, checkVersion11, checks.NotContainCSIObjects)
	defaultRegistry.Add(apiChecks.ImagesAreCertified, checkVersion11, checks.ImagesAreCertified_V1_1)
	defaultRegistry.Add(apiChecks.ChartTesting, checkVersion11, checks.ChartTesting)

	defaultRegistry.Add("BadHasReadme", checkVersion10, checks.HasReadme)
	defaultRegistry.Add("BadIsHelmV3Name", checkVersion10, checks.IsHelmV3)
	defaultRegistry.Add("BadContainsTestName", "v1.o", checks.ContainsTest)

	expectedChecks := FilteredRegistry{}
	expectedChecks[apiChecks.HasReadme] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.HasReadme, Version: checkVersion10}, Type: apiChecks.MandatoryCheckType, Func: checks.HasReadme}
	expectedChecks[apiChecks.IsHelmV3] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.IsHelmV3, Version: checkVersion10}, Type: apiChecks.MandatoryCheckType, Func: checks.IsHelmV3}
	expectedChecks[apiChecks.ContainsTest] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.ContainsTest, Version: checkVersion10}, Type: apiChecks.MandatoryCheckType, Func: checks.ContainsTest}
	expectedChecks[apiChecks.ContainsValues] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.ContainsValues, Version: checkVersion10}, Type: apiChecks.MandatoryCheckType, Func: checks.ContainsValues}
	expectedChecks[apiChecks.HasKubeVersion] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.HasKubeVersion, Version: checkVersion11}, Type: apiChecks.MandatoryCheckType, Func: checks.HasKubeVersion_V1_1}
	expectedChecks[apiChecks.ImagesAreCertified] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.ImagesAreCertified, Version: checkVersion11}, Type: apiChecks.MandatoryCheckType, Func: checks.ImagesAreCertified_V1_1}

	config := make(map[string]interface{})
	t.Run("Checks filtered using profile subset", func(t *testing.T) {
		filteredChecks := New(config).FilterChecks(defaultRegistry.AllChecks())
		CompareCheckMaps(t, expectedChecks, filteredChecks)
	})

	defaultRegistry.Add(apiChecks.ContainsValuesSchema, checkVersion10, checks.ContainsValuesSchema)
	defaultRegistry.Add(apiChecks.HasKubeVersion, checkVersion10, checks.HasKubeVersion)
	defaultRegistry.Add(apiChecks.NotContainsCRDs, checkVersion10, checks.NotContainCRDs)
	defaultRegistry.Add(apiChecks.HelmLint, checkVersion10, checks.HelmLint)
	defaultRegistry.Add(apiChecks.NotContainCsiObjects, checkVersion10, checks.NotContainCSIObjects)
	defaultRegistry.Add(apiChecks.ImagesAreCertified, checkVersion10, checks.ImagesAreCertified)
	defaultRegistry.Add(apiChecks.ChartTesting, checkVersion10, checks.ChartTesting)
	defaultRegistry.Add(apiChecks.RequiredAnnotationsPresent, checkVersion10, checks.RequiredAnnotationsPresent)

	expectedChecks[apiChecks.ContainsValuesSchema] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.ContainsValuesSchema, Version: checkVersion10}, Type: apiChecks.MandatoryCheckType, Func: checks.ContainsValuesSchema}
	expectedChecks[apiChecks.NotContainsCRDs] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.NotContainsCRDs, Version: checkVersion10}, Type: apiChecks.MandatoryCheckType, Func: checks.NotContainCRDs}
	expectedChecks[apiChecks.HelmLint] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.HelmLint, Version: checkVersion10}, Type: apiChecks.MandatoryCheckType, Func: checks.HelmLint}
	expectedChecks[apiChecks.NotContainCsiObjects] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.NotContainCsiObjects, Version: checkVersion10}, Type: apiChecks.MandatoryCheckType, Func: checks.NotContainCSIObjects}
	expectedChecks[apiChecks.ChartTesting] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.ChartTesting, Version: checkVersion10}, Type: apiChecks.MandatoryCheckType, Func: checks.ChartTesting}
	expectedChecks[apiChecks.RequiredAnnotationsPresent] = checks.Check{CheckID: checks.CheckID{Name: apiChecks.RequiredAnnotationsPresent, Version: checkVersion10}, Type: apiChecks.MandatoryCheckType, Func: checks.RequiredAnnotationsPresent}

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
			assert.Equal(t, v.CheckID.Name, expectedChecks[k].CheckID.Name, fmt.Sprintf("%s: Map names do not match! got:%s, expect:%s ", k, v.CheckID.Name, expectedChecks[k].CheckID.Name))
			assert.Equal(t, v.CheckID.Version, expectedChecks[k].CheckID.Version, fmt.Sprintf("%s: Map versions do not match! got:%s, expect:%s", k, v.CheckID.Version, expectedChecks[k].CheckID.Version))
			assert.Equal(t, v.Type, expectedChecks[k].Type, fmt.Sprintf("%s: Map types do not match! got:%s, expect:%s", k, v.Type, expectedChecks[k].Type))
			runFunc := filepath.Base(runtime.FuncForPC(reflect.ValueOf(v.Func).Pointer()).Name())
			expectFunc := filepath.Base(runtime.FuncForPC(reflect.ValueOf(expectedChecks[k].Func).Pointer()).Name())
			assert.Equal(t, runFunc, expectFunc, fmt.Sprintf("%s: Map funcs do not match! got:%v, expect:%v", k, runFunc, expectFunc))
		}
	}
}
