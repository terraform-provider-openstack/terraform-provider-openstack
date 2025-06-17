package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkingV2SecGroup_basic(t *testing.T) {
	var securityGroup groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(t.Context(), "openstack_networking_secgroup_v2.secgroup_1", &securityGroup),
					testAccCheckNetworkingV2SecGroupRuleCount(&securityGroup, 2),
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "stateful", "true"),
				),
			},
			{
				Config: testAccNetworkingV2SecGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("openstack_networking_secgroup_v2.secgroup_1", "id", &securityGroup.ID),
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "name", "security_group_2"),
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "stateful", "false"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_noDefaultRules(t *testing.T) {
	var securityGroup groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupNoDefaultRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(t.Context(),
						"openstack_networking_secgroup_v2.secgroup_1", &securityGroup),
					testAccCheckNetworkingV2SecGroupRuleCount(&securityGroup, 0),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_statefulNotSet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupStatefulNotSet,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "name", "security_group_1"),
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "stateful", "true"),
				),
			},
			{
				Config: testAccNetworkingV2SecGroupStatefulNotSetUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "name", "security_group_1"),
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "stateful", "false"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_statefulSetTrue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupStatefulSetTrue,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "name", "security_group_1"),
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "stateful", "true"),
				),
			},
			{
				Config: testAccNetworkingV2SecGroupStatefulSetTrueUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "name", "security_group_1"),
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "stateful", "false"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_statefulSetFalse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupStatefulSetFalse,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "name", "security_group_1"),
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "stateful", "false"),
				),
			},
			{
				Config: testAccNetworkingV2SecGroupStatefulSetFalseUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "name", "security_group_1"),
					resource.TestCheckResourceAttr("openstack_networking_secgroup_v2.secgroup_1", "stateful", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SecGroup_timeout(t *testing.T) {
	var securityGroup groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(t.Context(),
						"openstack_networking_secgroup_v2.secgroup_1", &securityGroup),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SecGroupDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_networking_secgroup_v2" {
				continue
			}

			_, err := groups.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Security group still exists")
			}
		}

		return nil
	}
}

func testAccCheckNetworkingV2SecGroupExists(ctx context.Context, n string, sg *groups.SecGroup) resource.TestCheckFunc {
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

		found, err := groups.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Security group not found")
		}

		*sg = *found

		return nil
	}
}

func testAccCheckNetworkingV2SecGroupRuleCount(sg *groups.SecGroup, count int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if len(sg.Rules) == count {
			return nil
		}

		return fmt.Errorf("Unexpected number of rules in group %s. Expected %d, got %d",
			sg.ID, count, len(sg.Rules))
	}
}

const testAccNetworkingV2SecGroupBasic = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}
`

const testAccNetworkingV2SecGroupUpdate = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group_2"
  description = "terraform security group acceptance test"
  stateful    = false
}
`

const testAccNetworkingV2SecGroupNoDefaultRules = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
	name = "security_group_1"
	description = "terraform security group acceptance test"
	delete_default_rules = true
}
`

const testAccNetworkingV2SecGroupTimeout = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"

  timeouts {
    delete = "5m"
  }
}
`

const testAccNetworkingV2SecGroupStatefulNotSet = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group_1"
  description = "stateful flag not set, expect true"
}
`

const testAccNetworkingV2SecGroupStatefulNotSetUpdate = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group_1"
  description = "stateful flag updated to false"
  stateful = false
}
`

const testAccNetworkingV2SecGroupStatefulSetTrue = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group_1"
  description = "stateful flag set to true"
  stateful = true
}
`

const testAccNetworkingV2SecGroupStatefulSetTrueUpdate = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group_1"
  description = "stateful flag updated to false"
  stateful = false
}
`

const testAccNetworkingV2SecGroupStatefulSetFalse = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group_1"
  description = "stateful flag explicitly set to false"
  stateful = false
}
`

const testAccNetworkingV2SecGroupStatefulSetFalseUpdate = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group_1"
  description = "stateful flag updated to true"
  stateful = true
}
`
