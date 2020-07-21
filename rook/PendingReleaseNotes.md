# Major Themes

## Action Required

## Notable Features

### Ceph

- Added a [toolbox job](Documentation/ceph-toolbox.md#toolbox-job) for running a script with Ceph commands, similar to running commands in the Rook toolbox.
- Ceph RBD Mirror daemon has been extracted to its own CRD, it has been removed from the `CephCluster` CRD, see the [rbd-mirror crd](Documentation/ceph-rbd-mirror-crd.html).
- CephCluster CRD has been converted to use the controller-runtime framework.
- CephBlockPool CRD has a new field called `parameters` which allows to set any property on a given [pool](Documentation/ceph-pool-crd.html#add-specific-pool-properties)
- OSD changes:
  - OSD on PVC now supports multipath device.
- Added [admission controller](Documentation/admission-controller-usage.md) support for CRD validations.
  - Support for Ceph CRDs is provided. Some validations for CephClusters are included and additional validations can be added for other CRDs
  - Can be extended to add support for other providers
- OBC changes:
  - Updated lib bucket provisioner version to support multithread and required change can be found in [operator.yaml](cluster/examples/kubernetes/ceph/operator.yaml#L449)
    - Can be extended to add support for other providers
- The Rook operator reflects the health of the CephObjectStore in its status field
- The CephObjectStore CR supports connecting to external Ceph Rados Gateways, refer to the [external object section](Documentation/ceph-object.html#connect-to-external-object-store)
- The CephObjectStore CR runs health checks on the object store endpoint, refer to the [health check section](Documentation/ceph-object-store-crd.html#health-settings)

### EdgeFS

### YugabyteDB

### Cassandra

- Updated Base image from Alpine 3.8 to 3.12 due to CVEs.

## Breaking Changes

### Ceph

- rbd-mirror daemons that were deployed through the CephCluster CR won't be managed anymore as they have their own CRD now.
To transition, you can inject the new rbd mirror CR with the desired `count` of daemons and delete the previously managed rbd mirror deployments manually.


## Known Issues

### <Storage Provider>

## Deprecations

### <Storage Provider>
