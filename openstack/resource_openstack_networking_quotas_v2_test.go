package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetworkingQuotasV2_basic(t *testing.T) {
	var project projects.Project

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingQuotasV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "floatingip", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "network", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "port", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "rbac_policy", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "router", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "security_group", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "security_group_rule", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "subnet", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "subnetpool", "1"),
				),
			},
			{
				Config: testAccNetworkingQuotasV2_update_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "floatingip", "3"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "network", "3"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "port", "4"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "rbac_policy", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "router", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "security_group", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "security_group_rule", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "subnet", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "subnetpool", "1"),
				),
			},
			{
				Config: testAccNetworkingQuotasV2_update_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "floatingip", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "network", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "port", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "rbac_policy", "4"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "router", "4"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "security_group", "3"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "security_group_rule", "3"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "subnet", "3"),
					resource.TestCheckResourceAttr(
						"openstack_networking_quotas_v2.quotas_1", "subnetpool", "3"),
				),
			},
		},
	})
}

const testAccNetworkingQuotasV2_basic = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_networking_quotas_v2" "quotas_1" {
  project_id          = "${openstack_identity_project_v3.project_1.id}"
  floatingip          = 2
  network             = 2
  port                = 2
  rbac_policy         = 1
  router              = 2
  security_group      = 2
  security_group_rule = 2
  subnet              = 1
  subnetpool          = 1
}
`

const testAccNetworkingQuotasV2_update_1 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_networking_quotas_v2" "quotas_1" {
  project_id          = "${openstack_identity_project_v3.project_1.id}"
  floatingip          = 3
  network             = 3
  port                = 4
  rbac_policy         = 1
  router              = 2
  security_group      = 2
  security_group_rule = 2
  subnet              = 1
  subnetpool          = 1
}
`

const testAccNetworkingQuotasV2_update_2 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_networking_quotas_v2" "quotas_1" {
  project_id          = "${openstack_identity_project_v3.project_1.id}"
  floatingip          = 2
  network             = 2
  port                = 2
  rbac_policy         = 4
  router              = 4
  security_group      = 3
  security_group_rule = 3
  subnet              = 3
  subnetpool          = 3
}
`
