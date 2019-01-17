---
layout: "openstack"
page_title: "OpenStack: openstack_networking_subnetpool_v2"
sidebar_current: "docs-openstack-datasource-networking-subnetpool-v2"
description: |-
  Get information on an OpenStack Subnetpool.
---

# openstack\_networking\_subnetpool\_v2

Use this data source to get the ID of an available OpenStack subnetpool.

## Example Usage

```hcl
data "openstack_networking_subnetpool_v2" "subnetpool_1" {
  name = "subnetpool_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to retrieve a subnetpool id. If omitted, the
    `region` argument of the provider is used.

* `name` - (Optional) The name of the subnetpool.

* `default_quota` - (Optional) The per-project quota on the prefix space that
    can be allocated from the subnetpool for project subnets.

* `project_id` - (Optional) The owner of the subnetpool.

* `prefixes` - (Optional) A list of subnet prefixes that are assigned to the
    subnetpool.

* `default_prefixlen` - (Optional) The size of the subnetpool default prefix
    length.

* `min_prefixlen` - (Optional) The size of the subnetpool min prefix length.

* `max_prefixlen` - (Optional) The size of the subnetpool max prefix length.

* `address_scope_id` - (Optional) The Neutron address scope that subnetpools
    is assigned to.

* `ip_version` - The IP protocol version.

* `shared` - (Optional) Whether this subnetpool is shared across all projects.

* `description` - (Optional) The human-readable description for the subnetpool.

* `is_default` - (Optional) Whether the subnetpool is default subnetpool or not.

* `tags` - (Optional) The list of subnetpool tags to filter.

## Attributes Reference

`id` is set to the ID of the found subnetpool. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `default_quota` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `created_at` -  The time at which subnetpool was created.
* `updated_at` - The time at which subnetpool was created.
* `prefixes` - See Argument Reference above.
* `default_prefixlen` - See Argument Reference above.
* `min_prefixlen` - See Argument Reference above.
* `max_prefixlen` - See Argument Reference above.
* `address_scope_id` - See Argument Reference above.
* `ip_version` -The IP protocol version.
* `shared` - See Argument Reference above.
* `description` - See Argument Reference above.
* `is_default` - See Argument Reference above.
* `revision_number` - The revision number of the subnetpool.
* `all_tags` - The set of string tags applied on the subnetpool.
