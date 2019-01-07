package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIKEPolicyV2_basic(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIKEPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "name", &policy.Name),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "description", &policy.Description),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "tenant_id", &policy.TenantID),
				),
			},
		},
	})
}

func TestAccIKEPolicyV2_withLifetime(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIKEPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2_withLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					testAccCheckLifetime("openstack_vpnaas_ike_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
				),
			},
		},
	})
}

func TestAccIKEPolicyV2_Update(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIKEPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "name", &policy.Name),
				),
			},
			{
				Config: testAccIKEPolicyV2_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "name", &policy.Name),
				),
			},
		},
	})
}

func TestAccIKEPolicyV2_withLifetimeUpdate(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIKEPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2_withLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					testAccCheckLifetime("openstack_vpnaas_ike_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "pfs", &policy.PFS),
				),
			},
			{
				Config: testAccIKEPolicyV2_withLifetimeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					testAccCheckLifetime("openstack_vpnaas_ike_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
				),
			},
		},
	})
}

func testAccCheckIKEPolicyV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_vpnaas_ike_policy_v2" {
			continue
		}
		_, err = ikepolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IKE policy (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckIKEPolicyV2Exists(n string, policy *ikepolicies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := ikepolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*policy = *found

		return nil
	}
}

const testAccIKEPolicyV2_basic = `
resource "openstack_vpnaas_ike_policy_v2" "policy_1" {
}
`

const testAccIKEPolicyV2_Update = `
resource "openstack_vpnaas_ike_policy_v2" "policy_1" {
	name = "updatedname"
}
`

const testAccIKEPolicyV2_withLifetime = `
resource "openstack_vpnaas_ike_policy_v2" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1200
	}
}
`

const testAccIKEPolicyV2_withLifetimeUpdate = `
resource "openstack_vpnaas_ike_policy_v2" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1400
	}
}
`
