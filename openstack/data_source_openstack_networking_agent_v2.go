package openstack

import (
	"context"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/agents"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworkingAgentV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingAgentV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"agent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"agent_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"alive": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"binary": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"host": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"topic": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"resources_synced": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"configurations": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"started_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"heartbeat_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkingAgentV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := agents.ListOpts{}

	if v, ok := d.GetOk("agent_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("agent_type"); ok {
		listOpts.AgentType = v.(string)
	}

	if v, ok := getOkExists(d, "alive"); ok {
		boolVal := v.(bool)
		listOpts.Alive = &boolVal
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		listOpts.AvailabilityZone = v.(string)
	}

	if v, ok := d.GetOk("binary"); ok {
		listOpts.Binary = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOk("host"); ok {
		listOpts.Host = v.(string)
	}

	if v, ok := d.GetOk("topic"); ok {
		listOpts.Topic = v.(string)
	}

	pages, err := agents.List(networkingClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to list openstack_networking_agent_v2: %s", err)
	}

	allAgents, err := agents.ExtractAgents(pages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_networking_agent_v2: %s", err)
	}

	if len(allAgents) < 1 {
		return diag.Errorf("No openstack_networking_agent_v2 found")
	}

	if len(allAgents) > 1 {
		return diag.Errorf("More than one openstack_networking_agent_v2 found")
	}

	a := allAgents[0]

	log.Printf("[DEBUG] Retrieved openstack_networking_agent_v2 %s: %+v", a.ID, a)
	d.SetId(a.ID)

	d.Set("region", GetRegion(d, config))
	d.Set("agent_type", a.AgentType)
	d.Set("alive", a.Alive)
	d.Set("availability_zone", a.AvailabilityZone)
	d.Set("binary", a.Binary)
	d.Set("description", a.Description)
	d.Set("host", a.Host)
	d.Set("topic", a.Topic)
	d.Set("admin_state_up", a.AdminStateUp)
	d.Set("resources_synced", a.ResourcesSynced)
	d.Set("configurations", mapNetworkAgentConfigurations(a.Configurations))
	d.Set("created_at", a.CreatedAt.String())
	d.Set("started_at", a.StartedAt.String())
	d.Set("heartbeat_timestamp", a.HeartbeatTimestamp.String())

	return nil
}

func mapNetworkAgentConfigurations(cfg map[string]any) map[string]string {
	m := make(map[string]string, len(cfg))

	for key, val := range cfg {
		m[key] = fmt.Sprintf("%v", val)
	}

	return m
}
