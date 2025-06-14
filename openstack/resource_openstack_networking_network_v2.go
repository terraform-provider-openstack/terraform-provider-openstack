package openstack

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/external"
	mtuext "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/mtu"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/provider"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vlantransparent"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkingNetworkV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingNetworkV2Create,
		ReadContext:   resourceNetworkingNetworkV2Read,
		UpdateContext: resourceNetworkingNetworkV2Update,
		DeleteContext: resourceNetworkingNetworkV2Delete,
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
				ForceNew: true,
				Computed: true,
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

			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"external": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"segments": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"physical_network": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"segmentation_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
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

			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"availability_zone_hints": {
				Type:     schema.TypeSet,
				Computed: true,
				ForceNew: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"transparent_vlan": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"port_security_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"mtu": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"dns_domain": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^$|\.$`), "fully-qualified (unambiguous) DNS domain names must have a dot at the end"),
			},

			"qos_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingNetworkV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	azHints := d.Get("availability_zone_hints").(*schema.Set)

	createOpts := NetworkCreateOpts{
		networks.CreateOpts{
			Name:                  d.Get("name").(string),
			Description:           d.Get("description").(string),
			TenantID:              d.Get("tenant_id").(string),
			AvailabilityZoneHints: expandToStringSlice(azHints.List()),
		},
		MapValueSpecs(d),
	}

	if v, ok := getOkExists(d, "admin_state_up"); ok {
		asu := v.(bool)
		createOpts.AdminStateUp = &asu
	}

	if v, ok := getOkExists(d, "shared"); ok {
		shared := v.(bool)
		createOpts.Shared = &shared
	}

	// Declare a finalCreateOpts interface.
	var finalCreateOpts networks.CreateOptsBuilder
	finalCreateOpts = createOpts

	// Add networking segments if specified.
	segments := expandNetworkingNetworkSegmentsV2(d.Get("segments").(*schema.Set))
	if len(segments) > 0 {
		finalCreateOpts = provider.CreateOptsExt{
			CreateOptsBuilder: finalCreateOpts,
			Segments:          segments,
		}
	}

	// Add the external attribute if specified.
	isExternal := d.Get("external").(bool)
	if isExternal {
		finalCreateOpts = external.CreateOptsExt{
			CreateOptsBuilder: finalCreateOpts,
			External:          &isExternal,
		}
	}

	// Add the transparent VLAN attribute if specified.
	isVLANTransparent := d.Get("transparent_vlan").(bool)
	if isVLANTransparent {
		finalCreateOpts = vlantransparent.CreateOptsExt{
			CreateOptsBuilder: finalCreateOpts,
			VLANTransparent:   &isVLANTransparent,
		}
	}

	// Add the port security attribute if specified.
	if v, ok := getOkExists(d, "port_security_enabled"); ok {
		portSecurityEnabled := v.(bool)
		finalCreateOpts = portsecurity.NetworkCreateOptsExt{
			CreateOptsBuilder:   finalCreateOpts,
			PortSecurityEnabled: &portSecurityEnabled,
		}
	}

	mtu := d.Get("mtu").(int)
	// Add the MTU attribute if specified.
	if mtu > 0 {
		finalCreateOpts = mtuext.CreateOptsExt{
			CreateOptsBuilder: finalCreateOpts,
			MTU:               mtu,
		}
	}

	// Add the DNS Domain attribute if specified.
	if dnsDomain := d.Get("dns_domain").(string); dnsDomain != "" {
		finalCreateOpts = dns.NetworkCreateOptsExt{
			CreateOptsBuilder: finalCreateOpts,
			DNSDomain:         dnsDomain,
		}
	}

	// Add the QoS policy ID attribute if specified.
	if qosPolicyID := d.Get("qos_policy_id").(string); qosPolicyID != "" {
		finalCreateOpts = policies.NetworkCreateOptsExt{
			CreateOptsBuilder: finalCreateOpts,
			QoSPolicyID:       qosPolicyID,
		}
	}

	log.Printf("[DEBUG] openstack_networking_network_v2 create options: %#v", finalCreateOpts)

	n, err := networks.Create(ctx, networkingClient, finalCreateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_network_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_network_v2 %s to become available.", n.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE", "DOWN"},
		Refresh:    resourceNetworkingNetworkV2StateRefreshFunc(ctx, networkingClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_network_v2 %s to become available: %s", n.ID, err)
	}

	d.SetId(n.ID)

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}

		tags, err := attributestags.ReplaceAll(ctx, networkingClient, "networks", n.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_network_v2 %s: %s", n.ID, err)
		}

		log.Printf("[DEBUG] Set tags %s on openstack_networking_network_v2 %s", tags, n.ID)
	}

	log.Printf("[DEBUG] Created openstack_networking_network_v2 %s: %#v", n.ID, n)

	return resourceNetworkingNetworkV2Read(ctx, d, meta)
}

