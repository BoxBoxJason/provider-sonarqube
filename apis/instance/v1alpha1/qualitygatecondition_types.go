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

// QualityGateConditionParameters are the configurable fields of a QualityGateCondition.
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.qualityGateName) || has(self.qualityGateName)", message="QualityGateName is required once set"
type QualityGateConditionParameters struct {
	// Name of the quality gate to which the condition belongs.
	// WARNING: QualityGateName is immutable once set.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="QualityGateName is immutable once set."
	// +kubebuilder:validation:MaxLength=100
	// +kubebuilder:validation:MinLength=1
	QualityGateName *string `json:"qualityGateName,omitempty"`

	// Reference to a QualityGate to which the condition belongs.
	// WARNING: QualityGateRef is immutable once set.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="QualityGateRef is immutable once set."
	QualityGateRef *xpv1.NamespacedReference `json:"qualityGateRef,omitempty"`

	// Selector for a QualityGate to which the condition belongs.
	// WARNING: QualityGateSelector is immutable once set.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="QualityGateSelector is immutable once set."
	QualityGateSelector *xpv1.NamespacedSelector `json:"qualityGateSelector,omitempty"`

	// Error is the Condition error threshold
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:MinLength=1
	Error string `json:"error,omitempty"`

	// Metric is the Condition metric that the condition applies to.
	// Only accepts metrics of the following types: INT, MILLISEC, RATING, WORK_DUR, FLOAT, PERCENT, LEVEL.
	// The following metrics are forbidden: alert_status, security_hotspots, new_security_hotspots.
	// +kubebuilder:validation:Pattern="^[a-zA-Z0-9_]+$"
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Metric string `json:"metric,omitempty"`

	// Op is the Condition operator.
	// Only LT (is lower than) and GT (is greater than) are supported.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=LT;GT
	Op *string `json:"op,omitempty"`
}

// QualityGateConditionObservation are the observable fields of a QualityGateCondition.
type QualityGateConditionObservation struct {
	// Error is the Condition error threshold
	Error string `json:"error,omitempty"`
	// ID is the Condition ID
	ID string `json:"id,omitempty"`
	// Metric is the Condition metric that the condition applies to.
	Metric string `json:"metric,omitempty"`
	// Op is the Condition operator.
	Op string `json:"op,omitempty"`
}

// A QualityGateConditionSpec defines the desired state of a QualityGateCondition.
type QualityGateConditionSpec struct {
	xpv2.ManagedResourceSpec `json:",inline"`
	ForProvider              QualityGateConditionParameters `json:"forProvider"`
}

// A QualityGateConditionStatus represents the observed state of a QualityGateCondition.
type QualityGateConditionStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          QualityGateConditionObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A QualityGateCondition is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={crossplane,managed,sonarqube}
type QualityGateCondition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QualityGateConditionSpec   `json:"spec"`
	Status QualityGateConditionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// QualityGateConditionList contains a list of QualityGateCondition
type QualityGateConditionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QualityGateCondition `json:"items"`
}

// QualityGateCondition type metadata.
var (
	QualityGateConditionKind             = reflect.TypeOf(QualityGateCondition{}).Name()
	QualityGateConditionGroupKind        = schema.GroupKind{Group: Group, Kind: QualityGateConditionKind}.String()
	QualityGateConditionKindAPIVersion   = QualityGateConditionKind + "." + SchemeGroupVersion.String()
	QualityGateConditionGroupVersionKind = SchemeGroupVersion.WithKind(QualityGateConditionKind)
)

func init() {
	SchemeBuilder.Register(&QualityGateCondition{}, &QualityGateConditionList{})
}
