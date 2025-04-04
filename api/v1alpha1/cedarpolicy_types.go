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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CedarPolicySpec defines the desired state of CedarPolicy
type CedarPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Policy is the CEDAR policy to be enforced
	Policy *Policy `json:"policy"`
}

type Policy struct {
	Principal   string            `json:"principal"`
	Action      string            `json:"action"`
	Resource    string            `json:"resource"`
	Effect      string            `json:"effect"`
	Conditions  []string          `json:"conditions"`
	Annotations map[string]string `json:"annotations"`
}

// CedarPolicyStatus defines the observed state of CedarPolicy
type CedarPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CedarPolicy is the Schema for the cedarpolicies API
type CedarPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CedarPolicySpec   `json:"spec,omitempty"`
	Status CedarPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CedarPolicyList contains a list of CedarPolicy
type CedarPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CedarPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CedarPolicy{}, &CedarPolicyList{})
}
