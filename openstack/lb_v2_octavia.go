package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/l7policies"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/monitors"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
)

func getListenerIDForL7PolicyOctavia(lbClient *gophercloud.ServiceClient, id string) (string, error) {
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

func waitForLBV2L7RuleOctavia(ctx context.Context, lbClient *gophercloud.ServiceClient, parentListener *listeners.Listener, parentL7policy *l7policies.L7Policy, l7rule *l7policies.Rule, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7rule %s to become %s.", l7rule.ID, target)

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("Unable to determine loadbalancer ID from listener %s", parentListener.ID)
	}

	lbID := parentListener.Loadbalancers[0].ID

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2L7RuleRefreshFuncOctavia(lbClient, lbID, parentL7policy.ID, l7rule),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
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

func resourceLBV2L7RuleRefreshFuncOctavia(lbClient *gophercloud.ServiceClient, lbID string, l7policyID string, l7rule *l7policies.Rule) resource.StateRefreshFunc {
	if l7rule.ProvisioningStatus == "" {
		return resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, "l7rule", l7rule.ID, l7policyID)
	}

	return func() (interface{}, string, error) {
		lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
		if err != nil {
			return lb, status, err
		}
		if !strSliceContains(getLbSkipStatuses(), status) {
			return lb, status, nil
		}

		l7rule, err := l7policies.GetRule(lbClient, l7policyID, l7rule.ID).Extract()
		if err != nil {
			return nil, "", err
		}

		return l7rule, l7rule.ProvisioningStatus, nil
	}
}

func resourceLBV2ListenerRefreshFuncOctavia(lbClient *gophercloud.ServiceClient, lbID string, listener *listeners.Listener) resource.StateRefreshFunc {
	if listener.ProvisioningStatus == "" {
		return resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, "listener", listener.ID, "")
	}

	return func() (interface{}, string, error) {
		lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
		if err != nil {
			return lb, status, err
		}
		if !strSliceContains(getLbSkipStatuses(), status) {
			return lb, status, nil
		}

		listener, err := listeners.Get(lbClient, listener.ID).Extract()
		if err != nil {
			return nil, "", err
		}

		return listener, listener.ProvisioningStatus, nil
	}
}

func waitForLBV2ListenerOctavia(ctx context.Context, lbClient *gophercloud.ServiceClient, listener *listeners.Listener, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for openstack_lb_listener_v2 %s to become %s.", listener.ID, target)

	if len(listener.Loadbalancers) == 0 {
		return fmt.Errorf("Failed to detect a openstack_lb_loadbalancer_v2 for the %s openstack_lb_listener_v2", listener.ID)
	}

	lbID := listener.Loadbalancers[0].ID

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2ListenerRefreshFuncOctavia(lbClient, lbID, listener),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}

		return fmt.Errorf("Error waiting for openstack_lb_listener_v2 %s to become %s: %s", listener.ID, target, err)
	}

	return nil
}

func lbV2FindLBIDviaPoolOctavia(lbClient *gophercloud.ServiceClient, pool *pools.Pool) (string, error) {
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

func resourceLBV2PoolRefreshFuncOctavia(lbClient *gophercloud.ServiceClient, lbID string, pool *pools.Pool) resource.StateRefreshFunc {
	if pool.ProvisioningStatus == "" {
		return resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, "pool", pool.ID, "")
	}

	return func() (interface{}, string, error) {
		lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
		if err != nil {
			return lb, status, err
		}
		if !strSliceContains(getLbSkipStatuses(), status) {
			return lb, status, nil
		}

		pool, err := pools.Get(lbClient, pool.ID).Extract()
		if err != nil {
			return nil, "", err
		}

		return pool, pool.ProvisioningStatus, nil
	}
}

func waitForLBV2PoolOctavia(ctx context.Context, lbClient *gophercloud.ServiceClient, pool *pools.Pool, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for pool %s to become %s.", pool.ID, target)

	lbID, err := lbV2FindLBIDviaPoolOctavia(lbClient, pool)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2PoolRefreshFuncOctavia(lbClient, lbID, pool),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
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

func resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient *gophercloud.ServiceClient, lbID, resourceType, resourceID string, parentID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		statuses, err := loadbalancers.GetStatuses(lbClient, lbID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nil, "", gophercloud.ErrDefault404{
					ErrUnexpectedResponseCode: gophercloud.ErrUnexpectedResponseCode{
						BaseError: gophercloud.BaseError{
							DefaultErrString: fmt.Sprintf("Unable to get statuses from the Load Balancer %s statuses tree: %s", lbID, err),
						},
					},
				}
			}
			return nil, "", fmt.Errorf("Unable to get statuses from the Load Balancer %s statuses tree: %s", lbID, err)
		}

		// Don't fail, when statuses returns "null"
		if statuses == nil || statuses.Loadbalancer == nil {
			statuses = new(loadbalancers.StatusTree)
			statuses.Loadbalancer = new(loadbalancers.LoadBalancer)
		} else if !strSliceContains(getLbSkipStatuses(), statuses.Loadbalancer.ProvisioningStatus) {
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
			member, err := pools.GetMember(lbClient, parentID, resourceID).Extract()
			return member, "ACTIVE", err

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
			l7Rule, err := l7policies.GetRule(lbClient, parentID, resourceID).Extract()
			return l7Rule, "ACTIVE", err
		}

		return nil, "", fmt.Errorf("An unexpected error occurred querying the status of %s %s by loadbalancer %s", resourceType, resourceID, lbID)
	}
}

