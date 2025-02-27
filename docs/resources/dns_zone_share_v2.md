---
subcategory: "DNS / Designate"
layout: "openstack"
page_title: "OpenStack: openstack_dns_zone_share_v2"
sidebar_current: "docs-openstack-resource-dns-zone-share-v2"
description: |-
  Manages a shared DNS zone in the OpenStack DNS Service (Designate V2).
---

# openstack\_dns\_zone\_share\_v2

Manages the sharing of a DNS zone in the OpenStack DNS Service (Designate V2).

## Example Usage

```hcl
resource "openstack_dns_zone_share_v2" "example" {
  zone_id           = "00000000-0000-0000-0000-000000000000"
  target_project_id = "11111111-1111-1111-1111-111111111111"
  # project_id is optional; if omitted, the provider derives it from the zone details.
  project_id        = "22222222-2222-2222-2222-222222222222"
}
```

## Argument Reference

The following arguments are supported:

- `zone_id` (Required) - The ID of the DNS zone to be shared.

- `target_project_id` (Required) - The ID of the target project with which the DNS zone will be shared.

- `project_id` (Optional) - The ID of the owner project authorizing the share. This corresponds to the `X-Auth-Sudo-Project-Id` header in the Designate API. If omitted, the provider's `project_id` is used.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the resource, in the format `<zone_id>/<share_id>`.

- `share_id` - The ID of the created share.

## Import

DNS zone shares can be imported using either of the following formats:

 - Simplified Format (when the zone owner is the same as the one you're working from):

```bash
$ terraform import openstack_dns_zone_share_v2.imported <zone_id>/<share_id>
```

 - Full Format (explicitly specifying the owner and target projects):

```bash
$ terraform import openstack_dns_zone_share_v2.imported <zone_id>/<project_id>/<target_project_id>/<share_id>
```

Replace `zone_id`, `project_id`, `target_project_id` and `share_id` with the appropriate IDs.

---

This documentation provides an overview of the `openstack_dns_zone_share_v2` resource, including its usage, arguments, attributes, and import instructions. 
