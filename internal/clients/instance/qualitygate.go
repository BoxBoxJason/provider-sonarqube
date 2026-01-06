package instance

import (
	"net/http"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"
	"github.com/crossplane/provider-sonarqube/apis/instance/v1alpha1"
	"github.com/crossplane/provider-sonarqube/internal/clients/common"
	"github.com/crossplane/provider-sonarqube/internal/helpers"
)

// QualityGatesClient is the interface for interacting with SonarQube Quality Gates API
// It handles all the operations related to Quality Gates in SonarQube, such as creating, updating, deleting, and retrieving Quality Gates and their conditions.
// It also handles users / groups / projects association with Quality Gates.
// It also interacts with Quality Gate Conditions.
type QualityGatesClient interface {
	AddGroup(opt *sonargo.QualitygatesAddGroupOption) (resp *http.Response, err error)
	AddUser(opt *sonargo.QualitygatesAddUserOption) (resp *http.Response, err error)
	Copy(opt *sonargo.QualitygatesCopyOption) (resp *http.Response, err error)
	Create(opt *sonargo.QualitygatesCreateOption) (v *sonargo.QualitygatesCreateObject, resp *http.Response, err error)
	CreateCondition(opt *sonargo.QualitygatesCreateConditionOption) (v *sonargo.QualitygatesCreateConditionObject, resp *http.Response, err error)
	DeleteCondition(opt *sonargo.QualitygatesDeleteConditionOption) (resp *http.Response, err error)
	Deselect(opt *sonargo.QualitygatesDeselectOption) (resp *http.Response, err error)
	Destroy(opt *sonargo.QualitygatesDestroyOption) (resp *http.Response, err error)
	GetByProject(opt *sonargo.QualitygatesGetByProjectOption) (v *sonargo.QualitygatesGetByProjectObject, resp *http.Response, err error)
	List() (v *sonargo.QualitygatesListObject, resp *http.Response, err error)
	ProjectStatus(opt *sonargo.QualitygatesProjectStatusOption) (v *sonargo.QualitygatesProjectStatusObject, resp *http.Response, err error)
	RemoveGroup(opt *sonargo.QualitygatesRemoveGroupOption) (resp *http.Response, err error)
	RemoveUser(opt *sonargo.QualitygatesRemoveUserOption) (resp *http.Response, err error)
	Rename(opt *sonargo.QualitygatesRenameOption) (resp *http.Response, err error)
	Search(opt *sonargo.QualitygatesSearchOption) (v *sonargo.QualitygatesSearchObject, resp *http.Response, err error)
	SearchGroups(opt *sonargo.QualitygatesSearchGroupsOption) (v *sonargo.QualitygatesSearchGroupsObject, resp *http.Response, err error)
	SearchUsers(opt *sonargo.QualitygatesSearchUsersOption) (v *sonargo.QualitygatesSearchUsersObject, resp *http.Response, err error)
	Select(opt *sonargo.QualitygatesSelectOption) (resp *http.Response, err error)
	SetAsDefault(opt *sonargo.QualitygatesSetAsDefaultOption) (resp *http.Response, err error)
	Show(opt *sonargo.QualitygatesShowOption) (v *sonargo.QualitygatesShowObject, resp *http.Response, err error)
	UpdateCondition(opt *sonargo.QualitygatesUpdateConditionOption) (resp *http.Response, err error)
}

// NewQualityGatesClient creates a new QualityGatesClient with the provided SonarQube client configuration.
func NewQualityGatesClient(clientConfig common.Config) QualityGatesClient {
	newClient := common.NewClient(clientConfig)
	return newClient.Qualitygates
}

// GenerateQualityGateCreateOptions generates SonarQube QualitygatesCreateOption from QualityGateParameters
func GenerateQualityGateCreateOptions(spec v1alpha1.QualityGateParameters) *sonargo.QualitygatesCreateOption {
	return &sonargo.QualitygatesCreateOption{
		Name: spec.Name,
	}
}

// GenerateQualityGateObservation generates QualityGateObservation from SonarQube QualitygatesShowObject
// observation should not be nil, else it will panic
func GenerateQualityGateObservation(observation *sonargo.QualitygatesShowObject) v1alpha1.QualityGateObservation {
	return v1alpha1.QualityGateObservation{
		Actions:           GenerateQualityGateActionsObservation(&observation.Actions),
		CaycStatus:        observation.CaycStatus,
		Conditions:        GenerateQualityGateConditionsObservation(observation.Conditions),
		IsAiCodeSupported: observation.IsAiCodeSupported,
		IsBuiltIn:         observation.IsBuiltIn,
		IsDefault:         observation.IsDefault,
		Name:              observation.Name,
	}
}

// GenerateQualityGateActionsObservation generates QualityGatesActions from SonarQube QualitygatesShowObject_sub1
// actions should not be nil, else it will panic
func GenerateQualityGateActionsObservation(actions *sonargo.QualitygatesShowObject_sub1) v1alpha1.QualityGatesActions {
	return v1alpha1.QualityGatesActions{
		AssociateProjects:     actions.AssociateProjects,
		Copy:                  actions.Copy,
		Delegate:              actions.Delegate,
		Delete:                actions.Delete,
		ManageAiCodeAssurance: actions.ManageAiCodeAssurance,
		ManageConditions:      actions.ManageConditions,
		Rename:                actions.Rename,
		SetAsDefault:          actions.SetAsDefault,
	}
}

// IsQualityGateUpToDate checks if the Quality Gate spec is up to date with the observed state
func IsQualityGateUpToDate(spec *v1alpha1.QualityGateParameters, observation *v1alpha1.QualityGateObservation) bool {
	if spec == nil {
		return true
	}
	if observation == nil {
		return false
	}

	if spec.Name != observation.Name {
		return false
	}

	if !helpers.IsComparablePtrEqualComparable(spec.Default, observation.IsDefault) {
		return false
	}

	return true
}

// LateInitializeQualityGate fills the spec with the observed state if the spec fields are nil
func LateInitializeQualityGate(spec *v1alpha1.QualityGateParameters, observation *v1alpha1.QualityGateObservation) {
	if spec == nil || observation == nil {
		return
	}

	helpers.AssignIfNil(&spec.Default, observation.IsDefault)
}
