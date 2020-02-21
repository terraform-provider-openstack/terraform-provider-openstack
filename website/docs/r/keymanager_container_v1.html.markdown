---
layout: "openstack"
page_title: "OpenStack: openstack_keymanager_container_v1"
sidebar_current: "docs-openstack-resource-keymanager-container-v1"
description: |-
  Manages a V1 Barbican container resource within OpenStack.
---

# openstack\_keymanager\_container\_v1

Manages a V1 Barbican container resource within OpenStack.

## Example Usage

### Simple secret

The container with the TLS certificates, which can be used by the loadbalancer HTTPS listener.

```hcl
resource "openstack_keymanager_secret_v1" "certificate_1" {
  name                 = "certificate"
  payload              = "${file("cert.pem")}"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "openstack_keymanager_secret_v1" "private_key_1" {
  name                 = "private_key"
  payload              = "${file("cert-key.pem")}"
  secret_type          = "private"
  payload_content_type = "text/plain"
}

resource "openstack_keymanager_secret_v1" "intermediate_1" {
  name                 = "intermediate"
  payload              = "${file("intermediate-ca.pem")}"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "openstack_keymanager_container_v1" "tls_1" {
  name = "tls"
  type = "certificate"

  secret_refs {
    name       = "certificate"
    secret_ref = "${openstack_keymanager_secret_v1.certificate_1.secret_ref}"
  }

  secret_refs {
    name       = "private_key"
    secret_ref = "${openstack_keymanager_secret_v1.private_key_1.secret_ref}"
  }

  secret_refs {
    name       = "intermediates"
    secret_ref = "${openstack_keymanager_secret_v1.intermediate_1.secret_ref}"
  }
}

data "openstack_networking_subnet_v2" "subnet_1" {
  name = "my-subnet"
}

resource "openstack_lb_loadbalancer_v2" "lb_1" {
  name          = "loadbalancer"
  vip_subnet_id = "${data.openstack_networking_subnet_v2.subnet_1.id}"
}

resource "openstack_lb_listener_v2" "listener_1" {
  name                      = "https"
  protocol                  = "TERMINATED_HTTPS"
  protocol_port             = 443
  loadbalancer_id           = "${openstack_lb_loadbalancer_v2.lb_1.id}"
  default_tls_container_ref = "${openstack_keymanager_container_v1.tls_1.container_ref}"
}
```

### Container with the ACL

~> **Note** Only read ACLs are supported

```hcl
resource "openstack_keymanager_container_v1" "tls_1" {
  name = "tls"
  type = "certificate"

  secret_refs {
    name       = "certificate"
    secret_ref = "${openstack_keymanager_secret_v1.certificate_1.secret_ref}"
  }

  secret_refs {
    name       = "private_key"
    secret_ref = "${openstack_keymanager_secret_v1.private_key_1.secret_ref}"
  }

  secret_refs {
    name       = "intermediates"
    secret_ref = "${openstack_keymanager_secret_v1.intermediate_1.secret_ref}"
  }

  acl {
    read {
      project_access = false
      users = [
        "userid1",
        "userid2",
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 KeyManager client.
    A KeyManager client is needed to create a container. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    V1 container.

* `name` - (Optional) Human-readable name for the Container. Does not have
    to be unique.

* `type` - (Required) Used to indicate the type of container. Must be one of `generic`, `rsa` or `certificate`.

* `secret_refs` - (Optional) A set of dictionaries containing references to secrets. The structure is described
    below.

* `acl` - (Optional) Allows to control an access to a container. Currently only
  the `read` operation is supported. If not specified, the container is
  accessible project wide. The `read` structure is described below.

The `secret_refs` block supports:

* `name` - (Optional) The name of the secret reference. The reference names must correspond the container type, more details are available [here](https://docs.openstack.org/barbican/stein/api/reference/containers.html).

* `secret_ref` - (Required) The secret reference / where to find the secret, URL.

The `acl` `read` block supports:

* `project_access` - (Optional) Whether the container is accessible project wide.
  Defaults to `true`.

* `users` - (Optional) The list of user IDs, which are allowed to access the
  container, when `project_access` is set to `false`.

* `created_at` - (Computed) The date the container ACL was created.

* `updated_at` - (Computed) The date the container ACL was last updated.

## Attributes Reference

The following attributes are exported:

* `container_ref` - The container reference / where to find the container.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `type` - See Argument Reference above.
* `secret_refs` - See Argument Reference above.
* `acl` - See Argument Reference above.
* `creator_id` - The creator of the container.
* `status` - The status of the container.
* `created_at` - The date the container was created.
* `updated_at` - The date the container was last updated.
* `consumers` - The list of the container consumers. The structure is described below.

The `consumers` block supports:

* `name` - The name of the consumer.

* `url` - The consumer URL.

## Import

Containers can be imported using the container id (the last part of the container reference), e.g.:

```
$ terraform import openstack_keymanager_container_v1.container_1 0c6cd26a-c012-4d7b-8034-057c0f1c2953
```
