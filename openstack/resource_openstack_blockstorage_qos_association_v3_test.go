package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/qos"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
)

func TestAccBlockstorageV3QosAssociation_basic(t *testing.T) {
	var qosTest qos.QoS
	var qosName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	var vt volumetypes.VolumeType
	var vtName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockstorageV3QosAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockstorageV3QosAssociationBasic(qosName, vtName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageQosV3Exists("openstack_blockstorage_qos_v3.qos", &qosTest),
					testAccCheckBlockStorageVolumeTypeV3Exists("openstack_blockstorage_volume_type_v3.volume_type_1", &vt),
					testAccCheckBlockstorageV3QosAssociationExists("openstack_blockstorage_qos_association_v3.qos_association"),
					resource.TestCheckResourceAttrPtr(
						"openstack_blockstorage_qos_association_v3.qos_association", "qos_id", &qosTest.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_blockstorage_qos_association_v3.qos_association", "volume_type_id", &vt.ID),
				),
			},
		},
	})
}

func testAccCheckBlockstorageV3QosAssociationDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_blockstorage_qos_association_v3" {
			continue
		}

		qosID, vtID, err := parseQosAssociationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		allPages, err := qos.ListAssociations(blockStorageClient, qosID).AllPages()
		if err == nil {
			allAssociations, err := qos.ExtractAssociations(allPages)
			if err == nil {
				for _, association := range allAssociations {
					if association.ID == vtID {
						return fmt.Errorf("Qos association still exists")
					}
				}
			}
		}
	}

	return nil
}

func testAccCheckBlockstorageV3QosAssociationExists(n string) resource.TestCheckFunc {
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

		qosID, vtID, err := parseQosAssociationID(rs.Primary.ID)
		if err != nil {
			return err
		}

		allPages, err := qos.ListAssociations(blockStorageClient, qosID).AllPages()
		if err != nil {
			return fmt.Errorf("Error retrieving associations for qos: %s", qosID)
		}

		allAssociations, err := qos.ExtractAssociations(allPages)
		if err != nil {
			return fmt.Errorf("Error extracting associations for qos: %s", qosID)
		}

		found := false
		for _, association := range allAssociations {
			if association.ID == vtID {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Qos association not found for qosID/vtID: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccBlockstorageV3QosAssociationBasic(qosName, vtName string) string {
	return fmt.Sprintf(`
resource "openstack_blockstorage_qos_v3" "qos" {
  name = "%s"
  consumer = "front-end"
  specs = {
	  read_iops_sec = "20000"
  }
}

resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
  name = "%s"
}

resource "openstack_blockstorage_qos_association_v3" "qos_association" {
  qos_id         = "${openstack_blockstorage_qos_v3.qos.id}"
  volume_type_id = "${openstack_blockstorage_volume_type_v3.volume_type_1.id}"
}
`, qosName, vtName)
}
