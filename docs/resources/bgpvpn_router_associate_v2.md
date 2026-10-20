---
subcategory: "BGP VPN / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_bgpvpn_router_associate_v2"
sidebar_current: "docs-openstack-resource-bgpvpn-router-associate-v2"
description: |-
  Manages a V2 BGP VPN router association resource within OpenStack.
---

# openstack\_bgpvpn\_router\_associate\_v2

Manages a V2 BGP VPN router association resource within OpenStack.

## Example Usage

```hcl
resource "openstack_bgpvpn_router_associate_v2" "association_1" {
  bgpvpn_id = "d57d39e1-dc63-44fd-8cbd-a4e1488100c5"
  router_id = "423fa80f-e0d7-4d02-a9a5-8b8c05812bf6"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a BGP VPN router association. If
  omitted, the `region` argument of the provider is used. Changing this creates
  a new BGP VPN router association.

* `bgpvpn_id` - (Required) The ID of the BGP VPN to which the router will be
  associated. Changing this creates a new BGP VPN router association.

* `router_id` - (Required) The ID of the router to be associated with the BGP
  VPN. Changing this creates a new BGP VPN router association.

* `project_id` - (Optional) The ID of the project that owns the BGP VPN router
  association. Only administrative and users with `advsvc` role can specify a
  project ID other than their own. Changing this creates a new BGP VPN router
  association.

* `advertise_extra_routes` - (Optional) A boolean flag indicating whether extra
  routes should be advertised. Defaults to true.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the BGP VPN router association.
* `region` - See Argument Reference above.
* `bgpvpn_id` - See Argument Reference above.
* `router_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `advertise_extra_routes` - See Argument Reference above.

## Import

BGP VPN router associations can be imported using the BGP VPN ID and BGP VPN
router association ID separated by a slash, e.g.:

```shell
terraform import openstack_bgpvpn_router_associate_v2.association_1 e26d509e-fc2d-4fb5-8562-619911a9a6bc/3cc9df2d-80db-4536-8ba6-295d1d0f723f
```
