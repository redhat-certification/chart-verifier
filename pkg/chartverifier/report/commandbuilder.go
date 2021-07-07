package report

import (
	"github.com/spf13/viper"
	"strings"
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
}

func (opt *ReportOptions) AddURI(uri string) {
	opt.URI = uri
}

func (opt *ReportOptions) AddConfig(config *viper.Viper) {
	opt.ViperConfig = config
}

func (opt *ReportOptions) AddValues(values []string) {
	// naively override values from the configuration
	for _, val := range values {
		parts := strings.Split(val, "=")
		opt.ViperConfig.Set(parts[0], parts[1])
	}
}

type ReportRegistry interface {
	Get(name string) ReportFunc
	Add(name string, reportFunc ReportFunc) ReportRegistry
	AllCommands() DefaultReportRegistry
}

type DefaultReportRegistry map[string]ReportFunc

func NewReportRegistry() ReportRegistry {
	return &DefaultReportRegistry{}
}

func (r *DefaultReportRegistry) AllCommands() DefaultReportRegistry {
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
