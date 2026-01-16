package openstack

import (
	"context"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/crontriggers"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWorkflowCronTriggerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowCronTriggerV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pattern": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"workflow_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"workflow_input": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"workflow_params": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceWorkflowCronTriggerV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	workflowClient, err := config.WorkflowV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack workflow client: %s", err)
	}

	listOpts := crontriggers.ListOpts{}

	if v, ok := d.GetOk("project_id"); ok {
		listOpts.ProjectID = v.(string)
	}

	if v, ok := d.GetOk("workflow_id"); ok {
		listOpts.WorkflowID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = &crontriggers.ListFilter{
			Filter: crontriggers.FilterEQ,
			Value:  v.(string),
		}
	}

	allPages, err := crontriggers.List(workflowClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query cron triggers: %s", err)
	}

	allCrontriggers, err := crontriggers.ExtractCronTriggers(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve cron triggers: %s", err)
	}

	if len(allCrontriggers) < 1 {
		return diag.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(allCrontriggers) > 1 {
		tflog.Debug(ctx, "Multiple results found", map[string]any{
			"cronTriggers": allCrontriggers,
		})

		return diag.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	dataSourceWorkflowCronTriggerV2Attributes(ctx, d, &allCrontriggers[0], GetRegion(d, config))

	return nil
}

func dataSourceWorkflowCronTriggerV2Attributes(ctx context.Context, d *schema.ResourceData, crontrigger *crontriggers.CronTrigger, region string) {
	d.SetId(crontrigger.ID)
	d.Set("region", region)
	d.Set("name", crontrigger.Name)
	d.Set("project_id", crontrigger.ProjectID)
	d.Set("pattern", crontrigger.Pattern)
	d.Set("workflow_id", crontrigger.WorkflowID)
	d.Set("workflow_input", crontrigger.WorkflowInput)
	d.Set("workflow_params", crontrigger.WorkflowParams)

	if err := d.Set("created_at", crontrigger.CreatedAt.Format(time.RFC3339)); err != nil {
		tflog.Warn(ctx, "Unable to set created_at for openstack_workflow_cron_trigger_v2", map[string]any{
			"id":    crontrigger.ID,
			"error": err,
		})
	}
}
