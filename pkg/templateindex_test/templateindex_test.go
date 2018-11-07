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

	"github.com/fromanirh/kubevirt-template-indexer/internal/pkg/testutils"
	"github.com/fromanirh/kubevirt-template-indexer/pkg/templateindex"
)

func TestTemplateIndexerCreatedEmpty(t *testing.T) {
	ti := templateindex.NewTemplateIndexer(testutils.NullLogger{})
	count := ti.Count()
	if count != 0 {
		t.Errorf("unexpected count: %v", count)
	}
}

func TestTemplateIndexerInvalidLedger(t *testing.T) {
	ld, err := templateindex.NewJSONLedger("foobar", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	ti := templateindex.NewTemplateIndexer(testutils.NullLogger{})
	ti.AddLedger("foobar", ld)

	summaries, err := ti.SummarizeBy("foobar")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	sort.Sort(testutils.ByID(summaries))

	expected := []templateindex.Summary{}
	testutils.CheckSummaries(t, summaries, expected)
}
