---
subcategory: "Networking / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_networking_bgp_speaker_v2"
sidebar_current: "docs-openstack-datasource-networking-bgp-speaker-v2"
description: |-
  Get information on an OpenStack BGP Speaker.
---

# openstack\_networking\_bgp\_speaker\_v2

Use this data source to get the ID of an available OpenStack BGP Speaker.

## Example Usage

```hcl
data "openstack_networking_bgp_speaker_v2" "speaker" {
  name = "speaker"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve BGP speaker. If omitted, the
  `region` argument of the provider is used.

* `name` - (Optional) The name of the BGP speaker.

* `speaker_id` - (Optional) The ID of the BGP speaker.

## Attributes Reference

`id` is set to the ID of the found BGP speaker. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `tenant_id` - The owner of the BGP speaker.
* `ip_version` - The IP version (4 or 6) of the BGP Speaker.
* `advertise_floating_ip_host_routes` - Whether the advertisement of floating ip host routes by the BGP Speaker is enabled.
* `advertise_tenant_networks` - Whether the advertisement of tenant network routes by the BGP Speaker is enabled.
* `local_as` - The local Autonomous System number of the BGP Speaker.
* `networks` - A list of network IDs to which the BGP Speaker is associated.
* `peers` - A list of BGP peer IDs to which the BGP speaker is associated.
* `dragents` - A list of dynamic routing agent IDs to which the BGP speaker is associated.
* `advertised_routes` - A list of dictionaries containing the `destination` and
  `next_hop` for each route advertised by the BGP speaker.
