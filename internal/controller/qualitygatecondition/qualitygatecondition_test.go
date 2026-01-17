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

package qualitygatecondition

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

type notQualityGateCondition struct {
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
		"NotQualityGateConditionError": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg:  &notQualityGateCondition{},
			},
			want: want{
				o:   managed.ExternalObservation{},
				err: errors.New(errNotQualityGateCondition),
			},
		},
		"EmptyExternalNameReturnsNotExists": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGateCondition{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "test-condition",
						Annotations: map[string]string{},
					},
				},
			},
			want: want{
				o:   managed.ExternalObservation{ResourceExists: false},
				err: nil,
			},
		},
		"ShowFailsReturnsError": {
			client: &fake.MockQualityGatesClient{
				ShowFn: func(opt *sonargo.QualitygatesShowOption) (*sonargo.QualitygatesShowObject, *http.Response, error) {
					return nil, nil, errors.New("api error")
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGateCondition {
					qgc := &v1alpha1.QualityGateCondition{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-condition",
							Annotations: map[string]string{},
						},
					}
					meta.SetExternalName(qgc, "cond-123")
					return qgc
				}(),
			},
			want: want{
				o:   managed.ExternalObservation{ResourceExists: false},
				err: errors.Wrap(errors.New("api error"), errShowQualityGateCondition),
			},
		},
		"ConditionNotFoundInQualityGate": {
			client: &fake.MockQualityGatesClient{
				ShowFn: func(opt *sonargo.QualitygatesShowOption) (*sonargo.QualitygatesShowObject, *http.Response, error) {
					return &sonargo.QualitygatesShowObject{
						Name:       "test-gate",
						Conditions: []sonargo.QualitygatesShowObject_sub2{},
					}, nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGateCondition {
					qgc := &v1alpha1.QualityGateCondition{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-condition",
							Annotations: map[string]string{},
						},
						Spec: v1alpha1.QualityGateConditionSpec{
							ForProvider: v1alpha1.QualityGateConditionParameters{
								QualityGateName: ptr.To("test-gate"),
								Metric:          "coverage",
								Error:           "80",
							},
						},
					}
					meta.SetExternalName(qgc, "cond-123")
					return qgc
				}(),
			},
			want: want{
				o:   managed.ExternalObservation{ResourceExists: false},
				err: errors.Wrap(errors.New("quality gate condition not found in observation"), errShowQualityGateCondition),
			},
		},
		"SuccessfulObserveResourceExists": {
			client: &fake.MockQualityGatesClient{
				ShowFn: func(opt *sonargo.QualitygatesShowOption) (*sonargo.QualitygatesShowObject, *http.Response, error) {
					return &sonargo.QualitygatesShowObject{
						Name: "test-gate",
						Conditions: []sonargo.QualitygatesShowObject_sub2{
							{
								ID:     "cond-123",
								Metric: "coverage",
								Op:     "LT",
								Error:  "80",
							},
						},
					}, nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGateCondition {
					qgc := &v1alpha1.QualityGateCondition{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-condition",
							Annotations: map[string]string{},
						},
						Spec: v1alpha1.QualityGateConditionSpec{
							ForProvider: v1alpha1.QualityGateConditionParameters{
								QualityGateName: ptr.To("test-gate"),
								Metric:          "coverage",
								Error:           "80",
								Op:              ptr.To("LT"),
							},
						},
					}
					meta.SetExternalName(qgc, "cond-123")
					return qgc
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
		"ResourceNotUpToDateWhenErrorsDiffer": {
			client: &fake.MockQualityGatesClient{
				ShowFn: func(opt *sonargo.QualitygatesShowOption) (*sonargo.QualitygatesShowObject, *http.Response, error) {
					return &sonargo.QualitygatesShowObject{
						Name: "test-gate",
						Conditions: []sonargo.QualitygatesShowObject_sub2{
							{
								ID:     "cond-123",
								Metric: "coverage",
								Op:     "LT",
								Error:  "85",
							},
						},
					}, nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGateCondition {
					qgc := &v1alpha1.QualityGateCondition{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-condition",
							Annotations: map[string]string{},
						},
						Spec: v1alpha1.QualityGateConditionSpec{
							ForProvider: v1alpha1.QualityGateConditionParameters{
								QualityGateName: ptr.To("test-gate"),
								Metric:          "coverage",
								Error:           "80",
								Op:              ptr.To("LT"),
							},
						},
					}
					meta.SetExternalName(qgc, "cond-123")
					return qgc
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
		"LateInitializeOp": {
			client: &fake.MockQualityGatesClient{
				ShowFn: func(opt *sonargo.QualitygatesShowOption) (*sonargo.QualitygatesShowObject, *http.Response, error) {
					return &sonargo.QualitygatesShowObject{
						Name: "test-gate",
						Conditions: []sonargo.QualitygatesShowObject_sub2{
							{
								ID:     "cond-123",
								Metric: "coverage",
								Op:     "LT",
								Error:  "80",
							},
						},
					}, nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGateCondition {
					qgc := &v1alpha1.QualityGateCondition{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-condition",
							Annotations: map[string]string{},
						},
						Spec: v1alpha1.QualityGateConditionSpec{
							ForProvider: v1alpha1.QualityGateConditionParameters{
								QualityGateName: ptr.To("test-gate"),
								Metric:          "coverage",
								Error:           "80",
								Op:              nil,
							},
						},
					}
					meta.SetExternalName(qgc, "cond-123")
					return qgc
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
		"NotQualityGateConditionError": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg:  &notQualityGateCondition{},
			},
			want: want{
				o:   managed.ExternalCreation{},
				err: errors.New(errNotQualityGateCondition),
			},
		},
		"CreateFails": {
			client: &fake.MockQualityGatesClient{
				CreateConditionFn: func(opt *sonargo.QualitygatesCreateConditionOption) (*sonargo.QualitygatesCreateConditionObject, *http.Response, error) {
					return nil, nil, errors.New("create error")
				},
			},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGateCondition{
					ObjectMeta: metav1.ObjectMeta{Name: "test-condition"},
					Spec: v1alpha1.QualityGateConditionSpec{
						ForProvider: v1alpha1.QualityGateConditionParameters{
							QualityGateName: ptr.To("test-gate"),
							Metric:          "coverage",
							Error:           "80",
						},
					},
				},
			},
			want: want{
				o:   managed.ExternalCreation{},
				err: errors.Wrap(errors.New("create error"), errCreateQualityGateCondition),
			},
		},
		"SuccessfulCreate": {
			client: &fake.MockQualityGatesClient{
				CreateConditionFn: func(opt *sonargo.QualitygatesCreateConditionOption) (*sonargo.QualitygatesCreateConditionObject, *http.Response, error) {
					return &sonargo.QualitygatesCreateConditionObject{
						ID:     "cond-123",
						Metric: opt.Metric,
						Error:  opt.Error,
					}, nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGateCondition{
					ObjectMeta: metav1.ObjectMeta{Name: "test-condition"},
					Spec: v1alpha1.QualityGateConditionSpec{
						ForProvider: v1alpha1.QualityGateConditionParameters{
							QualityGateName: ptr.To("test-gate"),
							Metric:          "coverage",
							Error:           "80",
						},
					},
				},
			},
			want: want{
				o:   managed.ExternalCreation{},
				err: nil,
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
		"NotQualityGateConditionError": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg:  &notQualityGateCondition{},
			},
			want: want{
				o:   managed.ExternalUpdate{},
				err: errors.New(errNotQualityGateCondition),
			},
		},
		"EmptyExternalNameReturnsError": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGateCondition{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "test-condition",
						Annotations: map[string]string{},
					},
				},
			},
			want: want{
				o:   managed.ExternalUpdate{},
				err: fmt.Errorf("external name is not set for Quality Gate Condition %s", "test-condition"),
			},
		},
		"UpdateFails": {
			client: &fake.MockQualityGatesClient{
				UpdateConditionFn: func(opt *sonargo.QualitygatesUpdateConditionOption) (*http.Response, error) {
					return nil, errors.New("update error")
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGateCondition {
					qgc := &v1alpha1.QualityGateCondition{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-condition",
							Annotations: map[string]string{},
						},
						Spec: v1alpha1.QualityGateConditionSpec{
							ForProvider: v1alpha1.QualityGateConditionParameters{
								QualityGateName: ptr.To("test-gate"),
								Metric:          "coverage",
								Error:           "85",
							},
						},
					}
					meta.SetExternalName(qgc, "cond-123")
					return qgc
				}(),
			},
			want: want{
				o:   managed.ExternalUpdate{},
				err: errors.Wrap(errors.New("update error"), errUpdateQualityGateCondition),
			},
		},
		"SuccessfulUpdate": {
			client: &fake.MockQualityGatesClient{
				UpdateConditionFn: func(opt *sonargo.QualitygatesUpdateConditionOption) (*http.Response, error) {
					return nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGateCondition {
					qgc := &v1alpha1.QualityGateCondition{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-condition",
							Annotations: map[string]string{},
						},
						Spec: v1alpha1.QualityGateConditionSpec{
							ForProvider: v1alpha1.QualityGateConditionParameters{
								QualityGateName: ptr.To("test-gate"),
								Metric:          "coverage",
								Error:           "85",
							},
						},
					}
					meta.SetExternalName(qgc, "cond-123")
					return qgc
				}(),
			},
			want: want{
				o:   managed.ExternalUpdate{},
				err: nil,
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
		"NotQualityGateConditionError": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg:  &notQualityGateCondition{},
			},
			want: want{
				o:   managed.ExternalDelete{},
				err: errors.New(errNotQualityGateCondition),
			},
		},
		"EmptyExternalNameDoesNothing": {
			client: &fake.MockQualityGatesClient{},
			args: args{
				ctx: context.Background(),
				mg: &v1alpha1.QualityGateCondition{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "test-condition",
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
				DeleteConditionFn: func(opt *sonargo.QualitygatesDeleteConditionOption) (*http.Response, error) {
					return nil, nil
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGateCondition {
					qgc := &v1alpha1.QualityGateCondition{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-condition",
							Annotations: map[string]string{},
						},
					}
					meta.SetExternalName(qgc, "cond-123")
					return qgc
				}(),
			},
			want: want{
				o:   managed.ExternalDelete{},
				err: nil,
			},
		},
		"DeleteFails": {
			client: &fake.MockQualityGatesClient{
				DeleteConditionFn: func(opt *sonargo.QualitygatesDeleteConditionOption) (*http.Response, error) {
					return nil, errors.New("delete error")
				},
			},
			args: args{
				ctx: context.Background(),
				mg: func() *v1alpha1.QualityGateCondition {
					qgc := &v1alpha1.QualityGateCondition{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test-condition",
							Annotations: map[string]string{},
						},
					}
					meta.SetExternalName(qgc, "cond-123")
					return qgc
				}(),
			},
			want: want{
				o:   managed.ExternalDelete{},
				err: errors.Wrap(errors.New("delete error"), errDeleteQualityGateCondition),
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
