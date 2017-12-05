package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/openstack/db/v1/users"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDatabaseInstanceV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseInstanceV1Create,
		Read:   resourceDatabaseInstanceV1Read,
		Delete: resourceDatabaseInstanceV1Delete,
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
			"flavor_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_FLAVOR_ID", nil),
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"datastore": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"port": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"fixed_ip_v4": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"fixed_ip_v6": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"database": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"charset": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"collate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"user": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"host": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"databases": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
			},
			"configuration": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
		},
	}
}

func resourceDatabaseInstanceV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	databaseV1Client, err := config.databaseV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cloud database client: %s", err)
	}

	var datastore instances.DatastoreOpts
	if p, ok := d.GetOk("datastore"); ok {
		pV := (p.([]interface{}))[0].(map[string]interface{})

		datastore = instances.DatastoreOpts{
			Version: pV["version"].(string),
			Type:    pV["type"].(string),
		}
	}

	createOpts := &instances.CreateOpts{
		FlavorRef: d.Get("flavor_id").(string),
		Name:      d.Get("name").(string),
		Size:      d.Get("size").(int),
	}

	createOpts.Datastore = &datastore

	// networks
	var networks []instances.NetworkOpts

	if p, ok := d.GetOk("network"); ok {
		if networkList, ok := p.([]interface{}); ok {

			for _, network := range networkList {
				networks = append(networks, instances.NetworkOpts{
					UUID:      network.(map[string]interface{})["uuid"].(string),
					Port:      network.(map[string]interface{})["port"].(string),
					V4FixedIP: network.(map[string]interface{})["fixed_ip_v4"].(string),
					V6FixedIP: network.(map[string]interface{})["fixed_ip_v6"].(string),
				})
			}

		}
	}

	createOpts.Networks = networks

	// databases
	var dbs databases.BatchCreateOpts

	if p, ok := d.GetOk("database"); ok {
		if databaseList, ok := p.([]interface{}); ok {

			for _, db := range databaseList {
				dbs = append(dbs, databases.CreateOpts{
					Name:    db.(map[string]interface{})["name"].(string),
					CharSet: db.(map[string]interface{})["charset"].(string),
					Collate: db.(map[string]interface{})["collate"].(string),
				})
			}

		}
	}

	createOpts.Databases = dbs

	// users
	var UserList users.BatchCreateOpts

	if p, ok := d.GetOk("user"); ok {
		if userList, ok := p.([]interface{}); ok {
			for _, user := range userList {
				UserList = append(UserList, users.CreateOpts{
					Name:      user.(map[string]interface{})["name"].(string),
					Password:  user.(map[string]interface{})["password"].(string),
					Databases: resourceDBv1GetDatabases(user.(map[string]interface{})["databases"].(*schema.Set).List()),
					Host:      user.(map[string]interface{})["host"].(string),
				})
			}
		}
	}

	createOpts.Users = UserList

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	instance, err := instances.Create(databaseV1Client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating cloud database instance: %s", err)
	}
	log.Printf("[INFO] instance ID: %s", instance.ID)

	// Wait for the volume to become available.
	log.Printf(
		"[DEBUG] Waiting for volume (%s) to become available",
		instance.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    DatabaseInstanceV1StateRefreshFunc(databaseV1Client, instance.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			instance.ID, err)
	}

	if configuration, ok := d.GetOk("configuration"); ok {
		instances.AttachConfigurationGroup(databaseV1Client, instance.ID, configuration.(string))
		log.Printf("Attaching configuration %v to the instance %v", configuration, instance.ID)
		instances.Restart(databaseV1Client, instance.ID)
		log.Printf("Restarting instance instance %v", instance.ID)
	}

	// Store the ID now
	d.SetId(instance.ID)

	return resourceDatabaseInstanceV1Read(d, meta)
}

func resourceDatabaseInstanceV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	databaseV1Client, err := config.databaseV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack cloud database client: %s", err)
	}

	instance, err := instances.Get(databaseV1Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "instance")
	}

	log.Printf("[DEBUG] Retrieved instance %s: %+v", d.Id(), instance)

	d.Set("name", instance.Name)
	d.Set("flavor_id", instance.Flavor)
	d.Set("datastore", instance.Datastore)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceDatabaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	databaseV1Client, err := config.databaseV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack cloud database client: %s", err)
	}

	if d.HasChange("configuration") {
		old, new := d.GetChange("configuration")
		instances.DetachConfigurationGroup(databaseV1Client, d.Id())
		log.Printf("Detaching configuration %v from the instance %v", old, d.Id())
		instances.AttachConfigurationGroup(databaseV1Client, d.Id(), new.(string))
		log.Printf("Attaching configuration %v to the instance %v", new, d.Id())
	}

	return resourceDatabaseInstanceV1Read(d, meta)
}

func resourceDatabaseInstanceV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	databaseV1Client, err := config.databaseV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating RS cloud instance client: %s", err)
	}

	log.Printf("[DEBUG] Deleting cloud database instance %s", d.Id())
	err = instances.Delete(databaseV1Client, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting cloud database instance: %s", err)
	}

	// Wait for the volume to delete before moving on.
	log.Printf("[DEBUG] Waiting for volume (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "SHUTDOWN"},
		Target:     []string{"DELETED"},
		Refresh:    DatabaseInstanceV1StateRefreshFunc(databaseV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

// DatabaseInstanceV1StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an cloud database instance.
func DatabaseInstanceV1StateRefreshFunc(client *gophercloud.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		i, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return i, "DELETED", nil
			}
			return nil, "", err
		}

		if i.Status == "error" {
			return i, i.Status, fmt.Errorf("There was an error creating the instance.")
		}

		return i, i.Status, nil
	}
}

func resourceDBv1GetDatabases(v []interface{}) databases.BatchCreateOpts {

	var dbs databases.BatchCreateOpts

	for _, db := range v {
		dbs = append(dbs, databases.CreateOpts{
			Name: db.(string),
		})
	}

	return dbs
}
