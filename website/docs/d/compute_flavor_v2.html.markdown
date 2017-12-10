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

* `name` - (Optional) The name of the flavor.

* `min_ram` - (Optional) The minimum amount of RAM (in megabytes).

* `ram` - (Optional) The exact amount of RAM (in megabytes).

* `min_disk` - (Optional) The minimum amount of disk (in gigabytes).

* `disk` - (Optional) The exact amount of disk (in gigabytes).

* `vcpus` - (Optional) The amount of VCPUs.

* `swap` - (Optional) The amount of swap (in gigabytes).

* `rx_tx_factor` - (Optional) The `rx_tx_factor` of the flavor.


## Attributes Reference

`id` is set to the ID of the found flavor. In addition, the following attributes
are exported:

* `is_public` - Whether the flavor is public or private.
