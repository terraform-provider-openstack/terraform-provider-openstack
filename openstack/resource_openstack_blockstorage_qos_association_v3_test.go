package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/qos"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBlockstorageV3QosAssociation_basic(t *testing.T) {
	var qosTest qos.QoS

	qosName := "ACCPTTEST-" + acctest.RandString(5)

	var vt volumetypes.VolumeType

	vtName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockstorageV3QosAssociationDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockstorageV3QosAssociationBasic(qosName, vtName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageQosV3Exists(t.Context(), "openstack_blockstorage_qos_v3.qos", &qosTest),
					testAccCheckBlockStorageVolumeTypeV3Exists(t.Context(), "openstack_blockstorage_volume_type_v3.volume_type_1", &vt),
					testAccCheckBlockstorageV3QosAssociationExists(t.Context(), "openstack_blockstorage_qos_association_v3.qos_association"),
					resource.TestCheckResourceAttrPtr(
						"openstack_blockstorage_qos_association_v3.qos_association", "qos_id", &qosTest.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_blockstorage_qos_association_v3.qos_association", "volume_type_id", &vt.ID),
				),
			},
		},
	})
}

func testAccCheckBlockstorageV3QosAssociationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		blockStorageClient, err := config.BlockStorageV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_blockstorage_qos_association_v3" {
				continue
			}

			qosID, vtID, err := parsePairedIDs(rs.Primary.ID, "openstack_blockstorage_qos_association_v3")
			if err != nil {
				return err
			}

			allPages, err := qos.ListAssociations(blockStorageClient, qosID).AllPages(ctx)
			if err == nil {
				allAssociations, err := qos.ExtractAssociations(allPages)
				if err == nil {
					for _, association := range allAssociations {
						if association.ID == vtID {
							return errors.New("Qos association still exists")
						}
					}
				}
			}
		}

		return nil
	}
}

func testAccCheckBlockstorageV3QosAssociationExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		blockStorageClient, err := config.BlockStorageV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %w", err)
		}

		qosID, vtID, err := parsePairedIDs(rs.Primary.ID, "openstack_blockstorage_qos_association_v3")
		if err != nil {
			return err
		}

		allPages, err := qos.ListAssociations(blockStorageClient, qosID).AllPages(ctx)
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
  qos_id         = openstack_blockstorage_qos_v3.qos.id
  volume_type_id = openstack_blockstorage_volume_type_v3.volume_type_1.id
}
`, qosName, vtName)
}
