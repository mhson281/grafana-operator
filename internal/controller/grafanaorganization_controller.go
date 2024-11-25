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
	"bytes"
	"context"
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	grafanav1alpha1 "github.com/mhson281/grafanaOperator/api/v1alpha1"
)

// GrafanaOrganizationReconciler reconciles a GrafanaOrganization object
type GrafanaOrganizationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=grafana.tcodelab.com,resources=grafanaorganizations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=grafana.tcodelab.com,resources=grafanaorganizations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=grafana.tcodelab.com,resources=grafanaorganizations/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *GrafanaOrganizationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("grafanaorganization", req.NamespacedName)

	var org grafanav1alpha1.GrafanaOrganization
	if err := r.Get(ctx, req.NamespacedName, &org); err != nil {
		if errors.IsNotFound(err) {
			// Resource not found, could have been deleted
			log.Info("A GrafanaOrganization resource has not been found.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get GrafanaOrganization")
		return ctrl.Result{}, err
	}

	if org.Status.OrganizationID != 0 {
		log.Info("Organization already exists in Grafana", "OrganizationID", org.Status.OrganizationID)
		return ctrl.Result{}, nil
	}

	grafanaOrgID, err := r.createGrafanaOrg(ctx, org.Spec.Name)
	if err != nil {
		log.Error(err, "Failed to create organization in Grafana")
		return ctrl.Result{}, err
	}

	org.Status.OrganizationID = grafanaOrgID
	if err := r.Status().Update(ctx, &org); err != nil {
		log.Error(err, "Failed to update GrafanaOrganization status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *GrafanaOrganizationReconciler) createGrafanaOrg(ctx context.Context, orgName string) (int, error) {
	apiURL := "http://172.18.0.3:32000/api/orgs"

	var secret corev1.Secret
	secretName := "grafana"
	namespace := "grafana"

	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: secretName}, &secret); err != nil {
		return 0, fmt.Errorf("Failed to fetch secret %s/%s: %v", namespace, secretName, err)
	}

	// extract admin username and pass from secret
	adminName := string(secret.Data["admin-user"])
	adminPassword := string(secret.Data["admin-password"])

	payload := map[string]string{"name": orgName}
	body, _ := json.Marshal(payload)

	// create the http Request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("Failed to send HTTP requeset: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(adminName, adminPassword)

	// Send the Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("Failed to create new organization with: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("Failed to decode response body: %v", err)
	}

	orgID, ok := result["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("Invalid reponse from grafana API: missing org id")
	}

	return int(orgID), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaOrganizationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&grafanav1alpha1.GrafanaOrganization{}).
		Complete(r)
}
