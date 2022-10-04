package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
)

func resourceComputeServerGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeServerGroupV2Create,
		ReadContext:   resourceComputeServerGroupV2Read,
		Update:        nil,
		DeleteContext: resourceComputeServerGroupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ForceNew: true,
				Required: true,
			},

			"policies": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"policy": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"policies"},
			},

			"rules": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"policies"},
				MinItems:      1,
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_server_per_host": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},

			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeServerGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	name := d.Get("name").(string)

	rawPolicies := d.Get("policies").([]interface{})
	policies := expandComputeServerGroupV2Policies(computeClient, rawPolicies)

	policy := d.Get("policy").(string)
	rulesVal, rulesPresent := d.GetOk("rules")

	var createOpts ComputeServerGroupV2CreateOpts

	// "policies" is replaced with "policy" and optional "rules" since microversion 2.64
	if policy == "" {
		createOpts = ComputeServerGroupV2CreateOpts{
			servergroups.CreateOpts{
				Name:     name,
				Policies: policies,
			},
			MapValueSpecs(d),
		}
	} else {
		computeClient.Microversion = "2.64"

		if policy == "anti-affinity" && rulesPresent {
			rules := rulesVal.([]map[string]interface{})

			var MaxServerPerHost int
			if v, ok := rules[0]["max_server_per_host"]; ok {
				MaxServerPerHost = v.(int)
			}

			createOpts = ComputeServerGroupV2CreateOpts{
				servergroups.CreateOpts{
					Name:   name,
					Policy: policy,
					Rules: &servergroups.Rules{
						MaxServerPerHost: MaxServerPerHost,
					},
				},
				MapValueSpecs(d),
			}
		} else {
			createOpts = ComputeServerGroupV2CreateOpts{
				servergroups.CreateOpts{
					Name:   name,
					Policy: policy,
				},
				MapValueSpecs(d),
			}
		}
	}

	log.Printf("[DEBUG] openstack_compute_servergroup_v2 create options: %#v", createOpts)
	newSG, err := servergroups.Create(computeClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_compute_servergroup_v2 %s: %s", name, err)
	}

	d.SetId(newSG.ID)

	return resourceComputeServerGroupV2Read(ctx, d, meta)
}

func resourceComputeServerGroupV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	// Attempt to read with unset microversion
	computeClient.Microversion = ""

	sg, err := servergroups.Get(computeClient, d.Id()).Extract()
	if err == nil {
		log.Printf("[DEBUG] Retrieved openstack_compute_servergroup_v2 %s: %#v", d.Id(), sg)

		d.Set("name", sg.Name)
		d.Set("members", sg.Members)
		d.Set("region", GetRegion(d, config))
		d.Set("policies", sg.Policies)

		return nil
	}

	// Attempt to read with microversion 2.64
	computeClient.Microversion = "2.64"

	sg, err = servergroups.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_compute_servergroup_v2"))
	}
	log.Printf("[DEBUG] Retrieved openstack_compute_servergroup_v2 %s: %#v", d.Id(), sg)

	d.Set("name", sg.Name)
	d.Set("members", sg.Members)
	d.Set("region", GetRegion(d, config))
	d.Set("policy", sg.Policy)

	rules := make(map[string]interface{})
	rules["max_server_per_host"] = sg.Rules.MaxServerPerHost
	rulesList := []map[string]interface{}{rules}
	d.Set("rules", rulesList)

	return nil
}

func resourceComputeServerGroupV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	if err := servergroups.Delete(computeClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_compute_servergroup_v2"))
	}

	return nil
}
