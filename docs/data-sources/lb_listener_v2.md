---
subcategory: "Load Balancing as a Service / Octavia"
layout: "openstack"
page_title: "OpenStack: openstack_lb_listener_v2"
sidebar_current: "docs-openstack-datasource-lb-listener-v2"
description: |-
  Get information on an OpenStack Load Balancer Listener.
---

# openstack\_lb\_listener\_v2

Use this data source to get the ID of an OpenStack Load Balancer listener.

## Example Usage

```hcl
data "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Load Balancer client.
    If omitted, the `region` argument of the provider is used.

* `listener_id` - (Optional) The ID of the listener. Exactly one of `name`,
  `listener_id` is required to be set.

* `name` - (Optional) The name of the listener. Exactly one of `name`,
  `listener_id` is required to be set.

* `tags` - (Optional) A set of tags applied to the loadbalancer's listener.
  The loadbalancer' listener will be returned if it has all of the specified tags.

* `loadbalancer_id` - (Optional) The ID of the load balancer associated with
  the requested listener.

* `protocol` - The protocol of the requested listener.

* `protocol_port` - The port on which the requested listener accepts client traffic.

## Attributes Reference

`id` is set to the ID of the found listener. In addition, the following attributes
are exported:

* `project_id` - The owner (project/tenant) ID of the listener.

* `name` - The name of the listener.

* `description` - The description of the listener.

* `protocol` - The protocol to loadbalance.

* `protocol_port` - The port on which to listen to client traffic that is
  associated with the Loadbalancer.

* `default_pool_id` - The UUID of default pool.

* `default_pool` - The default pool with which the Listener is associated.

* `loadbalancers` - A list of load balancer IDs.

* `connection_limit` - The maximum number of connections allowed for the Loadbalancer.

* `sni_container_refs` - The list of references to TLS secrets.

* `default_tls_container_ref` - A reference to a Barbican container of TLS secrets.

* `admin_state_up` - The administrative state of the Listener.

* `pools` - Pools are the pools which are part of this listener.

* `l7policies` - L7policies are the L7 policies which are part of this listener.

* `provisioning_status` - The provisioning status of the Listener.

* `timeout_client_data` - Frontend client inactivity timeout in milliseconds.

* `timeout_member_data` - Backend member inactivity timeout in milliseconds.

* `timeout_member_connect` - Backend member connection timeout in milliseconds.

* `timeout_tcp_inspect` - Time, in milliseconds, to wait for additional TCP
  packets for content inspection.

* `insert_headers` - A dictionary of optional headers to insert into the request
  before it is sent to the backend member.

* `allowed_cidrs` - A list of IPv4, IPv6 or mix of both CIDRs.

* `tls_ciphers` - List of ciphers in OpenSSL format (colon-separated).

* `tls_versions` - A list of TLS protocol versions.

* `tags` - Tags is a list of resource tags.

* `alpn_protocols` - A list of ALPN protocols.

* `client_authentication` - The TLS client authentication mode.

* `client_ca_tls_container_ref` - The ref of the key manager service secret
  containing a PEM format client CA certificate bundle for TERMINATED_HTTPS listeners.

* `client_crl_container_ref` - The URI of the key manager service secret
  containing a PEM format CA revocation list file for TERMINATED_HTTPS listeners.

* `hsts_include_subdomains` - Defines whether the includeSubDomains directive
  should be added to the Strict-Transport-Security HTTP response header.

* `hsts_max_age` - The value of the max_age directive for the
  Strict-Transport-Security HTTP response header.

* `hsts_preload` - Defines whether the preload directive should be added to the
  Strict-Transport-Security HTTP response header.

* `operating_status` - The operating status of the resource.
