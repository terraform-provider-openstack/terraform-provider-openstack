---
layout: "openstack"
page_title: "OpenStack: openstack_objectstorage_container_v1"
sidebar_current: "docs-openstack-resource-objectstorage-container-v1"
description: |-
  Manages a V1 container resource within OpenStack.
---

# openstack\_objectstorage\_container\_v1

Manages a V1 container resource within OpenStack.

## Example Usage

### Basic Container

```hcl
resource "openstack_objectstorage_container_v1" "container_1" {
  region = "RegionOne"
  name   = "tf-test-container-1"

  metadata = {
    test = "true"
  }

  content_type = "application/json"
  versioning   = true
}
```

### Basic Container with legacy versioning

```hcl
resource "openstack_objectstorage_container_v1" "container_1" {
  region = "RegionOne"
  name   = "tf-test-container-1"

  metadata = {
    test = "true"
  }

  content_type = "application/json"

  versioning_legacy {
    type     = "versions"
    location = "tf-test-container-versions"
  }
}
```

### Global Read Access

```hcl
# Requires that a user know the object name they are attempting to download

resource "openstack_objectstorage_container_v1" "container_1" {
  region = "RegionOne"
  name   = "tf-test-container-1"

  container_read = ".r:*"
}
```

### Global Read and List Access

```hcl
# Any user can read any object, and list all objects in the container

resource "openstack_objectstorage_container_v1" "container_1" {
  region = "RegionOne"
  name   = "tf-test-container-1"

  container_read = ".r:*,.rlistings"
}
```

### Write-Only Access for a User

```hcl
data "openstack_identity_auth_scope_v3" "current" {
  name = "current"
}

# The named user can only upload objects, not read objects or list the container

resource "openstack_objectstorage_container_v1" "container_1" {
  region = "RegionOne"
  name   = "tf-test-container-1"

  container_read  = ".r:-${var.username}"
  container_write = "${data.openstack_identity_auth_scope_v3.current.project_id}:${var.username}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the container. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new container.

* `name` - (Required) A unique name for the container. Changing this creates a
    new container.

* `container_read` - (Optional) Sets an access control list (ACL) that grants
    read access. This header can contain a comma-delimited list of users that
    can read the container (allows the GET method for all objects in the
    container). Changing this updates the access control list read access.

* `container_sync_to` - (Optional) The destination for container synchronization.
    Changing this updates container synchronization.

* `container_sync_key` - (Optional) The secret key for container synchronization.
    Changing this updates container synchronization.

* `container_write` - (Optional) Sets an ACL that grants write access.
    Changing this updates the access control list write access.

* `versioning` - (Optional) A boolean that can enable or disable object
  versioning. The default value is `false`. To use this feature, your Swift
  version must be 2.24 or higher (as described in the [OpenStack Swift Ussuri release notes](https://docs.openstack.org/releasenotes/swift/ussuri.html#relnotes-2-24-0-stable-ussuri)),
  and a cloud administrator must have set the `allow_object_versioning = true`
  configuration option in Swift. If you cannot set this versioning type, you may
  want to consider using `versioning_legacy` instead.

* `versioning_legacy` - (Deprecated) Enable legacy object versioning. The structure is described below.

* `metadata` - (Optional) Custom key/value pairs to associate with the container.
    Changing this updates the existing container metadata.

* `content_type` - (Optional) The MIME type for the container. Changing this
    updates the MIME type.

* `storage_policy` - (Optional) The storage policy to be used for the container. 
    Changing this creates a new container.

* `force_destroy` -  (Optional, Default:false ) A boolean that indicates all objects should be deleted from the container so that the container can be destroyed without error. These objects are not recoverable.

The `versioning_legacy` block supports:

  * `type` - (Required) Versioning type which can be `versions` or `history` according to [Openstack documentation](https://docs.openstack.org/swift/latest/api/object_versioning.html).
  * `location` - (Required) Container in which versions will be stored.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `container_read` - See Argument Reference above.
* `container_sync_to` - See Argument Reference above.
* `container_sync_key` - See Argument Reference above.
* `container_write` - See Argument Reference above.
* `versioning` - See Argument Reference above.
* `versioning_legacy` - See Argument Reference above.
* `metadata` - See Argument Reference above.
* `content_type` - See Argument Reference above.
* `storage_policy` - See Argument Reference above.

## Import

This resource can be imported by specifying the name of the container:

Some attributes can't be imported :
* `force_destroy`
* `content_type`
* `metadata`
* `container_sync_to`
* `container_sync_key`

So you'll have to `terraform plan` and `terraform apply` after the import to fix those missing attributes.

```
$ terraform import openstack_objectstorage_container_v1.container_1 <name>
```
