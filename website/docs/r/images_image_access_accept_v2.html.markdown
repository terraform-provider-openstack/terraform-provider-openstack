---
layout: "openstack"
page_title: "OpenStack: openstack_images_image_access_accept_v2"
sidebar_current: "docs-openstack-resource-images-image-access-accept-v2"
description: |-
  Manages a V2 Image membership proposal resource within OpenStack Glance.
---

# openstack\_images\_image\_access\_accept\_v2

Manages memberships status for the shared OpenStack Glance V2 Image within the
destination project, which has a member proposal.

## Example Usage

Accept a shared image membershipship proposal within the current project.

```hcl
data "openstack_images_image_v2" "rancheros" {
  name          = "RancherOS"
  visibility    = "shared"
  member_status = "all"
}

resource "openstack_images_image_access_accept_v2" "rancheros_member" {
  image_id = "${data.openstack_images_image_v2.rancheros.id}"
  status   = "accepted"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Glance client.
   A Glance client is needed to manage Image memberships. If omitted, the
  `region` argument of the provider is used. Changing this creates a new
  membership.

* `image_id` - (Required) The proposed image ID.

* `member_id` - (Optional) The member ID, e.g. the target project ID. Optional
  for admin accounts. Defaults to the current scope project ID.

* `status` - (Required) The membership proposal status. Can either be
  `accepted`, `rejected` or `pending`.

## Attributes Reference

The following attributes are exported:

* `created_at` - The date the image membership was created.
* `updated_at` - The date the image membership was last updated.
* `schema` - The membership schema.

## Import

Image access acceptance status can be imported using the `image_id`, e.g.

```
$ terraform import openstack_images_image_access_accept_v2 89c60255-9bd6-460c-822a-e2b959ede9d2
```
