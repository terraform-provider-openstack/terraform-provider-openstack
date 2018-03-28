package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func resourceIPSecPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPSecPolicyV2Create,
		Read:   resourceIPSecPolicyV2Read,
		Update: resourceIPSecPolicyV2Update,
		Delete: resourceIPSecPolicyV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_algorithm": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encapsulation_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pfs": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encryption_algorithm": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"transform_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"lifetime": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"units": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"value_specs": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIPSecPolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	encapsulationMode := IPSecPolicyV2GetEncapsulationMode(d.Get("encapsulation_mode").(string))
	authAlgorithm := IPSecPolicyV2GetAuthAlgorithm(d.Get("auth_algorithm").(string))
	encryptionAlgorithm := IPSecPolicyV2GetEncryptionAlgorithm(d.Get("encryption_algorithm").(string))
	pfs := IPSecPolicyV2GetPFS(d.Get("pfs").(string))
	transformProtocol := IPSecPolicyV2GetTransformProtocol(d.Get("transform_protocol").(string))
	lifetime := IPSecPolicyV2GetLifetimeCreateOpts(d.Get("lifetime").(map[string]interface{}))

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

	policy, err := ipsecpolicies.Create(networkingClient, opts).Extract()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForIPSecPolicyCreation(networkingClient, policy.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForState()

	log.Printf("[DEBUG] IPSec policy created: %#v", policy)

	d.SetId(policy.ID)

	return resourceIPSecPolicyV2Read(d, meta)
}

func resourceIPSecPolicyV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Retrieve information about IPSec policy: %s", d.Id())

	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	policy, err := ipsecpolicies.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "IPSec policy")
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
	var lifetimeMap map[string]interface{}
	lifetimeMap = make(map[string]interface{})
	lifetimeMap["units"] = policy.Lifetime.Units
	lifetimeMap["value"] = policy.Lifetime.Value
	d.Set("lifetime", &lifetimeMap)

	return nil
}

func resourceIPSecPolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
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
		opts.AuthAlgorithm = IPSecPolicyV2GetAuthAlgorithm(d.Get("auth_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("encryption_algorithm") {
		opts.EncryptionAlgorithm = IPSecPolicyV2GetEncryptionAlgorithm(d.Get("encryption_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("transform_protocol") {
		opts.TransformProtocol = IPSecPolicyV2GetTransformProtocol(d.Get("transform_protocol").(string))
		hasChange = true
	}

	if d.HasChange("pfs") {
		opts.PFS = IPSecPolicyV2GetPFS(d.Get("pfs").(string))
		hasChange = true
	}

	if d.HasChange("encapsulation_mode") {
		opts.EncapsulationMode = IPSecPolicyV2GetEncapsulationMode(d.Get("encapsulation_mode").(string))
		hasChange = true
	}

	if d.HasChange("lifetime") {
		opts.Lifetime = IPSecPolicyV2GetLifetimeUpdateOpts(d.Get("lifetime").(schema.ResourceData))
		hasChange = true
	}

	log.Printf("[DEBUG] Updating IPSec policy with id %s: %#v", d.Id(), opts)

	if hasChange {
		_, err = ipsecpolicies.Update(networkingClient, d.Id(), opts).Extract()
		if err != nil {
			return err
		}
	}
	return resourceIPSecPolicyV2Read(d, meta)
}

func resourceIPSecPolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy IPSec policy: %s", d.Id())

	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForIPSecPolicyDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return err
	}

	return nil
}

func waitForIPSecPolicyDeletion(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := ipsecpolicies.Delete(networkingClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if errCode, ok := err.(gophercloud.ErrUnexpectedResponseCode); ok {
			if errCode.Actual == 409 {
				return nil, "ACTIVE", nil
			}
		}

		return nil, "ACTIVE", err
	}
}

func waitForIPSecPolicyCreation(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ipsecpolicies.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func IPSecPolicyV2GetTransformProtocol(trp string) ipsecpolicies.TransformProtocol {
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
func IPSecPolicyV2GetPFS(pfsString string) ipsecpolicies.PFS {
	var pfs ipsecpolicies.PFS
	switch pfsString {
	case "group2":
		pfs = ipsecpolicies.PFSGroup2
	case "group5":
		pfs = ipsecpolicies.PFSGroup5
	case "group14":
		pfs = ipsecpolicies.PFSGroup14
	}
	return pfs

}
func IPSecPolicyV2GetEncryptionAlgorithm(encryptionAlgo string) ipsecpolicies.EncryptionAlgorithm {
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
	}
	return alg
}
func IPSecPolicyV2GetAuthAlgorithm(authAlgo string) ipsecpolicies.AuthAlgorithm {
	var alg ipsecpolicies.AuthAlgorithm
	switch authAlgo {
	case "sha1":
		alg = ipsecpolicies.AuthAlgorithmSHA1
	case "sha256":
		alg = ipsecpolicies.AuthAlgorithmSHA256
	case "sha384":
		alg = ipsecpolicies.AuthAlgorithmSHA384
	case "sha512":
		alg = "sha384"
	}
	return alg
}
func IPSecPolicyV2GetEncapsulationMode(encMode string) ipsecpolicies.EncapsulationMode {
	var mode ipsecpolicies.EncapsulationMode
	switch encMode {
	case "tunnel":
		mode = ipsecpolicies.EncapsulationModeTunnel
	case "transport":
		mode = ipsecpolicies.EncapsulationModeTransport
	}
	return mode
}

func IPSecPolicyV2GetLifetimeCreateOpts(d map[string]interface{}) ipsecpolicies.LifetimeCreateOpts {
	lifetime := ipsecpolicies.LifetimeCreateOpts{}
	if val, ok := d["units"]; ok {
		unit := IPSecPolicyV2GetUnit(val.(string))
		lifetime.Units = unit
	}
	if val, ok := d["value"]; ok {
		value, err := strconv.Atoi(val.(string))
		if err != nil {
			panic(err)
		}
		lifetime.Value = value
	}
	return lifetime
}

func IPSecPolicyV2GetUnit(units string) ipsecpolicies.Unit {
	var unit ipsecpolicies.Unit
	switch units {
	case "seconds":
		unit = ipsecpolicies.UnitSeconds
	case "kilobytes":
		unit = ipsecpolicies.UnitKilobytes
	}
	return unit
}

func IPSecPolicyV2GetLifetimeUpdateOpts(d schema.ResourceData) *ipsecpolicies.LifetimeUpdateOpts {
	var unit ipsecpolicies.Unit
	var value int
	unit = IPSecPolicyV2GetUnit(d.Get("units").(string))
	value = d.Get("value").(int)
	updateOpts := ipsecpolicies.LifetimeUpdateOpts{
		Units: unit,
		Value: value,
	}
	return &updateOpts

}
