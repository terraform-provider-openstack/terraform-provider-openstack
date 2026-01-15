---
subcategory: "Networking / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_networking_bgp_speaker_v2"
sidebar_current: "docs-openstack-resource-networking-bgp-speaker-v2"
description: |-
  Manages a V2 Neutron BGP Speaker resource within OpenStack.
---

# openstack\_networking\_bgp\_speaker\_v2

Manages a V2 Neutron BGP Speaker resource within OpenStack.

This resource allows you to configure a BGP speaker that can be associated with
a BGP peer to exchange routing information.

## Example Usage

```hcl
resource "openstack_networking_network_v2" "network1" {
  name = "network1"
}

resource "openstack_networking_bgp_peer_v2" "peer_1" {
  name      = "bgp_peer_1"
  peer_ip   = "192.0.2.10"
  remote_as = 65001
  auth_type = "md5"
  password  = "supersecret"
}

resource "openstack_networking_bgp_speaker_v2" "speaker_1" {
  name       = "bgp_speaker_1"
  ip_version = 4
  local_as   = 65000

  networks = [
    openstack_networking_network_v2.network1.id,
  ]

  peers = [
    opestack_networking_bgp_peer_v2.peer1.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a Neutron network. If omitted, the
  `region` argument of the provider is used. Changing this creates a new BGP
  speaker.

* `tenant_id` – (Optional) The tenant/project ID. Required if admin privileges
  are used. Changing this creates a new BGP speaker.

* `name` – (Optional) A name for the BGP speaker.

* `ip_version` – (Required) The IP version of the BGP speaker. Valid values are
  `4` or `6`. Defaults to `4`. Changing this creates a new BGP speaker.

* `advertise_floating_ip_host_routes` – (Optional) A boolean value indicating
  whether to advertise floating IP host routes. Defaults to `true`.

* `advertise_tenant_networks` - (Optional) A boolean value indicating whether to
  advertise tenant networks. Defaults to `true`.

* `local_as` – (Required) The local autonomous system number (ASN) for the BGP
  speaker. This is a mandatory field and must be specified. Changing this
  creates a new BGP speaker.

* `networks` - (Optional) A list of network IDs to associate with the BGP speaker.

* `peers` - (Optional) A list of BGP peer IDs to associate with the BGP speaker.

## Attributes Reference

The following attributes are exported:

* `id` – The ID of the BGP speaker.
* `tenant_id` – See Argument Reference above.
* `name` – See Argument Reference above.
* `ip_version` – See Argument Reference above.
* `advertise_floating_ip_host_routes` – See Argument Reference above.
* `advertise_tenant_networks` – See Argument Reference above.
* `local_as` – See Argument Reference above.
* `networks` – See Argument Reference above.
* `peers` – See Argument Reference above.
* `advertised_routes` – A list of dictionaries containing the `destination` and
  `next_hop` for each route advertised by the BGP speaker. This attribute is
  only populated after the BGP speaker has been created and has established BGP
  sessions with its peers.

## Import

BGP speakers can be imported using their ID:

```shell
terraform import openstack_networking_bgp_speaker_v2.speaker_1 8a2ad402-b805-46bf-a60b-008573ca2844
```
