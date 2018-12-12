package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
)

func resourceSharedfilesystemSecurityserviceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceSharedfilesystemSecurityserviceV2Create,
		Read:   resourceSharedfilesystemSecurityserviceV2Read,
		Update: resourceSharedfilesystemSecurityserviceV2Update,
		Delete: resourceSharedfilesystemSecurityserviceV2Delete,
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

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "active_directory" && value != "kerberos" && value != "ldap" {
						errors = append(errors, fmt.Errorf(
							"Only 'active_directory', 'kerberos' and 'ldap' are supported values for 'type'"))
					}
					return
				},
			},

			"dns_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"ou": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"user": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"server": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSharedfilesystemSecurityserviceV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	mvSet, ouErr := setManilaMicroversion(sfsClient)
	if !mvSet && ouErr != nil {
		return ouErr
	}

	createOpts := securityservices.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        securityservices.SecurityServiceType(d.Get("type").(string)),
		DNSIP:       d.Get("dns_ip").(string),
		User:        d.Get("user").(string),
		Domain:      d.Get("domain").(string),
		Server:      d.Get("server").(string),
	}

	ou := d.Get("ou").(string)
	if ou != "" {
		if ouErr == nil {
			createOpts.OU = ou
		} else {
			return ouErr
		}
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	createOpts.Password = d.Get("password").(string)
	securityservice, err := securityservices.Create(sfsClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating : %s", err)
	}

	d.SetId(securityservice.ID)

	return resourceSharedfilesystemSecurityserviceV2Read(d, meta)
}

func resourceSharedfilesystemSecurityserviceV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	mvSet, ouErr := setManilaMicroversion(sfsClient)
	if !mvSet && ouErr != nil {
		return ouErr
	}

	securityservice, err := securityservices.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "securityservice")
	}

	nopassword := securityservice
	nopassword.Password = ""
	log.Printf("[DEBUG] Retrieved securityservice %s: %#v", d.Id(), nopassword)

	d.Set("name", securityservice.Name)
	d.Set("description", securityservice.Description)
	d.Set("type", securityservice.Type)
	d.Set("project_id", securityservice.ProjectID)
	d.Set("domain", securityservice.Domain)
	d.Set("dns_ip", securityservice.DNSIP)
	d.Set("user", securityservice.User)
	d.Set("server", securityservice.Server)
	d.Set("region", GetRegion(d, config))

	if ouErr == nil {
		d.Set("ou", securityservice.OU)
	} else {
		d.Set("ou", "")
	}

	return nil
}

func resourceSharedfilesystemSecurityserviceV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	mvSet, ouErr := setManilaMicroversion(sfsClient)
	if !mvSet && ouErr != nil {
		return ouErr
	}

	var updateOpts securityservices.UpdateOpts
	// Name should always be sent, otherwise it is vanished by manila backend
	name := d.Get("name").(string)
	updateOpts.Name = &name
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("type") {
		updateOpts.Type = d.Get("type").(string)
	}
	if d.HasChange("dns_ip") {
		dnsIP := d.Get("dns_ip").(string)
		updateOpts.DNSIP = &dnsIP
	}
	if d.HasChange("ou") {
		if ouErr == nil {
			ou := d.Get("ou").(string)
			updateOpts.OU = &ou
		} else {
			return ouErr
		}
	}
	if d.HasChange("user") {
		user := d.Get("user").(string)
		updateOpts.User = &user
	}
	if d.HasChange("domain") {
		domain := d.Get("domain").(string)
		updateOpts.Domain = &domain
	}
	if d.HasChange("server") {
		server := d.Get("server").(string)
		updateOpts.Server = &server
	}

	log.Printf("[DEBUG] Updating securityservice %s with options: %#v", d.Id(), updateOpts)

	if d.HasChange("password") {
		password := d.Get("password").(string)
		updateOpts.Password = &password
	}

	_, err = securityservices.Update(sfsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Unable to update securityservice %s: %s", d.Id(), err)
	}

	return resourceSharedfilesystemSecurityserviceV2Read(d, meta)
}

func resourceSharedfilesystemSecurityserviceV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete securityservice %s", d.Id())
	err = securityservices.Delete(sfsClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting securityservice: %s", err)
	}

	d.SetId("")

	return nil
}
