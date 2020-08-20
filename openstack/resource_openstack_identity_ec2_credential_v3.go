package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v2/tenants"
	tokens2 "github.com/gophercloud/gophercloud/openstack/identity/v2/tokens"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/ec2credentials"
	tokens3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIdentityEc2CredentialV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityEc2CredentialV3Create,
		Read:   resourceIdentityEc2CredentialV3Read,
		Delete: resourceIdentityEc2CredentialV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"access": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},

			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
				ForceNew:  true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"trust_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIdentityEc2CredentialV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	user, project, err := GetTokenInfo(identityClient)
	if err != nil {
		return fmt.Errorf("Error getting token info: %s", err)
	}

	var tenantID string
	// Set tenant to user value if provided, get one from token otherwise
	if definedProject, ok := d.GetOk("project_id"); ok {
		tenantID = definedProject.(string)
	} else {
		tenantID = project
	}

	createOpts := ec2credentials.CreateOpts{
		TenantID: tenantID,
	}

	var userID string
	// Set userid to defined tf value, use one from token otherwise
	if definedUser, ok := d.GetOk("user_id"); ok {
		userID = definedUser.(string)
	} else {
		userID = user
	}

	log.Printf("[DEBUG] openstack_identity_ec2_credential_v3 create options: %#v", createOpts)

	ec2Credential, err := ec2credentials.Create(identityClient, userID, createOpts).Extract()

	if err != nil {
		if v, ok := err.(gophercloud.ErrDefault404); ok {
			return fmt.Errorf("Error creating openstack_identity_ec2_credential_v3: %s", v.ErrUnexpectedResponseCode.Body)
		}
		return fmt.Errorf("Error creating openstack_identity_ec2_credential_v3: %s", err)
	}

	d.SetId(ec2Credential.Access)

	d.Set("access", ec2Credential.Access)
	d.Set("secret", ec2Credential.Secret)
	d.Set("user_id", ec2Credential.UserID)
	d.Set("project_id", ec2Credential.TenantID)
	d.Set("trust_id", ec2Credential.TrustID)

	return resourceIdentityEc2CredentialV3Read(d, meta)
}

func resourceIdentityEc2CredentialV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	user, _, err := GetTokenInfo(identityClient)
	if err != nil {
		return fmt.Errorf("Error getting token info: %s", err)
	}

	var userID string
	// Set userid to defined tf value, use one from token otherwise
	if definedUser, ok := d.GetOk("user_id"); ok {
		userID = definedUser.(string)
	} else {
		userID = user
	}

	ec2Credential, err := ec2credentials.Get(identityClient, userID, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_identity_ec2_credential_v3")
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_ec2_credential_v3 %s: %#v", d.Id(), ec2Credential)

	d.Set("secret", ec2Credential.Secret)
	d.Set("user_id", ec2Credential.UserID)
	d.Set("project_id", ec2Credential.TenantID)
	d.Set("access", ec2Credential.Access)
	d.Set("trust_id", ec2Credential.TrustID)

	return nil
}

func resourceIdentityEc2CredentialV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	user, _, err := GetTokenInfo(identityClient)
	if err != nil {
		return fmt.Errorf("Error getting token info: %s", err)
	}

	var userID string
	// Set userid to defined tf value, use one from token otherwise
	if definedUser, ok := d.GetOk("user_id"); ok {
		userID = definedUser.(string)
	} else {
		userID = user
	}

	err = ec2credentials.Delete(identityClient, userID, d.Id()).ExtractErr()
	if err != nil {
		err = CheckDeleted(d, err, "Error deleting openstack_identity_ec2_credential_v3")
		if err != nil {
			return err
		}
	}
	return nil
}

func GetTokenInfo(sc *gophercloud.ServiceClient) (string, string, error) {
	r := sc.ProviderClient.GetAuthResult()
	switch r := r.(type) {
	case tokens2.CreateResult:
		user, project, err := GetTokenInfoV2(r)
		if err != nil {
			return "", "", err
		}
		return user, project, nil
	case tokens3.CreateResult:
		user, project, err := GetTokenInfoV3(r)
		if err != nil {
			return "", "", err
		}
		return user, project, nil
	default:
		return "", "", fmt.Errorf("got unexpected AuthResult type %t", r)
	}
}

func GetTokenInfoV3(t tokens3.CreateResult) (string, string, error) {
	user, err := t.ExtractUser()
	if err != nil {
		return "", "", err
	}
	project, err := t.ExtractProject()
	if err != nil {
		return "", "", err
	}
	return user.ID, project.ID, nil
}

func GetTokenInfoV2(t tokens2.CreateResult) (string, string, error) {
	var tokeninfo struct {
		Access struct {
			Token struct {
				Expires string         `json:"expires"`
				ID      string         `json:"id"`
				Tenant  tenants.Tenant `json:"tenant"`
			} `json:"token"`
			User tokens2.User `json:"user"`
		} `json:"access"`
	}

	err := t.ExtractInto(&tokeninfo)
	if err != nil {
		return "", "", err
	}
	return tokeninfo.Access.User.ID, tokeninfo.Access.Token.ID, nil
}
