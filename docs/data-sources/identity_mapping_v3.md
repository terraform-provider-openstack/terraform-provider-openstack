---
subcategory: "Identity / Keystone"
layout: "openstack"
page_title: "OpenStack: openstack_identity_mapping_v3"
sidebar_current: "docs-openstack-datasource-identity-mapping-v3"
description: |-
  Retrieves a V3 federation mapping from OpenStack Keystone.
---

# openstack\_identity\_mapping\_v3

Use this data source to retrieve a V3 federation mapping from OpenStack
Keystone.

~> **Note:** You _must_ have admin privileges in your OpenStack cloud and the
OS-FEDERATION extension must be enabled to use this data source.

## Example Usage

```hcl
data "openstack_identity_mapping_v3" "oidc" {
  mapping_id = "oidc"
}
```

## Argument Reference

The following arguments are supported:

* `mapping_id` - (Required) The unique ID of the federation mapping.

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
  If omitted, the `region` argument of the provider is used.

## Attributes Reference

The following attributes are exported:

* `rules` - A JSON array containing the Keystone federation mapping rules.
* `mapping_id` - See Argument Reference above.
* `region` - See Argument Reference above.
