---
layout: "openstack"
page_title: "OpenStack: openstack_identity_token_v3"
sidebar_current: "docs-openstack-datasource-identity-token-v3"
description: |-
  Get information about the current OpenStack token in use.
---

# openstack\_identity\_token\_v3

Use this data source to get authentication information about the current
token in use. This can be used as self-discovery or introspection of the
username or project name currently in use.

## Example Usage

```hcl
data "openstack_identity_token_v3" "token" {
  name = "my_token"
}
```

## Argument Reference

* `name` - (Required) The name of the token. This is an arbitrary name which is
  only used as a unique identifier so the actual token isn't used as the ID.

* `region` - (Optional) The region in which to obtain the V3 Identity client.
  A Identity client is needed to retrieve tokens ids. If omitted, the
  `region` argument of the provider is used.

## Attributes Reference

`id` is set to the name given to the token. In addition, the following attributes
are exported:

* `user_name` - The username the token is scoped to.
* `user_id` - The user ID the token is scoped to.
* `user_domain_name` - The domain name of the user.
* `user_domain_id` - The domain ID of the user.
* `project_name` - The project name the token is scoped to.
* `project_id` - The project ID the token is scoped to.
* `project_domain_name` - The domain name of the project.
* `project_domain_id` - The domain ID of the project.
* `roles` - A list of roles the token is scoped to. See reference below.

The `roles` block contains:

* `role_id` - The ID of the role.
* `role_name` - The name of the role.
