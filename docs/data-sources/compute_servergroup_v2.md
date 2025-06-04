---
subcategory: "Compute / Nova"
layout: "openstack"
page_title: "OpenStack: openstack_compute_servergroup_v2"
sidebar_current: "docs-openstack-datasource-compute-servergroup-v2"
description: |-
  Get information on an OpenStack server group.
---

# openstack\_compute\_servergroup\_v2

Use this data source to get information about server groups
by name.

## Example Usage

```hcl
data "openstack_compute_servergroup_v2" "test" {
  name = "test"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Compute client.
  If omitted, the `region` argument of the provider is used.

* `name` - (Required) The name of the server group.

## Attributes Reference

`id` is set to the ID of the found server group. In addition, the
following attributes are exported:

* `name` - See Argument Reference above.
* `user_id` - User ID of the server group.
* `project_id` - Project ID of the server group.
* `policies` - Policy name associated with the server group.
* `rules` - Rules which are applied to specified policy.
* `members` - The instances that are part of this server group.
* `metadata` - Metadata of the server group.
