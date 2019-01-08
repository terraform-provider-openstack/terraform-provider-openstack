package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	"github.com/hashicorp/terraform/flatmap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
)

func TestAccIPSecPolicyV2_basic(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"openstack_vpnaas_ipsec_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "name", &policy.Name),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "description", &policy.Description),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "tenant_id", &policy.TenantID),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "pfs", &policy.PFS),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "transform_protocol", &policy.TransformProtocol),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "encapsulation_mode", &policy.EncapsulationMode),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "encryption_algorithm", &policy.EncryptionAlgorithm),
				),
			},
		},
	})
}

func TestAccIPSecPolicyV2_withLifetime(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2_withLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"openstack_vpnaas_ipsec_policy_v2.policy_1", &policy),
					testAccCheckLifetime("openstack_vpnaas_ipsec_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
				),
			},
		},
	})
}

func TestAccIPSecPolicyV2_Update(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"openstack_vpnaas_ipsec_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "name", &policy.Name),
				),
			},
			{
				Config: testAccIPSecPolicyV2_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"openstack_vpnaas_ipsec_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "name", &policy.Name),
				),
			},
		},
	})
}

func TestAccIPSecPolicyV2_withLifetimeUpdate(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2_withLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"openstack_vpnaas_ipsec_policy_v2.policy_1", &policy),
					testAccCheckLifetime("openstack_vpnaas_ipsec_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "pfs", &policy.PFS),
				),
			},
			{
				Config: testAccIPSecPolicyV2_withLifetimeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"openstack_vpnaas_ipsec_policy_v2.policy_1", &policy),
					testAccCheckLifetime("openstack_vpnaas_ipsec_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
				),
			},
		},
	})
}

func testAccCheckIPSecPolicyV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_vpnaas_ipsec_policy_v2" {
			continue
		}
		_, err = ipsecpolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IPSec policy (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckIPSecPolicyV2Exists(n string, policy *ipsecpolicies.Policy) resource.TestCheckFunc {
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

		found, err := ipsecpolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*policy = *found

		return nil
	}
}

func testAccCheckLifetime(n string, unit *string, value *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		lifetime := flatmap.Expand(rs.Primary.Attributes, "lifetime")
		for _, raw := range lifetime.([]interface{}) {
			rawMap := raw.(map[string]interface{})

			expectedValue := rawMap["value"]
			expectedUnit := rawMap["units"]
			if expectedUnit != *unit {
				return fmt.Errorf("Expected lifetime unit %v but found %v", expectedUnit, *unit)
			}
			if expectedValue != strconv.Itoa(*value) {
				return fmt.Errorf("Expected lifetime value %v but found %v", expectedValue, *value)
			}
		}
		return nil
	}
}

const testAccIPSecPolicyV2_basic = `
resource "openstack_vpnaas_ipsec_policy_v2" "policy_1" {
}
`

const testAccIPSecPolicyV2_Update = `
resource "openstack_vpnaas_ipsec_policy_v2" "policy_1" {
	name = "updatedname"
}
`

const testAccIPSecPolicyV2_withLifetime = `
resource "openstack_vpnaas_ipsec_policy_v2" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1200
	}
}
`

const testAccIPSecPolicyV2_withLifetimeUpdate = `
resource "openstack_vpnaas_ipsec_policy_v2" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1400
	}
}
`
