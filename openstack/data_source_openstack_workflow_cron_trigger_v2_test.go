package openstack

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWorkflowV2CronTriggerDataSource_basic(t *testing.T) {
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
				Config: testAccWorkflowV2CronTriggerDataSourceBasic(workflowID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowV2CronTriggerDataSourceID("data.openstack_workflow_cron_trigger_v2.cron_trigger_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "name", "hello_cron_trigger"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "pattern", "0 5 * * *"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_id", workflowID),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_input.%", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_input.message", "Hello, OpenStack!"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_params.%", "2"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_params.priority", "high"),
					resource.TestCheckResourceAttr(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_params.notify", "mistral@openstack.org"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "project_id"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_workflow_cron_trigger_v2.cron_trigger_1", "created_at"),
				),
			},
		},
	})
}

func testAccCheckWorkflowV2CronTriggerDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cron trigger data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Cron trigger data source ID not set")
		}

		return nil
	}
}

func testAccWorkflowV2CronTriggerDataSourceBasic(workflowID string) string {
	return fmt.Sprintf(`
resource "openstack_workflow_cron_trigger_v2" "cron_trigger_1" {
  name        = "hello_cron_trigger"
  workflow_id = "%s"
  pattern     = "0 5 * * *"

  workflow_input = {
    message = "Hello, OpenStack!"
  }

  workflow_params = {
    priority = "high"
    notify   = "mistral@openstack.org"
  }
}

data "openstack_workflow_cron_trigger_v2" "cron_trigger_1" {
  name = openstack_workflow_cron_trigger_v2.cron_trigger_1.name
}
`, workflowID)
}
