---
layout: "openstack"
page_title: "OpenStack: openstack_compute_hypervisor_v2"
sidebar_current: "docs-openstack-datasource-compute-hypervisor-v2"
description: |-
  Get information on Openstack Hypervisor
---

# openstack\_compute\_hypervisor\_v2

Use this data source to get information about hypervisors
by hostname.

## Example Usage

```hcl
data "openstack_compute_hypervisor_v2" "host01" {
  hostname = "host01"
}
```

## Argument Reference

* `hostname` - The hostname of the hypervisor

## Attributes Reference

`id` is set to the ID of the found Hypervisor. In addition, the
following attributes are exported:

* `hostname` - See Argument Reference above.
* `host_ip` - The IP address of the Hypervisor
* `state` - The state of the hypervisor (`up` or `down`)
* `status` - The status of the hypervisor (`enabled` or `disabled`)
* `type` - The type of the hypervisor (example: `QEMU`)
* `vcpus` - The number of virtual CPUs the hypervisor can provide
* `memory` - The number in MegaBytes of memory the hypervisor can provide
* `disk` - The amount in GigaBytes of local storage the hypervisor can provide
