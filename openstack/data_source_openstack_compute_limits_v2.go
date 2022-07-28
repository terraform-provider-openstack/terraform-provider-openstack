package openstack

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/limits"
)

func dataSourceComputeLimitsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeLimitsV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"max_total_cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_image_meta": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_server_meta": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_personality": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_personality_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_total_keypairs": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_security_groups": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_security_group_rules": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_server_groups": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_server_group_members": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_total_floating_ips": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_total_instances": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_total_ram_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"total_cores_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"total_instances_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"total_floating_ips_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"total_ram_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"total_security_groups_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"total_server_groups_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeLimitsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)
	computeClient, err := config.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	projectID := d.Get("project_id").(string)
	getOpts := limits.GetOpts{
		TenantID: projectID,
	}

	q, err := limits.Get(computeClient, getOpts).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_compute_limits_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_compute_limits_v2 %s: %#v", d.Id(), q)

	id := fmt.Sprintf("%s/%s", projectID, region)
	d.SetId(id)
	d.Set("region", region)
	d.Set("project_id", projectID)

	d.Set("max_total_cores", q.Absolute.MaxTotalCores)
	d.Set("max_image_meta", q.Absolute.MaxImageMeta)
	d.Set("max_server_meta", q.Absolute.MaxServerMeta)
	d.Set("max_personality", q.Absolute.MaxPersonality)
	d.Set("max_personality_size", q.Absolute.MaxPersonalitySize)
	d.Set("max_total_keypairs", q.Absolute.MaxTotalKeypairs)
	d.Set("max_security_groups", q.Absolute.MaxSecurityGroups)
	d.Set("max_security_group_rules", q.Absolute.MaxSecurityGroupRules)
	d.Set("max_server_groups", q.Absolute.MaxServerGroups)
	d.Set("max_server_group_members", q.Absolute.MaxServerGroupMembers)
	d.Set("max_total_floating_ips", q.Absolute.MaxTotalFloatingIps)
	d.Set("max_total_instances", q.Absolute.MaxTotalInstances)
	d.Set("max_total_ram_size", q.Absolute.MaxTotalRAMSize)
	d.Set("total_cores_used", q.Absolute.TotalCoresUsed)
	d.Set("total_instances_used", q.Absolute.TotalInstancesUsed)
	d.Set("total_floating_ips_used", q.Absolute.TotalFloatingIpsUsed)
	d.Set("total_ram_used", q.Absolute.TotalRAMUsed)
	d.Set("total_security_groups_used", q.Absolute.TotalSecurityGroupsUsed)
	d.Set("total_server_groups_used", q.Absolute.TotalServerGroupsUsed)

	return nil
}
