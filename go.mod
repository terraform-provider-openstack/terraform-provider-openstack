module github.com/terraform-provider-openstack/terraform-provider-openstack

go 1.14

require (
    // TODO(iokiwi): waiting for new release of gophercloud
	// https://github.com/gophercloud/gophercloud/compare/v0.15.0...master
	github.com/gophercloud/gophercloud v0.15.0
	github.com/gophercloud/utils v0.0.0-20201101202656-8677e053dcf1
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/stretchr/testify v1.4.0
	gopkg.in/yaml.v2 v2.3.0
)
