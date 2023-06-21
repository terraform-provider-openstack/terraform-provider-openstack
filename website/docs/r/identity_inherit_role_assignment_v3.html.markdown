---
subcategory: "Identity / Keystone"
layout: "openstack"
page_title: "OpenStack: openstack_identity_inherit_role_assignment_v3"
sidebar_current: "docs-openstack-resource-identity-inherit-role-assignment-v3"
description: |-
  Manages a V3 Inherit Role assignment within OpenStack Keystone.
---

# openstack\_identity\_inherit\_role\_assignment\_v3

Manages a V3 Inherit Role assignment within OpenStack Keystone. This uses the
Openstack keystone `OS-INHERIT` api to created inherit roles within domains
and parent projects for users and groups.

~> **Note:** You _must_ have admin privileges in your OpenStack cloud to use
this resource.

## Example Usage

```hcl
resource "openstack_identity_user_v3" "user_1" {
  name = "user_1"
  domain_id = "default"
}

resource "openstack_identity_role_v3" "role_1" {
  name = "role_1"
  domain_id = "default"
}

resource "openstack_identity_inherit_role_assignment_v3" "role_assignment_1" {
  user_id = openstack_identity_user_v3.user_1.id
  domain_id = "default"
  role_id = openstack_identity_role_v3.role_1.id
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Optional; Required if `project_id` is empty) The domain to assign the role in.

* `group_id` - (Optional; Required if `user_id` is empty) The group to assign the role to.

* `project_id` - (Optional; Required if `domain_id` is empty) The project to assign the role in.
  The project should be able to containt child projects.

* `user_id` - (Optional; Required if `group_id` is empty) The user to assign the role to.

* `role_id` - (Required) The role to assign.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `group_id` - See Argument Reference above.
* `user_id` - See Argument Reference above.
* `role_id` - See Argument Reference above.

## Import

Inherit role assignments can be imported using a constructed id. The id should 
have the form of `domainID/projectID/groupID/userID/roleID`. When something is
not used then leave blank.

For example this will import the inherit role assignment for: 
projectID: 014395cd-89fc-4c9b-96b7-13d1ee79dad2,
userID: 4142e64b-1b35-44a0-9b1e-5affc7af1106,
roleID: ea257959-eeb1-4c10-8d33-26f0409a755d
( domainID and groupID are left blank)

```
$ terraform import openstack_identity_inherit_role_assignment_v3.role_assignment_1 /014395cd-89fc-4c9b-96b7-13d1ee79dad2//4142e64b-1b35-44a0-9b1e-5affc7af1106/ea257959-eeb1-4c10-8d33-26f0409a755d
```
