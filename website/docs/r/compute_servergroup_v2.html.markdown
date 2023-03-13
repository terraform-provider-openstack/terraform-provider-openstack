---
layout: "openstack"
page_title: "OpenStack: openstack_compute_servergroup_v2"
sidebar_current: "docs-openstack-resource-compute-servergroup-v2"
description: |-
  Manages a V2 Server Group resource within OpenStack.
---

# openstack\_compute\_servergroup\_v2

Manages a V2 Server Group resource within OpenStack.

## Example Usage

### Compute service API version 2.63 or below:

```hcl
resource "openstack_compute_servergroup_v2" "test-sg" {
  name     = "my-sg"
  policies = ["anti-affinity"]
}
```

### Compute service API version 2.64 or above:

```hcl
resource "openstack_compute_servergroup_v2" "test-sg" {
  name     = "my-sg"
  policies = ["anti-affinity"]
  rules {
      max_server_per_host = 3
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
  If omitted, the `region` argument of the provider is used. Changing
  this creates a new server group.

* `name` - (Required) A unique name for the server group. Changing this creates
  a new server group.

* `policies` - (Optional) A list of exactly one policy name to associate with
  the server group. See the Policies section for more information. Changing this
  creates a new server group.

* `value_specs` - (Optional) Map of additional options.

* `rules` - (Optional) The rules which are applied to specified `policy`. Currently,
  only the `max_server_per_host` rule is supported for the `anti-affinity` policy.

## Policies

* `affinity` - All instances/servers launched in this group will be hosted on
    the same compute node.

* `anti-affinity` - All instances/servers launched in this group will be
    hosted on different compute nodes.

* `soft-affinity` - All instances/servers launched in this group will be hosted
    on the same compute node if possible, but if not possible they
    still will be scheduled instead of failure. To use this policy your
    OpenStack environment should support Compute service API 2.15 or above.

* `soft-anti-affinity` - All instances/servers launched in this group will be
    hosted on different compute nodes if possible, but if not possible they
    still will be scheduled instead of failure. To use this policy your
    OpenStack environment should support Compute service API 2.15 or above.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `policies` - See Argument Reference above.
* `members` - The instances that are part of this server group.
* `rules` - See Argument Reference above.

## Import

Server Groups can be imported using the `id`, e.g.

```
$ terraform import openstack_compute_servergroup_v2.test-sg 1bc30ee9-9d5b-4c30-bdd5-7f1e663f5edf
```
