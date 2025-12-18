package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLBMonitorV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBMonitorV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"monitor_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name"},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"monitor_id"},
			},

			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"expected_codes": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"delay": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_retries": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_retries_down": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"http_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_name": {
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

func dataSourceLBMonitorV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancer client: %s", err)
	}

	listOpts := monitors.ListOpts{
		Tags: expandTagsList(d, "tags"),
	}

	if v, ok := d.GetOk("monitor_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("pool_id"); ok {
		listOpts.PoolID = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		listOpts.Type = v.(string)
	}

	if v, ok := d.GetOk("http_method"); ok {
		listOpts.HTTPMethod = v.(string)
	}

	if v, ok := d.GetOk("url_path"); ok {
		listOpts.URLPath = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	if v, ok := d.GetOk("expected_codes"); ok {
		listOpts.ExpectedCodes = v.(string)
	}

	allPages, err := monitors.List(lbClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query OpenStack loadbalancer monitors: %s", err)
	}

	allMonitors, err := monitors.ExtractMonitors(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve Openstack loadbalancer monitors: %s", err)
	}

	if len(allMonitors) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allMonitors) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allMonitors)

		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria")
	}

	dataSourceLBMonitorV2Attributes(d, &allMonitors[0])
	d.Set("region", GetRegion(d, config))

	return nil
}

func dataSourceLBMonitorV2Attributes(d *schema.ResourceData, monitor *monitors.Monitor) {
	log.Printf("[DEBUG] Retrieved openstack_lb_monitor_v2 %s: %#v", monitor.ID, monitor)

	d.SetId(monitor.ID)
	d.Set("project_id", monitor.ProjectID)
	d.Set("name", monitor.Name)
	d.Set("type", monitor.Type)
	d.Set("delay", monitor.Delay)
	d.Set("timeout", monitor.Timeout)
	d.Set("max_retries", monitor.MaxRetries)
	d.Set("max_retries_down", monitor.MaxRetriesDown)
	d.Set("http_method", monitor.HTTPMethod)
	d.Set("http_version", monitor.HTTPVersion)
	d.Set("url_path", monitor.URLPath)
	d.Set("expected_codes", monitor.ExpectedCodes)
	d.Set("domain_name", monitor.DomainName)
	d.Set("admin_state_up", monitor.AdminStateUp)
	d.Set("status", monitor.Status)
	d.Set("pools", flattenLBMonitorPoolsIDsV2(monitor.Pools))
	d.Set("provisioning_status", monitor.ProvisioningStatus)
	d.Set("operating_status", monitor.OperatingStatus)
	d.Set("tags", monitor.Tags)
}
