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

package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	configobserverv1 "github.com/ajaypp123/k8s-golang-projects/kubebuilder-configobserver/api/v1"
)

// ConfigObserverReconciler reconciles a ConfigObserver object
type ConfigObserverReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=configobserver.example.com,resources=configobservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=configobserver.example.com,resources=configobservers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=configobserver.example.com,resources=configobservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigObserver object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *ConfigObserverReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// TODO(user): your logic here
	// 1. Fetch the ConfigObserver CR
	var observer configobserverv1.ConfigObserver
	if err := r.Get(ctx, req.NamespacedName, &observer); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Fetch the Deployment
	var deploy appsv1.Deployment
	if err := r.Get(ctx, types.NamespacedName{Name: observer.Spec.DeploymentName, Namespace: req.Namespace}, &deploy); err != nil {
		log.Error(err, "unable to get Deployment")
		return ctrl.Result{}, err
	}

	replicas := *deploy.Spec.Replicas
	revision := deploy.Generation //deploy.Annotations["deployment.kubernetes.io/revision"]

	// 3. Update Replica ConfigMap
	if err := r.updateConfigMap(ctx, req.Namespace, observer.Spec.ReplicaConfigMap, "replicas", fmt.Sprintf("%d", replicas)); err != nil {
		log.Error(err, "failed to update replica ConfigMap")
	}

	// 4. Update Revision ConfigMap
	if err := r.updateConfigMap(ctx, req.Namespace, observer.Spec.VersionConfigMap, "version", fmt.Sprintf("%d", revision)); err != nil {
		log.Error(err, "failed to update version ConfigMap")
	}

	return ctrl.Result{RequeueAfter: time.Second * 10}, nil // Requeue periodically
}

func (r *ConfigObserverReconciler) updateConfigMap(ctx context.Context, namespace, name, key, value string) error {
	var cm corev1.ConfigMap
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, &cm)
	if err != nil {
		return err
	}

	if cm.Data == nil {
		cm.Data = map[string]string{}
	}
	cm.Data[key] = value

	return r.Update(ctx, &cm)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigObserverReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configobserverv1.ConfigObserver{}).
		Named("configobserver").
		Complete(r)
}
