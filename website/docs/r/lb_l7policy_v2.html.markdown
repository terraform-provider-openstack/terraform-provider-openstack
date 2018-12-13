---
layout: "openstack"
page_title: "OpenStack: openstack_lb_l7policy_v2"
sidebar_current: "docs-openstack-resource-lb-l7policy-v2"
description: |-
  Manages a V2 l7policy resource within OpenStack.
---

# openstack\_lb\_l7policy\_v2

Manages a V2 l7policy resource within OpenStack.

## Example Usage

```hcl
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_l7policy_v2" "l7policy_1" {
  name             = "test"
  action           = "REDIRECT_TO_POOL"
  description      = "test l7policy"
  position         = 1
  listener_id      = "${openstack_lb_listener_v2.listener_1.id}"
  redirect_pool_id = "${openstack_lb_pool_v2.pool_1.id}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create an . If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    L7policy.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
    the L7policy.  Only administrative users can specify a tenant UUID
    other than their own. Changing this creates a new L7policy.

* `name` - (Optional) Human-readable name for the L7policy. Does not have
    to be unique.

* `description` - (Optional) Human-readable description for the L7policy.

* `action` - (Required) The L7policy action - can either be REDIRECT\_TO\_POOL,
    REDIRECT\_TO\_URL or REJECT.

* `listener_id` - (Required) The Listener on which the L7policy will be associated with.
    Changing this creates a new L7policy.

* `position` - (Optional) The position of this policy on the listener. Positions start at 1.

* `redirect_pool_id` - (Optional) Requests matching this policy will be redirected to the
    pool with this ID. Only valid if action is REDIRECT\_TO\_POOL.

* `redirect_url` - (Optional) Requests matching this policy will be redirected to this URL.
    Only valid if action is REDIRECT\_TO\_URL.

* `admin_state_up` - (Optional) The administrative state of the L7policy.
    A valid value is true (UP) or false (DOWN).

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the L7policy.
* `region` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `action` - See Argument Reference above.
* `listener_id` - See Argument Reference above.
* `position` - See Argument Reference above.
* `redirect_pool_id` - See Argument Reference above.
* `redirect_url` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
