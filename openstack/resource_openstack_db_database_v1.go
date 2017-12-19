package openstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/db/v1/databases"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDatabaseDatabaseV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseDatabaseV1Create,
		Read:   resourceDatabaseDatabaseV1Read,
		Delete: resourceDatabaseDatabaseV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", ""),
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDatabaseDatabaseV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	databaseV1Client, err := config.databaseV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cloud database client: %s", err)
	}

	dbName := d.Get("name").(string)

	var dbs databases.BatchCreateOpts
	dbs = append(dbs, databases.CreateOpts{
		Name: dbName,
	})

	instanceID := d.Get("instance").(string)

	exists, err := DatabaseDatabaseV1State(databaseV1Client, instanceID, dbName)
	if err != nil {
		return fmt.Errorf("Error checking database status: %s", err)
	}
	if exists {
		return fmt.Errorf("Database %s exists on instance %s", dbName, instanceID)
	}

	databases.Create(databaseV1Client, instanceID, dbs)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    DatabaseDatabaseV1StateRefreshFunc(databaseV1Client, instanceID, dbName),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for database to become ready: %s", err)
	}

	// Store the ID now
	d.SetId(fmt.Sprintf("%s.%s", instanceID, dbName))

	return resourceDatabaseInstanceV1Read(d, meta)
}

func resourceDatabaseDatabaseV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	databaseV1Client, err := config.databaseV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cloud database client: %s", err)
	}

	dbID := strings.Split(d.Id(), ".")
	instanceID := dbID[0]
	dbName := dbID[1]

	exists, err := DatabaseDatabaseV1State(databaseV1Client, instanceID, dbName)
	if err != nil {
		return fmt.Errorf("Error checking database status: %s", err)
	}
	if !exists {
		return fmt.Errorf("Error, database %s was not found", err)
	}

	log.Printf("[DEBUG] Retrieved database %s", dbName)

	d.Set("name", dbName)

	return nil
}

func resourceDatabaseDatabaseV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	databaseV1Client, err := config.databaseV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cloud database client: %s", err)
	}

	dbID := strings.Split(d.Id(), ".")
	instanceID := dbID[0]
	dbName := dbID[1]

	exists, err := DatabaseDatabaseV1State(databaseV1Client, instanceID, dbName)
	if err != nil {
		return fmt.Errorf("Error checking database status: %s", err)
	}
	if !exists {
		return fmt.Errorf("Database %s does not exist on instance %s", dbName, instanceID)
	}

	databases.Delete(databaseV1Client, instanceID, dbName)

	d.SetId("")
	return nil
}

// DatabaseDatabaseV1StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an cloud database.
func DatabaseDatabaseV1StateRefreshFunc(client *gophercloud.ServiceClient, instanceID string, dbName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pages, err := databases.List(client, instanceID).AllPages()
		if err != nil {
			return nil, "", fmt.Errorf("Unable to retrieve databases, pages: %s", err)
		}

		allDatabases, err := databases.ExtractDBs(pages)
		if err != nil {
			return nil, "", fmt.Errorf("Unable to retrieve databases, extract: %s", err)
		}

		for _, v := range allDatabases {
			if v.Name == dbName {
				return v, "ACTIVE", nil
			}
		}

		return nil, "", fmt.Errorf("Error retrieving database %s status", dbName)
	}
}

func DatabaseDatabaseV1State(client *gophercloud.ServiceClient, instanceID string, dbName string) (exists bool, err error) {
	exists = false
	err = nil

	pages, err := databases.List(client, instanceID).AllPages()
	if err != nil {
		return
	}

	allDatabases, err := databases.ExtractDBs(pages)
	if err != nil {
		return
	}

	for _, v := range allDatabases {
		if v.Name == dbName {
			exists = true
			return
		}
	}

	return false, err

}
