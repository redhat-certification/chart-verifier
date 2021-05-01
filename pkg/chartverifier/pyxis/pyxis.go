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
package pyxis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var pyxisBaseUrl = "https://catalog.redhat.com/api/containers/v1/repositories"

type RepositoriesBody struct {
	PyxisRepositories []PyxisRepository `json:"data"`
}

type PyxisRepository struct {
	Id          string `json:"_id"`
	Repository  string `json:"repository"`
	VendorLabel string `json:"vendor_label"`
	Registry    string `json:"registry"`
}

type RegistriesBody struct {
	PyxisRegistries []PyxisRegistry `json:"data"`
}

type PyxisRegistry struct {
	Id           string               `json:"_id"`
	ParsedData   ImageData            `json:"parsed_data"`
	Repositories []RegistryRepository `json:"repositories"`
}

type ImageData struct {
	Digest string `json:"docker_image_digest"`
}

type RegistryRepository struct {
	Registry   string          `json:"registry"`
	Repository string          `json:"repository"`
	Tags       []RepositoryTag `json:"tags"`
}

type RepositoryTag struct {
	Digest string `json:"manifest_schema1_digest"`
	Name   string `json:"name"`
}

type ImageReference struct {
	Registries []string
	Repository string
	Tag        string
	Sha        string
}

func GetImageRegistries(repository string) ([]string, error) {
	var err error
	var registries []string

	req, _ := http.NewRequest("GET", pyxisBaseUrl, nil)
	queryString := req.URL.Query()
	queryString.Add("filter", fmt.Sprintf("repository==%s", repository))
	req.URL.RawQuery = queryString.Encode()
	req.Header.Set("X-API-KEY", "RedHatChartVerifier")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error getting repository %s : %v\n", repository, err))
	} else {
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			var repositoriesBody RepositoriesBody
			json.Unmarshal(body, &repositoriesBody)

			if len(repositoriesBody.PyxisRepositories) > 0 {
				for _, repo := range repositoriesBody.PyxisRepositories {
					registries = append(registries, repo.Registry)
				}
			} else {
				err = errors.New(fmt.Sprintf("Respository not found: %s", repository))
			}
		} else {
			err = errors.New(fmt.Sprintf("Bad response code from Pyxis: %d : %s", resp.StatusCode, req.URL))
		}
	}

	return registries, err
}

func IsImageInRegistry(imageRef ImageReference) (bool, error) {

	var err error
	found := false

	var tags []string
	var shas []string

Loops:
	for _, registry := range imageRef.Registries {

		requestUrl := fmt.Sprintf("%s/registry/%s/repository/%s/images", pyxisBaseUrl, registry, imageRef.Repository)
		req, _ := http.NewRequest("GET", requestUrl, nil)
		queryString := req.URL.Query()
		queryString.Add("filter", fmt.Sprintf("repositories=em=(repository==%s;registry==%s)", imageRef.Repository, registry))
		req.URL.RawQuery = queryString.Encode()
		req.Header.Set("X-API-KEY", "RedHatChartVerifier")
		client := &http.Client{}
		resp, reqErr := client.Do(req)

		if reqErr == nil {
			if resp.StatusCode == 200 {
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				var registriesBody RegistriesBody
				json.Unmarshal(body, &registriesBody)

				if len(registriesBody.PyxisRegistries) > 0 {
					found = false
					for _, reg := range registriesBody.PyxisRegistries {
						if len(imageRef.Sha) > 0 {
							if imageRef.Sha == reg.ParsedData.Digest {
								found = true
								err = nil
								continue Loops
							} else {
								shas = append(shas, reg.ParsedData.Digest)
							}
						} else {
							for _, repo := range reg.Repositories {
								if repo.Repository == imageRef.Repository && repo.Registry == registry {
									if len(imageRef.Sha) == 0 {
										for _, tag := range repo.Tags {
											if tag.Name == imageRef.Tag {
												found = true
												err = nil
												continue Loops
											} else {
												tags = append(tags, tag.Name)
											}
										}
									}
								}
							}
						}
					}
				} else {
					err = errors.New(fmt.Sprintf("No images found for Registry/Repository: %s/%s", registry, imageRef.Repository))
				}
			} else {
				err = errors.New(fmt.Sprintf("Bad response code %d from pyxis request : %s", resp.StatusCode, requestUrl))
			}
		} else {
			err = reqErr
		}
	}
	if !found {
		if err == nil {
			if len(imageRef.Sha) > 0 {
				err = errors.New(fmt.Sprintf("Digest %s not found. Found : %s", imageRef.Sha, strings.Join(shas, ", ")))
			} else {
				err = errors.New(fmt.Sprintf("Tag %s not found. Found : %s", imageRef.Tag, strings.Join(tags, ", ")))
			}
		}
	}
	return found, err
}
