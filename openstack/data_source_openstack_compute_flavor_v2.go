package openstack

import (
	"context"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceComputeFlavorV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeFlavorV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flavor_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name", "min_ram", "min_disk"},
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"flavor_id"},
			},

			"min_ram": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"flavor_id"},
			},

			"ram": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"min_disk": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"flavor_id"},
			},

			"disk": {
				Type:     schema.TypeInt,
				Optional: true,
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
			},

			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Computed values
			"extra_specs": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

// dataSourceComputeFlavorV2Read performs the flavor lookup.
func dataSourceComputeFlavorV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	var allFlavors []flavors.Flavor

	if v := d.Get("flavor_id").(string); v != "" {
		var flavor *flavors.Flavor
		// try and read flavor using microversion that includes description
		computeClient.Microversion = computeV2FlavorDescriptionMicroversion
		flavor, err = flavors.Get(ctx, computeClient, v).Extract()
		if err != nil {
			// reset microversion to 2.1 and try again
			computeClient.Microversion = "2.1"

			flavor, err = flavors.Get(ctx, computeClient, v).Extract()
			if err != nil {
				if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
					return diag.Errorf("No Flavor found")
				}

				return diag.Errorf("Unable to retrieve OpenStack %s flavor: %s", v, err)
			}
		}

		allFlavors = append(allFlavors, *flavor)
	} else {
		accessType := flavors.AllAccess

		if v, ok := getOkExists(d, "is_public"); ok {
			if v, ok := v.(bool); ok {
				if v {
					accessType = flavors.PublicAccess
				} else {
					accessType = flavors.PrivateAccess
				}
			}
		}

		listOpts := flavors.ListOpts{
			MinDisk:    d.Get("min_disk").(int),
			MinRAM:     d.Get("min_ram").(int),
			AccessType: accessType,
		}

		log.Printf("[DEBUG] openstack_compute_flavor_v2 ListOpts: %#v", listOpts)

		var allPages pagination.Page
		// try and read flavor using microversion that includes description
		computeClient.Microversion = computeV2FlavorDescriptionMicroversion
		allPages, err = flavors.ListDetail(computeClient, listOpts).AllPages(ctx)
		if err != nil {
			// reset microversion to 2.1 and try again
			computeClient.Microversion = "2.1"

			allPages, err = flavors.ListDetail(computeClient, listOpts).AllPages(ctx)
			if err != nil {
				return diag.Errorf("Unable to query OpenStack flavors: %s", err)
			}
		}

		allFlavors, err = flavors.ExtractFlavors(allPages)
		if err != nil {
			return diag.Errorf("Unable to retrieve OpenStack flavors: %s", err)
		}
	}

	// Loop through all flavors to find a more specific one.
	if len(allFlavors) > 0 {
		var filteredFlavors []flavors.Flavor

		for _, flavor := range allFlavors {
			if v := d.Get("name").(string); v != "" {
				if flavor.Name != v {
					continue
				}
			}

			if v := d.Get("description").(string); v != "" {
				if flavor.Description != v {
					continue
				}
			}

			// d.GetOk is used because 0 might be a valid choice.
			if v, ok := d.GetOk("ram"); ok {
				if flavor.RAM != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("vcpus"); ok {
				if flavor.VCPUs != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("disk"); ok {
				if flavor.Disk != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("swap"); ok {
				if flavor.Swap != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("rx_tx_factor"); ok {
				if flavor.RxTxFactor != v.(float64) {
					continue
				}
			}

			filteredFlavors = append(filteredFlavors, flavor)
		}

		allFlavors = filteredFlavors
	}

	if len(allFlavors) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allFlavors) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allFlavors)

		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria")
	}

	err = dataSourceComputeFlavorV2Attributes(ctx, d, computeClient, &allFlavors[0])
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("region", GetRegion(d, config))

	return nil
}

// dataSourceComputeFlavorV2Attributes populates the fields of a Flavor resource.
func dataSourceComputeFlavorV2Attributes(ctx context.Context, d *schema.ResourceData, computeClient *gophercloud.ServiceClient, flavor *flavors.Flavor) error {
	log.Printf("[DEBUG] Retrieved openstack_compute_flavor_v2 %s: %#v", flavor.ID, flavor)

	d.SetId(flavor.ID)
	d.Set("name", flavor.Name)
	d.Set("description", flavor.Description)
	d.Set("flavor_id", flavor.ID)
	d.Set("disk", flavor.Disk)
	d.Set("ram", flavor.RAM)
	d.Set("rx_tx_factor", flavor.RxTxFactor)
	d.Set("swap", flavor.Swap)
	d.Set("vcpus", flavor.VCPUs)
	d.Set("is_public", flavor.IsPublic)

	es, err := flavors.ListExtraSpecs(ctx, computeClient, d.Id()).Extract()
	if err != nil {
		return err
	}

	if err := d.Set("extra_specs", es); err != nil {
		log.Printf("[WARN] Unable to set extra_specs for openstack_compute_flavor_v2 %s: %s", d.Id(), err)
	}

	return nil
}
