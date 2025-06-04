package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
)

func dataSourceComputeServerGroupV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeServerGroupV2Read,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"rules": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeServerGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)
	computeClient, err := config.ComputeV2Client(ctx, region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	// Attempt to read with microversion 2.64
	computeClient.Microversion = "2.64"
	allPages, err := servergroups.List(computeClient, servergroups.ListOpts{}).AllPages(ctx)
	if err != nil {
		log.Printf("[DEBUG] Falling back to legacy API call due to: %#v", err)
		// fallback to legacy microversion
		computeClient.Microversion = ""

		allPages, err = servergroups.List(computeClient, servergroups.ListOpts{}).AllPages(ctx)
		if err != nil {
			return diag.Errorf("Error listing compute servergroups: %s", err)
		}
	}

	allServerGroups, err := servergroups.ExtractServerGroups(allPages)
	if err != nil {
		return diag.Errorf("Error extracting compute servergroups: %s", err)
	}

	name := d.Get("name").(string)

	var refinedServerGroups []servergroups.ServerGroup
	for _, servergroup := range allServerGroups {
		if servergroup.Name == name {
			refinedServerGroups = append(refinedServerGroups, servergroup)
		}
	}

	if len(refinedServerGroups) < 1 {
		return diag.Errorf("Could not find any servergroup with this name: %s", name)
	}
	if len(refinedServerGroups) > 1 {
		return diag.Errorf("More than one servergroup found with this name: %s", name)
	}

	sg := refinedServerGroups[0]

	d.SetId(sg.ID)
	d.Set("name", sg.Name)
	d.Set("user_id", sg.UserID)
	d.Set("project_id", sg.ProjectID)
	d.Set("members", sg.Members)
	d.Set("metadata", sg.Metadata)
	if sg.Policy != nil && *sg.Policy != "" {
		d.Set("policy", *sg.Policy)
	} else {
		d.Set("policy", sg.Policies[0])
	}
	if sg.Rules != nil {
		d.Set("rules", map[string]interface{}{"max_server_per_host": sg.Rules.MaxServerPerHost})
	} else {
		d.Set("rules", nil)
	}

	return nil
}
