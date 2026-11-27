---
subcategory: "Load Balancing as a Service / Octavia"
layout: "openstack"
page_title: "OpenStack: openstack_lb_pool_v2"
sidebar_current: "docs-openstack-resource-lb-pool-v2"
description: |-
  Manages a V2 pool resource within OpenStack.
---

# openstack\_lb\_pool\_v2

Manages a V2 pool resource within OpenStack.

~> **Note:** This resource has attributes that depend on octavia minor versions.
Please ensure your Openstack cloud supports the required [minor version](../#octavia-api-versioning).

## Example Usage

```hcl
resource "openstack_lb_pool_v2" "pool_1" {
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a pool. If omitted, the `region`
  argument of the provider is used. Changing this creates a new pool.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
  the pool.  Only administrative users can specify a tenant UUID other than
  their own. Changing this creates a new pool.

* `name` - (Optional) Human-readable name for the pool.

* `description` - (Optional) Human-readable description for the pool.

* `protocol` - (Required) The protocol - can either be TCP, HTTP, HTTPS, PROXY,
  UDP, PROXYV2 (**Octavia minor version >= 2.22**) or SCTP (**Octavia minor
  version >= 2.23**). Changing this creates a new pool.

* `loadbalancer_id` - (Optional) The load balancer on which to provision this
  pool. Changing this creates a new pool. Note: One of LoadbalancerID or
  ListenerID must be provided.

* `listener_id` - (Optional) The Listener on which the members of the pool will
  be associated with. Changing this creates a new pool. Note: One of
  LoadbalancerID or ListenerID must be provided.

* `lb_method` - (Required) The load balancing algorithm to distribute traffic
  to the pool's members. Must be one of ROUND_ROBIN, LEAST_CONNECTIONS,
  SOURCE_IP, or SOURCE_IP_PORT.

* `persistence` - (Optional) Omit this field to prevent session persistence.
  Indicates whether connections in the same session will be processed by the
  same Pool member or not. Changing this creates a new pool.

* `alpn_protocols` - (Optional) A list of ALPN protocols. Available protocols:
  `http/1.0`, `http/1.1`, `h2`. Supported only in **Octavia minor version >=
  2.24**.

* `ca_tls_container_ref` - (Optional) The reference of the key manager service
  secret containing a PEM format CA certificate bundle for `tls_enabled` pools.
  Supported only in **Octavia minor version >= 2.8**.

* `crl_container_ref` - (Optional) The reference of the key manager service
  secret containing a PEM format CA revocation list file for `tls_enabled`
  pools. Supported only in **Octavia minor version >= 2.8**.

* `tls_enabled` - (Optional) When true connections to backend member servers
  will use TLS encryption. Default is false. Supported only in **Octavia minor
  version >= 2.8**.

* `tls_container_ref` - (Optional) The reference to the key manager service
  secret containing a PKCS12 format certificate/key bundle for `tls_enabled`
  pools for TLS client authentication to the member servers. Supported only in
  **Octavia minor version >= 2.8**.

* `tls_ciphers` - (Optional) List of ciphers in OpenSSL format
  (colon-separated). See
  <https://docs.openssl.org/1.1.1/man1/ciphers/> for more information.
  Supported only in **Octavia minor version >= 2.15**.

* `tls_versions` - (Optional) A list of TLS protocol versions. Available
  versions: `TLSv1`, `TLSv1.1`, `TLSv1.2`, `TLSv1.3`. Supported only in
  **Octavia minor version >= 2.17**.

* `admin_state_up` - (Optional) The administrative state of the pool. A valid
  value is true (UP) or false (DOWN).

The `persistence` argument supports:

* `type` - (Required) The type of persistence mode. The current specification
  supports SOURCE_IP, HTTP_COOKIE, and APP_COOKIE.

* `cookie_name` - (Optional) The name of the cookie if persistence mode is set
  appropriately. Required if `type = APP_COOKIE`.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the pool.
* `tenant_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `loadbalancer_id` - See Argument Reference above.
* `listener_id` - See Argument Reference above.
* `lb_method` - See Argument Reference above.
* `persistence` - See Argument Reference above.
* `alpn_protocols` - See Argument Reference above.
* `ca_tls_container_ref` - See Argument Reference above.
* `crl_container_ref` - See Argument Reference above.
* `tls_enabled` - See Argument Reference above.
* `tls_container_ref` - See Argument Reference above.
* `tls_ciphers` - See Argument Reference above.
* `tls_versions` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.

## Import

Load Balancer Pool can be imported using the Pool ID, e.g.:

```shell
terraform import openstack_lb_pool_v2.pool_1 60ad9ee4-249a-4d60-a45b-aa60e046c513
```
