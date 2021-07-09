---
layout: "openstack"
page_title: "OpenStack: openstack_dns_transfer_accept_v2"
sidebar_current: "docs-openstack-resource-dns-transfer-accept-v2"
description: |-
  Manages a DNS zone Transfer accept in the OpenStack DNS Service
---

# openstack\_dns\_transfer\_accept\_v2

Manages a DNS zone transfer accept in the OpenStack DNS Service.

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
  description       = "a transfer accept"
}

resource "openstack_dns_transfer_accept_v2" "accept_1" {
  zone_transfer_request_id = "${openstack_dns_transfer_request_v2.request_1.id}"
  key                      = "${openstack_dns_transfer_request_v2.request_1.key}"
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    Keypairs are associated with accounts, but a Compute client is needed to
    create one. If omitted, the `region` argument of the provider is used.
    Changing this creates a new DNS zone.

* `zone_transfer_request_id` - (Required) The ID of the zone transfer request.

* `key` - (Required) The transfer key.

* `value_specs` - (Optional) Map of additional options. Changing this creates a
  new transfer accept.

* `disable_status_check` - (Optional) Disable wait for zone to reach ACTIVE
  status. The check is enabled by default. If this argument is true, zone
  will be considered as created/updated if OpenStack accept returned success.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `zone_transfer_request_id` - See Argument Reference above.
* `key` - See Argument Reference above.
* `value_specs` - See Argument Reference above.

## Import

This resource can be imported by specifying the transferAccept ID:

```
$ terraform import openstack_dns_transfer_accept_v2.accept_1 <accept_id>
```
