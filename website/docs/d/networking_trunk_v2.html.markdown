---
layout: "openstack"
page_title: "OpenStack: openstack_networking_trunk_v2"
sidebar_current: "docs-openstack-datasource-networking-trunk-v2"
description: |-
  Get information of an OpenStack Trunk.
---

# openstack\_networking\_trunk\_v2

Use this data source to get the ID of an available OpenStack trunk.

## Example Usage

```hcl
data "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve trunk ids. If omitted, the
  `region` argument of the provider is used.

* `project_id` - (Optional) The owner of the trunk.

* `trunk_id` - (Optional) The ID of the trunk.

* `name` - (Optional) The name of the trunk.

* `description` - (Optional) Human-readable description of the trunk.

* `port_id` - (Optional) The ID of the trunk parent port.

* `admin_state_up` - (Optional) The administrative state of the trunk.

* `status` - (Optional) The status of the trunk.

* `tags` - (Optional) The list of trunk tags to filter.

## Attributes Reference

`id` is set to the ID of the found trunk. In addition, the following attributes
are exported:

* `all_tags` - The set of string tags applied on the trunk.

* `sub_port` - The set of the trunk subports. The structure of each subport is
   described below.

The `sub_port` attribute has fields below:

* `port_id` - The ID of the trunk subport.

* `segmentation_type` - The segmenation tecnology used, e.g., "vlan".

* `segmentation_id` - The numeric id of the subport segment.
