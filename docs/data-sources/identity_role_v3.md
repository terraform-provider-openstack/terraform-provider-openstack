---
subcategory: "Identity / Keystone"
layout: "openstack"
page_title: "OpenStack: openstack_identity_role_v3"
sidebar_current: "docs-openstack-datasource-identity-role-v3"
description: |-
  Get information on an OpenStack Role.
---

# openstack\_identity\_role\_v3

Use this data source to get the ID of an OpenStack role.

## Example Usage

```hcl
data "openstack_identity_role_v3" "admin" {
  name = "admin"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
  If omitted, the `region` argument of the provider is used.

* `name` - (Required) The name of the role.

* `domain_id` - (Optional) The domain the role belongs to.

## Attributes Reference

`id` is set to the ID of the found role. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `domain_id` - See Argument Reference above.
