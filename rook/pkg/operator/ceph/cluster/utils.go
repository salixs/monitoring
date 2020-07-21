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
	"github.com/pkg/errors"
	"github.com/rook/rook/pkg/clusterd"
	"github.com/rook/rook/pkg/operator/k8sutil"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *cluster) setOrchestrationNeeded() {
	c.orchMux.Lock()
	c.orchestrationNeeded = true
	c.orchMux.Unlock()
}

// unsetOrchestrationStatus resets the orchestrationRunning-flag
func (c *cluster) unsetOrchestrationStatus() {
	c.orchMux.Lock()
	defer c.orchMux.Unlock()
	c.orchestrationRunning = false
}

// checkSetOrchestrationStatus is responsible to do orchestration as long as there is a request needed
func (c *cluster) checkSetOrchestrationStatus() bool {
	c.orchMux.Lock()
	defer c.orchMux.Unlock()

	// check if there is an orchestration needed currently
	if c.orchestrationNeeded == true && c.orchestrationRunning == false {
		// there is an orchestration needed
		// allow to enter the orchestration-loop
		c.orchestrationNeeded = false
		c.orchestrationRunning = true
		return true
	}

	return false
}

// populateConfigOverrideConfigMap creates the "rook-config-override" config map
// Its content allows modifying Ceph configuration flags
func populateConfigOverrideConfigMap(context *clusterd.Context, namespace string, ownerRef metav1.OwnerReference) error {
	placeholderConfig := map[string]string{
		k8sutil.ConfigOverrideVal: "",
	}

	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: k8sutil.ConfigOverrideName,
		},
		Data: placeholderConfig,
	}

	k8sutil.SetOwnerRef(&cm.ObjectMeta, &ownerRef)
	_, err := context.Clientset.CoreV1().ConfigMaps(namespace).Create(cm)
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "failed to create override configmap %s", namespace)
	}

	return nil
}
