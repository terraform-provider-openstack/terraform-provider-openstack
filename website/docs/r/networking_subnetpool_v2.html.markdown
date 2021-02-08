---
layout: "openstack"
page_title: "OpenStack: openstack_networking_subnetpool_v2"
sidebar_current: "docs-openstack-resource-networking-subnetpool-v2"
description: |-
  Manages a V2 Neutron subnetpool resource within OpenStack.
---

# openstack\_networking\_subnetpool\_v2

Manages a V2 Neutron subnetpool resource within OpenStack.

## Example Usage

### Create a Subnet Pool

```hcl
resource "openstack_networking_subnetpool_v2" "subnetpool_1" {
  name       = "subnetpool_1"
  ip_version = 6
  prefixes   = ["fdf7:b13d:dead:beef::/64", "fd65:86cc:a334:39b7::/64"]
}
```

### Create a Subnet from a Subnet Pool

```hcl
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnetpool_v2" "subnetpool_1" {
  name     = "subnetpool_1"
  prefixes = ["10.11.12.0/24"]
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name          = "subnet_1"
  cidr          = "10.11.12.0/25"
  network_id    = "${openstack_networking_network_v2.network_1.id}"
  subnetpool_id = "${openstack_networking_subnetpool_v2.subnetpool_1.id}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a Neutron subnetpool. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    subnetpool.

* `name` - (Required) The name of the subnetpool. Changing this updates the name of
    the existing subnetpool.

* `default_quota` - (Optional) The per-project quota on the prefix space that can be
    allocated from the subnetpool for project subnets. Changing this updates the
    default quota of the existing subnetpool.

* `project_id` - (Optional) The owner of the subnetpool. Required if admin wants to
    create a subnetpool for another project. Changing this creates a new subnetpool.

* `prefixes` - (Required) A list of subnet prefixes to assign to the subnetpool.
    Neutron API merges adjacent prefixes and treats them as a single prefix. Each
    subnet prefix must be unique among all subnet prefixes in all subnetpools that
    are associated with the address scope. Changing this updates the prefixes list
    of the existing subnetpool.

* `default_prefixlen` - (Optional) The size of the prefix to allocate when the cidr
    or prefixlen attributes are omitted when you create the subnet. Defaults to the
    MinPrefixLen. Changing this updates the default prefixlen of the existing
    subnetpool.

* `min_prefixlen` - (Optional) The smallest prefix that can be allocated from a
    subnetpool. For IPv4 subnetpools, default is 8. For IPv6 subnetpools, default
    is 64. Changing this updates the min prefixlen of the existing subnetpool.

* `max_prefixlen` - (Optional) The maximum prefix size that can be allocated from
    the subnetpool. For IPv4 subnetpools, default is 32. For IPv6 subnetpools,
    default is 128. Changing this updates the max prefixlen of the existing
    subnetpool.

* `address_scope_id` - (Optional) The Neutron address scope to assign to the
    subnetpool. Changing this updates the address scope id of the existing
    subnetpool.

* `shared` - (Optional) Indicates whether this subnetpool is shared across
    all projects. Changing this updates the shared status of the existing
    subnetpool.

* `description` - (Optional) The human-readable description for the subnetpool.
    Changing this updates the description of the existing subnetpool.

* `is_default` - (Optional) Indicates whether the subnetpool is default
    subnetpool or not. Changing this updates the default status of the existing
    subnetpool.

* `value_specs` - (Optional) Map of additional options.

* `tags` - (Optional) A set of string tags for the subnetpool.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `default_quota` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `created_at` - The time at which subnetpool was created.
* `updated_at` - The time at which subnetpool was created.
* `prefixes` - See Argument Reference above.
* `default_prefixlen` - See Argument Reference above.
* `min_prefixlen` - See Argument Reference above.
* `max_prefixlen` - See Argument Reference above.
* `address_scope_id` - See Argument Reference above.
* `ip_version` - The IP protocol version.
* `shared` - See Argument Reference above.
* `description` - See Argument Reference above.
* `is_default` - See Argument Reference above.
* `revision_number` - The revision number of the subnetpool.
* `value_specs` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `all_tags` - The collection of tags assigned on the subnetpool, which have been
  explicitly and implicitly added.

## Import

Subnetpools can be imported using the `id`, e.g.

```
$ terraform import openstack_networking_subnetpool_v2.subnetpool_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
