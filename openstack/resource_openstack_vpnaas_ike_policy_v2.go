package openstack

import (
	"context"
	"log"
	"maps"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/ikepolicies"
)

var ikePolicyV2PFSMap = map[string]ikepolicies.PFS{
	"group2":  ikepolicies.PFSGroup2,
	"group5":  ikepolicies.PFSGroup5,
	"group14": ikepolicies.PFSGroup14,
	"group15": ikepolicies.PFSGroup15,
	"group16": ikepolicies.PFSGroup16,
	"group17": ikepolicies.PFSGroup17,
	"group18": ikepolicies.PFSGroup18,
	"group19": ikepolicies.PFSGroup19,
	"group20": ikepolicies.PFSGroup20,
	"group21": ikepolicies.PFSGroup21,
	"group22": ikepolicies.PFSGroup22,
	"group23": ikepolicies.PFSGroup23,
	"group24": ikepolicies.PFSGroup24,
	"group25": ikepolicies.PFSGroup25,
	"group26": ikepolicies.PFSGroup26,
	"group27": ikepolicies.PFSGroup27,
	"group28": ikepolicies.PFSGroup28,
	"group29": ikepolicies.PFSGroup29,
	"group30": ikepolicies.PFSGroup30,
	"group31": ikepolicies.PFSGroup31,
}

var ikePolicyV2EncryptionAlgorithmMap = map[string]ikepolicies.EncryptionAlgorithm{
	"3des":           ikepolicies.EncryptionAlgorithm3DES,
	"aes-128":        ikepolicies.EncryptionAlgorithmAES128,
	"aes-256":        ikepolicies.EncryptionAlgorithmAES256,
	"aes-192":        ikepolicies.EncryptionAlgorithmAES192,
	"aes-128-ctr":    ikepolicies.EncryptionAlgorithmAES128CTR,
	"aes-192-ctr":    ikepolicies.EncryptionAlgorithmAES192CTR,
	"aes-256-ctr":    ikepolicies.EncryptionAlgorithmAES256CTR,
	"aes-128-ccm-8":  ikepolicies.EncryptionAlgorithmAES128CCM8,
	"aes-192-ccm-8":  ikepolicies.EncryptionAlgorithmAES192CCM8,
	"aes-256-ccm-8":  ikepolicies.EncryptionAlgorithmAES256CCM8,
	"aes-128-ccm-12": ikepolicies.EncryptionAlgorithmAES128CCM12,
	"aes-192-ccm-12": ikepolicies.EncryptionAlgorithmAES192CCM12,
	"aes-256-ccm-12": ikepolicies.EncryptionAlgorithmAES256CCM12,
	"aes-128-ccm-16": ikepolicies.EncryptionAlgorithmAES128CCM16,
	"aes-192-ccm-16": ikepolicies.EncryptionAlgorithmAES192CCM16,
	"aes-256-ccm-16": ikepolicies.EncryptionAlgorithmAES256CCM16,
	"aes-128-gcm-8":  ikepolicies.EncryptionAlgorithmAES128GCM8,
	"aes-192-gcm-8":  ikepolicies.EncryptionAlgorithmAES192GCM8,
	"aes-256-gcm-8":  ikepolicies.EncryptionAlgorithmAES256GCM8,
	"aes-128-gcm-12": ikepolicies.EncryptionAlgorithmAES128GCM12,
	"aes-192-gcm-12": ikepolicies.EncryptionAlgorithmAES192GCM12,
	"aes-256-gcm-12": ikepolicies.EncryptionAlgorithmAES256GCM12,
	"aes-128-gcm-16": ikepolicies.EncryptionAlgorithmAES128GCM16,
	"aes-192-gcm-16": ikepolicies.EncryptionAlgorithmAES192GCM16,
	"aes-256-gcm-16": ikepolicies.EncryptionAlgorithmAES256GCM16,
}
var ikePolicyV2AuthAlgorithmMap = map[string]ikepolicies.AuthAlgorithm{
	"sha1":     ikepolicies.AuthAlgorithmSHA1,
	"sha256":   ikepolicies.AuthAlgorithmSHA256,
	"sha384":   ikepolicies.AuthAlgorithmSHA384,
	"sha512":   ikepolicies.AuthAlgorithmSHA512,
	"aes-xcbc": ikepolicies.AuthAlgorithmAESXCBC,
	"aes-cmac": ikepolicies.AuthAlgorithmAESCMAC,
}

