---
subcategory: ""
layout: "openstack"
page_title: "Terraform Openstack Provider Version 2 Upgrade Guide"
description: |-
  Terraform Openstack Provider Version 2 Upgrade Guide
---

~> **Note:** The provider go module version has been updated to `v2` in the `2.1.0` release. Downstream users should check the changelog of the `2.1.0` release as well before changing to the `v2` module.

Version 2.0.0 of the Openstack provider for Terraform is a major release and includes some changes that you will need to consider when upgrading. We intend this guide to help with that process.

We previously marked most of the changes we outline in this guide as deprecated in the Terraform plan/apply output throughout previous provider releases. You can find these changes, including deprecation notices, in the [Terraform Openstack Provider CHANGELOG](https://github.com/terraform-provider-openstack/terraform-provider-openstack/blob/main/CHANGELOG.md).

Upgrade topics:
- [Provider Version Configuration](#provider-version-configuration)
- [Neutron LBaaS Deprecation](#neutron-lbaas-deprecation)
  - [Provider configuration option](#provider-configuration-option)
  - [Changes for users of Octavia](#changes-for-users-of-octavia)
  - [Changes for users of Neutron-LBaaS](#changes-for-users-of-neutron-lbaas)
- [Other Removals](#other-removals)
  - [Remove multiattach from volume_v2](#remove-multiattach-from-blockstoragevolumev3)
  - [Remove dhcp_disabled from subnet_v2 data source](#remove-dhcpdisabled-from-networkingsubnetv2-data-source)
  - [Remove update_at from image_v2](#remove-updateat-from-imagesimagev2)
  - [Remove instance_id from blockstorage_volume_attach_v2](#remove-instanceid-from-blockstoragevolumeattachv2)
  - [Remove member from lb_pool_v1](#remove-member-from-lbpoolv1)
  - [Remove allocation_pools from subnet_v2](#remove-allocationpools-from-networkingsubnetv2)
  - [Remove external_gateway from router_v2](#remove-externalgateway-from-networkingrouterv2)
  - [Remove floating_ip from instance_v2](#remove-floatingip-from-computeinstancev2)
  - [Remove volume from instance_v2](#remove-volume-from-computeinstancev2)
  - [Remove sort_key and sort_dir from glance data sources](#remove-sortkey-and-sortdir-from-glance-data-sources)


## Provider Version Configuration

-> Before upgrading to version 2.0.0, upgrade to the most recent 1.X.Y version of the provider and ensure that your environment successfully runs [`terraform plan`](https://www.terraform.io/docs/commands/plan.html). You should not see changes you don't expect.

Use [version constraints when configuring Terraform providers](https://www.terraform.io/docs/configuration/providers.html#provider-versions). If you are following that recommendation, update the version constraints in your Terraform configuration and run [`terraform init -upgrade`](https://www.terraform.io/docs/commands/init.html) to download the new version.

For example, given this previous configuration:

```terraform
terraform {
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.54.1"
    }
  }
}

provider "openstack" {
  # Configuration options
}
```

Update to the latest 2.0.0. version:

```terraform
terraform {
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 2.0.0"
    }
  }
}

provider "openstack" {
  # Configuration options
}
```

## Neutron LBaaS Deprecation

Neutron LBaaS has been deprecated since Queens(2018-02-28) and removed in Train(2019-10-16). Support for it is removed in `2.0.0`. `lb_XYZ_v2` resources support only Octavia.

### Provider configuration option

The provider configuration option `use_octavia` has been removed, since the provider always uses octavia. So the provider configuration below:

```terraform
provider "openstack" {
  ...
  use_octavia = false/true
}
```

Should be updated into:

```terraform
provider "openstack" {
  ...
}
```

### Changes for users of Octavia

The upgrade should be transparent for all loadbalancer resources that are already using Octavia. Therefore there is **no action required** by the user before/after the upgrade.

### Changes for users of Neutron-LBaaS

Neutron-LBaaS will not be supported, therefore any loadbalancer resources **cannot be managed** through terraform-provider-openstack in 2.0.0 and beyond. Cloud-administrator should consider upgrading from Neutron-LBaaS to Octavia.


## Other removals

Various deprecated attributes have been removed from resources and data sources

### Remove multiattach from blockstorage_volume_v3

`multiattach` is no longer possible on `blockstorage_volume_v3` resource. Instead the volume type that the volume uses has to be marked as multiattach as shown below:

```terraform
resource "openstack_blockstorage_volume_type_v3" "multiattach" {
  name        = "multiattach"
  description = "Multiattach-enabled volume type"
  extra_specs = {
      multiattach = "<is> True"
  }
}
```

More details on [multiattach volumes types](https://docs.openstack.org/cinder/latest/admin/volume-multiattach.html#multiattach-volume-type)

### Remove dhcp_disabled from networking_subnet_v2 data source

`dhcp_disabled` has been removed from `openstack_networking_subnet_v2` data source. `dhcp_enabled` should be used instead.

### Remove update_at from images_image_v2

`update_at` has been removed from `openstack_images_image_v2` resource. `updated_at` should be used instead.

### Remove instance_id from blockstorage_volume_attach_v2

`instance_id` has been removed from `openstack_blockstorage_volume_attach_v2` resource. This attribute was deprecated and did not perform any actions or have any functional impact. As such, its removal is transparent to the end-users and requires no changes to your existing configurations.

### Remove member from lb_pool_v1

`member` block has been removed from `openstack_lb_pool_v1` resource. `openstack_lb_member_v1` resource should be used instead. 

### Remove allocation_pools from networking_subnet_v2

`allocation_pools` has been removed from `openstack_networking_subnet_v2` resource. `allocation_pool` should be used instead. It is advisable to switch **before** the upgrade to `2.0.0`. The change will be transparent.

### Remove external_gateway from networking_router_v2

`external_gateway` has been removed from `openstack_networking_router_v2` resource. `external_network_id` should be used instead. It is advisable to switch **before** the upgrade to `2.0.0`. The change will be transparent.

### Remove floating_ip from compute_instance_v2

`floating_ip` has been removed from `openstack_compute_instance_v2` resource. This attribute was deprecated and did not perform any actions or have any functional impact. As such, its removal is transparent to the end-users and requires no changes to your existing configurations.

### Remove volume from compute_instance_v2

`volume` has been removed from `openstack_compute_instance_v2` resource. This attribute was deprecated and did not perform any actions or have any functional impact. As such, its removal is transparent to the end-users and requires no changes to your existing configurations.

### Remove sort_key and sort_dir from glance data sources

`sort_key` and `sort_dir` has been removed from `openstack_images_image_v2` and `openstack_images_image_ids_v2` data sources. `sort` should be used instead.

For example:
```terraform
data "openstack_images_image_ids_v2" "images" {
  name_regex = "^Ubuntu 16\\.04.*-amd64"
  sort_key   = "updated_at"
  sort_dir   = "asc"
}
```

Should be changed into:
```terraform
data "openstack_images_image_ids_v2" "images" {
  name_regex = "^Ubuntu 16\\.04.*-amd64"
  sort       = "updated_at:asc"
}
```

### Remove host_route from networking_subnet_v2

`host_route` block has been removed from `openstack_networking_subnet_v2` resource. This attribute has been deprecated and `openstack_networking_subnet_route_v2` resource should be used instead. The switch should be done before upgrading to `2.0.0` and it might require importing resources. Example:

```terraform
resource "openstack_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  network_id = openstack_networking_network_v2.network_1.id
  cidr       = "192.168.199.0/24"
  ip_version = 4

  host_routes {
    destination_cidr = "10.0.1.0/24"
    next_hop         = "192.168.199.254"
  }
}
```

should be changed into:
```terraform
resource "openstack_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  network_id = openstack_networking_network_v2.network_1.id
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "openstack_networking_subnet_route_v2" "subnet_route_1" {
  subnet_id        = openstack_networking_subnet_v2.subnet_1.id
  destination_cidr = "10.0.1.0/24"
  next_hop         = "192.168.199.254"
}
```

If the user does not want routes to be deleted and recreated, a removal of the subnet resource from the state should be done, subnet resource updated and then import the subnet and subnet_route resource.
