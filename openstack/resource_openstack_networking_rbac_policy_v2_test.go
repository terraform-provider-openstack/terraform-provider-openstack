package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/rbacpolicies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

func TestAccNetworkingV2RBACPolicy_basic(t *testing.T) {
	var rbac rbacpolicies.RBACPolicy
	var project projects.Project
	var network networks.Network
	var projectOneName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var projectTwoName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RBACPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RBACPolicyBasic(projectOneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2RBACPolicyExists("openstack_networking_rbac_policy_v2.rbac_policy_1", &rbac),
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
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_2", &project),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2RBACPolicyExists("openstack_networking_rbac_policy_v2.rbac_policy_1", &rbac),
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

func testAccCheckNetworkingV2RBACPolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_rbac_policy_v2" {
			continue
		}

		_, err := rbacpolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Project still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2RBACPolicyExists(n string, rbac *rbacpolicies.RBACPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := rbacpolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Project not found")
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
      object_id     = "${openstack_networking_network_v2.network_1.id}"
      object_type   = "network"
      target_tenant = "${openstack_identity_project_v3.project_1.id}"
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
      object_id     = "${openstack_networking_network_v2.network_1.id}"
      object_type   = "network"
      target_tenant = "${openstack_identity_project_v3.project_2.id}"
    }
  `, projectName)
}

func TestAccNetworkingV2RBACPolicy_search(t *testing.T) {
	var rbac1, rbac2 rbacpolicies.RBACPolicy
	var project1, project2, project3 projects.Project
	var network1, network2, network3 networks.Network
	var projectOneName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var projectTwoName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var projectThreeName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RBACPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RBACPolicy_shared_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network1),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_2", &network2),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_3", &network3),
				),
			},
			{
				Config: testAccNetworkingV2RBACPolicy_shared_update1(projectOneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project1),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network1),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_2", &network2),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_3", &network3),
					testAccCheckNetworkingV2RBACPolicyExists("openstack_networking_rbac_policy_v2.rbac_policy_1", &rbac1),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "action", "access_as_shared"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "action", (*string)(&rbac1.Action)),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_id", &network1.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_type", &rbac1.ObjectType),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_type", "network"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "target_tenant", &project1.ID),
				),
			},
			{
				Config: testAccNetworkingV2RBACPolicy_shared_update2(projectOneName, projectTwoName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project1),
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_2", &project2),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network1),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_2", &network2),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_3", &network3),
					testAccCheckNetworkingV2RBACPolicyExists("openstack_networking_rbac_policy_v2.rbac_policy_1", &rbac1),
					testAccCheckNetworkingV2RBACPolicyExists("openstack_networking_rbac_policy_v2.rbac_policy_2", &rbac2),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "action", "access_as_shared"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "action", (*string)(&rbac1.Action)),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_id", &network1.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_type", &rbac1.ObjectType),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_type", "network"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "target_tenant", &project1.ID),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "action", "access_as_shared"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "action", (*string)(&rbac2.Action)),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "object_id", &network2.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "object_type", &rbac2.ObjectType),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "object_type", "network"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "target_tenant", &project2.ID),
				),
			},
			{
				Config: testAccNetworkingV2RBACPolicy_shared_update3(projectThreeName, projectTwoName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_3", &project3),
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_2", &project2),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network1),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_2", &network2),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_3", &network3),
					testAccCheckNetworkingV2RBACPolicyExists("openstack_networking_rbac_policy_v2.rbac_policy_1", &rbac1),
					testAccCheckNetworkingV2RBACPolicyExists("openstack_networking_rbac_policy_v2.rbac_policy_2", &rbac2),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "action", "access_as_shared"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "action", (*string)(&rbac1.Action)),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_id", &network1.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_type", &rbac1.ObjectType),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "object_type", "network"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_1", "target_tenant", &project3.ID),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "action", "access_as_shared"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "action", (*string)(&rbac2.Action)),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "object_id", &network2.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "object_type", &rbac2.ObjectType),
					resource.TestCheckResourceAttr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "object_type", "network"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_rbac_policy_v2.rbac_policy_2", "target_tenant", &project2.ID),
				),
			},
		},
	})
}

func testAccNetworkingV2RBACPolicy_shared_basic() string {
	return `
    resource "openstack_networking_network_v2" "network_1" {
      name           = "network_0_24"
      admin_state_up = "false"
    }

    resource "openstack_networking_network_v2" "network_2" {
		name           = "network_1_24"
		admin_state_up = "false"
	}
	
    resource "openstack_networking_network_v2" "network_3" {
		name           = "network_0_32"
		admin_state_up = "false"
	}
    `
}

func testAccNetworkingV2RBACPolicy_shared_update1(projectOneName string) string {
	return fmt.Sprintf(`
    resource "openstack_networking_network_v2" "network_1" {
      name           = "network_0_24"
      admin_state_up = "false"
    }

    resource "openstack_networking_network_v2" "network_2" {
		name           = "network_1_24"
		admin_state_up = "false"
	}
	
    resource "openstack_networking_network_v2" "network_3" {
		name           = "network_0_32"
		admin_state_up = "false"
	}

    resource "openstack_identity_project_v3" "project_1" {
		name        = "%s"
		description = "The second project"
	}
	
	resource "openstack_networking_rbac_policy_v2" "rbac_policy_1" {
		action        = "access_as_shared"
		object_search {
			name_glob = "network_*_24"
		}
		object_type   = "network"
		target_tenant = "${openstack_identity_project_v3.project_1.id}"
	}
    `, projectOneName)
}

func testAccNetworkingV2RBACPolicy_shared_update2(projectOneName, projectTwoName string) string {
	return fmt.Sprintf(`
    resource "openstack_networking_network_v2" "network_1" {
      name           = "network_0_24"
      admin_state_up = "false"
    }

    resource "openstack_networking_network_v2" "network_2" {
		name           = "network_1_24"
		admin_state_up = "false"
	}
	
    resource "openstack_networking_network_v2" "network_3" {
		name           = "network_0_32"
		admin_state_up = "false"
	}

    resource "openstack_identity_project_v3" "project_1" {
		name        = "%s"
		description = "The second project"
	}
	
	resource "openstack_networking_rbac_policy_v2" "rbac_policy_1" {
		action        = "access_as_shared"
		object_search {
			name_glob = "network_*_24"
		}
		object_type   = "network"
		target_tenant = "${openstack_identity_project_v3.project_1.id}"
	}
	
	resource "openstack_identity_project_v3" "project_2" {
		name        = "%s"
		description = "The second project"
	}
	
	resource "openstack_networking_rbac_policy_v2" "rbac_policy_2" {
		action        = "access_as_shared"
		object_search {
			name_glob = "network_*_24"
		}
		object_type   = "network"
		target_tenant = "${openstack_identity_project_v3.project_2.id}"
	}
    `, projectOneName, projectTwoName)
}

func testAccNetworkingV2RBACPolicy_shared_update3(projectThreeName, projectTwoName string) string {
	return fmt.Sprintf(`
    resource "openstack_networking_network_v2" "network_1" {
      name           = "network_0_24"
      admin_state_up = "false"
    }

    resource "openstack_networking_network_v2" "network_2" {
		name           = "network_1_24"
		admin_state_up = "false"
	}
	
    resource "openstack_networking_network_v2" "network_3" {
		name           = "network_0_32"
		admin_state_up = "false"
	}

    resource "openstack_identity_project_v3" "project_3" {
		name        = "%s"
		description = "The second project"
	}
	
	resource "openstack_networking_rbac_policy_v2" "rbac_policy_1" {
		action        = "access_as_shared"
		object_search {
			name_glob = "network_*_24"
		}
		object_type   = "network"
		target_tenant = "${openstack_identity_project_v3.project_3.id}"
	}
	
	resource "openstack_identity_project_v3" "project_2" {
		name        = "%s"
		description = "The second project"
	}
	
	resource "openstack_networking_rbac_policy_v2" "rbac_policy_2" {
		action        = "access_as_shared"
		object_search {
			name_glob = "network_*_24"
		}
		object_type   = "network"
		target_tenant = "${openstack_identity_project_v3.project_2.id}"
	}
    `, projectThreeName, projectTwoName)
}
