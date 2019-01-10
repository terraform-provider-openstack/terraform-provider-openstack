package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

// networkingFloatingIPV2ID retrieves floating IP ID by the provided IP address.
func networkingFloatingIPV2ID(client *gophercloud.ServiceClient, floatingIP string) (string, error) {
	listOpts := floatingips.ListOpts{
		FloatingIP: floatingIP,
	}

	allPages, err := floatingips.List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return "", err
	}

	if len(allFloatingIPs) == 0 {
		return "", fmt.Errorf("there are no openstack_networking_floatingip_v2 with %s IP", floatingIP)
	}
	if len(allFloatingIPs) > 1 {
		return "", fmt.Errorf("there are more than one openstack_networking_floatingip_v2 with %s IP", floatingIP)
	}

	return allFloatingIPs[0].ID, nil
}
