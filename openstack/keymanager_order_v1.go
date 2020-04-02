package openstack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/orders"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func keyManagerOrderV1WaitForOrderDeletion(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := orders.Delete(kmClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if _, ok := err.(gophercloud.ErrDefault404); ok {
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

func keyManagerOrderV1WaitForOrderCreation(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		order, err := orders.Get(kmClient, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return "", "NOT_CREATED", nil
			}

			return "", "NOT_CREATED", err
		}

		if order.Status == "ERROR" {
			return "", order.Status, fmt.Errorf("Error creating order")
		}

		return order, order.Status, nil
	}
}

func keyManagerOrderV1GetUUIDfromOrderRef(ref string) string {
	// order ref has form https://{barbican_host}/v1/orders/{order_uuid}
	// so we are only interested in the last part
	ref_split := strings.Split(ref, "/")
	uuid := ref_split[len(ref_split)-1]
	return uuid
}

func expandKeyManagerOrderV1Meta(m map[string]interface{}) orders.MetaOpts {
	var meta orders.MetaOpts

	if v, ok := m["algorithm"]; ok {
		meta.Algorithm = v.(string)
	}

	if v, ok := m["bit_length"]; ok {
		i, _ := strconv.Atoi(v.(string))
		meta.BitLength = i
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

func flattenKeyManagerOrderV1Meta(m orders.Meta) map[string]interface{} {
	meta := make(map[string]interface{})

	if m.Algorithm != "" {
		meta["algorithm"] = m.Algorithm
	}

	if m.BitLength != 0 {
		meta["bit_length"] = strconv.Itoa(m.BitLength)
	}

	if !m.Expiration.IsZero() {
		meta["expiration"] = m.Expiration.UTC().Format(time.RFC3339)
	}

	if m.Mode != "" {
		meta["mode"] = m.Mode
	}

	if m.Name != "" {
		meta["name"] = m.Name
	}

	if m.PayloadContentType != "" {
		meta["payload_content_type"] = m.PayloadContentType

	}

	return meta
}
