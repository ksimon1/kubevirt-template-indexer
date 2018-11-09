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
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	templatev1 "github.com/openshift/api/template/v1"

	"github.com/fromanirh/kubevirt-template-indexer/pkg/reconciler"
	"github.com/fromanirh/kubevirt-template-indexer/pkg/templateindex"

	_ "github.com/fromanirh/kubevirt-template-indexer/pkg/okd"

	"github.com/fromanirh/kubevirt-template-indexer/internal/pkg/routes"
)

func zapLogger(development bool) logr.Logger {
	sink := zapcore.AddSync(os.Stderr)

	var encCfg zapcore.EncoderConfig
	var lvl zap.AtomicLevel

	opts := []zap.Option{
		zap.Development(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(1),
		zap.ErrorOutput(sink),
	}

	if development {
		encCfg = zap.NewDevelopmentEncoderConfig()
		lvl = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		encCfg = zap.NewProductionEncoderConfig()
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	enc := zapcore.NewConsoleEncoder(encCfg)
	log := zap.New(zapcore.NewCore(enc, sink, lvl))
	log = log.WithOptions(opts...)
	return zapr.NewLogger(log)
}

var log = logf.Log.WithName("kubevirt-template-indexer")

type ledgerDesc struct {
	Name  string
	Label string
}

func main() {
	develMode := flag.BoolP("develmode", "D", false, "enable development mode (more logs)")
	startupSync := flag.BoolP("skipsync", "s", true, "skip initial sync with cluster")
	namespace := flag.StringP("namespace", "N", "", "restrict namespace to watch (default: all)")
	iface := flag.StringP("interface", "I", "", "listen only on this interface for HTTP queries (default: all)")
	port := flag.IntP("port", "p", 8080, "listen on port for HTTP queries (default: 8080)")
	configDir := flag.StringP("confdir", "C", "/etc/template-index", "base directory for the config map files")
	flag.Parse()

	logf.SetLogger(zapLogger(*develMode))
	entryLog := log.WithName("entrypoint")

	index := templateindex.NewTemplateIndexer(log.WithName("indexer"))

	descs := []ledgerDesc{
		ledgerDesc{
			Name:  "os",
			Label: "os",
		},
		ledgerDesc{
			Name:  "workload",
			Label: "workload",
		},
		ledgerDesc{
			Name:  "size",
			Label: "flavor",
		},
	}
	for _, desc := range descs {
		ld := templateindex.NewJSONLedger(desc.Label)

		confPath := filepath.Join(*configDir, desc.Name)
		err := ld.ReadNameMap(confPath)
		if err != nil {
			entryLog.Error(err, fmt.Sprintf("unable read name map %s for ledger %->%s: %s", confPath, desc.Name, desc.Label, err))
			// we can carry on with less data
		}

		index.AddLedger(desc.Name, ld)
		entryLog.Info(fmt.Sprintf("added ledger %s for label=%s", desc.Name, desc.Label))
	}

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{Namespace: *namespace})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
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

	entryLog.Info("starting HTTP endpoints")
	go routes.Serve(*iface, *port, index, log.WithName("httpapi"))

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
