package profiles

import (
	"github.com/google/go-cmp/cmp"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProfile(t *testing.T) {

	testProfile := Profile{}

	testProfile.Apiversion = "v1"
	testProfile.Kind = "verifier-profile"
	testProfile.Name = "profile-1.0.0"

	testProfile.Annotations = []Annotation{DigestAnnotation, OCPVersionAnnotation, LastCertifiedTimestampAnnotation}

	testProfile.Checks = []*Check{
		{Name: checks.IsHelmV3Name, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.HasReadmeName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.ContainsTestName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.ContainsValuesName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.ContainsValuesSchemaName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.HasKubeversionName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.NotContainsCRDsName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.HelmLintName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.NotContainCsiObjectsName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.ImagesAreCertifiedName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.ChartTestingName, Type: checks.MandatoryCheckType, Version: "1.0"},
	}

	t.Run("Profile read from disk should match test profile", func(t *testing.T) {

		diskProfile := GetProfile()
		assert.True(t, cmp.Equal(diskProfile, &testProfile), "profiles do not match")

	})

}
