package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLBPoolV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBPoolV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pool_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name"},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"pool_id"},
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

			"lb_method": {
				Type:     schema.TypeString,
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

			"listeners": {
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

			"members": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"admin_state_up": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"protocol_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"provisioning_status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"operating_status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"monitor_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"monitor_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"tags": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"healthmonitor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
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

			"session_persistence": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cookie_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"alpn_protocols": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ca_tls_container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"crl_container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tls_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tls_ciphers": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tls_container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tls_versions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLBPoolV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancer client: %s", err)
	}

	listOpts := pools.ListOpts{
		Tags: expandTagsList(d, "tags"),
	}

	if v, ok := d.GetOk("pool_id"); ok {
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

	if v, ok := d.GetOk("lb_method"); ok {
		listOpts.LBMethod = v.(string)
	}

	allPages, err := pools.List(lbClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query OpenStack loadbalancer pools: %s", err)
	}

	allPools, err := pools.ExtractPools(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve Openstack loadbalancer pools: %s", err)
	}

	if len(allPools) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allPools) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allPools)

		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria")
	}

	dataSourceLBPoolV2Attributes(d, &allPools[0])
	d.Set("region", GetRegion(d, config))

	return nil
}

func dataSourceLBPoolV2Attributes(d *schema.ResourceData, pool *pools.Pool) {
	log.Printf("[DEBUG] Retrieved openstack_lb_pool_v2 %s: %#v", pool.ID, pool)

	d.SetId(pool.ID)
	d.Set("project_id", pool.ProjectID)
	d.Set("name", pool.Name)
	d.Set("description", pool.Description)
	d.Set("protocol", pool.Protocol)
	d.Set("lb_method", pool.LBMethod)
	d.Set("listeners", flattenLBListenerIDsV2(pool.Listeners))
	d.Set("members", flattenLBMembersV2(pool.Members))
	d.Set("healthmonitor_id", pool.MonitorID)
	d.Set("admin_state_up", pool.AdminStateUp)
	d.Set("loadbalancers", flattenLBPoolLoadbalancerIDsV2(pool.Loadbalancers))
	d.Set("session_persistence", flattenLBPoolPersistenceV2(pool.Persistence))
	d.Set("alpn_protocols", pool.ALPNProtocols)
	d.Set("ca_tls_container_ref", pool.CATLSContainerRef)
	d.Set("crl_container_ref", pool.CRLContainerRef)
	d.Set("tls_enabled", pool.TLSEnabled)
	d.Set("tls_ciphers", pool.TLSCiphers)
	d.Set("tls_container_ref", pool.TLSContainerRef)
	d.Set("tls_versions", pool.TLSVersions)
	d.Set("provisioning_status", pool.ProvisioningStatus)
	d.Set("operating_status", pool.OperatingStatus)
	d.Set("tags", pool.Tags)
}
