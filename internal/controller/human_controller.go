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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mammalsv1 "github.com/philomathesinc/k8s/api/v1"
)

// HumanReconciler reconciles a Human object
type HumanReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=mammals.example.com,resources=humans,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mammals.example.com,resources=humans/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mammals.example.com,resources=humans/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=create

func (r *HumanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// A variable to hold the CR data
	var human mammalsv1.Human

	// Call ubernetes API to get the CR data
	if err := r.Get(ctx, req.NamespacedName, &human); err != nil {
		log.Error(err, "unable to fetch Human")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// create the pod struct
	getPodForHuman := func(human *mammalsv1.Human) *corev1.Pod {
		message := fmt.Sprintf("%s has %d legs, %d hands and %d tails. Also, %s speaks in %s",
			human.Name,
			human.Spec.Legs,
			human.Spec.Hands,
			human.Spec.Tail,
			human.Name,
			human.Spec.MotherTongue)

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      human.Name + "-pod",
				Namespace: human.Namespace,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "default",
						Image: "busybox",
						Args: []string{
							"echo",
							fmt.Sprintf("%q", message),
						},
					},
				},
			},
		}

		if err := ctrl.SetControllerReference(human, pod, r.Scheme); err != nil {
			log.Error(err, "unable to set controller reference")
		}

		return pod
	}

	pod := getPodForHuman(&human)

	// call k8s api to create the pod
	if err := r.Create(ctx, pod); err != nil {
		log.Error(err, "unable to create pod")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HumanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mammalsv1.Human{}).
		Complete(r)
}
