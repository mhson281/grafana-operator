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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GrafanaTeamSpec defines the desired state of GrafanaTeam
type GrafanaTeamSpec struct {
	Name string `json:"name"`
}

// GrafanaTeamStatus defines the observed state of GrafanaTeam
type GrafanaTeamStatus struct {
	OrganizationID int64 `json:"org_id,omitempty"` // Team will be created associated to OrgID
	TeamID         int64 `json:"team_id,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=gt;team
// GrafanaTeam is the Schema for the grafanateams API
type GrafanaTeam struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrafanaTeamSpec   `json:"spec,omitempty"`
	Status GrafanaTeamStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// GrafanaTeamList contains a list of GrafanaTeam
type GrafanaTeamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GrafanaTeam `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GrafanaTeam{}, &GrafanaTeamList{})
}
