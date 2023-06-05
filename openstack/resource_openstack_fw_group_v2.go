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

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"ingress_firewall_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"egress_firewall_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
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

	shared := d.Get("shared").(bool)
	is_admin_state_up := d.Get("admin_state_up").(bool)
	createOpts := GroupCreateOpts{
		groups.CreateOpts{
			Name:                    d.Get("name").(string),
			TenantID:                d.Get("tenant_id").(string),
			Description:             d.Get("description").(string),
			IngressFirewallPolicyID: d.Get("ingress_firewall_policy_id").(string),
			EgressFirewallPolicyID:  d.Get("egress_firewall_policy_id").(string),
			Shared:                  &shared,
			AdminStateUp:            &is_admin_state_up,
		},
		MapValueSpecs(d),
	}

	associatedPortsRaw := d.Get("ports").(*schema.Set).List()
	if len(associatedPortsRaw) > 0 {
		var portIds []string
		for _, v := range associatedPortsRaw {
			portIds = append(portIds, v.(string))
		}

		createOpts.Ports = portIds
	}

	log.Printf("[DEBUG] openstack_fw_group_v2 create options: %#v", createOpts)

	group, err := groups.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_fw_group_v2: %s", err)
	}

	log.Printf("[DEBUG] Created openstack_fw_group_v2 %s: %#v", group.ID, group)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE", "INACTIVE", "DOWn"},
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
	d.Set("shared", group.Shared)
	d.Set("admin_state_up", group.AdminStateUp)
	d.Set("status", group.Status)
	d.Set("tenant_id", group.TenantID)
	d.Set("project_id", group.ProjectID)
	d.Set("ingress_firewall_policy_id", group.IngressFirewallPolicyID)
	d.Set("egress_firewall_policy_id", group.EgressFirewallPolicyID)
	d.Set("ports", group.Ports)
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

	var updateOpts groups.UpdateOpts

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("ingress_firewall_policy_id") {
		ingressFirewallPolicyID := d.Get("ingress_firewall_policy_id").(string)
		if ingressFirewallPolicyID == "" {
			_, err := groups.RemoveIngressPolicy(networkingClient, group.ID).Extract()
			if err != nil {
				return diag.Errorf("Error removing ingress firewall policy from openstack_fw_group_v2 %s: %s", d.Id(), err)
			}
		}
		if len(ingressFirewallPolicyID) > 0 {
			updateOpts.IngressFirewallPolicyID = &ingressFirewallPolicyID
		}
	}

	if d.HasChange("egress_firewall_policy_id") {
		egressFirewallPolicyID := d.Get("egress_firewall_policy_id").(string)
		if egressFirewallPolicyID == "" {
			_, err := groups.RemoveEgressPolicy(networkingClient, group.ID).Extract()
			if err != nil {
				return diag.Errorf("Error removing ingress firewall policy from openstack_fw_group_v2 %s: %s", d.Id(), err)
			}
		}
		if len(egressFirewallPolicyID) > 0 {
			updateOpts.EgressFirewallPolicyID = &egressFirewallPolicyID
		}
	}

	if d.HasChange("shared") {
		shared := d.Get("shared").(bool)
		updateOpts.Shared = &shared
	}

	if d.HasChange("admin_state_up") {
		admin_state_up := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &admin_state_up
	}

	var portIds []string
	if d.HasChange("ports") {
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

	log.Printf("[DEBUG] openstack_fw_group_v2 %s update options: %#v", d.Id(), updateOpts)
	err = groups.Update(networkingClient, d.Id(), updateOpts).Err
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
		newPorts := []string{}
		updateGroupOpts.Ports = &newPorts
		_, err := groups.Update(networkingClient, group.ID, updateGroupOpts).Extract()
		if err != nil {
			return diag.Errorf("Error removing ports from openstack_fw_group_v2 %s: %s", d.Id(), err)
		}
	}

	if group.IngressFirewallPolicyID != "" {
		_, err := groups.RemoveIngressPolicy(networkingClient, group.ID).Extract()
		if err != nil {
			return diag.Errorf("Error removing ingress firewall policy from openstack_fw_group_v2 %s: %s", d.Id(), err)
		}
	}

	if group.EgressFirewallPolicyID != "" {
		_, err := groups.RemoveEgressPolicy(networkingClient, group.ID).Extract()
		if err != nil {
			return diag.Errorf("Error removing egress firewall policy from openstack_fw_group_v2 %s: %s", d.Id(), err)
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
		return diag.Errorf("Error waiting for openstack_fw_firewall_v2 %s to Delete:  %s", d.Id(), err)
	}

	return nil
}
