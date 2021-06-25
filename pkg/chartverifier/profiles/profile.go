package profiles

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io"

	//"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type Annotation string
type VendorType string

const (
	DigestAnnotation                 Annotation = "Digest"
	OCPVersionAnnotation             Annotation = "OCPVersion"
	LastCertifiedTimestampAnnotation Annotation = "LastCertifiedTimestamp"

	PartnerVendorType   VendorType = "partner"
	CommunityVendorType VendorType = "community"
	RedhatVendorType    VendorType = "redhat"
	DefaultVendorType              = PartnerVendorType

	VendorTypeConfigName string = "verifier.vendortype"
	VersionConfigName    string = "verifier.version"
)

type Profile struct {
	Apiversion  string       `json:"apiversion" yaml:"apiversion"`
	Kind        string       `json:"kind" yaml:"kind"`
	Name        string       `json:"name" yaml:"name"`
	Vendor      VendorType   `json:"vendorType" yaml:"vendorType"`
	Version     string       `json:"version" yaml:"version"`
	Annotations []Annotation `json:"annotations" yaml:"annotations"`
	Checks      []*Check     `json:"checks" yaml:"checks"`
}

type Check struct {
	Name string           `json:"name" yaml:"name"`
	Type checks.CheckType `json:"type" yaml:"type"`
}

type FilteredRegistry map[checks.CheckName]checks.Check

var profile *Profile

func Get() *Profile {
	if profile == nil {
		return getDefaultProfile("No profile set for get")
	}
	return profile
}

func New(version string, config *viper.Viper) *Profile {

	profileVersion, _ := semver.NewVersion(version)
	profileVendorType := DefaultVendorType

	if config != nil {
		configVersion := config.GetString(VersionConfigName)
		if len(configVersion) > 0 {
			requestedVersion, err := semver.NewVersion(configVersion)
			if err != nil {
				if !requestedVersion.GreaterThan(profileVersion) {
					profileVersion = requestedVersion
				}
			}
		}
		configVendorType := config.GetString(VendorTypeConfigName)
		if len(configVendorType) > 0 {
			switch VendorType(configVendorType) {
			case PartnerVendorType:
				profileVendorType = PartnerVendorType
			case CommunityVendorType:
				profileVendorType = CommunityVendorType
			case RedhatVendorType:
				profileVendorType = RedhatVendorType
			}
		}
	}

	profile, err := getProfile(profileVendorType, profileVersion)

	if err != nil {
		profile = getDefaultProfile(err.Error())
	} else if profile == nil {
		profile = getDefaultProfile(fmt.Sprintf("No matching profile found : %s : %s :", profileVendorType, profileVersion))
	}

	return profile
}

func getProfile(vendor VendorType, version *semver.Version) (*Profile, error) {

	var configDir string
	if isRunningInDockerContainer() {
		configDir = filepath.Join("app", "config")
	} else {
		_, fn, _, ok := runtime.Caller(0)
		if !ok {
			return nil, errors.New("failed to get profile directory")
		}
		fileParts := strings.SplitAfter(fn, "chart-verifier")
		if len(fileParts) > 1 {
			// running as cli
			configDir = filepath.Join(fileParts[0], "config")
		}
	}

	var profile *Profile

	filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".yaml") {
			profileRead, err := readProfile(path)
			if err == nil {
				if strings.Compare(string(profileRead.Vendor), string(vendor)) == 0 {
					profileVersion, err := semver.NewVersion(string(profileRead.Version))
					if err == nil {
						if profileVersion.Major() == version.Major() && profileVersion.Minor() == version.Minor() {
							profile = profileRead
							profile.Name = strings.Split(info.Name(), ".yaml")[0]
							return io.EOF
						}
					}
				}
			}
		}
		return nil
	})

	return profile, nil
}

func (profile *Profile) FilterChecks(registry checks.DefaultRegistry) FilteredRegistry {

	filteredChecks := make(map[checks.CheckName]checks.Check)

	for _, check := range profile.Checks {
		splitter := regexp.MustCompile(`/`)
		splitCheck := splitter.Split(check.Name, -1)
		checkIndex := checks.CheckId{Name: checks.CheckName(splitCheck[1]), Version: splitCheck[0]}
		if newCheck, ok := registry[checkIndex]; ok {
			newCheck.Type = check.Type
			filteredChecks[checkIndex.Name] = newCheck
		}
	}

	return filteredChecks

}

func readProfile(fileName string) (*Profile, error) {

	// Open the json file which defines the tests to run
	profileYaml, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	profileBytes, err := ioutil.ReadAll(profileYaml)
	if err != nil {
		return nil, err
	}

	profile = &Profile{}
	err = yaml.Unmarshal(profileBytes, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil

}

func isRunningInDockerContainer() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then verifier is running
	// from inside a container
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}
