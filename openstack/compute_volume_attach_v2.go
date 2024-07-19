package openstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
)

func computeVolumeAttachV2AttachFunc(computeClient *gophercloud.ServiceClient, blockStorageClient *gophercloud.ServiceClient, instanceID, attachmentID string, volumeID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		va, err := volumeattach.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "ATTACHING", nil
			}
			return va, "", err
		}

		// Block Storage client will be empty if "ignore_volume_confirmation" == true.
		if blockStorageClient == nil {
			return va, "ATTACHED", nil
		}

		v, err := volumes.Get(blockStorageClient, volumeID).Extract()
		if err != nil {
			return va, "", err
		}
		if v.Status == "error" {
			return va, "", fmt.Errorf("volume entered unexpected error status")
		}
		if v.Status != "in-use" {
			return va, "ATTACHING", nil
		}

		return va, "ATTACHED", nil
	}
}

func computeVolumeAttachV2DetachFunc(computeClient *gophercloud.ServiceClient, instanceID, attachmentID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] openstack_compute_volume_attach_v2 attempting to detach OpenStack volume %s from instance %s",
			attachmentID, instanceID)

		va, err := volumeattach.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}
			return va, "", err
		}

		err = volumeattach.Delete(computeClient, instanceID, attachmentID).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}

			if _, ok := err.(gophercloud.ErrDefault400); ok {
				return nil, "", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] openstack_compute_volume_attach_v2 (%s/%s) is still active.", instanceID, attachmentID)
		return nil, "", nil
	}
}
