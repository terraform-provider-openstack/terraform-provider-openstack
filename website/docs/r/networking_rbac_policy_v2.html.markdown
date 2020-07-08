---
layout: "openstack"
page_title: "OpenStack: openstack_networking_rbac_policy_v2"
sidebar_current: "docs-openstack-resource-networking-rbac-policy-v2"
description: |-
  Creates an RBAC policy for an OpenStack V2 resource.
---

# openstack\_networking\_rbac\_policy\_v2

The RBAC policy resource contains functionality for working with Neutron RBAC
Policies. Role-Based Access Control (RBAC) policy framework enables both
operators and users to grant access to resources for specific projects.

Sharing an object with a specific project is accomplished by creating a
policy entry that permits the target project the `access_as_shared` action
on that object.

To make a network available as an external network for specific projects
rather than all projects, use the `access_as_external` action.
If a network is marked as external during creation, it now implicitly creates
a wildcard RBAC policy granting everyone access to preserve previous behavior
before this feature was added.

## Example Usage

```hcl
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_rbac_policy_v2" "rbac_policy_1" {
  action        = "access_as_shared"
  object_id     = "${openstack_networking_network_v2.network_1.id}"
  object_type   = "network"
  target_tenant = "20415a973c9e45d3917f078950644697"
}

resource "openstack_networking_rbac_policy_v2" "rbac_policy_1" {
  action        = "access_as_shared"
  object_search {
    name_glob   = "netowrk*"
  }
  object_type   = "network"
  target_tenant = "31526ba84daf56e4a280189a617557a8"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to configure a routing entry on a subnet. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    routing entry.

* `action` - (Required) Action for the RBAC policy. Can either be
  `access_as_external` or `access_as_shared`.

* `object_id` - (Optional) The ID of the `object_type` resource. An
  `object_type` of `network` returns a network ID and an `object_type` of
   `qos_policy` returns a QoS ID. One of `object_id` or `object_search`
   must be set.

* `object_search` - (Optional) A search criteria to find for the object to
  share. The searched for object is of the `object_type` resource. An
  `object_type` of `network` returns a network ID and an `object_type` of
   `qos_policy` returns a QoS ID. One of `object_id` or `object_search`
   must be set. This filed includes the following arguments:
  * `name_glob` - (Required) The glob expression that will be used to
    search for objects. The search will return all the `object_type`
    resources whose name matches the glob, ordered lèhabetically. The
    first one returned will be shared.
  * `owning_tenant_id` - (Optional) The search will be limited to resources
    owned by this tenant.
  * `unshared` - (Optional) The search will be limited to resources
    that are not shared with any project, i.e. that have no associated RBAC.
    The default value is `true`.

* `object_type` - (Required) The type of the object that the RBAC policy
  affects. Can either be `qos-policy` or `network`.

* `target_tenant` - (Required) The ID of the tenant to which the RBAC policy
  will be enforced.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `action` - See Argument Reference above.
* `object_id` - See Argument Reference above.
* `object_type` - See Argument Reference above.
* `target_tenant` - See Argument Reference above.
* `tenant_id` - The owner of the RBAC policy.

## Notes

## Import

RBAC policies can be imported using the `id`, e.g.

```
$ terraform import openstack_networking_rbac_policy_v2.rbac_policy_1 eae26a3e-1c33-4cc1-9c31-0cd729c438a1
```
