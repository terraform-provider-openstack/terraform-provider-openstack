package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIKEPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIKEPolicyV2Create,
		ReadContext:   resourceIKEPolicyV2Read,
		UpdateContext: resourceIKEPolicyV2Update,
		DeleteContext: resourceIKEPolicyV2Delete,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "sha1",
				ValidateFunc: resourceIKEPolicyV2AuthAlgorithm,
			},
			"encryption_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "aes-128",
				ValidateFunc: resourceIKEPolicyV2EncryptionAlgorithm,
			},
			"pfs": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "group5",
				ValidateFunc: resourceIKEPolicyV2PFS,
			},
			"phase1_negotiation_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "main",
				ValidateFunc: resourceIKEPolicyV2Phase1NegotiationMode,
			},
			"ike_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "v1",
				ValidateFunc: resourceIKEPolicyV2IKEVersion,
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
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIKEPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	lifetime := resourceIKEPolicyV2LifetimeCreateOpts(d.Get("lifetime").(*schema.Set))
	authAlgorithm := ikepolicies.AuthAlgorithm(d.Get("auth_algorithm").(string))
	encryptionAlgorithm := ikepolicies.EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
	pfs := ikepolicies.PFS(d.Get("pfs").(string))
	ikeVersion := ikepolicies.IKEVersion(d.Get("ike_version").(string))
	phase1NegotationMode := ikepolicies.Phase1NegotiationMode(d.Get("phase1_negotiation_mode").(string))

	opts := IKEPolicyCreateOpts{
		ikepolicies.CreateOpts{
			Name:                  d.Get("name").(string),
			Description:           d.Get("description").(string),
			TenantID:              d.Get("tenant_id").(string),
			Lifetime:              &lifetime,
			AuthAlgorithm:         authAlgorithm,
			EncryptionAlgorithm:   encryptionAlgorithm,
			PFS:                   pfs,
			IKEVersion:            ikeVersion,
			Phase1NegotiationMode: phase1NegotationMode,
		},
		MapValueSpecs(d),
	}
	log.Printf("[DEBUG] Create IKE policy: %#v", opts)

	policy, err := ikepolicies.Create(ctx, networkingClient, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForIKEPolicyCreation(ctx, networkingClient, policy.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_vpnaas_ike_policy_v2 %s to become active: %s", policy.ID, err)
	}

	log.Printf("[DEBUG] IKE policy created: %#v", policy)

	d.SetId(policy.ID)

	return resourceIKEPolicyV2Read(ctx, d, meta)
}

func resourceIKEPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about IKE policy: %s", d.Id())

	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	policy, err := ikepolicies.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "IKE policy"))
	}

	log.Printf("[DEBUG] Read OpenStack IKE Policy %s: %#v", d.Id(), policy)

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("auth_algorithm", policy.AuthAlgorithm)
	d.Set("encryption_algorithm", policy.EncryptionAlgorithm)
	d.Set("tenant_id", policy.TenantID)
	d.Set("pfs", policy.PFS)
	d.Set("phase1_negotiation_mode", policy.Phase1NegotiationMode)
	d.Set("ike_version", policy.IKEVersion)
	d.Set("region", GetRegion(d, config))

	// Set the lifetime
	lifetimeMap := make(map[string]any)
	lifetimeMap["units"] = policy.Lifetime.Units
	lifetimeMap["value"] = policy.Lifetime.Value

	var lifetime []map[string]any

	lifetime = append(lifetime, lifetimeMap)
	if err := d.Set("lifetime", &lifetime); err != nil {
		log.Printf("[WARN] unable to set IKE policy lifetime")
	}

	return nil
}

func resourceIKEPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	opts := ikepolicies.UpdateOpts{}

	var hasChange bool

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

	if d.HasChange("pfs") {
		opts.PFS = ikepolicies.PFS(d.Get("pfs").(string))
		hasChange = true
	}

	if d.HasChange("auth_algorithm") {
		opts.AuthAlgorithm = ikepolicies.AuthAlgorithm(d.Get("auth_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("encryption_algorithm") {
		opts.EncryptionAlgorithm = ikepolicies.EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("phase_1_negotiation_mode") {
		opts.Phase1NegotiationMode = ikepolicies.Phase1NegotiationMode(d.Get("phase_1_negotiation_mode").(string))
		hasChange = true
	}

	if d.HasChange("ike_version") {
		opts.IKEVersion = ikepolicies.IKEVersion(d.Get("ike_version").(string))
		hasChange = true
	}

	if d.HasChange("lifetime") {
		lifetime := resourceIKEPolicyV2LifetimeUpdateOpts(d.Get("lifetime").(*schema.Set))
		opts.Lifetime = &lifetime
		hasChange = true
	}

	log.Printf("[DEBUG] Updating IKE policy with id %s: %#v", d.Id(), opts)

	if hasChange {
		err = ikepolicies.Update(ctx, networkingClient, d.Id(), opts).Err
		if err != nil {
			return diag.FromErr(err)
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"ACTIVE"},
			Refresh:    waitForIKEPolicyUpdate(ctx, networkingClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIKEPolicyV2Read(ctx, d, meta)
}

func resourceIKEPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy IKE policy: %s", d.Id())

	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForIKEPolicyDeletion(ctx, networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForIKEPolicyDeletion(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		err := ikepolicies.Delete(ctx, networkingClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		return nil, "ACTIVE", err
	}
}

func waitForIKEPolicyCreation(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		policy, err := ikepolicies.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}

		return policy, "ACTIVE", nil
	}
}

func waitForIKEPolicyUpdate(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		policy, err := ikepolicies.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}

		return policy, "ACTIVE", nil
	}
}

func resourceIKEPolicyV2AuthAlgorithm(v any, k string) ([]string, []error) {
	switch ikepolicies.AuthAlgorithm(v.(string)) {
	case ikepolicies.AuthAlgorithmSHA1,
		ikepolicies.AuthAlgorithmSHA256,
		ikepolicies.AuthAlgorithmSHA384,
		ikepolicies.AuthAlgorithmSHA512,
		ikepolicies.AuthAlgorithmAESXCBC,
		ikepolicies.AuthAlgorithmAESCMAC:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_vpnaas_ike_policy_v2", k, v)}
}

func resourceIKEPolicyV2EncryptionAlgorithm(v any, k string) ([]string, []error) {
	switch ikepolicies.EncryptionAlgorithm(v.(string)) {
	case ikepolicies.EncryptionAlgorithm3DES,
		ikepolicies.EncryptionAlgorithmAES128,
		ikepolicies.EncryptionAlgorithmAES256,
		ikepolicies.EncryptionAlgorithmAES192,
		ikepolicies.EncryptionAlgorithmAES128CTR,
		ikepolicies.EncryptionAlgorithmAES192CTR,
		ikepolicies.EncryptionAlgorithmAES256CTR,
		ikepolicies.EncryptionAlgorithmAES128CCM8,
		ikepolicies.EncryptionAlgorithmAES192CCM8,
		ikepolicies.EncryptionAlgorithmAES256CCM8,
		ikepolicies.EncryptionAlgorithmAES128CCM12,
		ikepolicies.EncryptionAlgorithmAES192CCM12,
		ikepolicies.EncryptionAlgorithmAES256CCM12,
		ikepolicies.EncryptionAlgorithmAES128CCM16,
		ikepolicies.EncryptionAlgorithmAES192CCM16,
		ikepolicies.EncryptionAlgorithmAES256CCM16,
		ikepolicies.EncryptionAlgorithmAES128GCM8,
		ikepolicies.EncryptionAlgorithmAES192GCM8,
		ikepolicies.EncryptionAlgorithmAES256GCM8,
		ikepolicies.EncryptionAlgorithmAES128GCM12,
		ikepolicies.EncryptionAlgorithmAES192GCM12,
		ikepolicies.EncryptionAlgorithmAES256GCM12,
		ikepolicies.EncryptionAlgorithmAES128GCM16,
		ikepolicies.EncryptionAlgorithmAES192GCM16,
		ikepolicies.EncryptionAlgorithmAES256GCM16:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_vpnaas_ike_policy_v2", k, v)}
}

func resourceIKEPolicyV2PFS(v any, k string) ([]string, []error) {
	switch ikepolicies.PFS(v.(string)) {
	case ikepolicies.PFSGroup2,
		ikepolicies.PFSGroup5,
		ikepolicies.PFSGroup14,
		ikepolicies.PFSGroup15,
		ikepolicies.PFSGroup16,
		ikepolicies.PFSGroup17,
		ikepolicies.PFSGroup18,
		ikepolicies.PFSGroup19,
		ikepolicies.PFSGroup20,
		ikepolicies.PFSGroup21,
		ikepolicies.PFSGroup22,
		ikepolicies.PFSGroup23,
		ikepolicies.PFSGroup24,
		ikepolicies.PFSGroup25,
		ikepolicies.PFSGroup26,
		ikepolicies.PFSGroup27,
		ikepolicies.PFSGroup28,
		ikepolicies.PFSGroup29,
		ikepolicies.PFSGroup30,
		ikepolicies.PFSGroup31:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_vpnaas_ike_policy_v2", k, v)}
}

func resourceIKEPolicyV2IKEVersion(v any, k string) ([]string, []error) {
	switch ikepolicies.IKEVersion(v.(string)) {
	case ikepolicies.IKEVersionv1,
		ikepolicies.IKEVersionv2:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_vpnaas_ike_policy_v2", k, v)}
}

func resourceIKEPolicyV2Phase1NegotiationMode(v any, k string) ([]string, []error) {
	switch ikepolicies.Phase1NegotiationMode(v.(string)) {
	case ikepolicies.Phase1NegotiationModeMain:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_vpnaas_ike_policy_v2", k, v)}
}

func resourceIKEPolicyV2Unit(v string) ikepolicies.Unit {
	var unit ikepolicies.Unit

	switch v {
	case "kilobytes":
		unit = ikepolicies.UnitKilobytes
	case "seconds":
		unit = ikepolicies.UnitSeconds
	}

	return unit
}

func resourceIKEPolicyV2LifetimeCreateOpts(d *schema.Set) ikepolicies.LifetimeCreateOpts {
	lifetimeCreateOpts := ikepolicies.LifetimeCreateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]any)
		lifetimeCreateOpts.Units = resourceIKEPolicyV2Unit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetimeCreateOpts.Value = value
	}

	return lifetimeCreateOpts
}

func resourceIKEPolicyV2LifetimeUpdateOpts(d *schema.Set) ikepolicies.LifetimeUpdateOpts {
	lifetimeUpdateOpts := ikepolicies.LifetimeUpdateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]any)
		lifetimeUpdateOpts.Units = resourceIKEPolicyV2Unit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetimeUpdateOpts.Value = value
	}

	return lifetimeUpdateOpts
}
