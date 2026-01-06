/*
Copyright 2025 The Crossplane Authors.

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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	xpv2 "github.com/crossplane/crossplane-runtime/v2/apis/common/v2"
)

// QualityGateParameters represent the desired state of a QualityGate.
type QualityGateParameters struct {
	// Name is the Display name of the Quality Gate.
	// +kubebuilder:validation:MaxLength=100
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Default indicates whether this Quality Gate is the default one.
	// WARNING: It is currently not possible to unset the default Quality Gate in SonarQube once it is set. The only way to change the default Quality Gate is to set another Quality Gate as default.
	// +kubebuilder:validation:Optional
	Default *bool `json:"default,omitempty"`
}

// QualityGateObservation are the observable fields of a QualityGate.
type QualityGateObservation struct {
	// Actions represents the actions that can be performed on the Quality Gate.
	Actions QualityGatesActions `json:"actions"`
	// Defines the Clean as You Code status of the Quality Gate.
	CaycStatus string `json:"caycStatus"`
	// Conditions represents the list of conditions associated with the Quality Gate.
	Conditions []QualityGateConditionObservation `json:"conditions"`
	// IsAiCodeSupported indicates whether AI Code Assurance is supported for the Quality Gate.
	IsAiCodeSupported bool `json:"isAiCodeSupported"`
	// IsBuiltIn indicates whether the Quality Gate is built-in.
	IsBuiltIn bool `json:"isBuiltIn"`
	// IsDefault indicates whether the Quality Gate is the default one.
	IsDefault bool `json:"isDefault"`
	// Name represents the name of the Quality Gate.
	Name string `json:"name"`
}

// A QualityGateSpec defines the desired state of a QualityGate.
type QualityGateSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	// ForProvider represents the desired state of the Quality Gate.
	ForProvider QualityGateParameters `json:"forProvider"`
}

// A QualityGateStatus represents the observed state of a QualityGate.
type QualityGateStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	// AtProvider represents the observed state of the Quality Gate.
	AtProvider QualityGateObservation `json:"atProvider,omitempty"`
}

// QualityGatesActions represents the actions that can be performed on a Quality Gate.
type QualityGatesActions struct {
	// AssociateProjects defines whether projects can be associated with the Quality Gate.
	AssociateProjects bool `json:"associateProjects"`
	// Copy defines whether the Quality Gate can be copied.
	Copy bool `json:"copy"`
	// Delegate defines whether the Quality Gate can be delegated.
	Delegate bool `json:"delegate"`
	// Delete defines whether the Quality Gate can be deleted.
	Delete bool `json:"delete"`
	// ManageAiCodeAssurance defines whether AI Code Assurance settings can be managed.
	ManageAiCodeAssurance bool `json:"manageAiCodeAssurance"`
	// ManageConditions defines whether conditions of the Quality Gate can be managed.
	ManageConditions bool `json:"manageConditions"`
	// Rename defines whether the Quality Gate can be renamed.
	Rename bool `json:"rename"`
	// SetAsDefault defines whether the Quality Gate can be set as the default one.
	SetAsDefault bool `json:"setAsDefault"`
}

// +kubebuilder:object:root=true

// A QualityGate is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,sonarqube}
type QualityGate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QualityGateSpec   `json:"spec"`
	Status QualityGateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// QualityGateList contains a list of QualityGate
type QualityGateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QualityGate `json:"items"`
}

// QualityGate type metadata.
var (
	QualityGateKind             = reflect.TypeOf(QualityGate{}).Name()
	QualityGateGroupKind        = schema.GroupKind{Group: Group, Kind: QualityGateKind}.String()
	QualityGateKindAPIVersion   = QualityGateKind + "." + SchemeGroupVersion.String()
	QualityGateGroupVersionKind = SchemeGroupVersion.WithKind(QualityGateKind)
)

func init() {
	SchemeBuilder.Register(&QualityGate{}, &QualityGateList{})
}
