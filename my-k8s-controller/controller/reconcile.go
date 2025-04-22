package controller

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
)

func Reconcile(clientset *kubernetes.Clientset, observer *ConfigObserver) error {
	if observer == nil {
		return nil
	}

	deploy, err := clientset.AppsV1().Deployments(observer.Namespace).Get(context.TODO(), observer.Spec.DeploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting deployment: %v", err)
	}

	replicas := *deploy.Spec.Replicas
	revision := deploy.Generation //deploy.Annotations["deployment.kubernetes.io/revision"]

	if err := updateConfigMap(clientset, observer.Namespace, observer.Spec.VersionConfigMap, "version", fmt.Sprintf("%d", revision)); err != nil {
		return err
	}

	if err := updateConfigMap(clientset, observer.Namespace, observer.Spec.ReplicaConfigMap, "replicas", fmt.Sprintf("%d", replicas)); err != nil {
		return err
	}

	fmt.Printf("[Reconciled] %s -> revision=%d replicas=%d\n", observer.Spec.DeploymentName, revision, replicas)
	return nil
}

func ParseUnstructuredToObserver(u *unstructured.Unstructured) (*ConfigObserver, error) {
	data, err := json.Marshal(u.Object)
	if err != nil {
		return nil, err
	}
	var observer ConfigObserver
	if err := json.Unmarshal(data, &observer); err != nil {
		return nil, err
	}
	return &observer, nil
}
