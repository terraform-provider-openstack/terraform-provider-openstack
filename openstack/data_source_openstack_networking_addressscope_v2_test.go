package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOpenStackNetworkingAddressScopeV2DataSource_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSource_addressscope,
			},
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSource_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_addressscope_v2.addressscope_1"),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "name", "addressscope_1"),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "ip_version", "4"),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "shared", "false"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingAddressScopeV2DataSource_ipversion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSource_addressscope,
			},
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSource_ipversion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_addressscope_v2.addressscope_1"),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "name", "addressscope_1"),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "ip_version", "4"),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "shared", "false"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingAddressScopeV2DataSource_shared(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSource_addressscope,
			},
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSource_shared,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_addressscope_v2.addressscope_1"),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "name", "addressscope_1"),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "ip_version", "4"),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "shared", "false"),
				),
			},
		},
	})
}

const testAccOpenStackNetworkingAddressScopeV2DataSource_addressscope = `
resource "openstack_networking_addressscope_v2" "addressscope_1" {
  name       = "addressscope_1"
  ip_version = 4
  shared     = false
}`

var testAccOpenStackNetworkingAddressScopeV2DataSource_name = fmt.Sprintf(`
%s

data "openstack_networking_addressscope_v2" "addressscope_1" {
  name = "${openstack_networking_addressscope_v2.addressscope_1.name}"
}
`, testAccOpenStackNetworkingAddressScopeV2DataSource_addressscope)

var testAccOpenStackNetworkingAddressScopeV2DataSource_ipversion = fmt.Sprintf(`
%s

data "openstack_networking_addressscope_v2" "addressscope_1" {
  ip_version = "${openstack_networking_addressscope_v2.addressscope_1.ip_version}"
}
`, testAccOpenStackNetworkingAddressScopeV2DataSource_addressscope)

var testAccOpenStackNetworkingAddressScopeV2DataSource_shared = fmt.Sprintf(`
%s

data "openstack_networking_addressscope_v2" "addressscope_1" {
  shared = "${openstack_networking_addressscope_v2.addressscope_1.shared}"
}
`, testAccOpenStackNetworkingAddressScopeV2DataSource_addressscope)
