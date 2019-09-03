package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccKeyManagerContainerV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_keymanager_container_v1.container_1", "id",
						"openstack_keymanager_container_v1.container_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_keymanager_container_v1.container_1", "secret_refs",
						"openstack_keymanager_container_v1.container_1", "secret_refs"),
					resource.TestCheckResourceAttr(
						"data.openstack_keymanager_container_v1.container_1", "secret_refs.#", "3"),
				),
			},
		},
	})
}

const testAccKeyManagerContainerV1DataSource_basic = `
resource "openstack_keymanager_secret_v1" "certificate_1" {
  name                 = "certificate"
  payload              = "certificate"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "openstack_keymanager_secret_v1" "private_key_1" {
  name                 = "private_key"
  payload              = "private_key"
  secret_type          = "private"
  payload_content_type = "text/plain"
}

resource "openstack_keymanager_secret_v1" "intermediate_1" {
  name                 = "intermediate"
  payload              = "intermediate"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "openstack_keymanager_container_v1" "container_1" {
  name = "generic"
  type = "generic"

  secret_refs {
    name       = "certificate"
    secret_ref = "${openstack_keymanager_secret_v1.certificate_1.secret_ref}"
  }

  secret_refs {
    name       = "private_key"
    secret_ref = "${openstack_keymanager_secret_v1.private_key_1.secret_ref}"
  }

  secret_refs {
    name       = "intermediate"
    secret_ref = "${openstack_keymanager_secret_v1.intermediate_1.secret_ref}"
  }
}

data "openstack_keymanager_container_v1" "container_1" {
  name = "${openstack_keymanager_container_v1.container_1.name}"
}
`
