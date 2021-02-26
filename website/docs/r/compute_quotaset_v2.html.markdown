---
layout: "openstack"
page_title: "OpenStack: openstack_compute_quotaset_v2"
sidebar_current: "docs-openstack-resource-compute-quotaset-v2"
description: |-
  Manages a V2 compute quotaset resource within OpenStack.
---

# openstack\_compute\_quotaset\_v2

Manages a V2 compute quotaset resource within OpenStack.

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

resource "openstack_compute_quotaset_v2" "quotaset_1" {
  project_id           = "${openstack_identity_project_v3.project_1.id}"
  key_pairs            = 10
  ram                  = 40960
  cores                = 32
  instances            = 20
  server_groups        = 4
  server_group_members = 8
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the volume. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new quotaset.

* `project_id` - (Required) ID of the project to manage quotas.
    Changing this creates a new quotaset.

* `fixed_ips` - (Optional) Quota value for fixed IPs.
    Changing this updates the existing quotaset.

* `floating_ips` - (Optional) Quota value for floating IPs.
    Changing this updates the existing quotaset.

* `injected_file_content_bytes` - (Optional) Quota value for content bytes
    of injected files. Changing this updates the existing quotaset.

* `injected_file_path_bytes` - (Optional) Quota value for path bytes of
    injected files. Changing this updates the existing quotaset.

* `injected_files` - (Optional) Quota value for injected files.
    Changing this updates the existing quotaset.

* `key_pairs` - (Optional) Quota value for key pairs.
    Changing this updates the existing quotaset.

* `metadata_items` - (Optional) Quota value for metadata items.
    Changing this updates the existing quotaset.

* `ram` - (Optional) Quota value for RAM.
    Changing this updates the existing quotaset.

* `security_group_rules` - (Optional) Quota value for security group rules.
    Changing this updates the existing quotaset.

* `security_groups` - (Optional) Quota value for security groups.
    Changing this updates the existing quotaset.

* `cores` - (Optional) Quota value for cores.
    Changing this updates the existing quotaset.

* `instances` - (Optional) Quota value for instances.
    Changing this updates the existing quotaset.

* `server_groups` - (Optional) Quota value for server groups.
    Changing this updates the existing quotaset.

* `server_group_members` - (Optional) Quota value for server groups members.
    Changing this updates the existing quotaset.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `fixed_ips` - See Argument Reference above.
* `floating_ips` - See Argument Reference above.
* `injected_file_content_bytes` - See Argument Reference above.
* `injected_file_path_bytes` - See Argument Reference above.
* `injected_files` - See Argument Reference above.
* `key_pairs` - See Argument Reference above.
* `metadata_items` - See Argument Reference above.
* `ram` - See Argument Reference above.
* `security_group_rules` - See Argument Reference above.
* `security_groups` - See Argument Reference above.
* `cores` - See Argument Reference above.
* `instances` - See Argument Reference above.
* `server_groups` - See Argument Reference above.
* `server_group_members` - See Argument Reference above.

## Import

Quotasets can be imported using the `project_id/region_name`, e.g.

```
$ terraform import openstack_compute_quotaset_v2.quotaset_1 2a0f2240-c5e6-41de-896d-e80d97428d6b/region_1
```
