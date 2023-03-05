---
layout: "openstack"
page_title: "OpenStack: openstack_networking_network_v2"
sidebar_current: "docs-openstack-datasource-networking-network-v2"
description: |-
  Get information on an OpenStack Network.
---

# openstack\_networking\_network\_v2

Use this data source to get the ID of an available OpenStack network.

## Example Usage

```hcl
data "openstack_networking_network_v2" "network" {
  name = "tf_test_network"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve networks ids. If omitted, the
  `region` argument of the provider is used.

* `network_id` - (Optional) The ID of the network.

* `name` - (Optional) The name of the network.

* `description` - (Optional) Human-readable description of the network.

* `status` - (Optional) The status of the network.

* `external` - (Optional) The external routing facility of the network.

* `matching_subnet_cidr` - (Optional) The CIDR of a subnet within the network.

* `tenant_id` - (Optional) The owner of the network.

* `availability_zone_hints` - (Optional) The availability zone candidates for the network.

* `transparent_vlan` - (Optional) The VLAN transparent attribute for the
  network.

* `tags` - (Optional) The list of network tags to filter.

* `mtu` - (Optional) The network MTU to filter. Available, when Neutron `net-mtu`
  extension is enabled.

## Attributes Reference

`id` is set to the ID of the found network. In addition, the following attributes
are exported:

* `admin_state_up` - The administrative state of the network.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `region` - See Argument Reference above.
* `external` - See Argument Reference above.
* `shared` - Specifies whether the network resource can be accessed by any
   tenant or not.
* `availability_zone_hints` - The availability zone candidates for the network.
* `transparent_vlan` - See Argument Reference above.
* `segments` - An array of one or more provider segment objects.
* `mtu` - See Argument Reference above.
* `dns_domain` - The network DNS domain. Available, when Neutron DNS extension
  is enabled
* `subnets` - A list of subnet IDs belonging to the network.
* `all_tags` - The set of string tags applied on the network.
