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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HumanSpec defines the desired state of Human
type HumanSpec struct {
	Hands        int32  `json:"hands"`
	Legs         int32  `json:"legs"`
	Tail         int32  `json:"tail"`
	MotherTongue string `json:"mothertongue"`
}

// HumanStatus defines the observed state of Human
type HumanStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Human is the Schema for the humans API
type Human struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HumanSpec   `json:"spec,omitempty"`
	Status HumanStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HumanList contains a list of Human
type HumanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Human `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Human{}, &HumanList{})
}
