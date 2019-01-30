---
layout: "openstack"
page_title: "OpenStack: openstack_keymanager_secret_metadata_v1"
sidebar_current: "docs-openstack-resource-keymanager-secret-metadata-v1"
description: |-
  Manages a V1 Barbican secret metadata resource within OpenStack.
---

# openstack\_keymanager\_secret\_metadata\_v1

Manages a V1 Barbican secret metadata resource within OpenStack.

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
	}

resource "openstack_keymanager_secret_metadata_v1" "metadata_1" {
		secret_ref = "${openstack_keymanager_secret_v1.secret_1.secret_ref}"
		metadata {
			foo = "bar"
		}
	}
	
```

## Argument Reference

The following arguments are supported:

* `secret_ref` - (Required) The secret reference of the secret that this metadata belongs to.

* `metadata` - (Optional) A string -> string mapping of metadata values.

## Attributes Reference

The following attributes are exported:

* `metadata` - See Argument Reference above.

