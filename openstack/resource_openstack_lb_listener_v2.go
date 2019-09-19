package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	octavialisteners "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	neutronlisteners "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
)

func resourceListenerV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceListenerV2Create,
		Read:   resourceListenerV2Read,
		Update: resourceListenerV2Update,
		Delete: resourceListenerV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "TCP" && value != "HTTP" && value != "HTTPS" && value != "TERMINATED_HTTPS" {
						errors = append(errors, fmt.Errorf(
							"Only 'TCP', 'HTTP', 'HTTPS' and 'TERMINATED_HTTPS' are supported values for 'protocol'"))
					}
					return
				},
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
		},
	}
}

func resourceListenerV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	lbID := d.Get("loadbalancer_id").(string)
	adminStateUp := d.Get("admin_state_up").(bool)
	var sniContainerRefs []string
	if raw, ok := d.GetOk("sni_container_refs"); ok {
		for _, v := range raw.([]interface{}) {
			sniContainerRefs = append(sniContainerRefs, v.(string))
		}
	}

	timeout := d.Timeout(schema.TimeoutCreate)

	// Wait for LoadBalancer to become active before continuing.
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	// Choose either the Octavia or Neutron create options.
	var createOpts neutronlisteners.CreateOptsBuilder
	if config.useOctavia {
		// Use Octavia.
		opts := octavialisteners.CreateOpts{
			Protocol:               octavialisteners.Protocol(d.Get("protocol").(string)),
			ProtocolPort:           d.Get("protocol_port").(int),
			ProjectID:              d.Get("tenant_id").(string),
			LoadbalancerID:         lbID,
			Name:                   d.Get("name").(string),
			DefaultPoolID:          d.Get("default_pool_id").(string),
			Description:            d.Get("description").(string),
			DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
			SniContainerRefs:       sniContainerRefs,
			AdminStateUp:           &adminStateUp,
		}

		if v, ok := d.GetOk("connection_limit"); ok {
			connectionLimit := v.(int)
			opts.ConnLimit = &connectionLimit
		}

		if v, ok := d.GetOk("timeout_client_data"); ok {
			timeoutClientData := v.(int)
			opts.TimeoutClientData = &timeoutClientData
		}

		if v, ok := d.GetOk("timeout_member_connect"); ok {
			timeoutMemberConnect := v.(int)
			opts.TimeoutMemberConnect = &timeoutMemberConnect
		}

		if v, ok := d.GetOk("timeout_member_data"); ok {
			timeoutMemberData := v.(int)
			opts.TimeoutMemberData = &timeoutMemberData
		}

		if v, ok := d.GetOk("timeout_tcp_inspect"); ok {
			timeoutTCPInspect := v.(int)
			opts.TimeoutTCPInspect = &timeoutTCPInspect
		}

		log.Printf("[DEBUG] Create Options: %#v", opts)

		createOpts = opts
	} else {
		// Use Neutron.
		opts := neutronlisteners.CreateOpts{
			Protocol:               neutronlisteners.Protocol(d.Get("protocol").(string)),
			ProtocolPort:           d.Get("protocol_port").(int),
			TenantID:               d.Get("tenant_id").(string),
			LoadbalancerID:         lbID,
			Name:                   d.Get("name").(string),
			DefaultPoolID:          d.Get("default_pool_id").(string),
			Description:            d.Get("description").(string),
			DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
			SniContainerRefs:       sniContainerRefs,
			AdminStateUp:           &adminStateUp,
		}

		if v, ok := d.GetOk("connection_limit"); ok {
			connectionLimit := v.(int)
			opts.ConnLimit = &connectionLimit
		}

		log.Printf("[DEBUG] Create Options: %#v", opts)

		createOpts = opts
	}

	log.Printf("[DEBUG] Attempting to create listener")
	var listener *neutronlisteners.Listener
	err = resource.Retry(timeout, func() *resource.RetryError {
		listener, err = neutronlisteners.Create(lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating listener: %s", err)
	}

	// Wait for the listener to become ACTIVE.
	err = waitForLBV2Listener(lbClient, listener, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	d.SetId(listener.ID)

	return resourceListenerV2Read(d, meta)
}

func resourceListenerV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listenerResp := neutronlisteners.Get(lbClient, d.Id())

	// Choose either the Octavia or Neutron response type.
	if config.useOctavia {
		var listener octavialisteners.Listener
		err = listenerResp.ExtractInto(listener)
		if err != nil {
			return CheckDeleted(d, err, "listener")
		}

		log.Printf("[DEBUG] Retrieved listener %s: %#v", d.Id(), listener)

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
		d.Set("region", GetRegion(d, config))
	} else {
		var listener neutronlisteners.Listener
		err = listenerResp.ExtractInto(listener)
		if err != nil {
			return CheckDeleted(d, err, "listener")
		}

		log.Printf("[DEBUG] Retrieved listener %s: %#v", d.Id(), listener)

		// Required by import
		if len(listener.Loadbalancers) > 0 {
			d.Set("loadbalancer_id", listener.Loadbalancers[0].ID)
		}

		d.Set("name", listener.Name)
		d.Set("protocol", listener.Protocol)
		d.Set("tenant_id", listener.TenantID)
		d.Set("description", listener.Description)
		d.Set("protocol_port", listener.ProtocolPort)
		d.Set("admin_state_up", listener.AdminStateUp)
		d.Set("default_pool_id", listener.DefaultPoolID)
		d.Set("connection_limit", listener.ConnLimit)
		d.Set("sni_container_refs", listener.SniContainerRefs)
		d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef)
		d.Set("region", GetRegion(d, config))
	}

	return nil
}

