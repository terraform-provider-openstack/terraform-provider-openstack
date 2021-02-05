---
layout: "openstack"
page_title: "OpenStack: openstack_vpnaas_service_v2"
sidebar_current: "docs-openstack-resource-vpnaas-service-v2"
description: |-
  Manages a V2 Neutron VPN service resource within OpenStack.
---

# openstack\_vpnaas\_service\_v2

Manages a V2 Neutron VPN service resource within OpenStack.

## Example Usage

```hcl
resource "openstack_vpnaas_service_v2" "service_1" {
  name           = "my_service"
  router_id      = "14a75700-fc03-4602-9294-26ee44f366b3"
  admin_state_up = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a VPN service. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    service.

* `name` - (Optional) The name of the service. Changing this updates the name of
    the existing service.

* `tenant_id` - (Optional) The owner of the service. Required if admin wants to
    create a service for another project. Changing this creates a new service.

* `description` - (Optional) The human-readable description for the service.
    Changing this updates the description of the existing service.

* `admin_state_up` - (Optional) The administrative state of the resource. Can either be up(true) or down(false).
    Changing this updates the administrative state of the existing service.

* `subnet_id` - (Optional) SubnetID is the ID of the subnet. Default is null.

* `router_id` - (Required) The ID of the router. Changing this creates a new service.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `router_id` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `status` - Indicates whether IPsec VPN service is currently operational. Values are ACTIVE, DOWN, BUILD, ERROR, PENDING_CREATE, PENDING_UPDATE, or PENDING_DELETE.
* `external_v6_ip` - The read-only external (public) IPv6 address that is used for the VPN service.
* `external_v4_ip` - The read-only external (public) IPv4 address that is used for the VPN service.
* `description` - See Argument Reference above.
* `value_specs` - See Argument Reference above.

## Import

Services can be imported using the `id`, e.g.

```
$ terraform import openstack_vpnaas_service_v2.service_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
