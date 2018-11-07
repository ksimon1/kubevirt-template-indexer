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

type StringSet struct {
	data map[string]bool
}

func NewStringSet() *StringSet {
	return &StringSet{
		data: make(map[string]bool),
	}
}

func (sts *StringSet) Add(key string) {
	sts.data[key] = true
}

func (sts *StringSet) Remove(key string) {
	delete(sts.data, key)
}

func (sts *StringSet) Contains(key string) bool {
	return sts.data[key]
}
