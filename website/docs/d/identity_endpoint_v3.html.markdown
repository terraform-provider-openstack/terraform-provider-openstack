---
layout: "openstack"
page_title: "OpenStack: openstack_identity_endpoint_v3"
sidebar_current: "docs-openstack-datasource-identity-endpoint-v3"
description: |-
  Get information on an OpenStack Endpoint.
---

# openstack\_identity\_endpoint_v3

Use this data source to get the ID of an OpenStack endpoint.

Note: This usually requires admin privileges.

## Example Usage

```hcl
data "openstack_identity_endpoint_v3" "endpoint_1" {
  service_name = "demo"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Optional) The service id this endpoint belongs to.

* `service_name` - (Optional) The service name of the endpoint.

* `interface` - (Optional) The endpoint interface. Valid values are `public`,
  `internal`, and `admin`. Default value is `public`

* `region` - (Optional) The region the endpoint is located in.

## Attributes Reference

`id` is set to the ID of the found endpoint. In addition, the following attributes
are exported:

* `service_id` - See Argument Reference above.
* `service_name` - See Argument Reference above.
* `interface` - See Argument Reference above.
* `url` - The endpoint URL
* `region` - The region the endpoint is located in.
