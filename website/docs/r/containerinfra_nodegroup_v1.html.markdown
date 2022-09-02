---
layout: "openstack"
page_title: "OpenStack: openstack_containerinfra_cluster_v1"
sidebar_current: "docs-openstack-resource-containerinfra-cluster-v1"
description: |-
  Manages a V1 Magnum node group resource within OpenStack.
---

# openstack\_containerinfra\_nodegroup\_v1

Manages a V1 Magnum node group resource within OpenStack.

## Example Usage

### Create a Nodegroup

```hcl
resource "openstack_containerinfra_nodegroup_v1" "nodegroup_1" {
  name                = "nodegroup_1"
  cluster_id          = "b9a45c5c-cd03-4958-82aa-b80bf93cb922"
  node_count          = 5
}
```

## Argument reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 Container Infra
    client. A Container Infra client is needed to create a cluster. If omitted,
    the `region` argument of the provider is used. Changing this creates a new
    node group.

* `name` - (Required) The name of the node group. Changing this creates a new
    node group.

* `project_id` - (Optional) The project of the node group. Required if admin
    wants to create a cluster in another project. Changing this creates a new
    node group.

* `cluster_id` - (Required) The UUID of the V1 Container Infra cluster.
    Changing this creates a new node group.

* `docker_volume_size` - (Optional) The size (in GB) of the Docker volume.
    Changing this creates a new node group.

* `image_id` - (Required) The reference to an image that is used for nodes of the
    node group. Can be set via the `OS_MAGNUM_IMAGE` environment variable.
    Changing this updates the image attribute of the existing node group.

* `flavor_id` - (Optional) The flavor for the nodes of the node group. Can be set
    via the `OS_MAGNUM_FLAVOR` environment variable. Changing this creates a new
    node group.

* `labels` - (Optional) The list of key value pairs representing additional
    properties of the node group. Changing this creates a new node group.

* `merge_labels` - (Optional) Indicates whether the provided labels should be
    merged with cluster labels. Changing this creates a new nodegroup.

* `node_count` - (Optional) The number of nodes for the node group. Changing
    this update the number of nodes of the node group.

* `min_node_count` - (Optional) The minimum number of nodes for the node group.
    Changing this update the minimum number of nodes of the node group.

* `max_node_count` - (Optional) The maximum number of nodes for the node group.
    Changing this update the maximum number of nodes of the node group.

* `role` - (Optional) The role of nodes in the node group. Changing this
    creates a new node group.


## Attributes reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `created_at` - The time at which node group was created.
* `updated_at` - The time at which node group was created.
* `docker_volume_size` - See Argument Reference above.
* `role` - See Argument Reference above.
* `image_id` - See Argument Reference above.
* `flavor_id` - See Argument Reference above.
* `labels` - See Argument Reference above.
* `node_count` - See Argument Reference above.
* `min_node_count` - See Argument Reference above.
* `max_node_count` - See Argument Reference above.
* `role` - See Argument Reference above.

## Import

Node groups can be imported using the `id` (cluster_id/nodegroup_id), e.g.

```
$ terraform import openstack_containerinfra_nodegroup_v1.nodegroup_1 b9a45c5c-cd03-4958-82aa-b80bf93cb922/ce0f9463-dd25-474b-9fe8-94de63e5e42b
```
