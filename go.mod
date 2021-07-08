module github.com/terraform-provider-openstack/terraform-provider-openstack

go 1.15

replace github.com/gophercloud/gophercloud => github.com/randomswdev/gophercloud v0.17.1-0.20210609140925-9815af997180

require (
	github.com/gophercloud/gophercloud v0.17.1-0.20210429172743-0fc2c97ff6da
	github.com/gophercloud/utils v0.0.0-20210216074907-f6de111f2eae
	github.com/hashicorp/terraform-plugin-sdk v1.17.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
)
