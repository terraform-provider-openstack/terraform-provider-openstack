package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccComputeV2FlavorAccess_basic(t *testing.T) {
	var flavor flavors.Flavor

	flavorName := "ACCPTTEST-" + acctest.RandString(5)

	var project projects.Project

	projectName := "ACCPTTEST-" + acctest.RandString(5)

	var flavorAccess flavors.FlavorAccess

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2FlavorAccessDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FlavorAccessBasic(flavorName, projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists(t.Context(), "openstack_identity_project_v3.project_1", &project),
					testAccCheckComputeV2FlavorExists(t.Context(), "openstack_compute_flavor_v2.flavor_1", &flavor),
					testAccCheckComputeV2FlavorAccessExists(t.Context(), "openstack_compute_flavor_access_v2.access_1", &flavorAccess),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_flavor_access_v2.access_1", "flavor_id", &flavor.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_flavor_access_v2.access_1", "tenant_id", &project.ID),
				),
			},
		},
	})
}

func testAccCheckComputeV2FlavorAccessDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		computeClient, err := config.ComputeV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_compute_flavor_access_v2" {
				continue
			}

			fid, tid, err := parseComputeFlavorAccessID(rs.Primary.ID)
			if err != nil {
				return err
			}

			pager := flavors.ListAccesses(computeClient, fid)

			err = pager.EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
				accessList, err := flavors.ExtractAccesses(page)
				if err != nil {
					return false, err
				}

				for _, a := range accessList {
					if a.TenantID == tid {
						return false, errors.New("Flavor Access still exists")
					}
				}

				return true, nil
			})
			if err != nil {
				if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
					return nil
				}

				return err
			}
		}

		return nil
	}
}

func testAccCheckComputeV2FlavorAccessExists(ctx context.Context, n string, access *flavors.FlavorAccess) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		computeClient, err := config.ComputeV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %w", err)
		}

		fid, tid, err := parseComputeFlavorAccessID(rs.Primary.ID)
		if err != nil {
			return err
		}

		pager := flavors.ListAccesses(computeClient, fid)
		err = pager.EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
			accessList, err := flavors.ExtractAccesses(page)
			if err != nil {
				return false, err
			}

			for _, acc := range accessList {
				a := acc
				if a.TenantID == tid {
					access = &a

					return false, nil
				}
			}

			return true, nil
		})

		return err
	}
}

func testAccComputeV2FlavorAccessBasic(flavorName, tenantName string) string {
	return fmt.Sprintf(`
    resource "openstack_compute_flavor_v2" "flavor_1" {
      name = "%s"
      ram = 512
      vcpus = 1
      disk = 5

      is_public = false
    }

    resource "openstack_identity_project_v3" "project_1" {
      name = "%s"
    }

    resource "openstack_compute_flavor_access_v2" "access_1" {
      flavor_id = openstack_compute_flavor_v2.flavor_1.id
      tenant_id = openstack_identity_project_v3.project_1.id
    }
    `, flavorName, tenantName)
}
