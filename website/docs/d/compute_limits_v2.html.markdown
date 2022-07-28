---
layout: "openstack"
page_title: "OpenStack: openstack_compute_limits_v2"
sidebar_current: "docs-openstack-datasource-compute-limits-v2"
description: |-
  Get information on a Compute Limits of a project.
---

# openstack\_compute\_limits\_v2

Use this data source to get the compute limits of an OpenStack project.

## Example Usage

```hcl
data "openstack_compute_limits_v2" "limits" {
  project_id = "2e367a3d29f94fd988e6ec54e305ec9d"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    If omitted, the `region` argument of the provider is used.

* `project_id` - (Required) The id of the project to retrieve the limits.


## Attributes Reference

The following attributes are exported:

* `project_id` - See Argument Reference above.
* `region` - See Argument Reference above.
* `max_total_cores` - The number of allowed server cores for the tenant.
* `max_image_meta` - The number of allowed metadata items for each image. Starting from version 2.39 this field is dropped from ‘os-limits’ response, because ‘image-metadata’ proxy API was deprecated. Available until version 2.38.
* `max_server_meta` - The number of allowed server groups for the tenant.
* `max_personality` - The number of allowed injected files for the tenant. Available until version 2.56.
* `max_personality_size` - The number of allowed bytes of content for each injected file. Available until version 2.56.
* `max_total_keypairs` - The number of allowed key pairs for the user.
* `max_security_groups` - The number of allowed security groups for the tenant. Available until version 2.35.
* `max_security_group_rules` - The number of allowed rules for each security group. Available until version 2.35.
* `max_server_groups` - The number of allowed server groups for the tenant.
* `max_server_group_members` - The number of allowed members for each server group.
* `max_total_floating_ips` - The number of allowed floating IP addresses for each tenant. Available until version 2.35.
* `max_total_instances` - The number of allowed servers for the tenant.
* `max_total_ram_size` - The number of allowed floating IP addresses for the tenant. Available until version 2.35.
* `total_cores_used` - The number of used server cores in the tenant.
* `total_instances_used` - The number of used server cores in the tenant.
* `total_floating_ips_used` - The number of used floating IP addresses in the tenant.
* `total_ram_used` - The amount of used server RAM in the tenant.
* `total_security_groups_used` - The number of used security groups in the tenant. Available until version 2.35.
* `total_server_groups_used` - The number of used server groups in each tenant.
