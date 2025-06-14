package openstack

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/configurations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabaseConfigurationV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseConfigurationV1Create,
		ReadContext:   resourceDatabaseConfigurationV1Read,
		DeleteContext: resourceDatabaseConfigurationV1Delete,

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

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
				MaxItems: 1,
			},

			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"string_type": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceDatabaseConfigurationV1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	databaseV1Client, err := config.DatabaseV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	createOpts := &configurations.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	var datastore configurations.DatastoreOpts
	if v, ok := d.GetOk("datastore"); ok {
		datastore = expandDatabaseConfigurationV1Datastore(v.([]any))
	}

	createOpts.Datastore = &datastore

	values := make(map[string]any)
	if v, ok := d.GetOk("configuration"); ok {
		values = expandDatabaseConfigurationV1Values(v.([]any))
	}

	createOpts.Values = values

	log.Printf("[DEBUG] openstack_db_configuration_v1 create options: %#v", createOpts)

	cgroup, err := configurations.Create(ctx, databaseV1Client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_db_configuration_v1: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    databaseConfigurationV1StateRefreshFunc(ctx, databaseV1Client, cgroup.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_db_configuration_v1 %s to become ready: %s", cgroup.ID, err)
	}

	// Store the ID now
	d.SetId(cgroup.ID)

	return resourceDatabaseConfigurationV1Read(ctx, d, meta)
}

func resourceDatabaseConfigurationV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	databaseV1Client, err := config.DatabaseV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	cgroup, err := configurations.Get(ctx, databaseV1Client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_db_configuration_v1"))
	}

	log.Printf("[DEBUG] Retrieved openstack_db_configuration_v1 %s: %#v", d.Id(), cgroup)

	d.Set("name", cgroup.Name)
	d.Set("description", cgroup.Description)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceDatabaseConfigurationV1Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	databaseV1Client, err := config.DatabaseV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	err = configurations.Delete(ctx, databaseV1Client, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_db_configuration_v1"))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE", "SHUTOFF"},
		Target:     []string{"DELETED"},
		Refresh:    databaseConfigurationV1StateRefreshFunc(ctx, databaseV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_db_configuration_v1 %s to Delete:  %s", d.Id(), err)
	}

	return nil
}
