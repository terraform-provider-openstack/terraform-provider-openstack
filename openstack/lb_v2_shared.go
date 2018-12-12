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

func waitForLBV2Listener(lbClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for listener %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2ListenerRefreshFunc(lbClient, id),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: listener %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("Error waiting for listener %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceLBV2ListenerRefreshFunc(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		listener, err := listeners.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		// The listener resource has no Status attribute, so a successful Get is the best we can do
		return listener, "ACTIVE", nil
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

func waitForLBV2Member(lbClient *gophercloud.ServiceClient, poolID, memberID string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for member %s to become %s.", memberID, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2MemberRefreshFunc(lbClient, poolID, memberID),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: member %s not found: %s", memberID, err)
			}
		}
		return fmt.Errorf("Error waiting for member %s to become %s: %s", memberID, target, err)
	}

	return nil
}

func resourceLBV2MemberRefreshFunc(lbClient *gophercloud.ServiceClient, poolID, memberID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		member, err := pools.GetMember(lbClient, poolID, memberID).Extract()
		if err != nil {
			return nil, "", err
		}

		// The member resource has no Status attribute, so a successful Get is the best we can do
		return member, "ACTIVE", nil
	}
}

func waitForLBV2Monitor(lbClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for monitor %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2MonitorRefreshFunc(lbClient, id),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: monitor %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("Error waiting for monitor %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceLBV2MonitorRefreshFunc(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		monitor, err := monitors.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		// The monitor resource has no Status attribute, so a successful Get is the best we can do
		return monitor, "ACTIVE", nil
	}
}

func waitForLBV2Pool(lbClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for pool %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2PoolRefreshFunc(lbClient, id),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: pool %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("Error waiting for pool %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceLBV2PoolRefreshFunc(lbClient *gophercloud.ServiceClient, poolID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pool, err := pools.Get(lbClient, poolID).Extract()
		if err != nil {
			return nil, "", err
		}

		// The pool resource has no Status attribute, so a successful Get is the best we can do
		return pool, "ACTIVE", nil
	}
}

func waitForLBV2viaPool(lbClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	pool, err := pools.Get(lbClient, id).Extract()
	if err != nil {
		return err
	}

	lbID := ""
	listenerID := ""
	if len(pool.Loadbalancers) > 0 {
		lbID = pool.Loadbalancers[0].ID
	}
	if len(pool.Listeners) > 0 {
		listenerID = pool.Listeners[0].ID
	}

	return waitForLBV2viaLBorListener(lbClient, lbID, listenerID, target, pending, timeout)
}

func waitForLBV2viaLBorListener(lbClient *gophercloud.ServiceClient, lbID string, listenerID string, target string, pending []string, timeout time.Duration) error {
	if lbID != "" {
		return waitForLBV2LoadBalancer(lbClient, lbID, target, pending, timeout)
	}
	if listenerID != "" {
		return waitForLBV2viaListener(lbClient, listenerID, target, pending, timeout)
	}
	return fmt.Errorf("Neither Load Balancer ID nor Listener ID were provided")
}

func waitForLBV2viaListener(lbClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	listener, err := listeners.Get(lbClient, id).Extract()
	if err != nil {
		return err
	}
	if len(listener.Loadbalancers) > 0 {
		lbID := listener.Loadbalancers[0].ID
		return waitForLBV2LoadBalancer(lbClient, lbID, target, pending, timeout)
	}

	// got a listener but no LB - this is wrong
	return fmt.Errorf("No Load Balancer on listener %s", id)
}

func waitForLBV2L7policy(lbClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for l7policy %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceLBV2L7policyRefreshFunc(lbClient, id),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: l7policy %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("Error waiting for l7policy %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceLBV2L7policyRefreshFunc(lbClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		l7policy, err := l7policies.Get(lbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		// The l7policy resource has no Status attribute, so a successful Get is the best we can do
		return l7policy, "ACTIVE", nil
	}
}
