---
layout: "openstack"
page_title: "OpenStack: openstack_compute_aggregate_v2"
sidebar_current: "docs-openstack-resource-compute-aggregate-v2"
description: |-
  Manages a Host Aggregate within Openstack Nova
---

# openstack\_compute\_aggregate\_v2

Manages a Host Aggregate within Openstack Nova.

## Example Usage

### Full example

```hcl
resource "openstack_compute_aggregate_v2" "dell_servers" {
  name = "dell_servers"
  zone = "nova"
  metadata = {
    cpus = "56"
  }
  hosts = [
    "myhost01.example.com",
    "myhost02.example.com",
  ]
}
```

### Minimum required example

```hcl
resource "openstack_compute_aggregate_v2" "test" {
  name = "test"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name of the Host Aggregate
* `zone` - (Optional) The name of the Availability Zone to use. If ommited, it will take the default
  availability zone.
* `hosts` - (Optional) The list of hosts contained in the Host Aggregate. The hosts must be added
  to Openstack and visible in the web interface, or the provider will fail to add them to the host
  aggregate.
* `metadata` - (Optional) The metadata of the Host Aggregate. Can be useful to indicate scheduler hints.

## Import

You can import an existing Host Aggregate by their ID.
```
$ terraform import openstack_compute_aggregate_v2.myaggregate 24
```

The ID can be obtained with an openstack command:
```
$ openstack aggregate list
+----+------+-------------------+
| ID | Name | Availability Zone |
+----+------+-------------------+
| 59 | test | None              |
+----+------+-------------------+
```
