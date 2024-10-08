/*
Copyright 2024.

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

package v1alpha1

import (
	v1alpha1 "github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// ComponentLister helps list Components.
// All objects returned here must be treated as read-only.
type ComponentLister interface {
	// List lists all Components in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Component, err error)
	// Components returns an object that can list and get Components.
	Components(namespace string) ComponentNamespaceLister
	ComponentListerExpansion
}

// componentLister implements the ComponentLister interface.
type componentLister struct {
	listers.ResourceIndexer[*v1alpha1.Component]
}

// NewComponentLister returns a new ComponentLister.
func NewComponentLister(indexer cache.Indexer) ComponentLister {
	return &componentLister{listers.New[*v1alpha1.Component](indexer, v1alpha1.Resource("component"))}
}

// Components returns an object that can list and get Components.
func (s *componentLister) Components(namespace string) ComponentNamespaceLister {
	return componentNamespaceLister{listers.NewNamespaced[*v1alpha1.Component](s.ResourceIndexer, namespace)}
}

// ComponentNamespaceLister helps list and get Components.
// All objects returned here must be treated as read-only.
type ComponentNamespaceLister interface {
	// List lists all Components in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Component, err error)
	// Get retrieves the Component from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Component, error)
	ComponentNamespaceListerExpansion
}

// componentNamespaceLister implements the ComponentNamespaceLister
// interface.
type componentNamespaceLister struct {
	listers.ResourceIndexer[*v1alpha1.Component]
}
