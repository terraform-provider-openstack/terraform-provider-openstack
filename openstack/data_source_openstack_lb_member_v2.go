package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLBMemberV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBMemberV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"member_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name"},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"member_id"},
			},

			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"protocol_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"backup": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"monitor_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"monitor_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceLBMemberV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancer client: %s", err)
	}

	listOpts := pools.ListMembersOpts{
		ID:           d.Get("member_id").(string),
		Name:         d.Get("name").(string),
		Weight:       d.Get("weight").(int),
		Address:      d.Get("address").(string),
		ProtocolPort: d.Get("protocol_port").(int),
	}

	poolID := d.Get("pool_id").(string)

	allPages, err := pools.ListMembers(lbClient, poolID, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query OpenStack loadbalancer members: %s", err)
	}

	allMembers, err := pools.ExtractMembers(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve Openstack loadbalancer members: %s", err)
	}

	requestedTags := expandTagsList(d, "tags")
	if len(requestedTags) > 0 {
		var filtered []pools.Member

		for _, m := range allMembers {
			if hasAllRequestedTags(m.Tags, requestedTags) {
				filtered = append(filtered, m)
			}
		}

		allMembers = filtered
	}

	if len(allMembers) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allMembers) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allMembers)

		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria")
	}

	dataSourceLBMemberV2Attributes(d, &allMembers[0])
	d.Set("region", GetRegion(d, config))

	return nil
}

func dataSourceLBMemberV2Attributes(d *schema.ResourceData, member *pools.Member) {
	log.Printf("[DEBUG] Retrieved openstack_lb_member_v2 %s: %#v", member.ID, member)

	d.SetId(member.ID)
	d.Set("name", member.Name)
	d.Set("project_id", member.ProjectID)
	d.Set("weight", member.Weight)
	d.Set("admin_state_up", member.AdminStateUp)
	d.Set("subnet_id", member.SubnetID)
	d.Set("pool_id", member.PoolID)
	d.Set("address", member.Address)
	d.Set("protocol_port", member.ProtocolPort)
	d.Set("provisioning_status", member.ProvisioningStatus)
	d.Set("operating_status", member.OperatingStatus)
	d.Set("backup", member.Backup)
	d.Set("monitor_address", member.MonitorAddress)
	d.Set("monitor_port", member.MonitorPort)
	d.Set("tags", member.Tags)
}
