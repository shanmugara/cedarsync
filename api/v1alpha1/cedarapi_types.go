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

// CedarApiSpec defines the desired state of CedarApi
type CedarApiSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// cluster is the kubernetes cluster to be managed
	Cluster string `json:"cluster"`
	// apiUrl is the URL:port of the CEDAR API to fetch the policy from
	ApiUrl string `json:"apiUrl"`
}

// CedarApiStatus defines the observed state of CedarApi
type CedarApiStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CedarApi is the Schema for the cedarapis API
type CedarApi struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CedarApiSpec   `json:"spec,omitempty"`
	Status CedarApiStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CedarApiList contains a list of CedarApi
type CedarApiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CedarApi `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CedarApi{}, &CedarApiList{})
}
