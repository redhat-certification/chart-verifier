/*
 * Copyright 2021 Red Hat
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package checks

import (
	"fmt"
	"path"
	"strings"

	"helm.sh/helm/v3/pkg/lint"
	"helm.sh/helm/v3/pkg/lint/support"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/redhat-certification/chart-verifier/pkg/chartverifier/pyxis"
)

const (
	APIVersion2                  = "v2"
	ReadmeExist                  = "Chart has a README"
	ReadmeDoesNotExist           = "Chart does not have a README"
	NotHelm3Reason               = "API version is not V2, used in Helm 3"
	Helm3Reason                  = "API version is V2, used in Helm 3"
	TestTemplatePrefix           = "templates/tests/"
	ChartTestFilesExist          = "Chart test files exist"
	ChartTestFilesDoesNotExist   = "Chart test files do not exist"
	MinKuberVersionSpecified     = "Minimum Kubernetes version specified"
	MinKuberVersionNotSpecified  = "Minimum Kubernetes version is not specified"
	ValuesSchemaFileExist        = "Values schema file exist"
	ValuesSchemaFileDoesNotExist = "Values schema file does not exist"
	ValuesFileExist              = "Values file exist"
	ValuesFileDoesNotExist       = "Values file does not exist"
	ChartContainCRDs             = "Chart contains CRDs"
	ChartDoesNotContainCRDs      = "Chart does not contain CRDs"
	HelmLintSuccessful           = "Helm lint successful"
	HelmLintHasFailedPrefix      = "Helm lint has failed: "
	CSIObjectsExist              = "CSI objects exist"
	CSIObjectsDoesNotExist       = "CSI objects do not exist"
	NoImagesToCertify            = "No images to certify"
	ImageCertifyFailed           = "Failed to certify images"
	ImageCertified               = "Image is Red Hat certified"
	ImageNotCertified            = "Image is not Red Hat certified"
)

func notImplemented() (Result, error) {
	return Result{Ok: false}, errors.New("not implemented")
}

func IsHelmV3(uri string, _ *viper.Viper) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}
	isHelmV3 := c.Metadata.APIVersion == APIVersion2

	reason := NotHelm3Reason
	if isHelmV3 {
		reason = Helm3Reason
	}
	return NewResult(isHelmV3, reason), nil
}

func HasReadme(uri string, _ *viper.Viper) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	r := NewResult(false, ReadmeDoesNotExist)
	for _, f := range c.Files {
		if f.Name == "README.md" {
			r.SetResult(true, ReadmeExist)
		}
	}

	return r, nil
}

func ContainsTest(uri string, _ *viper.Viper) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	r := NewResult(false, ChartTestFilesDoesNotExist)
	for _, f := range c.Templates {
		if strings.HasPrefix(f.Name, TestTemplatePrefix) && strings.HasSuffix(f.Name, ".yaml") {
			r.Ok = true
			r.SetResult(true, ChartTestFilesExist)
			break
		}
	}

	return r, nil

}

func ContainsValues(uri string, _ *viper.Viper) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	r := NewResult(false, ValuesFileDoesNotExist)

	if len(c.Values) > 0 {
		r.SetResult(true, ValuesFileExist)
	}

	return r, nil
}

func ContainsValuesSchema(uri string, _ *viper.Viper) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}

	r := NewResult(false, ValuesSchemaFileDoesNotExist)

	if len(c.Schema) > 0 {
		r.SetResult(true, ValuesSchemaFileExist)
	}

	return r, nil
}

func KeywordsAreOpenshiftCategories(uri string, _ *viper.Viper) (Result, error) {
	return notImplemented()
}

func IsCommercialChart(uri string, _ *viper.Viper) (Result, error) {
	return notImplemented()
}

func IsCommunityChart(uri string, _ *viper.Viper) (Result, error) {
	return notImplemented()
}

func HasMinKubeVersion(uri string, _ *viper.Viper) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return NewResult(false, err.Error()), err
	}

	r := NewResult(false, MinKuberVersionNotSpecified)

	if c.Metadata.KubeVersion != "" {
		r.SetResult(true, MinKuberVersionSpecified)
	}

	return r, nil
}

func NotContainCRDs(uri string, _ *viper.Viper) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return NewResult(false, err.Error()), err
	}

	r := NewResult(true, ChartDoesNotContainCRDs)

	if len(c.CRDObjects()) > 0 {
		r.Ok = false
		r.SetResult(false, ChartContainCRDs)
	}

	return r, nil
}

func HelmLint(uri string, _ *viper.Viper) (Result, error) {
	c, p, err := LoadChartFromURI(uri)
	if err != nil {
		return NewResult(false, err.Error()), err
	}
	r := NewResult(true, HelmLintSuccessful)
	p = path.Join(p, c.Name())
	linter := lint.All(p, map[string]interface{}{}, "default", false)
	if linter.HighestSeverity > support.WarningSev {
		reason := ""
		for _, m := range linter.Messages {
			reason = reason + m.Error() + "\n"
		}
		r.SetResult(false, fmt.Sprintf("%s %s", HelmLintHasFailedPrefix, reason))
	}
	return r, nil
}

func NotContainsInfraPluginsAndDrivers(uri string, _ *viper.Viper) (Result, error) {
	return notImplemented()
}

func NotContainCSIObjects(uri string, _ *viper.Viper) (Result, error) {
	c, _, err := LoadChartFromURI(uri)
	if err != nil {
		return Result{}, err
	}
	r := NewResult(true, CSIObjectsDoesNotExist)
	for _, f := range c.Templates {
		if !strings.HasSuffix(f.Name, ".yaml") {
			continue
		}
		for _, v := range strings.Split(string(f.Data), "\n") {
			if strings.HasPrefix(v, "kind") {
				if strings.TrimSpace(strings.Split(v, ":")[1]) == "CSIDriver" {
					r.SetResult(false, CSIObjectsExist)
				}
			}
		}
	}

	return r, nil
}

func CanBeInstalledWithoutManualPreRequisites(uri string, _ *viper.Viper) (Result, error) {
	return notImplemented()
}

func CanBeInstalledWithoutClusterAdminPrivileges(uri string, _ *viper.Viper) (Result, error) {
	return notImplemented()
}

func ImagesAreCertified(uri string, _ *viper.Viper) (Result, error) {

	r := NewResult(false, "")

	images, err := getImageReferences(uri)

	if err != nil {
		r.SetResult(false, fmt.Sprintf("%s : Failed to get images : %v", ImageCertifyFailed, err))
	} else if len(images) == 0 {
		r.SetResult(true, NoImagesToCertify)
	} else {
		for _, image := range images {

			registries, repository, version := getImageParts(image)

			if len(registries) == 0 {
				registries, err = pyxis.GetImageRegistries(repository)
			}

			if err != nil {
				r.AddResult(false, fmt.Sprintf("%s : %s : %v", ImageNotCertified, image, err))
			} else if len(registries) == 0 {
				r.AddResult(false, fmt.Sprintf("%s : %s", ImageNotCertified, image))
			} else {
				certified := false
				for _, registry := range registries {
					found, checkImageErr := pyxis.IsImageInRegistry(repository, version, registry)
					if found {
						err = nil
						certified = true
						break
					} else if err == nil {
						err = checkImageErr
					}
				}
				if !certified {
					if err != nil {
						r.AddResult(false, fmt.Sprintf("%s : %s : %v", ImageNotCertified, image, err))
					} else {
						r.AddResult(false, fmt.Sprintf("%s : %s", ImageNotCertified, image))
					}
				} else {
					r.AddResult(true, fmt.Sprintf("%s : %s", ImageCertified, image))
				}
			}
		}
	}

	return r, nil
}

func getImageParts(image string) ([]string, string, string) {

	imageParts := strings.Split(image, "/")

	lastPart := imageParts[len(imageParts)-1]
	lastParts := strings.Split(lastPart, ":")
	var version string
	if len(lastParts) > 1 && len(lastParts[1]) > 0 {
		version = lastParts[1]
	} else {
		version = "latest"
	}

	imageParts[len(imageParts)-1] = lastParts[0]

	var registries []string
	var repository string
	if len(imageParts) > 2 && len(imageParts[0]) > 1 {
		registries = append(registries, imageParts[0])
		repository = strings.Join(imageParts[1:], "/")
	} else {
		repository = strings.Join(imageParts, "/")
	}
	return registries, repository, version
}
