package openstack

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/utils/terraform/hashcode"
)

func dataSourceIdentityProjectIdsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityProjectIdsV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"is_domain": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_regex"},
			},

			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				ConflictsWith: []string{"name"},
			},

			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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

// dataSourceIdentityProjectIdsV3Read performs the project lookup.
func dataSourceIdentityProjectIdsV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	isDomain := d.Get("is_domain").(bool)

	tags := []string{}
	tagList := d.Get("tags").(*schema.Set).List()
	for _, v := range tagList {
		tags = append(tags, fmt.Sprint(v))
	}
	joinedTags := strings.Join(tags, ",")

	listOpts := projects.ListOpts{
		DomainID: d.Get("domain_id").(string),
		Enabled:  &enabled,
		IsDomain: &isDomain,
		Name:     d.Get("name").(string),
		ParentID: d.Get("parent_id").(string),
		Tags:     joinedTags,
	}

	allPages, err := projects.List(identityClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to list projects in openstack_identity_project_ids_v3: %s", err)
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve projects in openstack_identity_project_ids_v3: %s", err)
	}

	log.Printf("[DEBUG] Retrieved %d images in openstack_identity_project_ids_v3: %+v", len(allProjects), allProjects)

	nameRegex, nameRegexOk := d.GetOk("name_regex")
	if nameRegexOk {
		allProjects = projectsFilterByRegex(allProjects, nameRegex.(string))
		log.Printf("[DEBUG] Project list filtered by regex: %s", d.Get("name_regex"))
	}

	log.Printf("[DEBUG] Got %d projects after filtering in openstack_identity_project_ids_v3: %+v", len(allProjects), allProjects)

	projectIDs := make([]string, len(allProjects))
	for i, image := range allProjects {
		projectIDs[i] = image.ID
	}

	d.SetId(fmt.Sprintf("%d", hashcode.String(strings.Join(projectIDs, ","))))
	d.Set("ids", projectIDs)

	return nil
}

func projectsFilterByRegex(projectArr []projects.Project, nameRegex string) []projects.Project {
	var result []projects.Project
	r := regexp.MustCompile(nameRegex)

	for _, project := range projectArr {
		if r.MatchString(project.Name) {
			result = append(result, project)
		}
	}

	return result
}
