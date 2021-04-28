package checks

import (
	"github.com/helm/chart-testing/pkg/chart"
	"github.com/helm/chart-testing/pkg/config"
)

func ChartTesting(opts *CheckOptions) (Result, error) {

	cfg := config.Configuration{}

	cfg.Namespace = opts.HelmEnvSettings.Namespace()

	testing := chart.NewTesting(cfg)

	results, err := testing.InstallCharts()
	if err != nil {
		return Result{}, err
	}

	// aggregate errors
	for _, r := range results {
		if r.Error != nil {
			return NewResult(false, r.Error.Error()), nil
		}
	}

	return NewResult(true, ""), nil
}
