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

package playground

import (
	"context"
	"fmt"
	"github.com/lburgazzoli/k8s-controller-playground/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	ctrlCli "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sync"

	playgroundApi "github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1"
	playgroundCli "github.com/lburgazzoli/k8s-controller-playground/pkg/controller/client"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
)

// ComponentReconciler reconciles a Component object
type ComponentReconciler struct {
	Client *playgroundCli.Client
	Scheme *runtime.Scheme

	m ctrl.Manager
	c controller.Controller
	o sync.Once
}

// +kubebuilder:rbac:groups=playground.lburgazzoli.github.io,resources=components,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=playground.lburgazzoli.github.io,resources=components/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=playground.lburgazzoli.github.io,resources=components/finalizers,verbs=update

func (r *ComponentReconciler) Reconcile(ctx context.Context, req *playgroundApi.Component) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.Info("rec")

	cm1 := corev1.ConfigMap{}
	cm1.Name = req.Name
	cm1.Namespace = req.Namespace
	cm1.Labels = map[string]string{
		"foo": "bar",
	}

	if err := ctrl.SetControllerReference(req, &cm1, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	cm1.Data = map[string]string{
		"cm1": "foo",
	}

	if err := Apply(ctx, r.Client, &cm1, ctrlCli.ForceOwnership, ctrlCli.FieldOwner("cm1")); err != nil {
		return ctrl.Result{}, err
	}

	cm1.OwnerReferences = nil
	cm1.Data = map[string]string{
		"cm2": "bar",
	}

	if err := Apply(ctx, r.Client, &cm1, ctrlCli.ForceOwnership, ctrlCli.FieldOwner("cm2")); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ComponentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.m = mgr

	rec := reconcile.AsReconciler(r.Client.Client, r)

	u := unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "maistra.io",
		Version: "v1",
		Kind:    "ServiceMeshMember",
	})

	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&playgroundApi.Component{}).
		Owns(&playgroundApi.Agent{}).
		Owns(
			&corev1.ConfigMap{},
			builder.WithPredicates(predicate.Funcs{
				UpdateFunc: func(e event.TypedUpdateEvent[ctrlCli.Object]) bool {
					ctrl.Log.WithName("cm").Info(">>>>>>>>>>>", "name", e.ObjectNew.GetName(), "rv", e.ObjectNew.GetResourceVersion())
					return true
				},
			}),
		).
		Build(rec)

	r.c = c

	if err != nil {
		return err
	}

	return nil
}

func Apply(ctx context.Context, cli ctrlCli.Client, in ctrlCli.Object, opts ...ctrlCli.PatchOption) error {
	u, err := resources.ToUnstructured(cli.Scheme(), in)
	if err != nil {
		return fmt.Errorf("failed to convert resource to unstructured: %w", err)
	}

	// safe copy
	u = u.DeepCopy()

	// remove not required fields
	unstructured.RemoveNestedField(u.Object, "metadata", "managedFields")
	unstructured.RemoveNestedField(u.Object, "metadata", "resourceVersion")
	unstructured.RemoveNestedField(u.Object, "status")

	err = cli.Patch(ctx, u, ctrlCli.Apply, opts...)
	switch {
	case k8serr.IsNotFound(err):
		return nil
	case err != nil:
		return fmt.Errorf("unable to patch object %s: %w", u, err)
	}

	// Write back the modified object so callers can access the patched object.
	err = cli.Scheme().Convert(u, in, ctx)
	if err != nil {
		return fmt.Errorf("failed to write modified object: %w", err)
	}

	return nil
}
