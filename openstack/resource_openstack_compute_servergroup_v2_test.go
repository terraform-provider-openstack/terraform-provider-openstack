package openstack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func TestAccComputeV2ServerGroup_basic(t *testing.T) {
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("openstack_compute_servergroup_v2.sg_1", &sg),
					resource.TestCheckResourceAttr(
						"openstack_compute_servergroup_v2.sg_1", "policies.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_servergroup_v2.sg_1", "policies.0", "affinity"),
				),
			},
		},
	})
}

func TestAccComputeV2ServerGroup_basic_v2_64(t *testing.T) {
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroupV264Policy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("openstack_compute_servergroup_v2.sg_1", &sg),
					resource.TestCheckResourceAttr(
						"openstack_compute_servergroup_v2.sg_1", "policies.0", "affinity"),
				),
			},
		},
	})
}

func TestAccComputeV2ServerGroup_v2_64_anti_affinity(t *testing.T) {
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroupV264PolicyAntiAffinity,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("openstack_compute_servergroup_v2.sg_1", &sg),
					resource.TestCheckResourceAttr(
						"openstack_compute_servergroup_v2.sg_1", "policies.0", "anti-affinity"),
				),
			},
		},
	})
}

func TestAccComputeV2ServerGroup_v2_64_with_rules(t *testing.T) {
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroupV264PolicyRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("openstack_compute_servergroup_v2.sg_1", &sg),
					resource.TestCheckResourceAttr(
						"openstack_compute_servergroup_v2.sg_1", "policies.0", "anti-affinity"),
					resource.TestCheckResourceAttr(
						"openstack_compute_servergroup_v2.sg_1", "rules.0.max_server_per_host", "2"),
				),
			},
		},
	})
}

func TestAccComputeV2ServerGroup_v2_64_with_invalid_rules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeV2ServerGroupV264InvalidPolicyRules,
				ExpectError: regexp.MustCompile(`expected rules\.0\.max_server_per_host to be at least \(1\), got .*`),
			},
		},
	})
}

func TestAccComputeV2ServerGroup_affinity(t *testing.T) {
	var instance servers.Server
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroupAffinity(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("openstack_compute_servergroup_v2.sg_1", &sg),
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(
						"openstack_compute_servergroup_v2.sg_1", "policies.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_servergroup_v2.sg_1", "policies.0", "affinity"),
				),
			},
		},
	})
}

func TestAccComputeV2ServerGroup_affinity_v2_64(t *testing.T) {
	var instance servers.Server
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroupAffinityV264(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("openstack_compute_servergroup_v2.sg_1", &sg),
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(
						"openstack_compute_servergroup_v2.sg_1", "policies.0", "affinity"),
				),
			},
		},
	})
}

func TestAccComputeV2ServerGroup_soft_affinity(t *testing.T) {
	var instance servers.Server
	var sg servergroups.ServerGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroupSoftAffinity(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2ServerGroupExists("openstack_compute_servergroup_v2.sg_1", &sg),
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(
						"openstack_compute_servergroup_v2.sg_1", "policies.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_servergroup_v2.sg_1", "policies.0", "soft-affinity"),
				),
			},
		},
	})
}

func testAccCheckComputeV2ServerGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_compute_servergroup_v2" {
			continue
		}

		_, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ServerGroup still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2ServerGroupExists(n string, kp *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		found, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ServerGroup not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckComputeV2InstanceInServerGroup(instance *servers.Server, sg *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Members) > 0 {
			for _, m := range sg.Members {
				if m == instance.ID {
					return nil
				}
			}
		}

		return fmt.Errorf("Instance %s is not part of Server Group %s", instance.ID, sg.ID)
	}
}

const testAccComputeV2ServerGroupBasic = `
resource "openstack_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["affinity"]
}
`

const testAccComputeV2ServerGroupV264Policy = `
resource "openstack_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["affinity"]
}
`

const testAccComputeV2ServerGroupV264PolicyAntiAffinity = `
resource "openstack_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["anti-affinity"]
}
`

const testAccComputeV2ServerGroupV264PolicyRules = `
resource "openstack_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["anti-affinity"]
  rules {
    max_server_per_host = 2
  }
}
`

const testAccComputeV2ServerGroupV264InvalidPolicyRules = `
resource "openstack_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["anti-affinity"]
  rules {
    max_server_per_host = -1
  }
}
`

func testAccComputeV2ServerGroupAffinity() string {
	return fmt.Sprintf(`
resource "openstack_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["affinity"]
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  scheduler_hints {
    group = "${openstack_compute_servergroup_v2.sg_1.id}"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2ServerGroupSoftAffinity() string {
	return fmt.Sprintf(`
resource "openstack_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["soft-affinity"]
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  scheduler_hints {
    group = "${openstack_compute_servergroup_v2.sg_1.id}"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2ServerGroupAffinityV264() string {
	return fmt.Sprintf(`
resource "openstack_compute_servergroup_v2" "sg_1" {
  name = "sg_1"
  policies = ["affinity"]
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  scheduler_hints {
    group = "${openstack_compute_servergroup_v2.sg_1.id}"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}
