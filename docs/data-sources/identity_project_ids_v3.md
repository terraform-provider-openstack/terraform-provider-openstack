---
subcategory: "Identity / Keystone"
layout: "openstack"
page_title: "OpenStack: openstack_identity_project_ids_v3"
sidebar_current: "docs-openstack-datasource-identity-project-ids-v3"
description: |-
  Provides a list of OpenStack Project IDs.
---

# openstack\_identity\_project\_ids\_v3

Use this data source to get a list of OpenStack Project IDs matching the
specified criteria.

~> **Note:** You _must_ have domain admin or cloud admin privileges in your OpenStack cloud to use
this datasource.

## Example Usage

```hcl
data "openstack_identity_project_ids_v3" "projects" {
  name_regex = "^prod.*"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
  If omitted, the `region` argument of the provider is used.

* `name` - (Optional) The name of the project. Cannot be used simultaneously with
  `name_regex`.

* `name_regex` - (Optional) The regular expression of the name of the project.
  Cannot be used simultaneously with `name`. Unlike filtering by `name` the
  `name_regex` filtering does by client on the result of OpenStack search
  query.

* `domain_id` - (Optional) The domain projects belongs to.

* `enabled` - (Optional) Whether the project is enabled or disabled. Valid
  values are `true` and `false`. Default is `true`.

* `name` - (Optional) The name of the project.

* `parent_id` - (Optional) The parent of the project.

* `tags` - (Optional) Tags for the project.

## Attributes Reference

`ids` is set to the list of Openstack Project IDs.
