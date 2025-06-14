package openstack

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkingSubnetV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSubnetV2Create,
		ReadContext:   resourceNetworkingSubnetV2Read,
		UpdateContext: resourceNetworkingSubnetV2Update,
		DeleteContext: resourceNetworkingSubnetV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cidr": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"prefix_length"},
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},

			"prefix_length": {
				Type:          schema.TypeInt,
				ConflictsWith: []string{"cidr"},
				Optional:      true,
				ForceNew:      true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"allocation_pool": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"gateway_ip": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"no_gateway"},
				Optional:      true,
				ForceNew:      false,
				Computed:      true,
			},

			"no_gateway": {
				Type:          schema.TypeBool,
				ConflictsWith: []string{"gateway_ip"},
				Optional:      true,
				Default:       false,
				ForceNew:      false,
			},

			"ip_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      4,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{4, 6}),
			},

			"enable_dhcp": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},

			"dns_nameservers": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"dns_publish_fixed_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"ipv6_address_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"slaac", "dhcpv6-stateful", "dhcpv6-stateless",
				}, false),
			},

			"ipv6_ra_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"slaac", "dhcpv6-stateful", "dhcpv6-stateless",
				}, false),
			},

			"segment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"subnetpool_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"service_types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNetworkingSubnetV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Check nameservers.
	if err := networkingSubnetV2DNSNameserverAreUnique(d.Get("dns_nameservers").([]any)); err != nil {
		return diag.Errorf("openstack_networking_subnet_v2 dns_nameservers argument is invalid: %s", err)
	}

	// Get raw allocation pool value.
	allocationPool := d.Get("allocation_pool").(*schema.Set).List()

	// Set basic options.
	createOpts := SubnetCreateOpts{
		subnets.CreateOpts{
			NetworkID:       d.Get("network_id").(string),
			Name:            d.Get("name").(string),
			Description:     d.Get("description").(string),
			TenantID:        d.Get("tenant_id").(string),
			IPv6AddressMode: d.Get("ipv6_address_mode").(string),
			IPv6RAMode:      d.Get("ipv6_ra_mode").(string),
			AllocationPools: expandNetworkingSubnetV2AllocationPools(allocationPool),
			DNSNameservers:  expandToStringSlice(d.Get("dns_nameservers").([]any)),
			ServiceTypes:    expandToStringSlice(d.Get("service_types").([]any)),
			SegmentID:       d.Get("segment_id").(string),
			SubnetPoolID:    d.Get("subnetpool_id").(string),
			IPVersion:       gophercloud.IPVersion(d.Get("ip_version").(int)),
		},
		MapValueSpecs(d),
	}

	if v, ok := d.GetOk("dns_publish_fixed_ip"); ok {
		v := v.(bool)
		createOpts.DNSPublishFixedIP = &v
	}

	// Set CIDR if provided. Check if inferred subnet would match the provided cidr.
	if v, ok := d.GetOk("cidr"); ok {
		cidr := v.(string)
		_, netAddr, err := net.ParseCIDR(cidr)
		if err != nil {
			return diag.Errorf("Invalid CIDR %s: %s", cidr, err)
		}

		if netAddr.String() != cidr {
			return diag.Errorf("cidr %s doesn't match subnet address %s for openstack_networking_subnet_v2", cidr, netAddr.String())
		}

		createOpts.CIDR = cidr
	}

	// Set gateway options if provided.
	if v, ok := d.GetOk("gateway_ip"); ok {
		gatewayIP := v.(string)
		createOpts.GatewayIP = &gatewayIP
	}

	noGateway := d.Get("no_gateway").(bool)
	if noGateway {
		gatewayIP := ""
		createOpts.GatewayIP = &gatewayIP
	}

	// Validate and set prefix options.
	if v, ok := d.GetOk("prefix_length"); ok {
		if d.Get("subnetpool_id").(string) == "" {
			return diag.Errorf("'prefix_length' is only valid if 'subnetpool_id' is set for openstack_networking_subnet_v2")
		}

		prefixLength := v.(int)
		createOpts.Prefixlen = prefixLength
	}

	// Set DHCP options if provided.
	enableDHCP := d.Get("enable_dhcp").(bool)
	createOpts.EnableDHCP = &enableDHCP

	log.Printf("[DEBUG] openstack_networking_subnet_v2 create options: %#v", createOpts)

	s, err := subnets.Create(ctx, networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_subnet_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_subnet_v2 %s to become available", s.ID)
	stateConf := &retry.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    networkingSubnetV2StateRefreshFunc(ctx, networkingClient, s.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_subnet_v2 %s to become available: %s", s.ID, err)
	}

	d.SetId(s.ID)

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}

		tags, err := attributestags.ReplaceAll(ctx, networkingClient, "subnets", s.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error creating tags on openstack_networking_subnet_v2 %s: %s", s.ID, err)
		}

		log.Printf("[DEBUG] Set tags %s on openstack_networking_subnet_v2 %s", tags, s.ID)
	}

	log.Printf("[DEBUG] Created openstack_networking_subnet_v2 %s: %#v", s.ID, s)

	return resourceNetworkingSubnetV2Read(ctx, d, meta)
}

