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

package qualitygate

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	v1alpha1 "github.com/crossplane/provider-sonarqube/apis/instance/v1alpha1"
	"github.com/crossplane/provider-sonarqube/internal/fake"
)

// Unlike many Kubernetes projects Crossplane does not use third party testing
// libraries, per the common Go test review comments. Crossplane encourages the
// use of table driven unit tests. The tests of the crossplane-runtime project
// are representative of the testing style Crossplane encourages.
//
// https://github.com/golang/go/wiki/TestComments
// https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md#contributing-code

type notQualityGate struct {
	resource.Managed
}

func TestObserve(t *testing.T) {
	type args struct {
		ctx context.Context
		mg  resource.Managed
	}
	type want struct {
		o   managed.ExternalObservation
		err error
	}

	cases := map[string]struct {
		client *fake.MockQualityGatesClient
		args   args
		want   want
	}{
		"NotQualityGateError": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg:  &notQualityGate{},
			},
			want: want{
				o:   managed.ExternalObservation{},
				err: errors.New(errNotQualityGate),
			},
		},
		"EmptyExternalNameReturnsNotExists": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGate{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "test-gate",
						Annotations: map[string]string{},
					},
				},
			},
			want: want{
				o:   managed.ExternalObservation{ResourceExists: false},
				err: nil,
			},
		},
		"ShowFailsReturnsNotExists": {
			client: &fake.MockQualityGatesClient{
				ShowFn: func(opt *sonargo.QualitygatesShowOption) (*sonargo.QualitygatesShowObject, *http.Response, error) {
					return nil, nil, errors.New("api error")
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGate {
					qg := &v1alpha1.QualityGate{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-gate",
							Annotations: map[string]string{},
						},
					}
					meta.SetExternalName(qg, "test-gate")
					return qg
				}(),
			},
			want: want{
				o:   managed.ExternalObservation{ResourceExists: false},
				err: nil,
			},
		},
		"SuccessfulObserveResourceExists": {
			client: &fake.MockQualityGatesClient{
				ShowFn: func(opt *sonargo.QualitygatesShowOption) (*sonargo.QualitygatesShowObject, *http.Response, error) {
					return &sonargo.QualitygatesShowObject{
						Name:       "test-gate",
						CaycStatus: "compliant",
						IsBuiltIn:  false,
						IsDefault:  false,
						Conditions: []sonargo.QualitygatesShowObject_sub2{},
						Actions:    sonargo.QualitygatesShowObject_sub1{},
					}, nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGate {
					qg := &v1alpha1.QualityGate{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-gate",
							Annotations: map[string]string{},
						},
						Spec: v1alpha1.QualityGateSpec{
							ForProvider: v1alpha1.QualityGateParameters{
								Name:    "test-gate",
								Default: ptr.To(false), // explicitly set to match observation and avoid late initialization
							},
						},
					}
					meta.SetExternalName(qg, "test-gate")
					return qg
				}(),
			},
			want: want{
				o: managed.ExternalObservation{
					ResourceExists:          true,
					ResourceUpToDate:        true,
					ResourceLateInitialized: false,
				},
				err: nil,
			},
		},
		"ResourceNotUpToDateWhenNamesDiffer": {
			client: &fake.MockQualityGatesClient{
				ShowFn: func(opt *sonargo.QualitygatesShowOption) (*sonargo.QualitygatesShowObject, *http.Response, error) {
					return &sonargo.QualitygatesShowObject{
						Name:       "different-name",
						CaycStatus: "compliant",
						IsBuiltIn:  false,
						IsDefault:  false,
						Conditions: []sonargo.QualitygatesShowObject_sub2{},
						Actions:    sonargo.QualitygatesShowObject_sub1{},
					}, nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGate {
					qg := &v1alpha1.QualityGate{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-gate",
							Annotations: map[string]string{},
						},
						Spec: v1alpha1.QualityGateSpec{
							ForProvider: v1alpha1.QualityGateParameters{
								Name:    "test-gate",
								Default: ptr.To(false), // explicitly set to match observation and avoid late initialization
							},
						},
					}
					meta.SetExternalName(qg, "test-gate")
					return qg
				}(),
			},
			want: want{
				o: managed.ExternalObservation{
					ResourceExists:          true,
					ResourceUpToDate:        false,
					ResourceLateInitialized: false,
				},
				err: nil,
			},
		},
		"LateInitializeDefault": {
			client: &fake.MockQualityGatesClient{
				ShowFn: func(opt *sonargo.QualitygatesShowOption) (*sonargo.QualitygatesShowObject, *http.Response, error) {
					return &sonargo.QualitygatesShowObject{
						Name:       "test-gate",
						CaycStatus: "compliant",
						IsBuiltIn:  false,
						IsDefault:  true,
						Conditions: []sonargo.QualitygatesShowObject_sub2{},
						Actions:    sonargo.QualitygatesShowObject_sub1{},
					}, nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGate {
					qg := &v1alpha1.QualityGate{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-gate",
							Annotations: map[string]string{},
						},
						Spec: v1alpha1.QualityGateSpec{
							ForProvider: v1alpha1.QualityGateParameters{
								Name:    "test-gate",
								Default: nil,
							},
						},
					}
					meta.SetExternalName(qg, "test-gate")
					return qg
				}(),
			},
			want: want{
				o: managed.ExternalObservation{
					ResourceExists:          true,
					ResourceUpToDate:        true,
					ResourceLateInitialized: true,
				},
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{qualityGatesClient: tc.client}
			got, err := e.Observe(tc.args.ctx, tc.args.mg)

			if diff := cmp.Diff(tc.want.err, err, cmp.Comparer(errComparer)); diff != "" {
				t.Errorf("Observe() error mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.o, got); diff != "" {
				t.Errorf("Observe() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		ctx context.Context
		mg  resource.Managed
	}
	type want struct {
		o   managed.ExternalCreation
		err error
	}

	cases := map[string]struct {
		client *fake.MockQualityGatesClient
		args   args
		want   want
	}{
		"NotQualityGateError": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg:  &notQualityGate{},
			},
			want: want{
				o:   managed.ExternalCreation{},
				err: errors.New(errNotQualityGate),
			},
		},
		"CreateFails": {
			client: &fake.MockQualityGatesClient{
				CreateFn: func(opt *sonargo.QualitygatesCreateOption) (*sonargo.QualitygatesCreateObject, *http.Response, error) {
					return nil, nil, errors.New("create error")
				},
			},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGate{
					ObjectMeta: metav1.ObjectMeta{Name: "test-gate"},
					Spec: v1alpha1.QualityGateSpec{
						ForProvider: v1alpha1.QualityGateParameters{
							Name: "test-gate",
						},
					},
				},
			},
			want: want{
				o:   managed.ExternalCreation{},
				err: errors.Wrap(errors.New("create error"), errCreateQualityGate),
			},
		},
		"SuccessfulCreate": {
			client: &fake.MockQualityGatesClient{
				CreateFn: func(opt *sonargo.QualitygatesCreateOption) (*sonargo.QualitygatesCreateObject, *http.Response, error) {
					return &sonargo.QualitygatesCreateObject{
						ID:   "gate-123",
						Name: opt.Name,
					}, nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGate{
					ObjectMeta: metav1.ObjectMeta{Name: "test-gate"},
					Spec: v1alpha1.QualityGateSpec{
						ForProvider: v1alpha1.QualityGateParameters{
							Name: "test-gate",
						},
					},
				},
			},
			want: want{
				o:   managed.ExternalCreation{},
				err: nil,
			},
		},
		"ExternalNameSetToSonarQubeName": {
			client: &fake.MockQualityGatesClient{
				CreateFn: func(opt *sonargo.QualitygatesCreateOption) (*sonargo.QualitygatesCreateObject, *http.Response, error) {
					return &sonargo.QualitygatesCreateObject{
						ID:   "some-generated-id",
						Name: "MySonarQubeGateName",
					}, nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGate{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-resource-name"},
					Spec: v1alpha1.QualityGateSpec{
						ForProvider: v1alpha1.QualityGateParameters{
							Name: "MySonarQubeGateName",
						},
					},
				},
			},
			want: want{
				o:   managed.ExternalCreation{},
				err: nil,
			},
		},
		"CreateWithDefaultTrue": {
			client: &fake.MockQualityGatesClient{
				CreateFn: func(opt *sonargo.QualitygatesCreateOption) (*sonargo.QualitygatesCreateObject, *http.Response, error) {
					return &sonargo.QualitygatesCreateObject{
						ID:   "gate-123",
						Name: "my-sonar-gate", // different from k8s resource name to test the fix
					}, nil, nil
				},
				SetAsDefaultFn: func(opt *sonargo.QualitygatesSetAsDefaultOption) (*http.Response, error) {
					// Verify the correct SonarQube quality gate name is used, not Kubernetes resource name
					if opt.Name != "my-sonar-gate" {
						return nil, errors.New("expected SonarQube gate name but got: " + opt.Name)
					}
					return nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGate{
					ObjectMeta: metav1.ObjectMeta{Name: "test-gate"}, // different from SonarQube name
					Spec: v1alpha1.QualityGateSpec{
						ForProvider: v1alpha1.QualityGateParameters{
							Name:    "my-sonar-gate",
							Default: ptr.To(true),
						},
					},
				},
			},
			want: want{
				o:   managed.ExternalCreation{},
				err: nil,
			},
		},
		"CreateWithDefaultTrueButSetDefaultFails": {
			client: &fake.MockQualityGatesClient{
				CreateFn: func(opt *sonargo.QualitygatesCreateOption) (*sonargo.QualitygatesCreateObject, *http.Response, error) {
					return &sonargo.QualitygatesCreateObject{
						ID:   "gate-123",
						Name: opt.Name,
					}, nil, nil
				},
				SetAsDefaultFn: func(opt *sonargo.QualitygatesSetAsDefaultOption) (*http.Response, error) {
					return nil, errors.New("set default error")
				},
			},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGate{
					ObjectMeta: metav1.ObjectMeta{Name: "test-gate"},
					Spec: v1alpha1.QualityGateSpec{
						ForProvider: v1alpha1.QualityGateParameters{
							Name:    "test-gate",
							Default: ptr.To(true),
						},
					},
				},
			},
			want: want{
				o:   managed.ExternalCreation{},
				err: errors.Wrap(errors.New("set default error"), errDefaultQualityGate),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{qualityGatesClient: tc.client}
			got, err := e.Create(tc.args.ctx, tc.args.mg)

			if diff := cmp.Diff(tc.want.err, err, cmp.Comparer(errComparer)); diff != "" {
				t.Errorf("Create() error mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.o, got); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		ctx context.Context
		mg  resource.Managed
	}
	type want struct {
		o   managed.ExternalUpdate
		err error
	}

	cases := map[string]struct {
		client *fake.MockQualityGatesClient
		args   args
		want   want
	}{
		"NotQualityGateError": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg:  &notQualityGate{},
			},
			want: want{
				o:   managed.ExternalUpdate{},
				err: errors.New(errNotQualityGate),
			},
		},
		"EmptyExternalNameReturnsError": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGate{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "test-gate",
						Annotations: map[string]string{},
					},
				},
			},
			want: want{
				o:   managed.ExternalUpdate{},
				err: fmt.Errorf("external name is not set for Quality Gate %s", "test-gate"),
			},
		},
		"SetAsDefaultWhenRequested": {
			client: &fake.MockQualityGatesClient{
				SetAsDefaultFn: func(opt *sonargo.QualitygatesSetAsDefaultOption) (*http.Response, error) {
					return nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGate {
					qg := &v1alpha1.QualityGate{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-gate",
							Annotations: map[string]string{},
						},
						Spec: v1alpha1.QualityGateSpec{
							ForProvider: v1alpha1.QualityGateParameters{
								Name:    "test-gate",
								Default: ptr.To(true),
							},
						},
					}
					meta.SetExternalName(qg, "test-gate")
					return qg
				}(),
			},
			want: want{
				o:   managed.ExternalUpdate{},
				err: nil,
			},
		},
		"SetAsDefaultFails": {
			client: &fake.MockQualityGatesClient{
				SetAsDefaultFn: func(opt *sonargo.QualitygatesSetAsDefaultOption) (*http.Response, error) {
					return nil, errors.New("set default error")
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGate {
					qg := &v1alpha1.QualityGate{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-gate",
							Annotations: map[string]string{},
						},
						Spec: v1alpha1.QualityGateSpec{
							ForProvider: v1alpha1.QualityGateParameters{
								Name:    "test-gate",
								Default: ptr.To(true),
							},
						},
					}
					meta.SetExternalName(qg, "test-gate")
					return qg
				}(),
			},
			want: want{
				o:   managed.ExternalUpdate{},
				err: errors.Wrap(errors.New("set default error"), errDefaultQualityGate),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{qualityGatesClient: tc.client}
			got, err := e.Update(tc.args.ctx, tc.args.mg)

			if diff := cmp.Diff(tc.want.err, err, cmp.Comparer(errComparer)); diff != "" {
				t.Errorf("Update() error mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.o, got); diff != "" {
				t.Errorf("Update() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		ctx context.Context
		mg  resource.Managed
	}
	type want struct {
		o   managed.ExternalDelete
		err error
	}

	cases := map[string]struct {
		client *fake.MockQualityGatesClient
		args   args
		want   want
	}{
		"NotQualityGateError": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg:  &notQualityGate{},
			},
			want: want{
				o:   managed.ExternalDelete{},
				err: errors.New(errNotQualityGate),
			},
		},
		"EmptyExternalNameDoesNothing": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGate{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "test-gate",
						Annotations: map[string]string{},
					},
				},
			},
			want: want{
				o:   managed.ExternalDelete{},
				err: nil,
			},
		},
		"SuccessfulDelete": {
			client: &fake.MockQualityGatesClient{
				DestroyFn: func(opt *sonargo.QualitygatesDestroyOption) (*http.Response, error) {
					// Verify the correct external name is used for deletion
					if opt.Name != "my-sonar-gate" {
						return nil, errors.New("expected external name 'my-sonar-gate' but got: " + opt.Name)
					}
					return nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGate {
					qg := &v1alpha1.QualityGate{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "k8s-resource-name", // different from external name to test the fix
							Annotations: map[string]string{},
						},
					}
					meta.SetExternalName(qg, "my-sonar-gate") // this should be used for deletion
					return qg
				}(),
			},
			want: want{
				o:   managed.ExternalDelete{},
				err: nil,
			},
		},
		"DeleteFails": {
			client: &fake.MockQualityGatesClient{
				DestroyFn: func(opt *sonargo.QualitygatesDestroyOption) (*http.Response, error) {
					return nil, errors.New("delete error")
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGate {
					qg := &v1alpha1.QualityGate{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-gate",
							Annotations: map[string]string{},
						},
					}
					meta.SetExternalName(qg, "test-gate")
					return qg
				}(),
			},
			want: want{
				o:   managed.ExternalDelete{},
				err: errors.Wrap(errors.New("delete error"), errDeleteQualityGate),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{qualityGatesClient: tc.client}
			got, err := e.Delete(tc.args.ctx, tc.args.mg)

			if diff := cmp.Diff(tc.want.err, err, cmp.Comparer(errComparer)); diff != "" {
				t.Errorf("Delete() error mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.o, got); diff != "" {
				t.Errorf("Delete() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDisconnect(t *testing.T) {
	e := &external{qualityGatesClient: &fake.MockQualityGatesClient{}}
	err := e.Disconnect(context.Background())
	if err != nil {
		t.Errorf("Disconnect() error = %v, want nil", err)
	}
}

func TestCreateSetsExternalNameToSonarQubeName(t *testing.T) {
	client := &fake.MockQualityGatesClient{
		CreateFn: func(opt *sonargo.QualitygatesCreateOption) (*sonargo.QualitygatesCreateObject, *http.Response, error) {
			return &sonargo.QualitygatesCreateObject{
				ID:   "generated-id-12345",
				Name: "ActualSonarQubeName",
			}, nil, nil
		},
	}

	qg := &v1alpha1.QualityGate{
		ObjectMeta: metav1.ObjectMeta{Name: "k8s-resource-name"},
		Spec: v1alpha1.QualityGateSpec{
			ForProvider: v1alpha1.QualityGateParameters{
				Name: "ActualSonarQubeName",
			},
		},
	}

	e := &external{qualityGatesClient: client}
	_, err := e.Create(context.Background(), qg)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify external name is set to the SonarQube name, not the ID
	externalName := meta.GetExternalName(qg)
	if externalName != "ActualSonarQubeName" {
		t.Errorf("Expected external name 'ActualSonarQubeName', got '%s'", externalName)
	}
}

// errComparer compares errors by their message
func errComparer(a, b error) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Error() == b.Error()
}
