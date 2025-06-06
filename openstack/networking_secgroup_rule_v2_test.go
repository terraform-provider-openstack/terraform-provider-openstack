package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/stretchr/testify/assert"
)

func TestUnitResourceNetworkingSecGroupRuleV2DirectionIngress(t *testing.T) {
	_, err := resourceNetworkingSecGroupRuleV2Direction("ingress", "direction")

	assert.Empty(t, err)
}

func TestUnitResourceNetworkingSecGroupRuleV2DirectionEgress(t *testing.T) {
	_, err := resourceNetworkingSecGroupRuleV2Direction("egress", "direction")

	assert.Empty(t, err)
}

func TestUnitResourceNetworkingSecGroupRuleV2DirectionUnknown(t *testing.T) {
	_, err := resourceNetworkingSecGroupRuleV2Direction("stuff", "direction")

	assert.Len(t, err, 1)
}

func TestUnitResourceNetworkingSecGroupRuleV2EtherType4(t *testing.T) {
	_, err := resourceNetworkingSecGroupRuleV2EtherType("IPv4", "ethertype")

	assert.Empty(t, err)
}

func TestUnitResourceNetworkingSecGroupRuleV2EtherType6(t *testing.T) {
	_, err := resourceNetworkingSecGroupRuleV2EtherType("IPv6", "ethertype")

	assert.Empty(t, err)
}

func TestUnitResourceNetworkingSecGroupRuleV2EtherTypeUnknown(t *testing.T) {
	_, err := resourceNetworkingSecGroupRuleV2EtherType("something", "ethertype")

	assert.Len(t, err, 1)
}

func TestUnitResourceNetworkingSecGroupRuleV2ProtocolString(t *testing.T) {
	protocols := map[string]rules.RuleProtocol{
		string(rules.ProtocolAH):        rules.ProtocolAH,
		string(rules.ProtocolDCCP):      rules.ProtocolDCCP,
		string(rules.ProtocolESP):       rules.ProtocolESP,
		string(rules.ProtocolEGP):       rules.ProtocolEGP,
		string(rules.ProtocolGRE):       rules.ProtocolGRE,
		string(rules.ProtocolICMP):      rules.ProtocolICMP,
		string(rules.ProtocolIGMP):      rules.ProtocolIGMP,
		string(rules.ProtocolIPv6Encap): rules.ProtocolIPv6Encap,
		string(rules.ProtocolIPv6Frag):  rules.ProtocolIPv6Frag,
		string(rules.ProtocolIPv6ICMP):  rules.ProtocolIPv6ICMP,
		string(rules.ProtocolIPv6NoNxt): rules.ProtocolIPv6NoNxt,
		string(rules.ProtocolIPv6Opts):  rules.ProtocolIPv6Opts,
		string(rules.ProtocolIPv6Route): rules.ProtocolIPv6Route,
		string(rules.ProtocolOSPF):      rules.ProtocolOSPF,
		string(rules.ProtocolPGM):       rules.ProtocolPGM,
		string(rules.ProtocolRSVP):      rules.ProtocolRSVP,
		string(rules.ProtocolSCTP):      rules.ProtocolSCTP,
		string(rules.ProtocolTCP):       rules.ProtocolTCP,
		string(rules.ProtocolUDP):       rules.ProtocolUDP,
		string(rules.ProtocolUDPLite):   rules.ProtocolUDPLite,
		string(rules.ProtocolVRRP):      rules.ProtocolVRRP,
		"":                              "",
	}

	for name := range protocols {
		_, err := resourceNetworkingSecGroupRuleV2Protocol(name, "protocol")

		assert.Empty(t, err)
	}
}

func TestUnitResourceNetworkingSecGroupRuleV2ProtocolNumber(t *testing.T) {
	_, err := resourceNetworkingSecGroupRuleV2Protocol("6", "protocol")

	assert.Empty(t, err)
}
