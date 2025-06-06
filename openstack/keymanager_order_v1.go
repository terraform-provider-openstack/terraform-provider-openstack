package openstack

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/orders"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func keyManagerOrderV1WaitForOrderDeletion(ctx context.Context, kmClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		err := orders.Delete(ctx, kmClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return "", "DELETED", nil
		}

		return nil, "ACTIVE", err
	}
}

func keyManagerOrderV1OrderType(v string) orders.OrderType {
	var otype orders.OrderType

	switch v {
	case "asymmetric":
		otype = orders.AsymmetricOrder
	case "key":
		otype = orders.KeyOrder
	}

	return otype
}

func keyManagerOrderV1WaitForOrderCreation(ctx context.Context, kmClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		order, err := orders.Get(ctx, kmClient, id).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return "", "NOT_CREATED", nil
			}

			return "", "NOT_CREATED", err
		}

		if order.Status == "ERROR" {
			return "", order.Status, errors.New("Error creating order")
		}

		return order, order.Status, nil
	}
}

func keyManagerOrderV1GetUUIDfromOrderRef(ref string) string {
	// order ref has form https://{barbican_host}/v1/orders/{order_uuid}
	// so we are only interested in the last part
	refSplit := strings.Split(ref, "/")
	uuid := refSplit[len(refSplit)-1]

	return uuid
}

func expandKeyManagerOrderV1Meta(s []any) orders.MetaOpts {
	var meta orders.MetaOpts

	m := s[0].(map[string]any)

	if v, ok := m["algorithm"]; ok {
		meta.Algorithm = v.(string)
	}

	if v, ok := m["bit_length"]; ok {
		meta.BitLength = v.(int)
	}

	if v, ok := m["expiration"]; ok {
		if t, _ := time.Parse(time.RFC3339, v.(string)); t != (time.Time{}) {
			meta.Expiration = &t
		}
	}

	if v, ok := m["mode"]; ok {
		meta.Mode = v.(string)
	}

	if v, ok := m["name"]; ok {
		meta.Name = v.(string)
	}

	if v, ok := m["payload_content_type"]; ok {
		meta.PayloadContentType = v.(string)
	}

	return meta
}

func flattenKeyManagerOrderV1Meta(m orders.Meta) []map[string]any {
	var meta []map[string]any

	s := make(map[string]any)

	if m.Algorithm != "" {
		s["algorithm"] = m.Algorithm
	}

	if m.BitLength != 0 {
		s["bit_length"] = m.BitLength
	}

	if !m.Expiration.IsZero() {
		s["expiration"] = m.Expiration.UTC().Format(time.RFC3339)
	}

	if m.Mode != "" {
		s["mode"] = m.Mode
	}

	if m.Name != "" {
		s["name"] = m.Name
	}

	if m.PayloadContentType != "" {
		s["payload_content_type"] = m.PayloadContentType
	}

	return append(meta, s)
}
