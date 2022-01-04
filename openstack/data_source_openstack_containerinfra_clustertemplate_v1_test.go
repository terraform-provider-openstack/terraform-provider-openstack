package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccContainerInfraV1ClusterTemplateDataSource_basic(t *testing.T) {
	resourceName := "data.openstack_containerinfra_clustertemplate_v1.clustertemplate_1"
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	imageName := acctest.RandomWithPrefix("tf-acc-image")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckContainerInfra(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterTemplateDataSource(clusterTemplateName, imageName),
			},
			{
				Config: testAccContainerInfraV1ClusterTemplateDataSourceBasic(
					testAccContainerInfraV1ClusterTemplateDataSource(clusterTemplateName, imageName),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterTemplateDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterTemplateName),
					resource.TestCheckResourceAttr(resourceName, "coe", "kubernetes"),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", "127.0.0.1:8801"),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_tag", "1.11.1"),
					resource.TestCheckResourceAttr(resourceName, "docker_storage_driver", "devicemapper"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", "5"),
					resource.TestCheckResourceAttr(resourceName, "labels.prometheus_monitoring", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.influx_grafana_dashboard_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_dashboard_enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckContainerInfraV1ClusterTemplateDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ct, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cluster template data source: %s", n)
		}

		if ct.Primary.ID == "" {
			return fmt.Errorf("Cluster template data source ID not set")
		}

		return nil
	}
}

func testAccContainerInfraV1ClusterTemplateDataSource(clusterTemplateName, imageName string) string {
	return fmt.Sprintf(`
resource "openstack_images_image_v2" "image_1" {
  name             = "%s"
  image_source_url = "https://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
  container_format = "bare"
  disk_format      = "raw"
  properties = {
    os_distro = "fedora-atomic"
  }

  timeouts {
    create = "10m"
  }
}

resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name       = "%s"
  image      = "${openstack_images_image_v2.image_1.id}"
  coe        = "kubernetes"
  http_proxy = "127.0.0.1:8801"
  docker_storage_driver = "devicemapper"
  docker_volume_size = 5
  labels = {
    kube_tag                         = "1.11.1"
    prometheus_monitoring            = "true"
    influx_grafana_dashboard_enabled = "true"
    kube_dashboard_enabled           = "true"
  }
}
`, imageName, clusterTemplateName)
}

func testAccContainerInfraV1ClusterTemplateDataSourceBasic(clusterTemplateResource string) string {
	return fmt.Sprintf(`
%s

data "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name = "${openstack_containerinfra_clustertemplate_v1.clustertemplate_1.name}"
}
`, clusterTemplateResource)
}
