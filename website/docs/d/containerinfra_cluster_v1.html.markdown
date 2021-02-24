---
layout: "openstack"
page_title: "OpenStack: openstack_containerinfra_cluster_v1"
sidebar_current: "docs-openstack-datasource-containerinfra-cluster-v1"
description: |-
  Get information on an OpenStack Magnum cluster.
---

# openstack\_containerinfra\_cluster\_v1

Use this data source to get the ID of an available OpenStack Magnum cluster.

## Example Usage

```hcl
data "openstack_containerinfra_cluster_v1" "cluster_1" {
  name = "cluster_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 Container Infra
    client.
    If omitted, the `region` argument of the provider is used.

* `name` - (Required) The name of the cluster.

## Attributes Reference

`id` is set to the ID of the found cluster. In addition, the following
attributes are exported:

* `region` - See Argument Reference above.

* `name` - See Argument Reference above.

* `project_id` - The project of the cluster.

* `user_id` - The user of the cluster.

* `created_at` - The time at which cluster was created.

* `updated_at` - The time at which cluster was updated.

* `api_address` - COE API address.

* `coe_version` - COE software version.

* `cluster_template_id` - The UUID of the V1 Container Infra cluster template.

* `create_timeout` - The timeout (in minutes) for creating the cluster.

* `discovery_url` - The URL used for cluster node discovery.

* `docker_volume_size` - The size (in GB) of the Docker volume.

* `flavor` - The flavor for the nodes of the cluster.

* `master_flavor` - The flavor for the master nodes.

* `keypair` - The name of the Compute service SSH keypair.

* `labels` - The list of key value pairs representing additional properties of
    the cluster.

* `master_count` - The number of master nodes for the cluster.

* `node_count` - The number of nodes for the cluster.

* `fixed_network` - The fixed network that is attached to the cluster.

* `fixed_subnet` - The fixed subnet that is attached to the cluster.

* `master_addresses` - IP addresses of the master node of the cluster.

* `node_addresses` - IP addresses of the node of the cluster.

* `stack_id` - UUID of the Orchestration service stack.
