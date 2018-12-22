package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/monitors"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/pools"
)

// lbPendingStatuses are the valid statuses a LoadBalancer will be in while
// it's updating.
var lbPendingStatuses = []string{"PENDING_CREATE", "PENDING_UPDATE"}

// lbPendingDeleteStatuses are the valid statuses a LoadBalancer will be before delete
var lbPendingDeleteStatuses = []string{"ERROR", "PENDING_UPDATE", "PENDING_DELETE", "ACTIVE"}

var lbSkipLBStatuses = []string{"ERROR", "ACTIVE"}

// chooseLBV2Client will determine which load balacing client to use:
// Either the Octavia/LBaaS client or the Neutron/Networking v2 client.
func chooseLBV2Client(d *schema.ResourceData, config *Config) (*gophercloud.ServiceClient, error) {
	if config.useOctavia {
		return config.loadBalancerV2Client(GetRegion(d, config))
	}
	return config.networkingV2Client(GetRegion(d, config))
}

// chooseLBV2AccTestClient will determine which load balacing client to use:
// Either the Octavia/LBaaS client or the Neutron/Networking v2 client.
// This is similar to the chooseLBV2Client function but specific for acceptance
// tests.
func chooseLBV2AccTestClient(config *Config, region string) (*gophercloud.ServiceClient, error) {
	if config.useOctavia {
		return config.loadBalancerV2Client(region)
	}
	return config.networkingV2Client(region)
}

func waitForLBV2Listener(lbClient *gophercloud.ServiceClient, listener *listeners.Listener, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for listener %s to become %s.", listener.ID, target)

	var refreshFunc resource.StateRefreshFunc

	if len(listener.Loadbalancers) == 0 {
		return fmt.Errorf("Failed to detect a Load Balancer for the %s Listener", listener.ID)
	}

	lbID := listener.Loadbalancers[0].ID

	if listener.ProvisioningStatus != "" {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, resourceLBV2ListenerRefreshFunc(lbClient, listener.ID))
	} else {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "listener", listener.ID)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    refreshFunc,
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("Error waiting for listener %s to become %s: %s", listener.ID, target, err)
	}

	return nil
}

func resourceLBV2ListenerRefreshFunc(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := listeners.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return lb, lb.ProvisioningStatus, nil
	}
}

func waitForLBV2LoadBalancer(lbClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for loadbalancer %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2LoadBalancerRefreshFunc(lbClient, id),
		Timeout:    timeout,
		Delay:      0,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: loadbalancer %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("Error waiting for loadbalancer %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceLBV2LoadBalancerRefreshFunc(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := loadbalancers.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return lb, lb.ProvisioningStatus, nil
	}
}

func waitForLBV2Member(lbClient *gophercloud.ServiceClient, parentPool *pools.Pool, member *pools.Member, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for member %s to become %s.", member.ID, target)

	var refreshFunc resource.StateRefreshFunc

	lbID, err := lbV2FindLBIDviaPool(lbClient, parentPool)
	if err != nil {
		return err
	}

	if member.ProvisioningStatus != "" {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, resourceLBV2MemberRefreshFunc(lbClient, parentPool.ID, member.ID))
	} else {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "member", member.ID)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    refreshFunc,
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("Error waiting for member %s to become %s: %s", member.ID, target, err)
	}

	return nil
}

func resourceLBV2MemberRefreshFunc(lbClient *gophercloud.ServiceClient, poolID, memberID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		member, err := pools.GetMember(lbClient, poolID, memberID).Extract()
		if err != nil {
			return nil, "", err
		}

		return member, member.ProvisioningStatus, nil
	}
}

func waitForLBV2Monitor(lbClient *gophercloud.ServiceClient, parentPool *pools.Pool, monitor *monitors.Monitor, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for monitor %s to become %s.", monitor.ID, target)

	var refreshFunc resource.StateRefreshFunc

	lbID, err := lbV2FindLBIDviaPool(lbClient, parentPool)
	if err != nil {
		return err
	}

	if monitor.ProvisioningStatus != "" {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, resourceLBV2MonitorRefreshFunc(lbClient, monitor.ID))
	} else {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "monitor", monitor.ID)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    refreshFunc,
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}
		return fmt.Errorf("Error waiting for monitor %s to become %s: %s", monitor.ID, target, err)
	}

	return nil
}

