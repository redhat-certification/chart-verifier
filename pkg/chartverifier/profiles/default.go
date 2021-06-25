package profiles

import (
	"fmt"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
)

const (
	CheckVersion10 = "v1.0"
)

func getDefaultProfile(msg string) *Profile {
	profile := Profile{}

	profile.Apiversion = "v1"
	profile.Kind = "verifier-profile"

	profile.Vendor = PartnerVendorType
	profile.Version = "1.1"

	profile.Annotations = []Annotation{DigestAnnotation, OCPVersionAnnotation, LastCertifiedTimestampAnnotation}

	profile.Checks = []*Check{
		{Name: fmt.Sprintf("%s/%s", CheckVersion10, checks.HasReadmeName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", CheckVersion10, checks.IsHelmV3Name), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", CheckVersion10, checks.ContainsTestName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", CheckVersion10, checks.ContainsValuesName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", CheckVersion10, checks.ContainsValuesSchemaName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", CheckVersion10, checks.HasKubeversionName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", CheckVersion10, checks.NotContainsCRDsName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", CheckVersion10, checks.HelmLintName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", CheckVersion10, checks.NotContainCsiObjectsName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", CheckVersion10, checks.ImagesAreCertifiedName), Type: checks.MandatoryCheckType},
		{Name: fmt.Sprintf("%s/%s", CheckVersion10, checks.ChartTestingName), Type: checks.MandatoryCheckType},
	}

	profile.Name = fmt.Sprintf("default-profile : %s", msg)

	return &profile
}
