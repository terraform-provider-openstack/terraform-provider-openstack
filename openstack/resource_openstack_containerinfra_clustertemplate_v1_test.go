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
					resource.TestCheckResourceAttr(resourceName, "name", clusterTemplateName),
					resource.TestCheckResourceAttr(resourceName, "image", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "coe", "kubernetes"),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "floating_ip_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "volume_driver", "cinder"),
					resource.TestCheckResourceAttr(resourceName, "docker_storage_driver", "overlay2"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", "5"),
					resource.TestCheckResourceAttr(resourceName, "external_network_id", osExtGwID),
					resource.TestCheckResourceAttr(resourceName, "network_driver", "flannel"),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", osMagnumHTTPProxy),
					resource.TestCheckResourceAttr(resourceName, "https_proxy", osMagnumHTTPSProxy),
					resource.TestCheckResourceAttr(resourceName, "no_proxy", osMagnumNoProxy),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_tag", "1.11.1"),
					resource.TestCheckResourceAttr(resourceName, "labels.prometheus_monitoring", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.influx_grafana_dashboard_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_dashboard_enabled", "true"),
				),
			},
			{
				Config: testAccContainerInfraV1ClusterTemplateUpdate(clusterTemplateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", clusterTemplateName),
					resource.TestCheckResourceAttr(resourceName, "image", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "coe", "kubernetes"),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "floating_ip_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "volume_driver", "cinder"),
					resource.TestCheckResourceAttr(resourceName, "docker_storage_driver", "overlay2"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", "10"),
					resource.TestCheckResourceAttr(resourceName, "external_network_id", osExtGwID),
					resource.TestCheckResourceAttr(resourceName, "network_driver", "calico"),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", osMagnumHTTPProxy),
					resource.TestCheckResourceAttr(resourceName, "https_proxy", osMagnumHTTPSProxy),
					resource.TestCheckResourceAttr(resourceName, "no_proxy", osMagnumNoProxy),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_tag", "1.12.1"),
					resource.TestCheckResourceAttr(resourceName, "labels.prometheus_monitoring", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.influx_grafana_dashboard_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_dashboard_enabled", "false"),
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
  name                  = "%s"
  image                 = "%s"
  coe                   = "kubernetes"
  master_flavor         = "%s"
  flavor                = "%s"
  floating_ip_enabled   = true
  volume_driver         = "cinder"
  docker_storage_driver = "overlay2"
  docker_volume_size    = 5
  external_network_id   = "%s"
  network_driver        = "flannel"
  http_proxy            = "%s"
  https_proxy           = "%s"
  no_proxy              = "%s"
  labels = {
    kube_tag                         = "1.11.1"
    prometheus_monitoring            = "true"
    influx_grafana_dashboard_enabled = "true"
    kube_dashboard_enabled           = "true"
  }
}
`, clusterTemplateName, osMagnumImage, osMagnumFlavor, osMagnumFlavor, osExtGwID, osMagnumHTTPProxy, osMagnumHTTPSProxy, osMagnumNoProxy)
}

func testAccContainerInfraV1ClusterTemplateUpdate(clusterTemplateName string) string {
	return fmt.Sprintf(`
resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name                  = "%s"
  image                 = "%s"
  coe                   = "kubernetes"
  master_flavor         = "%s"
  flavor                = "%s"
  floating_ip_enabled   = false
  volume_driver         = "cinder"
  docker_storage_driver = "overlay2"
  docker_volume_size    = 10
  external_network_id   = "%s"
  network_driver        = "calico"
  http_proxy            = "%s"
  https_proxy           = "%s"
  no_proxy              = "%s"
  labels = {
    kube_tag                         = "1.12.1"
    prometheus_monitoring            = "true"
    influx_grafana_dashboard_enabled = "true"
    kube_dashboard_enabled           = "false"
  }
}
`, clusterTemplateName, osMagnumImage, osMagnumFlavor, osMagnumFlavor, osExtGwID, osMagnumHTTPProxy, osMagnumHTTPSProxy, osMagnumNoProxy)
}
