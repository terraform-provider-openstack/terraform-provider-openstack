---
layout: "openstack"
page_title: "OpenStack: openstack_identity_project_v3"
sidebar_current: "docs-openstack-datasource-identity-project-v3"
description: |-
  Get information on an OpenStack Project.
---

# openstack\_identity\_project\_v3

Use this data source to get the ID of an OpenStack project.

## Example Usage

```hcl
data "openstack_identity_project_v3" "project_1" {
  name = "demo"
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Optional) The domain this project belongs to.

* `enabled` - (Optional) Whether the project is enabled or disabled. Valid
  values are `true` and `false`.

* `is_domain` - (Optional) Whether this project is a domain. Valid values
  are `true` and `false`.

* `name` - (Optional) The name of the project.

* `parent_id` - (Optional) The parent of this project.

* `tags` - (Optional) Tags for the project.

## Attributes Reference

`id` is set to the ID of the found project. In addition, the following attributes
are exported:

* `description` - The description of the project.
* `domain_id` - See Argument Reference above.
* `enabled` - See Argument Reference above.
* `is_domain` - See Argument Reference above.
* `name` - See Argument Reference above.
* `parent_id` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `region` - The region the project is located in.
