package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLBLoadbalancerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBLoadbalancerV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name"},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"loadbalancer_id"},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vip_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vip_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vip_network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vip_qos_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"loadbalancer_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"listeners": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"pools": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags_not": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags_not_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"additional_vips": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceLBLoadbalancerV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancer client: %s", err)
	}

	listOpts := loadbalancers.ListOpts{
		Tags:       expandTagsList(d, "tags"),
		TagsAny:    expandTagsList(d, "tags_any"),
		TagsNot:    expandTagsList(d, "tags_not"),
		TagsNotAny: expandTagsList(d, "tags_not_any"),
	}

	if v, ok := d.GetOk("loadbalancer_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("vip_address"); ok {
		listOpts.VipAddress = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	allPages, err := loadbalancers.List(lbClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query OpenStack loadbalancer: %s", err)
	}

	allLoadbalancers, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve Openstack loadbalancer: %s", err)
	}

	if len(allLoadbalancers) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allLoadbalancers) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allLoadbalancers)

		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria")
	}

	dataSourceLBLoadbalancerV2Attributes(d, &allLoadbalancers[0])
	d.Set("region", GetRegion(d, config))

	return nil
}

func dataSourceLBLoadbalancerV2Attributes(d *schema.ResourceData, loadbalancer *loadbalancers.LoadBalancer) {
	log.Printf("[DEBUG] Retrieved openstack_lb_loadbalancer_v2 %s: %#v", loadbalancer.ID, loadbalancer)

	d.SetId(loadbalancer.ID)
	d.Set("description", loadbalancer.Description)
	d.Set("admin_state_up", loadbalancer.AdminStateUp)
	d.Set("project_id", loadbalancer.ProjectID)
	d.Set("provisioning_status", loadbalancer.ProvisioningStatus)
	d.Set("vip_address", loadbalancer.VipAddress)
	d.Set("vip_port_id", loadbalancer.VipPortID)
	d.Set("vip_subnet_id", loadbalancer.VipSubnetID)
	d.Set("vip_network_id", loadbalancer.VipNetworkID)
	d.Set("vip_qos_policy_id", loadbalancer.VipQosPolicyID)
	d.Set("operating_status", loadbalancer.OperatingStatus)
	d.Set("name", loadbalancer.Name)
	d.Set("flavor_id", loadbalancer.FlavorID)
	d.Set("availability_zone", loadbalancer.AvailabilityZone)
	d.Set("loadbalancer_provider", loadbalancer.Provider)
	d.Set("listeners", flattenLBListenersV2(loadbalancer.Listeners))
	d.Set("pools", flattenLBPoolsV2(loadbalancer.Pools))
	d.Set("tags", loadbalancer.Tags)
	d.Set("additional_vips", flattenLBAdditionalVIPsV2(loadbalancer.AdditionalVips))
}
