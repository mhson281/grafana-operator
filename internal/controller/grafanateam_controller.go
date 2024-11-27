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
	"encoding/json"
	"fmt"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	grafanav1alpha1 "github.com/mhson281/grafanaOperator/api/v1alpha1"
)

// GrafanaTeamReconciler reconciles a GrafanaTeam object
type GrafanaTeamReconciler struct {
	client.Client
	Log    logr.Logger
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

	// Check if the organization already exists in Grafana
	if team.Status.TeamID != 0 {
		log.Info("Team already exists in Grafana", "TeamID", team.Status.TeamID)
		return ctrl.Result{}, nil
	}

	var org grafanav1alpha1.GrafanaOrganization
	if err := r.Get(ctx, client.ObjectKey{Name: team.Name, Namespace: team.Namespace}, &org); err != nil {
		if errors.IsNotFound(err) {
			log.Error(err, "Parent organization for GrafanaTeam does not exist")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get organization")
		return ctrl.Result{}, nil
	}

	// User the org_id from the organization's status to create the team
	if org.Status.OrganizationID == 0 {
		log.Error(nil, "GrafanaOrganization does not have an OrganizationID assigned")
		return ctrl.Result{}, nil
	}
	orgID := org.Status.OrganizationID

	// Call the Grafana API to create the team with the same name/org_id
  teamID, err := r.createGrafanaTeam(team.Spec.Name, orgID)
	if err != nil {
		log.Error(err, "Failed to create team in Grafana")
	}

	team.Status.TeamID = teamID
	if err := r.Status().Update(ctx, &team); err != nil {
		log.Error(err, "Failed to update GrafanaTeam status")
	}

	log.Info("Successfully created GrafanaTeam", "TeamID", teamID)

	return ctrl.Result{}, nil
}

func (r *GrafanaOrganizationReconciler) createGrafanaTeam(teamName string, orgID int64) (int64, error) {
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

	payload := map[string]interface{}{
		"name": teamName,
		"orgId": orgID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("Failed to marshal request payload: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("Failed to create HTTP request: %v", err)
	}

	// Set headers and basic auth
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(adminName, adminPassword)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Failed to send HTTP request: %v", err)
	}

	// Check response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("Failed to create new team: %s", resp.Status) 
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("Failed to parse response body, err")
	}

	teamID, ok := result["teamId"].(float64)
	if !ok {
		return 0, fmt.Errorf("Invalid response from Grafana API: missing or invalid teamID")
	}

	return int64(teamID), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaTeamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&grafanav1alpha1.GrafanaTeam{}).
		Complete(r)
}
