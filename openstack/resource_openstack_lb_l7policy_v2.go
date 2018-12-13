package openstack

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"
)

func resourceL7policyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceL7policyV2Create,
		Read:   resourceL7policyV2Read,
		Update: resourceL7policyV2Update,
		Delete: resourceL7policyV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"action": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "REDIRECT_TO_POOL" && value != "REDIRECT_TO_URL" && value != "REJECT" {
						errors = append(errors, fmt.Errorf(
							"Only 'REDIRECT_TO_POOL', 'REDIRECT_TO_URL' and 'REJECT' are supported values for 'action'"))
					}
					return
				},
			},

			"listener_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"position": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"redirect_pool_id": &schema.Schema{
				Type:          schema.TypeString,
				ConflictsWith: []string{"redirect_url"},
				Optional:      true,
			},

			"redirect_url": &schema.Schema{
				Type:          schema.TypeString,
				ConflictsWith: []string{"redirect_pool_id"},
				Optional:      true,
			},

			"admin_state_up": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceL7policyV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := l7policies.CreateOpts{
		TenantID:       d.Get("tenant_id").(string),
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Action:         l7policies.Action(d.Get("action").(string)),
		ListenerID:     d.Get("listener_id").(string),
		RedirectPoolID: d.Get("redirect_pool_id").(string),
		RedirectURL:    d.Get("redirect_url").(string),
		AdminStateUp:   &adminStateUp,
	}

	if v, ok := d.GetOk("position"); ok {
		createOpts.Position = int32(v.(int))
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)
	listenerID := createOpts.ListenerID

	err = checkL7policyAction(lbClient, &createOpts, timeout)
	if err != nil {
		return err
	}

	// Wait for Load Balancer via Listener to become active before continuing
	err = waitForLBV2viaListener(lbClient, listenerID, "ACTIVE", []string{"PENDING_CREATE", "PENDING_UPDATE"}, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Attempting to create l7policy")
	var l7policy *l7policies.L7Policy
	err = resource.Retry(timeout, func() *resource.RetryError {
		l7policy, err = l7policies.Create(lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating l7policy: %s", err)
	}

	// Wait for Load Balancer via Listener to become active before continuing
	err = waitForLBV2viaListener(lbClient, listenerID, "ACTIVE", []string{"PENDING_CREATE", "PENDING_UPDATE"}, timeout)
	if err != nil {
		return err
	}

	d.SetId(l7policy.ID)

	return resourceL7policyV2Read(d, meta)
}

func resourceL7policyV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	l7policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "l7policy")
	}

	log.Printf("[DEBUG] Retrieved l7policy %s: %#v", d.Id(), l7policy)

	// In certain cases LBaaSv2 extension doesn't return "ListenerID", skipping
	// d.Set("listener_id", l7policy.ListenerID)
	d.Set("action", l7policy.Action)
	d.Set("description", l7policy.Description)
	d.Set("tenant_id", l7policy.TenantID)
	d.Set("name", l7policy.Name)
	d.Set("position", int(l7policy.Position))
	d.Set("redirect_url", l7policy.RedirectURL)
	d.Set("redirect_pool_id", l7policy.RedirectPoolID)
	d.Set("region", GetRegion(d, config))
	d.Set("admin_state_up", l7policy.AdminStateUp)

	return nil
}

