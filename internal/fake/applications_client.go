/*
Copyright 2026 The Crossplane Authors.

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

package fake

import (
	"context"
	"net/http"

	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"

	"github.com/crossplane/provider-sonarqube/internal/clients/instance"
)

// MockApplicationsClient is a mock implementation of the
// ApplicationsClient interface.
type MockApplicationsClient struct {
	AddProjectFn        func(opt *sonar.ApplicationsAddProjectOptions) (*http.Response, error)
	CreateFn            func(opt *sonar.ApplicationsCreateOptions) (*sonar.ApplicationsCreate, *http.Response, error)
	CreateBranchFn      func(opt *sonar.ApplicationsCreateBranchOptions) (*http.Response, error)
	DeleteFn            func(opt *sonar.ApplicationsDeleteOptions) (*http.Response, error)
	DeleteBranchFn      func(opt *sonar.ApplicationsDeleteBranchOptions) (*http.Response, error)
	RefreshFn           func(opt *sonar.ApplicationsRefreshOptions) (*http.Response, error)
	RemoveProjectFn     func(opt *sonar.ApplicationsRemoveProjectOptions) (*http.Response, error)
	SearchAllProjectsFn func(opt *sonar.ApplicationsSearchProjectsOptions) ([]sonar.ApplicationProject, *http.Response, error)
	SearchProjectsFn    func(opt *sonar.ApplicationsSearchProjectsOptions) (*sonar.ApplicationsSearchProjects, *http.Response, error)
	SetTagsFn           func(opt *sonar.ApplicationsSetTagsOptions) (*http.Response, error)
	ShowFn              func(opt *sonar.ApplicationsShowOptions) (*sonar.ApplicationsShow, *http.Response, error)
	ShowLeakFn          func(opt *sonar.ApplicationsShowLeakOptions) (*sonar.ApplicationsShowLeak, *http.Response, error)
	UpdateFn            func(opt *sonar.ApplicationsUpdateOptions) (*http.Response, error)
	UpdateBranchFn      func(opt *sonar.ApplicationsUpdateBranchOptions) (*http.Response, error)
}

// Ensure MockApplicationsClient implements ApplicationsClient.
var _ instance.ApplicationsClient = &MockApplicationsClient{}

// AddProject implements ApplicationsClient.AddProject.
func (m *MockApplicationsClient) AddProject(_ context.Context, opt *sonar.ApplicationsAddProjectOptions) (*http.Response, error) {
	if m.AddProjectFn != nil {
		return m.AddProjectFn(opt)
	}

	return nil, errNotImplemented
}

// Create implements ApplicationsClient.Create.
func (m *MockApplicationsClient) Create(_ context.Context, opt *sonar.ApplicationsCreateOptions) (*sonar.ApplicationsCreate, *http.Response, error) {
	if m.CreateFn != nil {
		return m.CreateFn(opt)
	}

	return nil, nil, errNotImplemented
}

// CreateBranch implements ApplicationsClient.CreateBranch.
func (m *MockApplicationsClient) CreateBranch(_ context.Context, opt *sonar.ApplicationsCreateBranchOptions) (*http.Response, error) {
	if m.CreateBranchFn != nil {
		return m.CreateBranchFn(opt)
	}

	return nil, errNotImplemented
}

// Delete implements ApplicationsClient.Delete.
func (m *MockApplicationsClient) Delete(_ context.Context, opt *sonar.ApplicationsDeleteOptions) (*http.Response, error) {
	if m.DeleteFn != nil {
		return m.DeleteFn(opt)
	}

	return nil, errNotImplemented
}

// DeleteBranch implements ApplicationsClient.DeleteBranch.
func (m *MockApplicationsClient) DeleteBranch(_ context.Context, opt *sonar.ApplicationsDeleteBranchOptions) (*http.Response, error) {
	if m.DeleteBranchFn != nil {
		return m.DeleteBranchFn(opt)
	}

	return nil, errNotImplemented
}

// Refresh implements ApplicationsClient.Refresh.
func (m *MockApplicationsClient) Refresh(_ context.Context, opt *sonar.ApplicationsRefreshOptions) (*http.Response, error) {
	if m.RefreshFn != nil {
		return m.RefreshFn(opt)
	}

	return nil, errNotImplemented
}

// RemoveProject implements ApplicationsClient.RemoveProject.
func (m *MockApplicationsClient) RemoveProject(_ context.Context, opt *sonar.ApplicationsRemoveProjectOptions) (*http.Response, error) {
	if m.RemoveProjectFn != nil {
		return m.RemoveProjectFn(opt)
	}

	return nil, errNotImplemented
}

// SearchAllProjects implements ApplicationsClient.SearchAllProjects.
func (m *MockApplicationsClient) SearchAllProjects(_ context.Context, opt *sonar.ApplicationsSearchProjectsOptions) ([]sonar.ApplicationProject, *http.Response, error) {
	if m.SearchAllProjectsFn != nil {
		return m.SearchAllProjectsFn(opt)
	}

	return nil, nil, errNotImplemented
}

// SearchProjects implements ApplicationsClient.SearchProjects.
func (m *MockApplicationsClient) SearchProjects(_ context.Context, opt *sonar.ApplicationsSearchProjectsOptions) (*sonar.ApplicationsSearchProjects, *http.Response, error) {
	if m.SearchProjectsFn != nil {
		return m.SearchProjectsFn(opt)
	}

	return nil, nil, errNotImplemented
}

// SetTags implements ApplicationsClient.SetTags.
func (m *MockApplicationsClient) SetTags(_ context.Context, opt *sonar.ApplicationsSetTagsOptions) (*http.Response, error) {
	if m.SetTagsFn != nil {
		return m.SetTagsFn(opt)
	}

	return nil, errNotImplemented
}

// Show implements ApplicationsClient.Show.
func (m *MockApplicationsClient) Show(_ context.Context, opt *sonar.ApplicationsShowOptions) (*sonar.ApplicationsShow, *http.Response, error) {
	if m.ShowFn != nil {
		return m.ShowFn(opt)
	}

	return nil, nil, errNotImplemented
}

// ShowLeak implements ApplicationsClient.ShowLeak.
func (m *MockApplicationsClient) ShowLeak(_ context.Context, opt *sonar.ApplicationsShowLeakOptions) (*sonar.ApplicationsShowLeak, *http.Response, error) {
	if m.ShowLeakFn != nil {
		return m.ShowLeakFn(opt)
	}

	return nil, nil, errNotImplemented
}

// Update implements ApplicationsClient.Update.
func (m *MockApplicationsClient) Update(_ context.Context, opt *sonar.ApplicationsUpdateOptions) (*http.Response, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(opt)
	}

	return nil, errNotImplemented
}

// UpdateBranch implements ApplicationsClient.UpdateBranch.
func (m *MockApplicationsClient) UpdateBranch(_ context.Context, opt *sonar.ApplicationsUpdateBranchOptions) (*http.Response, error) {
	if m.UpdateBranchFn != nil {
		return m.UpdateBranchFn(opt)
	}

	return nil, errNotImplemented
}
