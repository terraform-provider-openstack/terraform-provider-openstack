---
layout: "openstack"
page_title: "OpenStack: openstack_networking_qos_policy_v2"
sidebar_current: "docs-openstack-resource-networking-qos-policy-v2"
description: |-
  Manages a V2 Neutron QoS policy resource within OpenStack.
---

# openstack\_networking\_qos\_policy\_v2

Manages a V2 Neutron QoS policy resource within OpenStack.

## Example Usage

### Create a QoS Policy

```hcl
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name        = "qos_policy_1"
  description = "bw_limit"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a Neutron Qos policy. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    QoS policy.

* `name` - (Required) The name of the QoS policy. Changing this updates the name of
    the existing QoS policy.

* `project_id` - (Optional) The owner of the QoS policy. Required if admin wants to
    create a QoS policy for another project. Changing this creates a new QoS policy.

* `shared` - (Optional) Indicates whether this QoS policy is shared across
    all projects. Changing this updates the shared status of the existing
    QoS policy.

* `description` - (Optional) The human-readable description for the QoS policy.
    Changing this updates the description of the existing QoS policy.

* `is_default` - (Optional) Indicates whether the QoS policy is default
    QoS policy or not. Changing this updates the default status of the existing
    QoS policy.

* `value_specs` - (Optional) Map of additional options.

* `tags` - (Optional) A set of string tags for the QoS policy.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `created_at` - The time at which QoS policy was created.
* `updated_at` - The time at which QoS policy was created.
* `shared` - See Argument Reference above.
* `description` - See Argument Reference above.
* `is_default` - See Argument Reference above.
* `revision_number` - The revision number of the QoS policy.
* `value_specs` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `all_tags` - The collection of tags assigned on the QoS policy, which have been
  explicitly and implicitly added.

## Import

QoS Policies can be imported using the `id`, e.g.

```
$ terraform import openstack_networking_qos_policy_v2.qos_policy_1 d6ae28ce-fcb5-4180-aa62-d260a27e09ae
```
