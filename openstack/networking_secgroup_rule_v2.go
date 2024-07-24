package openstack

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func resourceNetworkingSecGroupRuleV2StateRefreshFunc(client *gophercloud.ServiceClient, sgRuleID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		sgRule, err := rules.Get(client, sgRuleID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return sgRule, "DELETED", nil
			}

			return sgRule, "", err
		}

		return sgRule, "ACTIVE", nil
	}
}

func resourceNetworkingSecGroupRuleV2Direction(v interface{}, k string) ([]string, []error) {
	switch v.(string) {
	case string(rules.DirIngress):
		return nil, nil
	case string(rules.DirEgress):
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_networking_secgroup_rule_v2", v, k)}
}

func resourceNetworkingSecGroupRuleV2EtherType(v interface{}, k string) ([]string, []error) {
	switch v.(string) {
	case string(rules.EtherType4):
		return nil, nil
	case string(rules.EtherType6):
		return nil, nil
	}

	return nil, []error{fmt.Errorf("unknown %q %s for openstack_networking_secgroup_rule_v2s", v, k)}
}

func resourceNetworkingSecGroupRuleV2Protocol(v interface{}, k string) ([]string, []error) {
	switch v.(string) {
	case string(rules.ProtocolAH),
		string(rules.ProtocolDCCP),
		string(rules.ProtocolEGP),
		string(rules.ProtocolESP),
		string(rules.ProtocolGRE),
		string(rules.ProtocolICMP),
		string(rules.ProtocolIGMP),
		string(rules.ProtocolIPv6Encap),
		string(rules.ProtocolIPv6Frag),
		string(rules.ProtocolIPv6ICMP),
		string(rules.ProtocolIPv6NoNxt),
		string(rules.ProtocolIPv6Opts),
		string(rules.ProtocolIPv6Route),
		string(rules.ProtocolOSPF),
		string(rules.ProtocolPGM),
		string(rules.ProtocolRSVP),
		string(rules.ProtocolSCTP),
		string(rules.ProtocolTCP),
		string(rules.ProtocolUDP),
		string(rules.ProtocolUDPLite),
		string(rules.ProtocolVRRP),
		"":
		return nil, nil
	}

	// If the protocol wasn't matched above, see if it's an integer.
	p, err := strconv.Atoi(v.(string))
	if err != nil {
		return nil, []error{fmt.Errorf("unknown %q %s for openstack_networking_secgroup_rule_v2: %s", v, k, err)}
	}
	if p < 0 || p > 255 {
		return nil, []error{fmt.Errorf("unknown %q %s for openstack_networking_secgroup_rule_v2", v, k)}
	}

	return nil, nil
}
