---
layout: "openstack"
page_title: "OpenStack: openstack_containerinfra_nodegroup_v1"
sidebar_current: "docs-openstack-datasource-containerinfra-nodegroup-v1"
description: |-
  Get information on an OpenStack Magnum node group.
---

# openstack\_containerinfra\_nodegroup\_v1

Use this data source to get information of an available OpenStack Magnum node group.

## Example Usage

```hcl
data "openstack_containerinfra_nodegroup_v1" "nodegroup_1" {
  cluster_id = "cluster_1"
  name       = "nodegroup_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 Container Infra
    client.
    If omitted, the `region` argument of the provider is used.

* `cluster_id` - (Required) The name of the OpenStack Magnum cluster.

* `name` - (Required) The name of the node group.

## Attributes Reference

`id` is set to the ID of the found node group. In addition, the following
attributes are exported:

* `name` - See Argument Reference above.

* `region` - See Argument Reference above.

* `project_id` - The project of the node group.

* `created_at` - The time at which the node group was created.

* `updated_at` - The time at which the node group was updated.

* `docker_volume_size` - The size (in GB) of the Docker volume.

* `labels` - The list of key value pairs representing additional properties of
    the node group.

* `role` - The role of the node group.

* `node_count` - The number of nodes for the node group.

* `min_node_count` - The minimum number of nodes for the node group.

* `max_node_count` - The maximum number of nodes for the node group.

* `image` - The reference to an image that is used for nodes of the node group.

* `flavor` - The flavor for the nodes of the node group.
