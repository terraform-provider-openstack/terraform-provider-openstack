package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
)

func TestAccComputeV2KeypairDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2KeypairDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2KeypairDataSourceID("data.openstack_compute_keypair_v2.kp"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_keypair_v2.kp", "name", "the-key-name"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_keypair_v2.kp", "fingerprint", "78:a9:d0:f9:af:a8:1b:ca:bb:9f:65:88:47:af:1d:a9"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_keypair_v2.kp", "public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"),
				),
			},
		},
	})
}

func TestAccComputeV2KeypairDataSourceOtherUser(t *testing.T) {
	var user users.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2KeypairDataSourceOtherUser,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2KeypairDataSourceID("data.openstack_compute_keypair_v2.kp"),
					testAccCheckIdentityV3UserExists("openstack_identity_user_v3.user_1", &user),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_keypair_v2.kp", "name", "the-key-name"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_keypair_v2.kp", "fingerprint", "78:a9:d0:f9:af:a8:1b:ca:bb:9f:65:88:47:af:1d:a9"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_keypair_v2.kp", "public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"),
					resource.TestCheckResourceAttrPtr(
						"data.openstack_compute_keypair_v2.kp", "user_id", &user.ID),
				),
			},
		},
	})
}

func testAccCheckComputeV2KeypairDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find keypair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Keypair data source ID not set")
		}

		return nil
	}
}

const testAccComputeV2KeypairDataSourceBasic = `
resource "openstack_compute_keypair_v2" "kp" {
  name = "the-key-name"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

data "openstack_compute_keypair_v2" "kp" {
  name = "${openstack_compute_keypair_v2.kp.name}"
}
`

const testAccComputeV2KeypairDataSourceOtherUser = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}
	
resource "openstack_identity_user_v3" "user_1" {
  name = "user_1"
  default_project_id = "${openstack_identity_project_v3.project_1.id}"
}
  
resource "openstack_compute_keypair_v2" "kp" {
  name = "the-key-name"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
  user_id = "${openstack_identity_user_v3.user_1.id}"
}
data "openstack_compute_keypair_v2" "kp" {
  name = "${openstack_compute_keypair_v2.kp.name}"
  user_id = "${openstack_identity_user_v3.user_1.id}"
}
`
