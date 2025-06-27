package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLBV2Listener_basic(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2ListenerDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2ListenerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists(t.Context(), "openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "connection_limit", "-1"),
				),
			},
			{
				Config: testAccLbV2ListenerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "name", "listener_1_updated"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "connection_limit", "100"),
				),
			},
		},
	})
}

func TestAccLBV2Listener_octavia(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2ListenerDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2ListenerConfigOctavia,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists(t.Context(), "openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "connection_limit", "5"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_client_data", "1000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_member_connect", "2000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_member_data", "3000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_tcp_inspect", "4000"),
				),
			},
			{
				Config: testAccLbV2ListenerConfigOctaviaUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "name", "listener_1_updated"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "connection_limit", "100"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_client_data", "4000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_member_connect", "3000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_member_data", "2000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_tcp_inspect", "1000"),
				),
			},
		},
	})
}

func TestAccLBV2Listener_octavia_udp(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2ListenerDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2ListenerConfigOctaviaUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists(t.Context(), "openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "protocol", "UDP"),
				),
			},
		},
	})
}

func TestAccLBV2ListenerConfig_octavia_insert_headers(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2ListenerDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2ListenerConfigOctaviaInsertHeaders1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists(t.Context(), "openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-For", "true"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-Port", "false"),
				),
			},
			{
				Config: testAccLbV2ListenerConfigOctaviaInsertHeaders2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists(t.Context(), "openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-For", "false"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-Port", "true"),
				),
			},
			{
				Config: testAccLbV2ListenerConfigOctavia,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists(t.Context(), "openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckNoResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-For"),
					resource.TestCheckNoResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-Port"),
				),
			},
		},
	})
}

func TestAccLBV2Listener_hsts(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Skip("Secret creation attempt not allowed - please review your user/project privileges")
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2ListenerDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2ListenerConfigHSTS,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists(t.Context(), "openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "connection_limit", "-1"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "protocol", "TERMINATED_HTTPS"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "client_authentication", "OPTIONAL"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "alpn_protocols.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "tls_versions.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "tls_ciphers", "TLS13-CHACHA20-POLY1305-SHA256"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "hsts_preload", "true"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "hsts_include_subdomains", "true"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "hsts_max_age", "31536000"),
				),
			},
			{
				Config: testAccLbV2ListenerConfigHSTSUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "name", "listener_1_updated"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "connection_limit", "100"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "protocol", "HTTP"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "client_authentication", "NONE"),
					// alpn_protocols reset to Octavia default value
					// resource.TestCheckResourceAttr(
					//	"openstack_lb_listener_v2.listener_1", "alpn_protocols.#", "0"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "tls_versions.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "tls_ciphers", ""),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "hsts_preload", "false"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "hsts_include_subdomains", "false"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "hsts_max_age", "0"),
				),
			},
		},
	})
}

func testAccCheckLBV2ListenerDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		lbClient, err := config.LoadBalancerV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_lb_listener_v2" {
				continue
			}

			_, err := listeners.Get(ctx, lbClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("Listener still exists: %s", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckLBV2ListenerExists(ctx context.Context, n string, listener *listeners.Listener) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		lbClient, err := config.LoadBalancerV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %w", err)
		}

		found, err := listeners.Get(ctx, lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Member not found")
		}

		*listener = *found

		return nil
	}
}

