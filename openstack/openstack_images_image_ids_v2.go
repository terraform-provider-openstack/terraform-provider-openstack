package openstack

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceImagesImageIdsV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImagesImageIdsV2Read,

		Schema: map[string]*schema.Schema{

			"region": &schema.Schema{
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

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},

			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "name",
			},

			"sort_direction": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "asc",
				ValidateFunc: dataSourceImagesImageV2SortDirection,
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

			// Computed values
			"ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceImagesImageIdsV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	imageClient, err := config.imageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack image client: %s", err)
	}

	_, nameOk := d.GetOk("name")
	_, nameRegexOk := d.GetOk("name_regex")

	if nameOk == true && nameRegexOk == true {
		return fmt.Errorf("Attributes name and name_regexp can not be used at the same time")
	}

	visibility := resourceImagesImageV2VisibilityFromString(d.Get("visibility").(string))

	listOpts := images.ListOpts{
		Name:       d.Get("name").(string),
		Visibility: visibility,
		Owner:      d.Get("owner").(string),
		Status:     images.ImageStatusActive,
		SizeMin:    int64(d.Get("size_min").(int)),
		SizeMax:    int64(d.Get("size_max").(int)),
		SortKey:    d.Get("sort_key").(string),
		SortDir:    d.Get("sort_direction").(string),
		Tag:        d.Get("tag").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	allPages, err := images.List(imageClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query images: %s", err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve images: %s", err)
	}

	properties := d.Get("properties").(map[string]interface{})
	imageProperties := resourceImagesImageV2ExpandProperties(properties)
	if len(allImages) > 1 && len(imageProperties) > 0 {
		var filteredImages []images.Image
		for _, image := range allImages {
			if len(image.Properties) > 0 {
				match := true
				for searchKey, searchValue := range imageProperties {
					imageValue, ok := image.Properties[searchKey]
					if !ok {
						match = false
						break
					}

					if searchValue != imageValue {
						match = false
						break
					}
				}

				if match {
					filteredImages = append(filteredImages, image)
				}
			}
		}
		allImages = filteredImages
	}

	if nameRegexOk {
		var filteredImages []images.Image
		r := regexp.MustCompile(d.Get("name_regex").(string))
		for _, image := range allImages {
			// Check for a very rare case where the response would include no
			// image name. No name means nothing to attempt a match against,
			// therefore we are skipping such image.
			if image.Name == "" {
				log.Printf("[WARN] Unable to find AMI name to match against "+
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

	searchParams := dataSourceImagesImageIdsV2SearchParams(d)

	d.SetId(fmt.Sprintf("%d", hashcode.String(strings.Join(searchParams, ","))))
	d.Set("ids", imageIds)

	return nil
}

func dataSourceImagesImageIdsV2SearchParams(d *schema.ResourceData) []string {
	properties := resourceImagesImageV2ExpandProperties(d.Get("properties").(map[string]interface{}))
	propertyKeys := make([]string, 0)
	propertyValues := make([]string, 0)
	for key, value := range properties {
		propertyKeys = append(propertyKeys, key)
		propertyValues = append(propertyValues, value)
	}

	searchParams := []string{
		strParam(d.Get("region")),
		strParam(d.Get("name")),
		strParam(d.Get("nameRegex")),
		strParam(d.Get("visibility")),
		strParam(d.Get("owner")),
		fmt.Sprintf("%v", int64(d.Get("size_min").(int))),
		fmt.Sprintf("%v", int64(d.Get("size_max").(int))),
		strParam(d.Get("sort_key")),
		strParam(d.Get("sort_direction")),
		strParam(d.Get("tag")),
		strings.Join(propertyKeys, ","),
		strings.Join(propertyValues, ","),
	}

	return searchParams
}

func strParam(param interface{}) string {
	result := ""
	if param != nil {
		result = param.(string)
	}

	return result
}

func validateNameRegex(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if _, err := regexp.Compile(value); err != nil {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid regular expression: %s",
			k, err))
	}
	return
}
