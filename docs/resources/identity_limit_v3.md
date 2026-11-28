---
subcategory: "Identity / Keystone"
layout: "openstack"
page_title: "OpenStack: openstack_identity_limit_v3"
sidebar_current: "docs-openstack-resource-identity-limit-v3"
description: |-
  Manages a V3 Limit resource within OpenStack Keystone.
---

# openstack\_identity\_limit\_v3

Manages a V3 Limit resource within OpenStack Keystone.

~> **Note:** You _must_ have admin privileges in your OpenStack cloud to use
this resource.

## Example Usage

```hcl
data "openstack_identity_service_v3" "glance" {
  name = "glance"
}

resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_identity_limit_v3" "limit_1" {
  project_id     = openstack_identity_project_v3.project_1.id
  service_id     = data.openstack_identity_service_v3.glance.id
  resource_name  = "image_count_total"
  resource_limit = 10
  description    = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
  If omitted, the `region` argument of the provider is used. Changing this
  creates a new Limit.

* `project_id` - (Optional) The project the limit applies to. Changing this
  creates a new Limit.

* `domain_id` - (Optional) The domain the limit applies to. Changing this
  creates a new Limit.

* `service_id` - (Required) The service the limit applies to. Changing this
  creates a new Limit.

* `resource_name` - (Required) The resource that the limit applies to. Changing
  this creates a new Limit.

* `resource_limit` - (Required) Integer for the actual limit.

* `description` - (Optional) Description of the limit.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the limit
* `project_id` - See Argument Reference above.
* `domain_id` - See Argument Reference above.
* `region` - See Argument Reference above.
* `service_id` - See Argument Reference above.
* `resource_name` - See Argument Reference above.
* `resource_limit` - See Argument Reference above.
* `description` - See Argument Reference above.

## Import

Limits can be imported using the `id`, e.g.

```shell
terraform import openstack_identity_limit_v3.limit_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
