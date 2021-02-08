---
layout: "openstack"
page_title: "OpenStack: openstack_orchestration_stack_v1"
sidebar_current: "docs-openstack-resource-orchestration-stack-v1"
description: |-
  Manages a V1 stack resource within OpenStack.
---

# openstack\_orchestration\_stack\_v1

Manages a V1 stack resource within OpenStack.

## Example Usage

```hcl
resource "openstack_orchestration_stack_v1" "stack_1" {
  name = "stack_1"
  parameters = {
	length = 4
  }
  template_opts = {
	Bin = "heat_template_version: 2013-05-23\nparameters:\n  length:\n    type: number\nresources:\n  test_res:\n    type: OS::Heat::TestResource\n  random:\n    type: OS::Heat::RandomString\n    properties:\n      length: {get_param: length}\n"
  }
  environment_opts = {
	Bin = "\n"
  }
  disable_rollback = true
  timeout = 30
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the stack. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new stack.

* `name` - (Required) A unique name for the stack. It must start with an
    alphabetic character. Changing this updates the stack's name.

* `template_opts` - (Required) Template key/value pairs to associate with the
    stack which contains either the template file or url.
    Allowed keys: Bin, URL, Files. Changing this updates the existing stack
    Template Opts.

* `environment_opts` - (Optional) Environment key/value pairs to associate with
    the stack which contains details for the environment of the stack.
    Allowed keys: Bin, URL, Files. Changing this updates the existing stack
    Environment Opts.

* `disable_rollback` - (Optional) Enables or disables deletion of all stack
    resources when a stack creation fails. Default is true, meaning all
    resources are not deleted when stack creation fails.

* `parameters` - (Optional) User-defined key/value pairs as parameters to pass
    to the template. Changing this updates the existing stack parameters.

* `timeout` - (Optional) The timeout for stack action in minutes.

* `tags` - (Optional) A list of tags to assosciate with the Stack

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `disable_rollback` - See Argument Reference above.
* `timeout` - See Argument Reference above.
* `parameters` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `capabilities` - List of stack capabilities for stack.
* `description` - The description of the stack resource.
* `notification_topics` - List of notification topics for stack.
* `status` - The status of the stack.
* `status_reason` - The reason for the current status of the stack.
* `template_description` - The description of the stack template.
* `outputs` - A list of stack outputs.
* `creation_time` - The date and time when the resource was created. The date
    and time stamp format is ISO 8601: CCYY-MM-DDThh:mm:ss±hh:mm
    For example, 2015-08-27T09:49:58-05:00. The ±hh:mm value, if included,
    is the time zone as an offset from UTC.
* `updated_time` - The date and time when the resource was updated. The date
    and time stamp format is ISO 8601: CCYY-MM-DDThh:mm:ss±hh:mm
    For example, 2015-08-27T09:49:58-05:00. The ±hh:mm value, if included,
    is the time zone as an offset from UTC.

## Import

stacks can be imported using the `id`, e.g.

```
$ terraform import openstack_orchestration_stack_v1.stack_1 ea257959-eeb1-4c10-8d33-26f0409a755d
```
