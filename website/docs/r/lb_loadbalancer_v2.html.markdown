---
layout: "openstack"
page_title: "OpenStack: openstack_lb_loadbalancer_v2"
sidebar_current: "docs-openstack-resource-lb-loadbalancer-v2"
description: |-
  Manages a V2 loadbalancer resource within OpenStack.
---

# openstack\_lb\_loadbalancer\_v2

Manages a V2 loadbalancer resource within OpenStack.

~> **Note:** This resource has attributes that depend on octavia minor versions.
Please ensure your Openstack cloud supports the required [minor version](../#octavia-api-versioning).

## Example Usage

```hcl
resource "openstack_lb_loadbalancer_v2" "lb_1" {
  vip_subnet_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create an LB member. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    LB member.

* `vip_subnet_id` - (Optional) The subnet on which to allocate the
    Loadbalancer's address. A tenant can only create Loadbalancers on networks
    authorized by policy (e.g. networks that belong to them or networks that
    are shared).  Changing this creates a new loadbalancer.
    It is required to Neutron LBaaS but optional for Octavia.

* `vip_network_id` - (Optional) The network on which to allocate the
    Loadbalancer's address. A tenant can only create Loadbalancers on networks
    authorized by policy (e.g. networks that belong to them or networks that
    are shared).  Changing this creates a new loadbalancer.
    It is available only for Octavia.

* `vip_port_id` - (Optional) The port UUID that the loadbalancer will use.
  Changing this creates a new loadbalancer. It is available only for Octavia.

* `name` - (Optional) Human-readable name for the Loadbalancer. Does not have
    to be unique.

* `description` - (Optional) Human-readable description for the Loadbalancer.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
    the Loadbalancer.  Only administrative users can specify a tenant UUID
    other than their own.  Changing this creates a new loadbalancer.

* `vip_address` - (Optional) The ip address of the load balancer.
    Changing this creates a new loadbalancer.

* `admin_state_up` - (Optional) The administrative state of the Loadbalancer.
    A valid value is true (UP) or false (DOWN).

* `flavor_id` - (Optional) The UUID of a flavor. Changing this creates a new
    loadbalancer.

* `loadbalancer_provider` - (Optional) The name of the provider. Changing this
  creates a new loadbalancer.

* `availability_zone` - (Optional) The availability zone of the Loadbalancer.
  Changing this creates a new loadbalancer. Available only for Octavia
  **minor version 2.14 or later**.

* `security_group_ids` - (Optional) A list of security group IDs to apply to the
    loadbalancer. The security groups must be specified by ID and not name (as
    opposed to how they are configured with the Compute Instance).

* `tags` - (Optional) A list of simple strings assigned to the loadbalancer.
    Available only for Octavia **minor version 2.5 or later**.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `vip_subnet_id` - See Argument Reference above.
* `vip_network_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `vip_address` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `flavor_id` - See Argument Reference above.
* `loadbalancer_provider` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `security_group_ids` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `vip_port_id` - The Port ID of the Load Balancer IP.

## Import

Load Balancer can be imported using the Load Balancer ID, e.g.:

```
$ terraform import openstack_lb_loadbalancer_v2.loadbalancer_1 19bcfdc7-c521-4a7e-9459-6750bd16df76
```
