---
layout: "openstack"
page_title: "OpenStack: openstack_identity_group_v3"
sidebar_current: "docs-openstack-datasource-identity-group-v3"
description: |-
  Get information on an OpenStack Group.
---

# openstack\_identity\_group\_v3

Use this data source to get the ID of an OpenStack group.

~> **Note:** You _must_ have admin privileges in your OpenStack cloud to use
this resource.

## Example Usage

```hcl
data "openstack_identity_group_v3" "admins" {
  name = "admins"
}
```

## Argument Reference

* `name` - The name of the group.

* `domain_id` - (Optional) The domain the group belongs to.

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
    If omitted, the `region` argument of the provider is used.


## Attributes Reference

`id` is set to the ID of the found group. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `domain_id` - See Argument Reference above.
* `region` - See Argument Reference above.
* `description` - A description of the group.
