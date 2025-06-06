package openstack

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/members"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceImagesImageAccessAcceptV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImagesImageAccessAcceptV2Create,
		ReadContext:   resourceImagesImageAccessAcceptV2Read,
		UpdateContext: resourceImagesImageAccessAcceptV2Update,
		DeleteContext: resourceImagesImageAccessAcceptV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImagesImageAccessAcceptV2Import,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"member_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"accepted", "rejected", "pending",
				}, false),
			},

			// Computed-only
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceImagesImageAccessAcceptV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	imageClient, err := config.ImageV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID := d.Get("image_id").(string)
	memberID := d.Get("member_id").(string)
	status := d.Get("status").(string)

	if memberID == "" {
		memberID, err = resourceImagesImageAccessV2DetectMemberID(ctx, imageClient, imageID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// accept status on the consumer side
	opts := members.UpdateOpts{
		Status: status,
	}

	_, err = members.Update(ctx, imageClient, imageID, memberID, opts).Extract()
	if err != nil {
		return diag.Errorf("Error setting a member status to the %q image share for the %q member: %s", imageID, memberID, err)
	}

	id := fmt.Sprintf("%s/%s", imageID, memberID)
	d.SetId(id)

	return resourceImagesImageAccessAcceptV2Read(ctx, d, meta)
}

func resourceImagesImageAccessAcceptV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	imageClient, err := config.ImageV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID, memberID, err := parsePairedIDs(d.Id(), "openstack_images_image_access_accept_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	member, err := members.Get(ctx, imageClient, imageID, memberID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving the openstack_images_image_access_accept_v2"))
	}

	log.Printf("[DEBUG] Retrieved Image member %s: %#v", d.Id(), member)

	d.Set("status", member.Status)
	d.Set("image_id", member.ImageID)
	d.Set("member_id", member.MemberID)
	// Computed
	d.Set("schema", member.Schema)
	d.Set("created_at", member.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", member.UpdatedAt.Format(time.RFC3339))
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceImagesImageAccessAcceptV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	imageClient, err := config.ImageV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID, memberID, err := parsePairedIDs(d.Id(), "openstack_images_image_access_accept_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	status := d.Get("status").(string)

	opts := members.UpdateOpts{
		Status: status,
	}

	_, err = members.Update(ctx, imageClient, imageID, memberID, opts).Extract()
	if err != nil {
		return diag.Errorf("Error updateing the %q image with the %q member: %s", imageID, memberID, err)
	}

	return resourceImagesImageAccessAcceptV2Read(ctx, d, meta)
}

func resourceImagesImageAccessAcceptV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	imageClient, err := config.ImageV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	imageID, memberID, err := parsePairedIDs(d.Id(), "openstack_images_image_access_accept_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Rejecting Image membership %s", d.Id())
	// reject status on the consumer side
	opts := members.UpdateOpts{
		Status: "rejected",
	}
	if _, err := members.Update(ctx, imageClient, imageID, memberID, opts).Extract(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error rejecting the openstack_images_image_access_accept_v2"))
	}

	return nil
}

func resourceImagesImageAccessAcceptV2Import(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)

	config := meta.(*Config)

	imageClient, err := config.ImageV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack image client: %w", err)
	}

	imageID := parts[0]

	var memberID string
	if len(parts) > 1 {
		memberID = parts[1]
	} else {
		memberID, err = resourceImagesImageAccessV2DetectMemberID(ctx, imageClient, imageID)
		if err != nil {
			return nil, err
		}
	}

	id := fmt.Sprintf("%s/%s", imageID, memberID)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
