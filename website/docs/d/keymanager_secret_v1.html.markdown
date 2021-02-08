---
layout: "openstack"
page_title: "OpenStack: openstack_keymanager_secret_v1"
sidebar_current: "docs-openstack-datasource-keymanager-secret-v1"
description: |-
  Get information on a V1 Barbican secret resource within OpenStack.
---

# openstack\_keymanager\_secret\_v1

Use this data source to get the ID and the payload of an available Barbican
secret

~> **Important Security Notice** The payload of this data source will be stored
*unencrypted* in your Terraform state file. **Use of this resource for
production deployments is *not* recommended**. [Read more about sensitive data
in state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage

```hcl
data "openstack_keymanager_secret_v1" "example" {
  mode        = "cbc"
  secret_type = "passphrase"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 KeyManager client.
  A KeyManager client is needed to fetch a secret. If omitted, the `region`
  argument of the provider is used.

* `name` - (Optional) The Secret name.

* `bit_length` - (Optional) The Secret bit length.

* `algorithm` - (Optional) The Secret algorithm.

* `mode` - (Optional) The Secret mode.

* `secret_type` - (Optional) The Secret type. For more information see
  [Secret types](https://docs.openstack.org/barbican/latest/api/reference/secret_types.html).

* `acl_only` - (Optional) Select the Secret with an ACL that contains the user.
  Project scope is ignored. Defaults to `false`.

* `expiration_filter` - (Optional) Date filter to select the Secret with
  expiration matching the specified criteria. See Date Filters below for more
  detail.

* `created_at_filter` - (Optional) Date filter to select the Secret with
  created matching the specified criteria. See Date Filters below for more
  detail.

* `updated_at_filter` - (Optional) Date filter to select the Secret with
  updated matching the specified criteria. See Date Filters below for more
  detail.

## Date Filters

The values for the `expiration_filter`, `created_at_filter`, and
`updated_at_filter` parameters are comma-separated lists of time stamps in
RFC3339 format. The time stamps can be prefixed with any of these comparison
operators: *gt:* (greater-than), *gte:* (greater-than-or-equal), *lt:*
(less-than), *lte:* (less-than-or-equal).

For example, to get a passphrase a Secret with CBC moda, that will expire in
January of 2020:

```hcl
data "openstack_keymanager_secret_v1" "date_filter_example" {
  mode              = "cbc"
  secret_type       = "passphrase"
  expiration_filter = "gt:2020-01-01T00:00:00Z"
}
```

## Attributes Reference

The following attributes are exported:

* `secret_ref` - The secret reference / where to find the secret.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `bit_length` - See Argument Reference above.
* `algorithm` - See Argument Reference above.
* `mode` - See Argument Reference above.
* `secret_type` - See Argument Reference above.
* `acl_only` - See Argument Reference above.
* `expiration_filter` - See Argument Reference above.
* `created_at_filter` - See Argument Reference above.
* `updated_at_filter` - See Argument Reference above.
* `payload` - The secret payload.
* `payload_content_type` - The Secret content type.
* `payload_content_encoding` - The Secret encoding.
* `content_types` - The map of the content types, assigned on the secret.
* `creator_id` - The creator of the secret.
* `status` - The status of the secret.
* `expiration` - The date the secret will expire.
* `created_at` - The date the secret was created.
* `updated_at` - The date the secret was last updated.
* `metadata` - The map of metadata, assigned on the secret, which has been
  explicitly and implicitly added.
* `acl` - The list of ACLs assigned to a secret. The `read` structure is described below.

The `acl` `read` attribute supports:

* `project_access` - Whether the secret is accessible project wide.

* `users` - The list of user IDs, which are allowed to access the secret, when
  `project_access` is set to `false`.

* `created_at` - The date the secret ACL was created.

* `updated_at` - The date the secret ACL was last updated.
