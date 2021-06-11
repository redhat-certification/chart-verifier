package profiles

import (
	"errors"
	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/checks"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type Annotation string

const (
	DigestAnnotation                 Annotation = "Digest"
	OCPVersionAnnotation             Annotation = "OCPVersion"
	LastCertifiedTimestampAnnotation Annotation = "LastCertifiedTimestamp"
)

type Profile struct {
	Apiversion  string       `json:"apiversion" yaml:"apiversion"`
	Kind        string       `json:"kind" yaml:"kind"`
	Name        string       `json:"name" yaml:"name"`
	Annotations []Annotation `json:"annotations" yaml:"annotations"`
	Checks      []*Check     `json:"checks" yaml:"checks"`
}

type Check struct {
	Name    checks.CheckName `json:"name" yaml:"name"`
	Type    checks.CheckType `json:"type" yaml:"type"`
	Version string           `json:"version" yaml:"version"`
}

var profile *Profile

func GetProfile() *Profile {

	if profile != nil {
		return profile
	}

	fileName, err := getProfileFileName()
	if err != nil {
		profile = getDefaultProfile(err.Error())
		return profile
	}

	// Open the json file which defines the tests to run
	profileYaml, err := os.Open(fileName)
	if err != nil {
		profile = getDefaultProfile(err.Error())
		return profile
	}

	profileBytes, err := ioutil.ReadAll(profileYaml)
	if err != nil {
		profile = getDefaultProfile(err.Error())
		return profile
	}

	profile = &Profile{}
	err = yaml.Unmarshal(profileBytes, profile)
	if err != nil {
		profile = getDefaultProfile(err.Error())
		return profile
	}
	return profile
}

func getProfileFileName() (string, error) {

	_, fn, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to get profile directory")
	}

	// To be update when multiple profiles are supported
	return filepath.Join(filepath.Dir(fn), "profile-1.0.0.yaml"), nil
}