func resourceListenerV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := neutronlisteners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve listener %s: %s", d.Id(), err)
	}

	// Wait for the listener to become ACTIVE.
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBV2Listener(lbClient, listener, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	// Choose either the Octavia or Neutron update options.
	var updateOpts neutronlisteners.UpdateOptsBuilder
	if config.useOctavia {
		// Use Octavia.
		var opts octavialisteners.UpdateOpts
		if d.HasChange("name") {
			name := d.Get("name").(string)
			opts.Name = &name
		}
		if d.HasChange("description") {
			description := d.Get("description").(string)
			opts.Description = &description
		}
		if d.HasChange("connection_limit") {
			connLimit := d.Get("connection_limit").(int)
			opts.ConnLimit = &connLimit
		}
		if d.HasChange("timeout_client_data") {
			timeoutClientData := d.Get("timeout_client_data").(int)
			opts.ConnLimit = &timeoutClientData
		}
		if d.HasChange("timeout_member_connect") {
			timeoutMemberConnect := d.Get("timeout_member_connect").(int)
			opts.ConnLimit = &timeoutMemberConnect
		}
		if d.HasChange("timeout_member_data") {
			timeoutMemberData := d.Get("timeout_member_data").(int)
			opts.ConnLimit = &timeoutMemberData
		}
		if d.HasChange("timeout_tcp_inspect") {
			timeoutTCPInspect := d.Get("timeout_tcp_inspect").(int)
			opts.ConnLimit = &timeoutTCPInspect
		}
		if d.HasChange("default_pool_id") {
			defaultPoolID := d.Get("default_pool_id").(string)
			opts.DefaultPoolID = &defaultPoolID
		}
		if d.HasChange("default_tls_container_ref") {
			opts.DefaultTlsContainerRef = d.Get("default_tls_container_ref").(string)
		}
		if d.HasChange("sni_container_refs") {
			var sniContainerRefs []string
			if raw, ok := d.GetOk("sni_container_refs"); ok {
				for _, v := range raw.([]interface{}) {
					sniContainerRefs = append(sniContainerRefs, v.(string))
				}
			}
			opts.SniContainerRefs = sniContainerRefs
		}
		if d.HasChange("admin_state_up") {
			asu := d.Get("admin_state_up").(bool)
			opts.AdminStateUp = &asu
		}

		updateOpts = opts
	} else {
		// Use Neutron.
		var opts neutronlisteners.UpdateOpts
		if d.HasChange("name") {
			name := d.Get("name").(string)
			opts.Name = &name
		}
		if d.HasChange("description") {
			description := d.Get("description").(string)
			opts.Description = &description
		}
		if d.HasChange("connection_limit") {
			connLimit := d.Get("connection_limit").(int)
			opts.ConnLimit = &connLimit
		}
		if d.HasChange("default_pool_id") {
			defaultPoolID := d.Get("default_pool_id").(string)
			opts.DefaultPoolID = &defaultPoolID
		}
		if d.HasChange("default_tls_container_ref") {
			opts.DefaultTlsContainerRef = d.Get("default_tls_container_ref").(string)
		}
		if d.HasChange("sni_container_refs") {
			var sniContainerRefs []string
			if raw, ok := d.GetOk("sni_container_refs"); ok {
				for _, v := range raw.([]interface{}) {
					sniContainerRefs = append(sniContainerRefs, v.(string))
				}
			}
			opts.SniContainerRefs = sniContainerRefs
		}
		if d.HasChange("admin_state_up") {
			asu := d.Get("admin_state_up").(bool)
			opts.AdminStateUp = &asu
		}

		updateOpts = opts
	}

	log.Printf("[DEBUG] Updating listener %s with options: %#v", d.Id(), updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = neutronlisteners.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error updating listener %s: %s", d.Id(), err)
	}

	// Wait for the listener to become ACTIVE.
	err = waitForLBV2Listener(lbClient, listener, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	return resourceListenerV2Read(d, meta)

}

func resourceListenerV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := neutronlisteners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Unable to retrieve listener")
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Deleting listener %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = neutronlisteners.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return CheckDeleted(d, err, "Error deleting listener")
	}

	// Wait for the listener to become DELETED.
	err = waitForLBV2Listener(lbClient, listener, "DELETED", lbPendingDeleteStatuses, timeout)
	if err != nil {
		return err
	}

	return nil
}
