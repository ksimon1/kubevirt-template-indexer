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

package main

import (
	"os"

	flag "github.com/spf13/pflag"

	templatev1 "github.com/openshift/api/template/v1"

	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/fromanirh/kubevirt-template-indexer/pkg/reconciler"
	"github.com/fromanirh/kubevirt-template-indexer/pkg/templateindex"
)

// AddToSchemes may be used to add all resources defined in the project to a Scheme
var AddToSchemes runtime.SchemeBuilder

// AddToScheme adds all Resources to the Scheme
func AddToScheme(s *runtime.Scheme) error {
	return AddToSchemes.AddToScheme(s)
}

var log = logf.Log.WithName("kubevirt-template-indexer")

func main() {
	develMode := flag.BoolP("develmode", "D", false, "enable development mode (more logs)")
	startupSync := flag.BoolP("skipsync", "s", true, "skip initial sync with cluster")
	namespace := flag.StringP("namespace", "N", "", "restrict namespace to watch (default: all)")
	flag.Parse()

	logf.SetLogger(logf.ZapLogger(*develMode))
	entryLog := log.WithName("entrypoint")

	index := templateindex.NewTemplateIndex(log.WithName("indexer"))

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{Namespace: *namespace})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	entryLog.Info("registering Components")

	// Setup Scheme for all resources
	if err := AddToScheme(mgr.GetScheme()); err != nil {
		entryLog.Error(err, "unable to set up the basic scheme")
		os.Exit(1)
	}

	// Setup Scheme for templates
	if err := templatev1.AddToScheme(mgr.GetScheme()); err != nil {
		entryLog.Error(err, "unable to set up the template scheme")
		os.Exit(1)
	}

	entryLog.Info("setting up reconciler")
	tr := reconciler.NewTemplateReconciler(mgr.GetClient(), log.WithName("reconciler"), index)

	entryLog.Info("setting up controller")
	c, err := controller.New("foo-controller", mgr, controller.Options{
		Reconciler: tr,
	})
	if err != nil {
		entryLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}

	if *startupSync {
		entryLog.Info("syncing reconciler")
		err = tr.SyncWithCluster(*namespace)
		if err != nil {
			entryLog.Error(err, "unable to sync with cluster")
			os.Exit(1)
		}
	}

	if err := c.Watch(&source.Kind{Type: &templatev1.Template{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "unable to watch Templates")
		os.Exit(1)
	}

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
