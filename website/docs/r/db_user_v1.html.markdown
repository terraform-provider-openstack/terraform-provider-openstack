---
layout: "openstack"
page_title: "OpenStack: openstack_db_user_v1"
sidebar_current: "docs-openstack-resource-db-user-v1"
description: |-
  Manages a V1 database user resource within OpenStack.
---

# openstack\_db\_user_v1

Manages a V1 DB user resource within OpenStack.

## Example Usage

### User

```hcl
resource "openstack_db_user_v1" "basic" {
  name      = "basic"
  instance  = "${openstack_db_instance_v1.basic.id}"
  password  = "password"
  databases = ["testdb"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource.

* `instance` - (Required) The ID for the database instance.

* `password` - (Required) User's password.

* `databases` - (Optional) A list of database user should have access to.

## Attributes Reference

The following attributes are exported:

* `region` - Openstack region resource is created in.
* `name` - See Argument Reference above.
* `instance` - See Argument Reference above.
* `password` - See Argument Reference above.
* `databases` - See Argument Reference above.