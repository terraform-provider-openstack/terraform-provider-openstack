package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/rbacpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkingRBACPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRBACPolicyV2Create,
		ReadContext:   resourceNetworkingRBACPolicyV2Read,
		UpdateContext: resourceNetworkingRBACPolicyV2Update,
		DeleteContext: resourceNetworkingRBACPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"access_as_external", "access_as_shared",
				}, false),
			},

			"object_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"address_scope", "address_group", "network", "qos_policy", "security_group", "subnetpool", "bgpvpn",
				}, false),
			},

			"target_tenant": {
				Type:     schema.TypeString,
				Required: true,
			},

			"object_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingRBACPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := rbacpolicies.CreateOpts{
		Action:       rbacpolicies.PolicyAction(d.Get("action").(string)),
		ObjectType:   d.Get("object_type").(string),
		TargetTenant: d.Get("target_tenant").(string),
		ObjectID:     d.Get("object_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	rbac, err := rbacpolicies.Create(ctx, networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_rbac_policy_v2: %s", err)
	}

	d.SetId(rbac.ID)

	return resourceNetworkingRBACPolicyV2Read(ctx, d, meta)
}

func resourceNetworkingRBACPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	rbac, err := rbacpolicies.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_networking_rbac_policy_v2"))
	}

	log.Printf("[DEBUG] Retrieved RBAC policy %s: %+v", d.Id(), rbac)

	d.Set("action", string(rbac.Action))
	d.Set("object_type", rbac.ObjectType)
	d.Set("target_tenant", rbac.TargetTenant)
	d.Set("object_id", rbac.ObjectID)
	d.Set("project_id", rbac.ProjectID)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingRBACPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateOpts rbacpolicies.UpdateOpts

	if d.HasChange("target_tenant") {
		updateOpts.TargetTenant = d.Get("target_tenant").(string)

		_, err := rbacpolicies.Update(ctx, networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_rbac_policy_v2: %s", err)
		}
	}

	return resourceNetworkingRBACPolicyV2Read(ctx, d, meta)
}

func resourceNetworkingRBACPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	err = rbacpolicies.Delete(ctx, networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_rbac_policy_v2"))
	}

	return nil
}
