---
layout: "openstack"
page_title: "OpenStack: openstack_compute_volume_attach_v2"
sidebar_current: "docs-openstack-resource-compute-volume-attach-v2"
description: |-
  Attaches a Block Storage Volume to an Instance.
---

# openstack\_compute\_volume\_attach\_v2

Attaches a Block Storage Volume to an Instance using the OpenStack
Compute (Nova) v2 API.

## Example Usage

### Basic attachment of a single volume to a single instance

```hcl
resource "openstack_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "openstack_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "openstack_compute_volume_attach_v2" "va_1" {
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  volume_id   = "${openstack_blockstorage_volume_v2.volume_1.id}"
}
```

### Attaching multiple volumes to a single instance

```hcl
resource "openstack_blockstorage_volume_v2" "volumes" {
  count = 2
  name  = "${format("vol-%02d", count.index + 1)}"
  size  = 1
}

resource "openstack_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "openstack_compute_volume_attach_v2" "attachments" {
  count       = 2
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  volume_id   = "${openstack_blockstorage_volume_v2.volumes.*.id[count.index]}"
}

output "volume_devices" {
  value = "${openstack_compute_volume_attach_v2.attachments.*.device}"
}
```

Note that the above example will not guarantee that the volumes are attached in
a deterministic manner. The volumes will be attached in a seemingly random
order.

If you want to ensure that the volumes are attached in a given order, create
explicit dependencies between the volumes, such as:

```hcl
resource "openstack_blockstorage_volume_v2" "volumes" {
  count = 2
  name  = "${format("vol-%02d", count.index + 1)}"
  size  = 1
}

resource "openstack_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "openstack_compute_volume_attach_v2" "attach_1" {
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  volume_id   = "${openstack_blockstorage_volume_v2.volumes.0.id}"
}

resource "openstack_compute_volume_attach_v2" "attach_2" {
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  volume_id   = "${openstack_blockstorage_volume_v2.volumes.1.id}"

  depends_on = ["openstack_compute_volume_attach_v2.attach_1"]
}

output "volume_devices" {
  value = "${openstack_compute_volume_attach_v2.attachments.*.device}"
}
```

### Using Multiattach-enabled volumes

Multiattach Volumes are dependent upon your OpenStack cloud and not all
clouds support multiattach.

```hcl
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name        = "volume_1"
  size        = 1
  multiattach = true
}

resource "openstack_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "openstack_compute_instance_v2" "instance_2" {
  name            = "instance_2"
  security_groups = ["default"]
}

resource "openstack_compute_volume_attach_v2" "va_1" {
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  volume_id   = "${openstack_blockstorage_volume_v2.volume_1.id}"
  multiattach = true
}

resource "openstack_compute_volume_attach_v2" "va_2" {
  instance_id = "${openstack_compute_instance_v2.instance_2.id}"
  volume_id   = "${openstack_blockstorage_volume_v2.volume_1.id}"
  multiattach = true

  depends_on = ["openstack_compute_volume_attach_v2.va_1"]
}
```

It is recommended to use `depends_on` for the attach resources
to enforce the volume attachments to happen one at a time.

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    A Compute client is needed to create a volume attachment. If omitted, the
    `region` argument of the provider is used. Changing this creates a
    new volume attachment.

* `instance_id` - (Required) The ID of the Instance to attach the Volume to.

* `volume_id` - (Required) The ID of the Volume to attach to an Instance.

* `device` - (Optional) The device of the volume attachment (ex: `/dev/vdc`).
  _NOTE_: Being able to specify a device is dependent upon the hypervisor in
  use. There is a chance that the device specified in Terraform will not be
  the same device the hypervisor chose. If this happens, Terraform will wish
  to update the device upon subsequent applying which will cause the volume
  to be detached and reattached indefinitely. Please use with caution.

* `multiattach` - (Optional) Enable attachment of multiattach-capable volumes.

* `vendor_options` - (Optional) Map of additional vendor-specific options.
  Supported options are described below.

The `vendor_options` block supports:

* `ignore_volume_confirmation` - (Optional) Boolean to control whether
  to ignore volume status confirmation of the attached volume. This can be helpful
  to work with some OpenStack clouds which don't have the Block Storage V3 API available.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `instance_id` - See Argument Reference above.
* `volume_id` - See Argument Reference above.
* `device` - See Argument Reference above. _NOTE_: The correctness of this
  information is dependent upon the hypervisor in use. In some cases, this
  should not be used as an authoritative piece of information.
* `multiattach` - See Argument Reference above.

## Import

Volume Attachments can be imported using the Instance ID and Volume ID
separated by a slash, e.g.

```
$ terraform import openstack_compute_volume_attach_v2.va_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
