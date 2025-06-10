package openstack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type floatingIPExtended struct {
	floatingips.FloatingIP
	dns.FloatingIPDNSExt
}

// networkingFloatingIPV2ID retrieves floating IP ID by the provided IP address.
func networkingFloatingIPV2ID(ctx context.Context, client *gophercloud.ServiceClient, floatingIP string) (string, error) {
	listOpts := floatingips.ListOpts{
		FloatingIP: floatingIP,
	}

	allPages, err := floatingips.List(client, listOpts).AllPages(ctx)
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

func networkingFloatingIPV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, fipID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		fip, err := floatingips.Get(ctx, client, fipID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return fip, "DELETED", nil
			}

			return nil, "", err
		}

		return fip, fip.Status, nil
	}
}
