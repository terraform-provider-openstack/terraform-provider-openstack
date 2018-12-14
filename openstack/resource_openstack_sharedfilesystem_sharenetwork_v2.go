package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharenetworks"
)

func resourceSharedFilesystemShareNetworkV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceSharedFilesystemShareNetworkV2Create,
		Read:   resourceSharedFilesystemShareNetworkV2Read,
		Update: resourceSharedFilesystemShareNetworkV2Update,
		Delete: resourceSharedFilesystemShareNetworkV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"neutron_net_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"neutron_subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"security_service_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"network_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"segmentation_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSharedFilesystemShareNetworkV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	createOpts := sharenetworks.CreateOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		NeutronNetID:    d.Get("neutron_net_id").(string),
		NeutronSubnetID: d.Get("neutron_subnet_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	log.Printf("[DEBUG] Attempting to create sharenetwork")
	sharenetwork, err := sharenetworks.Create(sfsClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating sharenetwork: %s", err)
	}

	d.SetId(sharenetwork.ID)

	securityServiceIDs := resourceSharedFilesystemShareNetworkV2SecSvcToArray(d.Get("security_service_ids").(*schema.Set))
	for _, securityServiceID := range securityServiceIDs {
		log.Printf("[DEBUG] Adding %s security service to sharenetwork %s", securityServiceID, sharenetwork.ID)
		securityServiceOpts := sharenetworks.AddSecurityServiceOpts{SecurityServiceID: securityServiceID}
		_, err = sharenetworks.AddSecurityService(sfsClient, sharenetwork.ID, securityServiceOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error adding %s security service to sharenetwork: %s", securityServiceID, err)
		}
	}

	return resourceSharedFilesystemShareNetworkV2Read(d, meta)
}

func resourceSharedFilesystemShareNetworkV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sharenetwork, err := sharenetworks.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "sharenetwork")
	}

	log.Printf("[DEBUG] Retrieved sharenetwork %s: %#v", d.Id(), sharenetwork)

	securityServiceListOpts := securityservices.ListOpts{ShareNetworkID: d.Id()}
	securityServicePages, err := securityservices.List(sfsClient, securityServiceListOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to list security services for sharenetwork %s: %s", d.Id(), err)
	}
	securityServiceList, err := securityservices.ExtractSecurityServices(securityServicePages)
	if err != nil {
		return fmt.Errorf("Unable to extract security services for sharenetwork %s: %s", d.Id(), err)
	}
	log.Printf("[DEBUG] Retrieved security services for sharenetwork %s: %#v", d.Id(), securityServiceList)

	d.Set("security_service_ids", resourceSharedFilesystemShareNetworkV2SecSvcToArray(&securityServiceList))

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

func resourceSharedFilesystemShareNetworkV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
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
		_, err = sharenetworks.Update(sfsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Unable to update sharenetwork %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("security_service_ids") {
		old, new := d.GetChange("security_service_ids")

		oldList, newList := old.(*schema.Set), new.(*schema.Set)
		newSecurityServiceIDs := newList.Difference(oldList)
		oldSecurityServiceIDs := oldList.Difference(newList)

		for _, newSecurityServiceID := range newSecurityServiceIDs.List() {
			log.Printf("[DEBUG] Adding new %s security service to sharenetwork %s", newSecurityServiceID, d.Id())
			securityServiceOpts := sharenetworks.AddSecurityServiceOpts{SecurityServiceID: newSecurityServiceID.(string)}
			_, err = sharenetworks.AddSecurityService(sfsClient, d.Id(), securityServiceOpts).Extract()
			if err != nil {
				return fmt.Errorf("Error adding new %s security service to sharenetwork: %s", newSecurityServiceID, err)
			}
		}
		for _, oldSecurityServiceID := range oldSecurityServiceIDs.List() {
			log.Printf("[DEBUG] Removing old %s security service from sharenetwork %s", oldSecurityServiceID, d.Id())
			securityServiceOpts := sharenetworks.RemoveSecurityServiceOpts{SecurityServiceID: oldSecurityServiceID.(string)}
			_, err = sharenetworks.RemoveSecurityService(sfsClient, d.Id(), securityServiceOpts).Extract()
			if err != nil {
				return fmt.Errorf("Error removing old %s security service from sharenetwork: %s", oldSecurityServiceID, err)
			}
		}
	}

	return resourceSharedFilesystemShareNetworkV2Read(d, meta)
}

func resourceSharedFilesystemShareNetworkV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete sharenetwork %s", d.Id())
	err = sharenetworks.Delete(sfsClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting sharenetwork: %s", err)
	}

	return nil
}

func resourceSharedFilesystemShareNetworkV2SecSvcToArray(v interface{}) []string {
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
