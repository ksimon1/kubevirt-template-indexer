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
	"sync"

	"github.com/go-logr/logr"

	templatev1 "github.com/openshift/api/template/v1"

	"k8s.io/apimachinery/pkg/types"
)

type TemplateIndex struct {
	rwlock    sync.RWMutex
	log       logr.Logger
	templates map[types.UID]templatev1.Template
}

func NewTemplateIndex(log logr.Logger) *TemplateIndex {
	return &TemplateIndex{
		log:       log,
		templates: make(map[types.UID]templatev1.Template),
	}
}

// Set the initial state of the index. You must call this before to watch for updates.
func (ti *TemplateIndex) AddTemplates(ts []templatev1.Template) (int, error) {
	var err error
	var count int

	// unneeded, but better safe than sorry
	ti.rwlock.Lock()
	defer ti.rwlock.Unlock()

	for _, t := range ts {
		err = ti.add(&t)

		if err != nil {
			return count, err
		}
		count += 1
	}
	return count, nil
}

func (ti *TemplateIndex) Update(t *templatev1.Template) error {
	ti.rwlock.Lock()
	defer ti.rwlock.Unlock()

	ti.log.Info(fmt.Sprintf("handling template: %v", t.UID))

	_, ok := ti.templates[t.UID]
	if !ok {
		ti.add(t)
	} else {
		delete(ti.templates, t.UID)
		ti.log.Info(fmt.Sprintf("removed template: %v", t.UID))
	}
	return nil
}

func (ti *TemplateIndex) add(t *templatev1.Template) error {
	ti.templates[t.UID] = *t
	ti.log.Info(fmt.Sprintf("added template: %v", t.UID))
	return nil
}
