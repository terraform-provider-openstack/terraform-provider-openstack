package openstack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func TestAccLBV2LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer

	lbProvider := "haproxy"
	if osUseOctavia != "" {
		lbProvider = "octavia"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2LoadBalancerConfigBasic(lbProvider),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists("openstack_lb_loadbalancer_v2.loadbalancer_1", &lb),
				),
			},
			{
				Config: testAccLbV2LoadBalancerConfigUpdate(lbProvider),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_lb_loadbalancer_v2.loadbalancer_1", "name", "loadbalancer_1_updated"),
					resource.TestMatchResourceAttr(
						"openstack_lb_loadbalancer_v2.loadbalancer_1", "vip_port_id",
						regexp.MustCompile("^[a-f0-9-]+")),
				),
			},
		},
	})
}

func TestAccLBV2LoadBalancer_secGroup(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	var sg1, sg2 groups.SecGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2LoadBalancerSecGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(
						"openstack_lb_loadbalancer_v2.loadbalancer_1", &lb),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &sg1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &sg2),
					resource.TestCheckResourceAttr(
						"openstack_lb_loadbalancer_v2.loadbalancer_1", "security_group_ids.#", "1"),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg1),
				),
			},
			{
				Config: testAccLbV2LoadBalancerSecGroupUpdate1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(
						"openstack_lb_loadbalancer_v2.loadbalancer_1", &lb),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &sg1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &sg2),
					resource.TestCheckResourceAttr(
						"openstack_lb_loadbalancer_v2.loadbalancer_1", "security_group_ids.#", "2"),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg1),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg2),
				),
			},
			{
				Config: testAccLbV2LoadBalancerSecGroupUpdate2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(
						"openstack_lb_loadbalancer_v2.loadbalancer_1", &lb),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &sg1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &sg2),
					resource.TestCheckResourceAttr(
						"openstack_lb_loadbalancer_v2.loadbalancer_1", "security_group_ids.#", "1"),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg2),
				),
			},
		},
	})
}

func TestAccLBV2LoadBalancer_vip_network(t *testing.T) {
	var lb loadbalancers.LoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
			testAccPreCheckUseOctavia(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2LoadBalancerConfigVIPNetwork,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists("openstack_lb_loadbalancer_v2.loadbalancer_1", &lb),
				),
			},
		},
	})
}

func TestAccLBV2LoadBalancer_vip_port_id(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
			testAccPreCheckUseOctavia(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2LoadBalancerConfigVIPPortID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(
						"openstack_lb_loadbalancer_v2.loadbalancer_1", &lb),
					testAccCheckNetworkingV2PortExists(
						"openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttrPtr(
						"openstack_lb_loadbalancer_v2.loadbalancer_1", "vip_port_id", &port.ID),
				),
			},
		},
	})
}

func testAccCheckLBV2LoadBalancerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := chooseLBV2AccTestClient(config, osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_lb_loadbalancer_v2" {
			continue
		}

		lb, err := loadbalancers.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil && lb.ProvisioningStatus != "DELETED" {
			return fmt.Errorf("LoadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2LoadBalancerExists(
	n string, lb *loadbalancers.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := chooseLBV2AccTestClient(config, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
		}

		found, err := loadbalancers.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*lb = *found

		return nil
	}
}

func testAccCheckLBV2LoadBalancerHasSecGroup(
	lb *loadbalancers.LoadBalancer, sg *groups.SecGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		port, err := ports.Get(networkingClient, lb.VipPortID).Extract()
		if err != nil {
			return err
		}

		for _, p := range port.SecurityGroups {
			if p == sg.ID {
				return nil
			}
		}

		return fmt.Errorf("LoadBalancer does not have the security group")
	}
}

func testAccLbV2LoadBalancerConfigBasic(lbProvider string) string {
	return fmt.Sprintf(`
    resource "openstack_networking_network_v2" "network_1" {
      name = "network_1"
      admin_state_up = "true"
    }

    resource "openstack_networking_subnet_v2" "subnet_1" {
      name = "subnet_1"
      cidr = "192.168.199.0/24"
      ip_version = 4
      network_id = "${openstack_networking_network_v2.network_1.id}"
    }

    resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
      name = "loadbalancer_1"
      loadbalancer_provider = "%s"
      vip_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"

      timeouts {
        create = "15m"
        update = "15m"
        delete = "15m"
      }
    }`, lbProvider)
}

func testAccLbV2LoadBalancerConfigUpdate(lbProvider string) string {
	return fmt.Sprintf(`
    resource "openstack_networking_network_v2" "network_1" {
      name = "network_1"
      admin_state_up = "true"
    }

    resource "openstack_networking_subnet_v2" "subnet_1" {
      name = "subnet_1"
      cidr = "192.168.199.0/24"
      ip_version = 4
      network_id = "${openstack_networking_network_v2.network_1.id}"
    }

    resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
      name = "loadbalancer_1_updated"
      loadbalancer_provider = "%s"
      admin_state_up = "true"
      vip_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"

      timeouts {
        create = "15m"
        update = "15m"
        delete = "15m"
      }
    }`, lbProvider)
}

const testAccLbV2LoadBalancerSecGroup = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
    name = "loadbalancer_1"
    vip_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    security_group_ids = [
      "${openstack_networking_secgroup_v2.secgroup_1.id}"
    ]

    timeouts {
      create = "15m"
      update = "15m"
      delete = "15m"
    }
}
`

const testAccLbV2LoadBalancerSecGroupUpdate1 = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  security_group_ids = [
    "${openstack_networking_secgroup_v2.secgroup_1.id}",
    "${openstack_networking_secgroup_v2.secgroup_2.id}"
  ]

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
`

const testAccLbV2LoadBalancerSecGroupUpdate2 = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  security_group_ids = [
    "${openstack_networking_secgroup_v2.secgroup_2.id}"
  ]
  depends_on = ["openstack_networking_secgroup_v2.secgroup_1"]

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
`

const testAccLbV2LoadBalancerConfigVIPNetwork = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  loadbalancer_provider = "octavia"
  vip_network_id = "${openstack_networking_network_v2.network_1.id}"
  depends_on = ["openstack_networking_subnet_v2.subnet_1"]
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
`

const testAccLbV2LoadBalancerConfigVIPPortID = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"
  depends_on = ["openstack_networking_subnet_v2.subnet_1"]
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  loadbalancer_provider = "octavia"
  vip_port_id = "${openstack_networking_port_v2.port_1.id}"
  depends_on = ["openstack_networking_port_v2.port_1"]
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
`
