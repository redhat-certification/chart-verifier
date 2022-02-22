package report

import (
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/profiles"
	helmchart "helm.sh/helm/v3/pkg/chart"
)

type OutputReport struct {
	AnnotationsReport []Annotation    `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	DigestsReport     *DigestReport   `json:"digests,omitempty" yaml:"digests,omitempty"`
	MetadataReport    *MetadataReport `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	ResultsReport     *ResultsReport  `json:"results,omitempty" yaml:"results,omitempty"`
}

type Annotation struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type DigestReport struct {
	ChartDigest   string `json:"chart" yaml:"chart"`
	PackageDigest string `json:"package" yaml:"package"`
}

type MetadataReport struct {
	ProfileVendorType profiles.VendorType `json:"vendorType" yaml:"vendorType"`
	ProfileVersion    string              `json:"profileVersion" yaml:"profileVersion"`
	ProviderDelivery  bool                `json:"providerControlledDelivery" yaml:"providerControlledDelivery,omitempty"`
	ChartUri          string              `json:"chart-uri" yaml:"chart-uri"`
	Chart             *helmchart.Metadata `json:"chart" yaml:"chart"`
}

type ResultsReport struct {
	Passed   string   `json:"passed" yaml:"passed"`
	Failed   string   `json:"failed" yaml:"failed"`
	Messages []string `json:"message" yaml:"message"`
}
