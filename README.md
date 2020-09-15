Terraform OpenStack Provider
============================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by:

* Andrei Ozerov ([@ozerovandrei](https://github.com/ozerovandrei))
* Gavin Williams ([@fatmcgav](https://github.com/fatmcgav))
* Joe Topjian ([@jtopjian](https://github.com/jtopjian))

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

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
Please see the documentation at [terraform.io](https://www.terraform.io/docs/providers/openstack/index.html).

Or you can browse the documentation within this repo [here](https://github.com/terraform-provider-openstack/terraform-provider-openstack/tree/master/website/docs).

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](https://golang.org) installed on your machine (version 1.14+ is *required*).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the current directory.

```sh
$ make build
```

For further details on how to work on this provider, please see the [Testing and Development](https://www.terraform.io/docs/providers/openstack/index.html#testing-and-development) documentation.

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

Thank You
---------

We'd like to extend special thanks and appreciation to the following:

### OpenLab

<a href="http://openlabtesting.org/"><img src="assets/openlab.png" width="600px"></a>

OpenLab is providing a full CI environment to test each PR and merge for a variety of OpenStack releases.

### VEXXHOST

<a href="https://vexxhost.com/"><img src="assets/vexxhost.png" width="600px"></a>

VEXXHOST is providing their services to assist with the development and testing of this provider.


