---
subcategory: "Networking / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_networking_bgp_peer_v2"
sidebar_current: "docs-openstack-resource-networking-bgp-peer-v2"
description: |-
  Manages a V2 Neutron BGP Peer resource within OpenStack.
---

# openstack\_networking\_bgp\_peer\_v2

Manages a V2 Neutron BGP Peer resource within OpenStack.

This resource allows you to configure a BGP peer that can be associated with a
BGP speaker to exchange routing information.

## Example Usage

```hcl
resource "openstack_networking_bgp_peer_v2" "peer_1" {
  name      = "bgp_peer_1"
  peer_ip   = "192.0.2.10"
  remote_as = 65001
  auth_type = "md5"
  password  = "supersecret"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a Neutron network. If omitted, the
  `region` argument of the provider is used. Changing this creates a new BGP
  peer.

* `tenant_id` – (Optional) The tenant/project ID. Required if admin privileges
  are used. Changing this creates a new BGP peer.

* `name` – (Optional) A name for the BGP peer.

* `auth_type` – (Optional) The authentication type to use. Can be one of `none`
  or `md5`. Defaults to `none`. If set to not `none`, the `password` argument
  must also be provided. Changing this creates a new BGP peer.

* `password` – (Optional) The password used for MD5 authentication. Must be set
  only when `auth_type` is not `none`.

* `remote_as` – (Required) The AS number of the BGP peer. Changing this
  creates a new BGP peer.

* `peer_ip` – (Required) The IP address of the BGP peer. Must be a valid IP
  address. Changing this creates a new BGP peer.

## Attributes Reference

The following attributes are exported:

* `id` – The ID of the BGP peer.
* `tenant_id` – See Argument Reference above.
* `name` – See Argument Reference above.
* `auth_type` – See Argument Reference above.
* `password` – See Argument Reference above.
* `remote_as` – See Argument Reference above.
* `peer_ip` – See Argument Reference above.

## Import

BGP peers can be imported using their ID:

```shell
terraform import openstack_networking_bgp_peer_v2.peer1 a1b2c3d4-e5f6-7890-abcd-1234567890ef
```
