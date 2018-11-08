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

func TestTemplateIndexerDescribeBySimpleFilter(t *testing.T) {
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

	OS := "centos7.0"

	descs, err := ti.DescribeBy(FilterOptions{
		"os": OS,
	})
	if err != nil || len(descs) < 1 {
		t.Errorf("unexpected output: %v err=%v", len(descs), err)
		return
	}

	// first, the filtered output must NOT include unwanted data
	for _, desc := range descs {
		if desc.OS != OS {
			t.Errorf("OS mismatch: requested %v found %v", OS, desc.OS)
		}
	}
	// TODO: then, the filtered output must include ALL wanted data
}
