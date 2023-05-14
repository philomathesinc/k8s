package resources

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/go-logr/logr"
	mammalsv1 "github.com/philomathesinc/k8s/api/v1"
)

func PodForHuman(human *mammalsv1.Human, scheme *runtime.Scheme, log logr.Logger) *corev1.Pod {
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
					Name:    "default",
					Image:   "busybox",
					Command: []string{"sh", "-c"},
					Args: []string{
						fmt.Sprintf("echo %s && sleep 3600", message),
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(human, pod, scheme); err != nil {
		log.Error(err, "unable to set controller reference")
	}

	return pod
}
