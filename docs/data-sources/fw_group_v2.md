---
subcategory: "FWaaS / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_fw_group_v2"
sidebar_current: "docs-openstack-datasource-fw-group-v2"
description: |-
  Get information on an OpenStack Firewall Group V2.
---

# openstack\_fw\_group\_v2

Use this data source to get information of an available OpenStack firewall group v2.

## Example Usage

```hcl
data "openstack_fw_group_v2" "group" {
  name = "tf_test_group"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
    A Neutron client is needed to retrieve firewall group ids. If omitted, the
    `region` argument of the provider is used.

* `name` - (Optional) The name of the firewall group.

* `description` - (Optional) Human-readable description of the firewall group.

* `group_id` - (Optional) The ID of the firewall group.

* `tenant_id` - (Optional) - This argument conflicts and is interchangeable
    with `project_id`. The owner of the firewall group.

* `project_id` - (Optional) - This argument conflicts and is interchangeable
    with `tenant_id`. The owner of the firewall group.

* `shared` - (Optional) The sharing status of the firewall group.

* `admin_state_up` - (Optional) Administrative up/down status for the firewall group.

* `ingress_firewall_policy_id` - (Optional) The ingress policy ID of the firewall group.

* `egress_firewall_policy_id` - (Optional) The egress policy ID of the firewall group.

* `status` - (Optional) Enabled status for the firewall group.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `group_id` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `shared` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `ingress_firewall_policy_id` - See Argument Reference above.
* `egress_firewall_policy_id` - See Argument Reference above.
* `ports` - Ports associated with the firewall group.
* `status` - See Argument Reference above.
