/*
 * This file is part of the KubeVirt project
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
 *
 * Copyright 2018 Red Hat, Inc.
 */

package templateindex

import (
	"sort"

	templatev1 "github.com/openshift/api/template/v1"
)

type Summary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Description struct {
	Summary
	Description string `json:"description"`
	Icon        string `json:"icon-id"`
	OS          string `json:"osid"`
	Workload    string `json:"workload"`
	Size        string `json:"size"`
}

type FilterOptions map[string]string

func Describe(t *templatev1.Template, opts FilterOptions) Description {
	desc := Description{
		Summary: Summary{
			ID:   t.Name,
			Name: t.Annotations["openshift.io/display-name"],
		},
		Description: t.Annotations["description"],
		Icon:        t.Annotations["iconClass"],
		OS:          opts["os"],
		Workload:    opts["workload"],
		Size:        opts["size"],
	}
	for key, value := range t.Labels {
		if os, ok := tryToGetFlavour(key, value, "os.template.cnv.io"); ok && desc.OS == "" {
			desc.OS = os
		}
		if workload, ok := tryToGetFlavour(key, value, "workload.template.cnv.io"); ok && desc.Workload == "" {
			desc.Workload = workload
		}
		if size, ok := tryToGetFlavour(key, value, "flavor.template.cnv.io"); ok && desc.Size == "" {
			desc.Size = size
		}
	}
	return desc
}

type Ledger interface {
	Summarize([]templatev1.Template) []Summary
}

type JSONLedger struct {
	label string
	names map[string]string
}

func NewJSONLedger(label, path string) (*JSONLedger, error) {
	ld := &JSONLedger{
		label: label,
		names: make(map[string]string),
	}
	// TODO: read names
	return ld, nil
}

func (ld *JSONLedger) Summarize(templates []templatev1.Template) []Summary {
	seen := NewStringSet()
	summaries := []Summary{}

	for _, template := range templates {
		flavours := extractFlavours(&template, ld.label)

		for _, flavour := range flavours {
			if seen.Contains(flavour) {
				continue
			}

			summaries = append(summaries, Summary{
				ID:   flavour,
				Name: ld.names[flavour],
			})
			seen.Add(flavour)
		}
	}
	sort.Sort(byID(summaries))
	return summaries
}

type byID []Summary

func (a byID) Len() int           { return len(a) }
func (a byID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byID) Less(i, j int) bool { return a[i].ID < a[j].ID }
