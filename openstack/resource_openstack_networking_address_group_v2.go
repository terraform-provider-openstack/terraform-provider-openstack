package openstack

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/addressgroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkingAddressGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingAddressGroupV2Create,
		ReadContext:   resourceNetworkingAddressGroupV2Read,
		UpdateContext: resourceNetworkingAddressGroupV2Update,
		DeleteContext: resourceNetworkingAddressGroupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
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
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"addresses": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},
		},
	}
}

func resourceNetworkingAddressGroupV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	opts := addressgroups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ProjectID:   d.Get("project_id").(string),
		Addresses:   expandToStringSlice(d.Get("addresses").(*schema.Set).List()),
	}

	log.Printf("[DEBUG] openstack_networking_address_group_v2 create options: %#v", opts)

	ag, err := addressgroups.Create(ctx, networkingClient, opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_address_group_v2: %s", err)
	}

	d.SetId(ag.ID)

	return resourceNetworkingAddressGroupV2Read(ctx, d, meta)
}

func resourceNetworkingAddressGroupV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	ag, err := addressgroups.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_networking_address_group_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_address_group_v2: %#v", ag)

	d.Set("name", ag.Name)
	d.Set("description", ag.Description)
	d.Set("project_id", ag.ProjectID)
	d.Set("addresses", ag.Addresses)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingAddressGroupV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		updated    bool
		updateOpts addressgroups.UpdateOpts
	)

	if d.HasChange("name") {
		updated = true
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		updated = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if updated {
		log.Printf("[DEBUG] Updating openstack_networking_address_group_v2 %s with options: %#v", d.Id(), updateOpts)

		_, err = addressgroups.Update(ctx, networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_address_group_v2: %s", err)
		}
	}

	if d.HasChange("addresses") {
		o, n := d.GetChange("addresses")
		oldAddr, newAddr := o.(*schema.Set), n.(*schema.Set)
		addrToDel := oldAddr.Difference(newAddr)
		addrToAdd := newAddr.Difference(oldAddr)

		if v := addrToDel.List(); len(v) > 0 {
			log.Printf("[DEBUG] Removing addresses '%s' from openstack_networking_address_group_v2 '%s'", v, d.Get("name"))
			opts := addressgroups.UpdateAddressesOpts{
				Addresses: expandToStringSlice(v),
			}

			_, err = addressgroups.RemoveAddresses(ctx, networkingClient, d.Id(), opts).Extract()
			if err != nil {
				return diag.Errorf("Error deleting addresses '%s' from openstack_networking_address_group_v2: %s", v, err)
			}
		}

		if v := addrToAdd.List(); len(v) > 0 {
			log.Printf("[DEBUG] Adding addresses '%s' to openstack_networking_address_group_v2 '%s'", v, d.Get("name"))
			opts := addressgroups.UpdateAddressesOpts{
				Addresses: expandToStringSlice(v),
			}

			_, err = addressgroups.AddAddresses(ctx, networkingClient, d.Id(), opts).Extract()
			if err != nil {
				return diag.Errorf("Error adding addresses '%s' to openstack_networking_address_group_v2: %s", v, err)
			}
		}
	}

	return resourceNetworkingAddressGroupV2Read(ctx, d, meta)
}

func resourceNetworkingAddressGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingAddressGroupV2StateRefreshFuncDelete(ctx, networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting openstack_networking_address_group_v2: %s", err)
	}

	return nil
}
