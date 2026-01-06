package v1alpha1

import (
	"context"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reference"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	errResolveQualityGateRef = "cannot resolve qualityGateName reference"
)

// ResolveReferences resolves all the references of this QualityGateCondition
// Currently, it resolves the following references:
// - spec.forProvider.qualityGateName -> QualityGate
func (mg *QualityGateCondition) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPINamespacedResolver(c, mg)

	// resolve spec.forProvider.qualityGateName
	rsp, err := r.Resolve(ctx, reference.NamespacedResolutionRequest{
		CurrentValue: *mg.Spec.ForProvider.QualityGateName,
		Reference:    mg.Spec.ForProvider.QualityGateRef,
		Selector:     mg.Spec.ForProvider.QualityGateSelector,
		To:           reference.To{Managed: &QualityGate{}, List: &QualityGateList{}},
		Extract:      reference.ExternalName(),
	})

	if err != nil {
		return errors.Wrap(err, errResolveQualityGateRef)
	}

	resolvedName := &rsp.ResolvedValue

	mg.Spec.ForProvider.QualityGateName = resolvedName
	mg.Spec.ForProvider.QualityGateRef = rsp.ResolvedReference

	return nil
}
