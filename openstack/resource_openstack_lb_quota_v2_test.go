package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func TestAccLBQuotaV2_basic(t *testing.T) {
	var project projects.Project

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckLB(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBQuotaV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "loadbalancer", "1"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "listener", "2"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "member", "3"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "pool", "4"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "health_monitor", "5"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "l7_rule", "40"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "l7_policy", "41"),
				),
			},
			{
				Config: testAccLBQuotaV2Update1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "loadbalancer", "6"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "listener", "7"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "member", "8"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "pool", "9"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "health_monitor", "10"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "l7_rule", "42"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "l7_policy", "43"),
				),
			},
			{
				Config: testAccLBQuotaV2Update2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "loadbalancer", "11"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "listener", "12"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "member", "13"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "pool", "14"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "health_monitor", "15"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "l7_rule", "-1"),
					resource.TestCheckResourceAttr(
						"openstack_lb_quota_v2.quota_1", "l7_policy", "-1"),
				),
			},
		},
	})
}

const testAccLBQuotaV2Basic = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_lb_quota_v2" "quota_1" {
  project_id          = "${openstack_identity_project_v3.project_1.id}"
  loadbalancer        = 1
  listener            = 2
  member              = 3
  pool                = 4
  health_monitor      = 5
  l7_rule             = 40
  l7_policy           = 41
}
`

const testAccLBQuotaV2Update1 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_lb_quota_v2" "quota_1" {
  project_id          = "${openstack_identity_project_v3.project_1.id}"
  loadbalancer        = 6
  listener            = 7
  member              = 8
  pool                = 9
  health_monitor      = 10
  l7_rule             = 42
  l7_policy           = 43
}
`

const testAccLBQuotaV2Update2 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_lb_quota_v2" "quota_1" {
  project_id          = "${openstack_identity_project_v3.project_1.id}"
  loadbalancer        = 11
  listener            = 12
  member              = 13
  pool                = 14
  health_monitor      = 15
  l7_rule             = -1
  l7_policy           = -1
}
`
