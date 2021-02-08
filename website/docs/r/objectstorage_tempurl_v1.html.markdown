---
layout: "openstack"
page_title: "OpenStack: openstack_objectstorage_tempurl_v1"
sidebar_current: "docs-openstack-resource-objectstorage-tempurl-v1"
description: |-
  Generate a TempURL for a Swift container and object.
---

# openstack\_objectstorage\_tempurl\_v1

Use this resource to generate an OpenStack Object Storage temporary URL.

The temporary URL will be valid for as long as TTL is set to (in seconds).
Once the URL has expired, it will no longer be valid, but the resource
will remain in place. If you wish to automatically regenerate a URL, set
the `regenerate` argument to `true`. This will create a new resource with
a new ID and URL.

## Example Usage

```hcl
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "test"
  metadata = {
    Temp-URL-Key = "testkey"
  }
}

resource "openstack_objectstorage_object_v1" "object_1" {
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  name           = "test"
  content        = "Hello, world!"
}

resource "openstack_objectstorage_tempurl_v1" "obj_tempurl" {
  container = "${openstack_objectstorage_container_v1.container_1.name}"
  object    = "${openstack_objectstorage_object_v1.object_1.name}"
  method    = "post"
  ttl       = 20
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region the tempurl is located in.

* `container` - (Required) The container name the object belongs to.

* `object` - (Required) The object name the tempurl is for.

* `ttl` - (Required) The TTL, in seconds, for the URL. For how long it should
  be valid.

* `method` - (Optional) The method allowed when accessing this URL.
  Valid values are `GET`, and `POST`. Default is `GET`.

* `regenerate` - (Optional) Whether to automatically regenerate the URL when
  it has expired. If set to true, this will create a new resource with a new
  ID and new URL. Defaults to false.

## Attributes Reference

* `id` - Computed md5 hash based on the generated url
* `container` - See Argument Reference above.
* `object` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `method` - See Argument Reference above.
* `url` - The URL
* `region` - The region the endpoint is located in.
