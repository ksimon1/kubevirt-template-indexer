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
	"testing"

	"github.com/fromanirh/kubevirt-template-indexer/internal/pkg/testutils"
)

func TestDescribe(t *testing.T) {
	templates, err := testutils.LoadTemplates("test-data-alltemplates.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates! %v", err)
		return
	}
	for _, template := range templates {
		desc := Describe(&template, FilterOptions{})
		if desc.Name == "" || desc.ID == "" || desc.Icon == "" || desc.OS == "" || desc.Workload == "" || desc.Size == "" {
			t.Errorf("%#v", desc)
		}
	}
}

func TestCreateJSONLedgerWithoutFile(t *testing.T) {
	_, err := NewJSONLedger("os", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInvalidLedger(t *testing.T) {
	templates, err := testutils.LoadTemplates("test-data-alltemplates.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates! %v", err)
		return
	}

	ld, err := NewJSONLedger("foobar", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	summaries := ld.Summarize(templates)

	expected := []Summary{}
	checkSummaries(t, summaries, expected)
}

func TestOSLedger(t *testing.T) {
	templates, err := testutils.LoadTemplates("test-data-alltemplates.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates! %v", err)
		return
	}

	ld, err := NewJSONLedger("os", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	summaries := ld.Summarize(templates)

	expected := []Summary{
		Summary{ID: "centos7.0"},
		Summary{ID: "fedora26"},
		Summary{ID: "fedora27"},
		Summary{ID: "fedora28"},
		Summary{ID: "opensuse15.0"},
		Summary{ID: "rhel7.0"},
		Summary{ID: "rhel7.1"},
		Summary{ID: "rhel7.2"},
		Summary{ID: "rhel7.3"},
		Summary{ID: "rhel7.4"},
		Summary{ID: "rhel7.5"},
		Summary{ID: "ubuntu18.04"},
		Summary{ID: "win10"},
		Summary{ID: "win2k12r2"},
		Summary{ID: "win2k8"},
		Summary{ID: "win2k8r2"},
	}
	checkSummaries(t, summaries, expected)
}

func TestWorkloadLedger(t *testing.T) {
	templates, err := testutils.LoadTemplates("test-data-alltemplates.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates! %v", err)
		return
	}

	ld, err := NewJSONLedger("workload", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	summaries := ld.Summarize(templates)

	expected := []Summary{
		Summary{ID: "generic"},
		Summary{ID: "highperformance"},
	}
	checkSummaries(t, summaries, expected)
}

func TestSizeLedger(t *testing.T) {
	templates, err := testutils.LoadTemplates("test-data-alltemplates.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates! %v", err)
		return
	}

	ld, err := NewJSONLedger("flavor", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	summaries := ld.Summarize(templates)

	expected := []Summary{
		Summary{ID: "large"},
		Summary{ID: "medium"},
		Summary{ID: "small"},
		Summary{ID: "tiny"},
	}
	checkSummaries(t, summaries, expected)
}

func checkSummaries(t *testing.T, summaries, expected []Summary) {
	if len(expected) != len(summaries) {
		t.Errorf("expected %v summaries, received %v", len(expected), len(summaries))
	}
	for i, exp := range expected {
		if exp != summaries[i] {
			t.Errorf("expected=%#v received=%#v", exp, summaries[i])
		}
	}
}
