package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/taas/tapmirrors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccTaasTapMirrorV2_basic(t *testing.T) {
	var tapMirror tapmirrors.TapMirror

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckTaas(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckTapMirrorV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccTapMirrorV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTapMirrorV2Exists(t.Context(),
						"openstack_taas_tap_mirror_v2.tap_mirror_1", &tapMirror),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "name", "tap_mirror_1"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "description", "desc"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "mirror_type", "erspanv1"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "directions.0.in", "1000"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "directions.0.out", "1001"),
				),
			},
		},
	})
}

func TestAccTaasTapMirrorV2_update(t *testing.T) {
	var tapMirror tapmirrors.TapMirror

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckTaas(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckTapMirrorV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccTapMirrorV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTapMirrorV2Exists(t.Context(),
						"openstack_taas_tap_mirror_v2.tap_mirror_1", &tapMirror),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "name", "tap_mirror_1"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "description", "desc"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "mirror_type", "erspanv1"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "directions.0.in", "1000"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "directions.0.out", "1001"),
				),
			},
			{
				Config: testAccTapMirrorV2Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTapMirrorV2Exists(t.Context(),
						"openstack_taas_tap_mirror_v2.tap_mirror_1", &tapMirror),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "name", "updated tap_mirror_1"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "description", "updated desc"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "mirror_type", "erspanv1"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "directions.0.in", "1000"),
					resource.TestCheckResourceAttr("openstack_taas_tap_mirror_v2.tap_mirror_1", "directions.0.out", "1001"),
				),
			},
		},
	})
}

func testAccCheckTapMirrorV2Destroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_taas_tap_mirror_v2" {
				continue
			}

			_, err = tapmirrors.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("TapMirror (%s) still exists", rs.Primary.ID)
			}

			if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return err
			}
		}

		return nil
	}
}

func testAccCheckTapMirrorV2Exists(ctx context.Context, n string, tapMirror *tapmirrors.TapMirror) resource.TestCheckFunc {
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

		var found *tapmirrors.TapMirror

		found, err = tapmirrors.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		*tapMirror = *found

		return nil
	}
}

const testAccTapMirrorV2Basic = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = openstack_networking_network_v2.network_1.id

  fixed_ip {
	subnet_id = openstack_networking_subnet_v2.subnet_1.id
  }
}

resource "openstack_taas_tap_mirror_v2" "tap_mirror_1" {
    name = "tap_mirror_1"
    description = "desc"
	mirror_type = "erspanv1"
	port_id = openstack_networking_port_v2.port_1.id
	remote_ip = openstack_networking_port_v2.port_1.all_fixed_ips[0]
	directions {
		in = 1000
		out = 1001
	}
}
`

const testAccTapMirrorV2Update = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = openstack_networking_network_v2.network_1.id

  fixed_ip {
	subnet_id = openstack_networking_subnet_v2.subnet_1.id
  }
}

resource "openstack_taas_tap_mirror_v2" "tap_mirror_1" {
    name = "updated tap_mirror_1"
    description = "updated desc"
	mirror_type = "erspanv1"
	port_id = openstack_networking_port_v2.port_1.id
	remote_ip = openstack_networking_port_v2.port_1.all_fixed_ips[0]
	directions {
		in = 1000
		out = 1001
	}
}
`
