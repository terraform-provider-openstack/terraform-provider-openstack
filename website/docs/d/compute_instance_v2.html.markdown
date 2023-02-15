---
layout: "openstack"
page_title: "OpenStack: openstack_compute_instance_v2"
sidebar_current: "docs-openstack-datasource-compute-instance-v2"
description: |-
  Get information on an OpenStack Instance
---

# openstack\_compute\_instance\_v2

Use this data source to get the details of a running server

## Example Usage

```hcl
data "openstack_compute_instance_v2" "instance" {
  # Randomly generated UUID, for demonstration purposes
  id = "2ba26dc6-a12d-4889-8f25-794ea5bf4453"
}
```

## Argument Reference

* `id` - (Required) The UUID of the instance


## Attributes Reference

In addition to the above, the following attributes are exported:

* `name` - The name of the server.

* `image_id` - The image ID used to create the server.

* `image_name` - The image name used to create the server.

* `flavor_id` - The flavor ID used to create the server.

* `flavor_name` - The flavor name used to create the server.

* `user_data` - The user data added when the server was created.

* `security_groups` - An array of security group names associated with this server.

* `availability_zone` - The availability zone of this server.

* `network` - An array of maps, detailed below.

* `access_ip_v4` - The first IPv4 address assigned to this server.

* `access_ip_v6` - The first IPv6 address assigned to this server.

* `key_pair` - The name of the key pair assigned to this server.

* `tags` - A set of string tags assigned to this server.

* `metadata` - A set of key/value pairs made available to the server.

* `created` - The creation time of the instance.

* `updated` - The time when the instance was last updated.


The `network` block is defined as:

* `uuid` - The UUID of the network

* `name` - The name of the network

* `fixed_ip_v4` - The IPv4 address assigned to this network port.

* `fixed_ip_v6` - The IPv6 address assigned to this network port.

* `port` - The port UUID for this network

* `mac` - The MAC address assigned to this network interface.

