package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/rbacpolicies"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkingV2RBACPolicy_basic(t *testing.T) {
	var rbac rbacpolicies.RBACPolicy

	var project projects.Project

	var network networks.Network

	projectOneName := "ACCPTTEST-" + acctest.RandString(5)

	projectTwoName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RBACPolicyDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RBACPolicyBasic(projectOneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists(t.Context(), "openstack_identity_project_v3.project_1", &project),
					testAccCheckNetworkingV2NetworkExists(t.Context(), "openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2RBACPolicyExists(t.Context(), "openstack_networking_rbac_policy_v2.rbac_policy_1", &rbac),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "action", "access_as_shared"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "action", (*string)(&rbac.Action)),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_id", &network.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_type", &rbac.ObjectType),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_type", "network"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "target_tenant", &project.ID),
				),
			},
			{
				Config: testAccNetworkingV2RBACPolicyUpdate(projectTwoName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists(t.Context(), "openstack_identity_project_v3.project_2", &project),
					testAccCheckNetworkingV2NetworkExists(t.Context(), "openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2RBACPolicyExists(t.Context(), "openstack_networking_rbac_policy_v2.rbac_policy_1", &rbac),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "action", "access_as_shared"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "action", (*string)(&rbac.Action)),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_id", &network.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_type", &rbac.ObjectType),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_type", "network"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "target_tenant", &project.ID),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RBACPolicyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_networking_rbac_policy_v2" {
				continue
			}

			_, err := rbacpolicies.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Project still exists")
			}
		}

		return nil
	}
}

func testAccCheckNetworkingV2RBACPolicyExists(ctx context.Context, n string, rbac *rbacpolicies.RBACPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		found, err := rbacpolicies.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Project not found")
		}

		*rbac = *found

		return nil
	}
}

func testAccNetworkingV2RBACPolicyBasic(projectName string) string {
	return fmt.Sprintf(`
    resource "openstack_identity_project_v3" "project_1" {
      name        = "%s"
      description = "A project"
    }

    resource "openstack_networking_network_v2" "network_1" {
      name           = "network_1"
      admin_state_up = "false"
    }

    resource "openstack_networking_rbac_policy_v2" "rbac_policy_1" {
      action        = "access_as_shared"
      object_id     = openstack_networking_network_v2.network_1.id
      object_type   = "network"
      target_tenant = openstack_identity_project_v3.project_1.id
    }
  `, projectName)
}

func testAccNetworkingV2RBACPolicyUpdate(projectName string) string {
	return fmt.Sprintf(`
    resource "openstack_identity_project_v3" "project_2" {
      name        = "%s"
      description = "The second project"
    }

    resource "openstack_networking_network_v2" "network_1" {
      name           = "network_1"
      admin_state_up = "false"
    }

    resource "openstack_networking_rbac_policy_v2" "rbac_policy_1" {
      action        = "access_as_shared"
      object_id     = openstack_networking_network_v2.network_1.id
      object_type   = "network"
      target_tenant = openstack_identity_project_v3.project_2.id
    }
  `, projectName)
}
