---
layout: "openstack"
page_title: "OpenStack: openstack_containerinfra_cluster_v1"
sidebar_current: "docs-openstack-resource-containerinfra-cluster-v1"
description: |-
  Manages a V1 Magnum cluster resource within OpenStack.
---

# openstack\_containerinfra\_cluster_v1

Manages a V1 Magnum cluster resource within OpenStack.

## Example Usage

### Create a Cluster

```hcl
resource "openstack_containerinfra_cluster_v1" "cluster_1" {
  name                 = "cluster_1"
  cluster_template_id  = "b9a45c5c-cd03-4958-82aa-b80bf93cb922"
  master_count         = 3
  node_count           = 5
  keypair              = "ssh_keypair"
}
```

## Argument reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 Container Infra
    client. A Container Infra client is needed to create a cluster. If omitted,
    the `region` argument of the provider is used. Changing this creates a new
    cluster.

* `name` - (Required) The name of the cluster. Changing this updates the name
    of the existing cluster template.

* `project_id` - (Optional) The project of the cluster. Required if admin wants
    to create a cluster in another project. Changing this creates a new
    cluster.

* `user_id` - (Optional) The user of the cluster. Required if admin wants to
    create a cluster template for another user. Changing this creates a new
    cluster.

* `cluster_template_id` - (Required) The UUID of the V1 Container Infra cluster
    template. Changing this creates a new cluster.

* `create_timeout` - (Optional) The timeout (in minutes) for creating the
    cluster. Changing this creates a new cluster.

* `discovery_url` - (Optional) The URL used for cluster node discovery.
    Changing this creates a new cluster.

* `docker_volume_size` - (Optional) The size (in GB) of the Docker volume.
    Changing this creates a new cluster.

* `flavor` - (Optional) The flavor for the nodes of the cluster. Can be set via
    the `OS_MAGNUM_FLAVOR` environment variable. Changing this creates a new
    cluster.

* `master_flavor` - (Optional) The flavor for the master nodes. Can be set via
    the `OS_MAGNUM_MASTER_FLAVOR` environment variable. Changing this creates a
    new cluster.

* `keypair` - (Optional) The name of the Compute service SSH keypair. Changing
    this creates a new cluster.

* `labels` - (Optional) The list of key value pairs representing additional
    properties of the cluster. Changing this creates a new cluster.

* `master_count` - (Optional) The number of master nodes for the cluster.
    Changing this creates a new cluster.

* `node_count` - (Optional) The number of nodes for the cluster. Changing this
    creates a new cluster.

## Attributes reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `created_at` - The time at which cluster was created.
* `updated_at` - The time at which cluster was created.
* `api_address` - COE API address.
* `coe_version` - COE software version.
* `cluster_template_id` - See Argument Reference above.
* `container_version` - Container software version.
* `create_timeout` - See Argument Reference above.
* `discovery_url` - See Argument Reference above.
* `docker_volume_size` - See Argument Reference above.
* `flavor` - See Argument Reference above.
* `master_flavor` - See Argument Reference above.
* `keypair` - See Argument Reference above.
* `labels` - See Argument Reference above.
* `master_count` - See Argument Reference above.
* `node_count` - See Argument Reference above.
* `master_addresses` - IP addresses of the master node of the cluster.
* `node_addresses` - IP addresses of the node of the cluster.
* `stack_id` - UUID of the Orchestration service stack.

## Import

Clusters can be imported using the `id`, e.g.

```
$ terraform import openstack_containerinfra_cluster_v1.cluster_1 ce0f9463-dd25-474b-9fe8-94de63e5e42b
```
