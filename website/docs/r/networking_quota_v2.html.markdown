---
layout: "openstack"
page_title: "OpenStack: openstack_networking_quota_v2"
sidebar_current: "docs-openstack-resource-networking-quota-v2"
description: |-
  Manages a V2 networking quota resource within OpenStack.
---

# openstack\_networking\_quota\_v2

Manages a V2 networking quota resource within OpenStack.

~> **Note:** This usually requires admin privileges.

~> **Note:** This resource has a no-op deletion so no actual actions will be done against the OpenStack API
    in case of delete call.

~> **Note:** This resource has all-in creation so all optional quota arguments that were not specified are
    created with zero value.

## Example Usage

```hcl
resource "openstack_identity_project_v3" "project_1" {
  name = project_1
}

resource "openstack_networking_quota_v2" "quota_1" {
  project_id          = "${openstack_identity_project_v3.project_1.id}"
  floatingip          = 10
  network             = 4
  port                = 100
  rbac_policy         = 10
  router              = 4
  security_group      = 10
  security_group_rule = 100
  subnet              = 8
  subnetpool          = 2
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the quota. If
    omitted, the `region` argument of the provider is used. Changing this
    creates new quota.

* `project_id` - (Required) ID of the project to manage quota. Changing this
    creates new quota.

* `floatingip` - (Optional) Quota value for floating IPs. Changing this updates the
    existing quota.

* `network` - (Optional) Quota value for networks. Changing this updates the
    existing quota.

* `port` - (Optional) Quota value for ports. Changing this updates the
    existing quota.

* `rbac_policy` - (Optional) Quota value for RBAC policies.
    Changing this updates the existing quota.

* `router` - (Optional) Quota value for routers. Changing this updates the
    existing quota.

* `security_group` - (Optional) Quota value for security groups. Changing
    this updates the existing quota.

* `security_group_rule` - (Optional) Quota value for security group rules.
    Changing this updates the existing quota.

* `subnet` - (Optional) Quota value for subnets. Changing
    this updates the existing quota.

* `subnetpool` - (Optional) Quota value for subnetpools.
    Changing this updates the existing quota.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `floatingip` - See Argument Reference above.
* `network` - See Argument Reference above.
* `port` - See Argument Reference above.
* `rbac_policy` - See Argument Reference above.
* `router` - See Argument Reference above.
* `security_group` - See Argument Reference above.
* `security_group_rule` - See Argument Reference above.
* `subnet` - See Argument Reference above.
* `subnetpool` - See Argument Reference above.

## Import

Quotas can be imported using the `project_id/region_name`, e.g.

```
$ terraform import openstack_networking_quota_v2.quota_1 2a0f2240-c5e6-41de-896d-e80d97428d6b/region_1
```
