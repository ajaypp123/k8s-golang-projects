/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConfigObserverSpec defines the desired state of ConfigObserver.
type ConfigObserverSpec struct {
	DeploymentName   string `json:"deploymentName"`
	VersionConfigMap string `json:"versionConfigMap"`
	ReplicaConfigMap string `json:"replicaConfigMap"`
}

// ConfigObserverStatus defines the observed state of ConfigObserver.
type ConfigObserverStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ConfigObserver is the Schema for the configobservers API.
type ConfigObserver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigObserverSpec   `json:"spec,omitempty"`
	Status ConfigObserverStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ConfigObserverList contains a list of ConfigObserver.
type ConfigObserverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigObserver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigObserver{}, &ConfigObserverList{})
}
