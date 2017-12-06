---
layout: "openstack"
page_title: "OpenStack: openstack_db_database_v1"
sidebar_current: "docs-openstack-resource-db-database-v1"
description: |-
  Manages a V1 database resource within OpenStack.
---

# openstack\_db\_database_v1

Manages a V1 DB database resource within OpenStack.

## Example Usage

### Database

```hcl
resource "openstack_db_database_v1" "test" {
  name     = "testdb"
  instance = "${openstack_db_instance_v1.basic.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource.

* `instance` - (Required) The ID for the database instance.

## Attributes Reference

The following attributes are exported:

* `region` - Openstack region resource is created in.
* `name` - See Argument Reference above.
* `instance` - See Argument Reference above.
