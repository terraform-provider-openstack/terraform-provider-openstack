package openstack

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/fwaas_v2/groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func fwGroupV2RefreshFunc(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		group, err := groups.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return group, group.Status, nil
	}
}

func fwGroupV2DeleteFunc(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		group, err := groups.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return "", "DELETED", nil
			}

			return nil, "", fmt.Errorf("Unexpected error: %w", err)
		}

		return group, "DELETING", nil
	}
}

func fwGroupV2IngressPolicyDeleteFunc(ctx context.Context, networkingClient *gophercloud.ServiceClient, d *schema.ResourceData, groupID string) diag.Diagnostics {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE", "INACTIVE", "DOWN"},
		Refresh:    fwGroupV2RefreshFunc(ctx, networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_fw_group_v2 %s to become active: %s", d.Id(), err)
	}

	_, err = groups.RemoveIngressPolicy(ctx, networkingClient, groupID).Extract()
	if err != nil {
		return diag.Errorf("Error removing ingress firewall policy from openstack_fw_group_v2 %s: %s", d.Id(), err)
	}

	return nil
}

func fwGroupV2EgressPolicyDeleteFunc(ctx context.Context, networkingClient *gophercloud.ServiceClient, d *schema.ResourceData, groupID string) diag.Diagnostics {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE", "INACTIVE", "DOWN"},
		Refresh:    fwGroupV2RefreshFunc(ctx, networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_fw_group_v2 %s to become active: %s", d.Id(), err)
	}

	_, err = groups.RemoveEgressPolicy(ctx, networkingClient, groupID).Extract()
	if err != nil {
		return diag.Errorf("Error removing egress firewall policy from openstack_fw_group_v2 %s: %s", d.Id(), err)
	}

	return nil
}
