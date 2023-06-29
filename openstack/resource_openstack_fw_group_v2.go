package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/groups"
)

func resourceFWGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFWGroupV2Create,
		ReadContext:   resourceFWGroupV2Read,
		UpdateContext: resourceFWGroupV2Update,
		DeleteContext: resourceFWGroupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"project_id"},
			},

			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"tenant_id"},
			},

			"ingress_firewall_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"egress_firewall_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFWGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	groupcreateOpts := groups.CreateOpts{
		Name:                    d.Get("name").(string),
		TenantID:                d.Get("tenant_id").(string),
		ProjectID:               d.Get("project_id").(string),
		Description:             d.Get("description").(string),
		IngressFirewallPolicyID: d.Get("ingress_firewall_policy_id").(string),
		EgressFirewallPolicyID:  d.Get("egress_firewall_policy_id").(string),
	}

	if r, ok := d.GetOk("shared"); ok {
		shared := r.(bool)
		groupcreateOpts.Shared = &shared
	}

	if r, ok := d.GetOk("admin_state_up"); ok {
		adminStateUp := r.(bool)
		groupcreateOpts.AdminStateUp = &adminStateUp
	}

	associatedPortsRaw := d.Get("ports").(*schema.Set).List()
	if len(associatedPortsRaw) > 0 {
		var portIds []string
		for _, v := range associatedPortsRaw {
			portIds = append(portIds, v.(string))
		}

		groupcreateOpts.Ports = portIds
	}

	log.Printf("[DEBUG] openstack_fw_group_v2 create options: %#v", groupcreateOpts)

	group, err := groups.Create(networkingClient, groupcreateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_fw_group_v2: %s", err)
	}

	log.Printf("[DEBUG] Created openstack_fw_group_v2 %s: %#v", group.ID, group)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE", "INACTIVE", "DOWN"},
		Refresh:    fwGroupV2RefreshFunc(networkingClient, group.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_fw_group_v2 to become active: %s", err)
	}

	d.SetId(group.ID)

	return resourceFWGroupV2Read(ctx, d, meta)
}

func resourceFWGroupV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	group, err := groups.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_fw_group_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_fw_group_v2 %s: %#v", d.Id(), group)

	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("tenant_id", group.TenantID)
	d.Set("project_id", group.ProjectID)
	d.Set("ingress_firewall_policy_id", group.IngressFirewallPolicyID)
	d.Set("egress_firewall_policy_id", group.EgressFirewallPolicyID)
	d.Set("admin_state_up", group.AdminStateUp)
	d.Set("status", group.Status)
	d.Set("ports", group.Ports)
	d.Set("shared", group.Shared)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceFWGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	group, err := groups.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_fw_group_v2"))
	}

	var (
		hasChange  bool
		updateOpts groups.UpdateOpts
	)

	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("ingress_firewall_policy_id") {
		ingressFirewallPolicyID := d.Get("ingress_firewall_policy_id").(string)
		if ingressFirewallPolicyID == "" {
			log.Printf("[DEBUG] Attempting to clear ingress policy of openstack_fw_group_v2: %s.", group.ID)

			err := fwGroupV2IngressPolicyDeleteFunc(networkingClient, d, ctx, group.ID)
			if err != nil {
				return err
			}
		}
		if len(ingressFirewallPolicyID) > 0 {
			hasChange = true
			updateOpts.IngressFirewallPolicyID = &ingressFirewallPolicyID
		}
	}

	if d.HasChange("egress_firewall_policy_id") {
		egressFirewallPolicyID := d.Get("egress_firewall_policy_id").(string)
		if egressFirewallPolicyID == "" {
			log.Printf("[DEBUG] Attempting to clear egress policy of openstack_fw_group_v2: %s.", group.ID)

			err := fwGroupV2EgressPolicyDeleteFunc(networkingClient, d, ctx, group.ID)
			if err != nil {
				return err
			}
		}
		if len(egressFirewallPolicyID) > 0 {
			hasChange = true
			updateOpts.EgressFirewallPolicyID = &egressFirewallPolicyID
		}
	}

	if d.HasChange("shared") {
		hasChange = true
		shared := d.Get("shared").(bool)
		updateOpts.Shared = &shared
	}

	if d.HasChange("admin_state_up") {
		hasChange = true
		adminStateUp := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &adminStateUp
	}

	var portIds []string
	if d.HasChange("ports") {
		hasChange = true
		emptyList := make([]string, 0)
		updateOpts.Ports = &emptyList
		if _, ok := d.GetOk("ports"); ok {
			associatedPortsRaw := d.Get("ports").(*schema.Set).List()
			for _, v := range associatedPortsRaw {
				portIds = append(portIds, v.(string))
			}
			updateOpts.Ports = &portIds
		}
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_fw_group_v2 %s update options: %#v", d.Id(), updateOpts)

		_, err = groups.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_fw_group_v2 %s: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING_CREATE", "PENDING_UPDATE"},
			Target:     []string{"ACTIVE", "INACTIVE", "DOWN"},
			Refresh:    fwGroupV2RefreshFunc(networkingClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for openstack_fw_group_v2 %s to become active: %s", d.Id(), err)
		}
	}

	return resourceFWGroupV2Read(ctx, d, meta)
}

func resourceFWGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	group, err := groups.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_fw_group_v2"))
	}

	if len(group.Ports) > 0 {
		var updateGroupOpts groups.UpdateOpts
		emptyPorts := []string{}
		updateGroupOpts.Ports = &emptyPorts
		_, err := groups.Update(networkingClient, group.ID, updateGroupOpts).Extract()
		if err != nil {
			return diag.Errorf("Error removing ports from openstack_fw_group_v2 %s: %s", d.Id(), err)
		}
	}

	if group.IngressFirewallPolicyID != "" {
		diagErr := fwGroupV2IngressPolicyDeleteFunc(networkingClient, d, ctx, group.ID)
		if diagErr != nil {
			return diagErr
		}
	}

	if group.EgressFirewallPolicyID != "" {
		diagErr := fwGroupV2EgressPolicyDeleteFunc(networkingClient, d, ctx, group.ID)
		if diagErr != nil {
			return diagErr
		}
	}

	err = groups.Delete(networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting openstack_fw_group_v2 %s: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    fwGroupV2DeleteFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_fw_firewall_v2 %s to delete:  %s", d.Id(), err)
	}

	return nil
}
