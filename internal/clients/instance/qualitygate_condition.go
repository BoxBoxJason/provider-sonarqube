package instance

import (
	"errors"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"
	"github.com/crossplane/provider-sonarqube/apis/instance/v1alpha1"
	"github.com/crossplane/provider-sonarqube/internal/helpers"
)

// GenerateQualityGateConditionObservation generates QualityGateConditionObservation from SonarQube QualitygatesShowObject_sub2
func GenerateQualityGateConditionObservation(condition sonargo.QualitygatesShowObject_sub2) v1alpha1.QualityGateConditionObservation {
	return v1alpha1.QualityGateConditionObservation{
		Error:  condition.Error,
		ID:     condition.ID,
		Metric: condition.Metric,
		Op:     condition.Op,
	}
}

// GenerateQualityGateConditionsObservation generates a slice of QualityGateConditionObservation from a slice of SonarQube QualitygatesShowObject_sub2
func GenerateQualityGateConditionsObservation(conditions []sonargo.QualitygatesShowObject_sub2) []v1alpha1.QualityGateConditionObservation {
	conditionObservations := make([]v1alpha1.QualityGateConditionObservation, len(conditions))
	for i, condition := range conditions {
		conditionObservations[i] = GenerateQualityGateConditionObservation(condition)
	}
	return conditionObservations
}

// FindQualityGateConditionObservation finds a QualityGateConditionObservation by ID from a slice of SonarQube QualitygatesShowObject_sub2
func FindQualityGateConditionObservation(id string, condition []sonargo.QualitygatesShowObject_sub2) (v1alpha1.QualityGateConditionObservation, error) {
	for _, cond := range condition {
		if cond.ID == id {
			return GenerateQualityGateConditionObservation(cond), nil
		}
	}
	return v1alpha1.QualityGateConditionObservation{}, errors.New("quality gate condition not found in observation")
}

// GenerateCreateQualityGateConditionOption generates SonarQube QualitygatesCreateConditionOption from QualityGateConditionParameters
func GenerateCreateQualityGateConditionOption(params v1alpha1.QualityGateConditionParameters) sonargo.QualitygatesCreateConditionOption {
	option := sonargo.QualitygatesCreateConditionOption{
		GateName: *params.QualityGateName,
		Error:    params.Error,
		Metric:   params.Metric,
	}
	if params.Op != nil {
		option.Op = *params.Op
	}
	return option
}

// GenerateUpdateQualityGateConditionOption generates SonarQube QualitygatesUpdateConditionOption from QualityGateConditionParameters
func GenerateUpdateQualityGateConditionOption(id string, params v1alpha1.QualityGateConditionParameters) sonargo.QualitygatesUpdateConditionOption {
	option := sonargo.QualitygatesUpdateConditionOption{
		Id:     id,
		Error:  params.Error,
		Metric: params.Metric,
	}
	if params.Op != nil {
		option.Op = *params.Op
	}
	return option
}

// GenerateDeleteQualityGateConditionOption generates SonarQube QualitygatesDeleteConditionOption from QualityGateConditionObservation
func GenerateDeleteQualityGateConditionOption(id string) *sonargo.QualitygatesDeleteConditionOption {
	return &sonargo.QualitygatesDeleteConditionOption{
		Id: id,
	}
}

// IsQualityGateConditionUpToDate checks whether the observed QualityGateCondition is up to date with the desired QualityGateConditionParameters
func IsQualityGateConditionUpToDate(params *v1alpha1.QualityGateConditionParameters, observation *v1alpha1.QualityGateConditionObservation) bool {
	if params == nil {
		return true
	}
	if observation == nil {
		return false
	}

	if params.Error != observation.Error {
		return false
	}
	if params.Metric != observation.Metric {
		return false
	}
	if !helpers.IsComparablePtrEqualComparable(params.Op, observation.Op) {
		return false
	}

	return true
}

// LateInitializeQualityGateCondition fills the empty fields in *QualityGateConditionParameters with
// the values seen in QualityGateConditionObservation.
func LateInitializeQualityGateCondition(params *v1alpha1.QualityGateConditionParameters, observation *v1alpha1.QualityGateConditionObservation) {
	if params == nil || observation == nil {
		return
	}
	helpers.AssignIfNil(&params.Op, observation.Op)
}
