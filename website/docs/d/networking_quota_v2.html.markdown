---
layout: "openstack"
page_title: "OpenStack: openstack_networking_quota_v2"
sidebar_current: "docs-openstack-datasource-networking-quota-v2"
description: |-
  Get information on a NEtworking Quota of a project.
---

# openstack\_networking\_quota\_v2

Use this data source to get the networking quota of an OpenStack project.

## Example Usage

```hcl
data "openstack_networking_quota_v2" "quota" {
  project_id = "2e367a3d29f94fd988e6ec54e305ec9d"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Network client.
    If omitted, the `region` argument of the provider is used.

* `project_id` - (Required) The id of the project to retrieve the quota.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `floatingip` -  The number of allowed floating ips.
* `network` - The number of allowed networks.
* `port` - The number of allowed ports.
* `rbac_policy` - The number of allowed rbac policies.
* `router` - The amount of allowed routers.
* `security_group` - The number of allowed security groups.
* `security_group_rule` - The number of allowed security group rules.
* `subnet` - The number of allowed subnets.
* `subnetpool-` - The number of allowed subnet pools.
