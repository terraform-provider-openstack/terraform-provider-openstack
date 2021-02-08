---
layout: "openstack"
page_title: "OpenStack: openstack_objectstorage_object_v1"
sidebar_current: "docs-openstack-resource-objectstorage-object-v1"
description: |-
  Manages a V1 container object resource within OpenStack.
---

# openstack\_objectstorage\_object\_v1

Manages a V1 container object resource within OpenStack.

## Example Usage

### Example with simple content

```hcl
resource "openstack_objectstorage_container_v1" "container_1" {
  region = "RegionOne"
  name   = "tf-test-container-1"

  metadata {
    test = "true"
  }

  content_type = "application/json"
}

resource "openstack_objectstorage_object_v1" "doc_1" {
  region         = "RegionOne"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  name           = "test/default.json"
  metadata {
    test = "true"
  }

  content_type = "application/json"
  content      = <<JSON
               {
                 "foo" : "bar"
               }
JSON

}
```

### Example with content from file

```hcl
resource "openstack_objectstorage_container_v1" "container_1" {
  region = "RegionOne"
  name   = "tf-test-container-1"

  metadata {
    test = "true"
  }

  content_type = "application/json"
}

resource "openstack_objectstorage_object_v1" "doc_1" {
  region         = "RegionOne"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  name           = "test/default.json"
  metadata {
    test = "true"
  }

  content_type = "application/json"
  source       = "./default.json"
}
```

## Argument Reference

The following arguments are supported:

* `container_name` - (Required) A unique (within an account) name for the container. 
    The container name must be from 1 to 256 characters long and can start 
    with any character and contain any pattern. Character set must be UTF-8. 
    The container name cannot contain a slash (/) character because this 
    character delimits the container and object name. For example, the path 
    /v1/account/www/pages specifies the www container, not the www/pages container.

* `content` - (Optional) A string representing the content of the object. Conflicts with
   `source` and `copy_from`.

* `content_disposition` - (Optional) A string which specifies the override behavior for 
     the browser. For example, this header might specify that the browser use a download
     program to save this file rather than show the file, which is the default.
     
* `content_encoding` - (Optional) A string representing the value of the Content-Encoding
     metadata.

* `content_type` - (Optional) A string which sets the MIME type for the object.

* `copy_from` - (Optional) A string representing the name of an object 
    used to create the new object by copying the `copy_from` object. The value is in form 
    {container}/{object}. You must UTF-8-encode and then URL-encode the names of the 
    container and object before you include them in the header. Conflicts with `source` and 
    `content`.
             
* `delete_after` - (Optional) An integer representing the number of seconds after which the
    system removes the object. Internally, the Object Storage system stores this value in 
    the X-Delete-At metadata item.

* `delete_at` - (Optional) An string representing the date when the system removes the object. 
    For example, "2015-08-26" is equivalent to Mon, Wed, 26 Aug 2015 00:00:00 GMT.
    
* `detect_content_type` - (Optional) If set to true, Object Storage guesses the content 
    type based on the file extension and ignores the value sent in the Content-Type 
    header, if present.

* `etag` - (Optional) Used to trigger updates. The only meaningful value is ${md5(file("path/to/file"))}.

* `name` - (Required) A unique name for the object.

* `object_manifest` - (Optional) A string set to specify that this is a dynamic large 
    object manifest object. The value is the container and object name prefix of the
    segment objects in the form container/prefix. You must UTF-8-encode and then 
    URL-encode the names of the container and prefix before you include them in this 
    header.
    
* `region` - (Optional) The region in which to create the container. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new container.

* `source` - (Optional) A string representing the local path of a file which will be used
    as the object's content. Conflicts with `source` and `copy_from`.

## Attributes Reference

The following attributes are exported:

* `content_length` - If the operation succeeds, this value is zero (0) or the 
    length of informational or error text in the response body.
* `content_type` - If the operation succeeds, this value is the MIME type of the object. 
    If the operation fails, this value is the MIME type of the error text in the response 
    body.
* `date` - The date and time the system responded to the request, using the preferred 
    format of RFC 7231 as shown in this example Thu, 16 Jun 2016 15:10:38 GMT. The 
    time is always in UTC.
* `etag` - Whatever the value given in argument, will be overriden by the MD5 checksum of the uploaded object content. The value is not quoted. 
    If it is an SLO, it would be MD5 checksum of the segments’ etags.
* `last_modified` - The date and time when the object was last modified. The date and time 
    stamp format is ISO 8601:
       CCYY-MM-DDThh:mm:ss±hh:mm
    For example, 2015-08-27T09:49:58-05:00.
    The ±hh:mm value, if included, is the time zone as an offset from UTC. In the previous 
    example, the offset value is -05:00.
* `static_large_object` - True if object is a multipart_manifest.
* `trans_id` - A unique transaction ID for this request. Your service provider might 
    need this value if you report a problem.

* `container_name` - See Argument Reference above.
* `content` - See Argument Reference above.
* `content_disposition` - See Argument Reference above.
* `content_encoding` - See Argument Reference above.
* `copy_from` - See Argument Reference above.
* `delete_after` - See Argument Reference above.
* `delete_at` - See Argument Reference above.
* `detect_content_type` - See Argument Reference above.
* `name` - See Argument Reference above.
* `object_manifest` - See Argument Reference above.
* `region` - See Argument Reference above.
* `source` - See Argument Reference above.
