package openstack

import (
	"context"
	"log"
	"maps"
	"net/http"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
)

var ipsecPolicyV2PFSMap = map[string]ipsecpolicies.PFS{
	"group2":  ipsecpolicies.PFSGroup2,
	"group5":  ipsecpolicies.PFSGroup5,
	"group14": ipsecpolicies.PFSGroup14,
	"group15": ipsecpolicies.PFSGroup15,
	"group16": ipsecpolicies.PFSGroup16,
	"group17": ipsecpolicies.PFSGroup17,
	"group18": ipsecpolicies.PFSGroup18,
	"group19": ipsecpolicies.PFSGroup19,
	"group20": ipsecpolicies.PFSGroup20,
	"group21": ipsecpolicies.PFSGroup21,
	"group22": ipsecpolicies.PFSGroup22,
	"group23": ipsecpolicies.PFSGroup23,
	"group24": ipsecpolicies.PFSGroup24,
	"group25": ipsecpolicies.PFSGroup25,
	"group26": ipsecpolicies.PFSGroup26,
	"group27": ipsecpolicies.PFSGroup27,
	"group28": ipsecpolicies.PFSGroup28,
	"group29": ipsecpolicies.PFSGroup29,
	"group30": ipsecpolicies.PFSGroup30,
	"group31": ipsecpolicies.PFSGroup31,
}

var ipsecPolicyV2EncryptionAlgorithmMap = map[string]ipsecpolicies.EncryptionAlgorithm{
	"3des":           ipsecpolicies.EncryptionAlgorithm3DES,
	"aes-128":        ipsecpolicies.EncryptionAlgorithmAES128,
	"aes-256":        ipsecpolicies.EncryptionAlgorithmAES256,
	"aes-192":        ipsecpolicies.EncryptionAlgorithmAES192,
	"aes-128-ctr":    ipsecpolicies.EncryptionAlgorithmAES128CTR,
	"aes-192-ctr":    ipsecpolicies.EncryptionAlgorithmAES192CTR,
	"aes-256-ctr":    ipsecpolicies.EncryptionAlgorithmAES256CTR,
	"aes-128-ccm-8":  ipsecpolicies.EncryptionAlgorithmAES128CCM8,
	"aes-192-ccm-8":  ipsecpolicies.EncryptionAlgorithmAES192CCM8,
	"aes-256-ccm-8":  ipsecpolicies.EncryptionAlgorithmAES256CCM8,
	"aes-128-ccm-12": ipsecpolicies.EncryptionAlgorithmAES128CCM12,
	"aes-192-ccm-12": ipsecpolicies.EncryptionAlgorithmAES192CCM12,
	"aes-256-ccm-12": ipsecpolicies.EncryptionAlgorithmAES256CCM12,
	"aes-128-ccm-16": ipsecpolicies.EncryptionAlgorithmAES128CCM16,
	"aes-192-ccm-16": ipsecpolicies.EncryptionAlgorithmAES192CCM16,
	"aes-256-ccm-16": ipsecpolicies.EncryptionAlgorithmAES256CCM16,
	"aes-128-gcm-8":  ipsecpolicies.EncryptionAlgorithmAES128GCM8,
	"aes-192-gcm-8":  ipsecpolicies.EncryptionAlgorithmAES192GCM8,
	"aes-256-gcm-8":  ipsecpolicies.EncryptionAlgorithmAES256GCM8,
	"aes-128-gcm-12": ipsecpolicies.EncryptionAlgorithmAES128GCM12,
	"aes-192-gcm-12": ipsecpolicies.EncryptionAlgorithmAES192GCM12,
	"aes-256-gcm-12": ipsecpolicies.EncryptionAlgorithmAES256GCM12,
	"aes-128-gcm-16": ipsecpolicies.EncryptionAlgorithmAES128GCM16,
	"aes-192-gcm-16": ipsecpolicies.EncryptionAlgorithmAES192GCM16,
	"aes-256-gcm-16": ipsecpolicies.EncryptionAlgorithmAES256GCM16,
}
var ipsecPolicyV2AuthAlgorithmMap = map[string]ipsecpolicies.AuthAlgorithm{
	"sha1":     ipsecpolicies.AuthAlgorithmSHA1,
	"sha256":   ipsecpolicies.AuthAlgorithmSHA256,
	"sha384":   ipsecpolicies.AuthAlgorithmSHA384,
	"sha512":   ipsecpolicies.AuthAlgorithmSHA512,
	"aes-xcbc": ipsecpolicies.AuthAlgorithmAESXCBC,
	"aes-cmac": ipsecpolicies.AuthAlgorithmAESCMAC,
}

