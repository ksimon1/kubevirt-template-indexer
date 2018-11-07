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
	"fmt"
	"strings"

	templatev1 "github.com/openshift/api/template/v1"
)

const (
	suffix = "template.cnv.io"
)

func splitLabel(label string) (string, bool) {
	items := strings.Split(label, "/")
	if len(items) != 2 {
		return "", false
	}
	return items[1], true
}

func fixLabelKey(key string) string {
	if key == "size" {
		return "flavor"
	}
	return key
}

func makeLabel(key, value string) string {
	return fmt.Sprintf("%s.%s/%s", key, suffix, value)
}

func extractFlavours(t *templatev1.Template, label string) []string {
	flavours := []string{}
	label = fmt.Sprintf("%s.%s", label, suffix)
	for key, value := range t.Labels {
		if flavour, ok := tryToGetFlavour(key, value, label); ok {
			flavours = append(flavours, flavour)
		}

	}
	return flavours
}

func tryToGetFlavour(key, value, label string) (string, bool) {
	if strings.HasPrefix(key, label) && value == "true" {
		if flavour, ok := splitLabel(key); ok {
			return flavour, ok
		}
	}
	return "", false
}
