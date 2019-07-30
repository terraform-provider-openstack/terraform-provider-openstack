package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccKeyManagerContainerV1_basic(t *testing.T) {
	var container containers.Container
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1Exists(
						"openstack_keymanager_container_v1.container_1", &container),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "secret_refs.#", "3"),
				),
			},
			{
				Config: testAccKeyManagerContainerV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1Exists(
						"openstack_keymanager_container_v1.container_1", &container),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "secret_refs.#", "2"),
				),
			},
			{
				Config: testAccKeyManagerContainerV1_update1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1Exists(
						"openstack_keymanager_container_v1.container_1", &container),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "secret_refs.#", "3"),
				),
			},
			{
				Config: testAccKeyManagerContainerV1_update2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1Exists(
						"openstack_keymanager_container_v1.container_1", &container),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "secret_refs.#", "0"),
				),
			},
		},
	})
}

func testAccCheckContainerV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	kmClient, err := config.keyManagerV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack KeyManager client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_keymanager_container" {
			continue
		}
		_, err = containers.Get(kmClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Container (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckContainerV1Exists(n string, container *containers.Container) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		kmClient, err := config.keyManagerV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack KeyManager client: %s", err)
		}

		var found *containers.Container

		found, err = containers.Get(kmClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*container = *found

		return nil
	}
}

const testAccKeyManagerContainerV1 = `
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
`

var testAccKeyManagerContainerV1_basic = fmt.Sprintf(`
%s

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
`, testAccKeyManagerContainerV1)

var testAccKeyManagerContainerV1_update = fmt.Sprintf(`
%s

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
}
`, testAccKeyManagerContainerV1)

var testAccKeyManagerContainerV1_update1 = fmt.Sprintf(`
%s

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
    name       = "intermediate_new"
    secret_ref = "${openstack_keymanager_secret_v1.intermediate_1.secret_ref}"
  }
}
`, testAccKeyManagerContainerV1)

var testAccKeyManagerContainerV1_update2 = fmt.Sprintf(`
%s

resource "openstack_keymanager_container_v1" "container_1" {
  name = "generic"
  type = "generic"
}
`, testAccKeyManagerContainerV1)
