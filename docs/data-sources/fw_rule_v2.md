---
subcategory: "FWaaS / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_fw_rule_v2"
sidebar_current: "docs-openstack-datasource-fw-rule-v2"
description: |-
  Get information on an OpenStack Firewall Rule V2.
---

# openstack\_fw\_rule\_v2

Use this data source to get information of an available OpenStack firewall rule v2.

## Example Usage

```hcl
data "openstack_fw_rule_v2" "rule" {
  name = "tf_test_rule"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve firewall policy ids. If omitted, the
  `region` argument of the provider is used.

* `rule_id` - (Optional) The ID of the firewall rule.

* `name` - (Optional) The name of the firewall rule.

* `description` - (Optional) The description of the firewall rule.

* `tenant_id` - (Optional) - This argument conflicts and is interchangeable
    with `project_id`. The owner of the firewall rule.

* `project_id` - (Optional) - This argument conflicts and is interchangeable
    with `tenant_id`. The owner of the firewall rule.

* `protocol` - (Optional) The protocol type on which the firewall rule operates.

* `action` - (Optional) Action to be taken when the firewall rule matches.

* `ip_version` - (Optional) IP version, either 4 (default) or 6.

* `source_ip_address` - (Optional) The source IP address on which the firewall
    rule operates.

* `source_port` - (Optional) The source port on which the firewall
    rule operates.

* `destination_ip_address` - (Optional) The destination IP address on which the
    firewall rule operates.

* `destination_port` - (Optional) The destination port on which the firewall
    rule operates.

* `shared` - The sharing status of the firewall policy.

* `enabled` - Enabled status for the firewall rule.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `rule_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `action` - See Argument Reference above.
* `ip_version` - See Argument Reference above.
* `source_ip_address` - See Argument Reference above.
* `source_port` - See Argument Reference above.
* `destination_ip_address` - See Argument Reference above.
* `destination_port` - See Argument Reference above.
* `shared` - See Argument Reference above.
* `enabled` - See Argument Reference above.
* `firewall_policy_id` - The ID of the firewall policy the rule belongs to.
