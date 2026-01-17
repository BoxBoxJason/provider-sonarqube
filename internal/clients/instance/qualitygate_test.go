package instance

import (
	"testing"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"
	"github.com/google/go-cmp/cmp"
	"k8s.io/utils/ptr"

	"github.com/crossplane/provider-sonarqube/apis/instance/v1alpha1"
)

func TestGenerateQualityGateCreateOptions(t *testing.T) {
	tests := map[string]struct {
		spec v1alpha1.QualityGateParameters
		want *sonargo.QualitygatesCreateOption
	}{
		"BasicCreateOption": {
			spec: v1alpha1.QualityGateParameters{
				Name: "my-quality-gate",
			},
			want: &sonargo.QualitygatesCreateOption{
				Name: "my-quality-gate",
			},
		},
		"CreateOptionWithDefault": {
			spec: v1alpha1.QualityGateParameters{
				Name:    "default-gate",
				Default: ptr.To(true),
			},
			want: &sonargo.QualitygatesCreateOption{
				Name: "default-gate",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GenerateQualityGateCreateOptions(tc.spec)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateQualityGateCreateOptions() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGenerateQualityGateObservation(t *testing.T) {
	tests := map[string]struct {
		observation *sonargo.QualitygatesShowObject
		want        v1alpha1.QualityGateObservation
	}{
		"BasicObservation": {
			observation: &sonargo.QualitygatesShowObject{
				Name:              "test-gate",
				CaycStatus:        "compliant",
				IsBuiltIn:         false,
				IsDefault:         true,
				IsAiCodeSupported: false,
				Conditions:        []sonargo.QualitygatesShowObject_sub2{},
				Actions: sonargo.QualitygatesShowObject_sub1{
					AssociateProjects:     true,
					Copy:                  true,
					Delete:                true,
					ManageConditions:      true,
					Rename:                true,
					SetAsDefault:          true,
					Delegate:              false,
					ManageAiCodeAssurance: false,
				},
			},
			want: v1alpha1.QualityGateObservation{
				Name:              "test-gate",
				CaycStatus:        "compliant",
				IsBuiltIn:         false,
				IsDefault:         true,
				IsAiCodeSupported: false,
				Conditions:        []v1alpha1.QualityGateConditionObservation{},
				Actions: v1alpha1.QualityGatesActions{
					AssociateProjects:     true,
					Copy:                  true,
					Delete:                true,
					ManageConditions:      true,
					Rename:                true,
					SetAsDefault:          true,
					Delegate:              false,
					ManageAiCodeAssurance: false,
				},
			},
		},
		"ObservationWithConditions": {
			observation: &sonargo.QualitygatesShowObject{
				Name:       "gate-with-conditions",
				CaycStatus: "non_compliant",
				IsBuiltIn:  true,
				IsDefault:  false,
				Conditions: []sonargo.QualitygatesShowObject_sub2{
					{
						ID:     "1",
						Metric: "coverage",
						Op:     "LT",
						Error:  "80",
					},
					{
						ID:     "2",
						Metric: "duplicated_lines_density",
						Op:     "GT",
						Error:  "3",
					},
				},
				Actions: sonargo.QualitygatesShowObject_sub1{},
			},
			want: v1alpha1.QualityGateObservation{
				Name:       "gate-with-conditions",
				CaycStatus: "non_compliant",
				IsBuiltIn:  true,
				IsDefault:  false,
				Conditions: []v1alpha1.QualityGateConditionObservation{
					{
						ID:     "1",
						Metric: "coverage",
						Op:     "LT",
						Error:  "80",
					},
					{
						ID:     "2",
						Metric: "duplicated_lines_density",
						Op:     "GT",
						Error:  "3",
					},
				},
				Actions: v1alpha1.QualityGatesActions{},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GenerateQualityGateObservation(tc.observation)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateQualityGateObservation() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGenerateQualityGateActionsObservation(t *testing.T) {
	tests := map[string]struct {
		actions *sonargo.QualitygatesShowObject_sub1
		want    v1alpha1.QualityGatesActions
	}{
		"AllActionsEnabled": {
			actions: &sonargo.QualitygatesShowObject_sub1{
				AssociateProjects:     true,
				Copy:                  true,
				Delegate:              true,
				Delete:                true,
				ManageAiCodeAssurance: true,
				ManageConditions:      true,
				Rename:                true,
				SetAsDefault:          true,
			},
			want: v1alpha1.QualityGatesActions{
				AssociateProjects:     true,
				Copy:                  true,
				Delegate:              true,
				Delete:                true,
				ManageAiCodeAssurance: true,
				ManageConditions:      true,
				Rename:                true,
				SetAsDefault:          true,
			},
		},
		"NoActionsEnabled": {
			actions: &sonargo.QualitygatesShowObject_sub1{},
			want:    v1alpha1.QualityGatesActions{},
		},
		"PartialActionsEnabled": {
			actions: &sonargo.QualitygatesShowObject_sub1{
				Copy:   true,
				Rename: true,
			},
			want: v1alpha1.QualityGatesActions{
				Copy:   true,
				Rename: true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GenerateQualityGateActionsObservation(tc.actions)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateQualityGateActionsObservation() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIsQualityGateUpToDate(t *testing.T) {
	tests := map[string]struct {
		spec        *v1alpha1.QualityGateParameters
		observation *v1alpha1.QualityGateObservation
		want        bool
	}{
		"NilSpecReturnsTrue": {
			spec:        nil,
			observation: &v1alpha1.QualityGateObservation{Name: "test"},
			want:        true,
		},
		"NilObservationReturnsFalse": {
			spec:        &v1alpha1.QualityGateParameters{Name: "test"},
			observation: nil,
			want:        false,
		},
		"MatchingNameReturnsTrue": {
			spec:        &v1alpha1.QualityGateParameters{Name: "test"},
			observation: &v1alpha1.QualityGateObservation{Name: "test"},
			want:        true,
		},
		"DifferentNameReturnsFalse": {
			spec:        &v1alpha1.QualityGateParameters{Name: "test"},
			observation: &v1alpha1.QualityGateObservation{Name: "different"},
			want:        false,
		},
		"MatchingDefaultReturnsTrue": {
			spec:        &v1alpha1.QualityGateParameters{Name: "test", Default: ptr.To(true)},
			observation: &v1alpha1.QualityGateObservation{Name: "test", IsDefault: true},
			want:        true,
		},
		"DifferentDefaultReturnsFalse": {
			spec:        &v1alpha1.QualityGateParameters{Name: "test", Default: ptr.To(true)},
			observation: &v1alpha1.QualityGateObservation{Name: "test", IsDefault: false},
			want:        false,
		},
		"NilDefaultWithObservedFalseReturnsTrue": {
			spec:        &v1alpha1.QualityGateParameters{Name: "test", Default: nil},
			observation: &v1alpha1.QualityGateObservation{Name: "test", IsDefault: false},
			want:        true,
		},
		"NilDefaultWithObservedTrueReturnsTrue": {
			spec:        &v1alpha1.QualityGateParameters{Name: "test", Default: nil},
			observation: &v1alpha1.QualityGateObservation{Name: "test", IsDefault: true},
			want:        true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := IsQualityGateUpToDate(tc.spec, tc.observation)
			if got != tc.want {
				t.Errorf("IsQualityGateUpToDate() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLateInitializeQualityGate(t *testing.T) {
	tests := map[string]struct {
		spec        *v1alpha1.QualityGateParameters
		observation *v1alpha1.QualityGateObservation
		wantDefault *bool
	}{
		"NilSpecDoesNothing": {
			spec:        nil,
			observation: &v1alpha1.QualityGateObservation{IsDefault: true},
			wantDefault: nil,
		},
		"NilObservationDoesNothing": {
			spec:        &v1alpha1.QualityGateParameters{Name: "test"},
			observation: nil,
			wantDefault: nil,
		},
		"NilDefaultGetsInitialized": {
			spec:        &v1alpha1.QualityGateParameters{Name: "test", Default: nil},
			observation: &v1alpha1.QualityGateObservation{IsDefault: true},
			wantDefault: ptr.To(true),
		},
		"ExistingDefaultNotOverwritten": {
			spec:        &v1alpha1.QualityGateParameters{Name: "test", Default: ptr.To(false)},
			observation: &v1alpha1.QualityGateObservation{IsDefault: true},
			wantDefault: ptr.To(false),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			LateInitializeQualityGate(tc.spec, tc.observation)
			if tc.spec == nil {
				return
			}
			if tc.wantDefault == nil && tc.spec.Default != nil {
				t.Errorf("LateInitializeQualityGate() Default = %v, want nil", *tc.spec.Default)
				return
			}
			if tc.wantDefault != nil && tc.spec.Default == nil {
				t.Errorf("LateInitializeQualityGate() Default = nil, want %v", *tc.wantDefault)
				return
			}
			if tc.wantDefault != nil && tc.spec.Default != nil && *tc.spec.Default != *tc.wantDefault {
				t.Errorf("LateInitializeQualityGate() Default = %v, want %v", *tc.spec.Default, *tc.wantDefault)
			}
		})
	}
}
