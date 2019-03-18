package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccKeymanagerSecretV1_basic(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeymanagerSecretV1_basic,
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

func TestAccKeymanagerSecretV1_basicWithMetadata(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeymanagerSecretV1_basicWithMetadata,
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

func TestAccKeymanagerSecretV1_updateMetadata(t *testing.T) {
	var secret secrets.Secret
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeymanagerSecretV1_basicWithMetadata,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckMetadataEquals("foo", "bar", &secret),
				),
			},
			{
				Config: testAccKeymanagerSecretV1_updateMetadata,
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
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeymanagerSecretV1_noPayload,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckPayloadEquals("", &secret),
				),
			},
			{
				Config: testAccKeymanagerSecretV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckPayloadEquals("updatedfoobar", &secret),
				),
			},
		},
	})
}

func testAccCheckSecretV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	kmClient, err := config.keymanagerV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Keymanager client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_keymanager_secret" {
			continue
		}
		_, err = secrets.Get(kmClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Secret (%s) still exists.", rs.Primary.ID)
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
		kmClient, err := config.keymanagerV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack keymanager client: %s", err)
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
		kmClient, err := config.keymanagerV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack keymanager client: %s", err)
		}

		opts := secrets.GetPayloadOpts{
			PayloadContentType: "text/plain",
		}

		uuid := keymanagerSecretV1GetUUIDfromSecretRef(secret.SecretRef)
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
		kmClient, err := config.keymanagerV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		uuid := keymanagerSecretV1GetUUIDfromSecretRef(secret.SecretRef)
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

var testAccKeymanagerSecretV1_basic = fmt.Sprintf(`
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "foobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
}`)

var testAccKeymanagerSecretV1_basicWithMetadata = fmt.Sprintf(`
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
}`)

var testAccKeymanagerSecretV1_updateMetadata = fmt.Sprintf(`
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
}`)

var testAccKeymanagerSecretV1_noPayload = fmt.Sprintf(`
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  secret_type = "passphrase"
  payload = ""
}`)

var testAccKeymanagerSecretV1_update = fmt.Sprintf(`
resource "openstack_keymanager_secret_v1" "secret_1" {
  algorithm = "aes"
  bit_length = 256
  mode = "cbc"
  name = "mysecret"
  payload = "updatedfoobar"
  payload_content_type = "text/plain"
  secret_type = "passphrase"
}`)
