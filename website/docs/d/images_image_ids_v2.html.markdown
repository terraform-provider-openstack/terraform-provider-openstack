---
layout: "openstack"
page_title: "OpenStack: openstack_images_image_ids_v2"
sidebar_current: "docs-openstack-datasource-images-image-ids-v2"
description: |-
  Provides a list of Openstack Image IDs
---

# openstack\_images\_image\_ids\_v2

Use this data source to get a list of Openstack Image IDs matching the
specified criteria.

## Example Usage

```hcl
data "openstack_images_image_ids_v2" "images" {
  name_regex = "^Ubuntu 16\\.04.*-amd64"
  sort       = "updated_at"

  properties = {
    key = "value"
  }
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Glance client.
    A Glance client is needed to create an Image that can be used with
    a compute instance. If omitted, the `region` argument of the provider
    is used.

* `member_status` - (Optional) The status of the image. Must be one of
   "accepted", "pending", "rejected", or "all".

* `name` - (Optional) The name of the image. Cannot be used simultaneously
    with `name_regex`.

* `name_regex` - (Optional) The regular expressian of the name of the image.
    Cannot be used simultaneously with `name`. Unlike filtering by `name` the    
    `name_regex` filtering does by client on the result of OpenStack search
    query.

* `owner` - (Optional) The owner (UUID) of the image.

* `properties` - (Optional) a map of key/value pairs to match an image with.
    All specified properties must be matched. Unlike other options filtering    
    by `properties` does by client on the result of OpenStack search query.

* `size_min` - (Optional) The minimum size (in bytes) of the image to return.

* `size_max` - (Optional) The maximum size (in bytes) of the image to return.

* `sort` - (Optional) Sorts the response by one or more attribute and sort
    direction combinations. You can also set multiple sort keys and directions.
    Default direction is `desc`. Use the comma (,) character to separate
    multiple values. For example expression `sort = "name:asc,status"`
    sorts ascending by name and descending by status. `sort` cannot be used
    simultaneously with `sort_key`. If both are present in a configuration
    then only `sort` will be used.

* `sort_direction` - (Optional) Order the results in either `asc` or `desc`.
    Can be applied only with `sort_key`. Defaults to `asc`

* `sort_key` - (Optional) Sort images based on a certain key. Defaults to
    `name`. `sort_key` cannot be used simultaneously with `sort`. If both
    are present in a configuration then only `sort` will be used.

* `tag` - (Optional) Search for images with a specific tag.

* `tags` - (Optional) A list of tags required to be set on the image
      (all specified tags must be in the images tag list for it to be matched).

* `visibility` - (Optional) The visibility of the image. Must be one of
   "public", "private", "community", or "shared". Defaults to "private".

## Attributes Reference

`ids` is set to the list of Openstack Image IDs.
