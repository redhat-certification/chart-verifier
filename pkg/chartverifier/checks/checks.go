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

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/lint"
	"helm.sh/helm/v3/pkg/lint/support"

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
	KuberVersionSpecified        = "Kubernetes version specified"
	KuberVersionNotSpecified     = "Kubernetes version is not specified"
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

func IsHelmV3(opts *CheckOptions) (Result, error) {
	c, _, err := LoadChartFromURI(opts.URI)
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

func HasReadme(opts *CheckOptions) (Result, error) {
	c, _, err := LoadChartFromURI(opts.URI)
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

func ContainsTest(opts *CheckOptions) (Result, error) {
	c, _, err := LoadChartFromURI(opts.URI)
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

func ContainsValues(opts *CheckOptions) (Result, error) {
	c, _, err := LoadChartFromURI(opts.URI)
	if err != nil {
		return Result{}, err
	}

	r := NewResult(false, ValuesFileDoesNotExist)

	if len(c.Values) > 0 {
		r.SetResult(true, ValuesFileExist)
	}

	return r, nil
}

func ContainsValuesSchema(opts *CheckOptions) (Result, error) {
	c, _, err := LoadChartFromURI(opts.URI)
	if err != nil {
		return Result{}, err
	}

	r := NewResult(false, ValuesSchemaFileDoesNotExist)

	if len(c.Schema) > 0 {
		r.SetResult(true, ValuesSchemaFileExist)
	}

	return r, nil
}

func KeywordsAreOpenshiftCategories(opts *CheckOptions) (Result, error) {
	return notImplemented()
}

func IsCommercialChart(opts *CheckOptions) (Result, error) {
	return notImplemented()
}

func IsCommunityChart(opts *CheckOptions) (Result, error) {
	return notImplemented()
}

func HasKubeVersion(opts *CheckOptions) (Result, error) {
	c, _, err := LoadChartFromURI(opts.URI)
	if err != nil {
		return NewResult(false, err.Error()), err
	}

	r := NewResult(false, KuberVersionNotSpecified)

	if c.Metadata.KubeVersion != "" {
		r.SetResult(true, KuberVersionSpecified)
	}

	return r, nil
}

func NotContainCRDs(opts *CheckOptions) (Result, error) {
	c, _, err := LoadChartFromURI(opts.URI)
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

func HelmLint(opts *CheckOptions) (Result, error) {
	c, p, err := LoadChartFromURI(opts.URI)
	if err != nil {
		return NewResult(false, err.Error()), err
	}
	r := NewResult(true, HelmLintSuccessful)
	p = path.Join(p, c.Name())
	linter := lint.All(p, opts.Values, "default", false)
	if linter.HighestSeverity > support.WarningSev {
		reason := ""
		for _, m := range linter.Messages {
			reason = reason + m.Error() + "\n"
		}
		r.SetResult(false, fmt.Sprintf("%s %s", HelmLintHasFailedPrefix, reason))
	}
	return r, nil
}

func NotContainsInfraPluginsAndDrivers(opts *CheckOptions) (Result, error) {
	return notImplemented()
}

func NotContainCSIObjects(opts *CheckOptions) (Result, error) {
	c, _, err := LoadChartFromURI(opts.URI)
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

func CanBeInstalledWithoutManualPreRequisites(opts *CheckOptions) (Result, error) {
	return notImplemented()
}

func CanBeInstalledWithoutClusterAdminPrivileges(opts *CheckOptions) (Result, error) {
	return notImplemented()
}

func ImagesAreCertified(opts *CheckOptions) (Result, error) {

	r := NewResult(true, "")

	images, err := getImageReferences(opts.URI, opts.Values)

	if err != nil {
		r.SetResult(false, fmt.Sprintf("%s : Failed to get images, error running helm template : %v", ImageCertifyFailed, err))
	} else if len(images) == 0 {
		r.SetResult(true, NoImagesToCertify)
	} else {
		for _, image := range images {

			err = nil
			imageRef := parseImageReference(image)

			if len(imageRef.Registries) == 0 {
				imageRef.Registries, err = pyxis.GetImageRegistries(imageRef.Repository)
			}

			if err != nil {
				r.AddResult(false, fmt.Sprintf("%s : %s : %v", ImageNotCertified, image, err))
			} else if len(imageRef.Registries) == 0 {
				r.AddResult(false, fmt.Sprintf("%s : %s", ImageNotCertified, image))
			} else {
				certified, checkImageErr := pyxis.IsImageInRegistry(imageRef)
				if !certified {
					if checkImageErr != nil {
						r.AddResult(false, fmt.Sprintf("%s : %s : %v", ImageNotCertified, image, checkImageErr))
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

func parseImageReference(image string) pyxis.ImageReference {

	imageRef := pyxis.ImageReference{}
	imageParts := strings.Split(image, "/")

	lastPart := imageParts[len(imageParts)-1]
	var lastParts []string
	if strings.Contains(lastPart, "@sha") {
		lastParts = strings.Split(lastPart, "@")
		imageRef.Sha = lastParts[1]
	} else {
		lastParts = strings.Split(lastPart, ":")
		if len(lastParts) > 1 && len(lastParts[1]) > 0 {
			imageRef.Tag = lastParts[1]
		} else {
			imageRef.Tag = "latest"
		}
	}

	imageParts[len(imageParts)-1] = lastParts[0]

	if len(imageParts) > 1 && len(imageParts[0]) > 1 && strings.Contains(imageParts[0], ".") {
		imageRef.Registries = append(imageRef.Registries, imageParts[0])
		imageRef.Repository = strings.Join(imageParts[1:], "/")
	} else {
		imageRef.Repository = strings.Join(imageParts, "/")
	}

	return imageRef

}
