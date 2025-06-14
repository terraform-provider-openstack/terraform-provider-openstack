package openstack

import (
	"context"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabaseDatabaseV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseDatabaseV1Create,
		ReadContext:   resourceDatabaseDatabaseV1Read,
		DeleteContext: resourceDatabaseDatabaseV1Delete,
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDatabaseDatabaseV1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	databaseV1Client, err := config.DatabaseV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	dbName := d.Get("name").(string)
	instanceID := d.Get("instance_id").(string)

	var dbs databases.BatchCreateOpts
	dbs = append(dbs, databases.CreateOpts{
		Name: dbName,
	})

	exists, err := databaseDatabaseV1Exists(ctx, databaseV1Client, instanceID, dbName)
	if err != nil {
		return diag.Errorf("Error checking openstack_db_database_v1 %s status on %s: %s", dbName, instanceID, err)
	}

	if exists {
		return diag.Errorf("openstack_db_database_v1 %s already exists on instance %s", dbName, instanceID)
	}

	err = databases.Create(ctx, databaseV1Client, instanceID, dbs).ExtractErr()
	if err != nil {
		return diag.Errorf("Error creating openstack_db_database_v1 %s on %s: %s", dbName, instanceID, err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    databaseDatabaseV1StateRefreshFunc(ctx, databaseV1Client, instanceID, dbName),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_db_database_v1 %s on %s to become ready: %s", dbName, instanceID, err)
	}

	// Store the ID now
	d.SetId(fmt.Sprintf("%s/%s", instanceID, dbName))

	return resourceDatabaseDatabaseV1Read(ctx, d, meta)
}

func resourceDatabaseDatabaseV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	databaseV1Client, err := config.DatabaseV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	instanceID, dbName, err := parsePairedIDs(d.Id(), "openstack_db_database_v1")
	if err != nil {
		return diag.FromErr(err)
	}

	exists, err := databaseDatabaseV1Exists(ctx, databaseV1Client, instanceID, dbName)
	if err != nil {
		return diag.Errorf("Error checking if openstack_db_database_v1 %s exists: %s", d.Id(), err)
	}

	if !exists {
		d.SetId("")

		return nil
	}

	d.Set("instance_id", instanceID)
	d.Set("name", dbName)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceDatabaseDatabaseV1Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	databaseV1Client, err := config.DatabaseV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	instanceID, dbName, err := parsePairedIDs(d.Id(), "openstack_db_database_v1")
	if err != nil {
		return diag.FromErr(err)
	}

	exists, err := databaseDatabaseV1Exists(ctx, databaseV1Client, instanceID, dbName)
	if err != nil {
		return diag.Errorf("Error checking if openstack_db_database_v1 %s exists: %s", d.Id(), err)
	}

	if !exists {
		return nil
	}

	err = databases.Delete(ctx, databaseV1Client, instanceID, dbName).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_db_database_v1"))
	}

	return nil
}
