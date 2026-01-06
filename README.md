# provider-sonarqube

`provider-sonarqube` is a minimal [Crossplane](https://crossplane.io/) Provider
that is meant to be used as a sonarqube for implementing new Providers. It comes
with the following features that are meant to be refactored:

- A `ProviderConfig` type that only points to a credentials `Secret`.
- A `MyType` resource type that serves as an example managed resource.
- A managed resource controller that reconciles `MyType` objects and simply
  prints their configuration in its `Observe` method.

## Developing

1. Clone the repository using: `git clone https://github.com/boxboxjason/provider-sonarqube.git`
2. Run `make submodules` to initialize the "build" Make submodule we use for CI/CD.
3. Rename the provider by running the following command:

### Adding a new type

Add your new type by running the following command:

```shell
  export provider_name=SonarQube # Camel case, e.g. GitHub
  export group=instance # lower case e.g. core, cache, database, storage, etc.
  export type=QualityGate # Camel casee.g. Bucket, Database, CacheCluster, etc.
  make provider.addtype provider=${provider_name} group=${group} kind=${type}
```

1. Register your new type into `SetupGated` function in `internal/controller/register.go`
2. Run `make reviewable` to run code generation, linters, and tests.
3. Run `make build` to build the provider.

Refer to Crossplane's [CONTRIBUTING.md] file for more information on how the
Crossplane community prefers to work. The [Provider Development][provider-dev]
guide may also be of use.

[CONTRIBUTING.md]: https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md
[provider-dev]: https://github.com/crossplane/crossplane/blob/master/contributing/guide-provider-development.md
