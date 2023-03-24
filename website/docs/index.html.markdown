---
layout: "openstack"
page_title: "Provider: OpenStack"
sidebar_current: "docs-openstack-index"
description: |-
  The OpenStack provider is used to interact with the many resources supported by OpenStack. The provider needs to be configured with the proper credentials before it can be used.
---

# OpenStack Provider

The OpenStack provider is used to interact with the
many resources supported by OpenStack. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Define required providers
terraform {
required_version = ">= 0.14.0"
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.48.0"
    }
  }
}

# Configure the OpenStack Provider
provider "openstack" {
  user_name   = "admin"
  tenant_name = "admin"
  password    = "pwd"
  auth_url    = "http://myauthurl:5000/v2.0"
  region      = "RegionOne"
}

# Create a web server
resource "openstack_compute_instance_v2" "test-server" {
  # ...
}
```

## Configuration Reference

The following arguments are supported:

* `auth_url` - (Optional; required if `cloud` is not specified) The Identity
  authentication URL. If omitted, the `OS_AUTH_URL` environment variable is used.

* `cloud` - (Optional; required if `auth_url` is not specified) An entry in a
  `clouds.yaml` file. See the OpenStack `openstacksdk`
  [documentation](https://docs.openstack.org/openstacksdk/latest/user/config/configuration.html)
  for more information about `clouds.yaml` files. If omitted, the `OS_CLOUD`
  environment variable is used.

* `region` - (Optional) The region of the OpenStack cloud to use. If omitted,
  the `OS_REGION_NAME` environment variable is used. If `OS_REGION_NAME` is
  not set, then no region will be used. It should be possible to omit the
  region in single-region OpenStack environments, but this behavior may vary
  depending on the OpenStack environment being used.

* `user_name` - (Optional) The Username to login with. If omitted, the
  `OS_USERNAME` environment variable is used.

* `user_id` - (Optional) The User ID to login with. If omitted, the
  `OS_USER_ID` environment variable is used.

* `application_credential_id` - (Optional) (Identity v3 only) The ID of an
  application credential to authenticate with. An
  `application_credential_secret` has to bet set along with this parameter.
  If omitted, the `OS_APPLICATION_CREDENTIAL_ID` environment variable is used.

* `application_credential_name` - (Optional) (Identity v3 only) The name of an
  application credential to authenticate with. Requires `user_id`, or
  `user_name` and `user_domain_name` (or `user_domain_id`) to be set.
  If omitted, the `OS_APPLICATION_CREDENTIAL_NAME` environment variable is used.

* `application_credential_secret` - (Optional) (Identity v3 only) The secret of
  an application credential to authenticate with. Required by
  `application_credential_id` or `application_credential_name`.
  If omitted, the `OS_APPLICATION_CREDENTIAL_SECRET` environment variable is used.

* `tenant_id` - (Optional) The ID of the Tenant (Identity v2) or Project
  (Identity v3) to login with. If omitted, the `OS_TENANT_ID` or
  `OS_PROJECT_ID` environment variables are used.

* `tenant_name` - (Optional) The Name of the Tenant (Identity v2) or Project
  (Identity v3) to login with. If omitted, the `OS_TENANT_NAME` or
  `OS_PROJECT_NAME` environment variable are used.

* `password` - (Optional) The Password to login with. If omitted, the
  `OS_PASSWORD` environment variable is used.

* `token` - (Optional; Required if not using `user_name` and `password`)
  A token is an expiring, temporary means of access issued via the Keystone
  service. By specifying a token, you do not have to specify a username/password
  combination, since the token was already created by a username/password out of
  band of Terraform. If omitted, the `OS_TOKEN` or `OS_AUTH_TOKEN` environment
  variables are used.

* `user_domain_name` - (Optional) The domain name where the user is located. If
  omitted, the `OS_USER_DOMAIN_NAME` environment variable is checked.

* `user_domain_id` - (Optional) The domain ID where the user is located. If
  omitted, the `OS_USER_DOMAIN_ID` environment variable is checked.

* `project_domain_name` - (Optional) The domain name where the project is
  located. If omitted, the `OS_PROJECT_DOMAIN_NAME` environment variable is
  checked.

* `project_domain_id` - (Optional) The domain ID where the project is located
  If omitted, the `OS_PROJECT_DOMAIN_ID` environment variable is checked.

* `domain_id` - (Optional) The ID of the Domain to scope to (Identity v3). If
  omitted, the `OS_DOMAIN_ID` environment variable is checked.

* `domain_name` - (Optional) The Name of the Domain to scope to (Identity v3).
  If omitted, the following environment variables are checked (in this order):
  `OS_DOMAIN_NAME`.

* `default_domain` - (Optional) The ID of the Domain to scope to if no other
  domain is specified (Identity v3). If omitted, the environment variable
  `OS_DEFAULT_DOMAIN` is checked or a default value of "default" will be
  used.

* `system_scope` - (Optional) Set to `true` to enable system scoped authorization. If omitted, the `OS_SYSTEM_SCOPE` environment variable is used.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `OS_INSECURE` environment variable is used.

* `cacert_file` - (Optional) Specify a custom CA certificate when communicating
  over SSL. You can specify either a path to the file or the contents of the
  certificate. If omitted, the `OS_CACERT` environment variable is used.

* `cert` - (Optional) Specify client certificate file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the certificate. If omitted the `OS_CERT` environment variable is used.

* `key` - (Optional) Specify client private key file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the key. If omitted the `OS_KEY` environment variable is used.

* `endpoint_type` - (Optional) Specify which type of endpoint to use from the
  service catalog. It can be set using the `OS_ENDPOINT_TYPE` environment
  variable. If not set, public endpoints is used.

* `endpoint_overrides` - (Optional) A set of key/value pairs that can
  override an endpoint for a specified OpenStack service. Setting an override
  requires you to specify the full and complete endpoint URL. This might
  also invalidate any region you have set, too. Please see below for more details.
  Please use this at your own risk.

* `swauth` - (Optional) Set to `true` to authenticate against Swauth, a
  Swift-native authentication system. If omitted, the `OS_SWAUTH` environment
  variable is used. You must also set `username` to the Swauth/Swift username
  such as `username:project`. Set the `password` to the Swauth/Swift key.
  Finally, set `auth_url` as the location of the Swift service. Note that this
  will only work when used with the OpenStack Object Storage resources.

* `use_octavia` - (Optional) If set to `true`, API requests will go the Load Balancer
  service (Octavia) instead of the Networking service (Neutron).
  If omitted, the `OS_USE_OCTAVIA` environment variable is checked.

* `disable_no_cache_header` - (Optional) If set to `true`, the HTTP
  `Cache-Control: no-cache` header will not be added by default to all API requests.
  If omitted this header is added to all API requests to force HTTP caches (if any)
  to go upstream instead of serving cached responses.

* `delayed_auth` - (Optional) If set to `false`, OpenStack authorization will be perfomed,
  every time the service provider client is called. Defaults to `true`.
  If omitted, the `OS_DELAYED_AUTH` environment variable is checked.

* `allow_reauth` - (Optional) If set to `false`, OpenStack authorization won't be
  perfomed automatically, if the initial auth token get expired. Defaults to `true`.
  If omitted, the `OS_ALLOW_REAUTH` environment variable is checked.

* `max_retries` - (Optional) If set to a value greater than 0, the OpenStack
  client will retry failed HTTP connections and Too Many Requests (429 code)
  HTTP responses with a `Retry-After` header within the specified value.

## Overriding Service API Endpoints

There might be a situation in which you want or need to override an API endpoint
rather than use the endpoint which was returned to you in the service catalog.
You can do this by configuring the `endpoint_overrides` argument in the provider
configuration:

```hcl
provider "openstack" {

  endpoint_overrides = {
    "network"  = "https://example.com:9696/v2.0/"
    "volumev2" = "https://volumes.example.com:8776/v2/3eb25ae78e7b42d68276e9bca66c8e44/"
  }

}
```

Note how each URL ends in a "/" and the `volumev2` service includes the
tenant/project UUID. You must make sure you specify the full and complete
endpoint URL for this to work.

The service keys are the standard service entries used in the OpenStack
Identity/Keystone service catalog. This provider supports:

* `compute`: Compute / Nova v2
* `container-infra`: Container Infra / Magnum v1
* `database`: Database / Trove v1
* `dns`: DNS / Designate v2
* `identity`: Identity / Keystone v3
* `image`: Image / Glance v2
* `network`: Networking / Neutron v2
* `object-store`: Object Storage / Swift v1
* `octavia`: Load Balancing as a Service / Octavia v2
* `sharev2`: Shared Filesystem / Manila v2
* `volume`: Block Storage / Cinder v1
* `volumev2`: Block Storage / Cinder v2
* `volumev3`: Block Storage / Cinder v3

Please use this feature at your own risk. If you are unsure about needing
to override an endpoint, you most likely do not need to override one.

## Additional Logging

This provider has the ability to log all HTTP requests and responses between
Terraform and the OpenStack cloud which is useful for troubleshooting and
debugging.

To enable these logs, set the `OS_DEBUG` environment variable to `1` along
with the usual `TF_LOG=DEBUG` environment variable:

```shell
$ OS_DEBUG=1 TF_LOG=DEBUG terraform apply
```

If you submit these logs with a bug report, please ensure any sensitive
information has been scrubbed first!

## OpenStack Releases and Versions

This provider aims to support "vanilla" OpenStack. This means that we do all
testing and development using the official upstream OpenStack code. If your
OpenStack environment has patches or modifications, we do our best to
accommodate these modifications, but we can't guarantee this.

We try to support _all_ releases of OpenStack when we can. If your OpenStack
cloud is running an older release, we still should be able to support it.

### Octavia api versioning

[Octavia api](https://docs.openstack.org/api-ref/load-balancer/v2/) is using
minor versions when adding new features and functionality. The required minor
version of each feature are documented on the resource page. When using such a
feature ensure that your Openstack cloud supports the required minor version.
A simple way of checking which minor versions are supported on your Openstack
cloud is the following:

```shell
export OS_TOKEN=`openstack token issue -c id -f value`
curl -s -H "X-Auth-Token: $OS_TOKEN"  "https://example.com:9876/"
```



### Rackspace Compatibility

Using this OpenStack provider with Rackspace is not supported and not
guaranteed to work; however, users have reported success with the
following notes in mind:

* Interacting with instances and networks has been seen to work. Interacting
with all other resources is either untested or known to not work.

* Use your _password_ instead of your Rackspace API KEY.

* To use networks, override the endpoint in your provider configuration

```
provider "openstack" {
  user_name = "your-username"
  password  = "your-password"
  tenant_id = "your-tenant-id"
  region    = "IAD"
  auth_url  = "https://identity.api.rackspacecloud.com/v2.0/"
  endpoint_overrides = {
    "network" = "https://iad.networks.api.rackspacecloud.com/v2.0/"
  }
}
```

* Explicitly define the public and private networks in your
instances as shown below:

```
resource "openstack_compute_instance_v2" "my_instance" {
  name      = "my_instance"
  region    = "DFW"
  image_id  = "fabe045f-43f8-4991-9e6c-5cabd617538c"
  flavor_id = "general1-4"
  key_pair  = "provisioning_key"

  network {
    uuid = "00000000-0000-0000-0000-000000000000"
    name = "public"
  }

  network {
    uuid = "11111111-1111-1111-1111-111111111111"
    name = "private"
  }
}
```

If you try using this provider with Rackspace and run into bugs, you
are welcomed to open a bug report / issue on Github, but please keep
in mind that this is unsupported and the reported bug may not be
able to be fixed.

If you have successfully used this provider with Rackspace and can
add any additional comments, please let us know.

## Testing and Development

Thank you for your interest in further developing the OpenStack provider! Here
are a few notes which should help you get started. If you have any questions or
feel these notes need further details, please open an Issue and let us know.

### Coding and Style

This provider aims to adhere to the coding style and practices used in the
other major Terraform Providers. However, this is not a strict rule.

We're very mindful that not everyone is a full-time developer (most of the
OpenStack Provider contributors aren't!) and we're happy to provide
guidance. Don't be afraid if this is your first contribution to the OpenStack
provider or even your first contribution to an open source project!

### Testing Environment

In order to start fixing bugs or adding features, you need access to an
OpenStack environment. If it is safe to do, you can use a production OpenStack
cloud which you have access to. However, it's usually safer to work in a
development cloud.

[DevStack](https://docs.openstack.org/devstack/latest/) is a quick and easy way
to spin up an OpenStack cloud. All OpenStack services have DevStack plugins so
you can build a DevStack environment to test everything from Nova/Compute to
Designate/DNS.

Fully configuring a DevStack installation is outside the scope of this
document; however, we'll try to provide assistance where we can.

### Gophercloud

This provider uses [Gophercloud](https://github.com/gophercloud/gophercloud)
as the Go OpenStack SDK. All API interaction between this provider and an
OpenStack cloud is done exclusively with Gophercloud.

### Adding a Feature

If you'd like to add a new feature to this provider, it must first be supported
in Gophercloud. If Gophercloud is missing the feature, then it'll first have to
be added there before you can start working on the feature in Terraform.
Fortunately, most of the regular OpenStack Provider contributors also work on
Gophercloud, so we can try to get the feature added quickly.

If the feature is already included in Gophercloud, then you can begin work
directly in the OpenStack provider.

If you have any questions about whether Gophercloud currently supports a
certain feature, please feel free to ask and we can verify for you.

### Fixing a Bug

Similarly, if you find a bug in this provider, the bug might actually be a
Gophercloud bug. If this is the case, then we'll need to get the bug fixed in
Gophercloud first.

However, if the bug is with Terraform itself, then you can begin work directly
in the OpenStack provider.

Again, if you have any questions about whether the bug you're trying to fix is
a Gophercloud but, please ask.

### Dependencies

If you require pulling in changes from an external package, such as Gophercloud,
this provider uses [Go Modules](https://github.com/golang/go/wiki/Modules).

### Acceptance Tests

Acceptance Tests are a crucial part of adding features or fixing a bug. Please
make sure to read the core [testing](https://www.terraform.io/docs/extend/testing/index.html)
documentation for more information about how Acceptance Tests work.

In order to run the Acceptance Tests, you'll need to set the following
environment variables:

* `OS_IMAGE_ID` or `OS_IMAGE_NAME` - a UUID or name of an existing image in
    Glance.

* `OS_FLAVOR_ID` or `OS_FLAVOR_NAME` - an ID or name of an existing flavor.

* `OS_POOL_NAME` - The name of a Floating IP pool. In DevStack, this is
  called `public` and you should set this value to the word `public`.

* `OS_NETWORK_ID` - The UUID of the private network in your test environment.
  In DevStack, this network is called `private` and you should set this value
  to the UUID of the `private` network.

* `OS_EXTGW_ID` - The UUID of the public network in your test environment. In
  DevStack, this network is called `public` and you should set this value to
  the UUID of the `public` network.

The following additional environment variables might be required depending on
the feature or bug you're testing:

* `OS_DB_ENVIRONMENT` - Required if you're working on the `openstack_db_*`
  resources. Set this value to "1" to enable testing these resources.

* `OS_DB_DATASTORE_VERSION` - Required if you need to set a Trove/Database
  datastore version.

* `OS_DB_DATASTORE_TYPE` - Required if you need to set a Trove/Database
  datastore type.

* `OS_DNS_ENVIRONMENT` - Required if you're working on the `openstack_dns_*`
  resources. Set this value to "1" to enable testing these resources.

* `OS_SWIFT_ENVIRONMENT` - Required if you're working on an
  `openstack_objectstorage_*` resource. Set this value to "1" to enable testing
  these resources.

* `OS_LB_ENVIRONMENT` - Required if you're working on the `openstack_lb_*`
  resources. Set this value to "1" to enable testing these resources.

* `OS_FW_ENVIRONMENT` - Required if you're working on the `openstack_fw_*`
  resources. Set this value to "1" to enable testing these resources.

* `OS_VPN_ENVIRONMENT` - Required if your'e working on the `openstack_vpn_*`
  resources. Set this value to "1" to enable testing these resources.

* `OS_SFS_ENVIRONMENT` - Required if your'e working on the `openstack_openstack_sharedfilesystem_*`
  resources. Set this value to "1" to enable testing these resources.

* `OS_HYPERVISOR_HOSTNAME` - Required if you're working on the `openstack_compute_hypervisor_v2`
  data source. Set this value to one valid hypervisor hostname to test this data source.

We recommend only running the acceptance tests related to the feature or bug
you're working on. To do this, run:

```shell
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-openstack
$ make testacc TEST=./openstack TESTARGS="-run=<keyword> -count=1"
```

Where `<keyword>` is the full name or partial name of a test. For example:

```shell
$ make testacc TEST=./openstack TESTARGS="-run=TestAccComputeV2Keypair_basic -count=1"
```

We recommend running tests with logging set to `DEBUG`:

```shell
$ TF_LOG=DEBUG make testacc TEST=./openstack TESTARGS="-run=TestAccComputeV2Keypair_basic -count=1"
```

And you can even enable OpenStack debugging to see the actual HTTP API requests:

```shell
$ TF_LOG=DEBUG OS_DEBUG=1 make testacc TEST=./openstack TESTARGS="-run=TestAccComputeV2Keypair_basic -count=1"
```

### Creating a Pull Request

When you're ready to submit a Pull Request, create a branch, commit your code,
and push to your forked version of `terraform-provider-openstack`:

```shell
$ git remote add my-github-username https://github.com/my-github-username/terraform-provider-openstack
$ git checkout -b my-feature
$ git add .
$ git commit
$ git push -u my-github-username my-feature
```

Then navigate to https://github.com/terraform-provider-openstack/terraform-provider-openstack
and create a Pull Request.

### Testing with GitHub Actions

Once you have created a Pull Request, it will automatically be tested by
various GitHub Actions. GitHub Actions use Devstack to create a clean
vanilla Openstack Cloud. Multiple jobs are triggered to test various
versions of Openstack as well as different services. Testing will take
between 90-120 minutes and you will receive a notification with a test
report when testing has finished.

If there were any failures, check the provided logs.

There are a few reasons for test failures:

1. Your code changes worked in your environment but are not working in a
  different OpenStack environment.

2. Your code changes caused another test to fail.

3. DevStack is having issues.

If you are unable to determine why the failures happened, please ask and
we'll look into the cause.
