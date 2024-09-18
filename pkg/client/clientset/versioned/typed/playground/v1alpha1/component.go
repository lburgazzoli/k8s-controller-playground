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
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1alpha1 "github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1"
	playgroundv1alpha1 "github.com/lburgazzoli/k8s-controller-playground/pkg/client/applyconfiguration/playground/v1alpha1"
	scheme "github.com/lburgazzoli/k8s-controller-playground/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ComponentsGetter has a method to return a ComponentInterface.
// A group's client should implement this interface.
type ComponentsGetter interface {
	Components(namespace string) ComponentInterface
}

// ComponentInterface has methods to work with Component resources.
type ComponentInterface interface {
	Create(ctx context.Context, component *v1alpha1.Component, opts v1.CreateOptions) (*v1alpha1.Component, error)
	Update(ctx context.Context, component *v1alpha1.Component, opts v1.UpdateOptions) (*v1alpha1.Component, error)
	UpdateStatus(ctx context.Context, component *v1alpha1.Component, opts v1.UpdateOptions) (*v1alpha1.Component, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Component, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ComponentList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Component, err error)
	Apply(ctx context.Context, component *playgroundv1alpha1.ComponentApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Component, err error)
	ApplyStatus(ctx context.Context, component *playgroundv1alpha1.ComponentApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Component, err error)
	ComponentExpansion
}

// components implements ComponentInterface
type components struct {
	client rest.Interface
	ns     string
}

// newComponents returns a Components
func newComponents(c *PlaygroundV1alpha1Client, namespace string) *components {
	return &components{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the component, and returns the corresponding component object, and an error if there is any.
func (c *components) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Component, err error) {
	result = &v1alpha1.Component{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("components").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Components that match those selectors.
func (c *components) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ComponentList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ComponentList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("components").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested components.
func (c *components) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("components").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a component and creates it.  Returns the server's representation of the component, and an error, if there is any.
func (c *components) Create(ctx context.Context, component *v1alpha1.Component, opts v1.CreateOptions) (result *v1alpha1.Component, err error) {
	result = &v1alpha1.Component{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("components").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(component).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a component and updates it. Returns the server's representation of the component, and an error, if there is any.
func (c *components) Update(ctx context.Context, component *v1alpha1.Component, opts v1.UpdateOptions) (result *v1alpha1.Component, err error) {
	result = &v1alpha1.Component{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("components").
		Name(component.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(component).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *components) UpdateStatus(ctx context.Context, component *v1alpha1.Component, opts v1.UpdateOptions) (result *v1alpha1.Component, err error) {
	result = &v1alpha1.Component{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("components").
		Name(component.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(component).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the component and deletes it. Returns an error if one occurs.
func (c *components) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("components").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *components) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("components").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched component.
func (c *components) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Component, err error) {
	result = &v1alpha1.Component{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("components").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied component.
func (c *components) Apply(ctx context.Context, component *playgroundv1alpha1.ComponentApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Component, err error) {
	if component == nil {
		return nil, fmt.Errorf("component provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	name := component.Name
	if name == nil {
		return nil, fmt.Errorf("component.Name must be provided to Apply")
	}
	result = &v1alpha1.Component{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("components").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *components) ApplyStatus(ctx context.Context, component *playgroundv1alpha1.ComponentApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Component, err error) {
	if component == nil {
		return nil, fmt.Errorf("component provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}

	name := component.Name
	if name == nil {
		return nil, fmt.Errorf("component.Name must be provided to Apply")
	}

	result = &v1alpha1.Component{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("components").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
