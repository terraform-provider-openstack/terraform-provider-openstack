package openstack

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/volumeattach"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func computeVolumeAttachV2AttachFunc(ctx context.Context, computeClient *gophercloud.ServiceClient, blockStorageClient *gophercloud.ServiceClient, instanceID, attachmentID string, volumeID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		va, err := volumeattach.Get(ctx, computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return va, "ATTACHING", nil
			}

			return va, "", err
		}

		// Block Storage client will be empty if "ignore_volume_confirmation" == true.
		if blockStorageClient == nil {
			return va, "ATTACHED", nil
		}

		v, err := volumes.Get(ctx, blockStorageClient, volumeID).Extract()
		if err != nil {
			return va, "", err
		}

		if v.Status == "error" {
			return va, "", errors.New("volume entered unexpected error status")
		}

		if v.Status != "in-use" {
			return va, "ATTACHING", nil
		}

		return va, "ATTACHED", nil
	}
}

func computeVolumeAttachV2DetachFunc(ctx context.Context, computeClient *gophercloud.ServiceClient, instanceID, attachmentID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		log.Printf("[DEBUG] openstack_compute_volume_attach_v2 attempting to detach OpenStack volume %s from instance %s",
			attachmentID, instanceID)

		va, err := volumeattach.Get(ctx, computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return va, "DETACHED", nil
			}

			return va, "", err
		}

		err = volumeattach.Delete(ctx, computeClient, instanceID, attachmentID).ExtractErr()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return va, "DETACHED", nil
			}

			if gophercloud.ResponseCodeIs(err, http.StatusBadRequest) {
				return nil, "", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] openstack_compute_volume_attach_v2 (%s/%s) is still active.", instanceID, attachmentID)

		return nil, "", nil
	}
}
