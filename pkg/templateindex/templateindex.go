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
	"errors"
	"fmt"
	"sync"

	"github.com/go-logr/logr"

	templatev1 "github.com/openshift/api/template/v1"
)

type TemplateIndexer struct {
	rwlock sync.RWMutex
	log    logr.Logger
	// holds the real data
	// TODO: figure out if we are supposed to use UID,
	// or if Name is good enough.
	templates map[string]templatev1.Template
	ledgers   map[string]Ledger
}

func NewTemplateIndexer(log logr.Logger) *TemplateIndexer {
	return &TemplateIndexer{
		log:       log,
		templates: make(map[string]templatev1.Template),
		ledgers:   make(map[string]Ledger),
	}
}

func (ti *TemplateIndexer) Count() int {
	// unneeded, but better safe than sorry
	ti.rwlock.RLock()
	defer ti.rwlock.RUnlock()

	return len(ti.templates)
}

func (ti *TemplateIndexer) AddLedger(name string, ld Ledger) {
	// unneeded, but better safe than sorry
	ti.rwlock.Lock()
	defer ti.rwlock.Unlock()

	ti.ledgers[name] = ld
}

func (ti *TemplateIndexer) SummarizeBy(name string) ([]Summary, error) {
	ti.rwlock.RLock()
	defer ti.rwlock.RUnlock()

	ld, ok := ti.ledgers[name]
	if !ok {
		return []Summary{}, errors.New(fmt.Sprintf("invalid label: %v", name))
	}

	templates := make([]templatev1.Template, 0, len(ti.templates))
	for _, template := range ti.templates {
		templates = append(templates, template)
	}
	return ld.Summarize(templates), nil
}

func (ti *TemplateIndexer) DescribeBy(opts FilterOptions) ([]Description, error) {
	ti.rwlock.RLock()
	defer ti.rwlock.RUnlock()

	descriptions := []Description{}
	for _, template := range ti.templates {
		matched := 0
		for key, value := range opts {
			label := makeLabel(fixLabelKey(key), value)
			if _, ok := template.Labels[label]; ok {
				matched += 1
			} else {
				// TODO: log
			}
		}
		if matched == len(opts) {
			descriptions = append(descriptions, Describe(&template, opts))
		}
	}
	return descriptions, nil
}

// Set the initial state of the index. You must call this before to watch for updates.
func (ti *TemplateIndexer) AddTemplates(ts []templatev1.Template) (int, error) {
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

func (ti *TemplateIndexer) Update(t *templatev1.Template) error {
	ti.rwlock.Lock()
	defer ti.rwlock.Unlock()

	ti.log.Info(fmt.Sprintf("handling template: %v", t.Name))

	_, ok := ti.templates[t.Name]
	if !ok {
		ti.add(t)
	} else {
		ti.remove(t)
	}
	return nil
}

func (ti *TemplateIndexer) add(t *templatev1.Template) error {
	ti.templates[t.Name] = *t
	ti.log.Info(fmt.Sprintf("added template: %v", t.Name))
	return nil
}

func (ti *TemplateIndexer) remove(t *templatev1.Template) error {
	delete(ti.templates, t.Name)
	ti.log.Info(fmt.Sprintf("removed template: %v", t.Name))
	return nil
}
