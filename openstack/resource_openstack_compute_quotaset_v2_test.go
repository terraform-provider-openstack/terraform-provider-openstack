package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeQuotasetV2_basic(t *testing.T) {
	var (
		project  projects.Project
		quotaset quotasets.QuotaSet
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeQuotasetV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckComputeQuotasetV2Exists("openstack_compute_quotaset_v2.quotaset_1", &quotaset),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "fixed_ips", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "floating_ips", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "injected_file_content_bytes", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "injected_file_path_bytes", "1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "injected_files", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "key_pairs", "1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "metadata_items", "1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "ram", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "security_group_rules", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "security_groups", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "cores", "1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "instances", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "server_groups", "1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "server_group_members", "1"),
				),
			},
			{
				Config: testAccComputeQuotasetV2_update_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckComputeQuotasetV2Exists("openstack_compute_quotaset_v2.quotaset_1", &quotaset),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "fixed_ips", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "floating_ips", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "injected_file_content_bytes", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "injected_file_path_bytes", "3"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "injected_files", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "key_pairs", "3"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "metadata_items", "3"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "ram", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "security_group_rules", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "security_groups", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "cores", "3"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "instances", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "server_groups", "3"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "server_group_members", "3"),
				),
			},
			{
				Config: testAccComputeQuotasetV2_update_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckComputeQuotasetV2Exists("openstack_compute_quotaset_v2.quotaset_1", &quotaset),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "fixed_ips", "5"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "floating_ips", "5"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "injected_file_content_bytes", "5"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "injected_file_path_bytes", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "injected_files", "5"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "key_pairs", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "metadata_items", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "ram", "5"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "security_group_rules", "5"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "security_groups", "5"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "cores", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "instances", "5"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "server_groups", "4"),
					resource.TestCheckResourceAttr(
						"openstack_compute_quotaset_v2.quotaset_1", "server_group_members", "4"),
				),
			},
		},
	})
}

func testAccCheckComputeQuotasetV2Exists(n string, quotaset *quotasets.QuotaSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.BlockStorageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
		}

		found, err := quotasets.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Quotaset not found")
		}

		*quotaset = *found

		return nil
	}
}

const testAccComputeQuotasetV2_basic = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_compute_quotaset_v2" "quotaset_1" {
  project_id                  = "${openstack_identity_project_v3.project_1.id}"
  fixed_ips                   = 2
  floating_ips                = 2
  injected_file_content_bytes = 2
  injected_file_path_bytes    = 1
  injected_files              = 2
  key_pairs                   = 1
  metadata_items              = 1
  ram                         = 2
  security_group_rules        = 2
  security_groups             = 2
  cores                       = 1
  instances                   = 2
  server_groups               = 1
  server_group_members        = 1
}
`

const testAccComputeQuotasetV2_update_1 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_compute_quotaset_v2" "quotaset_1" {
  project_id           = "${openstack_identity_project_v3.project_1.id}"
  fixed_ips                   = 4
  floating_ips                = 4
  injected_file_content_bytes = 4
  injected_file_path_bytes    = 3
  injected_files              = 4
  key_pairs                   = 3
  metadata_items              = 3
  ram                         = 4
  security_group_rules        = 4
  security_groups             = 4
  cores                       = 3
  instances                   = 4
  server_groups               = 3
  server_group_members        = 3
}
`

const testAccComputeQuotasetV2_update_2 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_compute_quotaset_v2" "quotaset_1" {
  project_id           = "${openstack_identity_project_v3.project_1.id}"
  fixed_ips                   = 5
  floating_ips                = 5
  injected_file_content_bytes = 5
  injected_file_path_bytes    = 4
  injected_files              = 5
  key_pairs                   = 4
  metadata_items              = 4
  ram                         = 5
  security_group_rules        = 5
  security_groups             = 5
  cores                       = 4
  instances                   = 5
  server_groups               = 4
  server_group_members        = 4
}
`
