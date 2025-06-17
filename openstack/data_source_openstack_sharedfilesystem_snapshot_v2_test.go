package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/snapshots"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSFSV2SnapshotDataSource_basic(t *testing.T) {
	var shareID string

	if os.Getenv("TF_ACC") != "" {
		snapshot, err := testAccSFSV2SnapshotCreate(t, "test_snapshot")
		if err != nil {
			t.Fatal(err)
		}

		shareID = snapshot.ShareID

		defer testAccSFSV2SnapshotDelete(t, snapshot)
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

func waitForShareStatus(ctx context.Context, c *gophercloud.ServiceClient, id, status string, secs int) error {
	ctx1, cancel := context.WithTimeout(ctx, time.Duration(secs)*time.Second)
	defer cancel()

	return gophercloud.WaitFor(ctx1, func(ctx1 context.Context) (bool, error) {
		current, err := shares.Get(ctx1, c, id).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
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

func waitForSnapshotStatus(ctx context.Context, c *gophercloud.ServiceClient, id, status string, secs int) error {
	ctx1, cancel := context.WithTimeout(ctx, time.Duration(secs)*time.Second)
	defer cancel()

	return gophercloud.WaitFor(ctx1, func(ctx1 context.Context) (bool, error) {
		current, err := snapshots.Get(ctx1, c, id).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
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
	config, err := testAccAuthFromEnv(t.Context())
	if err != nil {
		return nil, err
	}

	client, err := config.SharedfilesystemV2Client(t.Context(), osRegionName)
	if err != nil {
		return nil, err
	}

	createShareOpts := shares.CreateOpts{
		Size:       1,
		Name:       "test",
		ShareProto: "NFS",
		ShareType:  "dhss_false",
	}

	share, err := shares.Create(t.Context(), client, createShareOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Share %s created, waiting for 'available' status", share.ID)

	err = waitForShareStatus(t.Context(), client, share.ID, "available", 600)
	if err != nil {
		nErr := shares.Delete(t.Context(), client, share.ID).ExtractErr()
		if nErr != nil {
			return nil, fmt.Errorf("Unable to get share available status (%w) and Delete:  %w)", err, nErr)
		}

		return nil, err
	}

	createOpts := snapshots.CreateOpts{
		ShareID: share.ID,
		Name:    snapshotName,
	}

	snapshot, err := snapshots.Create(t.Context(), client, createOpts).Extract()
	if err != nil {
		nErr := shares.Delete(t.Context(), client, share.ID).ExtractErr()
		if nErr != nil {
			return nil, fmt.Errorf("Unable to create snapshot (%w) and delete share (%s: %w)", err, share.ID, nErr)
		}

		return nil, err
	}

	t.Logf("Snapshot %s created, waiting for 'available' status", snapshot.ID)

	if err := waitForSnapshotStatus(t.Context(), client, snapshot.ID, "available", 600); err != nil {
		return nil, err
	}

	return snapshot, nil
}

func testAccSFSV2SnapshotDelete(t *testing.T, snapshot *snapshots.Snapshot) {
	config, err := testAccAuthFromEnv(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	client, err := config.SharedfilesystemV2Client(t.Context(), osRegionName)
	if err != nil {
		t.Fatal(err)
	}

	err = snapshots.Delete(t.Context(), client, snapshot.ID).ExtractErr()
	if err != nil {
		t.Fatal(err)
	}

	if err := waitForSnapshotStatus(t.Context(), client, snapshot.ID, "deleted", 600); err != nil {
		t.Fatal(err)
	}

	t.Logf("Snapshot %s deleted", snapshot.ID)

	err = shares.Delete(t.Context(), client, snapshot.ShareID).ExtractErr()
	if err != nil {
		t.Fatal(err)
	}

	if err := waitForShareStatus(t.Context(), client, snapshot.ShareID, "deleted", 600); err != nil {
		t.Fatal(err)
	}

	t.Logf("Share %s deleted", snapshot.ShareID)
}

func testAccCheckSFSV2SnapshotDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find snapshot data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Snapshot data source ID not set")
		}

		return nil
	}
}

const testAccSFSV2SnapshotDataSourceBasic = `
data "openstack_sharedfilesystem_snapshot_v2" "snapshot_1" {
  name = "test_snapshot"
}
`
