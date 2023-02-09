---
layout: "openstack"
page_title: "OpenStack: openstack_identity_ec2_credential_v3"
sidebar_current: "docs-openstack-resource-identity-ec2-credential-v3"
description: |-
  Manages a V3 EC2 Credential resource within OpenStack Keystone.
---

# openstack\_identity\_ec2\_credential\_v3

Manages a V3 EC2 Credential resource within OpenStack Keystone.
EC2 credentials in OpenStack are used to access S3 compatible Swift/RadosGW
endpoints or to authenticate against Keystone.

~> **Note:** All arguments including the EC2 credential access key and secret
will be stored in the raw state as plain-text. [Read more about sensitive data
in state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage

### EC2 credential in current project scope

```hcl
resource "openstack_identity_ec2_credential_v3" "ec2_key1" {}
```

### EC2 credential in pre-defined project scope

This allows administrative users to create EC2 credentials for a scope different
from the current auth scope.

```hcl
resource "openstack_identity_ec2_credential_v3" "ec2_key1" {
    project_id = "f7ac731cc11f40efbc03a9f9e1d1d21f"
}
```

## Arguments Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
   If omitted, the `region` argument of the provider is used. Changing this
   creates a new EC2 credential.
* `project_id` - (Optional) The ID of the project the EC2 credential is created
   for and that authentication requests using this EC2 credential will
   be scoped to. Only administrative users can specify a project ID different
   from the current auth scope.
* `user_id` - (Optional) The ID of the user the EC2 credential is created for.
   Only administrative users can specify a user ID different from the current
   auth scope.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `access` - contains an EC2 credential access UUID
* `secret` - contains an EC2 credential secret UUID
* `user_id` - contains a User ID of the EC2 credential owner
* `project_id` - contains an EC2 credential project scope
* `trust_id` - contains an EC2 credential trust ID scope

## Import

EC2 Credentials can be imported using the `access`, e.g.

```
$ terraform import openstack_identity_ec2_credential_v3.ec2_cred_1 2d0ac4a2f81b4b0f9513ee49e780647d
```
