---
layout: "openstack"
page_title: "OpenStack: sharedfilesystem_sharenetwork_v2"
sidebar_current: "docs-openstack-resource-sharedfilesystem-sharenetwork-v2"
description: |-
  Configure a Shared File System share network.
---

# sharedfilesystem\_sharenetwork\_v2

Use this resource to configure a share network.

A share network stores network information that share servers can use when
shares are created.

## Example Usage

### Basic share network

```hcl
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name              = "test_sharenetwork"
  description       = "test share network"
  neutron_net_id    = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}
```

### Share network with associated security services

```hcl
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_sharedfilesystem_securityservice_v2" "securityservice_1" {
  name        = "security"
  description = "created by terraform"
  type        = "active_directory"
  server      = "192.168.199.10"
  dns_ip      = "192.168.199.10"
  domain      = "example.com"
  ou          = "CN=Computers,DC=example,DC=com"
  user        = "joinDomainUser"
  password    = "s8cret"
}

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name              = "test_sharenetwork"
  description       = "test share network with security services"
  neutron_net_id    = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  security_service_ids = [
    "${openstack_sharedfilesystem_securityservice_v2.securityservice_1.id}",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Shared File System client.
    A Shared File System client is needed to create a share network. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    share network.

* `name` - (Optional) The name for the share network. Changing this updates the name
    of the existing share network.

* `description` - (Optional) The human-readable description for the share network.
    Changing this updates the description of the existing share network.

* `neutron_net_id` - (Required) The UUID of a neutron network when setting up or updating
    a share network. Changing this updates the existing share network if it's not used by
    shares.

* `neutron_subnet_id` - (Required) The UUID of the neutron subnet when setting up or
    updating a share network. Changing this updates the existing share network if it's
    not used by shares.

* `security_service_ids` - (Optional) The list of security service IDs to associate with
    the share network. The security service must be specified by ID and not name.

## Attributes Reference

* `id` - The unique ID for the Share Network.
* `region` - See Argument Reference above.
* `project_id` - The owner of the Share Network.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `neutron_net_id` - See Argument Reference above.
* `neutron_subnet_id` - See Argument Reference above.
* `security_service_ids` - See Argument Reference above.
* `network_type` - The share network type. Can either be VLAN, VXLAN, GRE, or flat.
* `segmentation_id` - The share network segmentation ID.
* `cidr` - The share network CIDR.
* `ip_version` - The IP version of the share network. Can either be 4 or 6.

## Import

This resource can be imported by specifying the ID of the share network:

```
$ terraform import openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1 <id>
```
