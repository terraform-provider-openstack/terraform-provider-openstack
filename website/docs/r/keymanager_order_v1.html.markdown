---
layout: "openstack"
page_title: "OpenStack: openstack_keymanager_order_v1"
sidebar_current: "docs-openstack-resource-keymanager-order-v1"
description: |-
  Manages a V1 Barbican order resource within OpenStack.
---

# openstack\_keymanager\_order\_v1

Manages a V1 Barbican order resource within OpenStack.

## Example Usage

### Basic order

```hcl
resource "openstack_keymanager_order_v1" "order_1" {
  type                 = "key"
  meta                 = {
    algorithm            = "aes"
    bit_length           = 256
    name                 = "mysecret"
    mode                 = "cbc"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 KeyManager client.
    A KeyManager client is needed to create a order. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    V1 order.

* `type` - (Required) The type of key to be generated. Must be one of `asymmetric`, `key`.

* `meta` - (Required) Dictionary containing the order metadata used to generate the order. The structure is described below.

The `meta` block supports:

* `algorithm` - Algorithm to use for key generation.

* `bit_length` - Bit lenght of key to be generated.

* `mode` - The mode to use for key generation.

* `name` - Name of the key to generate.

* `payload_content_type` - The media type for the content of the payload. Must be one of `text/plain`, `text/plain;charset=utf-8`, `text/plain; charset=utf-8`, `application/octet-stream`, `application/pkcs8`.

## Attributes Reference

The following attributes are exported:

* `created` - The date the order was created.
* `creator_id` - The creator of the order.
* `meta` - See Argument Reference above.
* `order_ref` - The order reference / where to find the order.
* `region` - See Argument Reference above.
* `secret_ref` - The secret reference / where to find the secret.
* `status` - The status of the order.
* `sub_status` - The sub status of the order.
* `sub_status_message` - The sub status message of the order.
* `type` - The type of the order.
* `updated` - The date the order was last updated.

## Import

Orders can be imported using the order id (the last part of the order reference), e.g.:

```
$ terraform import openstack_keymanager_order_v1.order_1 0c6cd26a-c012-4d7b-8034-057c0f1c2953

