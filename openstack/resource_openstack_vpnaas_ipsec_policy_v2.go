package openstack

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIPSecPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPSecPolicyV2Create,
		ReadContext:   resourceIPSecPolicyV2Read,
		UpdateContext: resourceIPSecPolicyV2Update,
		DeleteContext: resourceIPSecPolicyV2Delete,
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: resourceIPSecPolicyV2AuthAlgorithm,
			},
			"encapsulation_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: resourceIPSecPolicyV2EncapsulationMode,
			},
			"pfs": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: resourceIPSecPolicyV2PFS,
			},
			"encryption_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: resourceIPSecPolicyV2EncryptionAlgorithm,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transform_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: resourceIPSecPolicyV2TransformProtocol,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"lifetime": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"units": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIPSecPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	encapsulationMode := ipsecpolicies.EncapsulationMode(d.Get("encapsulation_mode").(string))
	authAlgorithm := ipsecpolicies.AuthAlgorithm(d.Get("auth_algorithm").(string))
	encryptionAlgorithm := ipsecpolicies.EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
	pfs := ipsecpolicies.PFS(d.Get("pfs").(string))
	transformProtocol := ipsecpolicies.TransformProtocol(d.Get("transform_protocol").(string))
	lifetime := resourceIPSecPolicyV2LifetimeCreateOpts(d.Get("lifetime").(*schema.Set))

	opts := IPSecPolicyCreateOpts{
		ipsecpolicies.CreateOpts{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			TenantID:            d.Get("tenant_id").(string),
			EncapsulationMode:   encapsulationMode,
			AuthAlgorithm:       authAlgorithm,
			EncryptionAlgorithm: encryptionAlgorithm,
			PFS:                 pfs,
			TransformProtocol:   transformProtocol,
			Lifetime:            &lifetime,
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create IPSec policy: %#v", opts)

	policy, err := ipsecpolicies.Create(ctx, networkingClient, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForIPSecPolicyCreation(ctx, networkingClient, policy.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_vpnaas_ipsec_policy_v2 %s to become active: %s", policy.ID, err)
	}

	log.Printf("[DEBUG] IPSec policy created: %#v", policy)

	d.SetId(policy.ID)

	return resourceIPSecPolicyV2Read(ctx, d, meta)
}

func resourceIPSecPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about IPSec policy: %s", d.Id())

	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	policy, err := ipsecpolicies.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "IPSec policy"))
	}

	log.Printf("[DEBUG] Read OpenStack IPSec policy %s: %#v", d.Id(), policy)

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("tenant_id", policy.TenantID)
	d.Set("encapsulation_mode", policy.EncapsulationMode)
	d.Set("encryption_algorithm", policy.EncryptionAlgorithm)
	d.Set("transform_protocol", policy.TransformProtocol)
	d.Set("pfs", policy.PFS)
	d.Set("auth_algorithm", policy.AuthAlgorithm)
	d.Set("region", GetRegion(d, config))

	// Set the lifetime
	lifetimeMap := make(map[string]any)
	lifetimeMap["units"] = policy.Lifetime.Units
	lifetimeMap["value"] = policy.Lifetime.Value

	var lifetime []map[string]any

	lifetime = append(lifetime, lifetimeMap)
	if err := d.Set("lifetime", &lifetime); err != nil {
		log.Printf("[WARN] unable to set IPSec policy lifetime")
	}

	return nil
}

func resourceIPSecPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var hasChange bool

	opts := ipsecpolicies.UpdateOpts{}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
		hasChange = true
	}

	if d.HasChange("auth_algorithm") {
		opts.AuthAlgorithm = ipsecpolicies.AuthAlgorithm(d.Get("auth_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("encryption_algorithm") {
		opts.EncryptionAlgorithm = ipsecpolicies.EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("transform_protocol") {
		opts.TransformProtocol = ipsecpolicies.TransformProtocol(d.Get("transform_protocol").(string))
		hasChange = true
	}

	if d.HasChange("pfs") {
		opts.PFS = ipsecpolicies.PFS(d.Get("pfs").(string))
		hasChange = true
	}

	if d.HasChange("encapsulation_mode") {
		opts.EncapsulationMode = ipsecpolicies.EncapsulationMode(d.Get("encapsulation_mode").(string))
		hasChange = true
	}

	if d.HasChange("lifetime") {
		lifetime := resourceIPSecPolicyV2LifetimeUpdateOpts(d.Get("lifetime").(*schema.Set))
		opts.Lifetime = &lifetime
		hasChange = true
	}

	log.Printf("[DEBUG] Updating IPSec policy with id %s: %#v", d.Id(), opts)

	if hasChange {
		_, err = ipsecpolicies.Update(ctx, networkingClient, d.Id(), opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"ACTIVE"},
			Refresh:    waitForIPSecPolicyUpdate(ctx, networkingClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIPSecPolicyV2Read(ctx, d, meta)
}

func resourceIPSecPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy IPSec policy: %s", d.Id())

	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForIPSecPolicyDeletion(ctx, networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForIPSecPolicyDeletion(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		err := ipsecpolicies.Delete(ctx, networkingClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
			return nil, "ACTIVE", nil
		}

		return nil, "ACTIVE", err
	}
}

func waitForIPSecPolicyCreation(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		policy, err := ipsecpolicies.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}

		return policy, "ACTIVE", nil
	}
}

func waitForIPSecPolicyUpdate(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		policy, err := ipsecpolicies.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}

		return policy, "ACTIVE", nil
	}
}

