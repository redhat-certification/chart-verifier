package chartverifier

import (
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"os"
	"path/filepath"
)

type ReportBuilder interface {
	SetCertificate(*Certificate) ReportBuilder
	SetChartUri(string) ReportBuilder
	AddChartYaml(*chart.File) ReportBuilder
	Generate() error
}

type reportBuilder struct {
	Certificate *Certificate
	ChartUri    string
	ChartYaml   *chart.File
}

type chartYamlMetadata struct {
	ChartYamlMetadata interface{} `json:"chart-yaml" yaml:"chart-yaml"`
}

func newChartYamlMetadata(data interface{}) *chartYamlMetadata {
	return &chartYamlMetadata{
		ChartYamlMetadata: data,
	}
}

func NewReportBuilder() ReportBuilder {
	return &reportBuilder{}
}

func (r *reportBuilder) SetCertificate(cert *Certificate) ReportBuilder {
	r.Certificate = cert
	return r
}

func (r *reportBuilder) SetChartUri(uri string) ReportBuilder {
	r.ChartUri = uri
	return r
}

func (r *reportBuilder) AddChartYaml(file *chart.File) ReportBuilder {
	r.ChartYaml = file
	return r
}

func (r *reportBuilder) Generate() error {

	var err error
	reportsDir := "reports"
	if _, err = os.Stat(reportsDir); os.IsNotExist(err) {
		err = os.Mkdir(reportsDir, 0755)
	}
	if err != nil {
		return err
	}

	reportDir := filepath.FromSlash(reportsDir + "/" + filepath.Base(r.ChartUri))

	if _, err = os.Stat(reportDir); !os.IsNotExist(err) {
		os.RemoveAll(reportDir)
	}
	if err = os.Mkdir(reportDir, 0755); err != nil {
		return err
	}

	f, err := os.Create(reportDir + "/verifier.report.yaml")
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(r.Certificate)
	if err != nil {
		return err
	}

	f.Write(b)

	c, _, err := checks.LoadChartFromURI(r.ChartUri)
	for _, chartFile := range c.Raw {
		if chartFile.Name == "Chart.yaml" {
			var unmarshalledChart interface{}
			err = yaml.Unmarshal(chartFile.Data, &unmarshalledChart)
			if err != nil {
				return err
			}
			b, err = yaml.Marshal(newChartYamlMetadata(unmarshalledChart))
			if err != nil {
				return err
			}
			f.Write(b)
			break
		}
	}

	f.Close()

	return err
}
