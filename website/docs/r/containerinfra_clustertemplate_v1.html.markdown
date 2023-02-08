---
layout: "openstack"
page_title: "OpenStack: openstack_containerinfra_clustertemplate_v1"
sidebar_current: "docs-openstack-resource-containerinfra-clustertemplate-v1"
description: |-
  Manages a V1 Magnum cluster template resource within OpenStack.
---

# openstack\_containerinfra\_clustertemplate\_v1

Manages a V1 Magnum cluster template resource within OpenStack.

## Example Usage

### Create a Cluster template

```hcl
resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name                  = "clustertemplate_1"
  image                 = "Fedora-Atomic-27"
  coe                   = "kubernetes"
  flavor                = "m1.small"
  master_flavor         = "m1.medium"
  dns_nameserver        = "1.1.1.1"
  docker_storage_driver = "devicemapper"
  docker_volume_size    = 10
  volume_driver         = "cinder"
  network_driver        = "flannel"
  server_type           = "vm"
  master_lb_enabled     = true
  floating_ip_enabled   = false

  labels = {
    kube_tag                         = "1.11.1"
    kube_dashboard_enabled           = "true"
    prometheus_monitoring            = "true"
    influx_grafana_dashboard_enabled = "true"
  }
}
```

## Argument reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 Container Infra
    client. A Container Infra client is needed to create a cluster template. If
    omitted,the `region` argument of the provider is used. Changing this
    creates a new cluster template.

* `name` - (Required) The name of the cluster template. Changing this updates
    the name of the existing cluster template.

* `project_id` - (Optional) The project of the cluster template. Required if
    admin wants to create a cluster template in another project. Changing this
    creates a new cluster template.

* `user_id` - (Optional) The user of the cluster template. Required if admin
    wants to create a cluster template for another user. Changing this creates
    a new cluster template.

* `apiserver_port` - (Optional) The API server port for the Container
    Orchestration Engine for this cluster template. Changing this updates the
    API server port of the existing cluster template.

* `coe` - (Required) The Container Orchestration Engine for this cluster
    template. Changing this updates the engine of the existing cluster
    template.

* `cluster_distro` - (Optional) The distro for the cluster (fedora-atomic,
    coreos, etc.). Changing this updates the cluster distro of the existing
    cluster template.

* `dns_nameserver` - (Optional) Address of the DNS nameserver that is used in
    nodes of the cluster. Changing this updates the DNS nameserver of the
    existing cluster template.

* `docker_storage_driver` - (Optional) Docker storage driver. Changing this
    updates the Docker storage driver of the existing cluster template.

* `docker_volume_size` - (Optional) The size (in GB) of the Docker volume.
    Changing this updates the Docker volume size of the existing cluster
    template.

* `external_network_id` - (Optional) The ID of the external network that will
    be used for the cluster. Changing this updates the external network ID of
    the existing cluster template.

* `fixed_network` - (Optional) The fixed network that will be attached to the
    cluster. Changing this updates the fixed network of the existing cluster
    template.

* `fixed_subnet` - (Optional) The fixed subnet that will be attached to the
    cluster. Changing this updates the fixed subnet of the existing cluster
    template.

* `flavor` - (Optional) The flavor for the nodes of the cluster. Can be set via
    the `OS_MAGNUM_FLAVOR` environment variable. Changing this updates the
    flavor of the existing cluster template.

* `master_flavor` - (Optional) The flavor for the master nodes. Can be set via
    the `OS_MAGNUM_MASTER_FLAVOR` environment variable. Changing this updates
    the master flavor of the existing cluster template.

* `floating_ip_enabled` - (Optional) Indicates whether created cluster should
    create floating IP for every node or not. Changing this updates the
    floating IP enabled attribute of the existing cluster template.

* `http_proxy` - (Optional) The address of a proxy for receiving all HTTP
    requests and relay them. Changing this updates the HTTP proxy address of
    the existing cluster template.

* `https_proxy` - (Optional) The address of a proxy for receiving all HTTPS
    requests and relay them. Changing this updates the HTTPS proxy address of
    the existing cluster template.

* `image` - (Required) The reference to an image that is used for nodes of the
    cluster. Can be set via the `OS_MAGNUM_IMAGE` environment variable.
    Changing this updates the image attribute of the existing cluster template.

* `insecure_registry` - (Optional) The insecure registry URL for the cluster
    template. Changing this updates the insecure registry attribute of the
    existing cluster template.

* `keypair_id` - (Optional) The name of the Compute service SSH keypair.
    Changing this updates the keypair of the existing cluster template.

* `labels` - (Optional) The list of key value pairs representing additional
    properties of the cluster template. Changing this updates the labels of the
    existing cluster template.

* `master_lb_enabled` - (Optional) Indicates whether created cluster should
    has a loadbalancer for master nodes or not. Changing this updates the
    attribute of the existing cluster template.

* `network_driver` - (Optional) The name of the driver for the container
    network. Changing this updates the network driver of the existing cluster
    template.

* `no_proxy` - (Optional) A comma-separated list of IP addresses that shouldn't
    be used in the cluster. Changing this updates the no proxy list of the
    existing cluster template.

* `public` - (Optional) Indicates whether cluster template should be public.
    Changing this updates the public attribute of the existing cluster
    template.

* `registry_enabled` - (Optional) Indicates whether Docker registry is enabled
    in the cluster. Changing this updates the registry enabled attribute of the
    existing cluster template.

* `server_type` - (Optional) The server type for the cluster template. Changing
    this updates the server type of the existing cluster template.

* `tls_disabled` - (Optional) Indicates whether the TLS should be disabled in
    the cluster. Changing this updates the attribute of the existing cluster.

* `volume_driver` - (Optional) The name of the driver that is used for the
    volumes of the cluster nodes. Changing this updates the volume driver of
    the existing cluster template.

* `hidden` - (Optional) Indicates whether the ClusterTemplate is hidden or not.
    Changing this updates the hidden attribute of the existing cluster
    template.

## Attributes reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `created_at` - The time at which cluster template was created.
* `updated_at` - The time at which cluster template was created.
* `apiserver_port` - See Argument Reference above.
* `coe` - See Argument Reference above.
* `cluster_distro` - See Argument Reference above.
* `dns_nameserver` - See Argument Reference above.
* `docker_storage_driver` - See Argument Reference above.
* `docker_volume_size` - See Argument Reference above.
* `external_network_id` - See Argument Reference above.
* `fixed_network` - See Argument Reference above.
* `fixed_subnet` - See Argument Reference above.
* `flavor` - See Argument Reference above.
* `master_flavor` - See Argument Reference above.
* `floating_ip_enabled` - See Argument Reference above.
* `http_proxy` - See Argument Reference above.
* `https_proxy` - See Argument Reference above.
* `image` - See Argument Reference above.
* `insecure_registry` - See Argument Reference above.
* `keypair_id` - See Argument Reference above.
* `labels` - See Argument Reference above.
* `links` - A list containing associated cluster template links.
* `master_lb_enabled` - See Argument Reference above.
* `network_driver` - See Argument Reference above.
* `no_proxy` - See Argument Reference above.
* `public` - See Argument Reference above.
* `registry_enabled` - See Argument Reference above.
* `server_type` - See Argument Reference above.
* `tls_disabled` - See Argument Reference above.
* `volume_driver` - See Argument Reference above.
* `hidden` - See Argument Reference above.

## Import

Cluster templates can be imported using the `id`, e.g.

```
$ terraform import openstack_containerinfra_clustertemplate_v1.clustertemplate_1 b9a45c5c-cd03-4958-82aa-b80bf93cb922
```
