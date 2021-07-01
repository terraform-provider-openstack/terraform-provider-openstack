---
layout: "openstack"
page_title: "OpenStack: openstack_dns_transfer_request_v2"
sidebar_current: "docs-openstack-resource-dns-transfer-request-v2"
description: |-
  Manages a DNS zone Transfer request in the OpenStack DNS Service
---

# openstack\_dns\_transfer\_request\_v2

Manages a DNS zone transfer request in the OpenStack DNS Service.

## Example Usage

### Automatically detect the correct network

```hcl
resource "openstack_dns_zone_v2" "example_zone" {
  name        = "example.com."
  email       = "jdoe@example.com"
  description = "An example zone"
  ttl         = 3000
  type        = "PRIMARY"
}

resource "openstack_dns_transfer_request_v2" "request_1" {
  zone_id           = "${openstack_dns_zone_v2.example_zone.id}"
  description       = "a transfer request"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    Keypairs are associated with accounts, but a Compute client is needed to
    create one. If omitted, the `region` argument of the provider is used.
    Changing this creates a new DNS zone.

* `zone_id` - (Required) The ID of the zone for which to create the transfer
  request.

* `target_project_id` - (Optional) The target Project ID to transfer to.

* `description` - (Optional) A description of the zone tranfer request.

* `value_specs` - (Optional) Map of additional options. Changing this creates a
  new transfer request.

* `disable_status_check` - (Optional) Disable wait for zone to reach ACTIVE
  status. The check is enabled by default. If this argument is true, zone
  will be considered as created/updated if OpenStack request returned success.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `zone_id` - See Argument Reference above.
* `target_project_id` - See Argument Reference above.
* `description` - See Argument Reference above.
* `value_specs` - See Argument Reference above.

## Import

This resource can be imported by specifying the transferRequest ID:

```
$ terraform import openstack_dns_transfer_request_v2.request_1 <request_id>
```
