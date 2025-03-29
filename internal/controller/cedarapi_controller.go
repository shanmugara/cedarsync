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
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	cedarsyncv1alpha1 "github.com/shanmugara/cedarsync/api/v1alpha1"
)

const (
	apiFinalizer = "cedarapi.cedarsync.omegahome.net/finalizer"
	requeueAfter = time.Second * 15
)

// ReconcileFn is the fucntion signature for the Reconcile function
type ReconcileFn func(ctx context.Context, capi *cedarsyncv1alpha1.CedarApi) error

// CedarPolicyReconciler reconciles a CedarPolicy object
type CedarPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cedarsync.omegahome.net,resources=cedarapis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cedarsync.omegahome.net,resources=cedarapis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cedarsync.omegahome.net,resources=cedarapis/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CedarPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *CedarPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling CedarPolicy")

	// First fetch the CedarApi instance
	cedarApi := &cedarsyncv1alpha1.CedarApi{}
	if err := r.Get(ctx, req.NamespacedName, cedarApi); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("CedarApi not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to fetch CedarApi")
		return ctrl.Result{}, err
	}

	// Check if the CedarApi instance is being deleted
	if !cedarApi.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("CedarApi is being deleted")

		// Delete all the CedarPolicy instances
		for _, step := range []struct {
			name string
			fn   ReconcileFn
		}{
			{"DeletePolicy", r.DeletePolicy},
		} {
			logger.Info("Deleting step", "step", step.name)
			if err := step.fn(ctx, cedarApi); err != nil {
				logger.Error(err, "Failed to delete step", "step", step.name)
				return ctrl.Result{Requeue: true}, err
			}
		}
		// Remove finalizer from the CedarApi instance
		if controllerutil.ContainsFinalizer(cedarApi, apiFinalizer) {
			controllerutil.RemoveFinalizer(cedarApi, apiFinalizer)
			if err := r.Update(ctx, cedarApi); err != nil {
				logger.Error(err, "Failed to remove finalizer")
				return ctrl.Result{Requeue: true}, err
			}
		}

		return ctrl.Result{}, nil

	} else {
		// Add finalizer to the CedarApi instance
		if !controllerutil.ContainsFinalizer(cedarApi, apiFinalizer) {
			controllerutil.AddFinalizer(cedarApi, apiFinalizer)
			if err := r.Update(ctx, cedarApi); err != nil {
				logger.Error(err, "Failed to add finalizer")
				return ctrl.Result{}, err
			}
		}
	}

	// Reconcile loop
	for _, step := range []struct {
		name string
		fn   ReconcileFn
	}{
		{"ReconcilePolicy", r.ReconcilePolicy},
	} {
		logger.Info("Reconciling step", "step", step.name)
		if err := step.fn(ctx, cedarApi); err != nil {
			logger.Error(err, "Failed to reconcile step", "step", step.name)
			return ctrl.Result{Requeue: true}, err
		}
	}

	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CedarPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cedarsyncv1alpha1.CedarApi{}).
		Owns(&cedarsyncv1alpha1.CedarPolicy{}).
		Complete(r)
}
