package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharenetworks"
)

func resourceSharedfilesystemSharenetworkV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceSharedfilesystemSharenetworkV2Create,
		Read:   resourceSharedfilesystemSharenetworkV2Read,
		Update: resourceSharedfilesystemSharenetworkV2Update,
		Delete: resourceSharedfilesystemSharenetworkV2Delete,
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
				Optional: true,
				Computed: true,
				ForceNew: true,
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

func resourceSharedfilesystemSharenetworkV2Create(d *schema.ResourceData, meta interface{}) error {
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

	securityServiceIDs := resourceSharedfilesystemSharenetworkSecurityServicesV2(d.Get("security_service_ids").(*schema.Set))
	if len(securityServiceIDs) > 0 {
		for _, securityServiceID := range securityServiceIDs {
			log.Printf("[DEBUG] Adding %s security service to sharenetwork %s", securityServiceID, sharenetwork.ID)
			securityServiceOpts := sharenetworks.AddSecurityServiceOpts{SecurityServiceID: securityServiceID}
			_, err = sharenetworks.AddSecurityService(sfsClient, sharenetwork.ID, securityServiceOpts).Extract()
			if err != nil {
				return fmt.Errorf("Error adding %s security service to sharenetwork: %s", securityServiceID, err)
			}
		}
	}

	return resourceSharedfilesystemSharenetworkV2Read(d, meta)
}

func resourceSharedfilesystemSharenetworkV2Read(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	securityServiceList, err := securityservices.ExtractSecurityServices(securityServicePages)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Retrieved security services for sharenetwork %s: %#v", d.Id(), securityServiceList)

	if len(securityServiceList) > 0 {
		d.Set("security_service_ids", resourceSharedfilesystemSharenetworkSecurityServices2IDsV2(&securityServiceList))
	} else {
		d.Set("security_service_ids", []string{})
	}

	d.Set("name", sharenetwork.Name)
	d.Set("description", sharenetwork.Description)
	d.Set("neutron_net_id", sharenetwork.NeutronNetID)
	d.Set("neutron_subnet_id", sharenetwork.NeutronSubnetID)
	d.Set("project_id", sharenetwork.ProjectID)
	d.Set("region", GetRegion(d, config))
	// Computed
	d.Set("network_type", sharenetwork.NetworkType)
	d.Set("segmentation_id", sharenetwork.SegmentationID)
	d.Set("cidr", sharenetwork.CIDR)
	d.Set("ip_version", sharenetwork.IPVersion)

	return nil
}

func resourceSharedfilesystemSharenetworkV2Update(d *schema.ResourceData, meta interface{}) error {
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
		oldSecurityServiceIDs := resourceSharedfilesystemSharenetworkSecurityServicesV2(old.(*schema.Set))
		newSecurityServiceIDs := resourceSharedfilesystemSharenetworkSecurityServicesV2(new.(*schema.Set))
		for _, newSecurityServiceID := range newSecurityServiceIDs {
			if !inArray(newSecurityServiceID, &oldSecurityServiceIDs) {
				log.Printf("[DEBUG] Adding new %s security service to sharenetwork %s", newSecurityServiceID, d.Id())
				securityServiceOpts := sharenetworks.AddSecurityServiceOpts{SecurityServiceID: newSecurityServiceID}
				_, err = sharenetworks.AddSecurityService(sfsClient, d.Id(), securityServiceOpts).Extract()
				if err != nil {
					return fmt.Errorf("Error adding new %s security service to sharenetwork: %s", newSecurityServiceID, err)
				}
			}
		}
		for _, oldSecurityServiceID := range oldSecurityServiceIDs {
			if !inArray(oldSecurityServiceID, &newSecurityServiceIDs) {
				log.Printf("[DEBUG] Removing old %s security service from sharenetwork %s", oldSecurityServiceID, d.Id())
				securityServiceOpts := sharenetworks.RemoveSecurityServiceOpts{SecurityServiceID: oldSecurityServiceID}
				_, err = sharenetworks.RemoveSecurityService(sfsClient, d.Id(), securityServiceOpts).Extract()
				if err != nil {
					return fmt.Errorf("Error removing old %s security service from sharenetwork: %s", oldSecurityServiceID, err)
				}
			}
		}
	}

	return resourceSharedfilesystemSharenetworkV2Read(d, meta)
}

func resourceSharedfilesystemSharenetworkV2Delete(d *schema.ResourceData, meta interface{}) error {
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

	d.SetId("")

	return nil
}

func resourceSharedfilesystemSharenetworkSecurityServicesV2(v *schema.Set) []string {
	var securityServices []string
	for _, v := range v.List() {
		securityServices = append(securityServices, v.(string))
	}
	return securityServices
}

func resourceSharedfilesystemSharenetworkSecurityServices2IDsV2(v *[]securityservices.SecurityService) []string {
	var securityServicesIDs []string
	for _, securityService := range *v {
		securityServicesIDs = append(securityServicesIDs, securityService.ID)
	}
	return securityServicesIDs
}

func inArray(a string, s *[]string) bool {
	for _, b := range *s {
		if a == b {
			return true
		}
	}
	return false
}
