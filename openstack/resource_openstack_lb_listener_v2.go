package openstack

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

			"client_authentication": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NONE", "OPTIONAL", "MANDATORY",
				}, false),
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					return o == "NONE" && n == ""
				},
				DiffSuppressOnRefresh: true,
			},

			"client_ca_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"client_crl_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"hsts_include_subdomains": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"hsts_max_age"},
			},

			"hsts_max_age": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"hsts_preload": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"hsts_max_age"},
			},

			"tls_ciphers": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true, // unsetting this parameter results in a default value
			},

			"tls_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true, // unsetting this parameter is not possible due to a bug in Octavia
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

func resourceListenerV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)

	// Wait for LoadBalancer to become active before continuing.
	err = waitForLBV2LoadBalancer(ctx, lbClient, d.Get("loadbalancer_id").(string), "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := listeners.CreateOpts{
		// Protocol SCTP requires octavia minor version 2.23
		Protocol:                listeners.Protocol(d.Get("protocol").(string)),
		ProtocolPort:            d.Get("protocol_port").(int),
		ProjectID:               d.Get("tenant_id").(string),
		LoadbalancerID:          d.Get("loadbalancer_id").(string),
		Name:                    d.Get("name").(string),
		DefaultPoolID:           d.Get("default_pool_id").(string),
		Description:             d.Get("description").(string),
		DefaultTlsContainerRef:  d.Get("default_tls_container_ref").(string),
		SniContainerRefs:        expandToStringSlice(d.Get("sni_container_refs").([]any)),
		ALPNProtocols:           expandToStringSlice(d.Get("alpn_protocols").(*schema.Set).List()),
		ClientAuthentication:    listeners.ClientAuthentication(d.Get("client_authentication").(string)),
		ClientCATLSContainerRef: d.Get("client_ca_tls_container_ref").(string),
		ClientCRLContainerRef:   d.Get("client_crl_container_ref").(string),
		HSTSIncludeSubdomains:   d.Get("hsts_include_subdomains").(bool),
		HSTSMaxAge:              d.Get("hsts_max_age").(int),
		HSTSPreload:             d.Get("hsts_preload").(bool),
		TLSCiphers:              d.Get("tls_ciphers").(string),
		InsertHeaders:           expandToMapStringString(d.Get("insert_headers").(map[string]any)),
		AllowedCIDRs:            expandToStringSlice(d.Get("allowed_cidrs").([]any)),
		AdminStateUp:            &adminStateUp,
		Tags:                    expandToStringSlice(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("tls_versions"); ok {
		createOpts.TLSVersions = expandLBListenerTLSVersionV2(v.(*schema.Set).List())
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

	log.Printf("[DEBUG] openstack_lb_listener_v2 create options: %#v", createOpts)

	var listener *listeners.Listener

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		listener, err = listeners.Create(ctx, lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Error creating openstack_lb_listener_v2: %s", err)
	}

	// Wait for the listener to become ACTIVE.
	err = waitForLBV2Listener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(listener.ID)

	return resourceListenerV2Read(ctx, d, meta)
}

func resourceListenerV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listener, err := listeners.Get(ctx, lbClient, d.Id()).Extract()
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
	d.Set("alpn_protocols", listener.ALPNProtocols)
	d.Set("client_authentication", listener.ClientAuthentication)
	d.Set("client_ca_tls_container_ref", listener.ClientCATLSContainerRef)
	d.Set("client_crl_container_ref", listener.ClientCRLContainerRef)
	d.Set("hsts_include_subdomains", listener.HSTSIncludeSubdomains)
	d.Set("hsts_max_age", listener.HSTSMaxAge)
	d.Set("hsts_preload", listener.HSTSPreload)
	d.Set("tls_ciphers", listener.TLSCiphers)
	d.Set("tls_versions", listener.TLSVersions)
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

func resourceListenerV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := listeners.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_lb_listener_v2 %s: %s", d.Id(), err)
	}

	// Wait for the listener to become ACTIVE.
	timeout := d.Timeout(schema.TimeoutUpdate)

	err = waitForLBV2Listener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
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
		v := expandToStringSlice(d.Get("sni_container_refs").([]any))
		updateOpts.SniContainerRefs = &v
	}

	if d.HasChange("admin_state_up") {
		hasChange = true
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	if d.HasChange("insert_headers") {
		hasChange = true
		v := expandToMapStringString(d.Get("insert_headers").(map[string]any))
		updateOpts.InsertHeaders = &v
	}

	if d.HasChange("allowed_cidrs") {
		hasChange = true
		v := expandToStringSlice(d.Get("allowed_cidrs").([]any))
		updateOpts.AllowedCIDRs = &v
	}

	if d.HasChange("alpn_protocols") {
		hasChange = true
		v := expandToStringSlice(d.Get("alpn_protocols").(*schema.Set).List())
		updateOpts.ALPNProtocols = &v
	}

	if d.HasChange("client_authentication") {
		hasChange = true

		v := listeners.ClientAuthentication(d.Get("client_authentication").(string))
		if v == "" {
			v = listeners.ClientAuthenticationNone
		}

		updateOpts.ClientAuthentication = &v
	}

	if d.HasChange("client_ca_tls_container_ref") {
		hasChange = true
		v := d.Get("client_ca_tls_container_ref").(string)
		updateOpts.ClientCATLSContainerRef = &v
	}

	if d.HasChange("client_crl_container_ref") {
		hasChange = true
		v := d.Get("client_crl_container_ref").(string)
		updateOpts.ClientCRLContainerRef = &v
	}

	if d.HasChange("hsts_include_subdomains") {
		hasChange = true
		v := d.Get("hsts_include_subdomains").(bool)
		updateOpts.HSTSIncludeSubdomains = &v
	}

	if d.HasChange("hsts_max_age") {
		hasChange = true
		v := d.Get("hsts_max_age").(int)
		updateOpts.HSTSMaxAge = &v
	}

	if d.HasChange("hsts_preload") {
		hasChange = true
		v := d.Get("hsts_preload").(bool)
		updateOpts.HSTSPreload = &v
	}

	if d.HasChange("tls_ciphers") {
		hasChange = true
		v := d.Get("tls_ciphers").(string)
		updateOpts.TLSCiphers = &v
	}

	if d.HasChange("tls_versions") {
		hasChange = true
		v := expandLBListenerTLSVersionV2(d.Get("tls_versions").(*schema.Set).List())
		updateOpts.TLSVersions = &v
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

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err = listeners.Update(ctx, lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("Error updating openstack_lb_listener_v2 %s: %s", d.Id(), err)
	}

	// Wait for the listener to become ACTIVE.
	err = waitForLBV2Listener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceListenerV2Read(ctx, d, meta)
}

func resourceListenerV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := listeners.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Unable to retrieve openstack_lb_listener_v2"))
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Deleting openstack_lb_listener_v2 %s", d.Id())

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = listeners.Delete(ctx, lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_lb_listener_v2"))
	}

	// Wait for the listener to become DELETED.
	err = waitForLBV2Listener(ctx, lbClient, listener, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
