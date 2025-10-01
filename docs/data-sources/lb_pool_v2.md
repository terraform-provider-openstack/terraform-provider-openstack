---
subcategory: "Load Balancing as a Service / Octavia"
layout: "openstack"
page_title: "OpenStack: openstack_lb_pool_v2"
sidebar_current: "docs-openstack-datasource-lb-pool-v2"
description: |-
  Get information on an OpenStack Load Balancer Pool.
---

# openstack\_lb\_pool\_v2

Use this data source to get the ID of an OpenStack Load Balancer pool.

## Example Usage

```hcl
data "openstack_lb_pool_v2" "pool_1" {
  name = "pool_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Load Balancer
  client. If omitted, the `region` argument of the provider is used.

* `pool_id` - (Optional) The ID of the pool. Exactly one of `name`, `pool_id`
  is required to be set.

* `name` - (Optional) The name of the pool. Exactly one of `name`, `pool_id`
  is required to be set.

* `tags` - (Optional) A set of tags applied to the loadbalancer's pool. The
  loadbalancer' pool will be returned if it has all of the specified tags.

* `loadbalancer_id` - (Optional) The ID of the load balancer associated with
  the requested pool.

* `protocol` - The protocol of the requested pool.

* `lb_method` - The load balancing algorithm to distribute traffic to the
  pool's members.

## Attributes Reference

`id` is set to the ID of the found pool. In addition, the following attributes
are exported:

* `project_id` - The owner (project/tenant) ID of the pool.

* `name` - The name of the pool.

* `description` - The description of the pool.

* `protocol` - The protocol to loadbalance.

* `lb_method` - The load-balancer algorithm, which is round-robin,
  least-connections, and so on.

* `listeners` - A list of listeners objects IDs.

* `members` - A list of member objects IDs.

* `healthmonitor_id` - The ID of associated health monitor.

* `admin_state_up` - The administrative state of the Pool, which is up (true)
  or down (false).

* `loadbalancers` - A list of load balancer objects IDs.

* `session_persistence` - Indicates whether connections in the same session
  will be processed by the same Pool member or not.

* `alpn_protocols` - A list of ALPN protocols.

* `ca_tls_container_ref` - The reference of the key manager service secret
  containing a PEM format CA certificate bundle for tls_enabled pools.

* `crl_container_ref` - The reference of the key manager service secret
  containing a PEM format CA revocation list file for tls_enabled pools.

* `tls_enabled` - When true connections to backend member servers will use
  TLS encryption.

* `tls_ciphers` - List of ciphers in OpenSSL format (colon-separated).

* `tls_container_ref` - The reference to the key manager service secret
  containing a PKCS12 format certificate/key bundle for tls_enabled pools for
  TLS client authentication to the member servers.

* `tls_versions` - A list of TLS protocol versions.

* `provisioning_status` - The provisioning status of the pool.

* `operating_status` - The operating status of the pool.

* `tags` - Tags is a list of resource tags.
