---
layout: "openstack"
page_title: "OpenStack: openstack_containerinfra_clustertemplate_v1"
sidebar_current: "docs-openstack-datasource-containerinfra-clustertemplate-v1"
description: |-
  Get information on an OpenStack Magnum cluster template.
---

# openstack\_containerinfra\_clustertemplate\_v1

Use this data source to get the ID of an available OpenStack Magnum cluster
template.

## Example Usage

```hcl
data "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name = "clustertemplate_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 Container Infra
    client.
    If omitted, the `region` argument of the provider is used.

* `name` - (Required) The name of the cluster template.

## Attributes Reference

`id` is set to the ID of the found cluster template. In addition, the following
attributes are exported:

* `region` - See Argument Reference above.

* `name` - See Argument Reference above.

* `project_id` - The project of the cluster template.

* `user_id` - The user of the cluster template.

* `created_at` - The time at which cluster template was created.

* `updated_at` - The time at which cluster template was updated.

* `apiserver_port` - The API server port for the Container Orchestration
    Engine for this cluster template.

* `coe` - The Container Orchestration Engine for this cluster template.

* `cluster_distro` - The distro for the cluster (fedora-atomic, coreos, etc.).

* `dns_nameserver` - Address of the DNS nameserver that is used in nodes of the
    cluster.

* `docker_storage_driver` - Docker storage driver. Changing this updates the
    Docker storage driver of the existing cluster template.

* `docker_volume_size` - The size (in GB) of the Docker volume.

* `external_network_id` - The ID of the external network that will be used for
    the cluster.

* `fixed_network` - The fixed network that will be attached to the cluster.

* `fixed_subnet` - =The fixed subnet that will be attached to the cluster.

* `flavor` - The flavor for the nodes of the cluster.

* `master_flavor` - The flavor for the master nodes.

* `floating_ip_enabled` - Indicates whether created cluster should create IP
    floating IP for every node or not.

* `http_proxy` - The address of a proxy for receiving all HTTP requests and
    relay them.

* `https_proxy` - The address of a proxy for receiving all HTTPS requests and
    relay them.

* `image` - The reference to an image that is used for nodes of the cluster.

* `insecure_registry` - The insecure registry URL for the cluster template.

* `keypair_id` - The name of the Compute service SSH keypair.

* `labels` - The list of key value pairs representing additional properties
    of the cluster template.

* `master_lb_enabled` - Indicates whether created cluster should has a
    loadbalancer for master nodes or not.

* `network_driver` - The name of the driver for the container network.

* `no_proxy` - A comma-separated list of IP addresses that shouldn't be used in
    the cluster.

* `public` - Indicates whether cluster template should be public.

* `registry_enabled` - Indicates whether Docker registry is enabled in the
    cluster.

* `server_type` - The server type for the cluster template.

* `tls_disabled` - Indicates whether the TLS should be disabled in the cluster.

* `volume_driver` - The name of the driver that is used for the volumes of the
    cluster nodes.

* `hidden` - Indicates whether the ClusterTemplate is hidden or not.
