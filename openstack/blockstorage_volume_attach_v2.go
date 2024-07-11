package openstack

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
)

func expandBlockStorageV2AttachMode(v string) (volumes.AttachMode, error) {
	var attachMode volumes.AttachMode
	var attachError error

	switch v {
	case "":
		attachMode = ""
	case "ro":
		attachMode = volumes.ReadOnly
	case "rw":
		attachMode = volumes.ReadWrite
	default:
		attachError = fmt.Errorf("Invalid attach_mode specified")
	}

	return attachMode, attachError
}

func blockStorageVolumeAttachV2ParseID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("Unable to determine openstack_blockstorage_volume_attach_v2 ID")
	}

	return parts[0], parts[1], nil
}