func resourceNetworkingNetworkV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var network networkExtended

	err = networks.Get(ctx, networkingClient, d.Id()).ExtractInto(&network)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_network_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_network_v2 %s: %#v", d.Id(), network)

	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("admin_state_up", network.AdminStateUp)
	d.Set("shared", network.Shared)
	d.Set("external", network.External)
	d.Set("tenant_id", network.TenantID)
	d.Set("segments", flattenNetworkingNetworkSegmentsV2(network))
	d.Set("transparent_vlan", network.VLANTransparent)
	d.Set("port_security_enabled", network.PortSecurityEnabled)
	d.Set("mtu", network.MTU)
	d.Set("dns_domain", network.DNSDomain)
	d.Set("qos_policy_id", network.QoSPolicyID)
	d.Set("region", GetRegion(d, config))

	networkingV2ReadAttributesTags(d, network.Tags)

	if err := d.Set("availability_zone_hints", network.AvailabilityZoneHints); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_network_v2 %s availability_zone_hints: %s", d.Id(), err)
	}

	return nil
}

func resourceNetworkingNetworkV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Declare finalUpdateOpts interface and basic updateOpts structure.
	var (
		finalUpdateOpts networks.UpdateOptsBuilder
		updateOpts      networks.UpdateOpts
	)

	// Populate basic updateOpts.
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	if d.HasChange("shared") {
		shared := d.Get("shared").(bool)
		updateOpts.Shared = &shared
	}

	// Change tags if needed.
	if d.HasChange("tags") {
		tags := networkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}

		tags, err := attributestags.ReplaceAll(ctx, networkingClient, "networks", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_network_v2 %s: %s", d.Id(), err)
		}

		log.Printf("[DEBUG] Set tags %s on openstack_networking_network_v2 %s", tags, d.Id())
	}

	// Save basic updateOpts into finalUpdateOpts.
	finalUpdateOpts = updateOpts

	// Populate extensions options.
	if d.HasChange("external") {
		isExternal := d.Get("external").(bool)
		finalUpdateOpts = external.UpdateOptsExt{
			UpdateOptsBuilder: finalUpdateOpts,
			External:          &isExternal,
		}
	}

	// Populate port security options.
	if d.HasChange("port_security_enabled") {
		portSecurityEnabled := d.Get("port_security_enabled").(bool)
		finalUpdateOpts = portsecurity.NetworkUpdateOptsExt{
			UpdateOptsBuilder:   finalUpdateOpts,
			PortSecurityEnabled: &portSecurityEnabled,
		}
	}

	if d.HasChange("mtu") {
		mtu := d.Get("mtu").(int)
		finalUpdateOpts = mtuext.UpdateOptsExt{
			UpdateOptsBuilder: finalUpdateOpts,
			MTU:               mtu,
		}
	}

	if d.HasChange("dns_domain") {
		dnsDomain := d.Get("dns_domain").(string)
		finalUpdateOpts = dns.NetworkUpdateOptsExt{
			UpdateOptsBuilder: finalUpdateOpts,
			DNSDomain:         &dnsDomain,
		}
	}

	if d.HasChange("qos_policy_id") {
		qosPolicyID := d.Get("qos_policy_id").(string)
		finalUpdateOpts = policies.NetworkUpdateOptsExt{
			UpdateOptsBuilder: finalUpdateOpts,
			QoSPolicyID:       &qosPolicyID,
		}
	}

	if d.HasChange("segments") {
		segments := expandNetworkingNetworkSegmentsV2(d.Get("segments").(*schema.Set))
		finalUpdateOpts = provider.UpdateOptsExt{
			UpdateOptsBuilder: finalUpdateOpts,
			Segments:          &segments,
		}
	}

	log.Printf("[DEBUG] openstack_networking_network_v2 %s update options: %#v", d.Id(), finalUpdateOpts)

	_, err = networks.Update(ctx, networkingClient, d.Id(), finalUpdateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_networking_network_v2 %s: %s", d.Id(), err)
	}

	return resourceNetworkingNetworkV2Read(ctx, d, meta)
}

func resourceNetworkingNetworkV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if err := networks.Delete(ctx, networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_network_v2"))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingNetworkV2StateRefreshFunc(ctx, networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_network_v2 %s to Delete:  %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}
