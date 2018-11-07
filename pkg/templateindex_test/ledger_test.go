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

package templateindex_test

import (
	"sort"
	"testing"

	templatev1 "github.com/openshift/api/template/v1"

	"github.com/fromanirh/kubevirt-template-indexer/internal/pkg/testutils"
	"github.com/fromanirh/kubevirt-template-indexer/pkg/templateindex"
)

var templates []templatev1.Template

func init() {
	var err error
	templates, err = testutils.LoadTemplates("template-data.yaml")
	if err != nil || len(templates) < 1 {
		panic("cannot load test templates!")
	}
}

func TestDescribe(t *testing.T) {
	templates, err := testutils.LoadTemplates("template-data.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates!")
		return
	}
	for _, template := range templates {
		desc := templateindex.Describe(&template)
		if desc.Name == "" || desc.ID == "" || desc.Icon == "" || desc.OS == "" || desc.Workload == "" || desc.Size == "" {
			t.Errorf("%#v", desc)
		}
	}
}

func TestCreateJSONLedgerWithoutFile(t *testing.T) {
	_, err := templateindex.NewJSONLedger("os", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInvalidLedger(t *testing.T) {
	ld, err := templateindex.NewJSONLedger("foobar", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	summaries := ld.Summarize(templates)
	sort.Sort(testutils.ByID(summaries))

	expected := []templateindex.Summary{}
	testutils.CheckSummaries(t, summaries, expected)
}

func TestOSLedger(t *testing.T) {
	ld, err := templateindex.NewJSONLedger("os", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	summaries := ld.Summarize(templates)
	sort.Sort(testutils.ByID(summaries))

	expected := []templateindex.Summary{
		templateindex.Summary{ID: "centos7.0"},
		templateindex.Summary{ID: "fedora26"},
		templateindex.Summary{ID: "fedora27"},
		templateindex.Summary{ID: "fedora28"},
		templateindex.Summary{ID: "opensuse15.0"},
		templateindex.Summary{ID: "rhel7.0"},
		templateindex.Summary{ID: "rhel7.1"},
		templateindex.Summary{ID: "rhel7.2"},
		templateindex.Summary{ID: "rhel7.3"},
		templateindex.Summary{ID: "rhel7.4"},
		templateindex.Summary{ID: "rhel7.5"},
		templateindex.Summary{ID: "ubuntu18.04"},
		templateindex.Summary{ID: "win10"},
		templateindex.Summary{ID: "win2k12r2"},
		templateindex.Summary{ID: "win2k8"},
		templateindex.Summary{ID: "win2k8r2"},
	}
	testutils.CheckSummaries(t, summaries, expected)
}

func TestWorkloadLedger(t *testing.T) {
	ld, err := templateindex.NewJSONLedger("workload", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	summaries := ld.Summarize(templates)
	sort.Sort(testutils.ByID(summaries))

	expected := []templateindex.Summary{
		templateindex.Summary{ID: "generic"},
		templateindex.Summary{ID: "highperformance"},
	}
	testutils.CheckSummaries(t, summaries, expected)
}

func TestSizeLedger(t *testing.T) {
	ld, err := templateindex.NewJSONLedger("flavor", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	summaries := ld.Summarize(templates)
	sort.Sort(testutils.ByID(summaries))

	expected := []templateindex.Summary{
		templateindex.Summary{ID: "large"},
		templateindex.Summary{ID: "medium"},
		templateindex.Summary{ID: "small"},
		templateindex.Summary{ID: "tiny"},
	}
	testutils.CheckSummaries(t, summaries, expected)
}
