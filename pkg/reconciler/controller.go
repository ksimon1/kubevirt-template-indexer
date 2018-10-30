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

package reconciler

import (
	"context"

	"github.com/go-logr/logr"

	templatev1 "github.com/openshift/api/template/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/fromanirh/kubevirt-template-indexer/pkg/templateindex"
)

type TemplateReconciler struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
	log    logr.Logger
	index  *templateindex.TemplateIndex
}

func NewTemplateReconciler(client client.Client, log logr.Logger, index *templateindex.TemplateIndex) *TemplateReconciler {
	return &TemplateReconciler{
		client: client,
		log:    log,
		index:  index,
	}
}

func (tr *TemplateReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// set up a convinient log object so we don't have to type request over and over again
	log := tr.log.WithValues("request", request)

	// Fetch the Template from the cache
	t := &templatev1.Template{}
	err := tr.client.Get(context.TODO(), request.NamespacedName, t)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find Template")
		return reconcile.Result{}, nil
	}

	if err != nil {
		log.Error(err, "Could not fetch Template")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, tr.index.Update(t)
}
