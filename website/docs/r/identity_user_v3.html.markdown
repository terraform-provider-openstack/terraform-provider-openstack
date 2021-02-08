---
layout: "openstack"
page_title: "OpenStack: openstack_identity_user_v3"
sidebar_current: "docs-openstack-resource-identity-user-v3"
description: |-
  Manages a V3 User resource within OpenStack Keystone.
---

# openstack\_identity\_user\_v3

Manages a V3 User resource within OpenStack Keystone.

~> **Note:** All arguments including the user password will be stored in the
raw state as plain-text. [Read more about sensitive data in
state](https://www.terraform.io/docs/language/state/sensitive-data.html).

~> **Note:** You _must_ have admin privileges in your OpenStack cloud to use
this resource.

## Example Usage

```hcl
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_identity_user_v3" "user_1" {
  default_project_id = "${openstack_identity_project_v3.project_1.id}"
  name               = "user_1"
  description        = "A user"

  password = "password123"

  ignore_change_password_upon_first_use = true

  multi_factor_auth_enabled = true

  multi_factor_auth_rule {
    rule = ["password", "totp"]
  }

  multi_factor_auth_rule {
    rule = ["password"]
  }

  extra = {
    email = "user_1@foobar.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description of the user.

* `default_project_id` - (Optional) The default project this user belongs to.

* `domain_id` - (Optional) The domain this user belongs to.

* `enabled` - (Optional) Whether the user is enabled or disabled. Valid
  values are `true` and `false`.

* `extra` - (Optional) Free-form key/value pairs of extra information.

* `ignore_change_password_upon_first_use` - (Optional) User will not have to
  change their password upon first use. Valid values are `true` and `false`.

* `ignore_password_expiry` - (Optional) User's password will not expire.
  Valid values are `true` and `false`.

* `ignore_lockout_failure_attempts` - (Optional) User will not have a failure
  lockout placed on their account. Valid values are `true` and `false`.

* `multi_factor_auth_enabled` - (Optional) Whether to enable multi-factor
  authentication. Valid values are `true` and `false`.

* `multi_factor_auth_rule` - (Optional) A multi-factor authentication rule.
  The structure is documented below. Please see the
  [Ocata release notes](https://docs.openstack.org/releasenotes/keystone/ocata.html)
  for more information on how to use mulit-factor rules.

* `name` - (Optional) The name of the user.

* `password` - (Optional) The password for the user.

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new User.

The `multi_factor_auth_rule` block supports:

* `rule` - (Required) A list of authentication plugins that the user must
  authenticate with.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.

## Import

Users can be imported using the `id`, e.g.

```
$ terraform import openstack_identity_user_v3.user_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
