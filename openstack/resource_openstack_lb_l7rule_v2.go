package openstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
)

func resourceL7RuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceL7RuleV2Create,
		Read:   resourceL7RuleV2Read,
		Update: resourceL7RuleV2Update,
		Delete: resourceL7RuleV2Delete,
		Importer: &schema.ResourceImporter{
			resourceL7RuleV2Import,
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

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"COOKIE", "FILE_TYPE", "HEADER", "HOST_NAME", "PATH",
				}, true),
			},

			"compare_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CONTAINS", "STARTS_WITH", "ENDS_WITH", "EQUAL_TO", "REGEX",
				}, true),
			},

			"l7policy_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"listener_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"invert": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"admin_state_up": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceL7RuleV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Assign some required variables for use in creation.
	l7policyID := d.Get("l7policy_id").(string)
	listenerID := ""
	ruleType := d.Get("type").(string)
	compareType := d.Get("compare_type").(string)

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := l7policies.CreateRuleOpts{
		TenantID:     d.Get("tenant_id").(string),
		RuleType:     l7policies.RuleType(ruleType),
		CompareType:  l7policies.CompareType(compareType),
		Value:        d.Get("value").(string),
		Key:          d.Get("key").(string),
		Invert:       d.Get("invert").(bool),
		AdminStateUp: &adminStateUp,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return fmt.Errorf("Unable to get parent L7 Policy: %s", err)
	}

	if parentL7Policy.ListenerID != "" {
		listenerID = parentL7Policy.ListenerID
	} else {
		// Fallback for the Neutron LBaaSv2 extension
		listenerID, err = getListenerIDForL7Policy(lbClient, l7policyID)
		if err != nil {
			return err
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve listener %s: %s", listenerID, err)
	}

	// Wait for parent L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(lbClient, parentListener, parentL7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Attempting to create L7 Rule")
	var l7Rule *l7policies.Rule
	err = resource.Retry(timeout, func() *resource.RetryError {
		l7Rule, err = l7policies.CreateRule(lbClient, l7policyID, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating L7 Rule: %s", err)
	}

	// Wait for L7 Rule to become active before continuing
	err = waitForLBV2L7Rule(lbClient, parentListener, parentL7Policy, l7Rule, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	d.SetId(l7Rule.ID)
	d.Set("listener_id", listenerID)

	return resourceL7RuleV2Read(d, meta)
}

func resourceL7RuleV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	l7policyID := d.Get("l7policy_id").(string)

	l7Rule, err := l7policies.GetRule(lbClient, l7policyID, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "L7 Rule")
	}

	log.Printf("[DEBUG] Retrieved L7 Rule %s: %#v", d.Id(), l7Rule)

	d.Set("l7policy_id", l7policyID)
	d.Set("type", l7Rule.RuleType)
	d.Set("compare_type", l7Rule.CompareType)
	d.Set("tenant_id", l7Rule.TenantID)
	d.Set("value", l7Rule.Value)
	d.Set("key", l7Rule.Key)
	d.Set("invert", l7Rule.Invert)
	d.Set("admin_state_up", l7Rule.AdminStateUp)

	return nil
}

func resourceL7RuleV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Assign some required variables for use in updating.
	l7policyID := d.Get("l7policy_id").(string)
	listenerID := d.Get("listener_id").(string)

	var updateOpts l7policies.UpdateRuleOpts

	if d.HasChange("type") {
		updateOpts.RuleType = l7policies.RuleType(d.Get("type").(string))
	}
	if d.HasChange("compare_type") {
		updateOpts.CompareType = l7policies.CompareType(d.Get("compare_type").(string))
	}
	if d.HasChange("value") {
		updateOpts.Value = d.Get("value").(string)
	}
	if d.HasChange("key") {
		key := d.Get("key").(string)
		updateOpts.Key = &key
	}
	if d.HasChange("invert") {
		invert := d.Get("invert").(bool)
		updateOpts.Invert = &invert
	}

	timeout := d.Timeout(schema.TimeoutUpdate)

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve listener %s: %s", listenerID, err)
	}

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return fmt.Errorf("Unable to get parent L7 Policy: %s", err)
	}

	// Get a clean copy of the L7 Rule.
	l7Rule, err := l7policies.GetRule(lbClient, l7policyID, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Unable to get L7 Rule: %s", err)
	}

	// Wait for parent L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(lbClient, parentListener, parentL7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	// Wait for L7 Rule to become active before continuing
	err = waitForLBV2L7Rule(lbClient, parentListener, parentL7Policy, l7Rule, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating L7 Rule %s with options: %#v", d.Id(), updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err := l7policies.UpdateRule(lbClient, l7policyID, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Unable to update L7 Rule %s: %s", d.Id(), err)
	}

	// Wait for L7 Rule to become active before continuing
	err = waitForLBV2L7Rule(lbClient, parentListener, parentL7Policy, l7Rule, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	return resourceL7RuleV2Read(d, meta)
}

func resourceL7RuleV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	l7policyID := d.Get("l7policy_id").(string)
	listenerID := d.Get("listener_id").(string)

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve listener %s: %s", listenerID, err)
	}

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return fmt.Errorf("Unable to get parent L7 Policy: %s", err)
	}

	// Get a clean copy of the L7 Rule.
	l7Rule, err := l7policies.GetRule(lbClient, l7policyID, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Unable to get L7 Rule: %s", err)
	}

	// Wait for parent L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(lbClient, parentListener, parentL7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Attempting to delete L7 Rule %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = l7policies.DeleteRule(lbClient, l7policyID, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting L7 Rule %s: %s", d.Id(), err)
	}

	err = waitForLBV2L7Rule(lbClient, parentListener, parentL7Policy, l7Rule, "DELETED", lbPendingDeleteStatuses, timeout)
	if err != nil {
		return err
	}

	return nil
}

func resourceL7RuleV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for L7 Rule. Format must be <policy id>/<rule id>")
		return nil, err
	}

	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listenerID := ""
	l7policyID := parts[0]
	l7ruleID := parts[1]

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return nil, fmt.Errorf("Unable to get parent L7 Policy: %s", err)
	}

	if parentL7Policy.ListenerID != "" {
		listenerID = parentL7Policy.ListenerID
	} else {
		// Fallback for the Neutron LBaaSv2 extension
		listenerID, err = getListenerIDForL7Policy(lbClient, l7policyID)
		if err != nil {
			return nil, err
		}
	}

	d.SetId(l7ruleID)
	d.Set("l7policy_id", l7policyID)
	d.Set("listener_id", listenerID)

	return []*schema.ResourceData{d}, nil
}
