---
layout: "openstack"
page_title: "OpenStack: openstack_networking_portforwarding_v2"
sidebar_current: "docs-openstack-resource-networking-portforwarding-v2"
description: |-
  Manages a V2 port forwarding resource within OpenStack.
---

# openstack\_networking\_portforwarding_v2

Manages a V2 portforwarding resource within OpenStack.

## Example Usage

### Simple portforwarding

```hcl
"port_forwarding": {
         "external_port": "7233",
         "internal_port": "22",
         "internal_port_id": "b930d7f6-ceb7-40a0-8b81-a425dd994ccf"
     }
```


## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to create a port. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    port.

* `internal_port_id` - The ID of the Neutron port associated to the floating IP port forwarding.

* `internal_ip_address` - The fixed IPv4 address of the Neutron port associated to the floating IP port forwarding.

* `internal_port` - The TCP/UDP/other protocol port number of the Neutron port fixed IP address associated to the floating ip port forwarding.

* `external_port` - The TCP/UDP/other protocol port number of the port forwarding’s floating IP address.

* `protocol` - The IP protocol used in the floating IP port forwarding.

* `tenant_id` - (Optional) The owner of the Port. Required if admin wants
    to create a port for another tenant. Changing this creates a new port.

* `description` - (Optional) A text describing the rule, which helps users to manage/find easily theirs rules.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `id` - The ID of the floating IP port forwarding.
* `internal_port_id` - See Argument Reference above.
* `internal_ip_address` - See Argument Reference above.
* `internal_port` - See Argument Reference above.
* `external_port` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `description` - See Argument Reference above.



## Notes
