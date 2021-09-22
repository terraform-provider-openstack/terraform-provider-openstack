package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func TestResourceNetworkingSecGroupRuleV2DirectionIngress(t *testing.T) {
	expected := rules.DirIngress

	actual, err := resourceNetworkingSecGroupRuleV2Direction("ingress")

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestResourceNetworkingSecGroupRuleV2DirectionEgress(t *testing.T) {
	expected := rules.DirEgress

	actual, err := resourceNetworkingSecGroupRuleV2Direction("egress")

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestResourceNetworkingSecGroupRuleV2DirectionUnknown(t *testing.T) {
	actual, err := resourceNetworkingSecGroupRuleV2Direction("stuff")

	assert.Error(t, err)
	assert.Empty(t, actual)
}

func TestResourceNetworkingSecGroupRuleV2EtherType4(t *testing.T) {
	expected := rules.EtherType4

	actual, err := resourceNetworkingSecGroupRuleV2EtherType("IPv4")

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestResourceNetworkingSecGroupRuleV2EtherType6(t *testing.T) {
	expected := rules.EtherType6

	actual, err := resourceNetworkingSecGroupRuleV2EtherType("IPv6")

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestResourceNetworkingSecGroupRuleV2EtherTypeUnknown(t *testing.T) {
	actual, err := resourceNetworkingSecGroupRuleV2EtherType("something")

	assert.Error(t, err)
	assert.Empty(t, actual)
}

func TestResourceNetworkingSecGroupRuleV2ProtocolString(t *testing.T) {
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
	}

	for name, expected := range protocols {
		actual, err := resourceNetworkingSecGroupRuleV2Protocol(name)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestResourceNetworkingSecGroupRuleV2ProtocolNumber(t *testing.T) {
	expected := rules.RuleProtocol("6")

	actual, err := resourceNetworkingSecGroupRuleV2Protocol("6")

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
