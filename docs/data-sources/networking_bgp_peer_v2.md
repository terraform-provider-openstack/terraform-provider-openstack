---
subcategory: "Networking / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_networking_bgp_peer_v2"
sidebar_current: "docs-openstack-datasource-networking-bgp-peer-v2"
description: |-
  Get information on an OpenStack BGP peer.
---

# openstack\_networking\_bgp\_peer\_v2

Use this data source to get the ID of an available OpenStack BGP peer.

## Example Usage

```hcl
data "openstack_networking_bgp_peer_v2" "peer" {
  name = "peer"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve BGP peers. If omitted, the
  `region` argument of the provider is used.

* `name` - (Optional) The name of the BGP peer.

* `peer_id` - (Optional) The ID of the BGP peer.

## Attributes Reference

`id` is set to the ID of the found BGP peer. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `tenant_id` - The owner of the BGP peer.
* `auth_type` - Authentication algorithm.
* `peer_ip` - Peer IP address.
* `remote_as` - Peer AS number.
