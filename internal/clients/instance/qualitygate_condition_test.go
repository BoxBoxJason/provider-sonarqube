package instance

import (
	"testing"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"
	"github.com/google/go-cmp/cmp"
	"k8s.io/utils/ptr"

	"github.com/crossplane/provider-sonarqube/apis/instance/v1alpha1"
)

func TestGenerateQualityGateConditionObservation(t *testing.T) {
	tests := map[string]struct {
		condition sonargo.QualitygatesShowObject_sub2
		want      v1alpha1.QualityGateConditionObservation
	}{
		"BasicCondition": {
			condition: sonargo.QualitygatesShowObject_sub2{
				ID:     "123",
				Metric: "coverage",
				Op:     "LT",
				Error:  "80",
			},
			want: v1alpha1.QualityGateConditionObservation{
				ID:     "123",
				Metric: "coverage",
				Op:     "LT",
				Error:  "80",
			},
		},
		"EmptyCondition": {
			condition: sonargo.QualitygatesShowObject_sub2{},
			want:      v1alpha1.QualityGateConditionObservation{},
		},
		"ConditionWithGTOperator": {
			condition: sonargo.QualitygatesShowObject_sub2{
				ID:     "456",
				Metric: "duplicated_lines_density",
				Op:     "GT",
				Error:  "3",
			},
			want: v1alpha1.QualityGateConditionObservation{
				ID:     "456",
				Metric: "duplicated_lines_density",
				Op:     "GT",
				Error:  "3",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GenerateQualityGateConditionObservation(tc.condition)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateQualityGateConditionObservation() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGenerateQualityGateConditionsObservation(t *testing.T) {
	tests := map[string]struct {
		conditions []sonargo.QualitygatesShowObject_sub2
		want       []v1alpha1.QualityGateConditionObservation
	}{
		"EmptySlice": {
			conditions: []sonargo.QualitygatesShowObject_sub2{},
			want:       []v1alpha1.QualityGateConditionObservation{},
		},
		"SingleCondition": {
			conditions: []sonargo.QualitygatesShowObject_sub2{
				{ID: "1", Metric: "coverage", Op: "LT", Error: "80"},
			},
			want: []v1alpha1.QualityGateConditionObservation{
				{ID: "1", Metric: "coverage", Op: "LT", Error: "80"},
			},
		},
		"MultipleConditions": {
			conditions: []sonargo.QualitygatesShowObject_sub2{
				{ID: "1", Metric: "coverage", Op: "LT", Error: "80"},
				{ID: "2", Metric: "duplicated_lines_density", Op: "GT", Error: "3"},
				{ID: "3", Metric: "new_coverage", Op: "LT", Error: "90"},
			},
			want: []v1alpha1.QualityGateConditionObservation{
				{ID: "1", Metric: "coverage", Op: "LT", Error: "80"},
				{ID: "2", Metric: "duplicated_lines_density", Op: "GT", Error: "3"},
				{ID: "3", Metric: "new_coverage", Op: "LT", Error: "90"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GenerateQualityGateConditionsObservation(tc.conditions)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateQualityGateConditionsObservation() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFindQualityGateConditionObservation(t *testing.T) {
	conditions := []sonargo.QualitygatesShowObject_sub2{
		{ID: "1", Metric: "coverage", Op: "LT", Error: "80"},
		{ID: "2", Metric: "duplicated_lines_density", Op: "GT", Error: "3"},
	}

	tests := map[string]struct {
		id        string
		condition []sonargo.QualitygatesShowObject_sub2
		want      v1alpha1.QualityGateConditionObservation
		wantErr   bool
	}{
		"FoundFirstCondition": {
			id:        "1",
			condition: conditions,
			want:      v1alpha1.QualityGateConditionObservation{ID: "1", Metric: "coverage", Op: "LT", Error: "80"},
			wantErr:   false,
		},
		"FoundSecondCondition": {
			id:        "2",
			condition: conditions,
			want:      v1alpha1.QualityGateConditionObservation{ID: "2", Metric: "duplicated_lines_density", Op: "GT", Error: "3"},
			wantErr:   false,
		},
		"NotFoundReturnsError": {
			id:        "999",
			condition: conditions,
			want:      v1alpha1.QualityGateConditionObservation{},
			wantErr:   true,
		},
		"EmptySliceReturnsError": {
			id:        "1",
			condition: []sonargo.QualitygatesShowObject_sub2{},
			want:      v1alpha1.QualityGateConditionObservation{},
			wantErr:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := FindQualityGateConditionObservation(tc.id, tc.condition)
			if (err != nil) != tc.wantErr {
				t.Errorf("FindQualityGateConditionObservation() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr {
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("FindQualityGateConditionObservation() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestGenerateCreateQualityGateConditionOption(t *testing.T) {
	tests := map[string]struct {
		params v1alpha1.QualityGateConditionParameters
		want   sonargo.QualitygatesCreateConditionOption
	}{
		"BasicCondition": {
			params: v1alpha1.QualityGateConditionParameters{
				QualityGateName: ptr.To("my-gate"),
				Metric:          "coverage",
				Error:           "80",
			},
			want: sonargo.QualitygatesCreateConditionOption{
				GateName: "my-gate",
				Metric:   "coverage",
				Error:    "80",
			},
		},
		"ConditionWithOperator": {
			params: v1alpha1.QualityGateConditionParameters{
				QualityGateName: ptr.To("my-gate"),
				Metric:          "coverage",
				Error:           "80",
				Op:              ptr.To("LT"),
			},
			want: sonargo.QualitygatesCreateConditionOption{
				GateName: "my-gate",
				Metric:   "coverage",
				Error:    "80",
				Op:       "LT",
			},
		},
		"ConditionWithGTOperator": {
			params: v1alpha1.QualityGateConditionParameters{
				QualityGateName: ptr.To("another-gate"),
				Metric:          "duplicated_lines_density",
				Error:           "3",
				Op:              ptr.To("GT"),
			},
			want: sonargo.QualitygatesCreateConditionOption{
				GateName: "another-gate",
				Metric:   "duplicated_lines_density",
				Error:    "3",
				Op:       "GT",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GenerateCreateQualityGateConditionOption(tc.params)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateCreateQualityGateConditionOption() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGenerateUpdateQualityGateConditionOption(t *testing.T) {
	tests := map[string]struct {
		id     string
		params v1alpha1.QualityGateConditionParameters
		want   sonargo.QualitygatesUpdateConditionOption
	}{
		"BasicUpdate": {
			id: "123",
			params: v1alpha1.QualityGateConditionParameters{
				QualityGateName: ptr.To("my-gate"),
				Metric:          "coverage",
				Error:           "85",
			},
			want: sonargo.QualitygatesUpdateConditionOption{
				Id:     "123",
				Metric: "coverage",
				Error:  "85",
			},
		},
		"UpdateWithOperator": {
			id: "456",
			params: v1alpha1.QualityGateConditionParameters{
				QualityGateName: ptr.To("my-gate"),
				Metric:          "duplicated_lines_density",
				Error:           "5",
				Op:              ptr.To("GT"),
			},
			want: sonargo.QualitygatesUpdateConditionOption{
				Id:     "456",
				Metric: "duplicated_lines_density",
				Error:  "5",
				Op:     "GT",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GenerateUpdateQualityGateConditionOption(tc.id, tc.params)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateUpdateQualityGateConditionOption() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGenerateDeleteQualityGateConditionOption(t *testing.T) {
	tests := map[string]struct {
		id   string
		want *sonargo.QualitygatesDeleteConditionOption
	}{
		"BasicDelete": {
			id:   "123",
			want: &sonargo.QualitygatesDeleteConditionOption{Id: "123"},
		},
		"EmptyID": {
			id:   "",
			want: &sonargo.QualitygatesDeleteConditionOption{Id: ""},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GenerateDeleteQualityGateConditionOption(tc.id)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateDeleteQualityGateConditionOption() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIsQualityGateConditionUpToDate(t *testing.T) {
	tests := map[string]struct {
		params      *v1alpha1.QualityGateConditionParameters
		observation *v1alpha1.QualityGateConditionObservation
		want        bool
	}{
		"NilParamsReturnsTrue": {
			params:      nil,
			observation: &v1alpha1.QualityGateConditionObservation{},
			want:        true,
		},
		"NilObservationReturnsFalse": {
			params:      &v1alpha1.QualityGateConditionParameters{},
			observation: nil,
			want:        false,
		},
		"MatchingValuesReturnsTrue": {
			params: &v1alpha1.QualityGateConditionParameters{
				Metric: "coverage",
				Error:  "80",
				Op:     ptr.To("LT"),
			},
			observation: &v1alpha1.QualityGateConditionObservation{
				Metric: "coverage",
				Error:  "80",
				Op:     "LT",
			},
			want: true,
		},
		"DifferentErrorReturnsFalse": {
			params: &v1alpha1.QualityGateConditionParameters{
				Metric: "coverage",
				Error:  "80",
				Op:     ptr.To("LT"),
			},
			observation: &v1alpha1.QualityGateConditionObservation{
				Metric: "coverage",
				Error:  "85",
				Op:     "LT",
			},
			want: false,
		},
		"DifferentMetricReturnsFalse": {
			params: &v1alpha1.QualityGateConditionParameters{
				Metric: "coverage",
				Error:  "80",
			},
			observation: &v1alpha1.QualityGateConditionObservation{
				Metric: "new_coverage",
				Error:  "80",
			},
			want: false,
		},
		"DifferentOpReturnsFalse": {
			params: &v1alpha1.QualityGateConditionParameters{
				Metric: "coverage",
				Error:  "80",
				Op:     ptr.To("LT"),
			},
			observation: &v1alpha1.QualityGateConditionObservation{
				Metric: "coverage",
				Error:  "80",
				Op:     "GT",
			},
			want: false,
		},
		"NilOpMatchesAnyObservedOp": {
			params: &v1alpha1.QualityGateConditionParameters{
				Metric: "coverage",
				Error:  "80",
				Op:     nil,
			},
			observation: &v1alpha1.QualityGateConditionObservation{
				Metric: "coverage",
				Error:  "80",
				Op:     "LT",
			},
			want: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := IsQualityGateConditionUpToDate(tc.params, tc.observation)
			if got != tc.want {
				t.Errorf("IsQualityGateConditionUpToDate() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLateInitializeQualityGateCondition(t *testing.T) {
	tests := map[string]struct {
		params      *v1alpha1.QualityGateConditionParameters
		observation *v1alpha1.QualityGateConditionObservation
		wantOp      *string
	}{
		"NilParamsDoesNothing": {
			params:      nil,
			observation: &v1alpha1.QualityGateConditionObservation{Op: "LT"},
			wantOp:      nil,
		},
		"NilObservationDoesNothing": {
			params:      &v1alpha1.QualityGateConditionParameters{Metric: "coverage"},
			observation: nil,
			wantOp:      nil,
		},
		"NilOpGetsInitialized": {
			params:      &v1alpha1.QualityGateConditionParameters{Metric: "coverage", Op: nil},
			observation: &v1alpha1.QualityGateConditionObservation{Op: "LT"},
			wantOp:      ptr.To("LT"),
		},
		"ExistingOpNotOverwritten": {
			params:      &v1alpha1.QualityGateConditionParameters{Metric: "coverage", Op: ptr.To("GT")},
			observation: &v1alpha1.QualityGateConditionObservation{Op: "LT"},
			wantOp:      ptr.To("GT"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			LateInitializeQualityGateCondition(tc.params, tc.observation)
			if tc.params == nil {
				return
			}
			if tc.wantOp == nil && tc.params.Op != nil {
				t.Errorf("LateInitializeQualityGateCondition() Op = %v, want nil", *tc.params.Op)
				return
			}
			if tc.wantOp != nil && tc.params.Op == nil {
				t.Errorf("LateInitializeQualityGateCondition() Op = nil, want %v", *tc.wantOp)
				return
			}
			if tc.wantOp != nil && tc.params.Op != nil && *tc.params.Op != *tc.wantOp {
				t.Errorf("LateInitializeQualityGateCondition() Op = %v, want %v", *tc.params.Op, *tc.wantOp)
			}
		})
	}
}
