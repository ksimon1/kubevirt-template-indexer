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
	"net/url"
	"testing"

	"github.com/fromanirh/kubevirt-template-indexer/internal/pkg/testutils"
)

func TestTemplateIndexerCreatedEmpty(t *testing.T) {
	ti := NewTemplateIndexer(testutils.NullLogger{})
	count := ti.Count()
	if count != 0 {
		t.Errorf("unexpected count: %v", count)
	}
}

func TestTemplateIndexerUnknownLedger(t *testing.T) {
	ti := NewTemplateIndexer(testutils.NullLogger{})
	summaries, err := ti.SummarizeBy("unknown")
	if err == nil {
		t.Errorf("unexpectedly succesful")
	}
	if len(summaries) != 0 {
		t.Errorf("unexpected output: %v", summaries)
	}

}

func TestTemplateIndexerEmptyLedger(t *testing.T) {
	ld, err := NewJSONLedger("foobar", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	ti := NewTemplateIndexer(testutils.NullLogger{})
	ti.AddLedger("foobar", ld)

	summaries, err := ti.SummarizeBy("foobar")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := []Summary{}
	checkSummaries(t, summaries, expected)
}

func TestTemplateIndexerWorkloadLedger(t *testing.T) {
	templates, err := testutils.LoadTemplates("test-data-alltemplates.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates! %v", err)
		return
	}

	ld, err := NewJSONLedger("workload", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	ti := NewTemplateIndexer(testutils.NullLogger{})
	ti.AddLedger("workload", ld)

	count, err := ti.AddTemplates(templates)
	if err != nil || count != len(templates) {
		t.Errorf("failed to add test templates! %v", err)
		return
	}

	summaries, err := ti.SummarizeBy("workload")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := []Summary{
		Summary{ID: "generic"},
		Summary{ID: "highperformance"},
	}
	checkSummaries(t, summaries, expected)
}

func TestTemplateIndexerDescribeBySimpleFilter(t *testing.T) {
	templates, err := testutils.LoadTemplates("test-data-alltemplates.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates! %v", err)
		return
	}

	ti := NewTemplateIndexer(testutils.NullLogger{})
	count, err := ti.AddTemplates(templates)
	if err != nil || count != len(templates) {
		t.Errorf("failed to add test templates! %v", err)
		return
	}

	size := "medium"

	descs, err := ti.DescribeBy(FilterOptions{
		"size": size,
	})
	if err != nil || len(descs) < 1 {
		t.Errorf("unexpected output: %v err=%v", len(descs), err)
		return
	}

	// first, the filtered output must NOT include unwanted data
	for _, desc := range descs {
		if desc.Size != size {
			t.Errorf("Size mismatch: requested %v found %v", size, desc.Size)
		}
	}
	// TODO: then, the filtered output must include ALL wanted data
}

func TestTemplateIndexerDescribeByFullFilter(t *testing.T) {
	templates, err := testutils.LoadTemplates("test-data-alltemplates.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates! %v", err)
		return
	}

	ti := NewTemplateIndexer(testutils.NullLogger{})
	count, err := ti.AddTemplates(templates)
	if err != nil || count != len(templates) {
		t.Errorf("cannot add test templates! %v", err)
		return
	}

	os := "centos7.0"
	size := "medium"
	workload := "generic"

	descs, err := ti.DescribeBy(FilterOptions{
		"size":     size,
		"os":       os,
		"workload": workload,
	})
	if err != nil || len(descs) < 1 {
		t.Errorf("unexpected output: %v err=%v", len(descs), err)
		return
	}

	// first, the filtered output must NOT include unwanted data
	for _, desc := range descs {
		if desc.Size != size {
			t.Errorf("Size mismatch: requested %v found %v", size, desc.Size)
		}
		if desc.OS != os {
			t.Errorf("OS mismatch: requested %v found %v", os, desc.OS)
		}
		if desc.Workload != workload {
			t.Errorf("OS mismatch: requested %v found %v", workload, desc.Workload)
		}
	}
	// TODO: then, the filtered output must include ALL wanted data
}

func TestTemplateIndexerDescribeByFullFilterFromURL(t *testing.T) {
	templates, err := testutils.LoadTemplates("test-data-alltemplates.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates! %v", err)
		return
	}

	ti := NewTemplateIndexer(testutils.NullLogger{})
	count, err := ti.AddTemplates(templates)
	if err != nil || count != len(templates) {
		t.Errorf("cannot add test templates! %v", err)
		return
	}

	u, err := url.Parse("http://localhost:18081/templates?size=medium&os=centos7.0&workload=generic")
	if err != nil {
		t.Errorf("cannot parse url %v", err)
		return
	}

	descs, err := ti.DescribeBy(FilterOptionsFromURL(u))
	if err != nil || len(descs) < 1 {
		t.Errorf("unexpected output: %v err=%v", len(descs), err)
		return
	}

	q := u.Query()
	// first, the filtered output must NOT include unwanted data
	for _, desc := range descs {
		if desc.Size != q.Get("size") {
			t.Errorf("Size mismatch: requested %v found %v", q.Get("size"), desc.Size)
		}
		if desc.OS != q.Get("os") {
			t.Errorf("OS mismatch: requested %v found %v", q.Get("os"), desc.OS)
		}
		if desc.Workload != q.Get("workload") {
			t.Errorf("OS mismatch: requested %v found %v", q.Get("workload"), desc.Workload)
		}
	}
	// TODO: then, the filtered output must include ALL wanted data
}

func TestTemplateIndexerAddUsingUpdate(t *testing.T) {
	templates, err := testutils.LoadTemplates("test-data-alltemplates.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates! %v", err)
		return
	}

	ld, err := NewJSONLedger("workload", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	ti := NewTemplateIndexer(testutils.NullLogger{})
	ti.AddLedger("workload", ld)

	for _, template := range templates {
		ti.Update(&template)
	}
	if ti.Count() != len(templates) {
		t.Errorf("failed to add test templates! %v", err)
		return
	}

	// we just need something.
	descs, err := ti.DescribeBy(FilterOptions{})
	if err != nil || len(descs) < 1 {
		t.Errorf("missing output: %v", err)
		return
	}

	summaries, err := ti.SummarizeBy("workload")
	if err != nil || len(summaries) < 1 {
		t.Errorf("missing output: %v", err)
		return
	}
}

func TestTemplateIndexerRemoveUsingUpdate(t *testing.T) {
	templates, err := testutils.LoadTemplates("test-data-alltemplates.yaml")
	if err != nil || len(templates) < 1 {
		t.Errorf("cannot load test templates! %v", err)
		return
	}

	ld, err := NewJSONLedger("workload", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	ti := NewTemplateIndexer(testutils.NullLogger{})
	ti.AddLedger("workload", ld)

	count, err := ti.AddTemplates(templates)
	if err != nil || count != len(templates) {
		t.Errorf("failed to add test templates! %v", err)
		return
	}

	for _, template := range templates {
		ti.Update(&template)
	}
	if ti.Count() != 0 {
		t.Errorf("failed to remove test templates! %v", err)
		return
	}

	// we just need something.
	descs, err := ti.DescribeBy(FilterOptions{})
	if err != nil || len(descs) != 0 {
		t.Errorf("unexpected output: %v", err)
		return
	}

	summaries, err := ti.SummarizeBy("workload")
	if err != nil || len(summaries) != 0 {
		t.Errorf("unexpected output: %v", err)
		return
	}
}

// TODO: test to remove just a specific template
