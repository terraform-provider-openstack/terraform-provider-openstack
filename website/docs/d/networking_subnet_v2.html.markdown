---
layout: "openstack"
page_title: "OpenStack: openstack_networking_subnet_v2"
sidebar_current: "docs-openstack-datasource-networking-subnet-v2"
description: |-
  Get information on an OpenStack Network Subnet.
---

# openstack\_networking\_subnet\_v2

Use this data source to get the ID of an available OpenStack network subnet.

## Example Usage

```hcl
data "openstack_networking_subnet_v2" "network" {
  name = "tf_test_subnet"
}
```

## Argument Reference

* `cidr` - (Optional) CIDR representing IP range for this subnet, based on IP
    version.

* `name` - (Optional) The name of the subnet.

* `network_id` - (Optional) The ID of the network. Can be useful when combining with "cidr" if
  two subnets in two different networks share the same CIDR block.

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve networks ids. If omitted, the
  `region` argument of the provider is used.

* `subnet_id` - (Optional) The ID of the subnet.

* `tenant_id` - (Optional) The owner of the subnet.

## Attributes Reference

`id` is set to the ID of the found network. In addition, the following attributes
are exported:

* `gateway_ip` - Default gateway used by devices in this subnet.
* `ip_version` - IP version, either 4 (default) or 6.
* `name` - See Argument Reference above.
* `network_id` - See Argument Reference above.
* `region` - See Argument Reference above.
* `shared` - (Optional)  Specifies whether the network resource can be accessed
    by any tenant or not.
