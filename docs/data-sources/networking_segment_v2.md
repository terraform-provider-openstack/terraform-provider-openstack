---
subcategory: "Networking / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_networking_segment_v2"
sidebar_current: "docs-openstack-datasource-networking-segment-v2"
description: |-
  Get information on an OpenStack Network Segment.
---

# openstack\_networking\_segment\_v2

Use this data source to get the ID of an available OpenStack network.

## Example Usage

```hcl
data "openstack_networking_segment_v2" "network" {
  name = "tf_test_segment"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve networks ids. If omitted, the
  `region` argument of the provider is used.

* `segment_id` - (Optional) The ID of the network segment

* `name` - (Optional) The name of the network segment.

* `description` - (Optional) Human-readable description of the network segment.

* `network_id` - (Optional) The ID of the network.

* `network_type` - (Optional) The type of the network, such as `vlan`, `vxlan`,
  `flat`, `gre`, `geneve`, or `local`.

* `physical_network` - (Optional) The name of the physical network.

* `segmentation_id` - (Optional) The segmentation ID of the network segment.

* `revision_sumber` - (Optional) The revision number of the network segment.

## Attributes Reference

`id` is set to the ID of the found network. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `segment_id` - See Argument Reference above.
* `network_id` - See Argument Reference above.
* `network_type` - See Argument Reference above.
* `physical_network` - See Argument Reference above.
* `segmentation_id` - See Argument Reference above.
* `revision_number` - See Argument Reference above.
* `created_at` - The date and time when the network segment was created.
* `updated_at` - The date and time when the network segment was last updated.
