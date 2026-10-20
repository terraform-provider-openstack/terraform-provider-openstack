---
subcategory: "BGP VPN / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_bgpvpn_network_associate_v2"
sidebar_current: "docs-openstack-resource-bgpvpn-network-associate-v2"
description: |-
  Manages a V2 BGP VPN network association resource within OpenStack.
---

# openstack\_bgpvpn\_network\_associate\_v2

Manages a V2 BGP VPN network association resource within OpenStack.

## Example Usage

```hcl
resource "openstack_bgpvpn_network_associate_v2" "association_1" {
  bgpvpn_id  = "e7189337-5684-46ee-bcb1-44f1a57066c9"
  network_id = "de83d56c-4d2f-44f7-ac24-af393252204f"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a BGP VPN network association. If
  omitted, the `region` argument of the provider is used. Changing this creates
  a new BGP VPN network association.

* `bgpvpn_id` - (Required) The ID of the BGP VPN to which the network will be
  associated. Changing this creates a new BGP VPN network association

* `network_id` - (Required) The ID of the network to be associated with the BGP
  VPN. Changing this creates a new BGP VPN network association.

* `project_id` - (Optional) The ID of the project that owns the BGP VPN network
  association. Only administrative and users with `advsvc` role can specify a
  project ID other than their own. Changing this creates a new BGP VPN network
  association.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the BGP VPN network association.
* `region` - See Argument Reference above.
* `bgpvpn_id` - See Argument Reference above.
* `network_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.

## Import

BGP VPN network associations can be imported using the BGP VPN ID and BGP VPN
network association ID separated by a slash, e.g.:

```shell
terraform import openstack_bgpvpn_network_associate_v2.association_1 2145aaa9-edaa-44fb-9815-e47a96677a72/67bb952a-f9d1-4fc8-ae84-082253a879d4
```
