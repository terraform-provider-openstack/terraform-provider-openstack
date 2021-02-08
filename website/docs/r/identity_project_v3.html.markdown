---
layout: "openstack"
page_title: "OpenStack: openstack_identity_project_v3"
sidebar_current: "docs-openstack-resource-identity-project-v3"
description: |-
  Manages a V3 Project resource within OpenStack Keystone.
---

# openstack\_identity\_project\_v3

Manages a V3 Project resource within OpenStack Keystone.

~> **Note:** You _must_ have admin privileges in your OpenStack cloud to use
this resource.

## Example Usage

```hcl
resource "openstack_identity_project_v3" "project_1" {
  name        = "project_1"
  description = "A project"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description of the project.

* `domain_id` - (Optional) The domain this project belongs to.

* `enabled` - (Optional) Whether the project is enabled or disabled. Valid
  values are `true` and `false`. Default is `true`.

* `is_domain` - (Optional) Whether this project is a domain. Valid values
  are `true` and `false`. Default is `false`. Changing this creates a new
  project/domain.

* `name` - (Optional) The name of the project.

* `parent_id` - (Optional) The parent of this project. Changing this creates
  a new project.

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new project.

* `tags` - (Optional) Tags for the project. Changing this updates the existing
    project.

## Attributes Reference

The following attributes are exported:

* `description` - The description of the project.
* `domain_id` - See Argument Reference above.
* `enabled` - See Argument Reference above.
* `is_domain` - See Argument Reference above.
* `name` - See Argument Reference above.
* `parent_id` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `region` - See Argument Reference above.

## Import

Projects can be imported using the `id`, e.g.

```
$ terraform import openstack_identity_project_v3.project_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
