package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func TestAccNetworkingV2Subnet_basic(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2SubnetDnsConsistency("openstack_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pools.0.start", "192.168.199.100"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "description", "my subnet description"),
				),
			},
			{
				Config: testAccNetworkingV2Subnet_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "name", "subnet_1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "gateway_ip", "192.168.199.1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "enable_dhcp", "true"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pools.0.start", "192.168.199.150"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "description", ""),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_enableDHCP(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_enableDHCP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "enable_dhcp", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_disableDHCP(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_disableDHCP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "enable_dhcp", "false"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_noGateway(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_noGateway,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "gateway_ip", ""),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_impliedGateway(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_impliedGateway,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "gateway_ip", "192.168.199.1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_timeout(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_subnetPool(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_subnetPool,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_subnetPoolNoCIDR(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_subnetPoolNoCIDR,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_subnetPrefixLength(t *testing.T) {
	var subnet [2]subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_subnetPrefixLength,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet[0]),
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_2", &subnet[1]),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "prefix_length", "27"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_2", "prefix_length", "32"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_multipleAllocationPools(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_multipleAllocationPools_1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pools.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2Subnet_multipleAllocationPools_2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pools.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2Subnet_multipleAllocationPools_3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pools.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_allocationPool(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_allocationPool_1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pool.1804036869.start", "10.3.0.2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pool.1804036869.end", "10.3.0.255"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pool.1586215448.start", "10.3.255.0"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pool.1586215448.end", "10.3.255.254"),
				),
			},
			{
				Config: testAccNetworkingV2Subnet_allocationPool_2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pool.1804036869.start", "10.3.0.2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pool.1804036869.end", "10.3.0.255"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pool.2876574713.start", "10.3.255.10"),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "allocation_pool.2876574713.end", "10.3.255.154"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_clearDNSNameservers(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_clearDNSNameservers_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2SubnetDnsConsistency("openstack_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "dns_nameservers.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2Subnet_clearDNSNameservers_2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_subnet_v2.subnet_1", "dns_nameservers.#", "0"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SubnetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_subnet_v2" {
			continue
		}

		_, err := subnets.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2SubnetExists(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := subnets.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		*subnet = *found

		return nil
	}
}

func testAccCheckNetworkingV2SubnetDnsConsistency(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		for i, dns := range subnet.DNSNameservers {
			if dns != rs.Primary.Attributes[fmt.Sprintf("dns_nameservers.%d", i)] {
				return fmt.Errorf("Dns Nameservers list elements or order is not consistent")
			}
		}

		return nil
	}
}

const testAccNetworkingV2Subnet_basic = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  description = "my subnet description"
  cidr = "192.168.199.0/24"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  dns_nameservers = ["10.0.16.4", "213.186.33.99"]

  allocation_pools {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }
}
`

const testAccNetworkingV2Subnet_update = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  dns_nameservers = ["10.0.16.4", "213.186.33.99"]

  allocation_pools {
    start = "192.168.199.150"
    end = "192.168.199.200"
  }
}
`

const testAccNetworkingV2Subnet_enableDHCP = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  enable_dhcp = true
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2Subnet_disableDHCP = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  enable_dhcp = false
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2Subnet_noGateway = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  no_gateway = true
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2Subnet_impliedGateway = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}
resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2Subnet_timeout = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  allocation_pools {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingV2Subnet_subnetPool = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnetpool_v2" "subnetpool_1" {
  name = "my_ipv4_pool"
  prefixes = ["10.11.12.0/24"]
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "10.11.12.0/25"
  no_gateway = true
	network_id = "${openstack_networking_network_v2.network_1.id}"
	subnetpool_id = "${openstack_networking_subnetpool_v2.subnetpool_1.id}"
}
`

const testAccNetworkingV2Subnet_subnetPoolNoCIDR = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnetpool_v2" "subnetpool_1" {
  name = "my_ipv4_pool"
  prefixes = ["10.11.12.0/24"]
  min_prefixlen = "24"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
	network_id = "${openstack_networking_network_v2.network_1.id}"
	subnetpool_id = "${openstack_networking_subnetpool_v2.subnetpool_1.id}"
}
`

const testAccNetworkingV2Subnet_subnetPrefixLength = `
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnetpool_v2" "subnetpool_1" {
  name     = "my_ipv4_pool"
  prefixes = ["10.11.12.0/24"]
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name          = "subnet_1"
  prefix_length = 27
  enable_dhcp   = false
  network_id    = "${openstack_networking_network_v2.network_1.id}"
  subnetpool_id = "${openstack_networking_subnetpool_v2.subnetpool_1.id}"
}

resource "openstack_networking_subnet_v2" "subnet_2" {
  name          = "subnet_2"
  prefix_length = 32
  enable_dhcp   = false
  network_id    = "${openstack_networking_network_v2.network_1.id}"
  subnetpool_id = "${openstack_networking_subnetpool_v2.subnetpool_1.id}"
}
`

const testAccNetworkingV2Subnet_multipleAllocationPools_1 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "10.3.0.0/16"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  allocation_pools {
    start = "10.3.0.2"
    end = "10.3.0.255"
  }

  allocation_pools {
    start = "10.3.255.0"
    end = "10.3.255.254"
  }
}
`

const testAccNetworkingV2Subnet_multipleAllocationPools_2 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "10.3.0.0/16"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  allocation_pools {
    start = "10.3.255.0"
    end = "10.3.255.254"
  }

  allocation_pools {
    start = "10.3.0.2"
    end = "10.3.0.255"
  }
}
`

const testAccNetworkingV2Subnet_multipleAllocationPools_3 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "10.3.0.0/16"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  allocation_pools {
    start = "10.3.255.10"
    end = "10.3.255.154"
  }

  allocation_pools {
    start = "10.3.0.2"
    end = "10.3.0.255"
  }
}
`

const testAccNetworkingV2Subnet_allocationPool_1 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "10.3.0.0/16"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  allocation_pool {
    start = "10.3.0.2"
    end = "10.3.0.255"
  }

  allocation_pool {
    start = "10.3.255.0"
    end = "10.3.255.254"
  }
}
`

const testAccNetworkingV2Subnet_allocationPool_2 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "10.3.0.0/16"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  allocation_pool {
    start = "10.3.255.10"
    end = "10.3.255.154"
  }

  allocation_pool {
    start = "10.3.0.2"
    end = "10.3.0.255"
  }
}
`

const testAccNetworkingV2Subnet_clearDNSNameservers_1 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  dns_nameservers = ["10.0.16.4", "213.186.33.99"]

  allocation_pools {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }
}
`

const testAccNetworkingV2Subnet_clearDNSNameservers_2 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  allocation_pools {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }
}
`
