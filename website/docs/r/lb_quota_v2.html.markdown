---
layout: "openstack"
page_title: "OpenStack: openstack_lb_quota_v2"
sidebar_current: "docs-openstack-resource-lb-quota-v2"
description: |-
  Manages a V2 quota resource within OpenStack.
---

# openstack\_lb\_quota\_v2

Manages a V2 load balancer quota resource within OpenStack.

~> **Note:** This usually requires admin privileges.

~> **Note:** This resource is only available for Octavia.

~> **Note:** This resource has a no-op deletion so no actual actions will be done against the OpenStack
   API in case of delete call.

~> **Note:** This resource has all-in creation so all optional quota arguments that were not specified are
   created with zero value.

~> **Note:** This resource has attributes that depend on octavia minor versions.
Please ensure your Openstack cloud supports the required [minor version](../#octavia-api-versioning).

## Example Usage

```hcl
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_lb_quota_v2" "quota_1" {
  project_id     = "${openstack_identity_project_v3.project_1.id}"
  loadbalancer   = 6
  listener       = 7
  member         = 8
  pool           = 9
  health_monitor = 10
  l7_policy      = 11
  l7_rule        = 12
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) ID of the project to manage quotas. Changing this
  creates a new quota.

* `region` - (Optional) Region in which to manage quotas. Changing this
  creates a new quota. If ommited, the region of the credentials is used.

* `loadbalancer` - (Optional) Quota value for loadbalancers. Changing this
  updates the existing quota. Omitting it sets it to 0.

* `listener` - (Optional) Quota value for listeners. Changing this updates
  the existing quota. Omitting it sets it to 0.

* `member` - (Optional) Quota value for members. Changing this updates
  the existing quota. Omitting it sets it to 0.

* `pool` - (Optional) Quota value for pools. Changing this updates the
  the existing quota. Omitting it sets it to 0.

* `health_monitor` - (Optional) Quota value for health_monitors. Changing
  this updates the existing quota. Omitting it sets it to 0.

* `l7_policy` - (Optional) Quota value for l7_policies. Changing this
  updates the existing quota. Omitting it sets it to 0. Available in
  **Octavia minor version 2.19**.

* `l7_rule` - (Optional) Quota value for l7_rules. Changing this
  updates the existing quota. Omitting it sets it to 0. Available in
  **Octavia minor version 2.19**.


## Attributes Reference

The following attributes are exported:

* `project_id` - See Argument Reference above.
* `loadbalancer` - See Argument Reference above.
* `listener` - See Argument Reference above.
* `member` - See Argument Reference above.
* `pool` - See Argument Reference above.
* `health_monitor` - See Argument Reference above.
* `l7_policy` - See Argument Reference above.
* `l7_rule` - See Argument Reference above.

## Import

Quotas can be imported using the `project_id/region_name`, where region_name is the
one defined is the Openstack credentials that are in use. E.g.

```
$ terraform import openstack_lb_quota_v2.quota_1 2a0f2240-c5e6-41de-896d-e80d97428d6b/region_1
```
