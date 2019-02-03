---
layout: "openstack"
page_title: "OpenStack: openstack_images_image_v2"
sidebar_current: "docs-openstack-resource-images-image-v2"
description: |-
  Manages a V2 Image resource within OpenStack Glance.
---

# openstack\_images\_image_v2

Manages a V2 Image resource within OpenStack Glance.

## Example Usage

```hcl
resource "openstack_images_image_v2" "rancheros" {
  name   = "RancherOS"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format = "qcow2"

  properties = {
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `container_format` - (Required) The container format. Must be one of
   "ami", "ari", "aki", "bare", "ovf".

* `disk_format` - (Required) The disk format. Must be one of
   "ami", "ari", "aki", "vhd", "vmdk", "raw", "qcow2", "vdi", "iso".

* `local_file_path` - (Optional) This is the filepath of the raw image file
   that will be uploaded to Glance. Conflicts with `image_source_url`.

* `image_cache_path` - (Optional) This is the directory where the images will
   be downloaded. Images will be stored with a filename corresponding to
   the url's md5 hash. Defaults to "$HOME/.terraform/image_cache"

* `image_source_url` - (Optional) This is the url of the raw image that will
   be downloaded in the `image_cache_path` before being uploaded to Glance.
   Glance is able to download image from internet but the `gophercloud` library
   does not yet provide a way to do so.
   Conflicts with `local_file_path`.

* `min_disk_gb` - (Optional) Amount of disk space (in GB) required to boot image.
   Defaults to 0.

* `min_ram_mb` - (Optional) Amount of ram (in MB) required to boot image.
   Defauts to 0.

* `name` - (Required) The name of the image.

* `properties` - (Optional) A map of key/value pairs to set freeform
    information about an image. See the "Notes" section for further
    information about properties.

* `protected` - (Optional) If true, image will not be deletable.
   Defaults to false.

* `region` - (Optional) The region in which to obtain the V2 Glance client.
    A Glance client is needed to create an Image that can be used with
    a compute instance. If omitted, the `region` argument of the provider
    is used. Changing this creates a new Image.

* `tags` - (Optional) The tags of the image. It must be a list of strings.
    At this time, it is not possible to delete all tags of an image.

* `verify_checksum` - (Optional) If false, the checksum will not be verified
    once the image is finished uploading. Defaults to true.

* `visibility` - (Optional) The visibility of the image. Must be one of
   "public", "private", "community", or "shared". The ability to set the
   visibility depends upon the configuration of the OpenStack cloud.

## Attributes Reference

The following attributes are exported:

* `checksum` - The checksum of the data associated with the image.
* `container_format` - See Argument Reference above.
* `created_at` - The date the image was created.
* `disk_format` - See Argument Reference above.
* `file` - the trailing path after the glance
   endpoint that represent the location of the image
   or the path to retrieve it.
* `id` - A unique ID assigned by Glance.
* `metadata` - The metadata associated with the image.
   Image metadata allow for meaningfully define the image properties
   and tags. See https://docs.openstack.org/glance/latest/user/metadefs-concepts.html.
* `min_disk_gb` - See Argument Reference above.
* `min_ram_mb` - See Argument Reference above.
* `name` - See Argument Reference above.
* `owner` - The id of the openstack user who owns the image.
* `properties` - See Argument Reference above.
* `protected` - See Argument Reference above.
* `region` - See Argument Reference above.
* `schema` - The path to the JSON-schema that represent
   the image or image
* `size_bytes` - The size in bytes of the data associated with the image.
* `status` - The status of the image. It can be "queued", "active"
   or "saving".
* `tags` - See Argument Reference above.
* `updated_at` - The date the image was last updated.
* `update_at` - (**Deprecated** - use `updated_at` instead)
* `visibility` - See Argument Reference above.

## Notes

### Properties

This resource supports the ability to add properties to a resource during
creation as well as add, update, and delete properties during an update of this
resource.

Newer versions of OpenStack are adding some read-only properties to each image.
These properties start with the prefix `os_`. If these properties are detected,
this resource will automatically reconcile these with the user-provided
properties.

In addition, the `direct_url` property is also automatically reconciled if the
Image Service set it.

## Import

Images can be imported using the `id`, e.g.

```
$ terraform import openstack_images_image_v2.rancheros 89c60255-9bd6-460c-822a-e2b959ede9d2
```
