---
subcategory: "Identity / Keystone"
layout: "openstack"
page_title: "OpenStack: openstack_identity_mapping_v3"
sidebar_current: "docs-openstack-resource-identity-mapping-v3"
description: |-
  Manages a V3 federation mapping resource within OpenStack Keystone.
---

# openstack\_identity\_mapping\_v3

Manages a V3 federation mapping resource within OpenStack Keystone. Federation
mappings translate remote identity provider attributes into local Keystone
users, groups, and projects.

~> **Note:** You _must_ have admin privileges in your OpenStack cloud and the
OS-FEDERATION extension must be enabled to use this resource.

## Example Usage

```hcl
resource "openstack_identity_group_v3" "federated" {
  name = "federated-users"
}

resource "openstack_identity_mapping_v3" "oidc" {
  mapping_id = "oidc"

  rules = jsonencode([
    {
      local = [
        {
          user = {
            name = "{0}"
            type = "ephemeral"
          }
        },
        {
          group = {
            id = openstack_identity_group_v3.federated.id
          }
        }
      ]
      remote = [
        {
          type = "OIDC-preferred_username"
        }
      ]
    }
  ])
}
```

## Argument Reference

The following arguments are supported:

* `mapping_id` - (Required) The unique ID of the federation mapping. Changing
  this creates a new mapping.

* `rules` - (Required) A JSON array of Keystone federation mapping rules. Use
  `jsonencode` to construct the value from Terraform expressions. Each rule
  contains a `local` array describing the Keystone resources and a `remote`
  array describing the identity provider attributes to match. Changes are
  applied in place.

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
  If omitted, the `region` argument of the provider is used. Changing this
  creates a new mapping.

## Attributes Reference

The following attributes are exported:

* `mapping_id` - See Argument Reference above.
* `rules` - See Argument Reference above.
* `region` - See Argument Reference above.

## Import

Federation mappings can be imported using the `mapping_id`, e.g.

```
$ terraform import openstack_identity_mapping_v3.oidc oidc
```
