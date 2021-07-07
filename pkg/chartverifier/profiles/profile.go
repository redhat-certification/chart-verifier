package profiles

import (
	"fmt"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"github.com/spf13/viper"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v3"
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

	VendorTypeConfigName string = "profile.vendortype"
	VersionConfigName    string = "profile.version"

	VendorTypeDefault      VendorType = "default"
	VendorTypeNotSpecified VendorType = "vendorTypeNotSpecified"
)

var profileMap map[VendorType][]*Profile

func init() {
	profileMap = make(map[VendorType][]*Profile)
	getProfiles()

	// add default profile to the map if a default profile was not found.
	if _, ok := profileMap[VendorTypeDefault]; !ok {
		profileMap[VendorTypeDefault] = append(profileMap[VendorTypeDefault], getDefaultProfile(""))
	}
}

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

var profileInUse *Profile

func Get() *Profile {
	if profileInUse == nil {
		return getDefaultProfile("No profile set for get")
	}
	return profileInUse
}

func New(config *viper.Viper) *Profile {

	profileVendorType := VendorTypeDefault
	var profileVersion string

	if config != nil {

		configVendorType := VendorType(config.GetString(VendorTypeConfigName))
		if len(configVendorType) > 0 {
			if _, ok := profileMap[configVendorType]; ok {
				profileVendorType = configVendorType
			}
		}

		configVersion := config.GetString(VersionConfigName)
		if len(configVersion) > 0 {
			if semver.IsValid(configVersion) {
				profileVersion = configVersion
			}
		}
	}

	profileInUse = getDefaultProfile(fmt.Sprintf("profile %s not found", profileVendorType))

	if vendorProfiles, ok := profileMap[profileVendorType]; ok {
		profileInUse = vendorProfiles[0]
		if len(vendorProfiles) > 1 {
			for _, vendorProfile := range vendorProfiles {
				if len(profileVersion) > 0 {
					if semver.Compare(semver.MajorMinor(vendorProfile.Version), semver.MajorMinor(profileVersion)) == 0 {
						profileInUse = vendorProfile
						break
					}
				}
				if semver.Compare(semver.MajorMinor(vendorProfile.Version), semver.MajorMinor(profileInUse.Version)) > 1 {
					profileInUse = vendorProfile
				}
			}
		}
	}

	return profileInUse
}

// Get all profiles in the profiles directory, and any subdirectories, and add each to the profile map
func getProfiles() {

	var configDir string
	if isRunningInDockerContainer() {
		configDir = filepath.Join("app", "config")
	} else {
		_, fn, _, ok := runtime.Caller(0)
		if !ok {
			return
		}
		index := strings.LastIndex(fn, "chart-verifier/")
		configDir = fn[0 : index+len("chart-verifier")]
		configDir = filepath.Join(configDir, "config")
	}

	filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if info != nil {
			if strings.HasSuffix(info.Name(), ".yaml") {
				profileRead, err := readProfile(path)
				if err == nil {
					// If version is not valid set to a default version
					if !semver.IsValid(profileRead.Version) {
						profileRead.Version = DefaultProfileVersion
					}
					if len(profileRead.Vendor) == 0 {
						profileRead.Vendor = VendorTypeNotSpecified
					}
					profileMap[profileRead.Vendor] = append(profileMap[profileRead.Vendor], profileRead)
					profileRead.Name = strings.Split(info.Name(), ".yaml")[0]
				}
			}
		}
		return nil
	})
	return
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

	// Open the yaml file which defines the tests to run
	profileYaml, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	profileBytes, err := ioutil.ReadAll(profileYaml)
	if err != nil {
		return nil, err
	}

	profile := &Profile{}
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