func resourceLBV2MonitorRefreshFunc(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		monitor, err := monitors.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return monitor, monitor.ProvisioningStatus, nil
	}
}

func waitForLBV2Pool(lbClient *gophercloud.ServiceClient, pool *pools.Pool, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for pool %s to become %s.", pool.ID, target)

	var refreshFunc resource.StateRefreshFunc

	lbID, err := lbV2FindLBIDviaPool(lbClient, pool)
	if err != nil {
		return err
	}

	if pool.ProvisioningStatus != "" {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, resourceLBV2PoolRefreshFunc(lbClient, pool.ID))
	} else {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "pool", pool.ID)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    refreshFunc,
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("Error waiting for pool %s to become %s: %s", pool.ID, target, err)
	}

	return nil
}

func resourceLBV2PoolRefreshFunc(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pool, err := pools.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return pool, pool.ProvisioningStatus, nil
	}
}

func lbV2FindLBIDviaPool(lbClient *gophercloud.ServiceClient, pool *pools.Pool) (string, error) {
	if len(pool.Loadbalancers) > 0 {
		return pool.Loadbalancers[0].ID, nil
	}

	if len(pool.Listeners) > 0 {
		listenerID := pool.Listeners[0].ID
		listener, err := listeners.Get(lbClient, listenerID).Extract()
		if err != nil {
			return "", err
		}

		if len(listener.Loadbalancers) > 0 {
			return listener.Loadbalancers[0].ID, nil
		}
	}

	return "", fmt.Errorf("Unable to determine loadbalancer ID from pool %s", pool.ID)
}

func resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient *gophercloud.ServiceClient, id string, resourceRefreshFunc func() (interface{}, string, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		refreshFunc := resourceLBV2LoadBalancerRefreshFunc(lbClient, id)
		lb, status, err := refreshFunc()
		if err != nil {
			return lb, status, err
		}
		if !strSliceContains(lbSkipLBStatuses, status) {
			return lb, status, nil
		}
		return resourceRefreshFunc()
	}
}

func resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient *gophercloud.ServiceClient, lbID, resourceType, resourceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		statuses, err := loadbalancers.GetStatuses(lbClient, lbID).Extract()
		if err != nil {
			return nil, "", fmt.Errorf("Unable to get statuses from the Load Balancer %s statuses tree: %s", lbID, err)
		}

		if !strSliceContains(lbSkipLBStatuses, statuses.Loadbalancer.ProvisioningStatus) {
			return statuses.Loadbalancer, statuses.Loadbalancer.ProvisioningStatus, nil
		}

		switch resourceType {
		case "listener":
			for _, listener := range statuses.Loadbalancer.Listeners {
				if listener.ID == resourceID {
					if listener.ProvisioningStatus != "" {
						return listener, listener.ProvisioningStatus, nil
					}
				}
			}
			listener, err := listeners.Get(lbClient, resourceID).Extract()
			return listener, "ACTIVE", err

		case "pool":
			for _, pool := range statuses.Loadbalancer.Pools {
				if pool.ID == resourceID {
					if pool.ProvisioningStatus != "" {
						return pool, pool.ProvisioningStatus, nil
					}
				}
			}
			pool, err := pools.Get(lbClient, resourceID).Extract()
			return pool, "ACTIVE", err

		case "monitor":
			for _, pool := range statuses.Loadbalancer.Pools {
				if pool.Monitor.ID == resourceID {
					if pool.Monitor.ProvisioningStatus != "" {
						return pool.Monitor, pool.Monitor.ProvisioningStatus, nil
					}
				}
			}
			monitor, err := monitors.Get(lbClient, resourceID).Extract()
			return monitor, "ACTIVE", err

		case "member":
			for _, pool := range statuses.Loadbalancer.Pools {
				for _, member := range pool.Members {
					if member.ID == resourceID {
						if member.ProvisioningStatus != "" {
							return member, member.ProvisioningStatus, nil
						}
					}
				}
			}
			return "", "DELETED", nil

		case "l7policy":
			for _, listener := range statuses.Loadbalancer.Listeners {
				for _, l7policy := range listener.L7Policies {
					if l7policy.ID == resourceID {
						if l7policy.ProvisioningStatus != "" {
							return l7policy, l7policy.ProvisioningStatus, nil
						}
					}
				}
			}
			l7policy, err := l7policies.Get(lbClient, resourceID).Extract()
			return l7policy, "ACTIVE", err

		case "l7rule":
			for _, listener := range statuses.Loadbalancer.Listeners {
				for _, l7policy := range listener.L7Policies {
					for _, l7rule := range l7policy.Rules {
						if l7rule.ID == resourceID {
							if l7rule.ProvisioningStatus != "" {
								return l7rule, l7rule.ProvisioningStatus, nil
							}
						}
					}
				}
			}
			return "", "DELETED", nil
		}

		return nil, "", fmt.Errorf("An unexpected error occurred querying the status of %s %s by loadbalancer %s", resourceType, resourceID, lbID)
	}
}

