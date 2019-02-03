---
layout: "openstack"
page_title: "OpenStack: openstack_objectstorage_container_v1"
sidebar_current: "docs-openstack-resource-objectstorage-container-v1"
description: |-
  Manages a V1 container resource within OpenStack.
---

# openstack\_objectstorage\_container_v1

Manages a V1 container resource within OpenStack.

## Example Usage

```hcl
resource "openstack_objectstorage_container_v1" "container_1" {
  region = "RegionOne"
  name   = "tf-test-container-1"

  metadata = {
    test = "true"
  }

  content_type = "application/json"

  versioning {
    type = "versions"
    location = "tf-test-container-versions"
  }
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

* `versioning` - (Optional) Enable object versioning. The structure is described below.

* `metadata` - (Optional) Custom key/value pairs to associate with the container.
    Changing this updates the existing container metadata.

* `content_type` - (Optional) The MIME type for the container. Changing this
    updates the MIME type.

* `force_destroy` -  (Optional, Default:false ) A boolean that indicates all objects should be deleted from the container so that the container can be destroyed without error. These objects are not recoverable.

The `versioning` block supports:

  * `type` - (Required) Versioning type which can be `versions` or `history` according to [Openstack documentation](https://docs.openstack.org/swift/latest/overview_object_versioning.html).
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
* `metadata` - See Argument Reference above.
* `content_type` - See Argument Reference above.
