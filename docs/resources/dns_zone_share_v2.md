---
subcategory: "DNS / Designate"
layout: "openstack"
page_title: "OpenStack: openstack_dns_zone_share_v2"
sidebar_current: "docs-openstack-resource-dns-zone-share-v2"
description: |-
  Manages a shared DNS zone in the OpenStack DNS Service
---

# openstack\_dns\_zone\_share\_v2

Manages the sharing of a DNS zone in the OpenStack DNS Service (Designate).

## Example Usage

```hcl
resource "openstack_dns_zone_share_v2" "example_share" {
  zone_id           = "00000000-0000-0000-0000-000000000000"
  target_project_id = "11111111111111111111111111111111"
  project_id        = "22222222222222222222222222222222"
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

DNS zone shares can be imported using a combination of the `zone_id`, `project_id`, `target_project_id` and `share_id`:

```bash
$ terraform import openstack_dns_zone_share_v2.imported <zone_id>:<project_id>:<target_project_id>/<share_id>
```

Replace `zone_id`, `project_id`, `target_project_id` and `share_id` with the appropriate IDs.

### Example of resource to be imported

```hcl
resource "openstack_dns_zone_share_v2" "imported" {
  zone_id           = "00000000-0000-0000-0000-000000000000"
  target_project_id = "33333333333333333333333333333333"
  project_id        = "22222222222222222222222222222222"
}
```
```bash
$ terraform import openstack_dns_zone_share_v2.imported 00000000-0000-0000-0000-000000000000:22222222222222222222222222222222:33333333333333333333333333333333/44444444-4444-4444-4444-444444444444
```

---

This documentation provides an overview of the `openstack_dns_zone_share_v2` resource, including its usage, arguments, attributes, and import instructions. 
