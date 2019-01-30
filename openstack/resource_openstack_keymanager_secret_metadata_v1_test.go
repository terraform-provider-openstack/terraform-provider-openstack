package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSecretMetadataV1_basic(t *testing.T) {
	var secret secrets.Secret
	var metadata map[string]string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretMetadataV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretMetadataV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					testAccCheckSecretMetadataV1Exists(
						"openstack_keymanager_secret_metadata_v1.metadata_1", &metadata),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckMetadataEquals("foo", "bar", &secret),
				),
			},
		},
	})
}

func TestAccSecretMetadataV1_Update(t *testing.T) {
	var secret secrets.Secret
	var metadata map[string]string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretMetadataV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecretMetadataV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					testAccCheckSecretMetadataV1Exists(
						"openstack_keymanager_secret_metadata_v1.metadata_1", &metadata),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckMetadataEquals("foo", "bar", &secret),
				),
			},
			{
				Config: testAccSecretMetadataV1_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecretV1Exists(
						"openstack_keymanager_secret_v1.secret_1", &secret),
					testAccCheckSecretMetadataV1Exists(
						"openstack_keymanager_secret_metadata_v1.metadata_1", &metadata),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "name", &secret.Name),
					resource.TestCheckResourceAttrPtr("openstack_keymanager_secret_v1.secret_1", "secret_type", &secret.SecretType),
					testAccCheckMetadataEquals("foo", "update", &secret),
					testAccCheckMetadataEquals("update", "true", &secret),
				),
			},
		},
	})
}

func testAccCheckSecretMetadataV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	kmClient, err := config.keymanagerV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Keymanager client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_keymanager_secret_metadata" {
			continue
		}
		_, err = secrets.GetMetadata(kmClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Secret metadata (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckSecretMetadataV1Exists(n string, metadata *map[string]string) resource.TestCheckFunc {
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

		var found map[string]string

		found, err = secrets.GetMetadata(kmClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*metadata = found

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

		uuid := getUUIDfromSecretRef(secret.SecretRef)
		metadatum, _ := secrets.GetMetadatum(kmClient, uuid, key).Extract()
		if err != nil {
			return err
		}
		if metadatum.Value != value {
			return fmt.Errorf("Metadata does not match. Expected %v but got %v", metadatum, value)
		}

		return nil
	}
}

var testAccSecretMetadataV1_basic = fmt.Sprintf(`
resource "openstack_keymanager_secret_v1" "secret_1" {
		algorithm = "aes"
		bit_length = 256
		mode = "cbc"
		name = "mysecret"
		payload = "foobar"
		payload_content_type = "text/plain"
		secret_type = "passphrase"
	}

resource "openstack_keymanager_secret_metadata_v1" "metadata_1" {
		secret_ref = "${openstack_keymanager_secret_v1.secret_1.secret_ref}"
		metadata {
			foo = "bar"
		}
	}`)

var testAccSecretMetadataV1_Update = fmt.Sprintf(`
resource "openstack_keymanager_secret_v1" "secret_1" {
		algorithm = "aes"
		bit_length = 256
		mode = "cbc"
		name = "mysecret"
		payload = "foobar"
		payload_content_type = "text/plain"
		secret_type = "passphrase"
	}

resource "openstack_keymanager_secret_metadata_v1" "metadata_1" {
		secret_ref = "${openstack_keymanager_secret_v1.secret_1.secret_ref}"
		metadata {
			foo = "update"
			update = "true"
		}
	}`)
