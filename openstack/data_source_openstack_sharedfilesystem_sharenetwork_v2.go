package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharenetworks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSharedFilesystemShareNetworkV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSharedFilesystemShareNetworkV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"neutron_net_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"neutron_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"security_service_id": {
				Type:     schema.TypeString,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"security_service_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"segmentation_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceSharedFilesystemShareNetworkV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem sfsClient: %s", err)
	}

	listOpts := sharenetworks.ListOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		ProjectID:       d.Get("project_id").(string),
		NeutronNetID:    d.Get("neutron_net_id").(string),
		NeutronSubnetID: d.Get("neutron_subnet_id").(string),
		NetworkType:     d.Get("network_type").(string),
	}

	if v, ok := getOkExists(d, "ip_version"); ok {
		listOpts.IPVersion = gophercloud.IPVersion(v.(int))
	}

	if v, ok := getOkExists(d, "segmentation_id"); ok {
		listOpts.SegmentationID = v.(int)
	}

	allPages, err := sharenetworks.ListDetail(sfsClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query share networks: %s", err)
	}

	allShareNetworks, err := sharenetworks.ExtractShareNetworks(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve share networks: %s", err)
	}

	if len(allShareNetworks) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var securityServiceID string

	var securityServiceIDs []string

	if v, ok := getOkExists(d, "security_service_id"); ok {
		// filtering by "security_service_id"
		securityServiceID = v.(string)

		var filteredShareNetworks []sharenetworks.ShareNetwork

		log.Printf("[DEBUG] Filtering share networks by a %s security service ID", securityServiceID)

		for _, shareNetwork := range allShareNetworks {
			tmp, err := resourceSharedFilesystemShareNetworkV2GetSvcByShareNetID(ctx, sfsClient, shareNetwork.ID)
			if err != nil {
				return diag.FromErr(err)
			}

			if strSliceContains(tmp, securityServiceID) {
				filteredShareNetworks = append(filteredShareNetworks, shareNetwork)
				securityServiceIDs = tmp
			}
		}

		if len(filteredShareNetworks) == 0 {
			return diag.Errorf("Your query returned no results after the security service ID filter. " +
				"Please change your search criteria and try again")
		}

		allShareNetworks = filteredShareNetworks
	}

	var shareNetwork sharenetworks.ShareNetwork

	if len(allShareNetworks) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allShareNetworks)

		return diag.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	shareNetwork = allShareNetworks[0]

	// skip extra calls if "security_service_id" filter was already used
	if securityServiceID == "" {
		securityServiceIDs, err = resourceSharedFilesystemShareNetworkV2GetSvcByShareNetID(ctx, sfsClient, shareNetwork.ID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(shareNetwork.ID)
	d.Set("name", shareNetwork.Name)
	d.Set("description", shareNetwork.Description)
	d.Set("project_id", shareNetwork.ProjectID)
	d.Set("neutron_net_id", shareNetwork.NeutronNetID)
	d.Set("neutron_subnet_id", shareNetwork.NeutronSubnetID)
	d.Set("security_service_ids", securityServiceIDs)
	d.Set("network_type", shareNetwork.NetworkType)
	d.Set("ip_version", shareNetwork.IPVersion)
	d.Set("segmentation_id", shareNetwork.SegmentationID)
	d.Set("cidr", shareNetwork.CIDR)
	d.Set("region", GetRegion(d, config))

	return nil
}
