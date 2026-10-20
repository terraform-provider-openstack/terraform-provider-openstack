---
subcategory: "TaaS / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_taas_tap_mirror_v2"
sidebar_current: "docs-openstack-resource-taas-tap-mirror-v2"
description: |-
  Manages a V2 Neutron Tap Mirror resource within OpenStack tap-as-a-service extension.
---

# openstack\_taas\_tap\_mirror\_v2

Manages a V2 Neutron Tap Mirror resource within OpenStack tap-as-a-service extension.

## Example Usage

```hcl
resource "openstack_taas_tap_mirror_v2" "tap_mirror_1" {
  mirror_type = "erspanv1"
  port_id     = "a25290e9-1a54-4c26-a5b3-34458d122acc"
  remote_ip   = "172.18.1.15"
  directions {
    in  = 1000
    out = 1001
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create an endpoint group. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    group.

* `name` - (Optional) The name of the Tap Mirror. Changing this updates the name of
    the existing Tap Mirror.

* `description` - (Optional) The human-readable description for the Tap Mirror.
    Changing this updates the description of the existing Tap Mirror.

* `tenant_id` - (Optional) The owner of the Tap Mirror. Required if admin wants to
    create a Tap Mirror for another project. Changing this creates a new Tap Mirror.

* `port_id` - (Required) The Port ID of the Tap Mirror, this will be the source of
    the mirrored traffic, and this traffic will be tunneled into the GRE or ERSPAN
    v1 tunnel. The tunnel itself is not starting from this port. Changing this
    creates a new Tap Mirror.

* `mirror_type` - (Required) The type of the mirroring, can be `gre` or `erspanv1`.
    Changing this creates a new Tap Mirror.

* `remote_ip` - (Required) The remote IP of the Tap Mirror, this will be the remote
    end of the GRE or ERSPAN v1 tunnel. Changing this creates a new Tap Mirror.

* `directions` - (Required) A block declaring the directions to be mirrored and their
    identifiers. One block has to be declared with at least one direction. Changing
    this creates a new Tap Mirror.

The `directions` block supports:

* `in` - (Optional) Declares ingress traffic to the port will be mirrored. The value
    is the identifier of the ERSPAN or GRE session between the source and destination,
    this must be unique within the project.

* `out` - (Optional) Declares egress traffic will be mirrored. The value is the
    identifier of the ERSPAN or GRE session between the source and destination,
    this must be unique within the project.

## Attributes Reference

The following attributes are exported:

* `id` - Id of the Tap Mirror.
* `project_id` - Id of the OpenStack project.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `port_id` - See Argument Reference above.
* `mirror_type` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `directions` - See Argument Reference above.

## Import

Tap Mirrors can be imported using the `id`, e.g.

```shell
terraform import openstack_taas_tap_mirror_v2.tap_mirror_1 0837b488-f0e2-4689-99b3-e3ed531f9b10
```
