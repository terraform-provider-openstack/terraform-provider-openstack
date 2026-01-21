package openstack

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workflows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWorkflowWorkflowV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowWorkflowV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"input": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"definition": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceWorkflowWorkflowV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	listOpts := workflows.ListOpts{
		ProjectID: d.Get("project_id").(string),
	}

	name := d.Get("name").(string)
	if name != "" {
		listOpts.Name = &workflows.ListFilter{
			Filter: workflows.FilterEQ,
			Value:  name,
		}
	}

	namespace := d.Get("namespace").(string)
	if namespace != "" {
		listOpts.Namespace = &workflows.ListFilter{
			Filter: workflows.FilterEQ,
			Value:  namespace,
		}
	}

	allPages, err := workflows.List(workflowClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query workflows: %s", err)
	}

	allWorkflows, err := workflows.ExtractWorkflows(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve workflows: %s", err)
	}

	if len(allWorkflows) < 1 {
		return diag.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	var workflow workflows.Workflow

	if len(allWorkflows) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allWorkflows)

		return diag.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	workflow = allWorkflows[0]

	dataSourceWorkflowWorkflowV2Attributes(d, &workflow, GetRegion(d, config))

	return nil
}

func dataSourceWorkflowWorkflowV2Attributes(d *schema.ResourceData, workflow *workflows.Workflow, region string) {
	d.SetId(workflow.ID)
	d.Set("region", region)
	d.Set("name", workflow.Name)
	d.Set("namespace", workflow.Namespace)
	d.Set("input", workflow.Input)
	d.Set("definition", workflow.Definition)
	d.Set("tags", workflow.Tags)
	d.Set("scope", workflow.Scope)
	d.Set("project_id", workflow.ProjectID)

	if err := d.Set("created_at", workflow.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set created_at for openstack_workflow_workflow_v2 %s: %s", workflow.ID, err)
	}
}
