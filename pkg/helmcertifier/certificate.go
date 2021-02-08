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

package helmcertifier

import "strconv"

type chartMetadata struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

type metadata struct {
	ChartMetadata chartMetadata `json:"chart" yaml:"chart"`
}

func newMetadata(name, version string) *metadata {
	return &metadata{
		ChartMetadata: chartMetadata{
			Name:    name,
			Version: version,
		},
	}
}

type certificate struct {
	Ok             bool           `json:"ok" yaml:"ok"`
	Metadata       *metadata      `json:"metadata" yaml:"metadata"`
	CheckResultMap checkResultMap `json:"results" yaml:"results"`
}

type checkResultMap map[string]checkResult

type checkResult struct {
	Ok     bool   `json:"ok" yaml:"ok"`
	Reason string `json:"reason" yaml:"reason"`
}

func newCertificate(name, version string, ok bool, resultMap checkResultMap) Certificate {
	return &certificate{
		Metadata:       newMetadata(name, version),
		Ok:             ok,
		CheckResultMap: resultMap,
	}
}

func (c *certificate) IsOk() bool {
	return c.Ok
}

func (c *certificate) String() string {
	report := "chart: " + c.Metadata.ChartMetadata.Name + "\n" +
		"version: " + c.Metadata.ChartMetadata.Version + "\n" +
		"ok: " + strconv.FormatBool(c.Ok) + "\n" +
		"\n"

	for k, v := range c.CheckResultMap {
		report += k + ":\n" +
			"\tok: " + strconv.FormatBool(v.Ok) + "\n" +
			"\treason: " + v.Reason + "\n"
	}

	return report
}
