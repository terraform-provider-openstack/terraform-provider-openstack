---
subcategory: "DNS / Designate"
layout: "openstack"
page_title: "OpenStack: openstack_dns_zone_share_v2"
sidebar_current: "docs-openstack-resource-dns-zone-share-v2"
description: |-
  Data source for retrieving shared DNS zones in the OpenStack DNS Service (Designate V2).
---

# openstack\_dns\_zone\_share\_v2

The `openstack_dns_zone_share_v2` data source retrieves a list of DNS zone shares for a given zone. It can be used to discover which projects a zone has been shared with.

## Example Usage

```hcl
data "openstack_dns_zone_share_v2" "example" {
  zone_id           = "00000000-0000-0000-0000-000000000000"
  # Optionally, filter by target project ID.
  target_project_id = "11111111-1111-1111-1111-111111111111"
  # Optionally, specify the owner project ID. Required if the zone is not in your default project.
  project_id        = "22222222-2222-2222-2222-222222222222"
}
```

## Argument Reference

The following arguments are supported:

- `zone_id` (Required) - The ID of the DNS zone for which to list shares.

- `target_project_id` (Optional) - If provided, the data source returns only the shares with this target project ID.

- `project_id` (Optional) - The owner project ID. If omitted, it is derived from the zone details.


## Attributes Reference

- Shares
  - A list of objects representing DNS zone shares. Each object includes:
    - `share_id`: The ID of the share.
    - `project_id`: The target project ID associated with the share.