func resourceIPSecPolicyV2TransformProtocol(v any, k string) ([]string, []error) {
	switch ipsecpolicies.TransformProtocol(v.(string)) {
	case ipsecpolicies.TransformProtocolESP,
		ipsecpolicies.TransformProtocolAH,
		ipsecpolicies.TransformProtocolAHESP:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_vpnaas_ipsec_policy_v2", k, v)}
}

func resourceIPSecPolicyV2PFS(v any, k string) ([]string, []error) {
	switch ipsecpolicies.PFS(v.(string)) {
	case ipsecpolicies.PFSGroup2,
		ipsecpolicies.PFSGroup5,
		ipsecpolicies.PFSGroup14,
		ipsecpolicies.PFSGroup15,
		ipsecpolicies.PFSGroup16,
		ipsecpolicies.PFSGroup17,
		ipsecpolicies.PFSGroup18,
		ipsecpolicies.PFSGroup19,
		ipsecpolicies.PFSGroup20,
		ipsecpolicies.PFSGroup21,
		ipsecpolicies.PFSGroup22,
		ipsecpolicies.PFSGroup23,
		ipsecpolicies.PFSGroup24,
		ipsecpolicies.PFSGroup25,
		ipsecpolicies.PFSGroup26,
		ipsecpolicies.PFSGroup27,
		ipsecpolicies.PFSGroup28,
		ipsecpolicies.PFSGroup29,
		ipsecpolicies.PFSGroup30,
		ipsecpolicies.PFSGroup31:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_vpnaas_ipsec_policy_v2", k, v)}
}

func resourceIPSecPolicyV2EncryptionAlgorithm(v any, k string) ([]string, []error) {
	switch ipsecpolicies.EncryptionAlgorithm(v.(string)) {
	case ipsecpolicies.EncryptionAlgorithm3DES,
		ipsecpolicies.EncryptionAlgorithmAES128,
		ipsecpolicies.EncryptionAlgorithmAES256,
		ipsecpolicies.EncryptionAlgorithmAES192,
		ipsecpolicies.EncryptionAlgorithmAES128CTR,
		ipsecpolicies.EncryptionAlgorithmAES192CTR,
		ipsecpolicies.EncryptionAlgorithmAES256CTR,
		ipsecpolicies.EncryptionAlgorithmAES128CCM8,
		ipsecpolicies.EncryptionAlgorithmAES192CCM8,
		ipsecpolicies.EncryptionAlgorithmAES256CCM8,
		ipsecpolicies.EncryptionAlgorithmAES128CCM12,
		ipsecpolicies.EncryptionAlgorithmAES192CCM12,
		ipsecpolicies.EncryptionAlgorithmAES256CCM12,
		ipsecpolicies.EncryptionAlgorithmAES128CCM16,
		ipsecpolicies.EncryptionAlgorithmAES192CCM16,
		ipsecpolicies.EncryptionAlgorithmAES256CCM16,
		ipsecpolicies.EncryptionAlgorithmAES128GCM8,
		ipsecpolicies.EncryptionAlgorithmAES192GCM8,
		ipsecpolicies.EncryptionAlgorithmAES256GCM8,
		ipsecpolicies.EncryptionAlgorithmAES128GCM12,
		ipsecpolicies.EncryptionAlgorithmAES192GCM12,
		ipsecpolicies.EncryptionAlgorithmAES256GCM12,
		ipsecpolicies.EncryptionAlgorithmAES128GCM16,
		ipsecpolicies.EncryptionAlgorithmAES192GCM16,
		ipsecpolicies.EncryptionAlgorithmAES256GCM16:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_vpnaas_ipsec_policy_v2", k, v)}
}

func resourceIPSecPolicyV2AuthAlgorithm(v any, k string) ([]string, []error) {
	switch ipsecpolicies.AuthAlgorithm(v.(string)) {
	case ipsecpolicies.AuthAlgorithmSHA1,
		ipsecpolicies.AuthAlgorithmSHA256,
		ipsecpolicies.AuthAlgorithmSHA384,
		ipsecpolicies.AuthAlgorithmSHA512,
		ipsecpolicies.AuthAlgorithmAESXCBC,
		ipsecpolicies.AuthAlgorithmAESCMAC:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_vpnaas_ipsec_policy_v2", k, v)}
}

func resourceIPSecPolicyV2EncapsulationMode(v any, k string) ([]string, []error) {
	switch ipsecpolicies.EncapsulationMode(v.(string)) {
	case ipsecpolicies.EncapsulationModeTunnel,
		ipsecpolicies.EncapsulationModeTransport:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_vpnaas_ipsec_policy_v2", k, v)}
}

func resourceIPSecPolicyV2LifetimeCreateOpts(d *schema.Set) ipsecpolicies.LifetimeCreateOpts {
	lifetime := ipsecpolicies.LifetimeCreateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]any)
		lifetime.Units = resourceIPSecPolicyV2Unit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetime.Value = value
	}

	return lifetime
}

func resourceIPSecPolicyV2Unit(units string) ipsecpolicies.Unit {
	var unit ipsecpolicies.Unit

	switch units {
	case "seconds":
		unit = ipsecpolicies.UnitSeconds
	case "kilobytes":
		unit = ipsecpolicies.UnitKilobytes
	}

	return unit
}

func resourceIPSecPolicyV2LifetimeUpdateOpts(d *schema.Set) ipsecpolicies.LifetimeUpdateOpts {
	lifetimeUpdateOpts := ipsecpolicies.LifetimeUpdateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]any)
		lifetimeUpdateOpts.Units = resourceIPSecPolicyV2Unit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetimeUpdateOpts.Value = value
	}

	return lifetimeUpdateOpts
}
