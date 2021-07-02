package report

import (
	"github.com/spf13/viper"
)

const (
	MetadataCommandName    string = "metadata"
	DigestsCommandName     string = "digests"
	ResultsCommandName     string = "results"
	AnnotationsCommandName string = "annotations"
	AllCommandsName        string = "all"
)

var reportCommandRegistry ReportRegistry

func init() {
	reportCommandRegistry = NewReportRegistry()

	reportCommandRegistry.Add(MetadataCommandName, Metadata)
	reportCommandRegistry.Add(DigestsCommandName, Digests)
	reportCommandRegistry.Add(ResultsCommandName, Results)
	reportCommandRegistry.Add(AnnotationsCommandName, Annotations)
	reportCommandRegistry.Add(AllCommandsName, All)

}

func ReportCommandRegistry() ReportRegistry {
	return reportCommandRegistry
}

type ReportFunc func(options *ReportOptions) (OutputReport, error)

type ReportOptions struct {
	// URI is the location of the chart to be checked.
	URI string
	// ViperConfig is the configuration collected by Viper.
	ViperConfig *viper.Viper
	// Values contains the values informed by the user through command line options.
	Values map[string]interface{}
}

type ReportRegistry interface {
	Get(name string) ReportFunc
	Add(name string, reportFunc ReportFunc) ReportRegistry
	AllChecks() DefaultReportRegistry
}

type DefaultReportRegistry map[string]ReportFunc

func NewReportRegistry() ReportRegistry {
	return &DefaultReportRegistry{}
}

func (r *DefaultReportRegistry) AllChecks() DefaultReportRegistry {
	return *r
}

func (r *DefaultReportRegistry) Get(name string) ReportFunc {
	v, _ := (*r)[name]
	return v
}

func (r *DefaultReportRegistry) Add(name string, reportFunc ReportFunc) ReportRegistry {
	(*r)[name] = reportFunc
	return r

}
