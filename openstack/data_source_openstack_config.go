package openstack

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOpenStackConfig() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOpenStackConfigRead,
		Schema: getProviderSchema(),
	}
}

// dataSourceOpenStackConfigRead performs the endpoint lookup.
func dataSourceOpenStackConfigRead(d *schema.ResourceData, meta interface{}) error {
	config, err := configureProvider(d)
	if err != nil {
		return err
	}

	(config.(*Config)).LoadAndValidate()
	d.SetId(d.Get("auth_url").(string))
	return nil
}
