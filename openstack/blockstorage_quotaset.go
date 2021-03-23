package openstack

import (
	"fmt"
	"strconv"
)

// blockStorageVolumeTypeQuotaConversion converts all values of the map to int.
func blockStorageVolumeTypeQuotaConversion(vtq map[string]interface{}) (map[string]interface{}, error) {
	newVTQ := make(map[string]interface{})
	for oldKey, oldVal := range vtq {
		tmp, ok := oldVal.(string)
		if !ok {
			return nil, fmt.Errorf("Error asserting type for %+v", oldVal)
		}
		newVal, err := strconv.Atoi(tmp)
		if err != nil {
			return nil, fmt.Errorf("Error converting to string for %s", tmp)
		}
		newVTQ[oldKey] = newVal
	}
	return newVTQ, nil
}
