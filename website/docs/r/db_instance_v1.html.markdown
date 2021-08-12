---
layout: "openstack"
page_title: "OpenStack: openstack_db_instance_v1"
sidebar_current: "docs-openstack-resource-db-instance-v1"
description: |-
  Manages a V1 DB instance resource within OpenStack.
---

# openstack\_db\_instance\_v1

Manages a V1 DB instance resource within OpenStack.

~> **Note:** All arguments including the instance user password will be stored
in the raw state as plain-text. [Read more about sensitive data in
state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage

### Instance

```hcl
resource "openstack_db_instance_v1" "test" {
  region    = "region-test"
  name      = "test"
  flavor_id = "31792d21-c355-4587-9290-56c1ed0ca376"
  size      = 8

  network {
    uuid = "c0612505-caf2-4fb0-b7cb-56a0240a2b12"
  }

  datastore {
    version = "mysql-5.7"
    type    = "mysql"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region in which to create the db instance. Changing this
    creates a new instance.

* `name` - (Required) A unique name for the resource.

* `flavor_id` - (Required) The flavor ID of the desired flavor for the instance.
    Changing this creates new instance.

* `configuration_id` - (Optional) Configuration ID to be attached to the instance. Database instance
   will be rebooted when configuration is detached.

* `size` - (Required) Specifies the volume size in GB. Changing this creates new instance.

* `datastore` - (Required) An array of database engine type and version. The datastore
    object structure is documented below. Changing this creates a new instance.

* `network` - (Optional) An array of one or more networks to attach to the
    instance. The network object structure is documented below. Changing this
    creates a new instance.

* `user` - (Optional) An array of username, password, host and databases. The user
    object structure is documented below.

* `database` - (Optional) An array of database name, charset and collate. The database
    object structure is documented below.

The `datastore` block supports:

* `type` - (Required) Database engine type to be used in new instance. Changing this
    creates a new instance.
* `version` - (Required) Version of database engine type to be used in new instance.
    Changing this creates a new instance.

The `network` block supports:

* `uuid` - (Required unless `port` is provided) The network UUID to
    attach to the instance. Changing this creates a new instance.

* `port` - (Required unless `uuid` is provided) The port UUID of a
    network to attach to the instance. Changing this creates a new instance.

* `fixed_ip_v4` - (Optional) Specifies a fixed IPv4 address to be used on this
    network. Changing this creates a new instance.

* `fixed_ip_v6` - (Optional) Specifies a fixed IPv6 address to be used on this
    network. Changing this creates a new instance.

The `user` block supports:

* `name` - (Optional) Username to be created on new instance. Changing this creates a
    new instance.

* `password` - (Optional) User's password. Changing this creates a
    new instance.

* `host` - (Optional) An ip address or % sign indicating what ip addresses can connect with
    this user credentials. Changing this creates a new instance.

* `databases` - (Optional) A list of databases that user will have access to. If not specified,
     user has access to all databases on th einstance. Changing this creates a new instance.

The `database` block supports:

* `name` - (Optional) Database to be created on new instance. Changing this creates a
    new instance.

* `collate` - (Optional) Database collation. Changing this creates a new instance.

* `charset` - (Optional) Database character set. Changing this creates a
    new instance.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `size` - See Argument Reference above.
* `flavor_id` - See Argument Reference above.
* `configuration_id` - See Argument Reference above.
* `datastore/type` - See Argument Reference above.
* `datastore/version` - See Argument Reference above.
* `network/uuid` - See Argument Reference above.
* `network/port` - See Argument Reference above.
* `network/fixed_ip_v4` - The Fixed IPv4 address of the Instance on that
    network.
* `network/fixed_ip_v6` - The Fixed IPv6 address of the Instance on that
* `database/name` - See Argument Reference above.
* `database/collate` - See Argument Reference above.
* `database/charset` - See Argument Reference above.
* `user/name` - See Argument Reference above.
* `user/password` - See Argument Reference above.
* `user/databases` - See Argument Reference above.
* `user/host` - See Argument Reference above.
* `addresses` - A list of IP addresses assigned to the instance.
