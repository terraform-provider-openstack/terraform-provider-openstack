---
layout: "openstack"
page_title: "OpenStack: sharedfilesystem_share_v2"
sidebar_current: "docs-openstack-resource-sharedfilesystem-share-v2"
description: |-
  Configure a Shared File System share.
---

# openstack\_sharedfilesystem\_share\_v2

Use this resource to configure a share.

## Example Usage

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
```

## Argument Reference

The following arguments are supported:

* `region` - The region in which to obtain the V2 Shared File System client.
    A Shared File System client is needed to create a share. Changing this
    creates a new share.

* `name` - (Optional) The name of the share. Changing this updates the name
    of the existing share.

* `description` - (Optional) The human-readable description for the share.
    Changing this updates the description of the existing share.

* `share_proto` - (Required) The share protocol - can either be NFS, CIFS,
    CEPHFS, GLUSTERFS, HDFS or MAPRFS. Changing this creates a new share.

* `size` - (Required) The share size, in GBs. The requested share size cannot be greater
    than the allowed GB quota. Changing this resizes the existing share.

* `share_type` - (Optional) The share type name. If you omit this parameter, the default
    share type is used.

* `snapshot_id` - (Optional) The UUID of the share's base snapshot. Changing this creates
    a new share.

* `is_public` - (Optional) The level of visibility for the share. Set to true to make
    share public. Set to false to make it private. Default value is false. Changing this
    updates the existing share.

* `metadata` - (Optional) One or more metadata key and value pairs as a dictionary of
    strings.

* `share_network_id` - (Optional) The UUID of a share network where the share server exists
    or will be created. If `share_network_id` is not set and you provide a `snapshot_id`,
    the share_network_id value from the snapshot is used. Changing this creates a new share.

* `availability_zone` - (Optional) The share availability zone. Changing this creates a
    new share.

## Attributes Reference

* `id` - The unique ID for the Share.
* `region` - See Argument Reference above.
* `project_id` - The owner of the Share.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `share_proto` - See Argument Reference above.
* `size` - See Argument Reference above.
* `share_type` - See Argument Reference above.
* `snapshot_id` - See Argument Reference above.
* `is_public` - See Argument Reference above.
* `metadata` - See Argument Reference above.
* `share_network_id` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `export_locations` - A list of export locations. For example, when a share server
    has more than one network interface, it can have multiple export locations.
* `has_replicas` - Indicates whether a share has replicas or not.
* `host` - The share host name.
* `replication_type` - The share replication type.
* `share_server_id` - The UUID of the share server.
* `all_metadata` - The map of metadata, assigned on the share, which has been
  explicitly and implicitly added.

## Import

This resource can be imported by specifying the ID of the share:

```
$ terraform import openstack_sharedfilesystem_share_v2.share_1 <id>
```
