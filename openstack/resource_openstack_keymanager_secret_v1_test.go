package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
)

func TestAccKeyManagerSecretV1_basic(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_1", "payload", "foobar"),
				),
			},
		},
	})
}

func TestAccKeyManagerSecretV1_basicWithMetadata(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretV1BasicWithMetadata,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
				),
			},
		},
	})
}

func TestAccKeyManagerSecretV1_updateMetadata(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretV1BasicWithMetadata,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckMetadataEquals("foo", "bar", &secret),
				),
			},
			{
				Config: testAccKeyManagerSecretV1UpdateMetadata,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckMetadataEquals("foo", "update", &secret),
				),
			},
		},
	})
}

func TestAccUpdateSecretV1_payload(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretV1NoPayload,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckPayloadEquals("", &secret),
				),
			},
			{
				Config: testAccKeyManagerSecretV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckPayloadEquals("updatedfoobar", &secret),
				),
			},
			{
				Config: testAccKeyManagerSecretV1UpdateWhitespace,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckPayloadEquals("updatedfoobar", &secret),
				),
			},
			{
				Config: testAccKeyManagerSecretV1UpdateBase64,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckPayloadEquals("base64foobar ", &secret),
				),
			},
		},
	})
}

func TestAccKeyManagerSecretV1_acls(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretV1Acls,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_1", "acl.0.read.0.users.#", "2"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_2", "acl.0.read.0.project_access", "true"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_2", "acl.0.read.0.users.#", "0"),
				),
			},
		},
	})
}

func TestAccKeyManagerSecretV1_acls_update(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretV1Acls,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_1", "acl.0.read.0.users.#", "2"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_2", "acl.0.read.0.project_access", "true"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_2", "acl.0.read.0.users.#", "0"),
				),
			},
			{
				Config: testAccKeyManagerSecretV1AclsUpdate1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_1", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_1", "acl.0.read.0.users.#", "2"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_2", "acl.0.read.0.project_access", "false"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_2", "acl.0.read.0.users.#", "1"),
				),
			},
			{
				Config: testAccKeyManagerSecretV1AclsUpdate2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_1", "acl.0.read.0.project_access", "true"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_1", "acl.0.read.0.users.#", "0"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_2", "acl.0.read.0.project_access", "true"),
					resource.TestCheckResourceAttr("openstack_keymanager_secret_v1.secret_2", "acl.0.read.0.users.#", "0"),
				),
			},
		},
	})
}

func testAccCheckSecretV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	kmClient, err := config.KeyManagerV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack KeyManager client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_keymanager_secret" {
			continue
		}
		_, err = secrets.Get(kmClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Secret (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckSecretV1Exists(n string, secret *secrets.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		kmClient, err := config.KeyManagerV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack KeyManager client: %s", err)
		}

		var found *secrets.Secret

		found, err = secrets.Get(kmClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*secret = *found

		return nil
	}
}

func testAccCheckPayloadEquals(payload string, secret *secrets.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		kmClient, err := config.KeyManagerV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack KeyManager client: %s", err)
		}

		opts := secrets.GetPayloadOpts{
			PayloadContentType: "text/plain",
		}

		uuid := keyManagerSecretV1GetUUIDfromSecretRef(secret.SecretRef)
		secretPayload, _ := secrets.GetPayload(kmClient, uuid, opts).Extract()
		if string(secretPayload) != payload {
			return fmt.Errorf("Payloads do not match. Expected %v but got %v", payload, secretPayload)
		}
		return nil
	}
}

func testAccCheckMetadataEquals(key string, value string, secret *secrets.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		kmClient, err := config.KeyManagerV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		uuid := keyManagerSecretV1GetUUIDfromSecretRef(secret.SecretRef)
		metadatum, err := secrets.GetMetadatum(kmClient, uuid, key).Extract()
		if err != nil {
			return err
		}
		if metadatum.Value != value {
			return fmt.Errorf("Metadata does not match. Expected %v but got %v", metadatum, value)
		}

		return nil
	}
}

const testAccKeyManagerSecretV1Basic = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
}`

const testAccKeyManagerSecretV1BasicWithMetadata = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
  metadata = {
    foo = "bar"
  }
}`

const testAccKeyManagerSecretV1UpdateMetadata = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
  metadata = {
    foo = "update"
  }
}`

const testAccKeyManagerSecretV1NoPayload = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  secret_type = "passphrase"
  payload = ""
}`

const testAccKeyManagerSecretV1Update = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "updatedfoobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
}`

const testAccKeyManagerSecretV1UpdateWhitespace = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = <<EOF
updatedfoobar
EOF
  payload_content_type = "text/plain"
  secret_type = "passphrase"
}`

const testAccKeyManagerSecretV1UpdateBase64 = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "${base64encode("base64foobar ")}"
  payload_content_type = "application/octet-stream"
  payload_content_encoding = "base64"
  secret_type = "passphrase"
}`

const testAccKeyManagerSecretV1Acls = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "${base64encode("base64foobar ")}"
  payload_content_type = "application/octet-stream"
  payload_content_encoding = "base64"
  secret_type = "passphrase"

  acl {
    read {
      project_access = false
      users = [
        "619e2ad074321cf246b03a89e95afee95fb26bb0b2d1fc7ba3bd30fcca25588a",
        "96b3ebddf275996285eae440e71227ba47c651be18391b0f2ebf1032ebae5dca",
      ]
    }
  }
}

resource "openstack_keymanager_secret_v1" "secret_2" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
}
`

const testAccKeyManagerSecretV1AclsUpdate1 = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "${base64encode("base64foobar ")}"
  payload_content_type = "application/octet-stream"
  payload_content_encoding = "base64"
  secret_type = "passphrase"

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

resource "openstack_keymanager_secret_v1" "secret_2" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"

  acl {
    read {
      project_access = false
      users = [
        "96b3ebddf275996285eae440e71227ba47c651be18391b0f2ebf1032ebae5dca",
      ]
    }
  }
}
`

const testAccKeyManagerSecretV1AclsUpdate2 = `
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "${base64encode("base64foobar ")}"
  payload_content_type = "application/octet-stream"
  payload_content_encoding = "base64"
  secret_type = "passphrase"

  acl {
    read {
      project_access = true
    }
  }
}

resource "openstack_keymanager_secret_v1" "secret_2" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"

  acl {
    read {
      project_access = true
    }
  }
}
`
