package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDNSQuotaV2_basic(t *testing.T) {
	var project projects.Project

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckDNS(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3ProjectDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSQuotaV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists(t.Context(), "openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "api_export_size", "5"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "recordset_records", "6"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "zone_records", "7"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "zone_recordsets", "8"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "zones", "9"),
				),
			},
			{
				Config: testAccDNSQuotaV2Update1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists(t.Context(), "openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "api_export_size", "9"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "recordset_records", "8"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "zone_records", "7"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "zone_recordsets", "6"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "zones", "5"),
				),
			},
			{
				Config: testAccDNSQuotaV2Update2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists(t.Context(), "openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "api_export_size", "11"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "recordset_records", "12"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "zone_records", "13"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "zone_recordsets", "14"),
					resource.TestCheckResourceAttr(
						"openstack_dns_quota_v2.quota_1", "zones", "15"),
				),
			},
		},
	})
}

const testAccDNSQuotaV2Basic = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_dns_quota_v2" "quota_1" {
  project_id        = openstack_identity_project_v3.project_1.id
  api_export_size   = 5
  recordset_records = 6
  zone_records      = 7
  zone_recordsets   = 8
  zones             = 9
}
`

const testAccDNSQuotaV2Update1 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_dns_quota_v2" "quota_1" {
  project_id        = openstack_identity_project_v3.project_1.id
  api_export_size   = 9
  recordset_records = 8
  zone_records      = 7
  zone_recordsets   = 6
  zones             = 5
}
`

const testAccDNSQuotaV2Update2 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_dns_quota_v2" "quota_1" {
  project_id        = openstack_identity_project_v3.project_1.id
  api_export_size   = 11
  recordset_records = 12
  zone_records      = 13
  zone_recordsets   = 14
  zones             = 15
}
`
