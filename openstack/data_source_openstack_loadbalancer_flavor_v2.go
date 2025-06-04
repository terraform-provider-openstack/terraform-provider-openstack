package openstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceLoadBalancerFlavorV2 uses the same Read function as
// dataSourceLBFlavorV2 but includes a deprecation message.
func dataSourceLoadBalancerFlavorV2() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceLBFlavorV2Read,
		DeprecationMessage: "Use openstack_lb_flavor_v2 instead.",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flavor_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name"},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"flavor_id"},
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"flavor_profile_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
