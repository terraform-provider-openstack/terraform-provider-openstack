package openstack

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/hashicorp/terraform/helper/resource"
	"strings"
)

func keymanagerSecretV1WaitForSecretDeletion(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := secrets.Delete(kmClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		return nil, "ACTIVE", err
	}
}

func keymanagerSecretV1SecretType(v string) secrets.SecretType {
	var stype secrets.SecretType
	switch v {
	case "symmetric":
		stype = secrets.SymmetricSecret
	case "public":
		stype = secrets.PublicSecret
	case "private":
		stype = secrets.PrivateSecret
	case "passphrase":
		stype = secrets.PassphraseSecret
	case "certificate":
		stype = secrets.CertificateSecret
	case "opaque":
		stype = secrets.OpaqueSecret
	}

	return stype
}

func keymanagerSecretV1WaitForSecretCreation(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		fmt.Println("ID is %v", id)
		secret, err := secrets.Get(kmClient, id).Extract()
		if err != nil {
			return "", "NOT_CREATED", nil
		}
		return secret, "ACTIVE", nil
	}
}

func keymanagerSecretV1GetUUIDfromSecretRef(ref string) string {
	// secret ref has form https://{barbican_host}/v1/secrets/{secret_uuid}
	// so we are only interested in the last part
	ref_split := strings.Split(ref, "/")
	uuid := ref_split[len(ref_split)-1]
	return uuid
}

// SecretCreateOpts represents the attributes used when creating a new Barbican secret.
type SecretCreateOpts struct {
	secrets.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}
