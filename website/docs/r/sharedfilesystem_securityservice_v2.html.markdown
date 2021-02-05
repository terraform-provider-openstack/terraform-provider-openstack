---
layout: "openstack"
page_title: "OpenStack: sharedfilesystem_securityservice_v2"
sidebar_current: "docs-openstack-resource-sharedfilesystem-securityservice-v2"
description: |-
  Configure a Shared File System security service.
---

# sharedfilesystem\_securityservice\_v2

Use this resource to configure a security service.

~> **Note:** All arguments including the security service password will be
stored in the raw state as plain-text. [Read more about sensitive data in
state](/docs/state/sensitive-data.html).

A security service stores configuration information for clients for
authentication and authorization (AuthN/AuthZ). For example, a share server
will be the client for an existing service such as LDAP, Kerberos, or
Microsoft Active Directory.

Minimum supported Manila microversion is 2.7.

## Example Usage

```hcl
resource "openstack_sharedfilesystem_securityservice_v2" "securityservice_1" {
  name        = "security"
  description = "created by terraform"
  type        = "active_directory"
  server      = "192.168.199.10"
  dns_ip      = "192.168.199.10"
  domain      = "example.com"
  ou          = "CN=Computers,DC=example,DC=com"
  user        = "joinDomainUser"
  password    = "s8cret"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Shared File System client.
    A Shared File System client is needed to create a security service. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    security service.

* `name` - (Optional) The name of the security service. Changing this updates the name
    of the existing security service.

* `description` - (Optional) The human-readable description for the security service.
    Changing this updates the description of the existing security service.

* `type` - (Required) The security service type - can either be active\_directory,
    kerberos or ldap.  Changing this updates the existing security service.

* `dns_ip` - (Optional) The security service DNS IP address that is used inside the
    tenant network.

* `ou` - (Optional) The security service ou. An organizational unit can be added to
    specify where the share ends up. New in Manila microversion 2.44.

* `user` - (Optional) The security service user or group name that is used by the
    tenant.

* `password` - (Optional) The user password, if you specify a user.

* `domain` - (Optional) The security service domain.

* `server` - (Optional) The security service host name or IP address.

## Attributes Reference

* `id` - The unique ID for the Security Service.
* `region` - See Argument Reference above.
* `project_id` - The owner of the Security Service.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `type` - See Argument Reference above.
* `dns_ip` - See Argument Reference above.
* `ou` - See Argument Reference above.
* `user` - See Argument Reference above.
* `password` - See Argument Reference above.
* `domain` - See Argument Reference above.
* `server` - See Argument Reference above.

## Import

This resource can be imported by specifying the ID of the security service:

```
$ terraform import openstack_sharedfilesystem_securityservice_v2.securityservice_1 <id>
```
