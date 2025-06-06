package openstack

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/l7policies"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceL7PolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceL7PolicyV2Create,
		ReadContext:   resourceL7PolicyV2Read,
		UpdateContext: resourceL7PolicyV2Update,
		DeleteContext: resourceL7PolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceL7PolicyV2Import,
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

			"tenant_id": {
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

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"REDIRECT_TO_POOL", "REDIRECT_TO_URL", "REJECT",
					"REDIRECT_PREFIX",
				}, true),
			},

			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"position": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"redirect_prefix": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"redirect_url", "redirect_pool_id"},
				Optional:      true,
			},

			"redirect_pool_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"redirect_url", "redirect_prefix"},
				Optional:      true,
			},

			"redirect_url": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"redirect_pool_id", "redirect_prefix"},
				Optional:      true,
				ValidateFunc: func(v any, _ string) (ws []string, errors []error) {
					value := v.(string)
					_, err := url.ParseRequestURI(value)
					if err != nil {
						errors = append(errors, fmt.Errorf("URL is not valid: %w", err))
					}

					return
				},
			},

			"redirect_http_code": {
				Type:          schema.TypeInt,
				ConflictsWith: []string{"redirect_url", "redirect_pool_id"},
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.IntInSlice([]int{301, 302, 303, 307, 308}),
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceL7PolicyV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Assign some required variables for use in creation.
	listenerID := d.Get("listener_id").(string)
	action := d.Get("action").(string)
	redirectPoolID := d.Get("redirect_pool_id").(string)
	redirectURL := d.Get("redirect_url").(string)
	redirectPrefix := d.Get("redirect_prefix").(string)
	redirectHTTPCodeInt := d.Get("redirect_http_code").(int)
	redirectHTTPCode := int32(redirectHTTPCodeInt)

	// Ensure the right combination of options have been specified.
	err = checkL7PolicyAction(action, redirectURL, redirectPoolID, redirectPrefix)
	if err != nil {
		return diag.Errorf("Unable to create L7 Policy: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := l7policies.CreateOpts{
		ProjectID:        d.Get("tenant_id").(string),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Action:           l7policies.Action(action),
		ListenerID:       listenerID,
		RedirectPoolID:   redirectPoolID,
		RedirectURL:      redirectURL,
		RedirectPrefix:   redirectPrefix,
		RedirectHttpCode: redirectHTTPCode,
		AdminStateUp:     &adminStateUp,
	}

	if v, ok := d.GetOk("position"); ok {
		createOpts.Position = int32(v.(int))
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	// Make sure the associated pool is active before proceeding.
	if redirectPoolID != "" {
		pool, err := pools.Get(ctx, lbClient, redirectPoolID).Extract()
		if err != nil {
			return diag.Errorf("Unable to retrieve %s: %s", redirectPoolID, err)
		}

		err = waitForLBV2Pool(ctx, lbClient, pool, "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(ctx, lbClient, listenerID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve listener %s: %s", listenerID, err)
	}

	// Wait for parent Listener to become active before continuing.
	err = waitForLBV2Listener(ctx, lbClient, parentListener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to create L7 Policy")

	var l7Policy *l7policies.L7Policy

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		l7Policy, err = l7policies.Create(ctx, lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Error creating L7 Policy: %s", err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, lbClient, parentListener, l7Policy, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l7Policy.ID)

	return resourceL7PolicyV2Read(ctx, d, meta)
}

func resourceL7PolicyV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	l7Policy, err := l7policies.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "L7 Policy"))
	}

	log.Printf("[DEBUG] Retrieved L7 Policy %s: %#v", d.Id(), l7Policy)

	d.Set("action", l7Policy.Action)
	d.Set("description", l7Policy.Description)
	d.Set("tenant_id", l7Policy.ProjectID)
	d.Set("name", l7Policy.Name)
	d.Set("position", int(l7Policy.Position))
	d.Set("redirect_url", l7Policy.RedirectURL)
	d.Set("redirect_pool_id", l7Policy.RedirectPoolID)
	d.Set("redirect_prefix", l7Policy.RedirectPrefix)
	d.Set("redirect_http_code", l7Policy.RedirectHttpCode)
	d.Set("region", GetRegion(d, config))
	d.Set("admin_state_up", l7Policy.AdminStateUp)

	return nil
}

