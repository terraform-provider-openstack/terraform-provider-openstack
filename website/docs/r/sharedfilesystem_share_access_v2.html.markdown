---
layout: "openstack"
page_title: "OpenStack: sharedfilesystem_share_access_v2"
sidebar_current: "docs-openstack-resource-sharedfilesystem-share-access-v2"
description: |-
  Configure a Shared File System share access list.
---

# openstack\_sharedfilesystem\_share\_access\_v2

Use this resource to control the share access lists.

~> **Important Security Notice** The access key assigned by this resource will
be stored *unencrypted* in your Terraform state file. If you use this resource
in production, please make sure your state file is sufficiently protected.
[Read more about sensitive data in
state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage

### NFS

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
  description       = "test share network with security services"
  neutron_net_id    = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}

resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  size             = 1
  share_network_id = "${openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1.id}"
}

resource "openstack_sharedfilesystem_share_access_v2" "share_access_1" {
  share_id     = "${openstack_sharedfilesystem_share_v2.share_1.id}"
  access_type  = "ip"
  access_to    = "192.168.199.10"
  access_level = "rw"
}
```

### CIFS

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
  name              = "test_sharenetwork_secure"
  description       = "share the secure love"
  neutron_net_id    = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  security_service_ids = [
    "${openstack_sharedfilesystem_securityservice_v2.securityservice_1.id}",
  ]
}

resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "cifs_share"
  share_proto      = "CIFS"
  size             = 1
  share_network_id = "${openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1.id}"
}

resource "openstack_sharedfilesystem_share_access_v2" "share_access_1" {
  share_id     = "${openstack_sharedfilesystem_share_v2.share_1.id}"
  access_type  = "user"
  access_to    = "windows"
  access_level = "ro"
}

resource "openstack_sharedfilesystem_share_access_v2" "share_access_2" {
  share_id     = "${openstack_sharedfilesystem_share_v2.share_1.id}"
  access_type  = "user"
  access_to    = "linux"
  access_level = "rw"
}

output "export_locations" {
  value = "${openstack_sharedfilesystem_share_v2.share_1.export_locations}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - The region in which to obtain the V2 Shared File System client.
    A Shared File System client is needed to create a share access. Changing this
    creates a new share access.

* `share_id` - (Required) The UUID of the share to which you are granted access.

* `access_type` - (Required) The access rule type. Can either be an ip, user,
  cert, or cephx. cephx support requires an OpenStack environment that supports
  Shared Filesystem microversion 2.13 (Mitaka) or later.

* `access_to` - (Required) The value that defines the access. Can either be an IP
    address or a username verified by configured Security Service of the Share Network.

* `access_level` - (Required) The access level to the share. Can either be `rw` or `ro`.

## Attributes Reference

* `id` - The unique ID for the Share Access.
* `region` - See Argument Reference above.
* `share_id` - See Argument Reference above.
* `access_type` - See Argument Reference above.
* `access_to` - See Argument Reference above.
* `access_level` - See Argument Reference above.
* `access_key` - The access credential of the entity granted access.

## Import

This resource can be imported by specifying the ID of the share and the ID of the
share access, separated by a slash, e.g.:

```
$ terraform import openstack_sharedfilesystem_share_access_v2.share_access_1 <share id>/<share access id>
```
