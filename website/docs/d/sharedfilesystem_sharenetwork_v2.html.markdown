---
layout: "openstack"
page_title: "OpenStack: openstack_sharedfilesystem_sharenetwork_v2"
sidebar_current: "docs-openstack-datasource-sharedfilesystem-sharenetwork-v2"
description: |-
  Get information on an Shared File System share network.
---

# openstack\_sharedfilesystem\_sharenetwork\_v2

Use this data source to get the ID of an available Shared File System share network.

## Example Usage

```hcl
data "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name = "sharenetwork_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Shared File System client.
    A Shared File System client is needed to read a share network. If omitted, the
    `region` argument of the provider is used.

* `name` - (Optional) The name of the share network.

* `description` - (Optional) The human-readable description of the share network.

* `project_id` - (Optional) The owner of the share network.

* `neutron_net_id` - (Optional) The neutron network UUID of the share network.

* `neutron_subnet_id` - (Optional) The neutron subnet UUID of the share network.

* `security_service_id` - (Optional) The security service IDs associated with
    the share network.

* `network_type` - (Optional) The share network type. Can either be VLAN, VXLAN,
    GRE, or flat.

* `segmentation_id` - (Optional) The share network segmentation ID.

* `cidr` - (Optional) The share network CIDR.

* `ip_version` - (Optional) The IP version of the share network. Can either be 4 or 6.

## Attributes Reference

`id` is set to the ID of the found share network . In addition, the following
attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - The owner of the Share Network.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `neutron_net_id` - See Argument Reference above.
* `neutron_subnet_id` - See Argument Reference above.
* `security_service_id` - See Argument Reference above.
* `network_type` - See Argument Reference above.
* `segmentation_id` - See Argument Reference above.
* `cidr` - See Argument Reference above.
* `ip_version` - See Argument Reference above.
* `security_service_ids` - The list of security service IDs associated with
    the share network.
