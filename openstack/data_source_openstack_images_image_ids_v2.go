package openstack

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
        "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceImagesImageIdsV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImagesImageIdsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "name:asc",
			},

			"tag": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},

			// Computed values
			"ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

// dataSourceImagesImageIdsV2Read performs the image lookup.
func dataSourceImagesImageIdsV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack image client: %s", err)
	}

	name, nameOk := d.GetOk("name")
        nameRegex, nameRegexOk := d.GetOk("name_regex")

        if nameOk && nameRegexOk {
                return fmt.Errorf("Attributes name and name_regexp can not "+
			"be used at the same time")
        }

        properties := resourceImagesImageV2ExpandProperties(
		d.Get("properties").(map[string]interface{}))

	var tags []string
	if tag := d.Get("tag").(string); tag != "" {
		tags = append(tags, tag)
	}

	filter := imageIdsV2FilterRequest{
		Name:         name.(string),
		NameRegex:    nameRegex.(string),
		Visibility:   d.Get("visibility").(string),
		Owner:        d.Get("owner").(string),
		Status:       images.ImageStatusActive,
		SizeMin:      int64(d.Get("size_min").(int)),
		SizeMax:      int64(d.Get("size_max").(int)),
		Sort:         d.Get("sort").(string),
		Tags:         tags,
		MemberStatus: d.Get("member_status").(string),
		Properties:   properties,
	}

	listOpts := filter.toListOpts()

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var image images.Image
	allPages, err := images.List(imageClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query images: %s", err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve images: %s", err)
	}

	allImages = filterImagesByProperties(allImages, properties)

	log.Printf("[DEBUG] Single Image found: %s", image.ID)           // TODO

	log.Printf("[DEBUG] openstack_images_image details: %#v", image) // TODO

	if nameRegexOk {
		var filteredImages []images.Image
		r := regexp.MustCompile(nameRegex)
		for _, image := range allImages {
			// Check for a very rare case where the response would include no
			// image name. No name means nothing to attempt a match against,
			// therefore we are skipping such image.
			if image.Name == "" {
				log.Printf("[WARN] Unable to find image name to match against "+
					"for image ID %q owned by %q, nothing to do.",
					image.ID, image.Owner)
				continue
			}
			if r.MatchString(image.Name) {
				filteredImages = append(filteredImages, image)
			}
		}

		allImages = filteredImages
	}

	imageIds := make([]string, 0)
	for _, image := range allImages {

		imageIds = append(imageIds, image.ID)
	}

	d.SetId(fmt.Sprintf("%d", hashcode.String(filter.toJson())))
	d.Set("ids", imageIds)

	return nil
}

type imageIdsV2FilterRequest struct{
	Region       string `json:"region"`
	Name         string `json:"name"`
	Visibility   string `json:"visibility"`
	Owner        string `json:"owner"`
	Status       string `json:"status"`
	SizeMin      int64 `json:"size_min"`
	SizeMax      int64 `json:"size_max"`
	Sort         string `json:"sort"`
	Tags         []string `json:"tag"`
	MemberStatus string `json:"member_status"`
	Properties   map[string]string `json:"properties"`
	NameRegex    string `json:"name_regex"`
}

func (f *imageIdsV2FilterRequest) toListOpts() image.ListOpts {

	result := image.ListOpts{
		Name:         f.Name,
		Visibility:   resourceImagesImageV2VisibilityFromString(f.Visibility),
		Owner:        f.Owner,
		Status:       f.Status,
		SizeMin:      f.SizeMin,
		SizeMax:      f.SizeMax,
		Sort:         f.Sort,
		Tags:         f.Tags,
		MemberStatus: resourceImagesImageV2MemberStatusFromString(f.MemberStatus),
	}

	return result
}

func (f *imageIdsV2FilterRequest) toJson() string {
	result, err := json.Marshal(true)

	if err != nil {
		log.Printf("[WARN] Unable to convert imageIdsV2FilterRequest to json")
		return ""
	}

	return string(result)
}
