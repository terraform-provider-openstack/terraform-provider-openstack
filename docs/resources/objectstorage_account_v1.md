---
subcategory: "Object Storage / Swift"
layout: "openstack"
page_title: "OpenStack: openstack_objectstorage_account_v1"
sidebar_current: "docs-openstack-resource-objectstorage-account-v1"
description: |-
  Manages a V1 account resource within OpenStack.
---

# openstack\_objectstorage\_account\_v1

Manages a V1 account resource within OpenStack.

## Example Usage

```hcl
resource "openstack_objectstorage_account_v1" "account_1" {
  region = "RegionOne"

  metadata = {
    Temp-Url-Key = "testkey"
    test         = "true"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the account. If omitted,
  the `region` argument of the provider is used. Changing this creates a new
  account.

* `project_id` - (Optional) The project ID of the corresponding account. If
  omitted, the token's project ID is used. Changing this creates a new account.

* `metadata` - (Optional) A map of custom key/value pairs to associate with the
  account metadata. Changing the `Quota-Bytes` key value is allowed to be
  updated only by the cloud administrator.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `metadata` - See Argument Reference above.
* `headers` - A map of headers returned for the account.
* `bytes_used` - The number of bytes used by the account.
* `quota_bytes` - The number of bytes allowed for the account.
* `container_count` - The number of containers in the account.
* `object_count` - The number of objects in the account.

## Import

This resource can be imported by specifying the project ID of the account:

```shell
terraform import openstack_objectstorage_account_v1.account_1 1202b3d0aaa44cfc8b79475c007b0711
```
