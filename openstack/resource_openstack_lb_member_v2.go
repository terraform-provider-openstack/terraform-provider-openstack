package openstack

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMemberV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberV2Create,
		ReadContext:   resourceMemberV2Read,
		UpdateContext: resourceMemberV2Update,
		DeleteContext: resourceMemberV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMemberV2Import,
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

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"protocol_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},

			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 256),
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"backup": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"monitor_address": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				ForceNew: false,
			},

			"monitor_port": {
				Type:         schema.TypeInt,
				Default:      nil,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntBetween(1, 65535),
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceMemberV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)

	createOpts := pools.CreateMemberOpts{
		Name:         d.Get("name").(string),
		ProjectID:    d.Get("tenant_id").(string),
		Address:      d.Get("address").(string),
		ProtocolPort: d.Get("protocol_port").(int),
		AdminStateUp: &adminStateUp,
	}

	// Must omit if not set
	if v, ok := d.GetOk("subnet_id"); ok {
		createOpts.SubnetID = v.(string)
	}

	// Set the weight only if it's defined in the configuration.
	// This prevents all members from being created with a default weight of 0.
	if v, ok := getOkExists(d, "weight"); ok {
		weight := v.(int)
		createOpts.Weight = &weight
	}

	if v, ok := d.GetOk("monitor_address"); ok {
		createOpts.MonitorAddress = v.(string)
	}

	if v, ok := d.GetOk("monitor_port"); ok {
		monitorPort := v.(int)
		createOpts.MonitorPort = &monitorPort
	}

	// Only set backup if it is defined by user as it requires
	// version 2.1 or later
	if v, ok := d.GetOk("backup"); ok {
		backup := v.(bool)
		createOpts.Backup = &backup
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(*schema.Set).List()
		createOpts.Tags = expandToStringSlice(tags)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := pools.Get(ctx, lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent pool %s: %s", poolID, err)
	}

	// Wait for parent pool to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)

	err = waitForLBV2Pool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to create member")

	var member *pools.Member

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		member, err = pools.CreateMember(ctx, lbClient, poolID, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Error creating member: %s", err)
	}

	// Wait for member to become active before continuing
	err = waitForLBV2Member(ctx, lbClient, parentPool, member, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(member.ID)

	return resourceMemberV2Read(ctx, d, meta)
}

func resourceMemberV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	poolID := d.Get("pool_id").(string)

	member, err := pools.GetMember(ctx, lbClient, poolID, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "member"))
	}

	log.Printf("[DEBUG] Retrieved member %s: %#v", d.Id(), member)

	d.Set("name", member.Name)
	d.Set("weight", member.Weight)
	d.Set("admin_state_up", member.AdminStateUp)
	d.Set("tenant_id", member.ProjectID)
	d.Set("subnet_id", member.SubnetID)
	d.Set("address", member.Address)
	d.Set("protocol_port", member.ProtocolPort)
	d.Set("region", GetRegion(d, config))
	d.Set("monitor_address", member.MonitorAddress)
	d.Set("monitor_port", member.MonitorPort)
	d.Set("backup", member.Backup)
	d.Set("tags", member.Tags)

	return nil
}

func resourceMemberV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateOpts pools.UpdateMemberOpts

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("weight") {
		weight := d.Get("weight").(int)
		updateOpts.Weight = &weight
	}

	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	if d.HasChange("monitor_address") {
		monitorAddress := d.Get("monitor_address").(string)
		updateOpts.MonitorAddress = &monitorAddress
	}

	if d.HasChange("monitor_port") {
		monitorPort := d.Get("monitor_port").(int)
		updateOpts.MonitorPort = &monitorPort
	}

	if d.HasChange("backup") {
		backup := d.Get("backup").(bool)
		updateOpts.Backup = &backup
	}

	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			tags := v.(*schema.Set).List()
			tagsToUpdate := expandToStringSlice(tags)
			updateOpts.Tags = tagsToUpdate
		} else {
			updateOpts.Tags = []string{}
		}
	}

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := pools.Get(ctx, lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent pool %s: %s", poolID, err)
	}

	// Get a clean copy of the member.
	member, err := pools.GetMember(ctx, lbClient, poolID, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve member: %s: %s", d.Id(), err)
	}

	// Wait for parent pool to become active before continuing.
	timeout := d.Timeout(schema.TimeoutUpdate)

	err = waitForLBV2Pool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for the member to become active before continuing.
	err = waitForLBV2Member(ctx, lbClient, parentPool, member, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating member %s with options: %#v", d.Id(), updateOpts)

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err = pools.UpdateMember(ctx, lbClient, poolID, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Unable to update member %s: %s", d.Id(), err)
	}

	// Wait for the member to become active before continuing.
	err = waitForLBV2Member(ctx, lbClient, parentPool, member, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMemberV2Read(ctx, d, meta)
}

func resourceMemberV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := pools.Get(ctx, lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent pool (%s) for the member: %s", poolID, err)
	}

	// Get a clean copy of the member.
	member, err := pools.GetMember(ctx, lbClient, poolID, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Unable to retrieve member"))
	}

	// Wait for parent pool to become active before continuing.
	timeout := d.Timeout(schema.TimeoutDelete)

	err = waitForLBV2Pool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error waiting for the members pool status"))
	}

	log.Printf("[DEBUG] Attempting to delete member %s", d.Id())

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = pools.DeleteMember(ctx, lbClient, poolID, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting member"))
	}

	// Wait for the member to become DELETED.
	err = waitForLBV2Member(ctx, lbClient, parentPool, member, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMemberV2Import(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := errors.New("Invalid format specified for Member. Format must be <pool id>/<member id>")

		return nil, err
	}

	poolID := parts[0]
	memberID := parts[1]

	d.SetId(memberID)
	d.Set("pool_id", poolID)

	return []*schema.ResourceData{d}, nil
}
