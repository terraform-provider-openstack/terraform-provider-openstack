package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
)

func TestAccContainerInfraV1ClusterTemplate_basic(t *testing.T) {
	var clusterTemplate clustertemplates.ClusterTemplate

	resourceName := "openstack_containerinfra_clustertemplate_v1.clustertemplate_1"
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterTemplateBasic(clusterTemplateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterTemplateExists(resourceName, &clusterTemplate),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterTemplateName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "apiserver_port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "coe", "kubernetes"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_distro"),
					resource.TestCheckResourceAttr(resourceName, "dns_nameserver", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "docker_storage_driver", "overlay2"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", "5"),
					resource.TestCheckResourceAttr(resourceName, "external_network_id", osExtGwID),
					resource.TestCheckResourceAttr(resourceName, "fixed_network", "cluster-network"),
					resource.TestCheckResourceAttr(resourceName, "fixed_subnet", "cluster-network-subnet"),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "floating_ip_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", osMagnumHTTPProxy),
					resource.TestCheckResourceAttr(resourceName, "https_proxy", osMagnumHTTPSProxy),
					resource.TestCheckResourceAttr(resourceName, "image", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "insecure_registry", "registry.example.com"),
					resource.TestCheckResourceAttr(resourceName, "keypair_id", "my-keypair"),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_tag", "1.11.1"),
					resource.TestCheckResourceAttr(resourceName, "labels.prometheus_monitoring", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.influx_grafana_dashboard_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_dashboard_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "master_lb_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "network_driver", "flannel"),
					resource.TestCheckResourceAttr(resourceName, "no_proxy", osMagnumNoProxy),
					resource.TestCheckResourceAttrSet(resourceName, "public"),
					resource.TestCheckResourceAttr(resourceName, "registry_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "server_type", "vm"),
					resource.TestCheckResourceAttr(resourceName, "tls_disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "volume_driver", "cinder"),
					resource.TestCheckResourceAttr(resourceName, "hidden", "false"),
				),
			},
			{
				Config: testAccContainerInfraV1ClusterTemplateUpdate(clusterTemplateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterTemplateName+"-updated"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "apiserver_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "coe", "kubernetes"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_distro"),
					resource.TestCheckResourceAttr(resourceName, "dns_nameserver", "9.9.9.9"),
					resource.TestCheckResourceAttr(resourceName, "docker_storage_driver", "overlay"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", "10"),
					resource.TestCheckResourceAttr(resourceName, "external_network_id", osExtGwID),
					resource.TestCheckResourceAttr(resourceName, "fixed_network", "cluster-network2"),
					resource.TestCheckResourceAttr(resourceName, "fixed_subnet", "cluster-network2-subnet"),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "floating_ip_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", "http://1.2.3.4:8080"),
					resource.TestCheckResourceAttr(resourceName, "https_proxy", "http://1.2.3.4:8080"),
					resource.TestCheckResourceAttr(resourceName, "image", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "insecure_registry", "registry-new.example.com"),
					resource.TestCheckResourceAttr(resourceName, "keypair_id", "my-new-keypair"),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_tag", "1.12.1"),
					resource.TestCheckResourceAttr(resourceName, "labels.prometheus_monitoring", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.influx_grafana_dashboard_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_dashboard_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "master_lb_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_driver", "calico"),
					resource.TestCheckResourceAttr(resourceName, "no_proxy", "localhost"),
					resource.TestCheckResourceAttrSet(resourceName, "public"),
					resource.TestCheckResourceAttr(resourceName, "registry_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "server_type", "vm"),
					resource.TestCheckResourceAttr(resourceName, "tls_disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "volume_driver", "cinder"),
					resource.TestCheckResourceAttr(resourceName, "hidden", "true"),
				),
			},
		},
	})
}

func testAccCheckContainerInfraV1ClusterTemplateExists(n string, clustertemplate *clustertemplates.ClusterTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		containerInfraClient, err := config.ContainerInfraV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
		}

		found, err := clustertemplates.Get(containerInfraClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.UUID != rs.Primary.ID {
			return fmt.Errorf("Cluster template not found")
		}

		*clustertemplate = *found

		return nil
	}
}

func testAccCheckContainerInfraV1ClusterTemplateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_containerinfra_clustertemplate_v1" {
			continue
		}

		_, err := clustertemplates.Get(containerInfraClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Cluster template still exists")
		}
	}

	return nil
}

func testAccContainerInfraV1ClusterTemplateBasic(clusterTemplateName string) string {
	return fmt.Sprintf(`
resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  region                = "%s"
  name                  = "%s"
  apiserver_port        = "8888"
  coe                   = "kubernetes"
  dns_nameserver        = "8.8.8.8"
  docker_storage_driver = "overlay2"
  docker_volume_size    = 5
  external_network_id   = "%s"
  fixed_network         = "cluster-network"
  fixed_subnet          = "cluster-network-subnet"
  flavor                = "%s"
  master_flavor         = "%s"
  floating_ip_enabled   = true
  http_proxy            = "%s"
  https_proxy           = "%s"
  image                 = "%s"
  insecure_registry     = "registry.example.com"
  keypair_id            = "my-keypair"
  labels = {
    kube_tag                         = "1.11.1"
    prometheus_monitoring            = "true"
    influx_grafana_dashboard_enabled = "true"
    kube_dashboard_enabled           = "true"
  }
  master_lb_enabled     = "true"
  network_driver        = "flannel"
  no_proxy              = "%s"
  registry_enabled      = "true"
  server_type           = "vm"
  tls_disabled          = "false"
  volume_driver         = "cinder"
  hidden                = "false"
}
`, osRegionName, clusterTemplateName, osExtGwID, osMagnumFlavor, osMagnumFlavor, osMagnumHTTPProxy, osMagnumHTTPSProxy, osMagnumImage, osMagnumNoProxy)
}

func testAccContainerInfraV1ClusterTemplateUpdate(clusterTemplateName string) string {
	return fmt.Sprintf(`
resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  region                = "%s"
  name                  = "%s-updated"
  apiserver_port        = "8080"
  coe                   = "kubernetes"
  dns_nameserver        = "9.9.9.9"
  docker_storage_driver = "overlay"
  docker_volume_size    = 10
  external_network_id   = "%s"
  fixed_network         = "cluster-network2"
  fixed_subnet          = "cluster-network2-subnet"
  flavor                = "%s"
  master_flavor         = "%s"
  floating_ip_enabled   = false
  http_proxy            = "http://1.2.3.4:8080"
  https_proxy           = "http://1.2.3.4:8080"
  image                 = "%s"
  insecure_registry     = "registry-new.example.com"
  keypair_id            = "my-new-keypair"
  labels = {
    kube_tag                         = "1.12.1"
    prometheus_monitoring            = "true"
    influx_grafana_dashboard_enabled = "true"
    kube_dashboard_enabled           = "false"
  }
  master_lb_enabled     = "false"
  network_driver        = "calico"
  no_proxy              = "localhost"
  registry_enabled      = "false"
  server_type           = "vm"
  tls_disabled          = "true"
  volume_driver         = "cinder"
  hidden                = "true"
}
`, osRegionName, clusterTemplateName, osExtGwID, osMagnumFlavor, osMagnumFlavor, osMagnumImage)
}
