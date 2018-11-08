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
)

func TestStringSetCreatedEmpty(t *testing.T) {
	s := NewStringSet()
	if s.Contains("foobar") {
		t.Errorf("unexpectedly contains foobar")
	}
}

func TestStringSetAdd(t *testing.T) {
	s := NewStringSet()
	s.Add("foobar")
	if !s.Contains("foobar") {
		t.Errorf("unexpectedly missing foobar")
	}
}

func TestStringSetAddRemove(t *testing.T) {
	s := NewStringSet()
	s.Add("foobar")
	s.Remove("foobar")
	if s.Contains("foobar") {
		t.Errorf("unexpectedly contains foobar")
	}
}
