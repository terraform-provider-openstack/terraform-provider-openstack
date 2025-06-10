package openstack

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/endpointgroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEndpointGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointGroupV2Create,
		ReadContext:   resourceEndpointGroupV2Read,
		UpdateContext: resourceEndpointGroupV2Update,
		DeleteContext: resourceEndpointGroupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
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
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"endpoints": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceEndpointGroupV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var createOpts endpointgroups.CreateOptsBuilder

	endpointType := resourceEndpointGroupV2EndpointType(d.Get("type").(string))
	endpoints := expandToStringSlice(d.Get("endpoints").(*schema.Set).List())

	createOpts = EndpointGroupCreateOpts{
		endpointgroups.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			TenantID:    d.Get("tenant_id").(string),
			Endpoints:   endpoints,
			Type:        endpointType,
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create group: %#v", createOpts)

	group, err := endpointgroups.Create(ctx, networkingClient, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForEndpointGroupCreation(ctx, networkingClient, group.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] EndpointGroup created: %#v", group)

	d.SetId(group.ID)

	return resourceEndpointGroupV2Read(ctx, d, meta)
}

func resourceEndpointGroupV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about group: %s", d.Id())

	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	group, err := endpointgroups.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "group"))
	}

	log.Printf("[DEBUG] Read OpenStack Endpoint EndpointGroup %s: %#v", d.Id(), group)

	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("tenant_id", group.TenantID)
	d.Set("type", group.Type)
	d.Set("endpoints", group.Endpoints)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceEndpointGroupV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	opts := endpointgroups.UpdateOpts{}

	var hasChange bool

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
		hasChange = true
	}

	var updateOpts endpointgroups.UpdateOptsBuilder = opts

	log.Printf("[DEBUG] Updating endpoint group with id %s: %#v", d.Id(), updateOpts)

	if hasChange {
		group, err := endpointgroups.Update(ctx, networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"UPDATED"},
			Refresh:    waitForEndpointGroupUpdate(ctx, networkingClient, group.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updated group with id %s", d.Id())
	}

	return resourceEndpointGroupV2Read(ctx, d, meta)
}

func resourceEndpointGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy group: %s", d.Id())

	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	err = endpointgroups.Delete(ctx, networkingClient, d.Id()).Err
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForEndpointGroupDeletion(ctx, networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)

	return diag.FromErr(err)
}

func waitForEndpointGroupDeletion(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		group, err := endpointgroups.Get(ctx, networkingClient, id).Extract()
		log.Printf("[DEBUG] Got group %s => %#v", id, group)

		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				log.Printf("[DEBUG] EndpointGroup %s is actually deleted", id)

				return "", "DELETED", nil
			}

			return nil, "", fmt.Errorf("Unexpected error: %w", err)
		}

		log.Printf("[DEBUG] EndpointGroup %s deletion is pending", id)

		return group, "DELETING", nil
	}
}

func waitForEndpointGroupCreation(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		group, err := endpointgroups.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}

		return group, "ACTIVE", nil
	}
}

func waitForEndpointGroupUpdate(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		group, err := endpointgroups.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}

		return group, "UPDATED", nil
	}
}

func resourceEndpointGroupV2EndpointType(epType string) endpointgroups.EndpointType {
	var et endpointgroups.EndpointType

	switch epType {
	case "subnet":
		et = endpointgroups.TypeSubnet
	case "cidr":
		et = endpointgroups.TypeCIDR
	case "vlan":
		et = endpointgroups.TypeVLAN
	case "router":
		et = endpointgroups.TypeRouter
	case "network":
		et = endpointgroups.TypeNetwork
	}

	return et
}