func resourceIKEPolicyV2() *schema.Resource {
	validPFSs := slices.Collect(maps.Keys(ipsecPolicyV2PFSMap))
	validEncryptionAlgorithms := slices.Collect(maps.Keys(ipsecPolicyV2EncryptionAlgorithmMap))
	validAuthAlgorithms := slices.Collect(maps.Keys(ipsecPolicyV2AuthAlgorithmMap))
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
				ValidateFunc: validation.StringInSlice(validAuthAlgorithms, false),
			},
			"encryption_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "aes-128",
				ValidateFunc: validation.StringInSlice(validEncryptionAlgorithms, false),
			},
			"pfs": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "group5",
				ValidateFunc: validation.StringInSlice(validPFSs, false),
			},
			"phase1_negotiation_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "main",
			},
			"ike_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "v1",
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

func resourceIKEPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	lifetime := resourceIKEPolicyV2LifetimeCreateOpts(d.Get("lifetime").(*schema.Set))
	authAlgorithm := resourceIKEPolicyV2AuthAlgorithm(d.Get("auth_algorithm").(string))
	encryptionAlgorithm := resourceIKEPolicyV2EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
	pfs := resourceIKEPolicyV2PFS(d.Get("pfs").(string))
	ikeVersion := resourceIKEPolicyV2IKEVersion(d.Get("ike_version").(string))
	phase1NegotationMode := resourceIKEPolicyV2Phase1NegotiationMode(d.Get("phase1_negotiation_mode").(string))

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

func resourceIKEPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var lifetimeMap = make(map[string]interface{})
	lifetimeMap["units"] = policy.Lifetime.Units
	lifetimeMap["value"] = policy.Lifetime.Value
	var lifetime []map[string]interface{}
	lifetime = append(lifetime, lifetimeMap)
	if err := d.Set("lifetime", &lifetime); err != nil {
		log.Printf("[WARN] unable to set IKE policy lifetime")
	}

	return nil
}

func resourceIKEPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		opts.PFS = resourceIKEPolicyV2PFS(d.Get("pfs").(string))
		hasChange = true
	}
	if d.HasChange("auth_algorithm") {
		opts.AuthAlgorithm = resourceIKEPolicyV2AuthAlgorithm(d.Get("auth_algorithm").(string))
		hasChange = true
	}
	if d.HasChange("encryption_algorithm") {
		opts.EncryptionAlgorithm = resourceIKEPolicyV2EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
		hasChange = true
	}
	if d.HasChange("phase_1_negotiation_mode") {
		opts.Phase1NegotiationMode = resourceIKEPolicyV2Phase1NegotiationMode(d.Get("phase_1_negotiation_mode").(string))
		hasChange = true
	}
	if d.HasChange("ike_version") {
		opts.IKEVersion = resourceIKEPolicyV2IKEVersion(d.Get("ike_version").(string))
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

func resourceIKEPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	return func() (interface{}, string, error) {
		err := ikepolicies.Delete(ctx, networkingClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		return nil, "ACTIVE", err
	}
}

func waitForIKEPolicyCreation(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ikepolicies.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func waitForIKEPolicyUpdate(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ikepolicies.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func resourceIKEPolicyV2AuthAlgorithm(v string) ikepolicies.AuthAlgorithm {
	return ikePolicyV2AuthAlgorithmMap[v]
}

func resourceIKEPolicyV2EncryptionAlgorithm(v string) ikepolicies.EncryptionAlgorithm {
	return ikePolicyV2EncryptionAlgorithmMap[v]
}

func resourceIKEPolicyV2PFS(v string) ikepolicies.PFS {
	return ikePolicyV2PFSMap[v]
}

func resourceIKEPolicyV2IKEVersion(v string) ikepolicies.IKEVersion {
	var ikeVersion ikepolicies.IKEVersion
	switch v {
	case "v1":
		ikeVersion = ikepolicies.IKEVersionv1
	case "v2":
		ikeVersion = ikepolicies.IKEVersionv2
	}
	return ikeVersion
}

func resourceIKEPolicyV2Phase1NegotiationMode(v string) ikepolicies.Phase1NegotiationMode {
	var phase1NegotiationMode ikepolicies.Phase1NegotiationMode
	switch v {
	case "main":
		phase1NegotiationMode = ikepolicies.Phase1NegotiationModeMain
	}
	return phase1NegotiationMode
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
		rawMap := raw.(map[string]interface{})
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
		rawMap := raw.(map[string]interface{})
		lifetimeUpdateOpts.Units = resourceIKEPolicyV2Unit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetimeUpdateOpts.Value = value
	}
	return lifetimeUpdateOpts
}
