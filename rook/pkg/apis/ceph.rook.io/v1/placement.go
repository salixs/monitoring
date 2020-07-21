/*
Copyright 2018 The Rook Authors. All rights reserved.

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
package v1

import (
	rookv1 "github.com/rook/rook/pkg/apis/rook.io/v1"
)

// GetMgrPlacement returns the placement for the MGR service
func GetMgrPlacement(p rookv1.PlacementSpec) rookv1.Placement {
	return p.All().Merge(p[KeyMgr])
}

// GetMonPlacement returns the placement for the MON service
func GetMonPlacement(p rookv1.PlacementSpec) rookv1.Placement {
	return p.All().Merge(p[KeyMon])
}

// GetOSDPlacement returns the placement for the OSD service
func GetOSDPlacement(p rookv1.PlacementSpec) rookv1.Placement {
	return p.All().Merge(p[KeyOSD])
}

// GetCleanupPlacement returns the placement the cleanup job
func GetCleanupPlacement(p rookv1.PlacementSpec) rookv1.Placement {
	return p.All().Merge(p[KeyCleanup])
}
