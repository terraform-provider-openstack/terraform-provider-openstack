---
layout: "openstack"
page_title: "OpenStack: openstack_identity_service_v3"
sidebar_current: "docs-openstack-resource-identity-service-v3"
description: |-
  Manages a V3 Service resource within OpenStack Keystone.
---

# openstack\_identity\_service\_v3

Manages a V3 Service resource within OpenStack Keystone.

~> **Note:** This usually requires admin privileges.

## Example Usage

```hcl
resource "openstack_identity_service_v3" "service_1" {
  name = "custom"
  type = "custom"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
  If omitted, the `region` argument of the provider is used.

* `name` - (Required) The service name.

* `description` - (Optional) The service description.

* `type` - (Required) The service type.

* `enabled` - (Optional) The service status. Defaults to `true`.

## Attributes Reference

`id` is set to the ID of the found service. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `type` - See Argument Reference above.
* `enabled` - See Argument Reference above.
* `description` - See Argument Reference above.

## Import

Services can be imported using the `id`, e.g.

```
$ terraform import openstack_identity_service_v3.service_1 6688e967-158a-496f-a224-cae3414e6b61
```
