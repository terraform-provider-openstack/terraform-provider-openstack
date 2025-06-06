package openstack

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePoolV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePoolV2Create,
		ReadContext:   resourcePoolV2Read,
		UpdateContext: resourcePoolV2Update,
		DeleteContext: resourcePoolV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePoolV2Import,
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

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP", "HTTPS", "PROXY", "SCTP", "PROXYV2",
				}, false),
			},

			// One of loadbalancer_id or listener_id must be provided
			"loadbalancer_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"loadbalancer_id", "listener_id"},
			},

			// One of loadbalancer_id or listener_id must be provided
			"listener_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"loadbalancer_id", "listener_id"},
			},

			"lb_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP", "SOURCE_IP_PORT",
				}, false),
			},

			"persistence": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SOURCE_IP", "HTTP_COOKIE", "APP_COOKIE",
							}, false),
						},

						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"alpn_protocols": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true, // unsetting this parameter results in a default value
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"http/1.0", "http/1.1", "h2",
					}, false),
				},
			},

			"ca_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"crl_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tls_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tls_ciphers": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true, // unsetting this parameter results in a default value
			},

			"tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tls_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true, // unsetting this parameter results in a default value
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"TLSv1", "TLSv1.1", "TLSv1.2", "TLSv1.3",
					}, false),
				},
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
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

func resourcePoolV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	lbID := d.Get("loadbalancer_id").(string)
	listenerID := d.Get("listener_id").(string)

	createOpts := pools.CreateOpts{
		ProjectID:         d.Get("tenant_id").(string),
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Protocol:          pools.Protocol(d.Get("protocol").(string)),
		LoadbalancerID:    lbID,
		ListenerID:        listenerID,
		LBMethod:          pools.LBMethod(d.Get("lb_method").(string)),
		ALPNProtocols:     expandToStringSlice(d.Get("alpn_protocols").(*schema.Set).List()),
		CATLSContainerRef: d.Get("ca_tls_container_ref").(string),
		CRLContainerRef:   d.Get("crl_container_ref").(string),
		TLSEnabled:        d.Get("tls_enabled").(bool),
		TLSCiphers:        d.Get("tls_ciphers").(string),
		TLSContainerRef:   d.Get("tls_container_ref").(string),
		AdminStateUp:      &adminStateUp,
	}

	if v, ok := d.GetOk("tls_versions"); ok {
		createOpts.TLSVersions = expandLBPoolTLSVersionV2(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("persistence"); ok {
		createOpts.Persistence, err = expandLBPoolPersistanceV2(v.([]any))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(*schema.Set).List()
		createOpts.Tags = expandToStringSlice(tags)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	// Wait for Listener or LoadBalancer to become active before continuing
	if listenerID != "" {
		listener, err := listeners.Get(ctx, lbClient, listenerID).Extract()
		if err != nil {
			return diag.Errorf("Unable to get openstack_lb_listener_v2 %s: %s", listenerID, err)
		}

		waitErr := waitForLBV2Listener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
		if waitErr != nil {
			return diag.Errorf("Error waiting for openstack_lb_listener_v2 %s to become active: %s", listenerID, waitErr)
		}
	} else {
		waitErr := waitForLBV2LoadBalancer(ctx, lbClient, lbID, "ACTIVE", getLbPendingStatuses(), timeout)
		if waitErr != nil {
			return diag.Errorf("Error waiting for openstack_lb_loadbalancer_v2 %s to become active: %s", lbID, waitErr)
		}
	}

	log.Printf("[DEBUG] Attempting to create pool")

	var pool *pools.Pool

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		pool, err = pools.Create(ctx, lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Error creating pool: %s", err)
	}

	// Pool was successfully created
	// Wait for pool to become active before continuing
	err = waitForLBV2Pool(ctx, lbClient, pool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pool.ID)

	return resourcePoolV2Read(ctx, d, meta)
}

func resourcePoolV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	pool, err := pools.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "pool"))
	}

	log.Printf("[DEBUG] Retrieved pool %s: %#v", d.Id(), pool)

	d.Set("lb_method", pool.LBMethod)
	d.Set("protocol", pool.Protocol)
	d.Set("description", pool.Description)
	d.Set("tenant_id", pool.ProjectID)
	d.Set("admin_state_up", pool.AdminStateUp)
	d.Set("name", pool.Name)
	d.Set("persistence", flattenLBPoolPersistenceV2(pool.Persistence))
	d.Set("alpn_protocols", pool.ALPNProtocols)
	d.Set("ca_tls_container_ref", pool.CATLSContainerRef)
	d.Set("crl_container_ref", pool.CRLContainerRef)
	d.Set("tls_enabled", pool.TLSEnabled)
	d.Set("tls_ciphers", pool.TLSCiphers)
	d.Set("tls_container_ref", pool.TLSContainerRef)
	d.Set("tls_versions", pool.TLSVersions)
	d.Set("region", GetRegion(d, config))
	d.Set("tags", pool.Tags)

	return nil
}

