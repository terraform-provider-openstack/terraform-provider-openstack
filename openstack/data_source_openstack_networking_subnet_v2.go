package openstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func dataSourceNetworkingSubnetV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkingSubnetV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: descriptions["tenant_id"],
			},
			"gateway_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkingSubnetV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))

	listOpts := subnets.ListOpts{
		ID:        d.Get("subnet_id").(string),
		NetworkID: d.Get("network_id").(string),
		Name:      d.Get("name").(string),
		CIDR:      d.Get("cidr").(string),
		TenantID:  d.Get("tenant_id").(string),
		Limit:     2,
	}

	pages, err := subnets.List(networkingClient, listOpts).AllPages()
	allSubnets, err := subnets.ExtractSubnets(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve subnets: %s", err)
	}

	if len(allSubnets) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allSubnets) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	subnet := allSubnets[0]

	log.Printf("[DEBUG] Retrieved Subnet %s: %+v", subnet.ID, subnet)

	d.SetId(subnet.ID)
	d.Set("network_id", subnet.NetworkID)
	d.Set("cidr", subnet.CIDR)
	d.Set("ip_version", subnet.IPVersion)
	d.Set("name", subnet.Name)
	d.Set("tenant_id", subnet.TenantID)
	d.Set("gateway_ip", subnet.GatewayIP)
	d.Set("region", GetRegion(d, config))
	return nil
}