func resourceL7PolicyV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Assign some required variables for use in updating.
	listenerID := d.Get("listener_id").(string)
	action := d.Get("action").(string)
	redirectPoolID := d.Get("redirect_pool_id").(string)
	redirectURL := d.Get("redirect_url").(string)
	redirectPrefix := d.Get("redirect_prefix").(string)
	redirectHTTPCodeInt := d.Get("redirect_http_code").(int)
	redirectHTTPCode := int32(redirectHTTPCodeInt)

	var updateOpts l7policies.UpdateOpts

	if d.HasChange("action") {
		updateOpts.Action = l7policies.Action(action)
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("redirect_pool_id") {
		updateOpts.RedirectPoolID = &redirectPoolID
	}

	if d.HasChange("redirect_url") {
		updateOpts.RedirectURL = &redirectURL
	}

	if d.HasChange("redirect_prefix") {
		updateOpts.RedirectPrefix = &redirectPrefix
	}

	if d.HasChange("redirect_http_code") {
		updateOpts.RedirectHttpCode = redirectHTTPCode
	}

	if d.HasChange("position") {
		updateOpts.Position = int32(d.Get("position").(int))
	}

	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &adminStateUp
	}

	// Ensure the right combination of options have been specified.
	err = checkL7PolicyAction(action, redirectURL, redirectPoolID, redirectPrefix)
	if err != nil {
		return diag.FromErr(err)
	}

	// Make sure the pool is active before continuing.
	timeout := d.Timeout(schema.TimeoutUpdate)

	if redirectPoolID != "" {
		pool, err := pools.Get(ctx, lbClient, redirectPoolID).Extract()
		if err != nil {
			return diag.Errorf("Unable to retrieve %s: %s", redirectPoolID, err)
		}

		err = waitForLBV2Pool(ctx, lbClient, pool, "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(ctx, lbClient, listenerID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent listener %s: %s", listenerID, err)
	}

	// Get a clean copy of the L7 Policy.
	l7Policy, err := l7policies.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve L7 Policy: %s: %s", d.Id(), err)
	}

	// Wait for parent Listener to become active before continuing.
	err = waitForLBV2Listener(ctx, lbClient, parentListener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, lbClient, parentListener, l7Policy, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating L7 Policy %s with options: %#v", d.Id(), updateOpts)

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err = l7policies.Update(ctx, lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Unable to update L7 Policy %s: %s", d.Id(), err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, lbClient, parentListener, l7Policy, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceL7PolicyV2Read(ctx, d, meta)
}

func resourceL7PolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	listenerID := d.Get("listener_id").(string)

	// Get a clean copy of the listener.
	listener, err := listeners.Get(ctx, lbClient, listenerID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent listener (%s) for the L7 Policy: %s", listenerID, err)
	}

	// Get a clean copy of the L7 Policy.
	l7Policy, err := l7policies.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Unable to retrieve L7 Policy"))
	}

	// Wait for Listener to become active before continuing.
	err = waitForLBV2Listener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to delete L7 Policy %s", d.Id())

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = l7policies.Delete(ctx, lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting L7 Policy"))
	}

	err = waitForLBV2L7Policy(ctx, lbClient, listener, l7Policy, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceL7PolicyV2Import(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack networking client: %w", err)
	}

	l7Policy, err := l7policies.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return nil, CheckDeleted(d, err, "L7 Policy")
	}

	log.Printf("[DEBUG] Retrieved L7 Policy %s during the import: %#v", d.Id(), l7Policy)

	if l7Policy.ListenerID != "" {
		d.Set("listener_id", l7Policy.ListenerID)
	} else {
		// Fallback for the Neutron LBaaSv2 extension
		listenerID, err := getListenerIDForL7Policy(ctx, lbClient, d.Id())
		if err != nil {
			return nil, err
		}

		d.Set("listener_id", listenerID)
	}

	return []*schema.ResourceData{d}, nil
}

func checkL7PolicyAction(action, redirectURL, redirectPoolID, redirectPrefix string) error {
	if action == "REJECT" {
		if redirectURL != "" || redirectPoolID != "" || redirectPrefix != "" {
			return fmt.Errorf(
				"redirect_url/pool_id/prefix must be empty when action is set to %s", action)
		}
	}

	if action == "REDIRECT_TO_POOL" && (redirectURL != "" || redirectPrefix != "") {
		return fmt.Errorf("redirect_url/prefix must be empty when action is set to %s", action)
	}

	if action == "REDIRECT_TO_URL" && (redirectPoolID != "" || redirectPrefix != "") {
		return fmt.Errorf("redirect_pool_id/prefix must be empty when action is set to %s", action)
	}

	if action == "REDIRECT_TO_PREFIX" && (redirectPoolID != "" || redirectURL != "") {
		return fmt.Errorf("redirect_pool_id/url must be empty when action is set to %s", action)
	}

	return nil
}
