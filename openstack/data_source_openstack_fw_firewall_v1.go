package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas/firewalls"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceFWFirewallV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFWFirewallV1Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				ConflictsWith: []string{"name"},
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				ConflictsWith: []string{"id"},
				Optional: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFWFirewallV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := firewalls.ListOpts{}
	errMessage := "id"
	if v, ok := d.GetOk("id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
		errMessage = "name"
	}

	pages, err := firewalls.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return err
	}

	allFWFirewalls, err := firewalls.ExtractFirewalls(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve openstack_fw_firewall_v1: %s", err)
	}

	if len(allFWFirewalls) < 1 {
		return fmt.Errorf("No openstack_fw_firewall_v1 found with %s: %s", errMessage ,d.Get(errMessage))
	}

	if len(allFWFirewalls) > 1 {
		return fmt.Errorf("More than one openstack_fw_firewall_v1 found with %s: %s",errMessage ,d.Get(errMessage))
	}

	firewall := allFWFirewalls[0]

	log.Printf("[DEBUG] Retrieved openstack_fw_firewall_v1 %s: %#v", firewall.ID, firewall)
	d.SetId(firewall.ID)

	d.SetId(firewall.ID)
	d.Set("name", firewall.Name)
	d.Set("tenant_id", firewall.TenantID)
	d.Set("description", firewall.Description)
	d.Set("admin_state_up", firewall.AdminStateUp)
	d.Set("status", firewall.Status)
	d.Set("policy_id", firewall.PolicyID)
	d.Set("project_id", firewall.ProjectID)
	d.Set("region", GetRegion(d, config))

	return nil
}
