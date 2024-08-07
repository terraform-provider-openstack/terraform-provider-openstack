package openstack

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
)

func computeInterfaceAttachV2AttachFunc(
	computeClient *gophercloud.ServiceClient, instanceID, attachmentID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		va, err := attachinterfaces.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "ATTACHING", nil
			}
			return va, "", err
		}

		return va, "ATTACHED", nil
	}
}

func computeInterfaceAttachV2DetachFunc(
	computeClient *gophercloud.ServiceClient, instanceID, attachmentID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to detach openstack_compute_interface_attach_v2 %s from instance %s",
			attachmentID, instanceID)

		va, err := attachinterfaces.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}
			return va, "", err
		}

		err = attachinterfaces.Delete(computeClient, instanceID, attachmentID).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}

			if _, ok := err.(gophercloud.ErrDefault400); ok {
				return nil, "", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] openstack_compute_interface_attach_v2 %s is still active.", attachmentID)
		return nil, "", nil
	}
}
