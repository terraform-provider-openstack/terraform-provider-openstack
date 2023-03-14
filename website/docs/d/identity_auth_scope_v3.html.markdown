---
layout: "openstack"
page_title: "OpenStack: openstack_identity_auth_scope_v3"
sidebar_current: "docs-openstack-datasource-identity-auth-scope-v3"
description: |-
  Get authentication information from the current authenticated scope.
---

# openstack\_identity\_auth\_scope\_v3

Use this data source to get authentication information about the current
auth scope in use. This can be used as self-discovery or introspection of
the username or project name currently in use as well as the service catalog.

~> **Important Security Notice** While the `set_token_id` is `true` this data
source will store an *unencrypted* session token in your Terraform state file.
**Use of this data source with `set_token_id = true` in production deployments
is *not* recommended**.
[Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage

### Simple

```hcl
data "openstack_identity_auth_scope_v3" "scope" {
  name = "my_scope"
}
```

To find the the public object storage endpoint for "region1" as listed in the
service catalog:

```hcl
locals {
  object_store_service    = [for entry in data.openstack_identity_auth_scope_v3.scope.service_catalog:
                                 entry if entry.type=="object-store"][0]
  object_store_endpoint   = [for endpoint in local.object_store_service.endpoints:
                                 endpoint if (endpoint.interface=="public" && endpoint.region=="region1")][0]
  object_store_public_url = local.object_store_endpoint.url
}
```

### In a combination with an http data source provider

See [http](/providers/hashicorp/http/latest/docs/data-sources/http) provider for reference.

```hcl
data "openstack_identity_auth_scope_v3" "scope" {
  name = "my_scope"
}
```

```hcl
locals {
  object_store_service    = [for entry in data.openstack_identity_auth_scope_v3.scope.service_catalog:
                                 entry if entry.type=="object-store"][0]
  object_store_endpoint   = [for endpoint in local.object_store_service.endpoints:
                                 endpoint if (endpoint.interface=="public" && endpoint.region=="region1")][0]
  object_store_public_url = local.object_store_endpoint.url
}

data "http" "example" {
  url = local.object_store_public_url

  request_headers = {
    "Accept"       = "application/json"
    "X-Auth-Token" = data.openstack_identity_auth_scope_v3.scope.token_id
  }
}

# print object storage containers in JSON format
output "containers" {
  value = data.http.example.response_body
}
```

## Argument Reference

* `name` - (Required) The name of the scope. This is an arbitrary name which is
  only used as a unique identifier so an actual token isn't used as the ID.

* `region` - (Optional) The region in which to obtain the V3 Identity client.
  A Identity client is needed to retrieve tokens IDs. If omitted, the
  `region` argument of the provider is used.

* `set_token_id` - (Optional) A boolean argument that determines whether to
  export the current auth scope token ID. When set to `true`, the `token_id`
  attribute will contain an unencrypted token that can be used for further API
  calls. **Warning**: please note that the leaked token may allow unauthorized
  access to other OpenStack services within the current auth scope, so use this
  option with caution.

## Attributes Reference

`id` is set to the name given to the scope. In addition, the following attributes
are exported:

* `user_name` - The username of the scope.
* `user_id` - The user ID the of the scope.
* `user_domain_name` - The domain name of the user.
* `user_domain_id` - The domain ID of the user.
* `domain_name` - The domain name of the scope.
* `domain_id` - The domain ID of the scope.
* `project_name` - The project name of the scope.
* `project_id` - The project ID of the scope.
* `project_domain_name` - The domain name of the project.
* `project_domain_id` - The domain ID of the project.
* `token_id` - The token ID of the scope.
* `roles` - A list of roles in the current scope. See reference below.
* `service_catalog` - A list of service catalog entries returned with the token.

The `roles` block contains:

* `role_id` - The ID of the role.
* `role_name` - The name of the role.

The `service_catalog` block contains:

* `id` - The ID of the service.
* `name` - The name of the service.
* `type` - The type of the service.
* `endpoints` - A list of endpoints for the service.

The `endpoints` block contains:

* `id` - The ID of the endpoint.
* `region` - The region of the endpoint.
* `region_id` - The region ID of the endpoint.
* `interface` - The interface of the endpoint.
* `url` - The URL of the endpoint.
