package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/securityservices"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharenetworks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSharedFilesystemShareNetworkV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFilesystemShareNetworkV2Create,
		ReadContext:   resourceSharedFilesystemShareNetworkV2Read,
		UpdateContext: resourceSharedFilesystemShareNetworkV2Update,
		DeleteContext: resourceSharedFilesystemShareNetworkV2Delete,
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

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"neutron_net_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"neutron_subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"security_service_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"segmentation_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSharedFilesystemShareNetworkV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	createOpts := sharenetworks.CreateOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		NeutronNetID:    d.Get("neutron_net_id").(string),
		NeutronSubnetID: d.Get("neutron_subnet_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	log.Printf("[DEBUG] Attempting to create sharenetwork")

	sharenetwork, err := sharenetworks.Create(ctx, sfsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating sharenetwork: %s", err)
	}

	d.SetId(sharenetwork.ID)

	securityServiceIDs := resourceSharedFilesystemShareNetworkV2SecSvcToArray(d.Get("security_service_ids").(*schema.Set))
	for _, securityServiceID := range securityServiceIDs {
		log.Printf("[DEBUG] Adding %s security service to sharenetwork %s", securityServiceID, sharenetwork.ID)
		securityServiceOpts := sharenetworks.AddSecurityServiceOpts{SecurityServiceID: securityServiceID}

		_, err = sharenetworks.AddSecurityService(ctx, sfsClient, sharenetwork.ID, securityServiceOpts).Extract()
		if err != nil {
			return diag.Errorf("Error adding %s security service to sharenetwork: %s", securityServiceID, err)
		}
	}

	return resourceSharedFilesystemShareNetworkV2Read(ctx, d, meta)
}

func resourceSharedFilesystemShareNetworkV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sharenetwork, err := sharenetworks.Get(ctx, sfsClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "sharenetwork"))
	}

	log.Printf("[DEBUG] Retrieved sharenetwork %s: %#v", d.Id(), sharenetwork)

	securityServiceIDs, err := resourceSharedFilesystemShareNetworkV2GetSvcByShareNetID(ctx, sfsClient, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("security_service_ids", securityServiceIDs)

	d.Set("name", sharenetwork.Name)
	d.Set("description", sharenetwork.Description)
	d.Set("neutron_net_id", sharenetwork.NeutronNetID)
	d.Set("neutron_subnet_id", sharenetwork.NeutronSubnetID)
	// Computed
	d.Set("project_id", sharenetwork.ProjectID)
	d.Set("region", GetRegion(d, config))
	d.Set("network_type", sharenetwork.NetworkType)
	d.Set("segmentation_id", sharenetwork.SegmentationID)
	d.Set("cidr", sharenetwork.CIDR)
	d.Set("ip_version", sharenetwork.IPVersion)

	return nil
}

func resourceSharedFilesystemShareNetworkV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	var updateOpts sharenetworks.UpdateOpts

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("neutron_net_id") {
		updateOpts.NeutronNetID = d.Get("neutron_net_id").(string)
	}

	if d.HasChange("neutron_subnet_id") {
		updateOpts.NeutronSubnetID = d.Get("neutron_subnet_id").(string)
	}

	if updateOpts != (sharenetworks.UpdateOpts{}) {
		log.Printf("[DEBUG] Updating sharenetwork %s with options: %#v", d.Id(), updateOpts)

		_, err = sharenetworks.Update(ctx, sfsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Unable to update sharenetwork %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("security_service_ids") {
		o, n := d.GetChange("security_service_ids")

		oldList, newList := o.(*schema.Set), n.(*schema.Set)
		newSecurityServiceIDs := newList.Difference(oldList)
		oldSecurityServiceIDs := oldList.Difference(newList)

		for _, newSecurityServiceID := range newSecurityServiceIDs.List() {
			id := newSecurityServiceID.(string)
			log.Printf("[DEBUG] Adding new %s security service to sharenetwork %s", id, d.Id())
			securityServiceOpts := sharenetworks.AddSecurityServiceOpts{SecurityServiceID: id}

			_, err = sharenetworks.AddSecurityService(ctx, sfsClient, d.Id(), securityServiceOpts).Extract()
			if err != nil {
				return diag.Errorf("Error adding new %s security service to sharenetwork: %s", id, err)
			}
		}

		for _, oldSecurityServiceID := range oldSecurityServiceIDs.List() {
			id := oldSecurityServiceID.(string)
			log.Printf("[DEBUG] Removing old %s security service from sharenetwork %s", id, d.Id())
			securityServiceOpts := sharenetworks.RemoveSecurityServiceOpts{SecurityServiceID: id}

			_, err = sharenetworks.RemoveSecurityService(ctx, sfsClient, d.Id(), securityServiceOpts).Extract()
			if err != nil {
				return diag.Errorf("Error removing old %s security service from sharenetwork: %s", id, err)
			}
		}
	}

	return resourceSharedFilesystemShareNetworkV2Read(ctx, d, meta)
}

func resourceSharedFilesystemShareNetworkV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete sharenetwork %s", d.Id())

	err = sharenetworks.Delete(ctx, sfsClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting sharenetwork"))
	}

	return nil
}

func resourceSharedFilesystemShareNetworkV2GetSvcByShareNetID(ctx context.Context, sfsClient *gophercloud.ServiceClient, shareNetworkID string) ([]string, error) {
	securityServiceListOpts := securityservices.ListOpts{ShareNetworkID: shareNetworkID}

	securityServicePages, err := securityservices.List(sfsClient, securityServiceListOpts).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to list security services for sharenetwork %s: %w", shareNetworkID, err)
	}

	securityServiceList, err := securityservices.ExtractSecurityServices(securityServicePages)
	if err != nil {
		return nil, fmt.Errorf("Unable to extract security services for sharenetwork %s: %w", shareNetworkID, err)
	}

	log.Printf("[DEBUG] Retrieved security services for sharenetwork %s: %#v", shareNetworkID, securityServiceList)

	return resourceSharedFilesystemShareNetworkV2SecSvcToArray(&securityServiceList), nil
}

func resourceSharedFilesystemShareNetworkV2SecSvcToArray(v any) []string {
	var securityServicesIDs []string

	switch t := v.(type) {
	case *schema.Set:
		for _, securityService := range (*v.(*schema.Set)).List() {
			securityServicesIDs = append(securityServicesIDs, securityService.(string))
		}
	case *[]securityservices.SecurityService:
		for _, securityService := range *v.(*[]securityservices.SecurityService) {
			securityServicesIDs = append(securityServicesIDs, securityService.ID)
		}
	default:
		log.Printf("[DEBUG] Invalid type provided to get the list of security service IDs: %s", t)
	}

	return securityServicesIDs
}
