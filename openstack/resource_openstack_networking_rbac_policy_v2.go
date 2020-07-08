package openstack

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/rbacpolicies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/ryanuber/go-glob"
)

func resourceNetworkingRBACPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingRBACPolicyV2Create,
		Read:   resourceNetworkingRBACPolicyV2Read,
		Update: resourceNetworkingRBACPolicyV2Update,
		Delete: resourceNetworkingRBACPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"access_as_external", "access_as_shared",
				}, false),
			},

			"object_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"address_scope", "network", "qos_policy", "security_group", "subnetpool",
				}, false),
			},

			"target_tenant": {
				Type:     schema.TypeString,
				Required: true,
			},

			"object_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"object_id", "object_search"},
			},

			"object_search": {
				Type:         schema.TypeSet,
				ForceNew:     true,
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"object_id", "object_search"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name_glob": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The blob expression to use for searching the object.",
						},
						"owning_tenant_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the tenant that owns the resource.",
						},
						"unshared": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Selet only resources that are not yet shared.",
						},
					},
				},
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func loadNetworksMatchingGlob(networkingClient *gophercloud.ServiceClient, owningTenantID *string, nameGlob string) ([]string, error) {
	filteredNetworks := []string{}

	listOpts := networks.ListOpts{
		SortKey: "name",
		SortDir: "asc",
	}

	if owningTenantID != nil {
		listOpts.TenantID = *owningTenantID
	}

	// Load all the objects
	allPages, err := networks.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return filteredNetworks, err
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return filteredNetworks, err
	}

	// List of all the networks matching the glob
	for _, network := range allNetworks {
		if glob.Glob(nameGlob, network.Name) {
			filteredNetworks = append(filteredNetworks, network.ID)
		}
	}

	return filteredNetworks, nil
}

func loadPoliciesMatchingGlob(networkingClient *gophercloud.ServiceClient, owningTenantID *string, nameGlob string) ([]string, error) {
	filteredPolicies := []string{}

	listOpts := policies.ListOpts{
		SortKey: "name",
		SortDir: "asc",
	}

	if owningTenantID != nil {
		listOpts.TenantID = *owningTenantID
	}

	// Load all the objects
	allPages, err := policies.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return filteredPolicies, err
	}

	allPolicies, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return filteredPolicies, err
	}

	// List of all the networks matching the glob
	for _, policy := range allPolicies {
		if glob.Glob(nameGlob, policy.Name) {
			filteredPolicies = append(filteredPolicies, policy.ID)
		}
	}

	return filteredPolicies, nil
}

func resourceNetworkingRBACPolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	action := rbacpolicies.PolicyAction(d.Get("action").(string))
	objectType := d.Get("object_type").(string)
	targetTenant := d.Get("target_tenant").(string)
	objectID := d.Get("object_id").(string)

	objectSearch, objectSearchExists := d.GetOk("object_search")

	if objectSearchExists {
		// we have to search for a suitable object and try to allocate it
		objectSearchMap := objectSearch.(*schema.Set).List()[0].(map[string]interface{})
		nameGlob := objectSearchMap["name_glob"].(string)
		wasUnshared := objectSearchMap["unshared"].(bool)
		owningTenantID, hasOwningTenantID := objectSearchMap["owningTenant"]

		var owningTenantIDPtr *string = nil
		if hasOwningTenantID {
			owningTenantIDString := owningTenantID.(string)
			owningTenantIDPtr = &owningTenantIDString
		}

		// Retrieve the objects
		var filteredObjects []string

		if objectType == "network" {
			filteredObjects, err = loadNetworksMatchingGlob(networkingClient, owningTenantIDPtr, nameGlob)
		} else {
			filteredObjects, err = loadPoliciesMatchingGlob(networkingClient, owningTenantIDPtr, nameGlob)
		}
		if err != nil {
			return err
		}

		if len(filteredObjects) == 0 {
			return fmt.Errorf("Unable to find any object matching the search criteria")
		}

		if wasUnshared {
			// loop until we completed the work (it is possible that we allocate multiple times a resource)
			for i := 0; i < 5; i++ {
				// load all the rbacs currently created
				allRBACPages, err := rbacpolicies.List(networkingClient, rbacpolicies.ListOpts{
					ObjectType: objectType,
				}).AllPages()
				if err != nil {
					return err
				}

				allRBACs, err := rbacpolicies.ExtractRBACPolicies(allRBACPages)
				if err != nil {
					return err
				}

				// Transform the list into a map
				allocatedObjects := make(map[string]bool, len(allRBACs))
				for _, rbac := range allRBACs {
					allocatedObjects[rbac.ObjectID] = true
				}

				// Iterate the list of objects discarding the allocated ones
				objectID := ""
				for _, fitleredObjectID := range filteredObjects {
					if allocatedObjects[fitleredObjectID] {
						continue
					}
					log.Printf("Found object: %s", fitleredObjectID)
					objectID = fitleredObjectID
					break
				}

				if objectID == "" {
					return fmt.Errorf("Unable to find an unshared object")
				}

				// Create the sharing and verify it is the only one
				createOpts := rbacpolicies.CreateOpts{
					Action:       action,
					ObjectType:   objectType,
					TargetTenant: targetTenant,
					ObjectID:     objectID,
				}

				log.Printf("[DEBUG] Create Options: %#v", createOpts)
				rbac, err := rbacpolicies.Create(networkingClient, createOpts).Extract()
				if err != nil {
					return fmt.Errorf("Error creating openstack_networking_rbac_policy_v2: %s", err)
				}
				log.Printf("Optimistic creation of RBAC succeeded: %s", rbac.ID)

				// Read the policy for the object, checking that it is shared only with one tenant
				allRBACPages, err = rbacpolicies.List(networkingClient, rbacpolicies.ListOpts{
					ObjectType: objectType,
					ObjectID:   rbac.ObjectID,
				}).AllPages()
				if err != nil {
					return err
				}

				allRBACs, err = rbacpolicies.ExtractRBACPolicies(allRBACPages)
				if err != nil {
					return err
				}

				// If the RBACs are different than one, let's remove he RBAC and try again
				if len(allRBACs) != 1 {
					log.Printf("Conflict while allocating object %s, removing sharing to tenant %s", objectID, targetTenant)
					err = rbacpolicies.Delete(networkingClient, rbac.ID).ExtractErr()
					if err != nil {
						return fmt.Errorf("Error while removing the conflicting RBAC: %s", rbac.ID)
					}

					// sleep randomly before retyring to avoid conflicting again.
					// Wait between 0 and 10000 milliseconds
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					n := r.Intn(10000)
					time.Sleep(time.Duration(n) * time.Millisecond)

					continue
				}

				// Commit the change
				d.SetId(rbac.ID)
				d.Set("object_id", objectID)

				log.Printf("Creation of RBAC completed: %s", rbac.ID)

				return resourceNetworkingRBACPolicyV2Read(d, meta)
			}

			return fmt.Errorf("Unable to allocate the rbac")
		}

		// Select the object ID from the list
		objectID = filteredObjects[0]
	}

	createOpts := rbacpolicies.CreateOpts{
		Action:       action,
		ObjectType:   objectType,
		TargetTenant: targetTenant,
		ObjectID:     objectID,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	rbac, err := rbacpolicies.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_networking_rbac_policy_v2: %s", err)
	}

	d.SetId(rbac.ID)
	d.Set("object_id", objectID)

	return resourceNetworkingRBACPolicyV2Read(d, meta)
}

func resourceNetworkingRBACPolicyV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	rbac, err := rbacpolicies.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_networking_rbac_policy_v2")
	}

	log.Printf("[DEBUG] Retrieved RBAC policy %s: %+v", d.Id(), rbac)

	d.Set("action", string(rbac.Action))
	d.Set("object_type", rbac.ObjectType)
	d.Set("target_tenant", rbac.TargetTenant)
	d.Set("object_id", rbac.ObjectID)
	d.Set("project_id", rbac.ProjectID)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingRBACPolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateOpts rbacpolicies.UpdateOpts

	if d.HasChange("target_tenant") {
		updateOpts.TargetTenant = d.Get("target_tenant").(string)

		_, err := rbacpolicies.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating openstack_networking_rbac_policy_v2: %s", err)
		}
	}

	return resourceNetworkingRBACPolicyV2Read(d, meta)
}

func resourceNetworkingRBACPolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	err = rbacpolicies.Delete(networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting openstack_networking_rbac_policy_v2")
	}

	return nil
}
