---
layout: "openstack"
page_title: "OpenStack: openstack_networking_subnet_v2"
sidebar_current: "docs-openstack-datasource-networking-subnet-v2"
description: |-
  Get information on an OpenStack Subnet.
---

# openstack\_networking\_subnet\_v2

Use this data source to get the ID of an available OpenStack subnet.

## Example Usage

```hcl
data "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve subnet ids. If omitted, the
  `region` argument of the provider is used.

* `name` - (Optional) The name of the subnet.

* `description` - (Optional) Human-readable description for the subnet.

* `dhcp_enabled` - (Optional) If the subnet has DHCP enabled.

* `dhcp_disabled` - (Optional) If the subnet has DHCP disabled.

* `ip_version` - (Optional) The IP version of the subnet (either 4 or 6).

* `ipv6_address_mode` - (Optional) The IPv6 address mode. Valid values are
  `dhcpv6-stateful`, `dhcpv6-stateless`, or `slaac`.

* `ipv6_ra_mode` - (Optional) The IPv6 Router Advertisement mode. Valid values
  are `dhcpv6-stateful`, `dhcpv6-stateless`, or `slaac`.

* `gateway_ip` - (Optional) The IP of the subnet's gateway.

* `cidr` - (Optional) The CIDR of the subnet.

* `subnet_id` - (Optional) The ID of the subnet.

* `subnetpool_id` - (Optional) The ID of the subnetpool associated with the subnet.

* `network_id` - (Optional) The ID of the network the subnet belongs to.

* `tenant_id` - (Optional) The owner of the subnet.

* `tags` - (Optional) The list of subnet tags to filter.

## Attributes Reference

`id` is set to the ID of the found subnet. In addition, the following attributes
are exported:

* `allocation_pools` - Allocation pools of the subnet.
* `enable_dhcp` - Whether the subnet has DHCP enabled or not.
* `dns_nameservers` - DNS Nameservers of the subnet.
* `host_routes` - Host Routes of the subnet.
* `region` - See Argument Reference above.
* `all_tags` - A set of string tags applied on the subnet.
