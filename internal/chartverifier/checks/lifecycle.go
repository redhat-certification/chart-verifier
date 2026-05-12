package checks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	ApiVersion   = "v1"
	Api          = "products"
	ApiEndpoint  = "access.redhat.com/product-life-cycles/api/" + ApiVersion + "/" + Api
	SearchString = "Openshift Container Platform"
)

type LifecycleData struct {
	UUID                  string `json:"uuid"`
	Name                  string `json:"name"`
	FormerNames           []any  `json:"former_names"`
	ShowLastMinorRelease  bool   `json:"show_last_minor_release"`
	ShowFinalMinorRelease bool   `json:"show_final_minor_release"`
	IsLayeredProduct      bool   `json:"is_layered_product"`
	AllPhases             []struct {
		Name           string `json:"name"`
		Ptype          string `json:"ptype"`
		Tooltip        any    `json:"tooltip"`
		DisplayName    string `json:"display_name"`
		AdditionalText string `json:"additional_text"`
	} `json:"all_phases"`
	Versions []struct {
		Name              string `json:"name"`
		Type              string `json:"type"`
		LastMinorRelease  any    `json:"last_minor_release"`
		FinalMinorRelease any    `json:"final_minor_release"`
		ExtraHeaderValue  any    `json:"extra_header_value"`
		AdditionalText    string `json:"additional_text"`
		Phases            []struct {
			Name           string `json:"name"`
			Date           string `json:"date"`
			DateFormat     string `json:"date_format"`
			AdditionalText string `json:"additional_text"`
		} `json:"phases"`
		Tier                   string `json:"tier"`
		OpenshiftCompatibility string `json:"openshift_compatibility"`
		ExtraDependences       []any  `json:"extra_dependences"`
	} `json:"versions"`
	Footnote                   string `json:"footnote"`
	IsOperator                 bool   `json:"is_operator"`
	ShowOpenshiftCompatibility bool   `json:"show_openshift_compatibility"`
	ReleaseCadence             string `json:"release_cadence"`
	Package                    string `json:"package"`
	Link                       string `json:"link"`
	Policies                   string `json:"policies"`
}

type Results struct {
	Data []LifecycleData `json:"data"`
}

type LifecycleDataGetter func(clusterVersion string) (string, error)

func (*LifecycleData) GetLifecycleStatus(clusterVersion string) (string, error) {
	// Get lifecycle data
	requestURL := fmt.Sprintf("https://%s?name=%s", ApiEndpoint, url.QueryEscape(SearchString))
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("Error getting request: %s\n", err)
		return "Unknown", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading body: %s\n", err)
	}

	// Filter to get lifecycle status
	var results Results
	err = json.Unmarshal(body, &results)
	if err != nil {
		fmt.Printf("Error unmarshalling json: %s\n", err.Error())
	}
	for _, product := range results.Data {
		for _, version := range product.Versions {
			if strings.Compare(version.Name, clusterVersion) == 0 {
				return version.Type, nil
			}
		}
	}

	return "Unknown", nil
}
