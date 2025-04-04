package openstack

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"sha1", "sha256", "sha384", "sha512", "aes-xcbc", "aes-cmac",
				}, false),
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"group2", "group5", "group14", "group15", "group16",
					"group17", "group18", "group19", "group20", "group21",
					"group22", "group23", "group24", "group25", "group26",
					"group27", "group28", "group29", "group30", "group31",
				}, false),
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
	var pfs ipsecpolicies.PFS
	switch pfsString {
	case "group2":
		pfs = ipsecpolicies.PFSGroup2
	case "group5":
		pfs = ipsecpolicies.PFSGroup5
	case "group14":
		pfs = ipsecpolicies.PFSGroup14
	case "group15":
		pfs = ipsecpolicies.PFSGroup15
	case "group16":
		pfs = ipsecpolicies.PFSGroup16
	case "group17":
		pfs = ipsecpolicies.PFSGroup17
	case "group18":
		pfs = ipsecpolicies.PFSGroup18
	case "group19":
		pfs = ipsecpolicies.PFSGroup19
	case "group20":
		pfs = ipsecpolicies.PFSGroup20
	case "group21":
		pfs = ipsecpolicies.PFSGroup21
	case "group22":
		pfs = ipsecpolicies.PFSGroup22
	case "group23":
		pfs = ipsecpolicies.PFSGroup23
	case "group24":
		pfs = ipsecpolicies.PFSGroup24
	case "group25":
		pfs = ipsecpolicies.PFSGroup25
	case "group26":
		pfs = ipsecpolicies.PFSGroup26
	case "group27":
		pfs = ipsecpolicies.PFSGroup27
	case "group28":
		pfs = ipsecpolicies.PFSGroup28
	case "group29":
		pfs = ipsecpolicies.PFSGroup29
	case "group30":
		pfs = ipsecpolicies.PFSGroup30
	case "group31":
		pfs = ipsecpolicies.PFSGroup31
	}
	return pfs
}
func resourceIPSecPolicyV2EncryptionAlgorithm(encryptionAlgo string) ipsecpolicies.EncryptionAlgorithm {
	var alg ipsecpolicies.EncryptionAlgorithm
	switch encryptionAlgo {
	case "3des":
		alg = ipsecpolicies.EncryptionAlgorithm3DES
	case "aes-128":
		alg = ipsecpolicies.EncryptionAlgorithmAES128
	case "aes-256":
		alg = ipsecpolicies.EncryptionAlgorithmAES256
	case "aes-192":
		alg = ipsecpolicies.EncryptionAlgorithmAES192
	case "aes-128-ctr":
		alg = ipsecpolicies.EncryptionAlgorithmAES128CTR
	case "aes-192-ctr":
		alg = ipsecpolicies.EncryptionAlgorithmAES192CTR
	case "aes-256-ctr":
		alg = ipsecpolicies.EncryptionAlgorithmAES256CTR
	case "aes-128-ccm-8":
		alg = ipsecpolicies.EncryptionAlgorithmAES128CCM8
	case "aes-192-ccm-8":
		alg = ipsecpolicies.EncryptionAlgorithmAES192CCM8
	case "aes-256-ccm-8":
		alg = ipsecpolicies.EncryptionAlgorithmAES256CCM8
	case "aes-128-ccm-12":
		alg = ipsecpolicies.EncryptionAlgorithmAES128CCM12
	case "aes-192-ccm-12":
		alg = ipsecpolicies.EncryptionAlgorithmAES192CCM12
	case "aes-256-ccm-12":
		alg = ipsecpolicies.EncryptionAlgorithmAES256CCM12
	case "aes-128-ccm-16":
		alg = ipsecpolicies.EncryptionAlgorithmAES128CCM16
	case "aes-192-ccm-16":
		alg = ipsecpolicies.EncryptionAlgorithmAES192CCM16
	case "aes-256-ccm-16":
		alg = ipsecpolicies.EncryptionAlgorithmAES256CCM16
	case "aes-128-gcm-8":
		alg = ipsecpolicies.EncryptionAlgorithmAES128GCM8
	case "aes-192-gcm-8":
		alg = ipsecpolicies.EncryptionAlgorithmAES192GCM8
	case "aes-256-gcm-8":
		alg = ipsecpolicies.EncryptionAlgorithmAES256GCM8
	case "aes-128-gcm-12":
		alg = ipsecpolicies.EncryptionAlgorithmAES128GCM12
	case "aes-192-gcm-12":
		alg = ipsecpolicies.EncryptionAlgorithmAES192GCM12
	case "aes-256-gcm-12":
		alg = ipsecpolicies.EncryptionAlgorithmAES256GCM12
	case "aes-128-gcm-16":
		alg = ipsecpolicies.EncryptionAlgorithmAES128GCM16
	case "aes-192-gcm-16":
		alg = ipsecpolicies.EncryptionAlgorithmAES192GCM16
	case "aes-256-gcm-16":
		alg = ipsecpolicies.EncryptionAlgorithmAES256GCM16
	}
	return alg
}
func resourceIPSecPolicyV2AuthAlgorithm(authAlgo string) ipsecpolicies.AuthAlgorithm {
	var alg ipsecpolicies.AuthAlgorithm
	switch authAlgo {
	case "sha1":
		alg = ipsecpolicies.AuthAlgorithmSHA1
	case "sha256":
		alg = ipsecpolicies.AuthAlgorithmSHA256
	case "sha384":
		alg = ipsecpolicies.AuthAlgorithmSHA384
	case "sha512":
		alg = ipsecpolicies.AuthAlgorithmSHA512
	case "aes-xcbc":
		alg = ipsecpolicies.AuthAlgorithmAESXCBC
	case "aes-cmac":
		alg = ipsecpolicies.AuthAlgorithmAESCMAC
	}
	return alg
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
