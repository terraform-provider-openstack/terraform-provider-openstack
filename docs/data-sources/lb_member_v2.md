---
subcategory: "Load Balancing as a Service / Octavia"
layout: "openstack"
page_title: "OpenStack: openstack_lb_member_v2"
sidebar_current: "docs-openstack-datasource-lb-member-v2"
description: |-
  Get information on an OpenStack Load Balancer Member.
---

# openstack\_lb\_member\_v2

Use this data source to get the ID of an OpenStack Load Balancer member.

## Example Usage

```hcl
data "openstack_lb_member_v2" "member_1" {
  name = "member_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Load Balancer
  client. If omitted, the `region` argument of the provider is used.

* `pool_id` - (Required) The Pool to which the Member belongs.

* `member_id` - (Optional) The ID of the member. Exactly one of `name`,
  `member_id` is required to be set.

* `name` - (Optional) The name of the member. Exactly one of `name`,
  `member_id` is required to be set.

* `weight` - (Optional) Weight of Member.

* `address` - (Optional) The IP address of the Member.

## Attributes Reference

`id` is set to the ID of the found member. In addition, the following attributes
are exported:

* `project_id` - The owner (project/tenant) ID of the member.

* `name` - The name of the member.

* `weight` - Weight of Member.

* `admin_state_up` - The administrative state of the member, which is up (true)
  or down (false).

* `subnet_id` - Parameter value for the subnet UUID.

* `pool_id` - The Pool to which the Member belongs.

* `address` - The IP address of the Member.

* `protocol_port` - The port on which the application is hosted.

* `provisioning_status` - The provisioning status of the member.

* `operating_status` - The operating status of the member.

* `backup` - Whether the member is a backup. A backup member receives traffic
  only if all non-backup members are unavailable.

* `monitor_address` - An alternate IP address used for health monitoring a backend member.

* `monitor_port` - An alternate protocol port used for health monitoring a backend member.

* `tags` - A list of simple strings assigned to the resource.
