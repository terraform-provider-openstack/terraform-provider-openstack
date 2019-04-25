---
layout: "openstack"
page_title: "OpenStack: openstack_networking_addressscope_v2"
sidebar_current: "docs-openstack-datasource-networking-addressscope-v2"
description: |-
  Get information on an OpenStack Address Scope.
---

# openstack\_networking\_addressscope\_v2

Use this data source to get the ID of an available OpenStack address-scope.

## Example Usage

```hcl
data "openstack_networking_addressscope_v2" "public_addressscope" {
  name       = "public_addressscope"
  shared     = true
  ip_version = 4
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve address-scopes. If omitted, the
  `region` argument of the provider is used.

* `name` - (Optional) Name of the address-scope.

* `ip_version` - (Optional) IP version.

* `shared` - (Optional) Indicates whether this address-scope is shared across
    all projects.

* `project_id` - (Optional) The owner of the address-scope.

## Attributes Reference

`id` is set to the ID of the found address-scope. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `ip_version` - See Argument Reference above.
* `shared` - See Argument Reference above.
* `project_id` - See Argument Reference above.
