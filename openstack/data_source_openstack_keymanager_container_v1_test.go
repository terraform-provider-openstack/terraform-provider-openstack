package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKeyManagerContainerV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckKeyManager(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerV1DataSourceBasic,
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

func TestAccKeyManagerContainerV1DataSource_acls(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckKeyManager(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerV1DataSourceAcls,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_keymanager_container_v1.container_1", "id",
						"openstack_keymanager_container_v1.container_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_keymanager_container_v1.container_1", "secret_refs",
						"openstack_keymanager_container_v1.container_1", "secret_refs"),
					resource.TestCheckResourceAttr(
						"data.openstack_keymanager_container_v1.container_1", "secret_refs.#", "3"),
					resource.TestCheckResourceAttr("data.openstack_keymanager_container_v1.container_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("data.openstack_keymanager_container_v1.container_1", "acl.0.read.0.users.#", "2"),
				),
			},
		},
	})
}

const testAccKeyManagerContainerV1DataSourceBasic = `
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

const testAccKeyManagerContainerV1DataSourceAcls = `
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

data "openstack_keymanager_container_v1" "container_1" {
  name = "${openstack_keymanager_container_v1.container_1.name}"
}
`
