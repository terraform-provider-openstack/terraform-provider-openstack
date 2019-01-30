package openstack

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/hashicorp/terraform/helper/resource"
)

func keymanagerSecretMetadataV1WaitForSecretMetadataCreation(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		fmt.Println("ID is %v", id)
		metadata, err := secrets.GetMetadata(kmClient, id).Extract()
		if err != nil {
			return "", "NOT_CREATED", nil
		}
		return metadata, "ACTIVE", nil
	}
}

// SecretMetadataCreateOpts represents the attributes used when creating a new Barbican secret.
type SecretMetadataCreateOpts struct {
	secrets.MetadataOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}
