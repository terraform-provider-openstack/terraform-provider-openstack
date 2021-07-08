package openstack

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/snapshots"
)

func TestAccSFSV2SnapshotDataSource_basic(t *testing.T) {
	var shareID string

	if os.Getenv("TF_ACC") != "" {
		snapshot, err := testAccSFSV2SnapshotCreate(t, "test_snapshot")
		if err != nil {
			t.Fatal(err)
		}
		shareID = snapshot.ShareID
		defer testAccSFSV2SnapshotDelete(t, snapshot) //nolint:errcheck
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSFS(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2SnapshotDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2SnapshotDataSourceID("data.openstack_sharedfilesystem_snapshot_v2.snapshot_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_snapshot_v2.snapshot_1", "name", "test_snapshot"),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_snapshot_v2.snapshot_1", "share_id", shareID),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_snapshot_v2.snapshot_1", "share_proto", "NFS"),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_snapshot_v2.snapshot_1", "share_size", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_snapshot_v2.snapshot_1", "size", "1"),
				),
			},
		},
	})
}

func waitForShareStatus(t *testing.T, c *gophercloud.ServiceClient, id, status string, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		current, err := shares.Get(c, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				switch status {
				case "deleted":
					return true, nil
				default:
					return false, err
				}
			}
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		if strings.Contains(current.Status, "error") {
			return true, fmt.Errorf("An error occurred, wrong status: %s", current.Status)
		}

		return false, nil
	})
}

func waitForSnapshotStatus(t *testing.T, c *gophercloud.ServiceClient, id, status string, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		current, err := snapshots.Get(c, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				switch status {
				case "deleted":
					return true, nil
				default:
					return false, err
				}
			}
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		if strings.Contains(current.Status, "error") {
			return true, fmt.Errorf("An error occurred, wrong status: %s", current.Status)
		}

		return false, nil
	})
}

func testAccSFSV2SnapshotCreate(t *testing.T, snapshotName string) (*snapshots.Snapshot, error) {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return nil, err
	}

	client, err := config.SharedfilesystemV2Client(osRegionName)
	if err != nil {
		return nil, err
	}

	createShareOpts := shares.CreateOpts{
		Size:       1,
		Name:       "test",
		ShareProto: "NFS",
		ShareType:  "dhss_false",
	}

	share, err := shares.Create(client, createShareOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Share %s created, waiting for 'available' status", share.ID)

	err = waitForShareStatus(t, client, share.ID, "available", 600)
	if err != nil {
		nErr := shares.Delete(client, share.ID).ExtractErr()
		if nErr != nil {
			return nil, fmt.Errorf("Unable to get share available status (%s) and Delete:  %s)", err, nErr)
		}
		return nil, err
	}

	createOpts := snapshots.CreateOpts{
		ShareID: share.ID,
		Name:    snapshotName,
	}

	snapshot, err := snapshots.Create(client, createOpts).Extract()
	if err != nil {
		nErr := shares.Delete(client, share.ID).ExtractErr()
		if nErr != nil {
			return nil, fmt.Errorf("Unable to create snapshot (%s) and delete share (%s: %s)", err, share.ID, nErr)
		}
		return nil, err
	}

	t.Logf("Snapshot %s created, waiting for 'available' status", snapshot.ID)

	if err := waitForSnapshotStatus(t, client, snapshot.ID, "available", 600); err != nil {
		t.Logf("%s", err)
	}

	return snapshot, nil
}

func testAccSFSV2SnapshotDelete(t *testing.T, snapshot *snapshots.Snapshot) error {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return err
	}

	client, err := config.SharedfilesystemV2Client(osRegionName)
	if err != nil {
		return err
	}

	err = snapshots.Delete(client, snapshot.ID).ExtractErr()
	if err != nil {
		return err
	}

	if err := waitForSnapshotStatus(t, client, snapshot.ID, "deleted", 600); err != nil {
		t.Logf("%s", err)
	}

	t.Logf("Snapshot %s deleted", snapshot.ID)

	err = shares.Delete(client, snapshot.ShareID).ExtractErr()
	if err != nil {
		return err
	}

	if err := waitForShareStatus(t, client, snapshot.ShareID, "deleted", 600); err != nil {
		t.Logf("%s", err)
	}

	t.Logf("Share %s deleted", snapshot.ShareID)

	return nil
}

func testAccCheckSFSV2SnapshotDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find snapshot data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Snapshot data source ID not set")
		}

		return nil
	}
}

const testAccSFSV2SnapshotDataSourceBasic = `
data "openstack_sharedfilesystem_snapshot_v2" "snapshot_1" {
  name = "test_snapshot"
}
`
