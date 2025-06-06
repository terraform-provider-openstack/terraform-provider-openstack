package openstack

import (
	"errors"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
)

func expandBlockStorageV3AttachMode(v string) (volumes.AttachMode, error) {
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
		attachError = errors.New("Invalid attach_mode specified")
	}

	return attachMode, attachError
}
