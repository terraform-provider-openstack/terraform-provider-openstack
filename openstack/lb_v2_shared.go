package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	octavialoadbalancers "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	octaviapools "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	neutronloadbalancers "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

const octaviaLBClientType = "load-balancer"

const (
	lbPendingCreate = "PENDING_CREATE"
	lbPendingUpdate = "PENDING_UPDATE"
	lbPendingDelete = "PENDING_DELETE"
	lbActive        = "ACTIVE"
	lbError         = "ERROR"
)

// lbPendingStatuses are the valid statuses a LoadBalancer will be in while
// it's updating.
func getLbPendingStatuses() []string {
	return []string{lbPendingCreate, lbPendingUpdate}
}

// lbPendingDeleteStatuses are the valid statuses a LoadBalancer will be before delete.
func getLbPendingDeleteStatuses() []string {
	return []string{lbError, lbPendingUpdate, lbPendingDelete, lbActive}
}

func getLbSkipStatuses() []string {
	return []string{lbError, lbActive}
}

// chooseLBV2Client will determine which load balacing client to use:
// either the Octavia/LBaaS client or the Neutron/Networking v2 client.
func chooseLBV2Client(d *schema.ResourceData, config *Config) (*gophercloud.ServiceClient, error) {
	if config.UseOctavia {
		return config.LoadBalancerV2Client(GetRegion(d, config))
	}
	return config.NetworkingV2Client(GetRegion(d, config))
}

// chooseLBV2AccTestClient will determine which load balacing client to use:
// either the Octavia/LBaaS client or the Neutron/Networking v2 client.
// This is similar to the chooseLBV2Client function but specific for acceptance
// tests.
func chooseLBV2AccTestClient(config *Config, region string) (*gophercloud.ServiceClient, error) {
	if config.UseOctavia {
		return config.LoadBalancerV2Client(region)
	}
	return config.NetworkingV2Client(region)
}

// chooseLBV2LoadbalancerUpdateOpts will determine which load balancer update options to use:
// either the Octavia/LBaaS or the Neutron/Networking v2.
func chooseLBV2LoadbalancerUpdateOpts(d *schema.ResourceData, config *Config) (neutronloadbalancers.UpdateOptsBuilder, error) {
	var hasChange bool

	if config.UseOctavia {
		// Use Octavia.
		var updateOpts octavialoadbalancers.UpdateOpts

		if d.HasChange("name") {
			hasChange = true
			name := d.Get("name").(string)
			updateOpts.Name = &name
		}
		if d.HasChange("description") {
			hasChange = true
			description := d.Get("description").(string)
			updateOpts.Description = &description
		}
		if d.HasChange("admin_state_up") {
			hasChange = true
			asu := d.Get("admin_state_up").(bool)
			updateOpts.AdminStateUp = &asu
		}

		if d.HasChange("tags") {
			hasChange = true
			if v, ok := d.GetOk("tags"); ok {
				tags := v.(*schema.Set).List()
				tagsToUpdate := expandToStringSlice(tags)
				updateOpts.Tags = &tagsToUpdate
			} else {
				updateOpts.Tags = &[]string{}
			}
		}

		if hasChange {
			return updateOpts, nil
		}
	}

	// Use Neutron.
	var updateOpts neutronloadbalancers.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("admin_state_up") {
		hasChange = true
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	if hasChange {
		return updateOpts, nil
	}

	return nil, nil
}

func waitForLBV2LoadBalancer(ctx context.Context, lbClient *gophercloud.ServiceClient, lbID string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for loadbalancer %s to become %s.", lbID, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID),
		Timeout:    timeout,
		Delay:      0,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: loadbalancer %s not found: %s", lbID, err)
			}
		}
		return fmt.Errorf("Error waiting for loadbalancer %s to become %s: %s", lbID, target, err)
	}

	return nil
}

func resourceLBV2LoadBalancerRefreshFunc(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := neutronloadbalancers.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return lb, lb.ProvisioningStatus, nil
	}
}

func flattenLBMembersV2(members []octaviapools.Member) []map[string]interface{} {
	m := make([]map[string]interface{}, len(members))

	for i, member := range members {
		m[i] = map[string]interface{}{
			"name":            member.Name,
			"weight":          member.Weight,
			"admin_state_up":  member.AdminStateUp,
			"subnet_id":       member.SubnetID,
			"address":         member.Address,
			"protocol_port":   member.ProtocolPort,
			"monitor_port":    member.MonitorPort,
			"monitor_address": member.MonitorAddress,
			"id":              member.ID,
			"backup":          member.Backup,
		}
	}

	return m
}

func expandLBMembersV2(members *schema.Set, lbClient *gophercloud.ServiceClient) []octaviapools.BatchUpdateMemberOpts {
	var m []octaviapools.BatchUpdateMemberOpts

	if members != nil {
		for _, raw := range members.List() {
			rawMap := raw.(map[string]interface{})
			name := rawMap["name"].(string)
			subnetID := rawMap["subnet_id"].(string)
			weight := rawMap["weight"].(int)
			adminStateUp := rawMap["admin_state_up"].(bool)

			member := octaviapools.BatchUpdateMemberOpts{
				Address:      rawMap["address"].(string),
				ProtocolPort: rawMap["protocol_port"].(int),
				Name:         &name,
				SubnetID:     &subnetID,
				Weight:       &weight,
				AdminStateUp: &adminStateUp,
			}

			// backup requires octavia minor version 2.1. Only set when specified
			if val, ok := rawMap["backup"]; ok {
				backup := val.(bool)
				member.Backup = &backup
			}

			// Only set monitor_port and monitor_address when explicitly specified, as they are optional arguments
			if val, ok := rawMap["monitor_port"]; ok {
				monitorPort := val.(int)
				if monitorPort > 0 {
					member.MonitorPort = &monitorPort
				}
			}

			if val, ok := rawMap["monitor_address"]; ok {
				monitorAddress := val.(string)
				if monitorAddress != "" {
					member.MonitorAddress = &monitorAddress
				}
			}

			m = append(m, member)
		}
	}

	return m
}

func resourceLoadBalancerV2SetSecurityGroups(networkingClient *gophercloud.ServiceClient, vipPortID string, d *schema.ResourceData) error {
	if vipPortID != "" {
		if v, ok := d.GetOk("security_group_ids"); ok {
			securityGroups := expandToStringSlice(v.(*schema.Set).List())
			updateOpts := ports.UpdateOpts{
				SecurityGroups: &securityGroups,
			}

			log.Printf("[DEBUG] Adding security groups to openstack_lb_loadbalancer_v2 "+
				"VIP port %s: %#v", vipPortID, updateOpts)

			_, err := ports.Update(networkingClient, vipPortID, updateOpts).Extract()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceLoadBalancerV2GetSecurityGroups(networkingClient *gophercloud.ServiceClient, vipPortID string, d *schema.ResourceData) error {
	port, err := ports.Get(networkingClient, vipPortID).Extract()
	if err != nil {
		return err
	}

	d.Set("security_group_ids", port.SecurityGroups)

	return nil
}
