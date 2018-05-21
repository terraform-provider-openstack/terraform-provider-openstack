---
layout: "openstack"
page_title: "OpenStack: openstack_identity_auth_scope_v3"
sidebar_current: "docs-openstack-datasource-identity-auth-scope-v3"
description: |-
  Get authentication information from the current authenticated scope.
---

# openstack\_identity\_auth\_scope\_v3

Use this data source to get authentication information about the current
auth scope in use. This can be used as self-discovery or introspection of
the username or project name currently in use.

## Example Usage

```hcl
data "openstack_identity_auth_scope_v3" "scope" {
  name = "my_scope"
}
```

## Argument Reference

* `name` - (Required) The name of the scope. This is an arbitrary name which is
  only used as a unique identifier so an actual token isn't used as the ID.

* `region` - (Optional) The region in which to obtain the V3 Identity client.
  A Identity client is needed to retrieve tokens IDs. If omitted, the
  `region` argument of the provider is used.

## Attributes Reference

`id` is set to the name given to the scope. In addition, the following attributes
are exported:

* `user_name` - The username of the scope.
* `user_id` - The user ID the of the scope.
* `user_domain_name` - The domain name of the user.
* `user_domain_id` - The domain ID of the user.
* `project_name` - The project name of the scope.
* `project_id` - The project ID of the scope.
* `project_domain_name` - The domain name of the project.
* `project_domain_id` - The domain ID of the project.
* `roles` - A list of roles in the current scope. See reference below.

The `roles` block contains:

* `role_id` - The ID of the role.
* `role_name` - The name of the role.
