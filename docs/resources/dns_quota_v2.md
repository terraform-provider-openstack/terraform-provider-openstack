---
subcategory: "DNS / Designate"
layout: "openstack"
page_title: "OpenStack: openstack_dns_quota_v2"
sidebar_current: "docs-openstack-resource-dns-quota-v2"
description: |-
  Manages DNS quota in OpenStack DNS Service.
---

# openstack\_dns\_quota\_v2

Manages DNS quota in OpenStack DNS Service.

~> **Note:** This usually requires admin privileges.

~> **Note:** This resource has a no-op deletion so no actual actions will be
done against the OpenStack API in case of delete call.

## Example Usage

```hcl
resource "openstack_identity_project_v3" "project_1" {
  name = project_1
}

resource "openstack_dns_quota_v2" "quota_1" {
  project_id        = openstack_identity_project_v3.project_1.id
  api_export_size   = 4
  recordset_records = 10
  zone_records      = 100
  zone_recordsets   = 8
  zones             = 2
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 DNS client. If
  omitted, the `region` argument of the provider is used. Changing this creates
  a new DNS quota.

* `project_id` - (Required) ID of the project to manage quota. Changing this
  creates new quota.

* `api_export_size` - (Optional) The maximum number of zones that can be
  exported via the API.

* `recordset_records` - (Optional) The maximum number of records in a
  recordset.

* `zone_records` - (Optional) The maximum number of records in a zone.

* `zone_recordsets` - (Optional) The maximum number of recordsets in a zone.

* `zones` - (Optional) The maximum number of zones that can be created.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `api_export_size` - See Argument Reference above.
* `recordset_records` - See Argument Reference above.
* `zone_records` - See Argument Reference above.
* `zone_recordsets` - See Argument Reference above.
* `zones` - See Argument Reference above.

## Import

Quotas can be imported using the `project_id/region_name`, e.g.

```shell
terraform import openstack_dns_quota_v2.quota_1 2a0f2240-c5e6-41de-896d-e80d97428d6b/region_1
```
