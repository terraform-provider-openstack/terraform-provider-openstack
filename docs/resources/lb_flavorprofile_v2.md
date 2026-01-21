---
subcategory: "Load Balancing as a Service / Octavia"
layout: "openstack"
page_title: "OpenStack: openstack_lb_flavorprofile_v2"
sidebar_current: "docs-openstack-resource-lb-flavorprofile-v2"
description: |-
  Manages a V2 flavorprofile resource within OpenStack.
---

# openstack\_lb\_flavorprofile\_v2

Manages a V2 load balancer flavorprofile resource within OpenStack.

~> **Note:** This usually requires admin privileges.

## Example Usage

### Using jsonencode

```hcl
resource "openstack_lb_flavorprofile_v2" "flavorprofile_1" {
  name          = "amphora-single-profile"
  provider_name = "amphora"
  flavor_data   = jsonencode({
    "loadbalancer_topology" : "SINGLE",
  })
}
```

### Using plain string

```hcl
resource "openstack_lb_flavorprofile_v2" "flavorprofile_1" {
  name          = "amphora-single-profile"
  provider_name = "amphora"
  flavor_data   = "{\"loadbalancer_topology\": \"SINGLE\"}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create an LB member. If omitted, the
  `region` argument of the provider is used. Changing this creates a new
  LB flavorprofile.

* `name` - (Required) Name of the flavorprofile. Changing this updates the existing
  flavorprofile.

* `provider_name` - (Required) The provider_name that the flavor_profile will use.
  Changing this updates the existing flavorprofile.

* `flavor_data` - (Required) String that passes the flavor_data for the flavorprofile.
  The data that are allowed depend on the `provider_name` that is passed. [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode)
  can be used for readability as shown in the example [above](#using-jsonencode).
  Changing this updates the existing flavorprofile.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `provider_name` - See Argument Reference above.
* `flavor_data` - See Argument Reference above.

## Import

flavorprofiles can be imported using their `id`. Example:

```shell
terraform import openstack_lb_flavorprofile_v2.flavorprofile_1 2a0f2240-c5e6-41de-896d-e80d97428d6b
```
