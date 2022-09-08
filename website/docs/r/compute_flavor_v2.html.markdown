---
layout: "openstack"
page_title: "OpenStack: openstack_compute_flavor_v2"
sidebar_current: "docs-openstack-resource-compute-flavor-v2"
description: |-
  Manages a V2 flavor resource within OpenStack.
---

# openstack\_compute\_flavor\_v2

Manages a V2 flavor resource within OpenStack.

## Example Usage

```hcl
resource "openstack_compute_flavor_v2" "test-flavor" {
  name  = "my-flavor"
  ram   = "8096"
  vcpus = "2"
  disk  = "20"

  extra_specs = {
    "hw:cpu_policy"        = "CPU-POLICY",
    "hw:cpu_thread_policy" = "CPU-THREAD-POLICY"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    Flavors are associated with accounts, but a Compute client is needed to
    create one. If omitted, the `region` argument of the provider is used.
    Changing this creates a new flavor.

* `name` - (Required) A unique name for the flavor. Changing this creates a new
    flavor.

* `description` - (Optional) The description of the flavor. Changing this
    updates the description of the flavor. Requires microversion >= 2.55.

* `ram` - (Required) The amount of RAM to use, in megabytes. Changing this
    creates a new flavor.

* `flavor_id` - (Optional) Unique ID (integer or UUID) of flavor to create. Changing
    this creates a new flavor.

* `vcpus` - (Required) The number of virtual CPUs to use. Changing this creates
    a new flavor.

* `disk` - (Required) The amount of disk space in GiB to use for the root
    (/) partition. Changing this creates a new flavor.

* `ephemeral` - (Optional) The amount of ephemeral in GiB. If unspecified,
    the default is 0. Changing this creates a new flavor.

* `swap` - (Optional) The amount of disk space in megabytes to use. If
    unspecified, the default is 0. Changing this creates a new flavor.

* `rx_tx_factor` - (Optional) RX/TX bandwith factor. The default is 1. Changing
    this creates a new flavor.

* `is_public` - (Optional) Whether the flavor is public. Changing this creates
    a new flavor.

* `extra_specs` - (Optional) Key/Value pairs of metadata for the flavor.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `ram` - See Argument Reference above.
* `vcpus` - See Argument Reference above.
* `disk` - See Argument Reference above.
* `ephemeral` - See Argument Reference above.
* `swap` - See Argument Reference above.
* `rx_tx_factor` - See Argument Reference above.
* `is_public` - See Argument Reference above.
* `extra_specs` - See Argument Reference above.

## Import

Flavors can be imported using the `ID`, e.g.

```
$ terraform import openstack_compute_flavor_v2.my-flavor 4142e64b-1b35-44a0-9b1e-5affc7af1106
```
