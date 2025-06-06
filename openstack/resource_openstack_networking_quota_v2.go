package openstack

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/quotas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkingQuotaV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingQuotaV2Create,
		ReadContext:   resourceNetworkingQuotaV2Read,
		UpdateContext: resourceNetworkingQuotaV2Update,
		Delete:        schema.RemoveFromState,
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

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"floatingip": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"network": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"rbac_policy": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"router": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"security_group": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"security_group_rule": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"subnet": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"subnetpool": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingQuotaV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)

	networkingClient, err := config.NetworkingV2Client(ctx, region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	updateOpts := quotas.UpdateOpts{}
	projectID := d.Get("project_id").(string)

	if v, ok := getOkExists(d, "floatingip"); ok {
		pfloatingIP := v.(int)
		updateOpts.FloatingIP = &pfloatingIP
	}

	if v, ok := getOkExists(d, "network"); ok {
		pnetwork := v.(int)
		updateOpts.Network = &pnetwork
	}

	if v, ok := getOkExists(d, "port"); ok {
		pport := v.(int)
		updateOpts.Port = &pport
	}

	if v, ok := getOkExists(d, "rbac_policy"); ok {
		prbacPolicy := v.(int)
		updateOpts.RBACPolicy = &prbacPolicy
	}

	if v, ok := getOkExists(d, "router"); ok {
		prouter := v.(int)
		updateOpts.Router = &prouter
	}

	if v, ok := getOkExists(d, "security_group"); ok {
		psecurityGroup := v.(int)
		updateOpts.SecurityGroup = &psecurityGroup
	}

	if v, ok := getOkExists(d, "security_group_rule"); ok {
		psecurityGroupRule := v.(int)
		updateOpts.SecurityGroupRule = &psecurityGroupRule
	}

	if v, ok := getOkExists(d, "subnet"); ok {
		psubnet := v.(int)
		updateOpts.Subnet = &psubnet
	}

	if v, ok := getOkExists(d, "subnetpool"); ok {
		psubnetPool := v.(int)
		updateOpts.SubnetPool = &psubnetPool
	}

	q, err := quotas.Update(ctx, networkingClient, projectID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_quota_v2: %s", err)
	}

	id := fmt.Sprintf("%s/%s", projectID, region)
	d.SetId(id)

	log.Printf("[DEBUG] Created openstack_networking_quota_v2 %#v", q)

	return resourceNetworkingQuotaV2Read(ctx, d, meta)
}

func resourceNetworkingQuotaV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)

	networkingClient, err := config.NetworkingV2Client(ctx, region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Depending on the provider version the resource was created, the resource id
	// can be either <project_id> or <project_id>/<region>. This parses the project_id
	// in both cases
	projectID := strings.Split(d.Id(), "/")[0]

	q, err := quotas.Get(ctx, networkingClient, projectID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_networking_quota_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_quota_v2 %s: %#v", d.Id(), q)

	d.Set("project_id", projectID)
	d.Set("region", region)
	d.Set("floatingip", q.FloatingIP)
	d.Set("network", q.Network)
	d.Set("port", q.Port)
	d.Set("rbac_policy", q.RBACPolicy)
	d.Set("router", q.Router)
	d.Set("security_group", q.SecurityGroup)
	d.Set("security_group_rule", q.SecurityGroupRule)
	d.Set("subnet", q.Subnet)
	d.Set("subnetpool", q.SubnetPool)

	return nil
}

func resourceNetworkingQuotaV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		hasChange  bool
		updateOpts quotas.UpdateOpts
	)

	if d.HasChange("floatingip") {
		hasChange = true
		floatingIP := d.Get("floatingip").(int)
		updateOpts.FloatingIP = &floatingIP
	}

	if d.HasChange("network") {
		hasChange = true
		network := d.Get("network").(int)
		updateOpts.Network = &network
	}

	if d.HasChange("port") {
		hasChange = true
		port := d.Get("port").(int)
		updateOpts.Port = &port
	}

	if d.HasChange("rbac_policy") {
		hasChange = true
		rbacPolicy := d.Get("rbac_policy").(int)
		updateOpts.RBACPolicy = &rbacPolicy
	}

	if d.HasChange("router") {
		hasChange = true
		router := d.Get("router").(int)
		updateOpts.Router = &router
	}

	if d.HasChange("security_group") {
		hasChange = true
		securityGroup := d.Get("security_group").(int)
		updateOpts.SecurityGroup = &securityGroup
	}

	if d.HasChange("security_group_rule") {
		hasChange = true
		securityGroupRule := d.Get("security_group_rule").(int)
		updateOpts.SecurityGroupRule = &securityGroupRule
	}

	if d.HasChange("subnet") {
		hasChange = true
		subnet := d.Get("subnet").(int)
		updateOpts.Subnet = &subnet
	}

	if d.HasChange("subnetpool") {
		hasChange = true
		subnetPool := d.Get("subnetpool").(int)
		updateOpts.SubnetPool = &subnetPool
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_networking_quota_v2 %s update options: %#v", d.Id(), updateOpts)
		projectID := d.Get("project_id").(string)

		_, err := quotas.Update(ctx, networkingClient, projectID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_quota_v2: %s", err)
		}
	}

	return resourceNetworkingQuotaV2Read(ctx, d, meta)
}