func resourceLBV2L7PolicyRefreshFunc(lbClient *gophercloud.ServiceClient, l7policyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		l7policy, err := l7policies.Get(lbClient, l7policyID).Extract()
		if err != nil {
			return nil, "", err
		}

		return l7policy, l7policy.ProvisioningStatus, nil
	}
}

func waitForLBV2L7Policy(lbClient *gophercloud.ServiceClient, parentListener *listeners.Listener, l7policy *l7policies.L7Policy, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7policy %s to become %s.", l7policy.ID, target)

	var refreshFunc resource.StateRefreshFunc

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("Unable to determine loadbalancer ID from listener %s", parentListener.ID)
	}

	lbID := parentListener.Loadbalancers[0].ID

	if l7policy.ProvisioningStatus != "" {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, resourceLBV2L7PolicyRefreshFunc(lbClient, l7policy.ID))
	} else {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "l7policy", l7policy.ID)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    refreshFunc,
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("Error waiting for l7policy %s to become %s: %s", l7policy.ID, target, err)
	}

	return nil
}

func getListenerIDForL7Policy(lbClient *gophercloud.ServiceClient, id string) (string, error) {
	log.Printf("[DEBUG] Trying to get Listener ID associated with the %s L7 Policy ID", id)
	lbsPages, err := loadbalancers.List(lbClient, loadbalancers.ListOpts{}).AllPages()
	if err != nil {
		return "", fmt.Errorf("No Load Balancers were found: %s", err)
	}

	lbs, err := loadbalancers.ExtractLoadBalancers(lbsPages)
	if err != nil {
		return "", fmt.Errorf("Unable to extract Load Balancers list: %s", err)
	}

	for _, lb := range lbs {
		statuses, err := loadbalancers.GetStatuses(lbClient, lb.ID).Extract()
		if err != nil {
			return "", fmt.Errorf("Failed to get Load Balancer statuses: %s", err)
		}
		for _, listener := range statuses.Loadbalancer.Listeners {
			for _, l7policy := range listener.L7Policies {
				if l7policy.ID == id {
					return listener.ID, nil
				}
			}
		}
	}

	return "", fmt.Errorf("Unable to find Listener ID associated with the %s L7 Policy ID", id)
}

func resourceLBV2L7RuleRefreshFunc(lbClient *gophercloud.ServiceClient, l7policyID string, l7ruleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		l7rule, err := l7policies.GetRule(lbClient, l7policyID, l7ruleID).Extract()
		if err != nil {
			return nil, "", err
		}

		return l7rule, l7rule.ProvisioningStatus, nil
	}
}

func waitForLBV2L7Rule(lbClient *gophercloud.ServiceClient, parentListener *listeners.Listener, parentL7policy *l7policies.L7Policy, l7rule *l7policies.Rule, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7rule %s to become %s.", l7rule.ID, target)

	var refreshFunc resource.StateRefreshFunc

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("Unable to determine loadbalancer ID from listener %s", parentListener.ID)
	}

	lbID := parentListener.Loadbalancers[0].ID

	if l7rule.ProvisioningStatus != "" {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, resourceLBV2L7RuleRefreshFunc(lbClient, parentL7policy.ID, l7rule.ID))
	} else {
		refreshFunc = resourceLBV2LoadBalancerStatusRefreshFuncNeutron(lbClient, lbID, "l7rule", l7rule.ID)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    refreshFunc,
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("Error waiting for l7rule %s to become %s: %s", l7rule.ID, target, err)
	}

	return nil
}

// strSliceContains checks if a given string is contained in a slice
// When anybody asks why Go needs generics, here you go.
func strSliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
