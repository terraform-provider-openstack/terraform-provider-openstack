---
subcategory: "Networking / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_networking_segment_v2"
sidebar_current: "docs-openstack-resource-networking-segment-v2"
description: |-
  Manages a Neutron network segment resource within OpenStack.
---

# openstack\_networking\_segment\_v2

Manages a Neutron network segment resource within OpenStack.

~> **Note:** This resource is only available if the Neutron service is
configured with the `segments` extension.

~> **Note:** This ussually requires admin privileges to create or manage
segments.

## Example Usage

```hcl
resource "openstack_networking_network_v2" "net_1" {
  name = "demo-net"
}

resource "openstack_networking_segment_v2" "segment_1" {
  name             = "flat-segment"
  description      = "Example flat segment"
  network_id       = openstack_networking_network_v2.net_1.id
  network_type     = "flat"
  physical_network = "public"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a Neutron network. If omitted, the
  `region` argument of the provider is used. Changing this creates a new
  segment.

* `name` – (Optional) A name for the segment.

* `description` – (Optional) A description for the segment.

* `network_id` – (Required) The UUID of the network this segment belongs to.
  Changing this will create a new segment.

* `network_type` – (Required) The network type. Valid values depend on the
  backend (e.g., `vlan`, `vxlan`, `flat`, `gre`, `geneve`, `local`). Changing
  this will create a new segment.

* `physical_network` – (Optional) The name of the physical network. Changing this
  will create a new segment.

* `segmentation_id` – (Optional) A segmentation identifier. Changing is allowed
  only for `vlan`.

## Attributes Reference

The following attributes are exported:

* `id` – The ID of the segment.
* `region` – See Argument Reference above.
* `name` – See Argument Reference above.
* `description` – See Argument Reference above.
* `network_id` – The ID of the network this segment belongs to.
* `network_type` – The type of the network segment.
* `physical_network` – The name of the physical network.
* `segmentation_id` – The segmentation identifier, if applicable.
* `revision_number` – The revision number of the segment.
* `created_at` – Creation timestamp (RFC3339 format).
* `updated_at` – Last update timestamp (RFC3339 format).

## Import

This resource can be imported by specifying the segment ID:

```shell
terraform import openstack_networking_segment_v2.segment1 a5e3a494-26ee-4fde-ad26-2d846c47072e
```
