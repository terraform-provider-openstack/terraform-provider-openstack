package openstack

import (
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
)

func expandBlockStorageVolumeTypeV3ExtraSpecs(raw map[string]any) volumetypes.ExtraSpecsOpts {
	extraSpecs := make(volumetypes.ExtraSpecsOpts, len(raw))
	for k, v := range raw {
		extraSpecs[k] = v.(string)
	}

	return extraSpecs
}
