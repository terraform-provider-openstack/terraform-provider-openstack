package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/osinherit"
)

func resourceIdentityInheritRoleAssignmentV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityInheritRoleAssignmentV3Create,
		ReadContext:   resourceIdentityInheritRoleAssignmentV3Read,
		DeleteContext: resourceIdentityInheritRoleAssignmentV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"domain_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"project_id"},
				Optional:      true,
				ForceNew:      true,
			},

			"group_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"user_id"},
				Optional:      true,
				ForceNew:      true,
			},

			"project_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"domain_id"},
				Optional:      true,
				ForceNew:      true,
			},

			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"user_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"group_id"},
				Optional:      true,
				ForceNew:      true,
			},
		},
	}
}

func resourceIdentityInheritRoleAssignmentV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	roleID := d.Get("role_id").(string)
	domainID := d.Get("domain_id").(string)
	groupID := d.Get("group_id").(string)
	projectID := d.Get("project_id").(string)
	userID := d.Get("user_id").(string)

	opts := osinherit.AssignOpts{
		DomainID:  domainID,
		GroupID:   groupID,
		ProjectID: projectID,
		UserID:    userID,
	}

	err = osinherit.Assign(identityClient, roleID, opts).ExtractErr()
	if err != nil {
		return diag.Errorf("Error creating openstack_identity_inherit_role_assignment_v3: %s", err)
	}

	id := identityRoleAssignmentV3ID(domainID, projectID, groupID, userID, roleID)
	d.SetId(id)

	return resourceIdentityInheritRoleAssignmentV3Read(ctx, d, meta)
}

func resourceIdentityInheritRoleAssignmentV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	domainID, projectID, groupID, userID, roleID, err := identityRoleAssignmentV3ParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error determining openstack_identity_inherit_role_assignment_v3 ID: %s", err)
	}

	validateOpts := osinherit.ValidateOpts{
		DomainID:  domainID,
		ProjectID: projectID,
		UserID:    userID,
		GroupID:   groupID,
	}

	err = osinherit.Validate(identityClient, roleID, validateOpts).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error vaalidatin openstack_identity_inherit_role_assignment_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_inherit_role_assignment_v3 %s", d.Id())
	d.Set("domain_id", domainID)
	d.Set("project_id", projectID)
	d.Set("group_id", groupID)
	d.Set("user_id", userID)
	d.Set("role_id", roleID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceIdentityInheritRoleAssignmentV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	domainID, projectID, groupID, userID, roleID, err := identityRoleAssignmentV3ParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error determining openstack_identity_inherit_role_assignment_v3 ID: %s", err)
	}

	opts := osinherit.UnassignOpts{
		DomainID:  domainID,
		GroupID:   groupID,
		ProjectID: projectID,
		UserID:    userID,
	}

	if err := osinherit.Unassign(identityClient, roleID, opts).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error unassigning openstack_identity_inherit_role_assignment_v3"))
	}

	return nil
}
