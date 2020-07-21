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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/rook/rook/pkg/apis/edgefs.rook.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// S3Lister helps list S3s.
type S3Lister interface {
	// List lists all S3s in the indexer.
	List(selector labels.Selector) (ret []*v1.S3, err error)
	// S3s returns an object that can list and get S3s.
	S3s(namespace string) S3NamespaceLister
	S3ListerExpansion
}

// s3Lister implements the S3Lister interface.
type s3Lister struct {
	indexer cache.Indexer
}

// NewS3Lister returns a new S3Lister.
func NewS3Lister(indexer cache.Indexer) S3Lister {
	return &s3Lister{indexer: indexer}
}

// List lists all S3s in the indexer.
func (s *s3Lister) List(selector labels.Selector) (ret []*v1.S3, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.S3))
	})
	return ret, err
}

// S3s returns an object that can list and get S3s.
func (s *s3Lister) S3s(namespace string) S3NamespaceLister {
	return s3NamespaceLister{indexer: s.indexer, namespace: namespace}
}

// S3NamespaceLister helps list and get S3s.
type S3NamespaceLister interface {
	// List lists all S3s in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.S3, err error)
	// Get retrieves the S3 from the indexer for a given namespace and name.
	Get(name string) (*v1.S3, error)
	S3NamespaceListerExpansion
}

// s3NamespaceLister implements the S3NamespaceLister
// interface.
type s3NamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all S3s in the indexer for a given namespace.
func (s s3NamespaceLister) List(selector labels.Selector) (ret []*v1.S3, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.S3))
	})
	return ret, err
}

// Get retrieves the S3 from the indexer for a given namespace and name.
func (s s3NamespaceLister) Get(name string) (*v1.S3, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("s3"), name)
	}
	return obj.(*v1.S3), nil
}
