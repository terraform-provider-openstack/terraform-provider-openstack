---
subcategory: "Workflow / Mistral"
layout: "openstack"
page_title: "OpenStack: openstack_workflow_workflow_v2"
sidebar_current: "docs-openstack-datasource-workflow-workflow-v2"
description: |-
  Get information on a workflow.
---

# openstack\_workflow\_workflow\_v2

Use this data source to get the ID of an available Mistral workflow.

## Example Usage

```hcl
data "openstack_workflow_workflow_v2" "hello_workflow" {
  name = "hello_workflow"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Workflow client.

* `name` - (Optional) The name of the workflow.

* `namespace` - (Optional) The namespace of the workflow.

* `project_id` - (Optional) The id of the project to retrieve the workflow.
    Requires admin privileges.

## Attributes Reference

`id` is set to the ID of the found workflow. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `namespace` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `input` - A set of input parameters required for workflow execution.
* `definition` - The workflow definition in Mistral v2 DSL.
* `tags` - A set of string tags for the workflow.
* `scope` - Scope (private or public).
* `created_at` - The date the workflow was created.
