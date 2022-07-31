package openstack

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	octaviapools "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	neutronpools "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/pools"
)

func resourceMemberV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberV2Create,
		ReadContext:   resourceMemberV2Read,
		UpdateContext: resourceMemberV2Update,
		DeleteContext: resourceMemberV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceMemberV2Import,
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
		},
	}
}

func resourceMemberV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)

	if config.UseOctavia {
		createOpts := octaviapools.CreateMemberOpts{
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
		if v, ok := d.GetOkExists("weight"); ok {
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

		log.Printf("[DEBUG] Create Options: %#v", createOpts)

		// Get a clean copy of the parent pool.
		poolID := d.Get("pool_id").(string)
		parentPool, err := neutronpools.Get(lbClient, poolID).Extract()
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
		var member *octaviapools.Member
		err = resource.Retry(timeout, func() *resource.RetryError {
			member, err = octaviapools.CreateMember(lbClient, poolID, createOpts).Extract()

			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return diag.Errorf("Error creating member: %s", err)
		}

		// Wait for member to become active before continuing
		err = waitForLBV2OctaviaMember(ctx, lbClient, parentPool, member, "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(member.ID)

		return resourceMemberV2Read(ctx, d, meta)
	}

	createOpts := neutronpools.CreateMemberOpts{
		Name:         d.Get("name").(string),
		TenantID:     d.Get("tenant_id").(string),
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
	if v, ok := d.GetOkExists("weight"); ok {
		weight := v.(int)
		createOpts.Weight = &weight
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := neutronpools.Get(lbClient, poolID).Extract()
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
	var member *neutronpools.Member
	err = resource.Retry(timeout, func() *resource.RetryError {
		member, err = neutronpools.CreateMember(lbClient, poolID, createOpts).Extract()

		// neutronpools.Create(lbClient, poolID, createOpts).Extract()
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

func resourceMemberV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	poolID := d.Get("pool_id").(string)

	if config.UseOctavia {
		member, err := octaviapools.GetMember(lbClient, poolID, d.Id()).Extract()
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

		return nil
	}

	member, err := neutronpools.GetMember(lbClient, poolID, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "member"))
	}

	log.Printf("[DEBUG] Retrieved member %s: %#v", d.Id(), member)

	d.Set("name", member.Name)
	d.Set("weight", member.Weight)
	d.Set("admin_state_up", member.AdminStateUp)
	d.Set("tenant_id", member.TenantID)
	d.Set("subnet_id", member.SubnetID)
	d.Set("address", member.Address)
	d.Set("protocol_port", member.ProtocolPort)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceMemberV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if config.UseOctavia {
		var updateOpts octaviapools.UpdateMemberOpts
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

		// Get a clean copy of the parent pool.
		poolID := d.Get("pool_id").(string)
		parentPool, err := neutronpools.Get(lbClient, poolID).Extract()
		if err != nil {
			return diag.Errorf("Unable to retrieve parent pool %s: %s", poolID, err)
		}

		// Get a clean copy of the member.
		member, err := octaviapools.GetMember(lbClient, poolID, d.Id()).Extract()
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
		err = waitForLBV2OctaviaMember(ctx, lbClient, parentPool, member, "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updating member %s with options: %#v", d.Id(), updateOpts)
		err = resource.Retry(timeout, func() *resource.RetryError {
			_, err = octaviapools.UpdateMember(lbClient, poolID, d.Id(), updateOpts).Extract()
			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return diag.Errorf("Unable to update member %s: %s", d.Id(), err)
		}

		// Wait for the member to become active before continuing.
		err = waitForLBV2OctaviaMember(ctx, lbClient, parentPool, member, "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		return resourceMemberV2Read(ctx, d, meta)
	}

	var updateOpts neutronpools.UpdateMemberOpts
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

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := neutronpools.Get(lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent pool %s: %s", poolID, err)
	}

	// Get a clean copy of the member.
	member, err := neutronpools.GetMember(lbClient, poolID, d.Id()).Extract()
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
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = neutronpools.UpdateMember(lbClient, poolID, d.Id(), updateOpts).Extract()
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

func resourceMemberV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if config.UseOctavia {
		// Get a clean copy of the parent pool.
		poolID := d.Get("pool_id").(string)
		parentPool, err := neutronpools.Get(lbClient, poolID).Extract()
		if err != nil {
			return diag.Errorf("Unable to retrieve parent pool (%s) for the member: %s", poolID, err)
		}

		// Get a clean copy of the member.
		member, err := octaviapools.GetMember(lbClient, poolID, d.Id()).Extract()
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
		err = resource.Retry(timeout, func() *resource.RetryError {
			err = neutronpools.DeleteMember(lbClient, poolID, d.Id()).ExtractErr()
			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return diag.FromErr(CheckDeleted(d, err, "Error deleting member"))
		}

		// Wait for the member to become DELETED.
		err = waitForLBV2OctaviaMember(ctx, lbClient, parentPool, member, "DELETED", getLbPendingDeleteStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		return nil
	}

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := neutronpools.Get(lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent pool (%s) for the member: %s", poolID, err)
	}

	// Get a clean copy of the member.
	member, err := neutronpools.GetMember(lbClient, poolID, d.Id()).Extract()
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
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = neutronpools.DeleteMember(lbClient, poolID, d.Id()).ExtractErr()
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

func resourceMemberV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for Member. Format must be <pool id>/<member id>")
		return nil, err
	}

	poolID := parts[0]
	memberID := parts[1]

	d.SetId(memberID)
	d.Set("pool_id", poolID)

	return []*schema.ResourceData{d}, nil
}
