---
layout: "openstack"
page_title: "OpenStack: openstack_dns_zone_v2"
sidebar_current: "docs-openstack-datasource-dns-zone-v2"
description: |-
  Get information on an OpenStack DNS Zone.
---

# openstack\_dns\_zone\_v2

Use this data source to get the ID of an available OpenStack DNS zone.

## Example Usage

```hcl
data "openstack_dns_zone_v2" "zone_1" {
  name = "example.com"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 DNS client.
  A DNS client is needed to retrieve zone ids. If omitted, the
  `region` argument of the provider is used.

* `name` - (Optional) The name of the zone.

* `project_id` - (Optional) The ID of the project the DNS zone is obtained from,
  sets `X-Auth-Sudo-Tenant-ID` header (requires an assigned user role in target project)

* `description` - (Optional) A description of the zone.

* `email` - (Optional) The email contact for the zone record.

* `status` - (Optional) The zone's status.

* `ttl` - (Optional) The time to live (TTL) of the zone.

* `type` - (Optional) The type of the zone. Can either be `PRIMARY` or `SECONDARY`.

* `all_projects` - (Optional) Try to obtain zone ID by listing all projects
  (requires admin role by default, depends on your policy configuration)

## Attributes Reference

`id` is set to the ID of the found zone. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `email` - See Argument Reference above.
* `type` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `description` - See Argument Reference above.
* `status` - See Argument Reference above.
* `attributes` - Attributes of the DNS Service scheduler.
* `masters` - An array of master DNS servers. When `type` is  `SECONDARY`.
* `created_at` - The time the zone was created.
* `updated_at` - The time the zone was last updated.
* `transferred_at` - The time the zone was transferred.
* `version` - The version of the zone.
* `serial` - The serial number of the zone.
* `pool_id` - The ID of the pool hosting the zone.
* `project_id` - The project ID that owns the zone.
