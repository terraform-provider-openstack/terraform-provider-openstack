package openstack

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workflows"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkflowV2WorkflowDataSource_basic(t *testing.T) {
	var workflowID string

	if os.Getenv("TF_ACC") != "" {
		workflow, err := testAccWorkflowV2WorkflowCreate(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		workflowID = workflow.ID

		defer testAccWorkflowV2WorkflowDelete(t, workflow.ID)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWorkflow(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowV2WorkflowDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "id", workflowID),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "name", "hello_workflow"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "namespace", "my_namespace"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "input", "message"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_workflow_v2.workflow_1", "definition"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_workflow_v2.workflow_1", "scope", "private"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_workflow_v2.workflow_1", "project_id"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_workflow_v2.workflow_1", "created_at"),
				),
			},
		},
	})
}

func testAccWorkflowV2WorkflowCreate(ctx context.Context) (*workflows.Workflow, error) {
	config, err := testAccAuthFromEnv(ctx)
	if err != nil {
		return nil, err
	}

	client, err := config.WorkflowV2Client(ctx, osRegionName)
	if err != nil {
		return nil, err
	}

	createWorkflowOpts := workflows.CreateOpts{
		Scope:      "private",
		Namespace:  "my_namespace",
		Definition: strings.NewReader(testAccWorkflowV2WorkflowDataSourceBasicDefinition),
	}

	workflows, err := workflows.Create(ctx, client, createWorkflowOpts).Extract()
	if err != nil {
		return nil, err
	}

	workflow := workflows[len(workflows)-1]

	return &workflow, nil
}

func testAccWorkflowV2WorkflowDelete(t *testing.T, workflowID string) {
	config, err := testAccAuthFromEnv(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	client, err := config.WorkflowV2Client(t.Context(), osRegionName)
	if err != nil {
		t.Fatal(err)
	}

	err = workflows.Delete(t.Context(), client, workflowID).ExtractErr()
	if err != nil {
		t.Fatal(err)
	}
}

const testAccWorkflowV2WorkflowDataSourceBasic = `
data "openstack_workflow_workflow_v2" "workflow_1" {
	name      = "hello_workflow"
	namespace = "my_namespace"
}
`

const testAccWorkflowV2WorkflowDataSourceBasicDefinition = `
version: '2.0'

hello_workflow:
  description: Simple echo example

  input:
    - message

  tags:
    - echo

  tasks:
    echo:
      action: std.echo
      input:
        output:
          my_message: <% $.message %>
`
