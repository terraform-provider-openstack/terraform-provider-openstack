---
layout: "openstack"
page_title: "OpenStack: openstack_db_configuration_v1"
sidebar_current: "docs-openstack-resource-db-configuration-v1"
description: |-
  Manages a V1 DB configuration resource within OpenStack.
---

# openstack\_db\_configuration\_v1

Manages a V1 DB configuration resource within OpenStack.

## Example Usage

### Configuration

```hcl
resource "openstack_db_configuration_v1" "test" {
  name        = "test"
  description = "description"

  datastore {
    version = "mysql-5.7"
    type    = "mysql"
  }

  configuration {
    name  = "max_connections"
    value = 200
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region in which to create the db instance. Changing this
    creates a new instance.

* `name` - (Required) A unique name for the resource.

* `description` - (Optional) Description of the resource.

* `datastore` - (Required) An array of database engine type and version. The datastore
    object structure is documented below. Changing this creates resource.

* `configuration` - (Optional) An array of configuration parameter name and value. Can be specified multiple times. The configuration object structure is documented below.

The `datastore` block supports:

* `type` - (Required) Database engine type to be used with this configuration. Changing this creates a new resource.
* `version` - (Required) Version of database engine type to be used with this configuration. Changing this creates a new resource.

The `configuration` block supports:

* `name` - (Optional) Configuration parameter name. Changing this creates a new resource.
* `value` - (Optional) Configuration parameter value. Changing this creates a new resource.
* `string_type` - (Optional) Whether or not to store configuration parameter value as string. Changing this creates a new resource. See the below note for more information.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `datastore/type` - See Argument Reference above.
* `datastore/version` - See Argument Reference above.
* `configuration/name` - See Argument Reference above.
* `configuration/value` - See Argument Reference above.
* `configuration/string_type` - See Argument Reference above.

## Types of configuration parameter values

Openstack API requires to store some database configuration parameter's values as strings, even if they contain numbers.
To force store their values as strings set `string_type` to `true`. Otherwise Terraform will try to store them as number what can cause error from Openstack API like below:
```
"The value provided for the configuration parameter log_min_duration_statement is not of type string."
```