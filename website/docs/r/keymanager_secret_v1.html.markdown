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
deployments is *not* recommended**.

## Example Usage

```hcl
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
  metadata = {
    key = "foo"
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
 
* `payload` - (Optional) The secret’s data to be stored. **payload_content_type** must also be supplied if **payload** is included.

* `payload_content_type` - (Optional) (required if **payload** is included) The media type for the content of the payload.

* `payload_content_encoding` - (Optional) (required if payload is encoded) The encoding used for the payload to be able to include it in the JSON request. Currently only base64 is supported.

* `metadata` - (Optional) Additional Metadata for the secret.
			
## Attributes Reference

The following attributes are exported:

* `secret_ref` - The secret reference / where to find the secret.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `bit_length` - See Argument Reference above.
* `algorithm` - See Argument Reference above.
* `mode` - See Argument Reference above.
* `secret_type` - See Argument Reference above.
* `creator_id` - The creator of the secret.
* `status` - The status of the secret.

## Import

Secrets can be imported using the secret id (the last part of the secret reference), e.g.:

```
$ terraform import openstack_keymanager_secret_v1.secret_1 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
