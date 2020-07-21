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

package ceph

import (
	"github.com/rook/rook/cmd/rook/rook"
	"github.com/rook/rook/pkg/daemon/admission"
	operator "github.com/rook/rook/pkg/operator/ceph"
	"github.com/spf13/cobra"
)

var (
	admissionCmd = &cobra.Command{
		Use:   "admission-controller",
		Short: "Starts admission controller",
	}
)

func init() {
	admissionCmd.Run = startAdmissionController
}

func startAdmissionController(cmd *cobra.Command, args []string) {
	rook.SetLogLevel()

	rook.LogStartupInfo(admissionCmd.Flags())

	context := rook.NewContext()

	a := admission.New(context, "ceph", operator.ValidateCephResource)
	a.StartServer()
}
