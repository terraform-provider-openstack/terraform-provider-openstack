package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacks"
)

func TestAccOrchestrationV1Stack_basic(t *testing.T) {
	var stack stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckOrchestrationV1StackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationV1StackBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationV1StackExists("openstack_orchestration_stack_v1.stack_1", &stack),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_1", "name", "stack_1"),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_1", "parameters.length", "4"),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_1", "timeout", "30"),
				),
			},
		},
	})
}

func TestAccOrchestrationV1Stack_tags(t *testing.T) {
	var stack stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckOrchestrationV1StackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationV1StackTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationV1StackExists("openstack_orchestration_stack_v1.stack_4", &stack),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_4", "name", "stack_4"),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_4", "tags.#", "2"),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_4", "tags.0", "foo"),
				),
			},
		},
	})
}

func TestAccOrchestrationV1Stack_update(t *testing.T) {
	var stack stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckOrchestrationV1StackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationV1StackPreUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationV1StackExists("openstack_orchestration_stack_v1.stack_3", &stack),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_3", "name", "stack_3"),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_3", "parameters.length", "4"),
				),
			},
			{
				Config: testAccOrchestrationV1StackUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationV1StackExists("openstack_orchestration_stack_v1.stack_3", &stack),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_3", "name", "stack_3"),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_3", "parameters.length", "5"),
					resource.TestCheckResourceAttrSet("openstack_orchestration_stack_v1.stack_3", "updated_time"),
				),
			},
		},
	})
}

func TestAccOrchestrationV1Stack_timeout(t *testing.T) {
	var stack stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckOrchestrationV1StackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationV1StackTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationV1StackExists("openstack_orchestration_stack_v1.stack_2", &stack),
				),
			},
		},
	})
}

func TestAccOrchestrationV1Stack_outputs(t *testing.T) {
	var stack stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckOrchestrationV1StackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationV1StackOutputs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrchestrationV1StackExists("openstack_orchestration_stack_v1.stack_5", &stack),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_5", "name", "stack_5"),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_5", "outputs.#", "1"),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_5", "outputs.0.output_value", "foo"),
					resource.TestCheckResourceAttr("openstack_orchestration_stack_v1.stack_5", "outputs.0.output_key", "value1"),
				),
			},
		},
	})
}

func testAccCheckOrchestrationV1StackDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	orchestrationClient, err := config.OrchestrationV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Orchestration client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_orchestration_stack_v1" {
			continue
		}

		stack, err := stacks.Find(orchestrationClient, rs.Primary.ID).Extract()
		if err == nil {
			if stack.Status != "DELETE_COMPLETE" {
				return fmt.Errorf("stack still exists")
			}
		}
	}

	return nil
}

func testAccCheckOrchestrationV1StackExists(n string, stack *stacks.RetrievedStack) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		orchestrationClient, err := config.OrchestrationV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack orchestration client: %s", err)
		}

		found, err := stacks.Find(orchestrationClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Stack not found")
		}

		*stack = *found

		return nil
	}
}

const testAccOrchestrationV1StackBasic = `
resource "openstack_orchestration_stack_v1" "stack_1" {
  name = "stack_1"
  parameters = {
	length = 4
  }
  template_opts = {
	Bin = "heat_template_version: 2013-05-23\nparameters:\n  length:\n    type: number\nresources:\n  test_res:\n    type: OS::Heat::TestResource\n  random:\n    type: OS::Heat::RandomString\n    properties:\n      length: {get_param: length}\n"
  }
  environment_opts = {
	Bin = "\n"
  }
  disable_rollback = true
  timeout = 30
}
`

const testAccOrchestrationV1StackPreUpdate = `
resource "openstack_orchestration_stack_v1" "stack_3" {
  name = "stack_3"
  parameters = {
	length = 4
  }
  template_opts = {
	Bin = "heat_template_version: 2013-05-23\nparameters:\n  length:\n    type: number\nresources:\n  test_res:\n    type: OS::Heat::TestResource\n  random:\n    type: OS::Heat::RandomString\n    properties:\n      length: {get_param: length}\n"
  }
  environment_opts = {
	Bin = "\n"
  }
  disable_rollback = true
}
`

const testAccOrchestrationV1StackUpdate = `
resource "openstack_orchestration_stack_v1" "stack_3" {
  name = "stack_3"
  parameters = {
	length = 5
  }
  template_opts = {
	Bin = "heat_template_version: 2013-05-23\nparameters:\n  length:\n    type: number\nresources:\n  test_res:\n    type: OS::Heat::TestResource\n  random:\n    type: OS::Heat::RandomString\n    properties:\n      length: {get_param: length}\n"
  }
  environment_opts = {
	Bin = "\n"
  }
  disable_rollback = true
}
`

const testAccOrchestrationV1StackTimeout = `
resource "openstack_orchestration_stack_v1" "stack_2" {
  name = "stack_2"
  parameters = {
	length = 4
  }
  template_opts = {
	Bin = "heat_template_version: 2013-05-23\nparameters:\n  length:\n    type: number\nresources:\n  test_res:\n    type: OS::Heat::TestResource\n  random:\n    type: OS::Heat::RandomString\n    properties:\n      length: {get_param: length}\n"
  }
  environment_opts = {
	Bin = "\n"
  }
  disable_rollback = true
  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccOrchestrationV1StackTags = `
resource "openstack_orchestration_stack_v1" "stack_4" {
  name = "stack_4"
  parameters = {
	length = 4
  }
  template_opts = {
	Bin = "heat_template_version: 2013-05-23\nparameters:\n  length:\n    type: number\nresources:\n  test_res:\n    type: OS::Heat::TestResource\n  random:\n    type: OS::Heat::RandomString\n    properties:\n      length: {get_param: length}\n"
  }
  environment_opts = {
	Bin = "\n"
  }
  disable_rollback = true
  tags = [
    "foo",
    "bar",
  ]
}
`

const testAccOrchestrationV1StackOutputs = `
resource "openstack_orchestration_stack_v1" "stack_5" {
  name = "stack_5"
  parameters = {
	length = 4
  }
  template_opts = {
	Bin = "heat_template_version: 2013-05-23\nparameters:\n  length:\n    type: number\nresources:\n  test_res:\n    type: OS::Heat::TestResource\n  random:\n    type: OS::Heat::RandomString\n    properties:\n      length: {get_param: length}\noutputs:\n  value1:\n    value: foo"
  }
  environment_opts = {
	Bin = "\n"
  }
  disable_rollback = true
}
`
