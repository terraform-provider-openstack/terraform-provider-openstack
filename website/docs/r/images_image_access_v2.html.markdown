---
layout: "openstack"
page_title: "OpenStack: openstack_images_image_access_v2"
sidebar_current: "docs-openstack-resource-images-image-access-v2"
description: |-
  Manages a V2 Image member resource within OpenStack Glance.
---

# openstack\_images\_image\_access\_v2

Manages members for the shared OpenStack Glance V2 Image within the source
project, which owns the Image.

## Example Usage

### Unprivileged user

Create a shared image and propose a membership to the
`bed6b6cbb86a4e2d8dc2735c2f1000e4` project ID.

```hcl
resource "openstack_images_image_v2" "rancheros" {
  name             = "RancherOS"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
  visibility       = "shared"

  properties = {
    key = "value"
  }
}

resource "openstack_images_image_access_v2" "rancheros_member" {
  image_id  = "${openstack_images_image_v2.rancheros.id}"
  member_id = "bed6b6cbb86a4e2d8dc2735c2f1000e4"
}
```

### Privileged user

Create a shared image and set a membership to the
`bed6b6cbb86a4e2d8dc2735c2f1000e4` project ID.

```hcl
resource "openstack_images_image_v2" "rancheros" {
  name             = "RancherOS"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"
  visibility       = "shared"

  properties = {
    key = "value"
  }
}

resource "openstack_images_image_access_v2" "rancheros_member" {
  image_id  = "${openstack_images_image_v2.rancheros.id}"
  member_id = "bed6b6cbb86a4e2d8dc2735c2f1000e4"
  status    = "accepted"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Glance client.
   A Glance client is needed to manage Image members. If omitted, the `region`
   argument of the provider is used. Changing this creates a new resource.

* `image_id` - (Required) The image ID.

* `member_id` - (Required) The member ID, e.g. the target project ID.

* `status` - (Optional) The member proposal status. Optional if admin wants to
  force the member proposal acceptance. Can either be `accepted`, `rejected` or
  `pending`. Defaults to `pending`. Foridden for non-admin users.

## Attributes Reference

The following attributes are exported:

* `created_at` - The date the image access was created.
* `updated_at` - The date the image access was last updated.
* `schema` - The member schema.

## Import

Image access can be imported using the `image_id` and the `member_id`,
separated by a slash, e.g.

```
$ terraform import openstack_images_image_access_v2 89c60255-9bd6-460c-822a-e2b959ede9d2/bed6b6cbb86a4e2d8dc2735c2f1000e4
```
