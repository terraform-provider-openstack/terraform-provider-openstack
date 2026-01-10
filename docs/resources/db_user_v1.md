---
subcategory: "Databases / Trove"
layout: "openstack"
page_title: "OpenStack: openstack_db_user_v1"
sidebar_current: "docs-openstack-resource-db-user-v1"
description: |-
  Manages a V1 database user resource within OpenStack.
---

# openstack\_db\_user\_v1

Manages a V1 DB user resource within OpenStack.

~> **Note:** All arguments including the database password will be stored in the
raw state as plain-text. [Read more about sensitive data in
state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage

### User

```hcl
resource "openstack_db_user_v1" "basic" {
  name        = "basic"
  instance_id = openstack_db_instance_v1.basic.id
  password    = "password"
  databases   = ["testdb"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the db user. Changing
  this creates a new user.

* `name` - (Required) A unique name for the resource.

* `instance_id` - (Required) The ID for the database instance.

* `password` - (Required) User's password.

* `databases` - (Optional) A list of database user should have access to.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `instance_id` - See Argument Reference above.
* `password` - See Argument Reference above.
* `databases` - See Argument Reference above.
