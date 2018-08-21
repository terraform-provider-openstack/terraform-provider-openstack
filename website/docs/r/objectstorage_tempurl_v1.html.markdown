---
layout: "openstack"
page_title: "OpenStack: openstack_objectstorage_tempurl_v1"
sidebar_current: "docs-openstack-resource-objectstorage-tempurl-v1"
description: |-
  Generate a TempURL for a Swift container and object.
---

# openstack\_objectstorage\_tempurl_v1

Use this resource to generate an OpenStack Object Storage temporary URL.

## Example Usage

```hcl
resource "openstack_objectstorage_tempurl_v1" "obj_tempurl" {
  container = "test"
  object = "container"
  method = "post"
  ttl = 20
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region the tempurl is located in.

* `container` - (Required) The container name the object belongs to.

* `object` - (Required) The object name the tempurl is for.

* `ttl` - (Required) The TTL, in seconds, for the URL. For how long it should
  be valid.

* `method` - (Optional) What methods are allowed for this URL.
  Valid values are `GET`, and `POST`. Default is `GET`.

## Attributes Reference

* `id` - Computed md5 hash based on the generated url
* `container` - See Argument Reference above.
* `object` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `method` - See Argument Reference above.
* `url` - The URL
* `region` - The region the endpoint is located in.
