package openstack

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceObjectStorageContainerV1V0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"container_read": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"container_sync_to": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"container_sync_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"container_write": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"versioning": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"versions", "history",
							}, true),
						},
						"location": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"storage_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceObjectStorageContainerStateUpgradeV0(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState["versioning_legacy"] = rawState["versioning"]
	rawState["versioning"] = false

	return rawState, nil
}
