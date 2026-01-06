package common

import (
	"context"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// ErrSecretNotFound is the error string used when a secret cannot be found.
	ErrSecretNotFound = "Cannot find referenced secret"
	// ErrSecretKeyNotFound is the error string used when a key in a secret cannot be found.
	ErrSecretKeyNotFound = "Cannot find key in referenced secret"
	// ErrSecretSelectorNil is the error string used when a secret selector is nil.
	ErrSecretSelectorNil = "Secret selector is nil"
)

// GetTokenValueFromSecret retrieves the token value from the referenced secret.
func GetTokenValueFromSecret(ctx context.Context, client client.Client, m resource.Managed, selector *xpv1.SecretKeySelector) (*string, error) {
	if selector == nil {
		return nil, errors.Errorf(ErrSecretSelectorNil)
	}

	secret := &corev1.Secret{}
	if err := client.Get(ctx, types.NamespacedName{Name: selector.Name, Namespace: selector.Namespace}, secret); err != nil {
		return nil, errors.Wrap(err, ErrSecretNotFound)
	}

	value := secret.Data[selector.Key]
	if value == nil {
		return nil, errors.Errorf(ErrSecretKeyNotFound)
	}

	data := string(value)
	return &data, nil
}

// GetTokenValueFromLocalSecret retrieves the token value from a local secret in the same namespace as the managed resource.
func GetTokenValueFromLocalSecret(ctx context.Context, client client.Client, m resource.Managed, l *xpv1.LocalSecretKeySelector) (*string, error) {
	if l == nil {
		return nil, errors.Errorf(ErrSecretSelectorNil)
	}

	return GetTokenValueFromSecret(ctx, client, m, &xpv1.SecretKeySelector{
		Key: l.Key,
		SecretReference: xpv1.SecretReference{
			Name:      l.Name,
			Namespace: m.GetNamespace(),
		},
	})
}
