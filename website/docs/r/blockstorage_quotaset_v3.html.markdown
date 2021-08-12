---
layout: "openstack"
page_title: "OpenStack: openstack_blockstorage_quotaset_v3"
sidebar_current: "docs-openstack-resource-blockstorage-quotaset-v3"
description: |-
  Manages a V3 quotaset resource within OpenStack.
---

# openstack\_blockstorage\_quotaset\_v3

Manages a V3 block storage quotaset resource within OpenStack.

~> **Note:** This usually requires admin privileges.

~> **Note:** This resource has a no-op deletion so no actual actions will be done against the OpenStack API
    in case of delete call.

~> **Note:** This resource has all-in creation so all optional quota arguments that were not specified are
    created with zero value. This excludes volume type quota.

## Example Usage

```hcl
resource "openstack_identity_project_v3" "project_1" {
  name = project_1
}

resource "openstack_blockstorage_quotaset_v3" "quotaset_1" {
  project_id = "${openstack_identity_project_v3.project_1.id}"
  volumes   = 10
  snapshots = 4
  gigabytes = 100
  per_volume_gigabytes = 10
  backups = 4
  backup_gigabytes = 10
  groups = 100
  volume_type_quota = {
    volumes_ssd = 30
    gigabytes_ssd = 500
    snapshots_ssd = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the volume. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new quotaset.

* `project_id` - (Required) ID of the project to manage quotas. Changing this
    creates a new quotaset.

* `volumes` - (Optional) Quota value for volumes. Changing this updates the
    existing quotaset.

* `snapshots` - (Optional) Quota value for snapshots. Changing this updates the
    existing quotaset.

* `gigabytes` - (Optional) Quota value for gigabytes. Changing this updates the
    existing quotaset.

* `per_volume_gigabytes` - (Optional) Quota value for gigabytes per volume .
    Changing this updates the existing quotaset.

* `backups` - (Optional) Quota value for backups. Changing this updates the
    existing quotaset.

* `backup_gigabytes` - (Optional) Quota value for backup gigabytes. Changing
    this updates the existing quotaset.

* `groups` - (Optional) Quota value for groups. Changing this updates the
    existing quotaset.

* `volume_type_quota` - (Optional)  Key/Value pairs for setting quota for
    volumes types. Possible keys are `snapshots_<volume_type_name>`,
    `volumes_<volume_type_name>` and `gigabytes_<volume_type_name>`.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `volumes` - See Argument Reference above.
* `snapshots` - See Argument Reference above.
* `gigabytes` - See Argument Reference above.
* `per_volume_gigabytes` - See Argument Reference above.
* `backups` - See Argument Reference above.
* `backup_gigabytes` - See Argument Reference above.
* `groups` - See Argument Reference above.
* `volume_type_quota` - See Argument Reference above.

## Import

Quotasets can be imported using the `project_id/region`, e.g.

```
$ terraform import openstack_blockstorage_quotaset_v3.quotaset_1 2a0f2240-c5e6-41de-896d-e80d97428d6b/region_1
```
