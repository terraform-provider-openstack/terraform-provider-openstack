package openstack

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// so far only "read" is supported
var aclOperations = []string{"read"}

var aclSchema = &schema.Schema{
	Type:     schema.TypeList, // the list, returned by Barbican, is always ordered
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true, // defaults to true in OpenStack Barbican code
			},
			"users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	},
}

func keyManagerSecretV1WaitForSecretDeletion(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := secrets.Delete(kmClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if _, ok := err.(gophercloud.ErrDefault404); ok {
			return "", "DELETED", nil
		}

		return nil, "ACTIVE", err
	}
}

func keyManagerSecretV1SecretType(v string) secrets.SecretType {
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

func keyManagerSecretV1WaitForSecretCreation(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		secret, err := secrets.Get(kmClient, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return "", "NOT_CREATED", nil
			}

			return "", "NOT_CREATED", err
		}

		if secret.Status == "ERROR" {
			return "", secret.Status, fmt.Errorf("Error creating secret")
		}

		return secret, secret.Status, nil
	}
}

func keyManagerSecretV1GetUUIDfromSecretRef(ref string) string {
	// secret ref has form https://{barbican_host}/v1/secrets/{secret_uuid}
	// so we are only interested in the last part
	ref_split := strings.Split(ref, "/")
	uuid := ref_split[len(ref_split)-1]
	return uuid
}

func flattenKeyManagerSecretV1Metadata(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func keyManagerSecretMetadataV1WaitForSecretMetadataCreation(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		metadata, err := secrets.GetMetadata(kmClient, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return "", "NOT_CREATED", nil
			}

			return "", "NOT_CREATED", err
		}
		return metadata, "ACTIVE", nil
	}
}

func keyManagerSecretV1GetPayload(kmClient *gophercloud.ServiceClient, id string) string {
	payload, err := secrets.GetPayload(kmClient, id, nil).Extract()
	if err != nil {
		fmt.Errorf("Could not retrieve payload for secret with id %s: %s", id, err)
	}
	return string(payload)
}

func resourceSecretV1PayloadBase64CustomizeDiff(diff *schema.ResourceDiff) error {
	encoding := diff.Get("payload_content_encoding").(string)
	if diff.Id() != "" && diff.HasChange("payload") && encoding == "base64" {
		o, n := diff.GetChange("payload")
		oldPayload := o.(string)
		newPayload := n.(string)

		v, err := base64.StdEncoding.DecodeString(newPayload)
		if err != nil {
			return fmt.Errorf("The Payload is not in the defined base64 format: %s", err)
		}
		newPayloadDecoded := string(v)

		if oldPayload == newPayloadDecoded {
			log.Printf("[DEBUG] payload has not changed. clearing diff")
			return diff.Clear("payload")
		}
	}

	return nil
}

func expandKeyManagerV1ACLsRaw(v interface{}, aclType string) acls.SetOpts {
	var res acls.SetOpts

	if v, ok := v.(map[string]interface{}); ok {
		if v, ok := v[aclType]; ok {
			if v, ok := v.([]interface{}); ok {
				for _, v := range v {
					if v, ok := v.(map[string]interface{}); ok {
						if v, ok := v["project_access"]; ok {
							if v, ok := v.(bool); ok {
								res.ProjectAccess = &v
							}
						}
						if v, ok := v["users"]; ok {
							if v, ok := v.(*schema.Set); ok {
								for _, v := range v.List() {
									if res.Users == nil {
										users := []string{}
										res.Users = &users
									}
									*res.Users = append(*res.Users, v.(string))
								}
							}
						}
					}
				}
			}
		}
	}
	return res
}

func expandKeyManagerV1ACLs(v interface{}, aclType string) acls.SetOpts {
	var res acls.SetOpts
	users := []string{}
	iTrue := true // set default value to true
	res.ProjectAccess = &iTrue
	res.Type = aclType

	raw := expandKeyManagerV1ACLsRaw(v, aclType)
	if raw.ProjectAccess != nil {
		res.ProjectAccess = raw.ProjectAccess
	}
	if raw.Users != nil {
		res.Users = raw.Users
	} else {
		res.Users = &users
	}

	return res
}

func flattenKeyManagerV1ACLs(acl *acls.ACL) []map[string][]map[string]interface{} {
	var m []map[string][]map[string]interface{}

	if acl != nil {
		allAcls := *acl
		for _, aclOp := range aclOperations {
			if v, ok := allAcls[aclOp]; ok {
				if m == nil {
					m = make([]map[string][]map[string]interface{}, 1)
					m[0] = make(map[string][]map[string]interface{})
				}
				if m[0][aclOp] == nil {
					m[0][aclOp] = make([]map[string]interface{}, 1)
				}
				m[0][aclOp][0] = map[string]interface{}{
					"project_access": v.ProjectAccess,
					"users":          v.Users,
					"created_at":     v.Created.UTC().Format(time.RFC3339),
					"updated_at":     v.Updated.UTC().Format(time.RFC3339),
				}
			}
		}
	}

	return m
}
