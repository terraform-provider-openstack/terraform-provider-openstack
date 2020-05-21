---
layout: "openstack"
page_title: "OpenStack: openstack_compute_keypair_v2"
sidebar_current: "docs-openstack-datasource-compute-keypair-v2"
description: |-
  Get information on an OpenStack Keypair.
---

# openstack\_compute\_keypair\_v2

Use this data source to get the ID and public key of an OpenStack keypair.

## Example Usage

```hcl
data "openstack_compute_keypair_v2" "kp" {
  name = "sand"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    If omitted, the `region` argument of the provider is used.

* `name` - (Required) The unique name of the keypair.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `fingerprint` - The fingerprint of the OpenSSH key.
* `public_key` - The OpenSSH-formatted public key of the keypair.
