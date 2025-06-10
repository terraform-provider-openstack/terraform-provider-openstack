package openstack

import (
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/acls"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// So far only "read" is supported.
func getSupportedACLOperations() [1]string {
	return [1]string{"read"}
}

func getACLSchema() *schema.Schema {
	return &schema.Schema{
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
}

func expandKeyManagerV1ACL(v any, aclType string) acls.SetOpt {
	users := []string{}
	iTrue := true // set default value to true
	res := acls.SetOpt{
		ProjectAccess: &iTrue,
		Users:         &users,
		Type:          aclType,
	}

	if v, ok := v.([]any); ok {
		for _, v := range v {
			if v, ok := v.(map[string]any); ok {
				if v, ok := v["project_access"]; ok {
					if v, ok := v.(bool); ok {
						res.ProjectAccess = &v
					}
				}

				if v, ok := v["users"]; ok {
					if v, ok := v.(*schema.Set); ok {
						for _, v := range v.List() {
							*res.Users = append(*res.Users, v.(string))
						}
					}
				}
			}
		}
	}

	return res
}

func expandKeyManagerV1ACLs(v any) acls.SetOpts {
	var res []acls.SetOpt

	if v, ok := v.([]any); ok {
		for _, v := range v {
			if v, ok := v.(map[string]any); ok {
				for aclType, v := range v {
					acl := expandKeyManagerV1ACL(v, aclType)
					res = append(res, acl)
				}
			}
		}
	}

	return res
}

func flattenKeyManagerV1ACLs(acl *acls.ACL) []map[string][]map[string]any {
	var m []map[string][]map[string]any

	if acl != nil {
		allAcls := *acl
		for _, aclOp := range getSupportedACLOperations() {
			if v, ok := allAcls[aclOp]; ok {
				if m == nil {
					m = make([]map[string][]map[string]any, 1)
					m[0] = make(map[string][]map[string]any)
				}

				if m[0][aclOp] == nil {
					m[0][aclOp] = make([]map[string]any, 1)
				}

				m[0][aclOp][0] = map[string]any{
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
