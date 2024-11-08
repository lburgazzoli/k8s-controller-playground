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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlCli "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sync"

	playgroundApi "github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1"
	playgroundCli "github.com/lburgazzoli/k8s-controller-playground/pkg/controller/client"
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

	if err := r.Client.Apply(
		ctx,
		&playgroundApi.Agent{
			TypeMeta: metav1.TypeMeta{
				APIVersion: playgroundApi.SchemeGroupVersion.String(),
				Kind:       "Agent",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
			Spec: playgroundApi.AgentSpec{
				Name: req.Spec.Name,
			},
		},
		ctrlCli.FieldOwner("playground-controller"),
		ctrlCli.ForceOwnership,
	); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Client.ApplyStatus(
		ctx,
		&playgroundApi.Component{
			TypeMeta: metav1.TypeMeta{
				APIVersion: playgroundApi.SchemeGroupVersion.String(),
				Kind:       "Component",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
			Status: playgroundApi.ComponentStatus{
				ObservedGeneration: req.Generation,
				Phase:              "Ready",
				Name:               req.Spec.Name,
			},
		},
		ctrlCli.FieldOwner("playground-controller"),
		ctrlCli.ForceOwnership,
	); err != nil {
		return ctrl.Result{}, err
	}

	r.o.Do(func() {

		u := unstructured.Unstructured{}
		u.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "maistra.io",
			Version: "v1",
			Kind:    "ServiceMeshMember",
		})

		err := r.c.Watch(source.Kind(
			r.m.GetCache(),
			&u,
			handler.TypedEnqueueRequestsFromMapFunc(func(context.Context, *unstructured.Unstructured) []reconcile.Request {
				return []reconcile.Request{{
					NamespacedName: types.NamespacedName{
						Name:      req.Name,
						Namespace: req.Namespace,
					},
				}}
			})),
		)

		if err != nil {
			panic(err)
		}
	})

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
		Build(rec)

	r.c = c

	if err != nil {
		return err
	}

	return nil
}
