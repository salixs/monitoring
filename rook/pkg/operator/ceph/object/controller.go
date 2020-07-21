/*
Copyright 2016 The Rook Authors. All rights reserved.

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

package object

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/pkg/capnslog"
	bktclient "github.com/kube-object-storage/lib-bucket-provisioner/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	cephv1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	"github.com/rook/rook/pkg/clusterd"
	cephclient "github.com/rook/rook/pkg/daemon/ceph/client"
	cephconfig "github.com/rook/rook/pkg/daemon/ceph/config"
	"github.com/rook/rook/pkg/operator/ceph/cluster/mon"
	opconfig "github.com/rook/rook/pkg/operator/ceph/config"
	opcontroller "github.com/rook/rook/pkg/operator/ceph/controller"
	"github.com/rook/rook/pkg/util/exec"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	controllerName = "ceph-object-controller"
)

var waitForRequeueIfObjectStoreNotReady = reconcile.Result{Requeue: true, RequeueAfter: 10 * time.Second}

var logger = capnslog.NewPackageLogger("github.com/rook/rook", controllerName)

// List of object resources to watch by the controller
var objectsToWatch = []runtime.Object{
	&corev1.Secret{TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: corev1.SchemeGroupVersion.String()}},
	&v1.Service{TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: v1.SchemeGroupVersion.String()}},
	&appsv1.Deployment{TypeMeta: metav1.TypeMeta{Kind: "Deployment", APIVersion: appsv1.SchemeGroupVersion.String()}},
}

var cephObjectStoreKind = reflect.TypeOf(cephv1.CephObjectStore{}).Name()

// Sets the type meta for the controller main object
var controllerTypeMeta = metav1.TypeMeta{
	Kind:       cephObjectStoreKind,
	APIVersion: fmt.Sprintf("%s/%s", cephv1.CustomResourceGroup, cephv1.Version),
}

// ReconcileCephObjectStore reconciles a cephObjectStore object
type ReconcileCephObjectStore struct {
	client              client.Client
	bktclient           bktclient.Interface
	scheme              *runtime.Scheme
	context             *clusterd.Context
	cephClusterSpec     *cephv1.ClusterSpec
	clusterInfo         *cephconfig.ClusterInfo
	objectStoreChannels map[string]*objectStoreHealth
}

type objectStoreHealth struct {
	stopChan          chan struct{}
	monitoringRunning bool
}

// Add creates a new cephObjectStore Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, context *clusterd.Context) error {
	return add(mgr, newReconciler(mgr, context))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, context *clusterd.Context) reconcile.Reconciler {
	// Add the cephv1 scheme to the manager scheme so that the controller knows about it
	mgrScheme := mgr.GetScheme()
	cephv1.AddToScheme(mgr.GetScheme())

	return &ReconcileCephObjectStore{
		client:              mgr.GetClient(),
		scheme:              mgrScheme,
		context:             context,
		bktclient:           bktclient.NewForConfigOrDie(context.KubeConfig),
		objectStoreChannels: make(map[string]*objectStoreHealth),
	}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	logger.Info("successfully started")

	// Watch for changes on the cephObjectStore CRD object
	err = c.Watch(&source.Kind{Type: &cephv1.CephObjectStore{TypeMeta: controllerTypeMeta}}, &handler.EnqueueRequestForObject{}, opcontroller.WatchControllerPredicate())
	if err != nil {
		return err
	}

	// Watch all other resources
	for _, t := range objectsToWatch {
		err = c.Watch(&source.Kind{Type: t}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &cephv1.CephObjectStore{},
		}, opcontroller.WatchPredicateForNonCRDObject(&cephv1.CephObjectStore{TypeMeta: controllerTypeMeta}, mgr.GetScheme()))
		if err != nil {
			return err
		}
	}

	return nil
}

// Reconcile reads that state of the cluster for a cephObjectStore object and makes changes based on the state read
// and what is in the cephObjectStore.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCephObjectStore) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// workaround because the rook logging mechanism is not compatible with the controller-runtime loggin interface
	reconcileResponse, err := r.reconcile(request)
	if err != nil {
		logger.Errorf("failed to reconcile %v", err)
	}

	return reconcileResponse, err
}

func (r *ReconcileCephObjectStore) reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the cephObjectStore instance
	cephObjectStore := &cephv1.CephObjectStore{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cephObjectStore)
	if err != nil {
		if kerrors.IsNotFound(err) {
			logger.Debug("cephObjectStore resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, errors.Wrap(err, "failed to get cephObjectStore")
	}

	// The CR was just created, initializing status fields
	if cephObjectStore.Status == nil {
		updateStatusPhase(r.client, request.NamespacedName, cephv1.ConditionProgressing)
	}

	// Make sure a CephCluster is present otherwise do nothing
	cephCluster, isReadyToReconcile, cephClusterExists, reconcileResponse := opcontroller.IsReadyToReconcile(r.client, r.context, request.NamespacedName, controllerName)
	if !isReadyToReconcile {
		// This handles the case where the Ceph Cluster is gone and we want to delete that CR
		// We skip the deleteStore() function since everything is gone already
		//
		// Also, only remove the finalizer if the CephCluster is gone
		// If not, we should wait for it to be ready
		// This handles the case where the operator is not ready to accept Ceph command but the cluster exists
		if !cephObjectStore.GetDeletionTimestamp().IsZero() && !cephClusterExists {
			// Remove finalizer
			err := opcontroller.RemoveFinalizer(r.client, cephObjectStore)
			if err != nil {
				return reconcile.Result{}, errors.Wrap(err, "failed to remove finalizer")
			}

			// Return and do not requeue. Successful deletion.
			return reconcile.Result{}, nil
		}

		return reconcileResponse, nil
	}
	r.cephClusterSpec = &cephCluster.Spec

	// Initialize the channel for this object store
	// This allows us to track multiple ObjectStores in the same namespace
	_, ok := r.objectStoreChannels[cephObjectStore.Name]
	if !ok {
		r.objectStoreChannels[cephObjectStore.Name] = &objectStoreHealth{
			stopChan:          make(chan struct{}),
			monitoringRunning: false,
		}
	}

	// Populate clusterInfo
	// Always populate it during each reconcile
	var clusterInfo *cephconfig.ClusterInfo
	if r.cephClusterSpec.External.Enable {
		clusterInfo = mon.PopulateExternalClusterInfo(r.context, request.NamespacedName.Namespace)
	} else {
		clusterInfo, _, _, err = mon.LoadClusterInfo(r.context, request.NamespacedName.Namespace)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to populate cluster info")
		}
	}
	r.clusterInfo = clusterInfo

	// Populate CephVersion
	// This only needed when bootstrapping rgw pods
	// Might improve if we set the version by running a 'ceph version' command on a goroutine and update it every now and then
	if !r.cephClusterSpec.External.Enable {
		currentCephVersion, err := cephclient.LeastUptodateDaemonVersion(r.context, r.clusterInfo.Name, opconfig.MonType)
		if err != nil {
			return reconcile.Result{}, errors.Wrapf(err, "failed to retrieve current ceph %q version", opconfig.MonType)
		}
		r.clusterInfo.CephVersion = currentCephVersion
	}

	// Set a finalizer so we can do cleanup before the object goes away
	err = opcontroller.AddFinalizerIfNotPresent(r.client, cephObjectStore)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to add finalizer")
	}

	// DELETE: the CR was deleted
	if !cephObjectStore.GetDeletionTimestamp().IsZero() {
		logger.Debugf("deleting store %q", cephObjectStore.Name)

		response, ok := r.verifyObjectBucketCleanup(cephObjectStore)
		if !ok {
			// If the object store cannot be deleted, requeue the request for deletion to see if the conditions
			// will eventually be satisfied such as the object buckets being removed
			return response, nil
		}

		// Close the channel to stop the healthcheck of the endpoint
		close(r.objectStoreChannels[cephObjectStore.Name].stopChan)

		// Remove object store from the map
		delete(r.objectStoreChannels, cephObjectStore.Name)

		cfg := clusterConfig{
			context:     r.context,
			store:       cephObjectStore,
			clusterSpec: r.cephClusterSpec,
		}
		err = cfg.deleteStore()
		if err != nil {
			return reconcile.Result{}, errors.Wrapf(err, "failed to delete store %q", cephObjectStore.Name)
		}

		// Remove finalizer
		err = opcontroller.RemoveFinalizer(r.client, cephObjectStore)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to remove finalizer")
		}

		// Return and do not requeue. Successful deletion.
		return reconcile.Result{}, nil
	}

	// validate the store settings
	if err := r.validateStore(cephObjectStore); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "invalid object store %q arguments", cephObjectStore.Name)
	}

	// CREATE/UPDATE
	reconcileResponse, err = r.reconcileCreateObjectStore(cephObjectStore, request.NamespacedName)
	if err != nil {
		return r.setFailedStatus(request.NamespacedName, "failed to create object store deployments", err)
	}

	// Set Ready status, we are done reconciling
	updateStatusPhase(r.client, request.NamespacedName, cephv1.ConditionReady)

	// Return and do not requeue
	logger.Debug("done reconciling")
	return reconcile.Result{}, nil
}

func (r *ReconcileCephObjectStore) reconcileCreateObjectStore(cephObjectStore *cephv1.CephObjectStore, namespacedName types.NamespacedName) (reconcile.Result, error) {
	cfg := clusterConfig{
		context:           r.context,
		clusterInfo:       r.clusterInfo,
		store:             cephObjectStore,
		rookVersion:       r.cephClusterSpec.CephVersion.Image,
		clusterSpec:       r.cephClusterSpec,
		DataPathMap:       opconfig.NewStatelessDaemonDataPathMap(opconfig.RgwType, cephObjectStore.Name, cephObjectStore.Namespace, r.cephClusterSpec.DataDirHostPath),
		client:            r.client,
		scheme:            r.scheme,
		skipUpgradeChecks: r.cephClusterSpec.SkipUpgradeChecks,
	}
	objContext := NewContext(r.context, cephObjectStore.Name, cephObjectStore.Namespace)
	objContext.UID = string(cephObjectStore.UID)

	var serviceIP string
	var err error

	if r.cephClusterSpec.External.Enable {
		// Assign the cephx user to run Ceph commands with
		objContext.RunAsUser = r.clusterInfo.ExternalCred.Username

		logger.Info("reconciling external object store")

		// RECONCILE SERVICE
		logger.Info("reconciling object store service")
		serviceIP, err = cfg.reconcileService(cephObjectStore)
		if err != nil {
			return r.setFailedStatus(namespacedName, "failed to reconcile service", err)
		}

		// RECONCILE ENDPOINTS
		// Always add the endpoint AFTER the service otherwise it will get overridden
		logger.Info("reconciling external object store endpoint")
		err = cfg.reconcileExternalEndpoint(cephObjectStore)
		if err != nil {
			return r.setFailedStatus(namespacedName, "failed to reconcile external endpoint", err)
		}

	} else {
		logger.Info("reconciling object store deployments")

		// Reconcile realm/zonegroup/zone CRs & update their names
		realmName, zoneGroupName, zoneName, reconcileResponse, err := r.reconcileMultisiteCRs(cephObjectStore)
		if err != nil {
			return reconcileResponse, err
		}

		// Reconcile Ceph Zone if Multisite
		if cephObjectStore.Spec.IsMultisite() {
			reconcileResponse, err := r.reconcileCephZone(cephObjectStore, zoneGroupName, realmName)
			if err != nil {
				return reconcileResponse, err
			}
		}

		// RECONCILE SERVICE
		logger.Debug("reconciling object store service")
		serviceIP, err := cfg.reconcileService(cephObjectStore)
		if err != nil {
			return r.setFailedStatus(namespacedName, "failed to reconcile service", err)
		}

		// RECONCILE POOLS
		logger.Info("reconciling object store pools")
		err = createPools(objContext, cephObjectStore.Spec)
		if err != nil {
			return r.setFailedStatus(namespacedName, "failed to create object pools", err)
		}

		// RECONCILE REALM
		logger.Infof("setting multisite settings for object store %q", cephObjectStore.Name)
		err = setMultisite(objContext, serviceIP, cephObjectStore.Spec, realmName, zoneGroupName, zoneName)
		if err != nil {
			return r.setFailedStatus(namespacedName, "failed to configure multisite for object store", err)
		}

		err = cfg.createOrUpdateStore(realmName, zoneGroupName, zoneName)
		if err != nil {
			return reconcile.Result{}, errors.Wrapf(err, "failed to create object store %q", cephObjectStore.Name)
		}
	}

	// Start monitoring
	if cephObjectStore.Spec.HealthCheck.Bucket.Enabled {
		r.startMonitoring(cephObjectStore, objContext, serviceIP, namespacedName)
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileCephObjectStore) reconcileCephZone(store *cephv1.CephObjectStore, zoneGroupName string, realmName string) (reconcile.Result, error) {
	realmArg := fmt.Sprintf("--rgw-realm=%s", realmName)
	zoneGroupArg := fmt.Sprintf("--rgw-zonegroup=%s", zoneGroupName)
	zoneArg := fmt.Sprintf("--rgw-zone=%s", store.Spec.Zone.Name)
	objContext := NewContext(r.context, store.Name, store.Namespace)

	_, err := RunAdminCommandNoRealm(objContext, "zone", "get", realmArg, zoneGroupArg, zoneArg)
	if err != nil {
		if code, ok := exec.ExitStatus(err); ok && code == int(syscall.ENOENT) {
			return waitForRequeueIfObjectStoreNotReady, errors.Wrapf(err, "ceph zone %q not found", store.Spec.Zone.Name)
		} else {
			return waitForRequeueIfObjectStoreNotReady, errors.Wrapf(err, "radosgw-admin zone get failed with code %d", code)
		}
	}

	logger.Infof("Zone %q found in Ceph cluster will include object store %q", store.Spec.Zone.Name, store.Name)
	return reconcile.Result{}, nil
}

func (r *ReconcileCephObjectStore) reconcileMultisiteCRs(cephObjectStore *cephv1.CephObjectStore) (string, string, string, reconcile.Result, error) {
	if cephObjectStore.Spec.IsMultisite() {
		zoneName := cephObjectStore.Spec.Zone.Name
		zone, err := r.context.RookClientset.CephV1().CephObjectZones(cephObjectStore.Namespace).Get(zoneName, metav1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return "", "", "", waitForRequeueIfObjectStoreNotReady, err
			}
			return "", "", "", waitForRequeueIfObjectStoreNotReady, errors.Wrapf(err, "error getting CephObjectZone %q", cephObjectStore.Spec.Zone.Name)
		}
		logger.Infof("CephObjectZone resource %s found", zone.Name)

		zonegroup, err := r.context.RookClientset.CephV1().CephObjectZoneGroups(cephObjectStore.Namespace).Get(zone.Spec.ZoneGroup, metav1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return "", "", "", waitForRequeueIfObjectStoreNotReady, err
			}
			return "", "", "", waitForRequeueIfObjectStoreNotReady, errors.Wrapf(err, "error getting CephObjectZoneGroup %q", zone.Spec.ZoneGroup)
		}
		logger.Infof("CephObjectZoneGroup resource %s found", zonegroup.Name)

		realm, err := r.context.RookClientset.CephV1().CephObjectRealms(cephObjectStore.Namespace).Get(zonegroup.Spec.Realm, metav1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return "", "", "", waitForRequeueIfObjectStoreNotReady, err
			}
			return "", "", "", waitForRequeueIfObjectStoreNotReady, errors.Wrapf(err, "error getting CephObjectRealm %q", zonegroup.Spec.Realm)
		}
		logger.Infof("CephObjectRealm resource %s found", realm.Name)

		return realm.Name, zonegroup.Name, zone.Name, reconcile.Result{}, nil
	}

	return cephObjectStore.Name, cephObjectStore.Name, cephObjectStore.Name, reconcile.Result{}, nil
}

func (r *ReconcileCephObjectStore) verifyObjectBucketCleanup(objectstore *cephv1.CephObjectStore) (reconcile.Result, bool) {
	bktProvisioner := GetObjectBucketProvisioner(r.context, objectstore.Namespace)
	bktProvisioner = strings.Replace(bktProvisioner, "/", "-", -1)
	selector := fmt.Sprintf("bucket-provisioner=%s", bktProvisioner)
	objectBuckets, err := r.bktclient.ObjectbucketV1alpha1().ObjectBuckets().List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		logger.Errorf("failed to delete object store. failed to list buckets for objectstore %q in namespace %q", objectstore.Name, objectstore.Namespace)
		return opcontroller.WaitForRequeueIfFinalizerBlocked, false
	}

	if len(objectBuckets.Items) == 0 {
		logger.Infof("no buckets found for objectstore %q in namespace %q", objectstore.Name, objectstore.Namespace)
		return reconcile.Result{}, true
	}

	bucketNames := make([]string, 0)
	for _, bucket := range objectBuckets.Items {
		bucketNames = append(bucketNames, bucket.Name)
	}

	logger.Errorf("failed to delete object store. buckets for objectstore %q in namespace %q are not cleaned up. remaining buckets: %+v", objectstore.Name, objectstore.Namespace, bucketNames)
	return opcontroller.WaitForRequeueIfFinalizerBlocked, false
}

func (r *ReconcileCephObjectStore) startMonitoring(objectstore *cephv1.CephObjectStore, objContext *Context, serviceIP string, namespacedName types.NamespacedName) {
	// Start monitoring object store
	if r.objectStoreChannels[objectstore.Name].monitoringRunning {
		logger.Debug("external rgw endpoint monitoring go routine already running!")
		return
	}

	// Set the monitoring flag so we don't start more than one go routine
	r.objectStoreChannels[objectstore.Name].monitoringRunning = true

	var port string

	if objectstore.Spec.Gateway.Port != 0 {
		port = strconv.Itoa(int(objectstore.Spec.Gateway.Port))
	} else if objectstore.Spec.Gateway.SecurePort != 0 {
		port = strconv.Itoa(int(objectstore.Spec.Gateway.SecurePort))
	}

	rgwChecker := newBucketChecker(r.context, objContext, serviceIP, port, r.client, namespacedName, &objectstore.Spec.HealthCheck)
	logger.Info("starting rgw healthcheck")
	go rgwChecker.checkObjectStore(r.objectStoreChannels[objectstore.Name].stopChan)
}
