# provider-sonarqube

## Overview

`provider-sonarqube` brings SonarQube configuration under Crossplane management. Install
it into a control plane and declare SonarQube resources as Kubernetes manifests instead of
managing them through ad hoc scripts or manual clicks.

The provider currently covers the areas teams usually need to standardize first:

* Instance resources such as projects, quality gates, quality profiles, rules, and settings
* IAM resources such as groups and permissions templates
* Integration resources such as GitLab ALM configuration

## Why This Provider

Use `provider-sonarqube` when you want SonarQube configuration to behave like the rest of
your platform:

* GitOps-friendly reconciliation for SonarQube state
* Crossplane compositions and abstractions for reusable platform building blocks
* A single declarative model for shared quality policies and project configuration

## Getting Started and Documentation

1. Create a [User Token](https://docs.sonarsource.com/sonarqube-server/user-guide/managing-tokens) on your SonarQube instance. It should preferably have admin permissions to ensure you can manage all resource types.
2. Store it in a Kubernetes secret:

    ```bash
    kubectl create secret generic example-provider-secret -n default --from-literal=token="<USER_TOKEN>"
    ```

3. Configure a `ProviderConfig` that points at your SonarQube instance. You can use token-based authentication or basic auth.

    ```yaml
    apiVersion: sonarqube.crossplane.io/v1alpha1
    kind: ProviderConfig
    metadata:
      name: example
      namespace: default
    spec:
      baseUrl: http://sonarqube.example.com/api
      token:
        source: Secret
        secretRef:
          namespace: default
          name: example-provider-secret
          key: token
    ```

4. Apply the example manifests in this repository to get started quickly.

    ```bash
    kubectl apply -f examples/providerconfig.yaml
    ```

Available examples live under `examples/` and are a practical starting point for
understanding the resource shapes supported by the provider.

## Documentation

* [CRD documentation](https://marketplace.upbound.io/providers/crossplane-contrib/provider-sonarqube/latest/crds)
* [Contributing guide](CONTRIBUTING.md)
* [Issue tracker](https://github.com/crossplane-contrib/provider-sonarqube/issues)

## Support

If you run into a problem or want to request a feature, please open an issue in the
[provider repository](https://github.com/crossplane-contrib/provider-sonarqube/issues).

## Licensing

`provider-sonarqube` is licensed under Apache 2.0. See [LICENSE](LICENSE) for full license text.

[![FOSSA
Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fcrossplane-contrib%2Fprovider-sonarqube.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fcrossplane-contrib%2Fprovider-sonarqube?ref=badge_large)
