---
layout: "openstack"
page_title: "OpenStack: openstack_dns_recordset_v2"
sidebar_current: "docs-openstack-resource-dns-recordset-v2"
description: |-
  Manages a DNS record set in the OpenStack DNS Service
---

# openstack\_dns\_recordset\_v2

Manages a DNS record set in the OpenStack DNS Service.

## Example Usage

### Automatically detect the correct network

```hcl
resource "openstack_dns_zone_v2" "example_zone" {
  name        = "example.com."
  email       = "email2@example.com"
  description = "a zone"
  ttl         = 6000
  type        = "PRIMARY"
}

resource "openstack_dns_recordset_v2" "rs_example_com" {
  zone_id     = "${openstack_dns_zone_v2.example_zone.id}"
  name        = "rs.example.com."
  description = "An example record set"
  ttl         = 3000
  type        = "A"
  records     = ["10.0.0.1"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 DNS client.
    If omitted, the `region` argument of the provider is used.
    Changing this creates a new DNS  record set.

* `zone_id` - (Required) The ID of the zone in which to create the record set.
  Changing this creates a new DNS  record set.

* `name` - (Required) The name of the record set. Note the `.` at the end of the name.
  Changing this creates a new DNS  record set.

* `project_id` - (Optional) The ID of the project DNS zone is created
  for, sets `X-Auth-Sudo-Tenant-ID` header (requires an assigned 
  user role in target project)

* `type` - (Optional) The type of record set. Examples: "A", "MX".
  Changing this creates a new DNS  record set.

* `ttl` - (Optional) The time to live (TTL) of the record set.

* `description` - (Optional) A description of the  record set.

* `records` - (Optional) An array of DNS records. _Note:_ if an IPv6 address
  contains brackets (`[ ]`), the brackets will be stripped and the modified
  address will be recorded in the state.

* `value_specs` - (Optional) Map of additional options. Changing this creates a
  new record set.

* `disable_status_check` - (Optional) Disable wait for recordset to reach ACTIVE
  status. This argumen is disabled by default. If it is set to true, the recordset
  will be considered as created/updated/deleted if OpenStack request returned success.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `type` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `description` - See Argument Reference above.
* `records` - See Argument Reference above.
* `zone_id` - See Argument Reference above.
* `value_specs` - See Argument Reference above.

## Import

This resource can be imported by specifying the zone ID and recordset ID,
separated by a forward slash.

```
$ terraform import openstack_dns_recordset_v2.recordset_1 <zone_id>/<recordset_id>
```
