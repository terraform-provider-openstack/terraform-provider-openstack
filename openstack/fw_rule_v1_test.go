package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas/rules"
)

func TestExpandFWRuleV1IPVersion(t *testing.T) {
	ipv := 4

	expected := gophercloud.IPv4
	actual := expandFWRuleV1IPVersion(ipv)
	assert.Equal(t, expected, actual)
}

func TestExpandFWRuleV1Protocol(t *testing.T) {
	proto := "tcp"

	expected := rules.ProtocolTCP
	actual := expandFWRuleV1Protocol(proto)
	assert.Equal(t, expected, actual)
}
