package openstack

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/transfer/accept"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// TransferAcceptCreateOpts represents the attributes used when creating a new transfer accept.
type TransferAcceptCreateOpts struct {
	accept.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToTransferAcceptCreateMap casts a CreateOpts struct to a map.
// It overrides accept.ToTransferAcceptCreateMap to add the ValueSpecs field.
func (opts TransferAcceptCreateOpts) ToTransferAcceptCreateMap() (map[string]any, error) {
	b, err := BuildRequest(opts, "")
	if err != nil {
		return nil, err
	}

	if m, ok := b[""].(map[string]any); ok {
		return m, nil
	}

	return nil, fmt.Errorf("Expected map but got %T", b[""])
}

func dnsTransferAcceptV2RefreshFunc(ctx context.Context, dnsClient *gophercloud.ServiceClient, transferAcceptID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		transferAccept, err := accept.Get(ctx, dnsClient, transferAcceptID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return transferAccept, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] openstack_dns_transfer_accept_v2 %s current status: %s", transferAccept.ID, transferAccept.Status)

		return transferAccept, transferAccept.Status, nil
	}
}
