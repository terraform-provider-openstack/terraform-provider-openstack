package openstack

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOpenStackComputeInstanceV2ImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1)

	err := resourceComputeInstanceV2Read(d, meta)
	if err != nil {
		return nil, err
	}

	results[0] = d

	return results, nil
}
