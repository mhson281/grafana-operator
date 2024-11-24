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

// GrafanaOrganizationSpec defines the desired state of GrafanaOrganization
type GrafanaOrganizationSpec struct {
	Name string `json:"name"`
}

// GrafanaOrganizationStatus defines the observed state of GrafanaOrganization
type GrafanaOrganizationStatus struct {
	OrganizationID int `json:"organizationID,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=gorg;org
// GrafanaOrganization is the Schema for the grafanaorganizations API
type GrafanaOrganization struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrafanaOrganizationSpec   `json:"spec,omitempty"`
	Status GrafanaOrganizationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// GrafanaOrganizationList contains a list of GrafanaOrganization
type GrafanaOrganizationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GrafanaOrganization `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GrafanaOrganization{}, &GrafanaOrganizationList{})
}
