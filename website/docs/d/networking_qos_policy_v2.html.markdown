---
layout: "openstack"
page_title: "OpenStack: openstack_networking_qos_policy_v2"
sidebar_current: "docs-openstack-datasource-networking-qos-policy-v2"
description: |-
  Get information on an OpenStack QoS Policy.
---

# openstack\_networking\_qos\_policy\_v2

Use this data source to get the ID of an available OpenStack QoS policy.

## Example Usage

```hcl
data "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to retrieve a QoS policy ID. If omitted, the
    `region` argument of the provider is used.

* `name` - (Optional) The name of the QoS policy.

* `project_id` - (Optional) The owner of the QoS policy.

* `shared` - (Optional) Whether this QoS policy is shared across all projects.

* `description` - (Optional) The human-readable description for the QoS policy.

* `is_default` - (Optional) Whether the QoS policy is default policy or not.

* `tags` - (Optional) The list of QoS policy tags to filter.

## Attributes Reference

`id` is set to the ID of the found QoS policy. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `created_at` -  The time at which QoS policy was created.
* `updated_at` - The time at which QoS policy was created.
* `shared` - See Argument Reference above.
* `description` - See Argument Reference above.
* `is_default` - See Argument Reference above.
* `revision_number` - The revision number of the QoS policy.
* `all_tags` - The set of string tags applied on the QoS policy.
