---
subcategory: ""
layout: "openstack"
page_title: "Terraform Openstack Provider Version 3.0 Upgrade Guide"
description: |-
  Terraform Openstack Provider Version 3.0 Upgrade Guide
---

Version 3.0.0 of the OpenStack provider for Terraform is a major release that includes several important changes. This guide is intended to assist you in managing the upgrade process.

Most of the changes described in this guide were previously marked as deprecated in the `terraform plan/apply` output in earlier provider releases. You can review these changes, including deprecation notices, in the [Terraform OpenStack Provider CHANGELOG](https://github.com/terraform-provider-openstack/terraform-provider-openstack/blob/main/CHANGELOG.md).

### Table of Contents:

- [Provider Version Configuration](#provider-version-configuration)
- [Compute Floating IPs API Removal](#compute-floating-ips-api-removal)
- [Compute Security Groups API Removal](#compute-security-groups-api-removal)
- [Compute Project Networks API Removal](#compute-project-networks-api-removal)
- [Other Resource and Data Source Removals](#other-resource-and-data-source-removals)

## Provider Version Configuration

When configuring Terraform providers, always use [version constraints](https://www.terraform.io/docs/configuration/providers.html#provider-versions). If you have implemented version constraints, update them accordingly in your configuration and run [`terraform init -upgrade`](https://www.terraform.io/docs/commands/init.html) to download the latest version.

For example, if your previous configuration looked like this:

```terraform
terraform {
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 2.1.0"
    }
  }
}

provider "openstack" {
  # Configuration options
}
```

You should update it to the following:

```terraform
terraform {
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 3.0.0"
    }
  }
}

provider "openstack" {
  # Configuration options
}
```

## Compute Floating IPs API Removal

The Compute Floating IPs API has been [deprecated](https://docs.openstack.org/api-ref/compute/#floating-ips-os-floating-ips-deprecated) in the OpenStack Nova service and [removed](https://github.com/gophercloud/gophercloud/blob/master/docs/MIGRATING.md) from the Gophercloud v2 SDK. This removal requires manual changes in your existing Terraform configuration.

For example, the following configuration:

```terraform
resource "openstack_compute_floatingip_v2" "floatip_1" {
  pool = "public"
}

resource "openstack_compute_floatingip_associate_v2" "fip" {
  instance_id = openstack_compute_instance_v2.my_instance.id
  floating_ip = openstack_compute_floatingip_v2.floatip_1.address
}
```

should be changed to:

```terraform
resource "openstack_networking_floatingip_v2" "floatip_1" {
  pool = "public"
}

data "openstack_networking_port_v2" "port" {
  device_id  = openstack_compute_instance_v2.my_instance.id
  network_id = openstack_compute_instance_v2.my_instance.network.0.uuid
}

resource "openstack_networking_floatingip_associate_v2" "fip_associate" {
  floating_ip = openstack_networking_floatingip_v2.floatip_1.address
  port_id     = data.openstack_networking_port_v2.port.id
}
```

## Compute Security Groups API Removal

The Compute Security Groups API has been [deprecated](https://docs.openstack.org/api-ref/compute/#security-groups-os-security-groups-deprecated) in the OpenStack Nova service and [removed](https://github.com/gophercloud/gophercloud/blob/master/docs/MIGRATING.md) from the Gophercloud v2 SDK. This change is not transparent to users and requires updates to your configuration.

For example, the following configuration:

```terraform
resource "openstack_compute_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "a security group"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}
```

should be updated to:

```terraform
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "a security group"
}

resource "openstack_networking_secgroup_rule_v2" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.secgroup_1.id
}
```

## Compute Project Networks API Removal

The Compute Project Networks API has been [deprecated](https://docs.openstack.org/api-ref/compute/#project-networks-os-tenant-networks-deprecated) in the OpenStack Nova service and [removed](https://github.com/gophercloud/gophercloud/blob/master/docs/MIGRATING.md) from the Gophercloud v2 SDK. This change is transparent to users, and no modifications to existing configurations are required.

This API was used in the `openstack_compute_instance_v2` resource when the `OS_NOVA_NETWORK` environment variable was set. If you are using this variable, it should be removed from your configuration, and you should switch to using the Neutron API.

## Other Resource and Data Source Removals

Several data sources have been removed in this version:

| Data Source | Replacement |
|-------------|-------------|
| `openstack_blockstorage_snapshot_v2` | `openstack_blockstorage_snapshot_v3` |
| `openstack_blockstorage_volume_v2`   | `openstack_blockstorage_volume_v3`   |
| `openstack_fw_policy_v1`             | `openstack_fw_policy_v2`             |

In addition, the following resources have been removed:

| Resource | Replacement |
|----------|-------------|
| `openstack_blockstorage_quotaset_v2`           | `openstack_blockstorage_quotaset_v3`           |
| `openstack_blockstorage_volume_v1`             | `openstack_blockstorage_volume_v3`             |
| `openstack_blockstorage_volume_v2`             | `openstack_blockstorage_volume_v3`             |
| `openstack_blockstorage_volume_attach_v2`      | `openstack_blockstorage_volume_attach_v3`      |
| `openstack_fw_firewall_v1`                     | `openstack_fw_group_v2`                        |
| `openstack_fw_policy_v1`                       | `openstack_fw_policy_v2`                       |
| `openstack_fw_rule_v1`                         | `openstack_fw_rule_v2`                         |
| `openstack_lb_member_v1`                       | `openstack_lb_member_v2`                       |
| `openstack_lb_monitor_v1`                      | `openstack_lb_monitor_v2`                      |
| `openstack_lb_pool_v1`                         | `openstack_lb_pool_v2`                         |
| `openstack_lb_vip_v1`                          | `openstack_lb_loadbalancer_v2`, `openstack_lb_listener_v2`, `openstack_networking_floatingip_associate_v2` |
