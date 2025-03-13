package openstack

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/containers"
)

func TestAccKeyManagerContainerV1_basic(t *testing.T) {
	var container containers.Container
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerV1Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1Exists(
						"openstack_keymanager_container_v1.container_1", &container),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "secret_refs.#", "3"),
				),
			},
			{
				Config: testAccKeyManagerContainerV1Update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1Exists(
						"openstack_keymanager_container_v1.container_1", &container),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "secret_refs.#", "2"),
				),
			},
			{
				Config: testAccKeyManagerContainerV1Update1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1Exists(
						"openstack_keymanager_container_v1.container_1", &container),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "secret_refs.#", "3"),
				),
			},
			{
				Config: testAccKeyManagerContainerV1Update2(),
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

func TestAccKeyManagerContainerV1_acls(t *testing.T) {
	var container containers.Container
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerV1Acls(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1Exists(
						"openstack_keymanager_container_v1.container_1", &container),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "secret_refs.#", "3"),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "acl.0.read.0.users.#", "2"),
				),
			},
		},
	})
}

func TestAccKeyManagerContainerV1_certificate_type(t *testing.T) {
	var container containers.Container
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerV1CertificateType(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1Exists(
						"openstack_keymanager_container_v1.container_1", &container),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "secret_refs.#", "3"),
				),
			},
		},
	})
}

func TestAccKeyManagerContainerV1_acls_update(t *testing.T) {
	var container containers.Container
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerV1Acls(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1Exists(
						"openstack_keymanager_container_v1.container_1", &container),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "secret_refs.#", "3"),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "acl.0.read.0.users.#", "2"),
				),
			},
			{
				Config: testAccKeyManagerContainerV1AclsUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1Exists(
						"openstack_keymanager_container_v1.container_1", &container),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "name", &container.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_container_v1.container_1", "type", &container.Type),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "secret_refs.#", "3"),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "acl.0.read.0.project_access", "true"),
					resource.TestCheckResourceAttr("openstack_keymanager_container_v1.container_1", "acl.0.read.0.users.#", "0"),
				),
			},
		},
	})
}

func testAccCheckContainerV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	kmClient, err := config.KeyManagerV1Client(context.TODO(), osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack KeyManager client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_keymanager_container" {
			continue
		}
		_, err = containers.Get(context.TODO(), kmClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Container (%s) still exists", rs.Primary.ID)
		}
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
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
		kmClient, err := config.KeyManagerV1Client(context.TODO(), osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack KeyManager client: %s", err)
		}

		var found *containers.Container

		found, err = containers.Get(context.TODO(), kmClient, rs.Primary.ID).Extract()
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

func testAccKeyManagerContainerV1Basic() string {
	return fmt.Sprintf(`
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
    name       = "intermediates"
    secret_ref = "${openstack_keymanager_secret_v1.intermediate_1.secret_ref}"
  }
}
`, testAccKeyManagerContainerV1)
}

func testAccKeyManagerContainerV1Update() string {
	return fmt.Sprintf(`
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
}

func testAccKeyManagerContainerV1Update1() string {
	return fmt.Sprintf(`
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
}

func testAccKeyManagerContainerV1Update2() string {
	return fmt.Sprintf(`
%s

resource "openstack_keymanager_container_v1" "container_1" {
  name = "generic"
  type = "certificate"
}
`, testAccKeyManagerContainerV1)
}

func testAccKeyManagerContainerV1Acls() string {
	return fmt.Sprintf(`
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
    name       = "intermediates"
    secret_ref = "${openstack_keymanager_secret_v1.intermediate_1.secret_ref}"
  }

  acl {
    read {
      project_access = false
      users = [
        "96b3ebddf275996285eae440e71227ba47c651be18391b0f2ebf1032ebae5dca",
        "619e2ad074321cf246b03a89e95afee95fb26bb0b2d1fc7ba3bd30fcca25588a",
      ]
    }
  }
}
`, testAccKeyManagerContainerV1)
}

func testAccKeyManagerContainerV1AclsUpdate() string {
	return fmt.Sprintf(`
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
    name       = "intermediates"
    secret_ref = "${openstack_keymanager_secret_v1.intermediate_1.secret_ref}"
  }

  acl {
    read {}
  }
}
`, testAccKeyManagerContainerV1)
}

func testAccKeyManagerContainerV1CertificateType() string {
	return fmt.Sprintf(`
%s

resource "openstack_keymanager_container_v1" "container_1" {
  name = "generic"
  type = "certificate"

  secret_refs {
    name       = "certificate"
    secret_ref = "${openstack_keymanager_secret_v1.certificate_1.secret_ref}"
  }

  secret_refs {
    name       = "private_key"
    secret_ref = "${openstack_keymanager_secret_v1.private_key_1.secret_ref}"
  }

  secret_refs {
    name       = "intermediates"
    secret_ref = "${openstack_keymanager_secret_v1.intermediate_1.secret_ref}"
  }
}
`, testAccKeyManagerContainerV1)
}
