package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/orders"
	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/secrets"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccKeyManagerOrderV1_basic(t *testing.T) {
	var order orders.Order

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckOrderV1Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerOrderV1Symmetric,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrderV1Exists(t.Context(),
						"openstack_keymanager_order_v1.test-acc-basic", &order),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_order_v1.test-acc-basic", "type", &order.Type),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_order_v1.test-acc-basic", "meta.0.name", &order.Meta.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_order_v1.test-acc-basic", "meta.0.algorithm", &order.Meta.Algorithm),
					resource.TestCheckResourceAttrSet("openstack_keymanager_order_v1.test-acc-basic", "meta.0.bit_length"),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_order_v1.test-acc-basic", "meta.0.mode", &order.Meta.Mode),
				),
			},
		},
	})
}

func testAccCheckOrderV1Destroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		kmClient, err := config.KeyManagerV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack KeyManager client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_keymanager_order_v1" {
				continue
			}

			_, err = orders.Get(ctx, kmClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("Order (%s) still exists", rs.Primary.ID)
			}

			if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return err
			}

			secretRefSplit := strings.Split(rs.Primary.Attributes["secret_ref"], "/")
			uuid := secretRefSplit[len(secretRefSplit)-1]
			result := secrets.Delete(ctx, kmClient, uuid)

			if result.ExtractErr() != nil {
				return fmt.Errorf("Secret (%s) still exists", uuid)
			}
		}

		return nil
	}
}

func testAccCheckOrderV1Exists(ctx context.Context, n string, order *orders.Order) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		kmClient, err := config.KeyManagerV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack KeyManager client: %w", err)
		}

		var found *orders.Order

		found, err = orders.Get(ctx, kmClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		*order = *found

		return nil
	}
}

const testAccKeyManagerOrderV1Symmetric = `
resource "openstack_keymanager_order_v1" "test-acc-basic" {
  type = "key"
  meta {
    name = "test-acc-basic"
    algorithm = "aes"
    bit_length = 256
    mode = "cbc"
  }
}`
