package openstack

import (
	"fmt"
	"log"
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

	dbname := d.Get("name").(string)

	var dbs databases.BatchCreateOpts
	dbs = append(dbs, databases.CreateOpts{
		Name: dbname,
	})

	instance_id := d.Get("instance").(string)
	databases.Create(databaseV1Client, instance_id, dbs)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    DatabaseDatabaseV1StateRefreshFunc(databaseV1Client, instance_id, dbname),
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
	d.SetId(instance_id)

	return resourceDatabaseInstanceV1Read(d, meta)
}

func resourceDatabaseDatabaseV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	databaseV1Client, err := config.databaseV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cloud database client: %s", err)
	}

	dbname := d.Get("name").(string)

	pages, err := databases.List(databaseV1Client, d.Id()).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to retrieve databases, pages: %s", err)
	}
	allDatabases, err := databases.ExtractDBs(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve databases, extract: %s", err)
	}

	for _, v := range allDatabases {
		if v.Name == dbname {
			d.Set("name", v.Name)
			break
		}
	}
	log.Printf("[DEBUG] Retrieved database %s", dbname)

	return nil
}

func resourceDatabaseDatabaseV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	databaseV1Client, err := config.databaseV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cloud database client: %s", err)
	}

	dbname := d.Get("name").(string)

	pages, err := databases.List(databaseV1Client, d.Id()).AllPages()
	allDatabases, err := databases.ExtractDBs(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve databases: %s", err)
	}

	log.Println("Retrieved databases", allDatabases)
	log.Println("Looking for db", dbname)

	dbExists := false

	for _, v := range allDatabases {
		if v.Name == dbname {
			dbExists = true
			break
		}
	}

	if !dbExists {
		log.Printf("Database %s was not found on instance %s", dbname, d.Id())
	}

	databases.Delete(databaseV1Client, d.Id(), dbname)

	d.SetId("")
	return nil
}

// DatabaseDatabaseV1StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an cloud database.
func DatabaseDatabaseV1StateRefreshFunc(client *gophercloud.ServiceClient, instance_id string, dbname string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pages, err := databases.List(client, instance_id).AllPages()
		if err != nil {
			return nil, "", fmt.Errorf("Unable to retrieve databases, pages: %s", err)
		}

		allDatabases, err := databases.ExtractDBs(pages)
		if err != nil {
			return nil, "", fmt.Errorf("Unable to retrieve databases, extract: %s", err)
		}

		for _, v := range allDatabases {
			if v.Name == dbname {
				return v, "ACTIVE", nil
			}
		}

		return nil, "", fmt.Errorf("Error retrieving database %s status", dbname)
	}
}
