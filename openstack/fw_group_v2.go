package openstack

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/groups"
)

func fwGroupV2RefreshFunc(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := groups.Get(networkingClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return group, group.Status, nil
	}
}

func fwGroupV2DeleteFunc(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := groups.Get(networkingClient, id).Extract()

		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("Unexpected error: %s", err)
		}

		return group, "DELETING", nil
	}
}

func fwGroupV2IngressPolicyDeleteFunc(networkingClient *gophercloud.ServiceClient, d *schema.ResourceData, ctx context.Context, groupID string) diag.Diagnostics {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE", "INACTIVE", "DOWN"},
		Refresh:    fwGroupV2RefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_fw_group_v2 %s to become active: %s", d.Id(), err)
	}

	_, err = groups.RemoveIngressPolicy(networkingClient, groupID).Extract()
	if err != nil {
		return diag.Errorf("Error removing ingress firewall policy from openstack_fw_group_v2 %s: %s", d.Id(), err)
	}

	return nil
}

func fwGroupV2EgressPolicyDeleteFunc(networkingClient *gophercloud.ServiceClient, d *schema.ResourceData, ctx context.Context, groupID string) diag.Diagnostics {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE", "INACTIVE", "DOWN"},
		Refresh:    fwGroupV2RefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_fw_group_v2 %s to become active: %s", d.Id(), err)
	}

	_, err = groups.RemoveEgressPolicy(networkingClient, groupID).Extract()
	if err != nil {
		return diag.Errorf("Error removing egress firewall policy from openstack_fw_group_v2 %s: %s", d.Id(), err)
	}

	return nil
}
