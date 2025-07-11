package openstack

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-provider-openstack/utils/v2/hashcode"
)

func dataSourceImagesImageIDsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImagesImageIDsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_regex"},
			},

			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(images.ImageVisibilityPublic),
					string(images.ImageVisibilityPrivate),
					string(images.ImageVisibilityShared),
					string(images.ImageVisibilityCommunity),
				}, false),
			},

			"member_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(images.ImageMemberStatusAccepted),
					string(images.ImageMemberStatusPending),
					string(images.ImageMemberStatusRejected),
					string(images.ImageMemberStatusAll),
				}, false),
			},

			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"size_min": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"size_max": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"sort": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "name:asc",
				ValidateFunc: dataSourceValidateImageSortFilter,
			},

			"tag": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"hidden": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},

			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				ConflictsWith: []string{"name"},
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"container_format": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"disk_format": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// dataSourceImagesImageIDsV2Read performs the image lookup.
func dataSourceImagesImageIDsV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	imageClient, err := config.ImageV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	visibility := resourceImagesImageV2VisibilityFromString(
		d.Get("visibility").(string))
	memberStatus := resourceImagesImageV2MemberStatusFromString(
		d.Get("member_status").(string))
	properties := resourceImagesImageV2ExpandProperties(
		d.Get("properties").(map[string]any))

	tags := []string{}
	tagList := d.Get("tags").(*schema.Set).List()

	for _, v := range tagList {
		tags = append(tags, fmt.Sprint(v))
	}

	tag := d.Get("tag").(string)
	if tag != "" {
		tags = append(tags, tag)
	}

	listOpts := images.ListOpts{
		Name:            d.Get("name").(string),
		Visibility:      visibility,
		Hidden:          d.Get("hidden").(bool),
		Owner:           d.Get("owner").(string),
		Status:          images.ImageStatusActive,
		SizeMin:         int64(d.Get("size_min").(int)),
		SizeMax:         int64(d.Get("size_max").(int)),
		Sort:            d.Get("sort").(string),
		ContainerFormat: d.Get("container_format").(string),
		DiskFormat:      d.Get("disk_format").(string),
		Tags:            tags,
		MemberStatus:    memberStatus,
	}

	log.Printf("[DEBUG] List Options in openstack_images_image_ids_v2: %#v", listOpts)

	allPages, err := images.List(imageClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to list images in openstack_images_image_ids_v2: %s", err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve images in openstack_images_image_ids_v2: %s", err)
	}

	log.Printf("[DEBUG] Retrieved %d images in openstack_images_image_ids_v2: %+v", len(allImages), allImages)

	allImages = imagesFilterByProperties(allImages, properties)

	log.Printf("[DEBUG] Image list filtered by properties: %#v", properties)

	nameRegex, nameRegexOk := d.GetOk("name_regex")
	if nameRegexOk {
		allImages = imagesFilterByRegex(allImages, nameRegex.(string))

		log.Printf("[DEBUG] Image list filtered by regex: %s", d.Get("name_regex"))
	}

	log.Printf("[DEBUG] Got %d images after filtering in openstack_images_image_ids_v2: %+v", len(allImages), allImages)

	imageIDs := make([]string, len(allImages))
	for i, image := range allImages {
		imageIDs[i] = image.ID
	}

	d.SetId(strconv.Itoa(hashcode.String(strings.Join(imageIDs, ","))))
	d.Set("ids", imageIDs)
	d.Set("region", GetRegion(d, config))

	return nil
}
