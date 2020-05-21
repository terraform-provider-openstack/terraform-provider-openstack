---
layout: "openstack"
page_title: "OpenStack: openstack_identity_user_v3"
sidebar_current: "docs-openstack-datasource-identity-user-v3"
description: |-
  Get information on an OpenStack User.
---

# openstack\_identity\_user\_v3

Use this data source to get the ID of an OpenStack user.

## Example Usage

```hcl
data "openstack_identity_user_v3" "user_1" {
  name = "user_1"
}
```

## Argument Reference

The following arguments are supported:

* `default_project_id` - (Optional) The default project this user belongs to.

* `domain_id` - (Optional) The domain this user belongs to.

* `enabled` - (Optional) Whether the user is enabled or disabled. Valid
  values are `true` and `false`.

* `idp_id` - (Optional) The identity provider ID of the user.

* `name` - (Optional) The name of the user.

* `password_expires_at` - (Optional) Query for expired passwords. See the [OpenStack API docs](https://developer.openstack.org/api-ref/identity/v3/#list-users) for more information on the query format.

* `protocol_id` - (Optional) The protocol ID of the user.

* `unique_id` - (Optional) The unique ID of the user.

## Attributes Reference

The following attributes are exported:

* `default_project_id` - See Argument Reference above.
* `domain_id` - See Argument Reference above.
* `enabled` - See Argument Reference above.
* `idp_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `password_expires_at` - See Argument Reference above.
* `protocol_id` - See Argument Reference above.
* `region` - The region the user is located in.
* `unique_id` - See Argument Reference above.
* `description` - A description of the user.
