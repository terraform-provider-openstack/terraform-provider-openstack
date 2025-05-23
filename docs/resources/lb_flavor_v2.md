---
subcategory: "Load Balancing as a Service / Octavia"
layout: "openstack"
page_title: "OpenStack: openstack_lb_flavor_v2"
sidebar_current: "docs-openstack-resource-lb-flavorp-v2"
description: |-
  Manages a V2 flavor resource within OpenStack.
---

# openstack\_lb\_flavor\_v2

Manages a V2 load balancer flavor resource within OpenStack.

~> **Note:** This usually requires admin privileges.

## Example Usage 


```hcl
resource "openstack_lb_flavorprofile_v2" "fp_1" {
	name          = "test"
	provider_name = "amphora"
	flavor_data   = jsonencode({
	  "loadbalancer_topology": "ACTIVE_STANDBY",
	})
}

resource "openstack_lb_flavor_v2" "flavor_1" {
	name              = "test"
	description       = "This is a test flavor"
	flavor_profile_id = openstack_lb_flavorprofile_v2.fp_1.id	 
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the flavor. Changing this updates the existing
  flavor.

* `flavor_profile_id` - (Required) The flavor_profile_id that the flavor
  will use. Changing this creates a new flavor.

* `description` - (Optional) The description of the flavor. Changing this 
  updates the existing flavor.

* `enabled` - (Optional) Whether the flavor is enabled or not. Defaults to `true`.
  Changing this updates the existing flavor. 

## Attributes Reference

The following attributes are exported:

* `id` - The id of the flavor
* `name` - See Argument Reference above.
* `flavor_profile_id` - See Argument Reference above.
* `description` - See Argument Reference above.
* `enabled` -  See Argument Reference above.

## Import

flavors can be imported using their `id`. Example:
```
$ terraform import openstack_lb_flavor_v2.flavor_1 2a0f2240-c5e6-41de-896d-e80d97428d6b
```
