package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigObserverSpec defines the desired state of the custom resource
type ConfigObserverSpec struct {
	DeploymentName   string `json:"deploymentName"`
	VersionConfigMap string `json:"versionConfigMap"`
	ReplicaConfigMap string `json:"replicaConfigMap"`
}

// ConfigObserver is the Schema for the custom resource
type ConfigObserver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ConfigObserverSpec `json:"spec,omitempty"`
}

// ConfigObserverList contains a list of ConfigObserver
type ConfigObserverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigObserver `json:"items"`
}
