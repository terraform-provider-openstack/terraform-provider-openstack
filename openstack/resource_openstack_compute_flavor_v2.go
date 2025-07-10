package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComputeFlavorV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeFlavorV2Create,
		ReadContext:   resourceComputeFlavorV2Read,
		UpdateContext: resourceComputeFlavorV2Update,
		DeleteContext: resourceComputeFlavorV2Delete,
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

			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ram": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"vcpus": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"disk": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"swap": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"rx_tx_factor": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},

			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"ephemeral": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"extra_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceComputeFlavorV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	name := d.Get("name").(string)
	disk := d.Get("disk").(int)
	swap := d.Get("swap").(int)
	isPublic := d.Get("is_public").(bool)
	ephemeral := d.Get("ephemeral").(int)
	createOpts := flavors.CreateOpts{
		Name:       name,
		RAM:        d.Get("ram").(int),
		VCPUs:      d.Get("vcpus").(int),
		Disk:       &disk,
		ID:         d.Get("flavor_id").(string),
		Swap:       &swap,
		RxTxFactor: d.Get("rx_tx_factor").(float64),
		IsPublic:   &isPublic,
		Ephemeral:  &ephemeral,
	}

	// description requires nova microversion 2.55. Only set when specified.
	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		createOpts.Description = description
		computeClient.Microversion = computeV2FlavorDescriptionMicroversion
	}

	log.Printf("[DEBUG] openstack_compute_flavor_v2 create options: %#v", createOpts)

	fl, err := flavors.Create(ctx, computeClient, &createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_compute_flavor_v2 %s: %s", name, err)
	}

	d.SetId(fl.ID)

	extraSpecsRaw := d.Get("extra_specs").(map[string]any)
	if len(extraSpecsRaw) > 0 {
		extraSpecs := expandComputeFlavorV2ExtraSpecs(extraSpecsRaw)

		_, err := flavors.CreateExtraSpecs(ctx, computeClient, fl.ID, extraSpecs).Extract()
		if err != nil {
			return diag.Errorf("Error creating extra_specs for openstack_compute_flavor_v2 %s: %s", fl.ID, err)
		}
	}

	return resourceComputeFlavorV2Read(ctx, d, meta)
}

func resourceComputeFlavorV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	var fl *flavors.Flavor

	// try and read flavor using microversion that includes description
	computeClient.Microversion = computeV2FlavorDescriptionMicroversion
	fl, err = flavors.Get(ctx, computeClient, d.Id()).Extract()
	if err != nil {
		// reset microversion to 2.1 and try again
		computeClient.Microversion = "2.1"

		fl, err = flavors.Get(ctx, computeClient, d.Id()).Extract()
		if err != nil {
			return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_compute_flavor_v2"))
		}
	}

	log.Printf("[DEBUG] Retrieved openstack_compute_flavor_v2 %s: %#v", d.Id(), fl)

	d.Set("name", fl.Name)
	d.Set("description", fl.Description)
	d.Set("ram", fl.RAM)
	d.Set("vcpus", fl.VCPUs)
	d.Set("disk", fl.Disk)
	d.Set("flavor_id", fl.ID)
	d.Set("swap", fl.Swap)
	d.Set("rx_tx_factor", fl.RxTxFactor)
	d.Set("is_public", fl.IsPublic)
	d.Set("ephemeral", fl.Ephemeral)
	d.Set("region", GetRegion(d, config))

	es, err := flavors.ListExtraSpecs(ctx, computeClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Error reading extra_specs for openstack_compute_flavor_v2 %s: %s", d.Id(), err)
	}

	if err := d.Set("extra_specs", es); err != nil {
		log.Printf("[WARN] Unable to set extra_specs for openstack_compute_flavor_v2 %s: %s", d.Id(), err)
	}

	return nil
}

func resourceComputeFlavorV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	var hasChange bool

	var updateOpts flavors.UpdateOpts

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
		computeClient.Microversion = computeV2FlavorDescriptionMicroversion
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_compute_flavor_v2 %s update options: %#v", d.Id(), updateOpts)

		_, err = flavors.Update(ctx, computeClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error openstack_compute_flavor_v2 %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("extra_specs") {
		oldES, newES := d.GetChange("extra_specs")

		// Delete all old extra specs.
		for oldKey := range oldES.(map[string]any) {
			if err := flavors.DeleteExtraSpec(ctx, computeClient, d.Id(), oldKey).ExtractErr(); err != nil {
				return diag.Errorf("Error deleting extra_spec %s from openstack_compute_flavor_v2 %s: %s", oldKey, d.Id(), err)
			}
		}

		// Add new extra specs.
		newESRaw := newES.(map[string]any)
		if len(newESRaw) > 0 {
			extraSpecs := expandComputeFlavorV2ExtraSpecs(newESRaw)

			_, err := flavors.CreateExtraSpecs(ctx, computeClient, d.Id(), extraSpecs).Extract()
			if err != nil {
				return diag.Errorf("Error creating extra_specs for openstack_compute_flavor_v2 %s: %s", d.Id(), err)
			}
		}
	}

	return resourceComputeFlavorV2Read(ctx, d, meta)
}

func resourceComputeFlavorV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	err = flavors.Delete(ctx, computeClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_compute_flavor_v2"))
	}

	return nil
}
