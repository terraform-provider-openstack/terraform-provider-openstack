---
subcategory: "Networking / Neutron"
layout: "openstack"
page_title: "OpenStack: openstack_networking_agent_v2"
sidebar_current: "docs-openstack-datasource-networking-agent-v2"
description: |-
  Get information on an OpenStack Networking Agent.
---

# openstack\_networking\_agent\_v2

Use this data source to get the ID of an available OpenStack Networking agent.

## Example Usage

```hcl
data "openstack_networking_agent_v2" "agent" {
  agent_type = "BGP dynamic routing agent"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve networking agents. If omitted, the
  `region` argument of the provider is used.

* `agent_id` - (Optional) The ID of the agent.

* `agent_type` - (Optional) The type of the agent.

* `alive` - (Optional) Indicates whether this agent is alive and running.

* `availability_zone` - (Optional) The availability zone of the agent.

* `binary` - (Optional) The executable command used to start the agent.

* `description` - (Optional) The description of the agent.

* `host` - (Optional) The hostname of the system the agent is running on.

* `topic` - (Optional) The name of AMQP topic the agent is listening on.

## Attributes Reference

`id` is set to the ID of the found agent. In addition, the following attributes
are exported:

* `agent_type` - See Argument Reference above.
* `alive` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `binary` - See Argument Reference above.
* `description` - See Argument Reference above.
* `host` - See Argument Reference above.
* `topic` - See Argument Reference above.
* `admin_state_up` - The administrative state of the resource.
* `resources_synced` - Indicates the success of the last synchronization attempt to Placement.
* `configurations` - An object containing configuration specific key/value pairs; the semantics of which are determined by the binary name and type.
* `created_at` - Time at which the resource has been created.
* `started_at` - Time at which the agent was started.
* `heartbeat_timestamp` - Time at which the last heartbeat was received.