func waitForLBV2L7PolicyOctavia(ctx context.Context, lbClient *gophercloud.ServiceClient, parentListener *listeners.Listener, l7policy *l7policies.L7Policy, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7policy %s to become %s.", l7policy.ID, target)

	if len(parentListener.Loadbalancers) == 0 {
		return fmt.Errorf("Unable to determine loadbalancer ID from listener %s", parentListener.ID)
	}

	lbID := parentListener.Loadbalancers[0].ID

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2L7PolicyRefreshFuncOctavia(lbClient, lbID, l7policy),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
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

func resourceLBV2L7PolicyRefreshFuncOctavia(lbClient *gophercloud.ServiceClient, lbID string, l7policy *l7policies.L7Policy) resource.StateRefreshFunc {
	if l7policy.ProvisioningStatus == "" {
		return resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, "l7policy", l7policy.ID, "")
	}

	return func() (interface{}, string, error) {
		lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
		if err != nil {
			return lb, status, err
		}
		if !strSliceContains(getLbSkipStatuses(), status) {
			return lb, status, nil
		}

		l7policy, err := l7policies.Get(lbClient, l7policy.ID).Extract()
		if err != nil {
			return nil, "", err
		}

		return l7policy, l7policy.ProvisioningStatus, nil
	}
}

func waitForLBV2MemberOctavia(ctx context.Context, lbClient *gophercloud.ServiceClient, parentPool *pools.Pool, member *pools.Member, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for member %s to become %s.", member.ID, target)

	lbID, err := lbV2FindLBIDviaPoolOctavia(lbClient, parentPool)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2MemberRefreshFuncOctavia(lbClient, lbID, parentPool.ID, member),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
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

func resourceLBV2MemberRefreshFuncOctavia(lbClient *gophercloud.ServiceClient, lbID string, poolID string, member *pools.Member) resource.StateRefreshFunc {
	if member.ProvisioningStatus == "" {
		return resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, "member", member.ID, poolID)
	}

	return func() (interface{}, string, error) {
		lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
		if err != nil {
			return lb, status, err
		}
		if !strSliceContains(getLbSkipStatuses(), status) {
			return lb, status, nil
		}

		member, err := pools.GetMember(lbClient, poolID, member.ID).Extract()
		if err != nil {
			return nil, "", err
		}

		return member, member.ProvisioningStatus, nil
	}
}

func waitForLBV2MonitorOctavia(ctx context.Context, lbClient *gophercloud.ServiceClient, parentPool *pools.Pool, monitor *monitors.Monitor, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for openstack_lb_monitor_v2 %s to become %s.", monitor.ID, target)

	lbID, err := lbV2FindLBIDviaPoolOctavia(lbClient, parentPool)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2MonitorRefreshFuncOctavia(lbClient, lbID, monitor),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
		}
		return fmt.Errorf("Error waiting for openstack_lb_monitor_v2 %s to become %s: %s", monitor.ID, target, err)
	}

	return nil
}

func resourceLBV2MonitorRefreshFuncOctavia(lbClient *gophercloud.ServiceClient, lbID string, monitor *monitors.Monitor) resource.StateRefreshFunc {
	if monitor.ProvisioningStatus == "" {
		return resourceLBV2LoadBalancerStatusRefreshFuncOctavia(lbClient, lbID, "monitor", monitor.ID, "")
	}

	return func() (interface{}, string, error) {
		lb, status, err := resourceLBV2LoadBalancerRefreshFunc(lbClient, lbID)()
		if err != nil {
			return lb, status, err
		}
		if !strSliceContains(getLbSkipStatuses(), status) {
			return lb, status, nil
		}

		monitor, err := monitors.Get(lbClient, monitor.ID).Extract()
		if err != nil {
			return nil, "", err
		}

		return monitor, monitor.ProvisioningStatus, nil
	}
}

func flattenLBPoolPersistenceV2(p pools.SessionPersistence) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":        p.Type,
			"cookie_name": p.CookieName,
		},
	}
}

func waitForLBV2LoadBalancerOctavia(ctx context.Context, lbClient *gophercloud.ServiceClient, lbID string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for loadbalancer %s to become %s.", lbID, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2LoadBalancerRefreshFuncOctavia(lbClient, lbID),
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

func resourceLBV2LoadBalancerRefreshFuncOctavia(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := loadbalancers.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return lb, lb.ProvisioningStatus, nil
	}
}

func expandLBV2ListenerHeadersMap(raw map[string]interface{}) (map[string]string, error) {
	m := make(map[string]string, len(raw))
	for key, val := range raw {
		labelValue, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("label %s value should be string", key)
		}

		m[key] = labelValue
	}

	return m, nil
}
