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

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	grafanav1alpha1 "github.com/mhson281/grafanaOperator/api/v1alpha1"
)

// GrafanaTeamReconciler reconciles a GrafanaTeam object
type GrafanaTeamReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=grafana.tcodelab.com,resources=grafanateams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=grafana.tcodelab.com,resources=grafanateams/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=grafana.tcodelab.com,resources=grafanateams/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *GrafanaTeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the GrafanaTeam instance
	var team grafanav1alpha1.GrafanaTeam
	if err := r.Get(ctx, req.NamespacedName, &team); err != nil {
		if errors.IsNotFound(err) {
			log.Info("A Grafana Team resource has not been found")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get GrafanaTeam")
		return ctrl.Result{}, nil
	}



	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaTeamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&grafanav1alpha1.GrafanaTeam{}).
		Complete(r)
}
