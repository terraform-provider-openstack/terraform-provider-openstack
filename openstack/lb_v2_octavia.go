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
