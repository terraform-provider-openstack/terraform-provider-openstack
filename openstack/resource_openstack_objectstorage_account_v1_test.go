package openstack

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/accounts"
)

func TestAccObjectStorageV1Account_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSwift(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckObjectStorageV1AccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1AccountBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_account_v1.account_1", "metadata.%", "2"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_account_v1.account_1", "metadata.Test", "true"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_account_v1.account_1", "metadata.Uppertest", "true"),
				),
			},
			{
				Config: testAccObjectStorageV1AccountUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_account_v1.account_1", "metadata.%", "1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_account_v1.account_1", "metadata.Test", "true"),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Account_quota(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckSwift(t)
			t.Skip("Works only in selected environments")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckObjectStorageV1AccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1AccountQuota,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_account_v1.account_1", "metadata.%", "2"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_account_v1.account_1", "metadata.Test", "true"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_account_v1.account_1", "metadata.Quota-Bytes", "1024000"),
				),
			},
		},
	})
}

func testAccCheckObjectStorageV1AccountDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	objectStorageClient, err := config.ObjectStorageV1Client(context.TODO(), osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack object storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_objectstorage_account_v1" {
			continue
		}

		res := accounts.Get(context.TODO(), objectStorageClient, nil)
		metadata, err := res.ExtractMetadata()
		if err != nil {
			return fmt.Errorf("failed to retrieve account metadata: %s", err)
		}

		if len(metadata) > 1 {
			return fmt.Errorf("account metadata still exists")
		}
	}

	return nil
}

const testAccObjectStorageV1AccountBasic = `
resource "openstack_objectstorage_account_v1" "account_1" {
  metadata = {
    test = "true"
    upperTest = "true"
  }
}
`

const testAccObjectStorageV1AccountUpdate = `
resource "openstack_objectstorage_account_v1" "account_1" {
  metadata = {
    test = "true"
  }
}
`

const testAccObjectStorageV1AccountQuota = `
resource "openstack_objectstorage_account_v1" "account_1" {
  metadata = {
    test = "true"
    quota-bytes = "1024000"
  }
}
`
