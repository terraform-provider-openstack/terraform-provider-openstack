package openstack

import "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"

const computeV2FlavorDescriptionMicroversion = "2.55"

func expandComputeFlavorV2ExtraSpecs(raw map[string]any) flavors.ExtraSpecsOpts {
	extraSpecs := make(flavors.ExtraSpecsOpts, len(raw))
	for k, v := range raw {
		extraSpecs[k] = v.(string)
	}

	return extraSpecs
}
