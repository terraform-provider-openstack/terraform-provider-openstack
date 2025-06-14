package openstack

import (
	"context"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabaseUserV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseUserV1Create,
		ReadContext:   resourceDatabaseUserV1Read,
		DeleteContext: resourceDatabaseUserV1Delete,

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

			"password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"host": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"databases": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceDatabaseUserV1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	databaseV1Client, err := config.DatabaseV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	userName := d.Get("name").(string)
	rawDatabases := d.Get("databases").(*schema.Set).List()
	instanceID := d.Get("instance_id").(string)

	var usersList users.BatchCreateOpts
	usersList = append(usersList, users.CreateOpts{
		Name:      userName,
		Password:  d.Get("password").(string),
		Host:      d.Get("host").(string),
		Databases: expandDatabaseUserV1Databases(rawDatabases),
	})

	err = users.Create(ctx, databaseV1Client, instanceID, usersList).ExtractErr()
	if err != nil {
		return diag.Errorf("Error creating openstack_db_user_v1: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    databaseUserV1StateRefreshFunc(ctx, databaseV1Client, instanceID, userName),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_db_user_v1 %s to be created: %s", userName, err)
	}

	// Store the ID now
	d.SetId(fmt.Sprintf("%s/%s", instanceID, userName))

	return resourceDatabaseUserV1Read(ctx, d, meta)
}

func resourceDatabaseUserV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	databaseV1Client, err := config.DatabaseV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	instanceID, userName, err := parsePairedIDs(d.Id(), "openstack_db_user_v1")
	if err != nil {
		return diag.FromErr(err)
	}

	exists, userObj, err := databaseUserV1Exists(ctx, databaseV1Client, instanceID, userName)
	if err != nil {
		return diag.Errorf("Error checking if openstack_db_user_v1 %s exists: %s", d.Id(), err)
	}

	if !exists {
		d.SetId("")

		return nil
	}

	d.Set("name", userName)
	d.Set("region", GetRegion(d, config))

	databases := flattenDatabaseUserV1Databases(userObj.Databases)
	if err := d.Set("databases", databases); err != nil {
		return diag.Errorf("Unable to set databases: %s", err)
	}

	return nil
}

func resourceDatabaseUserV1Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	databaseV1Client, err := config.DatabaseV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	instanceID, userName, err := parsePairedIDs(d.Id(), "openstack_db_user_v1")
	if err != nil {
		return diag.FromErr(err)
	}

	exists, _, err := databaseUserV1Exists(ctx, databaseV1Client, instanceID, userName)
	if err != nil {
		return diag.Errorf("Error checking if openstack_db_user_v1 %s exists: %s", d.Id(), err)
	}

	if !exists {
		return nil
	}

	err = users.Delete(ctx, databaseV1Client, instanceID, userName).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_db_user_v1"))
	}

	return nil
}
