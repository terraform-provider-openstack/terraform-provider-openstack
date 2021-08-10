---
layout: "openstack"
page_title: "OpenStack: openstack_identity_domain_v3"
sidebar_current: "docs-openstack-datasource-identity-domain-v3"
description: |-
  Get information on an OpenStack Domain.
---

# openstack\_identity\_domain\_v3

Use this data source to get the ID of an OpenStack domain.

~> **Note:** This usually requires admin privileges.

## Example Usage

```hcl
data "openstack_identity_domain_v3" "Default" {
  name = "Default"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
  If omitted, the `region` argument of the provider is used.

* `name` - (Optional) The domain name.

* `enabled` - (Optional) The domain status.

## Attributes Reference

`id` is set to the ID of the found domain. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `enabled` - See Argument Reference above.
