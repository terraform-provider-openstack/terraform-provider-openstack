package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
)

func TestAccComputeV2SecGroup_basic(t *testing.T) {
	var secgroup secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroupBasicOrig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_1", &secgroup),
				),
			},
		},
	})
}

func TestAccComputeV2SecGroup_update(t *testing.T) {
	var secgroup secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroupBasicOrig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_1", &secgroup),
				),
			},
			{
				Config: testAccComputeV2SecGroupBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_1", &secgroup),
					testAccCheckComputeV2SecGroupRuleCount(&secgroup, 2),
				),
			},
		},
	})
}

func TestAccComputeV2SecGroup_groupID(t *testing.T) {
	var secgroup1, secgroup2, secgroup3 secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroupIDOrig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_1", &secgroup1),
					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_2", &secgroup2),
					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_3", &secgroup3),
					testAccCheckComputeV2SecGroupGroupIDMatch(&secgroup1, &secgroup3),
				),
			},
			{
				Config: testAccComputeV2SecGroupIDUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_1", &secgroup1),
					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_2", &secgroup2),
					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_3", &secgroup3),
					testAccCheckComputeV2SecGroupGroupIDMatch(&secgroup2, &secgroup3),
				),
			},
		},
	})
}

// NOTE: disabled since Terraform SDK V2 currently does not support indexes into TypeSet.
//func TestAccComputeV2SecGroup_self(t *testing.T) {
//	var secgroup secgroups.SecurityGroup
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:          func() { testAccPreCheck(t) },
//		ProviderFactories: testAccProviders,
//		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccComputeV2SecGroupSelf,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_1", &secgroup),
//					testAccCheckComputeV2SecGroupGroupIDMatch(&secgroup, &secgroup),
//					resource.TestCheckResourceAttr(
//						"openstack_compute_secgroup_v2.sg_1", "rule.3170486100.self", "true"),
//					resource.TestCheckResourceAttr(
//						"openstack_compute_secgroup_v2.sg_1", "rule.3170486100.from_group_id", ""),
//				),
//			},
//		},
//	})
//}

func TestAccComputeV2SecGroup_icmpZero(t *testing.T) {
	var secgroup secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroupIcmpZero,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_1", &secgroup),
				),
			},
		},
	})
}

// NOTE: disabled since Terraform SDK V2 currently does not support indexes into TypeSet.
//func TestAccComputeV2SecGroup_lowerCaseCIDR(t *testing.T) {
//	var secgroup secgroups.SecurityGroup
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:          func() { testAccPreCheck(t) },
//		ProviderFactories: testAccProviders,
//		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccComputeV2SecGroupLowerCaseCIDR,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.sg_1", &secgroup),
//					resource.TestCheckResourceAttr(
//						"openstack_compute_secgroup_v2.sg_1", "rule.768649014.cidr", "2001:558:fc00::/39"),
//				),
//			},
//		},
//	})
//}

func testAccCheckComputeV2SecGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_compute_secgroup_v2" {
			continue
		}

		_, err := secgroups.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Security group still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2SecGroupExists(n string, secgroup *secgroups.SecurityGroup) resource.TestCheckFunc {
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

		found, err := secgroups.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group not found")
		}

		*secgroup = *found

		return nil
	}
}

func testAccCheckComputeV2SecGroupRuleCount(secgroup *secgroups.SecurityGroup, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(secgroup.Rules) != count {
			return fmt.Errorf("Security group rule count does not match. Expected %d, got %d", count, len(secgroup.Rules))
		}

		return nil
	}
}

func testAccCheckComputeV2SecGroupGroupIDMatch(sg1, sg2 *secgroups.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg2.Rules) == 1 {
			if sg1.Name != sg2.Rules[0].Group.Name || sg1.TenantID != sg2.Rules[0].Group.TenantID {
				return fmt.Errorf("%s was not correctly applied to %s", sg1.Name, sg2.Name)
			}
		} else {
			return fmt.Errorf("%s rule count is incorrect", sg2.Name)
		}

		return nil
	}
}

const testAccComputeV2SecGroupBasicOrig = `
resource "openstack_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
  rule {
    from_port = 1
    to_port = 65535
    ip_protocol = "udp"
    cidr = "0.0.0.0/0"
  }
  rule {
    from_port = -1
    to_port = -1
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }
}
`

const testAccComputeV2SecGroupBasicUpdate = `
resource "openstack_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 2200
    to_port = 2200
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
  rule {
    from_port = -1
    to_port = -1
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }
}
`

const testAccComputeV2SecGroupIDOrig = `
resource "openstack_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "sg_2" {
  name = "sg_2"
  description = "second test security group"
  rule {
    from_port = -1
    to_port = -1
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "sg_3" {
  name = "sg_3"
  description = "third test security group"
  rule {
    from_port = 80
    to_port = 80
    ip_protocol = "tcp"
    from_group_id = "${openstack_compute_secgroup_v2.sg_1.id}"
  }
}
`

const testAccComputeV2SecGroupIDUpdate = `
resource "openstack_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "sg_2" {
  name = "sg_2"
  description = "second test security group"
  rule {
    from_port = -1
    to_port = -1
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "sg_3" {
  name = "sg_3"
  description = "third test security group"
  rule {
    from_port = 80
    to_port = 80
    ip_protocol = "tcp"
    from_group_id = "${openstack_compute_secgroup_v2.sg_2.id}"
  }
}
`

//const testAccComputeV2SecGroupSelf = `
//resource "openstack_compute_secgroup_v2" "sg_1" {
//  name = "sg_1"
//  description = "first test security group"
//  rule {
//    from_port = 22
//    to_port = 22
//    ip_protocol = "tcp"
//    self = true
//  }
//}
//`

const testAccComputeV2SecGroupIcmpZero = `
resource "openstack_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 0
    to_port = 0
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }
}
`

//const testAccComputeV2SecGroupLowerCaseCIDR = `
//resource "openstack_compute_secgroup_v2" "sg_1" {
//  name = "sg_1"
//  description = "first test security group"
//  rule {
//    from_port = 22
//    to_port = 22
//    ip_protocol = "tcp"
//    cidr = "2001:558:FC00::/39"
//  }
//}
//`
