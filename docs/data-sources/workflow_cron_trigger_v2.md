---
subcategory: "Workflow / Mistral"
layout: "openstack"
page_title: "OpenStack: openstack_workflow_cron_trigger_v2"
sidebar_current: "docs-openstack-datasource-workflow-cron-trigger-v2"
description: |-
  Retrieve information about an existing Mistral Cron Trigger.
---

# openstack\_workflow\_cron\_trigger\_v2

Use this data source to retrieve information about an existing Mistral Cron Trigger.

## Example Usage

```hcl
data "openstack_workflow_cron_trigger_v2" "cron_trigger" {
  name = "cron_trigger"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Workflow V2 client.
    If omitted, the `region` argument of the provider is used.

* `name` - (Optional) The name of the cron trigger.

* `workflow_id` - (Optional) The ID of the workflow associated with the cron trigger.

* `project_id` - (Optional) The ID of the project from which to retrieve the cron trigger.
    Requires admin privileges.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the cron trigger.
* `region` - The region in which the cron trigger was found.
* `project_id` - The ID of the project owning the cron trigger.
* `name` - The name of the cron trigger.
* `pattern` - The cron-like schedule pattern that defines when the workflow is executed.
* `workflow_id` - The ID of the workflow associated with the cron trigger.
* `workflow_input` - A map of input parameters passed to the workflow when it is executed.
* `workflow_params` - A map of additional workflow execution parameters.
* `created_at` - The time at which the cron trigger was created.
