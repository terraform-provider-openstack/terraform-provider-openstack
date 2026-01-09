package openstack

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/crontriggers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWorkflowV2CronTrigger_basic(t *testing.T) {
	var workflowID string

	if os.Getenv("TF_ACC") != "" {
		workflow, err := testAccWorkflowV2WorkflowCreate(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		workflowID = workflow.ID
		defer testAccWorkflowV2WorkflowDelete(t, workflowID)
	}

	var crontrigger crontriggers.CronTrigger

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckWorkflow(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckWorkflowV2CronTriggerDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowV2CronTriggerBasic(workflowID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowV2CronTriggerExists(t.Context(),
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", &crontrigger),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "id"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "name", "hello_cron_trigger"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_id", workflowID),
					resource.TestCheckResourceAttr(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "pattern", "0 5 * * *"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_input.%", "1"),
					resource.TestCheckResourceAttr(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "workflow_params.%", "2"),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "project_id"),
					resource.TestCheckResourceAttrSet(
						"openstack_workflow_cron_trigger_v2.cron_trigger_1", "created_at"),
				),
			},
		},
	})
}

func testAccCheckWorkflowV2CronTriggerDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		workflowClient, err := config.WorkflowV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack workflow client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_workflow_cron_trigger_v2" {
				continue
			}

			_, err := crontriggers.Get(ctx, workflowClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("CronTrigger still exists")
			}
		}

		return nil
	}
}

func testAccCheckWorkflowV2CronTriggerExists(ctx context.Context, n string, crontrigger *crontriggers.CronTrigger) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		workflowClient, err := config.WorkflowV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack workflow client: %w", err)
		}

		found, err := crontriggers.Get(ctx, workflowClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("CronTrigger not found")
		}

		*crontrigger = *found

		return nil
	}
}

func testAccWorkflowV2CronTriggerBasic(workflowID string) string {
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
`, workflowID)
}
