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

package testutils

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"

	templatev1 "github.com/openshift/api/template/v1"

	// turns out we need the double legacy registration here
	// while in cluster we only need the modern way, otherwise
	// it breaks. TODO: investigate
	_ "github.com/fromanirh/kubevirt-template-indexer/pkg/okdlegacy"
	"github.com/fromanirh/kubevirt-template-indexer/pkg/templateindex"
)

type ByID []templateindex.Summary

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

func CheckSummaries(t *testing.T, summaries, expected []templateindex.Summary) {
	if len(expected) != len(summaries) {
		t.Errorf("expected %v summaries, received %v", len(expected), len(summaries))
	}
	for i, exp := range expected {
		if exp != summaries[i] {
			t.Errorf("expected=%#v received=%#v", exp, summaries[i])
		}
	}
}

func LoadYAML(path string) ([][]byte, error) {
	src, err := os.Open(path)
	defer src.Close()
	if err != nil {
		return [][]byte{}, err
	}

	data, err := ioutil.ReadAll(src)
	if err != nil {
		return [][]byte{}, err
	}

	return bytes.Split(data, []byte("---\n")), nil
}

func LoadTemplates(path string) ([]templatev1.Template, error) {
	var err error
	templates := []templatev1.Template{}
	yamls, err := LoadYAML(path)
	if err != nil {
		return templates, err
	}

	// Create a YAML serializer.  JSON is a subset of YAML, so is supported too.
	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme,
		scheme.Scheme)

	for _, yaml := range yamls {
		obj, _, err := s.Decode(yaml, nil, nil)
		if err != nil {
			return templates, err
		}

		templates = append(templates, *(obj.(*templatev1.Template)))
	}
	return templates, nil
}

// NullLogger is a logr.Logger that does nothing.
type NullLogger struct{}

func (_ NullLogger) Info(_ string, _ ...interface{}) {
}

func (_ NullLogger) Enabled() bool {
	return false
}

func (_ NullLogger) Error(_ error, _ string, _ ...interface{}) {
}

func (log NullLogger) V(_ int) logr.InfoLogger {
	return log
}

func (log NullLogger) WithName(_ string) logr.Logger {
	return log
}

func (log NullLogger) WithValues(_ ...interface{}) logr.Logger {
	return log
}