const testAccLbV2ListenerConfigBasic = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  default_pool_id = openstack_lb_pool_v2.pool_1.id
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigUpdate = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1_updated"
  protocol = "HTTP"
  protocol_port = 8080
  connection_limit = 100
  admin_state_up = "true"
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigOctavia = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  connection_limit = 5
  timeout_client_data = 1000
  timeout_member_connect = 2000
  timeout_member_data = 3000
  timeout_tcp_inspect = 4000
  default_pool_id = openstack_lb_pool_v2.pool_1.id
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
  tags = ["tag1"]

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigOctaviaUpdate = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1_updated"
  protocol = "HTTP"
  protocol_port = 8080
  connection_limit = 100
  timeout_client_data = 4000
  timeout_member_connect = 3000
  timeout_member_data = 2000
  timeout_tcp_inspect = 1000
  admin_state_up = "true"
  default_pool_id = openstack_lb_pool_v2.pool_1.id
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
  tags = ["tag1", "tag2"]

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigOctaviaUDP = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "UDP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "UDP"
  protocol_port = 53
  default_pool_id = openstack_lb_pool_v2.pool_1.id
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigOctaviaInsertHeaders1 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  default_pool_id = openstack_lb_pool_v2.pool_1.id
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id

  insert_headers = {
    X-Forwarded-For = "true"
    X-Forwarded-Port = "false"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigOctaviaInsertHeaders2 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  default_pool_id = openstack_lb_pool_v2.pool_1.id
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id

  insert_headers = {
    X-Forwarded-For = "false"
    X-Forwarded-Port = "true"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigHSTS = `
resource "openstack_keymanager_secret_v1" "certificate_1" {
  name                     = "certificate"
  payload                  = "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"
  payload_content_encoding = "base64"
  payload_content_type     = "application/octet-stream"
}

resource "openstack_keymanager_secret_v1" "certificate_2" {
  name                 = "certificate"
  payload              = <<EOF
-----BEGIN CERTIFICATE-----
MIIDbTCCAlWgAwIBAgIUEGBRJ6o6kXFbZWFO2sIkWNANAs4wDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAgFw0yNDA5MjQxNTQwMTRaGA8yMTI0
MDgzMTE1NDAxNFowRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUx
ITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAKR4zT0mnYTv/pZtPLHCgYZAuJR8VVy4Su9XuszE
mfDKSdg/UQmAQ1DCv0I+t3Wp3i2VAPW3734f4fcunYMFqbUn+v9wRk7N6hgnPOFW
P4fai93a1kcEjRqBzh0HfE/w4j+GC+AF4rDojJVX9SrR7TohF3oXa6ysjvB7Cjnx
Bg6uTGFIsPdAoPrgAIWT00pDnNPZBmX6IQGkQ2DNr/7O0ibfnWJVmvLaWsbSE20c
ks65iMYwIUgVOyLFs/KWAgi9j9mCRbYQYpvx8C8bz3jVpBhT9AEehssZ9ahX3BOj
AT1rsfA1zCTU3J6NUikWxZSjUW8WqbcQrf73ozkgM/ooOC8CAwEAAaNTMFEwHQYD
VR0OBBYEFGI++tfJ6IilXjqoOEHpcbRoO07NMB8GA1UdIwQYMBaAFGI++tfJ6Iil
XjqoOEHpcbRoO07NMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEB
ACYryWPD4ArYRvErj8ZbYvqxXGM4d+UtP5OdlsKL2y+6+J+XwCZ9rfcTHFQoQS1R
BnOR2LcIDR8yIJG01UnXmCjruRelsdAhXPBHK46KDaKkYc8l6V9iguzwpNqmfq+3
8PM9wp40cJaOHOriNltYw9PEw7HhF3A5vL0vrZy81cS+xPg+9ifhuSO+aBDpc9Lu
ortBUAxHqN3pPofupxpeJXIDT/M/xBw+fope90md0WVR/S+3bvdtEANPL2gGMKWY
0utZuMcHvuDc56XDTe/NaviLR5J63j5CjovFrntdW5607l+oHqLDlv/uR8EgKIBU
feFvlwLZ+xpvAhR+EkLsKLo=
-----END CERTIFICATE-----
EOF
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "TERMINATED_HTTPS"
  protocol_port = 8080
  default_pool_id = openstack_lb_pool_v2.pool_1.id
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id

  client_authentication = "OPTIONAL"
  client_ca_tls_container_ref = openstack_keymanager_secret_v1.certificate_2.secret_ref

  default_tls_container_ref = openstack_keymanager_secret_v1.certificate_1.secret_ref

  tls_versions   = ["TLSv1.2"]
  alpn_protocols = ["http/1.1"]
  tls_ciphers    = "TLS13-CHACHA20-POLY1305-SHA256"

  hsts_preload            = true
  hsts_max_age            = 31536000
  hsts_include_subdomains = true

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigHSTSUpdate = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1_updated"
  protocol = "HTTP"
  protocol_port = 8080
  connection_limit = 100
  admin_state_up = "true"
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id

  // we need to keep them, since Octavia has a bug
  tls_versions   = ["TLSv1.2"]

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`
