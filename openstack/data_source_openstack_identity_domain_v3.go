package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIdentityDomainV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdentityDomainV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// dataSourceIdentityDomainV3Read performs the domain lookup.
func dataSourceIdentityDomainV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	listOpts := domains.ListOpts{
		Enabled: &enabled,
		Name:    d.Get("name").(string),
	}

	log.Printf("[DEBUG] openstack_identity_domain_v3 list options: %#v", listOpts)

	var domain domains.Domain
	allPages, err := domains.List(identityClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query openstack_identity_domain_v3: %s", err)
	}

	allDomains, err := domains.ExtractDomains(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve openstack_identity_domain_v3: %s", err)
	}

	if len(allDomains) < 1 {
		return fmt.Errorf("Your openstack_identity_domain_v3 query returned no results")
	}

	if len(allDomains) > 1 {
		return fmt.Errorf("Your openstack_identity_domain_v3 query returned more than one result")
	}

	domain = allDomains[0]

	return dataSourceIdentityDomainV3Attributes(d, config, &domain)
}

// dataSourceIdentityDomainV3Attributes populates the fields of an Domain resource.
func dataSourceIdentityDomainV3Attributes(d *schema.ResourceData, config *Config, domain *domains.Domain) error {
	log.Printf("[DEBUG] openstack_identity_domain_v3 details: %#v", domain)

	d.SetId(domain.ID)
	d.Set("name", domain.Name)
	d.Set("domain_id", domain.ID)
	d.Set("region", GetRegion(d, config))
	d.Set("enabled", domain.Enabled)

	return nil
}
