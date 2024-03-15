package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
)

func resourceListenerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListenerV2Create,
		ReadContext:   resourceListenerV2Read,
		UpdateContext: resourceListenerV2Update,
		DeleteContext: resourceListenerV2Delete,
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

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "SCTP", "HTTP", "HTTPS", "TERMINATED_HTTPS", "PROMETHEUS",
				}, false),
			},

			"protocol_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"connection_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"sni_container_refs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"timeout_client_data": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"timeout_member_connect": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"timeout_member_data": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"timeout_tcp_inspect": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"insert_headers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},

			"allowed_cidrs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceListenerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := config.LoadBalancerV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)

	// Wait for LoadBalancer to become active before continuing.
	err = waitForLBV2LoadBalancerOctavia(ctx, lbClient, d.Get("loadbalancer_id").(string), "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	var sniContainerRefs []string
	if raw, ok := d.GetOk("sni_container_refs"); ok {
		for _, v := range raw.([]interface{}) {
			sniContainerRefs = append(sniContainerRefs, v.(string))
		}
	}

	createOpts := listeners.CreateOpts{
		// Protocol SCTP requires octavia minor version 2.23
		Protocol:               listeners.Protocol(d.Get("protocol").(string)),
		ProtocolPort:           d.Get("protocol_port").(int),
		ProjectID:              d.Get("tenant_id").(string),
		LoadbalancerID:         d.Get("loadbalancer_id").(string),
		Name:                   d.Get("name").(string),
		DefaultPoolID:          d.Get("default_pool_id").(string),
		Description:            d.Get("description").(string),
		DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
		SniContainerRefs:       sniContainerRefs,
		AdminStateUp:           &adminStateUp,
	}

	if v, ok := d.GetOk("connection_limit"); ok {
		connectionLimit := v.(int)
		createOpts.ConnLimit = &connectionLimit
	}

	if v, ok := d.GetOk("timeout_client_data"); ok {
		timeoutClientData := v.(int)
		createOpts.TimeoutClientData = &timeoutClientData
	}

	if v, ok := d.GetOk("timeout_member_connect"); ok {
		timeoutMemberConnect := v.(int)
		createOpts.TimeoutMemberConnect = &timeoutMemberConnect
	}

	if v, ok := d.GetOk("timeout_member_data"); ok {
		timeoutMemberData := v.(int)
		createOpts.TimeoutMemberData = &timeoutMemberData
	}

	if v, ok := d.GetOk("timeout_tcp_inspect"); ok {
		timeoutTCPInspect := v.(int)
		createOpts.TimeoutTCPInspect = &timeoutTCPInspect
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(*schema.Set).List()
		createOpts.Tags = expandToStringSlice(tags)
	}

	// Get and check insert  headers map.
	rawHeaders := d.Get("insert_headers").(map[string]interface{})
	headers, err := expandLBV2ListenerHeadersMap(rawHeaders)
	if err != nil {
		return diag.Errorf("Unable to parse insert_headers argument for openstack_lb_listener_v2: %s", err)
	}

	createOpts.InsertHeaders = headers

	if raw, ok := d.GetOk("allowed_cidrs"); ok {
		allowedCidrs := make([]string, len(raw.([]interface{})))
		for i, v := range raw.([]interface{}) {
			allowedCidrs[i] = v.(string)
		}
		createOpts.AllowedCIDRs = allowedCidrs
	}

	log.Printf("[DEBUG] openstack_lb_listener_v2 create options: %#v", createOpts)
	var listener *listeners.Listener
	err = resource.Retry(timeout, func() *resource.RetryError {
		listener, err = listeners.Create(lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Error creating openstack_lb_listener_v2: %s", err)
	}

	// Wait for the listener to become ACTIVE.
	err = waitForLBV2ListenerOctavia(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(listener.ID)

	return resourceListenerV2Read(ctx, d, meta)
}

func resourceListenerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := config.LoadBalancerV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listener, err := listeners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "openstack_lb_listener_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_lb_listener_v2 %s: %#v", d.Id(), listener)

	d.Set("name", listener.Name)
	d.Set("protocol", listener.Protocol)
	d.Set("tenant_id", listener.ProjectID)
	d.Set("description", listener.Description)
	d.Set("protocol_port", listener.ProtocolPort)
	d.Set("admin_state_up", listener.AdminStateUp)
	d.Set("default_pool_id", listener.DefaultPoolID)
	d.Set("connection_limit", listener.ConnLimit)
	d.Set("timeout_client_data", listener.TimeoutClientData)
	d.Set("timeout_member_connect", listener.TimeoutMemberConnect)
	d.Set("timeout_member_data", listener.TimeoutMemberData)
	d.Set("timeout_tcp_inspect", listener.TimeoutTCPInspect)
	d.Set("sni_container_refs", listener.SniContainerRefs)
	d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef)
	d.Set("allowed_cidrs", listener.AllowedCIDRs)
	d.Set("region", GetRegion(d, config))
	d.Set("tags", listener.Tags)

	// Required by import.
	if len(listener.Loadbalancers) > 0 {
		d.Set("loadbalancer_id", listener.Loadbalancers[0].ID)
	}

	if err := d.Set("insert_headers", listener.InsertHeaders); err != nil {
		return diag.Errorf("Unable to set openstack_lb_listener_v2 insert_headers: %s", err)
	}

	return nil
}

func resourceListenerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := config.LoadBalancerV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := listeners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_lb_listener_v2 %s: %s", d.Id(), err)
	}

	// Wait for the listener to become ACTIVE.
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBV2ListenerOctavia(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	var updateOpts listeners.UpdateOpts
	var hasChange bool

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

	if d.HasChange("connection_limit") {
		hasChange = true
		connLimit := d.Get("connection_limit").(int)
		updateOpts.ConnLimit = &connLimit
	}

	if d.HasChange("timeout_client_data") {
		hasChange = true
		timeoutClientData := d.Get("timeout_client_data").(int)
		updateOpts.TimeoutClientData = &timeoutClientData
	}

	if d.HasChange("timeout_member_connect") {
		hasChange = true
		timeoutMemberConnect := d.Get("timeout_member_connect").(int)
		updateOpts.TimeoutMemberConnect = &timeoutMemberConnect
	}

	if d.HasChange("timeout_member_data") {
		hasChange = true
		timeoutMemberData := d.Get("timeout_member_data").(int)
		updateOpts.TimeoutMemberData = &timeoutMemberData
	}

	if d.HasChange("timeout_tcp_inspect") {
		hasChange = true
		timeoutTCPInspect := d.Get("timeout_tcp_inspect").(int)
		updateOpts.TimeoutTCPInspect = &timeoutTCPInspect
	}

	if d.HasChange("default_pool_id") {
		hasChange = true
		defaultPoolID := d.Get("default_pool_id").(string)
		updateOpts.DefaultPoolID = &defaultPoolID
	}

	if d.HasChange("default_tls_container_ref") {
		hasChange = true
		defaultTLSContainerRef := d.Get("default_tls_container_ref").(string)
		updateOpts.DefaultTlsContainerRef = &defaultTLSContainerRef
	}

	if d.HasChange("sni_container_refs") {
		hasChange = true
		var sniContainerRefs []string
		if raw, ok := d.GetOk("sni_container_refs"); ok {
			for _, v := range raw.([]interface{}) {
				sniContainerRefs = append(sniContainerRefs, v.(string))
			}
		}
		updateOpts.SniContainerRefs = &sniContainerRefs
	}

	if d.HasChange("admin_state_up") {
		hasChange = true
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	if d.HasChange("insert_headers") {
		hasChange = true

		// Get and check insert headers map.
		rawHeaders := d.Get("insert_headers").(map[string]interface{})
		headers, err := expandLBV2ListenerHeadersMap(rawHeaders)
		if err != nil {
			return diag.Errorf("Error parsing insert header for openstack_lb_listener_v2 %s: %s", d.Id(), err)
		}

		updateOpts.InsertHeaders = &headers
	}

	if d.HasChange("allowed_cidrs") {
		hasChange = true
		var allowedCidrs []string
		if raw, ok := d.GetOk("allowed_cidrs"); ok {
			for _, v := range raw.([]interface{}) {
				allowedCidrs = append(allowedCidrs, v.(string))
			}
		}
		updateOpts.AllowedCIDRs = &allowedCidrs
	}

	if d.HasChange("tags") {
		hasChange = true
		if v, ok := d.GetOk("tags"); ok {
			tags := v.(*schema.Set).List()
			tagsToUpdate := expandToStringSlice(tags)
			updateOpts.Tags = &tagsToUpdate
		} else {
			updateOpts.Tags = &[]string{}
		}
	}

	if !hasChange {
		log.Printf("[DEBUG] openstack_lb_listener_v2 %s: nothing to update", d.Id())
		return resourceListenerV2Read(ctx, d, meta)
	}

	log.Printf("[DEBUG] openstack_lb_listener_v2 %s update options: %#v", d.Id(), updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = listeners.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Error updating openstack_lb_listener_v2 %s: %s", d.Id(), err)
	}

	// Wait for the listener to become ACTIVE.
	err = waitForLBV2ListenerOctavia(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceListenerV2Read(ctx, d, meta)
}

func resourceListenerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := config.LoadBalancerV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := listeners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Unable to retrieve openstack_lb_listener_v2"))
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Deleting openstack_lb_listener_v2 %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = listeners.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_lb_listener_v2"))
	}

	// Wait for the listener to become DELETED.
	err = waitForLBV2ListenerOctavia(ctx, lbClient, listener, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
