---
subcategory: "BGP VPN / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_bgpvpn_port_associate_v2"
sidebar_current: "docs-openstack-resource-bgpvpn-port-associate-v2"
description: |-
  Manages a V2 BGP VPN port association resource within OpenStack.
---

# openstack\_bgpvpn\_port\_associate\_v2

Manages a V2 BGP VPN port association resource within OpenStack.

## Example Usage

```hcl
resource "openstack_bgpvpn_port_associate_v2" "association_1" {
  bgpvpn_id = "19382ec5-8098-47d9-a9c6-6270c91103f4"
  port_id   = "b83a95b8-c2c8-4eac-9a9e-ddc85bd1266f"

  routes {
    type   = "prefix"
    prefix = "192.168.170.1/32"
  }
  routes {
    type      = "bgpvpn"
    bgpvpn_id = "35af1cc6-3d0f-4c5d-86f8-8cdb508d3f0c"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a BGP VPN port association. If
  omitted, the `region` argument of the provider is used. Changing this creates
  a new BGP VPN port association.

* `bgpvpn_id` - (Required) The ID of the BGP VPN to which the port will be
  associated. Changing this creates a new BGP VPN port association.

* `port_id` - (Required) The ID of the port to be associated with the BGP VPN.
  Changing this creates a new BGP VPN port association.

* `project_id` - (Optional) The ID of the project that owns the port
  association. Only administrative and users with `advsvc` role can specify a
  project ID other than their own. Changing this creates a new BGP VPN port
  association.

* `advertise_fixed_ips` - (Optional) A boolean flag indicating whether fixed
  IPs should be advertised. Defaults to true.

* `routes` - (Optional) A list of dictionaries containing the following keys:
  * `type` - (Required) Can be `prefix` or `bgpvpn`. For the `prefix` type, the
    CIDR prefix (v4 or v6) must be specified in the `prefix` key. For the
    `bgpvpn` type, the BGP VPN ID must be specified in the `bgpvpn_id` key.
  * `prefix` - (Optional) The CIDR prefix (v4 or v6) to be advertised. Required
    if `type` is `prefix`. Conflicts with `bgpvpn_id`.
  * `bgpvpn_id` - (Optional) The ID of the BGP VPN to be advertised. Required
    if `type` is `bgpvpn`. Conflicts with `prefix`.
  * `local_pref` - (Optional) The BGP LOCAL\_PREF value of the routes that will
    be advertised.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the BGP VPN port association.
* `region` - See Argument Reference above.
* `bgpvpn_id` - See Argument Reference above.
* `port_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `advertise_fixed_ips` - See Argument Reference above.
* `routes` - See Argument Reference above.

## Import

BGP VPN port associations can be imported using the BGP VPN ID and BGP VPN port
association ID separated by a slash, e.g.:

```shell
terraform import openstack_bgpvpn_port_associate_v2.association_1 5bb44ecf-f8fe-4d75-8fc5-313f96ee2696/8f8fc660-3f28-414e-896a-0c7c51162fcf
```
