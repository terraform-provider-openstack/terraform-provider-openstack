---
subcategory: "Load Balancing as a Service / Octavia"
layout: "openstack"
page_title: "OpenStack: openstack_lb_listener_v2"
sidebar_current: "docs-openstack-resource-lb-listener-v2"
description: |-
  Manages a V2 listener resource within OpenStack.
---

# openstack\_lb\_listener\_v2

Manages a V2 listener resource within OpenStack.

~> **Note:** This resource has attributes that depend on octavia minor versions.
Please ensure your Openstack cloud supports the required [minor version](../#octavia-api-versioning).

## Example Usage

### Simple listener

```hcl
resource "openstack_lb_listener_v2" "listener_1" {
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

  insert_headers = {
    X-Forwarded-For = "true"
  }
}
```

### Listener with TLS and client certificate authentication

```hcl
resource "openstack_keymanager_secret_v1" "certificate_1" {
  name                     = "certificate"
  payload                  = filebase64("snakeoil.p12")
  payload_content_encoding = "base64"
  payload_content_type     = "application/octet-stream"
}

resource "openstack_keymanager_secret_v1" "ca_certificate_1" {
  name                 = "certificate"
  payload              = file("CA.pem")
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

data "openstack_networking_subnet_v2" "subnet_1" {
  name = "my-subnet"
}

resource "openstack_lb_loadbalancer_v2" "lb_1" {
  name          = "loadbalancer"
  vip_subnet_id = data.openstack_networking_subnet_v2.subnet_1.id
}

resource "openstack_lb_listener_v2" "listener_1" {
  name                        = "https"
  protocol                    = "TERMINATED_HTTPS"
  protocol_port               = 443
  loadbalancer_id             = openstack_lb_loadbalancer_v2.lb_1.id
  default_tls_container_ref   = openstack_keymanager_secret_v1.certificate_1
  client_authentication       = "OPTIONAL"
  client_ca_tls_container_ref = openstack_keymanager_secret_v1.ca_certificate_2.secret_ref
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a listener. If omitted, the `region`
argument of the provider is used. Changing this creates a new Listener.

* `protocol` - (Required) The protocol can be either `TCP`, `HTTP`, `HTTPS`,
  `TERMINATED_HTTPS`, `UDP`, `SCTP` (supported only in **Octavia minor version
  \>= 2.23**), or `PROMETHEUS` (supported only in **Octavia minor version >=
  2.25**). Changing this creates a new Listener.

* `protocol_port` - (Required) The port on which to listen for client traffic.
* Changing this creates a new Listener.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
  the Listener.  Only administrative users can specify a tenant UUID other than
  their own. Changing this creates a new Listener.

* `loadbalancer_id` - (Required) The load balancer on which to provision this
  Listener. Changing this creates a new Listener.

* `name` - (Optional) Human-readable name for the Listener. Does not have to be
  unique.

* `default_pool_id` - (Optional) The ID of the default pool with which the
  Listener is associated.

* `description` - (Optional) Human-readable description for the Listener.

* `connection_limit` - (Optional) The maximum number of connections allowed for
  the Listener.

* `timeout_client_data` - (Optional) The client inactivity timeout in
  milliseconds.

* `timeout_member_connect` - (Optional) The member connection timeout in
  milliseconds.

* `timeout_member_data` - (Optional) The member inactivity timeout in
  milliseconds.

* `timeout_tcp_inspect` - (Optional) The time in milliseconds, to wait for
  additional TCP packets for content inspection.

* `default_tls_container_ref` – (Optional) A reference to a Barbican Secrets
  container that stores TLS information. This is required when the protocol is
  `TERMINATED_HTTPS`. For more information, see the
  [Octavia TLS-terminated HTTPS load balancer guide][octavia-tls-guide].

* `sni_container_refs` – (Optional) A list of references to Barbican Secrets
  containers that store SNI information. For more information, see the
  [Octavia TLS-terminated HTTPS load balancer guide][octavia-tls-guide].

* `admin_state_up` - (Optional) The administrative state of the Listener. A
  valid value is true (UP) or false (DOWN).

* `insert_headers` - (Optional) The list of key value pairs representing
  headers to insert into the request before it is sent to the backend members.
  Changing this updates the headers of the existing listener.

* `allowed_cidrs` - (Optional) A list of CIDR blocks that are permitted to
  connect to this listener, denying all other source addresses. If not present,
  defaults to allow all.

* `alpn_protocols` - (Optional) A list of ALPN protocols. Available protocols:
  `http/1.0`, `http/1.1`, `h2`. Supported only in **Octavia minor version >=
  2.20**.

* `client_authentication` - (Optional) The TLS client authentication mode.
  Available options: `NONE`, `OPTIONAL` or `MANDATORY`. Requires
  `TERMINATED_HTTPS` listener protocol and the `client_ca_tls_container_ref`.
  Supported only in **Octavia minor version >= 2.8**.

* `client_ca_tls_container_ref` - (Optional) The ref of the key manager service
  secret containing a PEM format client CA certificate bundle for
  `TERMINATED_HTTPS` listeners. Required if `client_authentication` is
  `OPTIONAL` or `MANDATORY`. Supported only in **Octavia minor version >=
  2.8**.

* `client_crl_container_ref` - (Optional) The URI of the key manager service
  secret containing a PEM format CA revocation list file for `TERMINATED_HTTPS`
  listeners. Supported only in **Octavia minor version >= 2.8**.

* `hsts_include_subdomains` - (Optional) Defines whether the
  **includeSubDomains** directive should be added to the
  Strict-Transport-Security HTTP response header. This requires setting the
  `hsts_max_age` option as well in order to become effective. Requires
  `TERMINATED_HTTPS` listener protocol. Supported only in **Octavia minor
  version >= 2.27**.

* `hsts_max_age` - (Optional) The value of the **max_age** directive for the
  Strict-Transport-Security HTTP response header. Setting this enables HTTP
  Strict Transport Security (HSTS) for the TLS-terminated listener. Requires
  `TERMINATED_HTTPS` listener protocol. Supported only in **Octavia minor
  version >= 2.27**.

* `hsts_preload` - (Optional) Defines whether the **preload** directive should
  be added to the Strict-Transport-Security HTTP response header. This requires
  setting the `hsts_max_age` option as well in order to become effective.
  Requires `TERMINATED_HTTPS` listener protocol. Supported only in **Octavia
  minor version >= 2.27**.

* `tls_ciphers` - (Optional) List of ciphers in OpenSSL format
  (colon-separated). See
  <https://docs.openssl.org/1.1.1/man1/ciphers/> for more information.
  Supported only in **Octavia minor version >= 2.15**.

* `tls_versions` - (Optional) A list of TLS protocol versions. Available
  versions: `TLSv1`, `TLSv1.1`, `TLSv1.2`, `TLSv1.3`. Supported only in
  **Octavia minor version >= 2.17**.

* `tags` - (Optional) A list of simple strings assigned to the pool. Available
    for Octavia **minor version 2.5 or later**.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the Listener.
* `protocol` - See Argument Reference above.
* `protocol_port` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `loadbalancer_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `default_pool_id` - See Argument Reference above.
* `description` - See Argument Reference above.
* `connection_limit` - See Argument Reference above.
* `timeout_client_data` - See Argument Reference above.
* `timeout_member_connect` - See Argument Reference above.
* `timeout_member_data` - See Argument Reference above.
* `timeout_tcp_inspect` - See Argument Reference above.
* `default_tls_container_ref` - See Argument Reference above.
* `sni_container_refs` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `insert_headers` - See Argument Reference above.
* `allowed_cidrs` - See Argument Reference above.
* `alpn_protocols` - See Argument Reference above.
* `client_authentication` - See Argument Reference above.
* `client_ca_tls_container_ref` - See Argument Reference above.
* `client_crl_container_ref` - See Argument Reference above.
* `hsts_include_subdomains` - See Argument Reference above.
* `hsts_max_age` - See Argument Reference above.
* `hsts_preload` - See Argument Reference above.
* `tls_ciphers` - See Argument Reference above.
* `tls_versions` - See Argument Reference above
* `tags` - See Argument Reference above.

## Import

Load Balancer Listener can be imported using the Listener ID, e.g.:

```shell
terraform import openstack_lb_listener_v2.listener_1 b67ce64e-8b26-405d-afeb-4a078901f15a
```

[octavia-tls-guide]: https://docs.openstack.org/octavia/latest/user/guides/basic-cookbook.html#deploy-a-tls-terminated-https-load-balancer
