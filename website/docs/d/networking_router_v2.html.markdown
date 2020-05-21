---
layout: "openstack"
page_title: "OpenStack: openstack_networking_router_v2"
sidebar_current: "docs-openstack-datasource-networking-router-v2"
description: |-
  Get information on an OpenStack Floating IP.
---

# openstack\_networking\_router\_v2

Use this data source to get the ID of an available OpenStack router.

## Example Usage

```hcl
data "openstack_networking_router_v2" "router" {
  name = "router_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve router ids. If omitted, the
  `region` argument of the provider is used.

* `router_id` - (Optional) The UUID of the router resource.

* `name` - (Optional) The name of the router.

* `description` - (Optional) Human-readable description of the router.

* `admin_state_up` - (Optional) Administrative up/down status for the router (must be "true" or "false" if provided).

* `distributed` - (Optional) Indicates whether or not to get a distributed router.

* `status` - (Optional) The status of the router (ACTIVE/DOWN).

* `tags` - (Optional) The list of router tags to filter.

* `tenant_id` - (Optional) The owner of the router.

## Attributes Reference

`id` is set to the ID of the found router. In addition, the following attributes
are exported:

* `enable_snat` - The value that points out if the Source NAT is enabled on the router.

* `external_network_id` - The network UUID of an external gateway for the router.

* `availability_zone_hints` - The availability zone that is used to make router resources highly available.

* `external_fixed_ip` - The external fixed IPs of the router.

The `external_fixed_ip` block supports:

* `subnet_id`- Subnet in which the fixed IP belongs to.

* `ip_address` - The IP address to set on the router.

* `all_tags` - The set of string tags applied on the router.
