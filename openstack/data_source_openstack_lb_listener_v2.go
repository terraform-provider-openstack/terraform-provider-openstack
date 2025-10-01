package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLBListenerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBListenerV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"listener_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name"},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"listener_id"},
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"protocol_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_pool": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"loadbalancers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"connection_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"sni_container_refs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"pools": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"l7policies": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"timeout_client_data": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"timeout_member_data": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"timeout_member_connect": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"timeout_tcp_inspect": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"insert_headers": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"allowed_cidrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tls_ciphers": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tls_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"alpn_protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"client_authentication": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"client_ca_tls_container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"client_crl_container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hsts_include_subdomains": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"hsts_max_age": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"hsts_preload": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLBListenerV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancer client: %s", err)
	}

	listOpts := listeners.ListOpts{
		Tags: expandTagsList(d, "tags"),
	}

	if v, ok := d.GetOk("listener_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("loadbalancer_id"); ok {
		listOpts.LoadbalancerID = v.(string)
	}

	if v, ok := d.GetOk("protocol"); ok {
		listOpts.Protocol = v.(string)
	}

	if v, ok := d.GetOk("protocol_port"); ok {
		listOpts.ProtocolPort = v.(int)
	}

	allPages, err := listeners.List(lbClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query OpenStack loadbalancer listeners: %s", err)
	}

	allListeners, err := listeners.ExtractListeners(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve Openstack loadbalancer listeners: %s", err)
	}

	if len(allListeners) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allListeners) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allListeners)

		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria")
	}

	dataSourceLBListenerV2Attributes(d, &allListeners[0])
	d.Set("region", GetRegion(d, config))

	return nil
}

func dataSourceLBListenerV2Attributes(d *schema.ResourceData, listener *listeners.Listener) {
	log.Printf("[DEBUG] Retrieved openstack_lb_listener_v2 %s: %#v", listener.ID, listener)

	d.SetId(listener.ID)
	d.Set("project_id", listener.ProjectID)
	d.Set("name", listener.Name)
	d.Set("description", listener.Description)
	d.Set("protocol", listener.Protocol)
	d.Set("protocol_port", listener.ProtocolPort)
	d.Set("default_pool_id", listener.DefaultPoolID)
	d.Set("default_pool", listener.DefaultPool)
	d.Set("loadbalancers", flattenLBListenerLoadbalancerIDsV2(listener.Loadbalancers))
	d.Set("connection_limit", listener.ConnLimit)
	d.Set("sni_container_refs", listener.SniContainerRefs)
	d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef)
	d.Set("admin_state_up", listener.AdminStateUp)
	d.Set("pools", flattenLBPoolsV2(listener.Pools))
	d.Set("l7policies", flattenLBPoliciesV2(listener.L7Policies))
	d.Set("provisioning_status", listener.ProvisioningStatus)
	d.Set("timeout_client_data", listener.TimeoutClientData)
	d.Set("timeout_member_data", listener.TimeoutMemberData)
	d.Set("timeout_member_connect", listener.TimeoutMemberConnect)
	d.Set("timeout_tcp_inspect", listener.TimeoutTCPInspect)
	d.Set("insert_headers", listener.InsertHeaders)
	d.Set("allowed_cidrs", listener.AllowedCIDRs)
	d.Set("tls_ciphers", listener.TLSCiphers)
	d.Set("tls_versions", listener.TLSVersions)
	d.Set("tags", listener.Tags)
	d.Set("alpn_protocols", listener.ALPNProtocols)
	d.Set("client_authentication", listener.ClientAuthentication)
	d.Set("client_ca_tls_container_ref", listener.ClientCATLSContainerRef)
	d.Set("client_crl_container_ref", listener.ClientCRLContainerRef)
	d.Set("hsts_include_subdomains", listener.HSTSIncludeSubdomains)
	d.Set("hsts_max_age", listener.HSTSMaxAge)
	d.Set("hsts_preload", listener.HSTSPreload)
	d.Set("operating_status", listener.OperatingStatus)
}
