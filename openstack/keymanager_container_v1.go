package openstack

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/containers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func keyManagerContainerV1WaitForContainerDeletion(ctx context.Context, kmClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		err := containers.Delete(ctx, kmClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return "", "DELETED", nil
		}

		return nil, "ACTIVE", err
	}
}

func keyManagerContainerV1Type(v string) containers.ContainerType {
	var ctype containers.ContainerType

	switch v {
	case "generic":
		ctype = containers.GenericContainer
	case "rsa":
		ctype = containers.RSAContainer
	case "certificate":
		ctype = containers.CertificateContainer
	}

	return ctype
}

func keyManagerContainerV1WaitForContainerCreation(ctx context.Context, kmClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		container, err := containers.Get(ctx, kmClient, id).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return "", "NOT_CREATED", nil
			}

			return "", "NOT_CREATED", err
		}

		if container.Status == "ERROR" {
			return "", container.Status, errors.New("Error creating container")
		}

		return container, container.Status, nil
	}
}

func keyManagerContainerV1GetUUIDfromContainerRef(ref string) string {
	// container ref has form https://{barbican_host}/v1/containers/{container_uuid}
	// so we are only interested in the last part
	refSplit := strings.Split(ref, "/")
	uuid := refSplit[len(refSplit)-1]

	return uuid
}

func expandKeyManagerContainerV1SecretRefs(secretRefs *schema.Set) []containers.SecretRef {
	l := make([]containers.SecretRef, 0, len(secretRefs.List()))

	for _, v := range secretRefs.List() {
		if v, ok := v.(map[string]any); ok {
			var s containers.SecretRef

			if v, ok := v["secret_ref"]; ok {
				s.SecretRef = v.(string)
			}

			if v, ok := v["name"]; ok {
				s.Name = v.(string)
			}

			l = append(l, s)
		}
	}

	return l
}

func flattenKeyManagerContainerV1SecretRefs(sr []containers.SecretRef) []map[string]any {
	m := make([]map[string]any, 0, len(sr))

	for _, v := range sr {
		m = append(m, map[string]any{
			"name":       v.Name,
			"secret_ref": v.SecretRef,
		})
	}

	return m
}

func flattenKeyManagerContainerV1Consumers(cr []containers.ConsumerRef) []map[string]any {
	m := make([]map[string]any, 0, len(cr))

	for _, v := range cr {
		m = append(m, map[string]any{
			"name": v.Name,
			"url":  v.URL,
		})
	}

	return m
}
