---
subcategory: "Load Balancing as a Service / Octavia"
layout: "openstack"
page_title: "OpenStack: openstack_lb_monitor_v2"
sidebar_current: "docs-openstack-datasource-lb-monitor-v2"
description: |-
  Get information on an OpenStack Load Balancer Monitor.
---

# openstack\_lb\_monitor\_v2

Use this data source to get the ID of an OpenStack Load Balancer monitor.

## Example Usage

```hcl
data "openstack_lb_monitor_v2" "monitor_1" {
  name = "monitor_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Load Balancer
  client. If omitted, the `region` argument of the provider is used.

* `monitor_id` - (Optional) The ID of the monitor. Exactly one of `name`,
  `monitor_id` is required to be set.

* `name` - (Optional) The name of the monitor. Exactly one of `name`,
  `monitor_id` is required to be set.

* `pool_id` - (Optional) The Pool to Monitor.

* `type` - (Optional) The type of probe, which is PING, TCP, HTTP, or HTTPS,
  that is sent by the load balancer to verify the member state.

* `http_method` - (Optional) The HTTP method used for requests by the Monitor.

* `url_path` - (Optional) URI path that will be accessed if Monitor type
  is HTTP or HTTPS.

* `status` - (Optional) The status of the health monitor. Indicates whether
  the health monitor is operational.

* `expected_codes` - (Optional) Expected HTTP codes for a passing HTTP(S)
  monitor.

* `tags` - (Optional) Tags is a list of resource tags. Tags are arbitrarily
  defined strings attached to the resource.

## Attributes Reference

`id` is set to the ID of the found monitor. In addition, the following
attributes are exported:

* `project_id` - The owner (project/tenant) ID of the monitor.

* `name` - See Argument Reference above.

* `type` - See Argument Reference above.

* `delay` - The time, in seconds, between sending probes to members.

* `timeout` - The maximum number of seconds for a monitor to wait for a
  connection to be established before it times out.

* `max_retries` - Number of allowed connection failures before changing the
  status of the member to INACTIVE.

* `max_retries_down` - Number of allowed connection failures before changing
  the status of the member to Error.

* `http_method` - See Argument Reference above.

* `http_version` - The HTTP version that the monitor uses for requests.

* `url_path` - See Argument Reference above.

* `expected_codes` - See Argument Reference above.

* `domain_name` - The HTTP host header that the monitor uses for requests.

* `admin_state_up` - The administrative state of the health monitor, which is
  up (true) or down (false).

* `status` - See Argument Reference above.

* `pools` - List of pools that are associated with the health monitor.

* `provisioning_status` - The provisioning status of the Monitor.

* `operating_status` - The operating status of the monitor.

* `tags` - See Argument Reference above.
