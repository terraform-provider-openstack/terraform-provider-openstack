---
layout: "openstack"
page_title: "OpenStack: openstack_identity_group_v3"
sidebar_current: "docs-openstack-resource-identity-group-v3"
description: |-
  Manages a V3 group resource within OpenStack Keystone.
---

# openstack\_identity\_group\_v3

Manages a V3 group resource within OpenStack Keystone.

~> **Note:** You _must_ have admin privileges in your OpenStack cloud to use
this resource.

## Example Usage

```hcl
resource "openstack_identity_group_v3" "group_1" {
  name        = "group_1"
  description = "group 1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the group.

* `description` - (Optional) A description of the group.

* `domain_id` - (Optional) The domain the group belongs to.

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new group.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `domain_id` - See Argument Reference above.
* `region` - See Argument Reference above.

## Import

groups can be imported using the `id`, e.g.

```
$ terraform import openstack_identity_group_v3.group_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
