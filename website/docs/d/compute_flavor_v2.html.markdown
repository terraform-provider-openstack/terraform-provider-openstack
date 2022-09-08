---
layout: "openstack"
page_title: "OpenStack: openstack_compute_flavor_v2"
sidebar_current: "docs-openstack-datasource-compute-flavor-v2"
description: |-
  Get information on an OpenStack Flavor.
---

# openstack\_compute\_flavor\_v2

Use this data source to get the ID of an available OpenStack flavor.

## Example Usage

```hcl
data "openstack_compute_flavor_v2" "small" {
  vcpus = 1
  ram   = 512
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    If omitted, the `region` argument of the provider is used.

* `flavor_id` - (Optional) The ID of the flavor. Conflicts with the `name`,
    `min_ram` and `min_disk`

* `name` - (Optional) The name of the flavor. Conflicts with the `flavor_id`.

* `min_ram` - (Optional) The minimum amount of RAM (in megabytes). Conflicts
   with the `flavor_id`.

* `ram` - (Optional) The exact amount of RAM (in megabytes).

* `min_disk` - (Optional) The minimum amount of disk (in gigabytes). Conflicts
   with the `flavor_id`.

* `disk` - (Optional) The exact amount of disk (in gigabytes).

* `vcpus` - (Optional) The amount of VCPUs.

* `description` - (Optional) The description of the flavor.

* `swap` - (Optional) The amount of swap (in gigabytes).

* `rx_tx_factor` - (Optional) The `rx_tx_factor` of the flavor.

* `is_public` - (Optional) The flavor visibility.


## Attributes Reference

`id` is set to the ID of the found flavor. In addition, the following attributes
are exported:

* `extra_specs` - Key/Value pairs of metadata for the flavor.
