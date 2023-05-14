/*
Copyright 2023.

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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mammalsv1 "github.com/philomathesinc/k8s/api/v1"
	"github.com/philomathesinc/k8s/internal/resources"
	"k8s.io/apimachinery/pkg/api/errors"
)

// HumanReconciler reconciles a Human object
type HumanReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=mammals.example.com,resources=humans,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mammals.example.com,resources=humans/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mammals.example.com,resources=humans/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=create;list;get;watch

func (r *HumanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// A variable to hold the CR data
	var human mammalsv1.Human

	// Call kubernetes API to get the CR data
	if err := r.Get(ctx, req.NamespacedName, &human); err != nil {
		log.Error(err, "unable to fetch Human")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var found corev1.Pod
	key := types.NamespacedName{
		Name:      human.Name + "-pod",
		Namespace: human.Namespace,
	}
	// Call kubernetes API to get the pod
	if err := r.Get(ctx, key, &found); err != nil {
		if errors.IsNotFound(err) {
			pod := resources.PodForHuman(&human, r.Scheme, log)
			// call k8s api to create the pod
			if err := r.Create(ctx, pod); err != nil {
				log.Error(err, "unable to create pod")
			}
			log.Info("Pod created", "name", key.Name, "namespace", key.Namespace)
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to get pod")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HumanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mammalsv1.Human{}).
		Complete(r)
}
