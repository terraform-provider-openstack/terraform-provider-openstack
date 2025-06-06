package openstack

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/transfer/request"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// TransferRequestCreateOpts represents the attributes used when creating a new transfer request.
type TransferRequestCreateOpts struct {
	request.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToTransferRequestCreateMap casts a CreateOpts struct to a map.
// It overrides request.ToTransferRequestCreateMap to add the ValueSpecs field.
func (opts TransferRequestCreateOpts) ToTransferRequestCreateMap() (map[string]any, error) {
	b, err := BuildRequest(opts, "")
	if err != nil {
		return nil, err
	}

	if m, ok := b[""].(map[string]any); ok {
		return m, nil
	}

	return nil, fmt.Errorf("Expected map but got %T", b[""])
}

func dnsTransferRequestV2RefreshFunc(ctx context.Context, dnsClient *gophercloud.ServiceClient, transferRequestID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		transferRequest, err := request.Get(ctx, dnsClient, transferRequestID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return transferRequest, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] openstack_dns_transfer_request_v2 %s current status: %s", transferRequest.ID, transferRequest.Status)

		return transferRequest, transferRequest.Status, nil
	}
}
