package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIKEPolicyVPNaaSV2_basic(t *testing.T) {
	var policy ikepolicies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIKEPolicyV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(t.Context(),
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "name", &policy.Name),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "description", &policy.Description),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "tenant_id", &policy.TenantID),
				),
			},
		},
	})
}

func TestAccIKEPolicyVPNaaSV2_withLifetime(t *testing.T) {
	var policy ikepolicies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIKEPolicyV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2WithLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(t.Context(),
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					// testAccCheckLifetime("openstack_vpnaas_ike_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
				),
			},
		},
	})
}

func TestAccIKEPolicyVPNaaSV2_Update(t *testing.T) {
	var policy ikepolicies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIKEPolicyV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(t.Context(),
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "name", &policy.Name),
				),
			},
			{
				Config: testAccIKEPolicyV2Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(t.Context(),
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "name", &policy.Name),
				),
			},
		},
	})
}

func TestAccIKEPolicyVPNaaSV2_withLifetimeUpdate(t *testing.T) {
	var policy ikepolicies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIKEPolicyV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2WithLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(t.Context(),
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					// testAccCheckLifetime("openstack_vpnaas_ike_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_ike_policy_v2.policy_1", "pfs", &policy.PFS),
				),
			},
			{
				Config: testAccIKEPolicyV2WithLifetimeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(t.Context(),
						"openstack_vpnaas_ike_policy_v2.policy_1", &policy),
					// testAccCheckLifetime("openstack_vpnaas_ike_policy_v2.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
				),
			},
		},
	})
}

func testAccCheckIKEPolicyV2Destroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_vpnaas_ike_policy_v2" {
				continue
			}

			_, err = ikepolicies.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("IKE policy (%s) still exists", rs.Primary.ID)
			}

			if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return err
			}
		}

		return nil
	}
}

func testAccCheckIKEPolicyV2Exists(ctx context.Context, n string, policy *ikepolicies.Policy) resource.TestCheckFunc {
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

		found, err := ikepolicies.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		*policy = *found

		return nil
	}
}

const testAccIKEPolicyV2Basic = `
resource "openstack_vpnaas_ike_policy_v2" "policy_1" {
}
`

const testAccIKEPolicyV2Update = `
resource "openstack_vpnaas_ike_policy_v2" "policy_1" {
	name = "updatedname"
}
`

const testAccIKEPolicyV2WithLifetime = `
resource "openstack_vpnaas_ike_policy_v2" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1200
	}
}
`

const testAccIKEPolicyV2WithLifetimeUpdate = `
resource "openstack_vpnaas_ike_policy_v2" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1400
	}
}
`
