package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentityUserMembershipV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityUserMembershipV3Create,
		ReadContext:   resourceIdentityUserMembershipV3Read,
		DeleteContext: resourceIdentityUserMembershipV3Delete,
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
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIdentityUserMembershipV3Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	userID := d.Get("user_id").(string)
	groupID := d.Get("group_id").(string)

	if err := users.AddToGroup(ctx, identityClient, groupID, userID).ExtractErr(); err != nil {
		return diag.Errorf("Error creating openstack_identity_user_membership_v3: %s", err)
	}

	id := fmt.Sprintf("%s/%s", userID, groupID)
	d.SetId(id)

	return resourceIdentityUserMembershipV3Read(ctx, d, meta)
}

func resourceIdentityUserMembershipV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	userID, groupID, err := parsePairedIDs(d.Id(), "openstack_identity_user_membership_v3")
	if err != nil {
		return diag.FromErr(err)
	}

	userMembership, err := users.IsMemberOfGroup(ctx, identityClient, groupID, userID).Extract()
	if err != nil || !userMembership {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_identity_user_membership_v3"))
	}

	d.Set("region", GetRegion(d, config))
	d.Set("user_id", userID)
	d.Set("group_id", groupID)

	return nil
}

func resourceIdentityUserMembershipV3Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	userID, groupID, err := parsePairedIDs(d.Id(), "openstack_identity_user_membership_v3")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := users.RemoveFromGroup(ctx, identityClient, groupID, userID).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error removing openstack_identity_user_membership_v3"))
	}

	return nil
}
