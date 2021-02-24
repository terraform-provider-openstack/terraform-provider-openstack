---
layout: "openstack"
page_title: "OpenStack: openstack_networking_subnet_ids_v2"
sidebar_current: "docs-openstack-datasource-networking-subnet-ids-v2"
description: |-
  Provides a list of Openstack Subnet IDs.
---

# openstack\_networking\_subnet\_ids\_v2

Use this data source to get a list of Openstack Subnet IDs matching the
specified criteria.

## Example Usage

```hcl
data "openstack_networking_subnet_ids_v2" "subnets" {
  name_regex = "public"
  tags = [
    "public"
  ]
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve subnet ids. If omitted, the
  `region` argument of the provider is used.

* `name` - (Optional) The name of the subnet.

* `description` - (Optional) Human-readable description of the subnet.

* `dhcp_enabled` - (Optional) If the subnet has DHCP enabled.

* `network_id` - (Optional) The ID of the network the subnet belongs to.

* `tenant_id` - (Optional) The owner of the subnet.

* `ip_version` - (Optional) The IP version of the subnet (either 4 or 6).

* `gateway_ip` - (Optional) The IP of the subnet's gateway.

* `cidr` - (Optional) The CIDR of the subnet.

* `ipv6_address_mode` - (Optional) The IPv6 address mode. Valid values are
  `dhcpv6-stateful`, `dhcpv6-stateless`, or `slaac`.

* `ipv6_ra_mode` - (Optional) The IPv6 Router Advertisement mode. Valid values
  are `dhcpv6-stateful`, `dhcpv6-stateless`, or `slaac`.

* `subnetpool_id` - (Optional) The ID of the subnetpool associated with the subnet.

* `tags` - (Optional) The list of subnet tags to filter.

* `sort_key` - (Optional) Sort subnets based on a certain key. Defaults to none.

* `sort_direction` - (Optional) Order the results in either `asc` or `desc`.
  Defaults to none.

## Attributes Reference

`ids` is set to the list of Openstack Subnet IDs.
