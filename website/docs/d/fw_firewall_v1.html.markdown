---
layout: "openstack"
page_title: "OpenStack: openstack_fw_firewall_v1"
sidebar_current: "docs-openstack-datasource-fw-firewall-v1"
description: |-
  Get information on an OpenStack Firewall.
---

# openstack\_fw\_firewall_v1

Use this data source to get information of an available OpenStack firewall.

## Example Usage

```hcl
data "openstack_fw_firewall_v1" "firewall" {
  name = "tf_test_firewall"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve firewall ids. If omitted, the
  `region` argument of the provider is used.

* `id` - (Optional) The ID of the firewall firewall.
  Required if `name` is not used.

* `name` - (Optional) The name of the firewall firewall.
  Required if `id` is not used.

* `tenant_id` - (Optional) The owner of the firewall firewall.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `policy_id` - The active policy ID associated with the firewall.
* `description` - The description of the firewall.
* `admin_state_up` - The administrative state of the firewall.
* `status` - The status of the firewall service.
  Values are ACTIVE, INACTIVE, ERROR, DOWN, PENDING_CREATE,
  PENDING_UPDATE, or PENDING_DELETE.
