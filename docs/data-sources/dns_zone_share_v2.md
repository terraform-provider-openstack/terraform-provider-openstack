---
subcategory: "DNS / Designate"
layout: "openstack"
page_title: "OpenStack: openstack_dns_zone_share_v2"
sidebar_current: "docs-openstack-resource-dns-zone-share-v2"
description: |-
  Data source for retrieving shared DNS zone in the OpenStack DNS Service (Designate V2).
---

# openstack\_dns\_zone\_share\_v2

Use this data source to get information about a DNS zone share.

## Example Usage

```hcl
data "openstack_dns_zone_share_v2" "example" {
  zone_id = "00000000-0000-0000-0000-000000000000"

  # Optionally, filter by target project ID.
  target_project_id = "11111111-1111-1111-1111-111111111111"

  # Optionally, specify the owner project ID. Required if the zone is not in your default project.
  project_id = "22222222-2222-2222-2222-222222222222"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 DNS client. If
  omitted, the `region` argument of the provider is used. Changing this creates
  a new DNS zone share data source.

* `zone_id` - (Required) The ID of the DNS zone for which to get share.

* `share_id` - (Optional) The ID of the DNS zone share to retrieve. If
  provided, the data source returns only the share with this ID.

* `all_projects` - (Optional) If set to `true`, the data source will search
  across all projects. If set to `false`, it will only search within the
  current project. Defaults to `false`.

* `target_project_id` - (Optional) If provided, the data source returns the
  share with this target project ID.

* `project_id` - (Optional) The owner project ID. If omitted, it is derived
  from the zone share details.

## Attributes Reference

`id` is set to the ID of the found DNS share. In addition, the following
attributes are exported:

* `region` - See Argument Reference above.
* `zone_id` - See Argument Reference above.
* `share_id` - See Argument Reference above.
* `all_projects` - See Argument Reference above.
* `target_project_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `share_id` - The ID of the zone share.
