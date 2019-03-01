---
layout: "openstack"
page_title: "OpenStack: openstack_networking_port_v2"
sidebar_current: "docs-openstack-datasource-networking-port-v2"
description: |-
  Get information of an OpenStack Port.
---

# openstack\_networking\_port\_v2

Use this data source to get the ID of an available OpenStack port.

## Example Usage

```hcl
data "openstack_networking_port_v2" "port_1" {
  name = "port_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve port ids. If omitted, the
  `region` argument of the provider is used.

* `project_id` - (Optional) The owner of the port.

* `port_id` - (Optional) The ID of the port.

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

* `dns_name` - (Optional) The port DNS name to filter. Available, when Neutron
    DNS extension is enabled.

## Attributes Reference

`id` is set to the ID of the found port. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.

* `project_id` - See Argument Reference above.

* `port_id` - See Argument Reference above.

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `admin_state_up` - See Argument Reference above.

* `network_id` - See Argument Reference above.

* `device_owner` - See Argument Reference above.

* `mac_address` - See Argument Reference above.

* `device_id` - See Argument Reference above.

* `allowed_address_pairs` - An IP/MAC Address pair of additional IP
    addresses that can be active on this port. The structure is described
    below.

* `all_fixed_ips` - The collection of Fixed IP addresses on the port in the
  order returned by the Network v2 API.

* `all_security_group_ids` - The set of security group IDs applied on the port.

* `all_tags` - The set of string tags applied on the port.

* `extra_dhcp_option` - An extra DHCP option configured on the port.
    The structure is described below.

* `binding` - The port binding information. The structure is described below.

* `dns_name` - See Argument Reference above.

* `dns_assignment` - The list of maps representing port DNS assignments.

The `allowed_address_pairs` attribute has fields below:

* `ip_address` - The additional IP address.

* `mac_address` - The additional MAC address.

The `extra_dhcp_option` attribute has fields below:

* `name` - Name of the DHCP option.

* `value` - Value of the DHCP option.

* `ip_version` - IP protocol version

The `binding` attribute has fields below:

* `host_id` - The ID of the host, which has the allocatee port.

* `profile` - A JSON string containing the binding profile information.

* `vnic_type` - VNIC type for the port.

* `vif_details` - A map of JSON strings containing additional details for this
    specific binding.

* `vif_type` - The VNIC type of the port binding.
