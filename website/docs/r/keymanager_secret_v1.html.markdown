---
layout: "openstack"
page_title: "OpenStack: openstack_keymanager_secret_v1"
sidebar_current: "docs-openstack-resource-keymanager-secret-v1"
description: |-
  Manages a V1 Barbican secret resource within OpenStack.
---

# openstack\_keymanager\_secret\_v1

Manages a V1 Barbican secret resource within OpenStack.

~> **Important Security Notice** The payload of this resource will be stored
*unencrypted* in your Terraform state file. **Use of this resource for production
deployments is *not* recommended**. [Read more about sensitive data in
state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage

### Simple secret

```hcl
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm            = "aes"
  bit_length           = 256
  mode                 = "cbc"
  name                 = "mysecret"
  payload              = "foobar"
  payload_content_type = "text/plain"
  secret_type          = "passphrase"

  metadata = {
    key = "foo"
  }
}
```

### Secret with whitespaces

~> **Note** If you want to store payload with leading or trailing whitespaces,
it's recommended to store it in a base64 encoding. Plain text payload can also
work, but further addind or removing of the leading or trailing whitespaces
won't be detected as a state change, e.g. changing plain text payload from
`password ` to `password` won't recreate the secret.

```hcl
resource "openstack_keymanager_secret_v1" "secret_1" {
  name                     = "password"
  payload                  = "${base64encode("password with the whitespace at the end ")}"
  secret_type              = "passphrase"
  payload_content_type     = "application/octet-stream"
  payload_content_encoding = "base64"
}
```

### Secret with the expiration date

```hcl
resource "openstack_keymanager_secret_v1" "secret_1" {
  name                 = "certificate"
  payload              = "${file("certificate.pem")}"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
  expiration           = "${timeadd(timestamp(), format("%dh", 8760))}" # one year in hours

  lifecycle {
    ignore_changes = [
      expiration
    ]
  }
}
```

### Secret with the ACL

~> **Note** Only read ACLs are supported

```hcl
resource "openstack_keymanager_secret_v1" "secret_1" {
  name                 = "certificate"
  payload              = "${file("certificate.pem")}"
  secret_type          = "certificate"
  payload_content_type = "text/plain"

  acl {
    read {
      project_access = false
      users = [
        "userid1",
        "userid2",
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 KeyManager client.
    A KeyManager client is needed to create a secret. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    V1 secret.

* `name` - (Optional) Human-readable name for the Secret. Does not have
    to be unique.
    
* `bit_length` - (Optional) Metadata provided by a user or system for informational purposes.

* `algorithm` - (Optional) Metadata provided by a user or system for informational purposes.

* `mode` - (Optional) Metadata provided by a user or system for informational purposes.

* `secret_type` - (Optional) Used to indicate the type of secret being stored. For more information see [Secret types](https://docs.openstack.org/barbican/latest/api/reference/secret_types.html).
 
* `payload` - (Optional) The secret's data to be stored. **payload\_content\_type** must also be supplied if **payload** is included.

* `payload_content_type` - (Optional) (required if **payload** is included) The media type for the content of the payload. Must be one of `text/plain`, `text/plain;charset=utf-8`, `text/plain; charset=utf-8`, `application/octet-stream`, `application/pkcs8`.

* `payload_content_encoding` - (Optional) (required if **payload** is encoded) The encoding used for the payload to be able to include it in the JSON request. Must be either `base64` or `binary`.

* `expiration` - (Optional) The expiration time of the secret in the RFC3339 timestamp format (e.g. `2019-03-09T12:58:49Z`). If omitted, a secret will never expire. Changing this creates a new secret.

* `metadata` - (Optional) Additional Metadata for the secret.

* `acl` - (Optional) Allows to control an access to a secret. Currently only the
  `read` operation is supported. If not specified, the secret is accessible
  project wide.

The `acl` `read` block supports:

* `project_access` - (Optional) Whether the secret is accessible project wide.
  Defaults to `true`.

* `users` - (Optional) The list of user IDs, which are allowed to access the
  secret, when `project_access` is set to `false`.

* `created_at` - (Computed) The date the secret ACL was created.

* `updated_at` - (Computed) The date the secret ACL was last updated.

## Attributes Reference

The following attributes are exported:

* `secret_ref` - The secret reference / where to find the secret.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `bit_length` - See Argument Reference above.
* `algorithm` - See Argument Reference above.
* `mode` - See Argument Reference above.
* `secret_type` - See Argument Reference above.
* `payload` - See Argument Reference above.
* `payload_content_type` - See Argument Reference above.
* `acl` - See Argument Reference above.
* `payload_content_encoding` - See Argument Reference above.
* `expiration` - See Argument Reference above.
* `content_types` - The map of the content types, assigned on the secret.
* `creator_id` - The creator of the secret.
* `status` - The status of the secret.
* `created_at` - The date the secret was created.
* `updated_at` - The date the secret was last updated.
* `all_metadata` - The map of metadata, assigned on the secret, which has been
  explicitly and implicitly added.

## Import

Secrets can be imported using the secret id (the last part of the secret reference), e.g.:

```
$ terraform import openstack_keymanager_secret_v1.secret_1 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
