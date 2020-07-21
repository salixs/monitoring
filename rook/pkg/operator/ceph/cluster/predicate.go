/*
Copyright 2020 The Rook Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package cluster to manage a Ceph cluster.
package cluster

import (
	"github.com/rook/rook/pkg/clusterd"
	discoverDaemon "github.com/rook/rook/pkg/daemon/discover"
	"github.com/rook/rook/pkg/operator/k8sutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// predicateForNodeWatcher is the predicate function to trigger reconcile on Node events
func predicateForNodeWatcher(client client.Client, context *clusterd.Context) predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			clientCluster := newClientCluster(client, e.Meta.GetNamespace(), context)
			return clientCluster.onK8sNode(e.Object)
		},

		UpdateFunc: func(e event.UpdateEvent) bool {
			clientCluster := newClientCluster(client, e.MetaNew.GetNamespace(), context)
			return clientCluster.onK8sNode(e.ObjectNew)
		},

		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},

		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}
}

// predicateForHotPlugCMWatcher is the predicate function to trigger reconcile on ConfigMap events (hot-plug)
func predicateForHotPlugCMWatcher(client client.Client) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			isHotPlugCM := isHotPlugCM(e.ObjectNew)
			if !isHotPlugCM {
				logger.Debugf("hot-plug cm watcher: only reconcile on hot plug cm changes, this %q cm is handled by another watcher", e.MetaNew.GetName())
				return false
			}

			clientCluster := newClientCluster(client, e.MetaNew.GetNamespace(), &clusterd.Context{})
			return clientCluster.onDeviceCMUpdate(e.ObjectOld, e.ObjectNew)
		},

		DeleteFunc: func(e event.DeleteEvent) bool {
			// TODO: if the configmap goes away we could retrigger rook-discover DS
			// However at this point the returned bool can only trigger a reconcile of the CephCluster object
			// Definitely non-trivial but nice to have in the future
			return false
		},

		CreateFunc: func(e event.CreateEvent) bool {
			return false
		},

		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}
}

// isHotPlugCM informs whether the object is the cm for hot-plug disk
func isHotPlugCM(obj runtime.Object) bool {
	// If not a ConfigMap, let's not reconcile
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return false
	}

	// Get the labels
	labels := cm.GetLabels()

	labelVal, labelKeyExist := labels[k8sutil.AppAttr]
	if labelKeyExist && labelVal == discoverDaemon.AppName {
		return true
	}

	return false
}
