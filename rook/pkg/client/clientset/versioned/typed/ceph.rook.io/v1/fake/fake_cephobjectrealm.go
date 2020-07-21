/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	cephrookiov1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeCephObjectRealms implements CephObjectRealmInterface
type FakeCephObjectRealms struct {
	Fake *FakeCephV1
	ns   string
}

var cephobjectrealmsResource = schema.GroupVersionResource{Group: "ceph.rook.io", Version: "v1", Resource: "cephobjectrealms"}

var cephobjectrealmsKind = schema.GroupVersionKind{Group: "ceph.rook.io", Version: "v1", Kind: "CephObjectRealm"}

// Get takes name of the cephObjectRealm, and returns the corresponding cephObjectRealm object, and an error if there is any.
func (c *FakeCephObjectRealms) Get(name string, options v1.GetOptions) (result *cephrookiov1.CephObjectRealm, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(cephobjectrealmsResource, c.ns, name), &cephrookiov1.CephObjectRealm{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cephrookiov1.CephObjectRealm), err
}

// List takes label and field selectors, and returns the list of CephObjectRealms that match those selectors.
func (c *FakeCephObjectRealms) List(opts v1.ListOptions) (result *cephrookiov1.CephObjectRealmList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(cephobjectrealmsResource, cephobjectrealmsKind, c.ns, opts), &cephrookiov1.CephObjectRealmList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &cephrookiov1.CephObjectRealmList{ListMeta: obj.(*cephrookiov1.CephObjectRealmList).ListMeta}
	for _, item := range obj.(*cephrookiov1.CephObjectRealmList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested cephObjectRealms.
func (c *FakeCephObjectRealms) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(cephobjectrealmsResource, c.ns, opts))

}

// Create takes the representation of a cephObjectRealm and creates it.  Returns the server's representation of the cephObjectRealm, and an error, if there is any.
func (c *FakeCephObjectRealms) Create(cephObjectRealm *cephrookiov1.CephObjectRealm) (result *cephrookiov1.CephObjectRealm, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(cephobjectrealmsResource, c.ns, cephObjectRealm), &cephrookiov1.CephObjectRealm{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cephrookiov1.CephObjectRealm), err
}

// Update takes the representation of a cephObjectRealm and updates it. Returns the server's representation of the cephObjectRealm, and an error, if there is any.
func (c *FakeCephObjectRealms) Update(cephObjectRealm *cephrookiov1.CephObjectRealm) (result *cephrookiov1.CephObjectRealm, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(cephobjectrealmsResource, c.ns, cephObjectRealm), &cephrookiov1.CephObjectRealm{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cephrookiov1.CephObjectRealm), err
}

// Delete takes name of the cephObjectRealm and deletes it. Returns an error if one occurs.
func (c *FakeCephObjectRealms) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(cephobjectrealmsResource, c.ns, name), &cephrookiov1.CephObjectRealm{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCephObjectRealms) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(cephobjectrealmsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &cephrookiov1.CephObjectRealmList{})
	return err
}

// Patch applies the patch and returns the patched cephObjectRealm.
func (c *FakeCephObjectRealms) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *cephrookiov1.CephObjectRealm, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(cephobjectrealmsResource, c.ns, name, pt, data, subresources...), &cephrookiov1.CephObjectRealm{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cephrookiov1.CephObjectRealm), err
}
