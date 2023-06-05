package openstack

import (
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
)

func expandFWRuleV2Action(a string) rules.Action {
	var action rules.Action
	switch strings.ToLower(a) {
	case "allow":
		action = rules.ActionAllow
	case "deny":
		action = rules.ActionDeny
	case "reject":
		action = rules.ActionReject
	}

	return action
}

func expandFWRuleV2IPVersion(ipv int) gophercloud.IPVersion {
	var ipVersion gophercloud.IPVersion
	switch ipv {
	case 4:
		ipVersion = gophercloud.IPv4
	case 6:
		ipVersion = gophercloud.IPv6
	}

	return ipVersion
}

func expandFWRuleV2Protocol(p string) rules.Protocol {
	var protocol rules.Protocol
	switch strings.ToLower(p) {
	case "any":
		protocol = rules.ProtocolAny
	case "icmp":
		protocol = rules.ProtocolICMP
	case "tcp":
		protocol = rules.ProtocolTCP
	case "udp":
		protocol = rules.ProtocolUDP
	}

	return protocol
}
