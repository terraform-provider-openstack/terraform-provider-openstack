package openstack

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/endpoints"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	InterfaceRegex = "admin|public|internal"
)

func resourceIdentityEndpointV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityEndpointV3Create,
		Read:   resourceIdentityEndpointV3Read,
		Update: resourceIdentityEndpointV3Update,
		Delete: resourceIdentityEndpointV3Delete,

		Schema: map[string]*schema.Schema{
			"interface": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateInterface(),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func validateInterface() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)

		if !regexp.MustCompile(InterfaceRegex).MatchString(value) {
			errors = append(errors, fmt.Errorf(
				"%q name must be of [admin|public|internal]", value))
		}
		return
	}
}

func resourceIdentityEndpointV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	availability := gophercloud.AvailabilityPublic
	if d.Get("interface").(string) == "admin" {
		availability = gophercloud.AvailabilityAdmin
	}
	if d.Get("interface").(string) == "internal" {
		availability = gophercloud.AvailabilityInternal
	}

	createOpts := endpoints.CreateOpts{
		Availability: availability,
		Name:         d.Get("name").(string),
		Region:       d.Get("region").(string),
		URL:          d.Get("url").(string),
		ServiceID:    d.Get("service_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	endpoint, err := endpoints.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenStack endpoint: %s", err)
	}

	d.SetId(endpoint.ID)

	return resourceIdentityEndpointV3Read(d, meta)
}

func resourceIdentityEndpointV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var found endpoints.Endpoint
	err = endpoints.List(identityClient, endpoints.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		if endpointList, err := endpoints.ExtractEndpoints(page); err != nil {
			return false, err
		} else {
			for _, endpoint := range endpointList {
				if endpoint.ID == d.Id() {
					found = endpoint
					break
				}
			}
		}
		return true, nil
	})

	if err != nil {
		return CheckDeleted(d, err, "endpoint")
	}

	log.Printf("[DEBUG] Retrieved OpenStack endpoint: %#v", found)

	d.Set("name", found.Name)
	d.Set("region", found.Region)
	d.Set("url", found.URL)
	d.Set("service_id", found.ServiceID)
	d.Set("interface", string(found.Availability))

	return nil
}

func resourceIdentityEndpointV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var hasChange bool
	var updateOpts endpoints.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("region") {
		hasChange = true
		updateOpts.Region = d.Get("region").(string)
	}

	if d.HasChange("url") {
		hasChange = true
		updateOpts.URL = d.Get("url").(string)
	}

	if d.HasChange("service_id") {
		hasChange = true
		updateOpts.ServiceID = d.Get("service_id").(string)
	}

	if d.HasChange("interface") {
		hasChange = true

		availability := gophercloud.AvailabilityPublic
		if d.Get("interface").(string) == "admin" {
			availability = gophercloud.AvailabilityAdmin
		}
		if d.Get("interface").(string) == "internal" {
			availability = gophercloud.AvailabilityInternal
		}
		updateOpts.Availability = availability
	}

	if hasChange {
		_, err := endpoints.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenStack endpoint: %s", err)
		}
	}

	return resourceIdentityEndpointV3Read(d, meta)
}

func resourceIdentityEndpointV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	err = endpoints.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenStack endpoint: %s", err)
	}

	d.SetId("")
	return nil
}
