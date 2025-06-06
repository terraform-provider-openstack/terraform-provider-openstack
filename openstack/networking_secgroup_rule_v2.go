package openstack

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func resourceNetworkingSecGroupRuleV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, sgRuleID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		sgRule, err := rules.Get(ctx, client, sgRuleID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return sgRule, "DELETED", nil
			}

			return sgRule, "", err
		}

		return sgRule, "ACTIVE", nil
	}
}

func resourceNetworkingSecGroupRuleV2Direction(v any, k string) ([]string, []error) {
	switch rules.RuleDirection(v.(string)) {
	case rules.DirIngress:
		return nil, nil
	case rules.DirEgress:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_networking_secgroup_rule_v2", v, k)}
}

func resourceNetworkingSecGroupRuleV2EtherType(v any, k string) ([]string, []error) {
	switch rules.RuleEtherType(v.(string)) {
	case rules.EtherType4:
		return nil, nil
	case rules.EtherType6:
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_networking_secgroup_rule_v2s", v, k)}
}

func resourceNetworkingSecGroupRuleV2Protocol(v any, k string) ([]string, []error) {
	switch rules.RuleProtocol(v.(string)) {
	case rules.ProtocolAH,
		rules.ProtocolDCCP,
		rules.ProtocolEGP,
		rules.ProtocolESP,
		rules.ProtocolGRE,
		rules.ProtocolICMP,
		rules.ProtocolIGMP,
		rules.ProtocolIPv6Encap,
		rules.ProtocolIPv6Frag,
		rules.ProtocolIPv6ICMP,
		rules.ProtocolIPv6NoNxt,
		rules.ProtocolIPv6Opts,
		rules.ProtocolIPv6Route,
		rules.ProtocolOSPF,
		rules.ProtocolPGM,
		rules.ProtocolRSVP,
		rules.ProtocolSCTP,
		rules.ProtocolTCP,
		rules.ProtocolUDP,
		rules.ProtocolUDPLite,
		rules.ProtocolVRRP,
		rules.ProtocolIPIP,
		rules.ProtocolAny:
		return nil, nil
	}

	// If the protocol wasn't matched above, see if it's an integer.
	p, err := strconv.Atoi(v.(string))
	if err != nil {
		return nil, []error{fmt.Errorf("unknown %q %s for openstack_networking_secgroup_rule_v2: %w", v, k, err)}
	}

	if p < 0 || p > 255 {
		return nil, []error{fmt.Errorf("unknown %q %s for openstack_networking_secgroup_rule_v2", v, k)}
	}

	return nil, nil
}
