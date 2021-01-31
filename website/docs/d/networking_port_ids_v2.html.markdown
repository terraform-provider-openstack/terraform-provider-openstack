---
layout: "openstack"
page_title: "OpenStack: openstack_networking_port_ids_v2"
sidebar_current: "docs-openstack-datasource-networking-port-ids-v2"
description: |-
  Provides a list of Openstack Port IDs.
---

# openstack\_networking\_port\_ids\_v2

Use this data source to get a list of Openstack Port IDs matching the
specified criteria.

## Example Usage

```hcl
data "openstack_networking_port_ids_v2" "ports" {
  name = "port"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve port ids. If omitted, the
  `region` argument of the provider is used.

* `project_id` - (Optional) The owner of the port.

* `name` - (Optional) The name of the port.

* `description` - (Optional) Human-readable description of the port.

* `admin_state_up` - (Optional) The administrative state of the port.

* `network_id` - (Optional) The ID of the network the port belongs to.

* `device_owner` - (Optional) The device owner of the port.

* `mac_address` - (Optional) The MAC address of the port.

* `device_id` - (Optional) The ID of the device the port belongs to.

* `fixed_ip` - (Optional) The port IP address filter.

* `status` - (Optional) The status of the port.

* `security_group_ids` - (Optional) The list of port security group IDs to filter.

* `tags` - (Optional) The list of port tags to filter.

* `sort_key` - (Optional) Sort ports based on a certain key. Defaults to none.

* `sort_direction` - (Optional) Order the results in either `asc` or `desc`.
  Defaults to none.

## Attributes Reference

`ids` is set to the list of Openstack Port IDs.