func resourceL7policyV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var doCheckAction bool
	var updateOpts l7policies.UpdateOpts
	redirectPoolID := ""
	redirectURL := ""
	if d.HasChange("action") {
		doCheckAction = true
		updateOpts.Action = l7policies.Action(d.Get("action").(string))
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
		doCheckAction = true
		redirectPoolID = d.Get("redirect_pool_id").(string)
		updateOpts.RedirectPoolID = &redirectPoolID
	}
	if d.HasChange("redirect_url") {
		doCheckAction = true
		redirectURL = d.Get("redirect_url").(string)
		updateOpts.RedirectURL = &redirectURL
	}
	if d.HasChange("position") {
		updateOpts.Position = d.Get("position").(int32)
	}
	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &adminStateUp
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	listenerID := d.Get("listener_id").(string)

	if doCheckAction {
		err = checkL7policyAction(lbClient, &updateOpts, timeout)
		if err != nil {
			return err
		}
	}

	// Wait for Load Balancer via Listener to become active before continuing
	err = waitForLBV2viaListener(lbClient, listenerID, "ACTIVE", []string{"PENDING_CREATE", "PENDING_UPDATE"}, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating l7policy %s with options: %#v", d.Id(), updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = l7policies.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Unable to update l7policy %s: %s", d.Id(), err)
	}

	// Wait for Load Balancer via Listener to become active before continuing
	err = waitForLBV2viaListener(lbClient, listenerID, "ACTIVE", []string{"PENDING_CREATE", "PENDING_UPDATE"}, timeout)
	if err != nil {
		return err
	}

	return resourceL7policyV2Read(d, meta)
}

func resourceL7policyV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	listenerID := d.Get("listener_id").(string)
	// Wait for Load Balancer via Listener to become active before continuing
	err = waitForLBV2viaListener(lbClient, listenerID, "ACTIVE", []string{"PENDING_CREATE", "PENDING_UPDATE"}, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Attempting to delete l7policy %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = l7policies.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting l7policy %s: %s", d.Id(), err)
	}

	err = waitForLBV2L7policy(lbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return err
	}

	// Wait for Load Balancer via Listener to become active before continuing
	err = waitForLBV2viaListener(lbClient, listenerID, "ACTIVE", []string{"PENDING_CREATE", "PENDING_UPDATE"}, timeout)
	if err != nil {
		return err
	}

	return nil
}

func checkL7policyAction(lbClient *gophercloud.ServiceClient, opts interface{}, timeout time.Duration) (err error) {
	var action l7policies.Action
	var redirectURL *string
	var redirectPoolID *string
	var update *l7policies.UpdateOpts

	switch t := opts.(type) {
	case *l7policies.CreateOpts:
		action = (*opts.(*l7policies.CreateOpts)).Action
		tmp1 := (*opts.(*l7policies.CreateOpts)).RedirectURL
		redirectURL = &tmp1
		tmp2 := (*opts.(*l7policies.CreateOpts)).RedirectPoolID
		redirectPoolID = &tmp2
	case *l7policies.UpdateOpts:
		action = (*opts.(*l7policies.UpdateOpts)).Action
		redirectURL = (*opts.(*l7policies.UpdateOpts)).RedirectURL
		redirectPoolID = (*opts.(*l7policies.UpdateOpts)).RedirectPoolID
		update = opts.(*l7policies.UpdateOpts)
	default:
		return fmt.Errorf(`Invalid type: %s`, t)
	}

	if action == "REJECT" {
		if *redirectURL != "" {
			return fmt.Errorf(`"redirect_url" should be empty, when "action" is set to %s`, action)
		}
		if *redirectPoolID != "" {
			return fmt.Errorf(`"redirect_pool_id" should be empty, when "action" is set to %s`, action)
		}
	} else {
		if action == "REDIRECT_TO_POOL" && *redirectURL == "" {
			if update != nil {
				// We have to unset the value in order to change the "action" type
				update.RedirectURL = nil
			}
			err = waitForLBV2viaPool(lbClient, *redirectPoolID, "ACTIVE", []string{"PENDING_CREATE", "PENDING_UPDATE"}, timeout)
			if err != nil {
				return err
			}
		} else if action == "REDIRECT_TO_POOL" {
			return fmt.Errorf(`"redirect_url" should be empty, when "action" is set to %s`, action)
		}

		if action == "REDIRECT_TO_URL" && *redirectPoolID == "" {
			if update != nil {
				// We have to unset the value in order to change the "action" type
				update.RedirectPoolID = nil
			}
			_, err = url.ParseRequestURI(*redirectURL)
			if err != nil {
				return err
			}
		} else if action == "REDIRECT_TO_URL" {
			return fmt.Errorf(`"redirect_pool_id" should be empty, when "action" is set to %s`, action)
		}
	}
	return nil
}
