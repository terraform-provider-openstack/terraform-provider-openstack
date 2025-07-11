package openstack

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-provider-openstack/utils/v2/hashcode"
)

func dataSourceIdentityProjectIDsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityProjectIDsV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				ConflictsWith: []string{"name_regex"},
			},

			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
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

// dataSourceIdentityProjectIDsV3Read performs the project lookup.
func dataSourceIdentityProjectIDsV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	isDomain := d.Get("is_domain").(bool)

	listOpts := projects.ListOpts{
		DomainID: d.Get("domain_id").(string),
		Enabled:  &enabled,
		IsDomain: &isDomain,
		Name:     d.Get("name").(string),
		ParentID: d.Get("parent_id").(string),
		Tags:     strings.Join(expandObjectTags(d), ","),
	}

	allPages, err := projects.List(identityClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to list projects in openstack_identity_project_ids_v3: %s", err)
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve projects in openstack_identity_project_ids_v3: %s", err)
	}

	log.Printf("[DEBUG] Retrieved %d projects in openstack_identity_project_ids_v3: %+v", len(allProjects), allProjects)

	if v, ok := d.GetOk("name_regex"); ok {
		allProjects, err = projectsFilterByRegex(allProjects, v.(string))
		if err != nil {
			return diag.Errorf("Error while compiling regex: %s", err)
		}

		log.Printf("[DEBUG] Project list filtered by regex: %s", v)
	}

	log.Printf("[DEBUG] Got %d projects after filtering in openstack_identity_project_ids_v3: %+v", len(allProjects), allProjects)

	projectIDs := make([]string, len(allProjects))
	for i, project := range allProjects {
		projectIDs[i] = project.ID
	}

	d.SetId(strconv.Itoa(hashcode.String(strings.Join(projectIDs, ","))))
	d.Set("ids", projectIDs)
	d.Set("region", GetRegion(d, config))

	return nil
}

func projectsFilterByRegex(projectArr []projects.Project, nameRegex string) ([]projects.Project, error) {
	r, err := regexp.Compile(nameRegex)
	if err != nil {
		return nil, err
	}

	result := make([]projects.Project, 0, len(projectArr))

	for _, project := range projectArr {
		if r.MatchString(project.Name) {
			result = append(result, project)
		}
	}

	return result, nil
}