func resourceIPSecPolicyV2() *schema.Resource {
	validPFSs := slices.Collect(maps.Keys(ipsecPolicyV2PFSMap))
	validEncryptionAlgorithms := slices.Collect(maps.Keys(ipsecPolicyV2EncryptionAlgorithmMap))
	validAuthAlgorithms := slices.Collect(maps.Keys(ipsecPolicyV2AuthAlgorithmMap))
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
				ValidateFunc: validation.StringInSlice(validAuthAlgorithms, false),
			},
			"encapsulation_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tunnel", "transport",
				}, false),
			},
			"pfs": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(validPFSs, false),
			},
			"encryption_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(validEncryptionAlgorithms, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transform_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"esp", "ah", "ah-esp",
				}, false),
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

func resourceIPSecPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	encapsulationMode := resourceIPSecPolicyV2EncapsulationMode(d.Get("encapsulation_mode").(string))
	authAlgorithm := resourceIPSecPolicyV2AuthAlgorithm(d.Get("auth_algorithm").(string))
	encryptionAlgorithm := resourceIPSecPolicyV2EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
	pfs := resourceIPSecPolicyV2PFS(d.Get("pfs").(string))
	transformProtocol := resourceIPSecPolicyV2TransformProtocol(d.Get("transform_protocol").(string))
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

func resourceIPSecPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var lifetimeMap = make(map[string]interface{})
	lifetimeMap["units"] = policy.Lifetime.Units
	lifetimeMap["value"] = policy.Lifetime.Value
	var lifetime []map[string]interface{}
	lifetime = append(lifetime, lifetimeMap)
	if err := d.Set("lifetime", &lifetime); err != nil {
		log.Printf("[WARN] unable to set IPSec policy lifetime")
	}

	return nil
}

func resourceIPSecPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		opts.AuthAlgorithm = resourceIPSecPolicyV2AuthAlgorithm(d.Get("auth_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("encryption_algorithm") {
		opts.EncryptionAlgorithm = resourceIPSecPolicyV2EncryptionAlgorithm(d.Get("encryption_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("transform_protocol") {
		opts.TransformProtocol = resourceIPSecPolicyV2TransformProtocol(d.Get("transform_protocol").(string))
		hasChange = true
	}

	if d.HasChange("pfs") {
		opts.PFS = resourceIPSecPolicyV2PFS(d.Get("pfs").(string))
		hasChange = true
	}

	if d.HasChange("encapsulation_mode") {
		opts.EncapsulationMode = resourceIPSecPolicyV2EncapsulationMode(d.Get("encapsulation_mode").(string))
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

func resourceIPSecPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	return func() (interface{}, string, error) {
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
	return func() (interface{}, string, error) {
		policy, err := ipsecpolicies.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func waitForIPSecPolicyUpdate(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ipsecpolicies.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func resourceIPSecPolicyV2TransformProtocol(trp string) ipsecpolicies.TransformProtocol {
	var protocol ipsecpolicies.TransformProtocol
	switch trp {
	case "esp":
		protocol = ipsecpolicies.TransformProtocolESP
	case "ah":
		protocol = ipsecpolicies.TransformProtocolAH
	case "ah-esp":
		protocol = ipsecpolicies.TransformProtocolAHESP
	}
	return protocol
}
func resourceIPSecPolicyV2PFS(pfsString string) ipsecpolicies.PFS {
	return ipsecPolicyV2PFSMap[pfsString]
}
func resourceIPSecPolicyV2EncryptionAlgorithm(encryptionAlgo string) ipsecpolicies.EncryptionAlgorithm {
	return ipsecPolicyV2EncryptionAlgorithmMap[encryptionAlgo]
}
func resourceIPSecPolicyV2AuthAlgorithm(authAlgo string) ipsecpolicies.AuthAlgorithm {
	return ipsecPolicyV2AuthAlgorithmMap[authAlgo]
}
func resourceIPSecPolicyV2EncapsulationMode(encMode string) ipsecpolicies.EncapsulationMode {
	var mode ipsecpolicies.EncapsulationMode
	switch encMode {
	case "tunnel":
		mode = ipsecpolicies.EncapsulationModeTunnel
	case "transport":
		mode = ipsecpolicies.EncapsulationModeTransport
	}
	return mode
}

func resourceIPSecPolicyV2LifetimeCreateOpts(d *schema.Set) ipsecpolicies.LifetimeCreateOpts {
	lifetime := ipsecpolicies.LifetimeCreateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
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
		rawMap := raw.(map[string]interface{})
		lifetimeUpdateOpts.Units = resourceIPSecPolicyV2Unit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetimeUpdateOpts.Value = value
	}
	return lifetimeUpdateOpts
}
