---
subcategory: "Load Balancing as a Service / Octavia"
layout: "openstack"
page_title: "OpenStack: openstack_lb_flavorprofile_v2"
sidebar_current: "docs-openstack-datasource-lb-flavorprofile-v2"
description: |-
  Get information on an OpenStack Load Balancer Flavorprofile.
---

# openstack\_lb\_flavorprofile\_v2

Use this data source to get the ID of an OpenStack Load Balancer flavorprofile.

~> **Note:** This usually requires admin privileges.

## Example Usage

```hcl
data "openstack_lb_flavorprofile_v2" "fp_1" {
  name = "flavorprofile_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Load Balancer client.
  If omitted, the `region` argument of the provider is used.

* `flavorprofile_id` - (Optional) The ID of the flavorprofile. Conflicts with `name` and
  `provider_name`.

* `name` - (Optional) The name of the flavorprofile. Conflicts with `flavorprofile_id`.

* `provider_name` - (Optional) The name of the provider that the flavorprofile uses. Conflicts
  with `flavorprofile_id`.

## Attributes Reference

`id` is set to the ID of the found flavorprofile. In addition, the following attributes
are exported:

* `name` - The name of the flavorprofile.

* `provider_name` - The name of the provider that the flavorprofile uses.

* `flavor_data` - Extra data of the flavorprofile depending on the provider.
