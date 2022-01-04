---
layout: "openstack"
page_title: "OpenStack: openstack_compute_quotaset_v2"
sidebar_current: "docs-openstack-datasource-compute-quotaset-v2"
description: |-
  Get information on a Compute Quotaset of a project.
---

# openstack\_compute\_quotaset\_v2

Use this data source to get the compute quotaset of an OpenStack project.

## Example Usage

```hcl
data "openstack_compute_quotaset_v2" "quota" {
  project_id = "2e367a3d29f94fd988e6ec54e305ec9d"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    If omitted, the `region` argument of the provider is used.

* `project_id` - (Required) The id of the project to retrieve the quotaset.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `cores` -  The number of allowed server cores.
* `instances` - The number of allowed servers.
* `key_pairs` - The number of allowed key pairs for each user.
* `metadata_items` - The number of allowed metadata items for each server.
* `ram` - The amount of allowed server RAM, in MiB.
* `server_groups` - The number of allowed server groups.
* `server_group_members` - The number of allowed members for each server group.
* `fixed_ips` - The number of allowed fixed IP addresses. Available until version 2.35.
* `floating_ips` - The number of allowed floating IP addresses. Available until version 2.35.
* `security_group_rules` - The number of allowed rules for each security group. Available until version 2.35.
* `security_groups` - The number of allowed security groups. Available until version 2.35.
* `injected_file_content_bytes` - The number of allowed bytes of content for each injected file. Available until version 2.56.
* `injected_file_path_bytes` - The number of allowed bytes for each injected file path. Available until version 2.56.
* `injected_files` - The number of allowed injected files. Available until version 2.56.
