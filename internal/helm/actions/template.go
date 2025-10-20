package actions

import (
	"bytes"
	"fmt"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

func RenderManifests(name string, url string, vals map[string]interface{}, conf *action.Configuration) (string, error) {
	validate := false
	client := action.NewInstall(conf)
	client.DryRun = true
	includeCrds := true
	client.ReleaseName = "RELEASE-NAME"
	client.Replace = true // Skip the releaseName check
	client.ClientOnly = !validate
	emptyResponse := ""

	// Must set the capabilities on *action.Install{}.KubeVersion directly
	// because Helm will replace our capabilities with the defaults they
	// configure.
	if conf.Capabilities != nil {
		client.KubeVersion = &conf.Capabilities.KubeVersion
	}

	name, chart, err := client.NameAndChart([]string{name, url})
	if err != nil {
		return emptyResponse, err
	}
	client.ReleaseName = name

	cp, err := client.LocateChart(chart, cli.New())
	if err != nil {
		return emptyResponse, err
	}

	ch, err := loader.Load(cp)
	if err != nil {
		return emptyResponse, err
	}

	rel, err := client.Run(ch, vals)
	if err != nil {
		return emptyResponse, err
	}

	var manifests bytes.Buffer
	var output bytes.Buffer

	if includeCrds {
		for _, f := range rel.Chart.CRDs() {
			fmt.Fprintf(&manifests, "---\n# Source: %s\n%s\n", f.Name, f.Data)
		}
	}

	fmt.Fprintln(&manifests, strings.TrimSpace(rel.Manifest))

	if !client.DisableHooks {
		for _, m := range rel.Hooks {
			fmt.Fprintf(&manifests, "---\n# Source: %s\n%s\n", m.Path, m.Manifest)
		}
	}

	fmt.Fprintf(&output, "%s", manifests.String())
	return output.String(), nil
}
