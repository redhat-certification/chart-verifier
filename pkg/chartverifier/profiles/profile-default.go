package profiles

import (
	"fmt"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
)

func getDefaultProfile(msg string) *Profile {
	profile := Profile{}

	profile.Name = fmt.Sprintf("default-profile : %s", msg)

	profile.Annotations = []Annotation{DigestAnnotation, OCPVersionAnnotation, LastCertifiedTimestampAnnotation}

	profile.Checks = []*Check{
		{Name: checks.HasReadmeName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.IsHelmV3Name, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.ContainsTestName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.ContainsValuesName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.ContainsValuesName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.HasKubeversionName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.NotContainsCRDsName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.HelmLintName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.NotContainCsiObjectsName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.ImagesAreCertifiedName, Type: checks.MandatoryCheckType, Version: "1.0"},
		{Name: checks.ChartTestingName, Type: checks.MandatoryCheckType, Version: "1.0"},
	}

	return &profile
}
