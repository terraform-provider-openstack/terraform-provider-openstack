---
subcategory: "Block Storage / Cinder"
layout: "openstack"
page_title: "OpenStack: openstack_blockstorage_volume_type_v3"
sidebar_current: "docs-openstack-resource-blockstorage-volume-type-v3"
description: |-
  Manages a V3 volume type resource within OpenStack.
---

# openstack\_blockstorage\_volume\_type\_v3

Manages a V3 block storage volume type resource within OpenStack.

~> **Note:** This usually requires admin privileges.


## Example Usage

### Basic Volume Type

```hcl
resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
  name        = "volume_type_1"
  description = "Volume type 1"
  extra_specs = {
      capabilities        = "gpu"
      volume_backend_name = "ssd"
  }
}

```

### Volume Type with multiattach enabled

```hcl
resource "openstack_blockstorage_volume_type_v3" "multiattach" {
  name        = "multiattach"
  description = "Multiattach-enabled volume type"
  extra_specs = {
      multiattach = "<is> True"
  }
}

```

### Encrypted Volume Type

```hcl
resource "openstack_blockstorage_volume_type_v3" "encrypted" {
  name = "encrypted"

  encryption {
    provider         = "luks"
    cipher           = "aes-xts-plain64"
    key_size         = 256
    control_location = "front-end"
  }
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the volume. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new quotaset.

* `name` - (Required) Name of the volume type.  Changing this
    updates the `name` of an existing volume type.

* `description` - (Optional) Human-readable description of the port. Changing
    this updates the `description` of an existing volume type.

* `is_public` - (Optional) Whether the volume type is public. Changing
    this updates the `is_public` of an existing volume type.

* `extra_specs` - (Optional) Key/Value pairs of metadata for the volume type.

* `encryption` - (Optional) The encryption spec for volumes created from
    this type. The `encryption` object structure is documented below.
    Cinder rejects creating, changing or removing the spec while any volume
    uses the type.

The `encryption` block supports:

* `provider` - (Required) The class that provides encryption support, such
    as `luks` or `plain`.

* `cipher` - (Optional) The encryption algorithm or mode, such as
    `aes-xts-plain64`.

* `key_size` - (Optional) Size of the encryption key, in bits.

* `control_location` - (Optional) Notional service where encryption is
    performed. Valid values are `front-end` and `back-end`. Defaults to
    `front-end`.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `is_public` - See Argument Reference above.
* `extra_specs` - See Argument Reference above.
* `encryption` - See Argument Reference above.

## Import

Volume types can be imported using the `volume_type_id`, e.g.

```
$ terraform import openstack_blockstorage_volume_type_v3.volume_type_1 941793f0-0a34-4bc4-b72e-a6326ae58283
```
