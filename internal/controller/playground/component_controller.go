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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	playgroundApi "github.com/lburgazzoli/k8s-controller-playground/api/playground/v1alpha1"
	playgroundAc "github.com/lburgazzoli/k8s-controller-playground/pkg/client/applyconfiguration/playground/v1alpha1"
	playgroundCli "github.com/lburgazzoli/k8s-controller-playground/pkg/controller/client"
)

// ComponentReconciler reconciles a Component object
type ComponentReconciler struct {
	Client *playgroundCli.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=playground.lburgazzoli.github.io,resources=components,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=playground.lburgazzoli.github.io,resources=components/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=playground.lburgazzoli.github.io,resources=components/finalizers,verbs=update

func (r *ComponentReconciler) Reconcile(ctx context.Context, req *playgroundApi.Component) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.Info("rec")

	_, err := r.Client.PlaygroundV1alpha1().Components(req.Namespace).ApplyStatus(
		ctx,
		playgroundAc.Component(req.Name, req.Namespace).
			WithStatus(playgroundAc.ComponentStatus().
				WithObservedGeneration(req.Generation).
				WithPhase("Ready"),
			),
		metav1.ApplyOptions{
			FieldManager: "playground-controller",
			Force:        true,
		},
	)

	if err != nil {
		return ctrl.Result{}, err
	}

	// playgroundClient

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ComponentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	rec := reconcile.AsReconciler(r.Client.Client, r)

	return ctrl.NewControllerManagedBy(mgr).
		For(&playgroundApi.Component{}).
		Complete(rec)
}
