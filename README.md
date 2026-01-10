# Terraform OpenStack Provider

The **Terraform OpenStack Provider** allows you to manage OpenStack resources using Terraform or OpenTofu.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.x
- [OpenTofu](https://opentofu.org/docs/intro/install) >= 1.x
- [Go](https://golang.org/doc/install) >= 1.24 (to build the provider plugin)

## Building the Provider

Clone the repository:

```shell
git clone git@github.com:terraform-provider-openstack/terraform-provider-openstack.git
```

Enter the provider directory and build the provider:

```shell
cd terraform-provider-openstack
make build
```

## Using the provider

For usage instructions and examples, refer to the documentation:

- [Terraform Registry](https://registry.terraform.io/providers/terraform-provider-openstack/openstack/latest/docs)
- [OpenTofu Registry](https://search.opentofu.org/provider/terraform-provider-openstack/openstack/latest)
- Or browse the [documentation in this repository](https://github.com/terraform-provider-openstack/terraform-provider-openstack/tree/main/docs)

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the current directory.

```shell
make build
```

For further details on how to work on this provider, please see the [Testing and Development](https://github.com/terraform-provider-openstack/terraform-provider-openstack/blob/main/docs/index.md#testing-and-development) documentation.

## Releasing the Provider

This repository contains a GitHub Action configured to automatically build and
publish assets for release when a tag is pushed that matches the pattern `v*`
(ie. `v0.1.0`).

A [Goreleaser](https://goreleaser.com/) configuration is provided that produce
build artifacts matching the [layout required](https://developer.hashicorp.com/terraform/registry/providers/publishing)
to publish the provider in the Terraform Registry.

Releases will as drafts. Once marked as published on the GitHub Releases page,
they will become available via the Terraform Registry.

Before releasing, a PR updating the changelog should be made to trigger the CI
for all services and ensure that everything is OK. Moreover, update the example
on `docs/index.md` to point to the new version.
