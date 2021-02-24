---
layout: "openstack"
page_title: "OpenStack: openstack_compute_flavor_access_v2"
sidebar_current: "docs-openstack-resource-compute-flavor-access-v2"
description: |-
  Manages a project access for flavor V2 resource within OpenStack.
---

# openstack\_compute\_flavor\_access\_v2

Manages a project access for flavor V2 resource within OpenStack.

~> **Note:** You _must_ have admin privileges in your OpenStack cloud to use
this resource.

---

## Example Usage

```hcl
resource "openstack_identity_project_v3" "project_1" {
  name = "my-project"
}

resource "openstack_compute_flavor_v2" "flavor_1" {
  name      = "my-flavor"
  ram       = "8096"
  vcpus     = "2"
  disk      = "20"
  is_public = false
}

resource "openstack_compute_flavor_access_v2" "access_1" {
  tenant_id = "${openstack_identity_project_v3.project_1.id}"
  flavor_id = "${openstack_compute_flavor_v2.flavor_1.id}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    If omitted, the `region` argument of the provider is used.
    Changing this creates a new flavor access.

* `flavor_id` - (Required) The UUID of flavor to use. Changing this creates a new flavor access.

* `tenant_id` - (Required) The UUID of tenant which is allowed to use the flavor.
    Changing this creates a new flavor access.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `flavor_id` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.

## Import

This resource can be imported by specifying all two arguments, separated
by a forward slash:

```
$ terraform import openstack_compute_flavor_access_v2.access_1 <flavor_id>/<tenant_id>
```