func resourceNetworkingSubnetV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	s, err := subnets.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_subnet_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_subnet_v2 %s: %#v", d.Id(), s)

	d.Set("network_id", s.NetworkID)
	d.Set("cidr", s.CIDR)
	d.Set("ip_version", s.IPVersion)
	d.Set("name", s.Name)
	d.Set("description", s.Description)
	d.Set("tenant_id", s.TenantID)
	d.Set("dns_nameservers", s.DNSNameservers)
	d.Set("service_types", s.ServiceTypes)
	d.Set("enable_dhcp", s.EnableDHCP)
	d.Set("network_id", s.NetworkID)
	d.Set("ipv6_address_mode", s.IPv6AddressMode)
	d.Set("ipv6_ra_mode", s.IPv6RAMode)
	d.Set("segment_id", s.SegmentID)
	d.Set("subnetpool_id", s.SubnetPoolID)
	d.Set("dns_publish_fixed_ip", s.DNSPublishFixedIP)

	networkingV2ReadAttributesTags(d, s.Tags)

	// Set the allocation_pool attribute
	allocationPools := flattenNetworkingSubnetV2AllocationPools(s.AllocationPools)
	if err := d.Set("allocation_pool", allocationPools); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_subnet_v2 allocation_pool: %s", err)
	}

	// Set the subnet's "gateway_ip" and "no_gateway" attributes.
	d.Set("gateway_ip", s.GatewayIP)
	d.Set("no_gateway", false)

	if s.GatewayIP != "" {
		d.Set("no_gateway", false)
	} else {
		d.Set("no_gateway", true)
	}

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingSubnetV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var hasChange bool

	var updateOpts subnets.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("gateway_ip") {
		hasChange = true
		updateOpts.GatewayIP = nil

		if v, ok := d.GetOk("gateway_ip"); ok {
			gatewayIP := v.(string)
			updateOpts.GatewayIP = &gatewayIP
		}
	}

	if d.HasChange("no_gateway") {
		if d.Get("no_gateway").(bool) {
			hasChange = true
			gatewayIP := ""
			updateOpts.GatewayIP = &gatewayIP
		}
	}

	if d.HasChange("dns_nameservers") {
		if err := networkingSubnetV2DNSNameserverAreUnique(d.Get("dns_nameservers").([]any)); err != nil {
			return diag.Errorf("openstack_networking_subnet_v2 dns_nameservers argument is invalid: %s", err)
		}

		hasChange = true
		nameservers := expandToStringSlice(d.Get("dns_nameservers").([]any))
		updateOpts.DNSNameservers = &nameservers
	}

	if d.HasChange("service_types") {
		hasChange = true
		serviceTypes := expandToStringSlice(d.Get("service_types").([]any))
		updateOpts.ServiceTypes = &serviceTypes
	}

	if d.HasChange("enable_dhcp") {
		hasChange = true
		v := d.Get("enable_dhcp").(bool)
		updateOpts.EnableDHCP = &v
	}

	if d.HasChange("allocation_pool") {
		hasChange = true
		updateOpts.AllocationPools = expandNetworkingSubnetV2AllocationPools(d.Get("allocation_pool").(*schema.Set).List())
	}

	if d.HasChange("segment_id") {
		hasChange = true
		v := d.Get("segment_id").(string)
		updateOpts.SegmentID = &v
	}

	if d.HasChange("dns_publish_fixed_ip") {
		hasChange = true
		v := d.Get("dns_publish_fixed_ip").(bool)
		updateOpts.DNSPublishFixedIP = &v
	}

	if hasChange {
		log.Printf("[DEBUG] Updating openstack_networking_subnet_v2 %s with options: %#v", d.Id(), updateOpts)

		_, err = subnets.Update(ctx, networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating OpenStack Neutron openstack_networking_subnet_v2 %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		tags := networkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}

		tags, err := attributestags.ReplaceAll(ctx, networkingClient, "subnets", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating tags on openstack_networking_subnet_v2 %s: %s", d.Id(), err)
		}

		log.Printf("[DEBUG] Updated tags %s on openstack_networking_subnet_v2 %s", tags, d.Id())
	}

	return resourceNetworkingSubnetV2Read(ctx, d, meta)
}

func resourceNetworkingSubnetV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingSubnetV2StateRefreshFuncDelete(ctx, networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_subnet_v2 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}
