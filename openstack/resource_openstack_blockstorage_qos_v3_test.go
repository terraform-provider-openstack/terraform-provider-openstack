package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/qos"
)

func TestAccBlockStorageQosV3_basic(t *testing.T) {
	var qosTest qos.QoS

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageQosV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageQosV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageQosV3Exists("openstack_blockstorage_qos_v3.qos", &qosTest),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "name", "foo"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "consumer", "front-end"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "specs.%", "1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "specs.read_iops_sec", "20000"),
				),
			},
			{
				Config: testAccBlockStorageQosV3Update1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageQosV3Exists("openstack_blockstorage_qos_v3.qos", &qosTest),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "name", "foo"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "consumer", "back-end"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "specs.%", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "specs.read_iops_sec", "40000"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "specs.write_iops_sec", "40000"),
				),
			},
			{
				Config: testAccBlockStorageQosV3Update2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageQosV3Exists("openstack_blockstorage_qos_v3.qos", &qosTest),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "name", "foo"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "consumer", "back-end"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_qos_v3.qos", "specs.%", "0"),
				),
			},
		},
	})
}

func testAccCheckBlockStorageQosV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_blockstorage_qos_v3" {
			continue
		}

		_, err := qos.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Qos still exists")
		}
	}

	return nil
}

func testAccCheckBlockStorageQosV3Exists(n string, qosTest *qos.QoS) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.BlockStorageV3Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
		}

		found, err := qos.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Qos not found")
		}

		*qosTest = *found

		return nil
	}
}

const testAccBlockStorageQosV3Basic = `
resource "openstack_blockstorage_qos_v3" "qos" {
	name = "foo"
	consumer = "front-end"
    specs = {
		read_iops_sec = "20000"
	}

}
`

const testAccBlockStorageQosV3Update1 = `
resource "openstack_blockstorage_qos_v3" "qos" {
	name = "foo"
	consumer = "back-end"
    specs = {
		read_iops_sec = "40000"
		write_iops_sec = "40000"
	}

}
`

const testAccBlockStorageQosV3Update2 = `
resource "openstack_blockstorage_qos_v3" "qos" {
	name = "foo"
	consumer = "back-end"
    specs = {
	}
}
`
