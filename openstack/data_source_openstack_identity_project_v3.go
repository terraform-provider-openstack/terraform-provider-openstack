package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
)

func dataSourceIdentityProjectV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityProjectV3Read,

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
				Type:     schema.TypeString,
				Optional: true,
			},

			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"name", "domain_id", "description", "enabled", "is_domain", "parent_id", "tags",
				},
			},
		},
	}
}

func filterProjects(allProjects []projects.Project, listOpts projects.ListOpts) (results []projects.Project) {
	for _, p := range allProjects {
		// TODO: Test for all fields :)
		if p.Name == listOpts.Name {
			results = append(results, p)
		}
	}
	return
}

// dataSourceIdentityProjectV3Read performs the project lookup.
func dataSourceIdentityProjectV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var allProjects []projects.Project
	if v := d.Get("project_id").(string); v != "" {
		project, err := projects.Get(ctx, identityClient, v).Extract()
		if err != nil {
			return diag.Errorf("Unable to query openstack_identity_project_v3: %s", err)
		}
		allProjects = append(allProjects, *project)
	} else {
		enabled := d.Get("enabled").(bool)
		isDomain := d.Get("is_domain").(bool)
		listOpts := projects.ListOpts{
			DomainID: d.Get("domain_id").(string),
			Enabled:  &enabled,
			IsDomain: &isDomain,
			Name:     d.Get("name").(string),
			ParentID: d.Get("parent_id").(string),
		}

		allPages, err := projects.List(identityClient, listOpts).AllPages(ctx)
		if err != nil {
			userID := config.UserID
			log.Printf("[DEBUG] Will try to find project with users.ListProjects as I am unable to query openstack_identity_project_v3: %s. Trying listing userprojects.", err)
			if userID == "" {
				tokenInfo, tokenErr := getTokenInfo(ctx, identityClient)
				if tokenErr != nil {
					return diag.Errorf("Error when getting token info: %s", err)
				}
				userID = tokenInfo.userID
			}
			// Search for all the projects using the users.ListProjects API call and filter them
			allPages, err = users.ListProjects(identityClient, userID).AllPages(ctx)
			if err != nil {
				return diag.Errorf("Unable to query openstack_identity_project_v3: %s", err)
			}
			allProjects, err = projects.ExtractProjects(allPages)
			if err != nil {
				return diag.Errorf("Unable to retrieve openstack_identity_project_v3: %s", err)
			}
			allProjects = filterProjects(allProjects, listOpts)
		} else {
			allProjects, err = projects.ExtractProjects(allPages)
			if err != nil {
				return diag.Errorf("Unable to retrieve openstack_identity_project_v3: %s", err)
			}
		}
	}

	if len(allProjects) < 1 {
		return diag.Errorf("Your openstack_identity_project_v3 query returned no results. " +
			"Please change your search criteria and try again")
	}

	if len(allProjects) > 1 {
		return diag.Errorf("Your openstack_identity_project_v3 query returned more than one result %#v", allProjects)
	}

	dataSourceIdentityProjectV3Attributes(d, &allProjects[0])
	d.Set("region", GetRegion(d, config))

	return nil
}

// dataSourceIdentityProjectV3Attributes populates the fields of an Project resource.
func dataSourceIdentityProjectV3Attributes(d *schema.ResourceData, project *projects.Project) {
	log.Printf("[DEBUG] openstack_identity_project_v3 details: %#v", project)

	d.SetId(project.ID)
	d.Set("project_id", project.ID)
	d.Set("is_domain", project.IsDomain)
	d.Set("description", project.Description)
	d.Set("domain_id", project.DomainID)
	d.Set("enabled", project.Enabled)
	d.Set("name", project.Name)
	d.Set("parent_id", project.ParentID)
	d.Set("tags", project.Tags)
}
