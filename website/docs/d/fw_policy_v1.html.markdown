---
layout: "openstack"
page_title: "OpenStack: openstack_fw_policy_v1"
sidebar_current: "docs-openstack-datasource-fw-policy-v1"
description: |-
  Get information on an OpenStack Firewall Policy.
---

# openstack\_fw\_policy\_v1

Use this data source to get firewall policy information of an available OpenStack firewall policy.

## Example Usage

```hcl
data "openstack_fw_policy_v1" "policy" {
  name = "tf_test_policy"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve firewall policy ids. If omitted, the
  `region` argument of the provider is used.

* `policy_id` - (Optional) The ID of the firewall policy.

* `name` - (Optional) The name of the firewall policy.

* `tenant_id` - (Optional) The owner of the firewall policy.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `policy_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `description` - The description of the firewall policy.
* `audited` - The audit status of the firewall policy.
* `shared` - The sharing status of the firewall policy.
* `rules` - The array of one or more firewall rules that comprise the policy.
