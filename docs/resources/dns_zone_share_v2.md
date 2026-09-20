---
subcategory: "DNS / Designate"
layout: "openstack"
page_title: "OpenStack: openstack_dns_zone_share_v2"
sidebar_current: "docs-openstack-resource-dns-zone-share-v2"
description: |-
  Manages the sharing of a DNS zone in the OpenStack DNS Service.
---

# openstack\_dns\_zone\_share\_v2

Manages the sharing of a DNS zone in the OpenStack DNS Service.

## Example Usage

```hcl
resource "openstack_dns_zone_share_v2" "example" {
  zone_id           = "00000000-0000-0000-0000-000000000000"
  target_project_id = "11111111-1111-1111-1111-111111111111"
  # project_id is optional; if omitted, the provider derives it from the zone details.
  project_id = "22222222-2222-2222-2222-222222222222"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 DNS client. If
  omitted, the `region` argument of the provider is used. Changing this creates
  a new DNS zone share.

* `zone_id` - (Required) The ID of the DNS zone to be shared.

* `target_project_id` - (Required) The ID of the target project with which the
  DNS zone will be shared.

* `project_id` - (Optional) The ID of the project DNS zone is created for, sets
  `X-Auth-Sudo-Tenant-ID` header (requires an assigned user role in target
  project).

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `zone_id` - See Argument Reference above.
* `target_project_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.

## Import

DNS zone share can be imported by specifying the zone ID with share ID and optional project ID:

```shell
terraform import openstack_dns_zone_share_v2.share_1 60cbdc69-64f9-49ee-b294-352e71e22827/0e1dae51-aee2-4b44-962f-885bb69f3a5c
terraform import openstack_dns_zone_share_v2.share_1 60cbdc69-64f9-49ee-b294-352e71e22827/0e1dae51-aee2-4b44-962f-885bb69f3a5c/eb92139f6c054a878852ac9e8cbe612a
```
