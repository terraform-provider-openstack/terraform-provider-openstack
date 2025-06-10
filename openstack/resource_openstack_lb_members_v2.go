package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMembersV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMembersV2Create,
		ReadContext:   resourceMembersV2Read,
		UpdateContext: resourceMembersV2Update,
		DeleteContext: resourceMembersV2Delete,
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

			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"member": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"address": {
							Type:     schema.TypeString,
							Required: true,
						},

						"protocol_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(0, 256),
						},

						"monitor_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"monitor_address": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"backup": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"admin_state_up": {
							Type:     schema.TypeBool,
							Default:  true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMembersV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := expandLBMembersV2(d.Get("member").(*schema.Set))
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

	log.Printf("[DEBUG] Attempting to create members")

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = pools.BatchUpdateMembers(ctx, lbClient, poolID, createOpts).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Error creating members: %s", err)
	}

	// Wait for parent pool to become active before continuing
	err = waitForLBV2Pool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(poolID)

	return resourceMembersV2Read(ctx, d, meta)
}

func resourceMembersV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	allPages, err := pools.ListMembers(lbClient, d.Id(), pools.ListMembersOpts{}).AllPages(ctx)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_lb_members_v2"))
	}

	members, err := pools.ExtractMembers(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_lb_members_v2: %s", err)
	}

	log.Printf("[DEBUG] Retrieved members for the %s pool: %#v", d.Id(), members)

	d.Set("pool_id", d.Id())
	d.Set("member", flattenLBMembersV2(members))
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceMembersV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if d.HasChange("member") {
		updateOpts := expandLBMembersV2(d.Get("member").(*schema.Set))

		// Get a clean copy of the parent pool.
		parentPool, err := pools.Get(ctx, lbClient, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("Unable to retrieve parent pool %s: %s", d.Id(), err)
		}

		// Wait for parent pool to become active before continuing.
		timeout := d.Timeout(schema.TimeoutUpdate)

		err = waitForLBV2Pool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updating %s pool members with options: %#v", d.Id(), updateOpts)

		err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			err = pools.BatchUpdateMembers(ctx, lbClient, d.Id(), updateOpts).ExtractErr()
			if err != nil {
				return checkForRetryableError(err)
			}

			return nil
		})
		if err != nil {
			return diag.Errorf("Unable to update member %s: %s", d.Id(), err)
		}

		// Wait for parent pool to become active before continuing
		err = waitForLBV2Pool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMembersV2Read(ctx, d, meta)
}

func resourceMembersV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the parent pool.
	parentPool, err := pools.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, fmt.Sprintf("Unable to retrieve parent pool (%s) for the member", d.Id())))
	}

	// Wait for parent pool to become active before continuing.
	timeout := d.Timeout(schema.TimeoutDelete)

	err = waitForLBV2Pool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error waiting for the members' pool status"))
	}

	log.Printf("[DEBUG] Attempting to delete %s pool members", d.Id())

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = pools.BatchUpdateMembers(ctx, lbClient, d.Id(), []pools.BatchUpdateMemberOpts{}).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting members"))
	}

	// Wait for parent pool to become active before continuing.
	err = waitForLBV2Pool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error waiting for the members' pool status"))
	}

	return nil
}
