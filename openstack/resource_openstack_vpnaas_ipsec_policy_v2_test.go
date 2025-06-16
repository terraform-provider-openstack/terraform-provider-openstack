package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIPSecPolicyVPNaaSV2_basic(t *testing.T) {
	var policy ipsecpolicies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPSecPolicyV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(t.Context(),
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

func TestAccIPSecPolicyVPNaaSV2_withLifetime(t *testing.T) {
	var policy ipsecpolicies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPSecPolicyV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2WithLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(t.Context(),
						"openstack_vpnaas_ipsec_policy_v2.policy_1", &policy),
					testAccCheckLifetime(t, "openstack_vpnaas_ipsec_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
				),
			},
		},
	})
}

func TestAccIPSecPolicyVPNaaSV2_Update(t *testing.T) {
	var policy ipsecpolicies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPSecPolicyV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(t.Context(),
						"openstack_vpnaas_ipsec_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "name", &policy.Name),
				),
			},
			{
				Config: testAccIPSecPolicyV2Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(t.Context(),
						"openstack_vpnaas_ipsec_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "name", &policy.Name),
				),
			},
		},
	})
}

func TestAccIPSecPolicyVPNaaSV2_withLifetimeUpdate(t *testing.T) {
	var policy ipsecpolicies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPSecPolicyV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2WithLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(t.Context(),
						"openstack_vpnaas_ipsec_policy_v2.policy_1", &policy),
					testAccCheckLifetime(t, "openstack_vpnaas_ipsec_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ipsec_policy_v2.policy_1", "pfs", &policy.PFS),
				),
			},
			{
				Config: testAccIPSecPolicyV2WithLifetimeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(t.Context(),
						"openstack_vpnaas_ipsec_policy_v2.policy_1", &policy),
					testAccCheckLifetime(t, "openstack_vpnaas_ipsec_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
				),
			},
		},
	})
}

func testAccCheckIPSecPolicyV2Destroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_vpnaas_ipsec_policy_v2" {
				continue
			}

			_, err = ipsecpolicies.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("IPSec policy (%s) still exists", rs.Primary.ID)
			}

			if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return err
			}
		}

		return nil
	}
}

func testAccCheckIPSecPolicyV2Exists(ctx context.Context, n string, policy *ipsecpolicies.Policy) resource.TestCheckFunc {
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

		found, err := ipsecpolicies.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		*policy = *found

		return nil
	}
}

func testAccCheckLifetime(t *testing.T, n string, unit *string, value *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		// [DEBUG] key: lifetime.452086442.units value: seconds
		// [DEBUG] key: lifetime.452086442.value value: 1200
		// [DEBUG] key: lifetime.# value: 1
		for k, v := range rs.Primary.Attributes {
			t.Logf("[DEBUG] key: %s value: %s", k, v)
			// Do one check for each time a key like "lifetime.<number>.units" is seen.
			// If more than one exists they apparently must all have the same values.
			if strings.HasPrefix(k, "lifetime.") && k[9] >= '0' && k[9] <= '9' && strings.HasSuffix(k, ".units") {
				// Find "lifetime.<number>" so we can append ".value"
				index := strings.LastIndex(k, ".")
				base := k[:index]
				expectedValue := rs.Primary.Attributes[base+".value"]
				expectedUnit := rs.Primary.Attributes[k]
				t.Logf("[DEBUG] expectedValue: %s expectedUnit: %s", expectedValue, expectedUnit)

				if expectedUnit != *unit {
					return fmt.Errorf("Expected lifetime unit %v but found %v", expectedUnit, *unit)
				}

				if expectedValue != strconv.Itoa(*value) {
					return fmt.Errorf("Expected lifetime value %v but found %v", expectedValue, *value)
				}
			}
		}

		return nil
	}
}

const testAccIPSecPolicyV2Basic = `
resource "openstack_vpnaas_ipsec_policy_v2" "policy_1" {
}
`

const testAccIPSecPolicyV2Update = `
resource "openstack_vpnaas_ipsec_policy_v2" "policy_1" {
	name = "updatedname"
}
`

const testAccIPSecPolicyV2WithLifetime = `
resource "openstack_vpnaas_ipsec_policy_v2" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1200
	}
}
`

const testAccIPSecPolicyV2WithLifetimeUpdate = `
resource "openstack_vpnaas_ipsec_policy_v2" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1400
	}
}
`