func resourcePoolV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	var updateOpts pools.UpdateOpts
	if d.HasChange("lb_method") {
		updateOpts.LBMethod = pools.LBMethod(d.Get("lb_method").(string))
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	if d.HasChange("persistence") {
		updateOpts.Persistence, err = expandLBPoolPersistanceV2(d.Get("persistence").([]any))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("alpn_protocols") {
		v := expandToStringSlice(d.Get("alpn_protocols").(*schema.Set).List())
		updateOpts.ALPNProtocols = &v
	}

	if d.HasChange("ca_tls_container_ref") {
		v := d.Get("ca_tls_container_ref").(string)
		updateOpts.CATLSContainerRef = &v
	}

	if d.HasChange("crl_container_ref") {
		v := d.Get("crl_container_ref").(string)
		updateOpts.CRLContainerRef = &v
	}

	if d.HasChange("tls_enabled") {
		v := d.Get("tls_enabled").(bool)
		updateOpts.TLSEnabled = &v
	}

	if d.HasChange("tls_ciphers") {
		v := d.Get("tls_ciphers").(string)
		updateOpts.TLSCiphers = &v
	}

	if d.HasChange("tls_container_ref") {
		v := d.Get("tls_container_ref").(string)
		updateOpts.TLSContainerRef = &v
	}

	if d.HasChange("tls_versions") {
		v := expandLBPoolTLSVersionV2(d.Get("tls_versions").(*schema.Set).List())
		updateOpts.TLSVersions = &v
	}

	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			tags := v.(*schema.Set).List()
			tagsToUpdate := expandToStringSlice(tags)
			updateOpts.Tags = &tagsToUpdate
		} else {
			updateOpts.Tags = &[]string{}
		}
	}

	timeout := d.Timeout(schema.TimeoutUpdate)

	// Get a clean copy of the pool.
	pool, err := pools.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve pool %s: %s", d.Id(), err)
	}

	// Wait for pool to become active before continuing
	err = waitForLBV2Pool(ctx, lbClient, pool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating pool %s with options: %#v", d.Id(), updateOpts)

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err = pools.Update(ctx, lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Unable to update pool %s: %s", d.Id(), err)
	}

	// Wait for pool to become active before continuing
	err = waitForLBV2Pool(ctx, lbClient, pool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePoolV2Read(ctx, d, meta)
}

func resourcePoolV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	// Get a clean copy of the pool.
	pool, err := pools.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Unable to retrieve pool"))
	}

	log.Printf("[DEBUG] Attempting to delete pool %s", d.Id())

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = pools.Delete(ctx, lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting pool"))
	}

	// Wait for Pool to delete
	err = waitForLBV2Pool(ctx, lbClient, pool, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePoolV2Import(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack loadbalancing client: %w", err)
	}

	pool, err := pools.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return nil, CheckDeleted(d, err, "pool")
	}

	log.Printf("[DEBUG] Retrieved pool %s during the import: %#v", d.Id(), pool)

	if len(pool.Listeners) > 0 && pool.Listeners[0].ID != "" {
		d.Set("listener_id", pool.Listeners[0].ID)
	} else if len(pool.Loadbalancers) > 0 && pool.Loadbalancers[0].ID != "" {
		d.Set("loadbalancer_id", pool.Loadbalancers[0].ID)
	} else {
		return nil, errors.New("Unable to detect pool's Listener ID or Load Balancer ID")
	}

	return []*schema.ResourceData{d}, nil
}
