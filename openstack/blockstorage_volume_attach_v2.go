package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
)

func expandBlockStorageV2AttachMode(v string) (volumeactions.AttachMode, error) {
	var attachMode volumeactions.AttachMode
	var attachError error

	switch v {
	case "":
		attachMode = ""
	case "ro":
		attachMode = volumeactions.ReadOnly
	case "rw":
		attachMode = volumeactions.ReadWrite
	default:
		attachError = fmt.Errorf("Invalid attach_mode specified")
	}

	return attachMode, attachError
}
