Terraform OpenStack Provider
============================

Documentation:
- [registry.terraform.io](https://registry.terraform.io/providers/terraform-provider-openstack/openstack/latest/docs)
- [search.opentofu.org](https://search.opentofu.org/provider/terraform-provider-openstack/openstack/latest)

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 1.x
- [OpenTofu](https://opentofu.org/docs/intro/install) 1.x
- [Go](https://golang.org/doc/install) 1.23 (to build the provider plugin)

Building The Provider
---------------------

Clone the repository

```sh
$ git clone git@github.com:terraform-provider-openstack/terraform-provider-openstack.git
```

Enter the provider directory and build the provider

```sh
$ cd terraform-provider-openstack
$ make build
```

Using the provider
----------------------
Please see the documentation at [registry.terraform.io](https://registry.terraform.io/providers/terraform-provider-openstack/openstack/latest/docs) or [search.opentofu.org](https://search.opentofu.org/provider/terraform-provider-openstack/openstack/latest).

Or you can browse the documentation within this repo [here](https://github.com/terraform-provider-openstack/terraform-provider-openstack/tree/main/website/docs).

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](https://golang.org) installed on your machine (version 1.23+ is *required*).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the current directory.

```sh
$ make build
```

For further details on how to work on this provider, please see the [Testing and Development](https://github.com/terraform-provider-openstack/terraform-provider-openstack/blob/main/docs/index.md#testing-and-development) documentation.

Releasing the Provider
----------------------

This repository contains a GitHub Action configured to automatically build and
publish assets for release when a tag is pushed that matches the pattern `v*`
(ie. `v0.1.0`).

A [Gorelaser](https://goreleaser.com/) configuration is provided that produce
build artifacts matching the [layout required](https://www.terraform.io/docs/registry/providers/publishing.html#manually-preparing-a-release)
to publish the provider in the Terraform Registry.

Releases will as drafts. Once marked as published on the GitHub Releases page,
they will become available via the Terraform Registry.

Before releasing, a PR updating the changelog should be made to trigger the CI 
for all services and ensure that everything is OK. Moreover, update the example
on `website/docs/index.html.markdown` to point to the new version.
