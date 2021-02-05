---
layout: "openstack"
page_title: "OpenStack: openstack_networking_addressscope_v2"
sidebar_current: "docs-openstack-resource-networking-addressscope-v2"
description: |-
  Manages a V2 Neutron addressscope resource within OpenStack.
---

# openstack\_networking\_addressscope\_v2

Manages a V2 Neutron addressscope resource within OpenStack.

## Example Usage

### Create an Address-scope

```hcl
resource "openstack_networking_addressscope_v2" "addressscope_1" {
  name       = "addressscope_1"
  ip_version = 6
}
```

### Create a Subnet Pool from an Address-scope

```hcl
resource "openstack_networking_addressscope_v2" "addressscope_1" {
  name       = "addressscope_1"
  ip_version = 6
}

resource "openstack_networking_subnetpool_v2" "subnetpool_1" {
  name             = "subnetpool_1"
  prefixes         = ["fdf7:b13d:dead:beef::/64", "fd65:86cc:a334:39b7::/64"]
  address_scope_id = "${openstack_networking_addressscope_v2.addressscope_1.id}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a Neutron address-scope. If omitted,
    the `region` argument of the provider is used. Changing this creates a new
    address-scope.

* `name` - (Required) The name of the address-scope. Changing this updates the
    name of the existing address-scope.

* `ip_version` - (Optional) IP version, either 4 (default) or 6. Changing this
    creates a new address-scope.

* `shared` - (Optional) Indicates whether this address-scope is shared across
    all projects. Changing this updates the shared status of the existing
    address-scope.

* `project_id` - (Optional) The owner of the address-scope. Required if admin
    wants to create a address-scope for another project. Changing this creates a
    new address-scope.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `ip_version` - See Argument Reference above.
* `shared` - See Argument Reference above.
* `project_id` - See Argument Reference above.

## Import

Address-scopes can be imported using the `id`, e.g.

```
$ terraform import openstack_networking_addressscope_v2.addressscope_1 9cc35860-522a-4d35-974d-51d4b011801e
```
