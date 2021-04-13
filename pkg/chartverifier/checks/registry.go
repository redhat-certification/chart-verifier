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

import "github.com/spf13/viper"

type Result struct {
	// Ok indicates whether the result was successful or not.
	Ok bool
	// Reason for the result value.  This is a message indicating
	// the reason for the value of Ok became true or false.
	Reason string
}

func NewResult(outcome bool, reason string) Result {
	result := Result{}
	result.Ok = outcome
	result.Reason = reason
	return result
}

func (r *Result) SetResult(outcome bool, reason string) Result {
	r.Ok = outcome
	r.Reason = reason
	return *r
}

func (r *Result) AddResult(outcome bool, reason string) Result {
	r.Ok = r.Ok && outcome
	if len(r.Reason) > 0 {
		r.Reason += "\n\t\t"
	}
	r.Reason += reason
	return *r
}

type CheckFunc func(uri string, config *viper.Viper) (Result, error)

type Registry interface {
	Get(name string) (CheckFunc, bool)
	Add(name string, checkFunc CheckFunc) Registry
	AllChecks() []string
}

type defaultRegistry map[string]CheckFunc

func (r *defaultRegistry) AllChecks() []string {
	allChecks := make([]string, 0)
	for k, _ := range *r {
		allChecks = append(allChecks, k)
	}
	return allChecks
}

func NewRegistry() Registry {
	return &defaultRegistry{}
}

func (r *defaultRegistry) Get(name string) (CheckFunc, bool) {
	v, ok := (*r)[name]
	return v, ok
}

func (r *defaultRegistry) Add(name string, checkFunc CheckFunc) Registry {
	(*r)[name] = checkFunc
	return r
}
