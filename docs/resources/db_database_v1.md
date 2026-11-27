---
subcategory: "Databases / Trove"
layout: "openstack"
page_title: "OpenStack: openstack_db_database_v1"
sidebar_current: "docs-openstack-resource-db-database-v1"
description: |-
  Manages a V1 database resource within OpenStack.
---

# openstack\_db\_database\_v1

Manages a V1 DB database resource within OpenStack.

## Example Usage

### Database

```hcl
resource "openstack_db_database_v1" "mydb" {
  name        = "mydb"
  instance_id = openstack_db_instance_v1.basic.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the database. Changing
  this creates a new database.

* `name` - (Required) A unique name for the resource.

* `instance_id` - (Required) The ID for the database instance.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `instance_id` - See Argument Reference above.

## Import

Databases can be imported by using `instance-id/db-name`, e.g.

```shell
terraform import openstack_db_database_v1.mydb 7b9e3cd3-00d9-449c-b074-8439f8e274fa/mydb
```
