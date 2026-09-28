---
subcategory: "Block Storage / Cinder"
layout: "openstack"
page_title: "OpenStack: openstack_blockstorage_volume_type_v3"
sidebar_current: "docs-openstack-datasource-blockstorage-volume-type-v3"
description: |-
  Get information on an OpenStack Volume Type.
---

# openstack\_blockstorage\_volume\_type\_v3

Use this data source to get information about an existing volume type.

~> **Note:** This usually requires admin privileges.

## Example Usage

```hcl
data "openstack_blockstorage_volume_type_v3" "volume_type_1" {
  name = "volume_type_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V3 Block Storage
    client. If omitted, the `region` argument of the provider is used.

* `name` - (Optional) The name of the volume type.

* `is_public` - (Optional) Whether the volume type is public.

* `extra_specs` - (Optional) Key/Value pairs of metadata for the volume type.

## Attributes Reference

`id` is set to the ID of the found volume type. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - Human-readable description for the volume type.
* `is_public` - See Argument Reference above.
* `extra_specs` - See Argument Reference above.
* `qos_specs_id` - Qos Spec ID
* `public_access` - Volume Type access public attribute
