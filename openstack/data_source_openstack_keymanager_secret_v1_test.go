package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKeyManagerSecretV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_keymanager_secret_v1.secret_1", "id",
						"openstack_keymanager_secret_v1.secret_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_keymanager_secret_v1.secret_2", "id",
						"openstack_keymanager_secret_v1.secret_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_keymanager_secret_v1.secret_1", "payload",
						"openstack_keymanager_secret_v1.secret_1", "payload"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_keymanager_secret_v1.secret_2", "payload",
						"openstack_keymanager_secret_v1.secret_2", "payload"),
					resource.TestCheckResourceAttr(
						"data.openstack_keymanager_secret_v1.secret_1", "metadata.foo", "update"),
				),
			},
		},
	})
}

const testAccKeyManagerSecretV1DataSource_basic = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm   = "aes"
  bit_length  = 192
  mode        = "cbc"
  name        = "mysecret"
  payload     = "foobar"
  secret_type = "passphrase"
  payload_content_type = "text/plain"
  metadata = {
    foo = "update"
  }
}

resource "openstack_keymanager_secret_v1" "secret_2" {
  algorithm   = "aes"
  bit_length  = 256
  mode        = "cbc"
  name        = "mysecret"
  secret_type = "passphrase"
  payload     = "foo"
  expiration  = "3000-07-31T12:02:46Z"
  payload_content_type = "text/plain"
}

data "openstack_keymanager_secret_v1" "secret_1" {
  bit_length  = "${openstack_keymanager_secret_v1.secret_1.bit_length}"
  secret_type = "passphrase"
}

data "openstack_keymanager_secret_v1" "secret_2" {
  mode              = "cbc"
  secret_type       = "passphrase"
  expiration_filter = "${openstack_keymanager_secret_v1.secret_2.expiration}"
}
`
