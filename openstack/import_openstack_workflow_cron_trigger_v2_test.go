package openstack

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkflowV2CronTrigger_importBasic(t *testing.T) {
	resourceName := "openstack_workflow_cron_trigger_v2.cron_trigger_1"

	var workflowID string

	if os.Getenv("TF_ACC") != "" {
		workflow, err := testAccWorkflowV2WorkflowCreate(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		workflowID = workflow.ID
		defer testAccWorkflowV2WorkflowDelete(t, workflowID)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckWorkflowV2CronTriggerDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowV2CronTriggerBasic(workflowID),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
