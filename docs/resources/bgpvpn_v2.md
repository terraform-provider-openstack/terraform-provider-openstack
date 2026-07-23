---
subcategory: "BGP VPN / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_bgpvpn_v2"
sidebar_current: "docs-openstack-resource-bgpvpn-v2"
description: |-
  Manages a V2 BGP VPN service resource within OpenStack.
---

# openstack\_bgpvpn\_v2

Manages a V2 BGP VPN service resource within OpenStack.

## Example Usage

```hcl
resource "openstack_bgpvpn_v2" "bgpvpn_1" {
  name                 = "bgpvpn1"
  route_distinguishers = ["64512:1"]
  route_targets        = ["64512:1"]
  import_targets       = ["64512:2"]
  export_targets       = ["64512:3"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a BGP VPN service. If omitted, the
  `region` argument of the provider is used. Changing this creates a new
  BGP VPN.

* `name` - (Optional) The name of the BGP VPN. Changing this updates the name of
  the existing BGP VPN.

* `type` - (Optional) The type of the BGP VPN (either `l2` or `l3`). Changing this
  creates a new BGP VPN. Defaults to `l3`.

* `project_id` - (Optional) The ID of the project that owns the BGPVPN. Only
  administrative and users with `advsvc` role can specify a project ID other
  than their own. Changing this creates a new BGP VPN.

* `vni` - (Optional) The globally-assigned VXLAN VNI for the BGP VPN. Changing
  this creates a new BGP VPN.

* `local_pref` - (Optional) The default BGP LOCAL\_PREF of routes that will be
  advertised to the BGP VPN, unless overridden per-route.

* `route_distinguishers` - (Optional) A list of route distinguisher strings. If
 specified, one of these RDs will be used to advertise VPN routes.

* `route_targets` - (Optional) A list of Route Targets that will be both
  imported and used for export.

* `import_targets` - (Optional) A list of additional Route Targets that will be
  imported.

* `export_targets` - (Optional) A list of additional Route Targets that will be
  used for export.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the BGP VPN.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `type` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `vni` - See Argument Reference above.
* `local_pref` - See Argument Reference above.
* `route_distinguishers` - See Argument Reference above.
* `route_targets` - See Argument Reference above.
* `import_targets` - See Argument Reference above.
* `export_targets` - See Argument Reference above.
* `networks` - A list of network IDs that are associated with the BGP VPN.
* `routers` - A list of router IDs that are associated with the BGP VPN.
* `ports` - A list of port IDs that are associated with the BGP VPN.
* `shared` - Indicates whether the BGP VPN is shared across projects.

## Import

BGP VPNs can be imported using the `id`, e.g.

```shell
terraform import openstack_bgpvpn_v2.bgpvpn_1 1eec2c66-6be2-4305-af3f-354c9b81f18c
```
