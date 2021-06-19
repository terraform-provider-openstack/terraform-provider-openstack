package openstack

import (
	"fmt"
	"strconv"
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
	imageName := acctest.RandomWithPrefix("tf-acc-image")
	dockerVolumeSize := 5

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
				Config: testAccContainerInfraV1ClusterTemplateBasic(clusterTemplateName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterTemplateExists(resourceName, &clusterTemplate),
					resource.TestCheckResourceAttr(resourceName, "name", clusterTemplateName),
					resource.TestCheckResourceAttr(resourceName, "coe", "kubernetes"),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", "127.0.0.1:8801"),
				),
			},
			{
				Config: testAccContainerInfraV1ClusterTemplateUpdate(clusterTemplateName, imageName, dockerVolumeSize),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", clusterTemplateName),
					resource.TestCheckResourceAttr(resourceName, "coe", "kubernetes"),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", ""),
					resource.TestCheckResourceAttr(resourceName, "docker_storage_driver", "devicemapper"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(dockerVolumeSize)),
				),
			},
		},
	})
}

func TestAccContainerInfraV1ClusterTemplate_labels(t *testing.T) {
	resourceName := "openstack_containerinfra_clustertemplate_v1.clustertemplate_1"
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	imageName := acctest.RandomWithPrefix("tf-acc-image")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckContainerInfra(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterTemplateLabels(clusterTemplateName, imageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "labels.kube_tag", "1.11.1"),
					resource.TestCheckResourceAttr(resourceName, "labels.prometheus_monitoring", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.influx_grafana_dashboard_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_dashboard_enabled", "true"),
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

func testAccContainerInfraV1ClusterTemplateBasic(clusterTemplateName, imageName string) string {
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
}
`, imageName, clusterTemplateName)
}

func testAccContainerInfraV1ClusterTemplateUpdate(clusterTemplateName, imageName string, dockerVolumeSize int) string {
	return fmt.Sprintf(`
resource "openstack_images_image_v2" "image_1" {
  name   = "%s"
  image_source_url = "https://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
  container_format = "bare"
  disk_format = "raw"
  properties = {
    os_distro = "fedora-atomic"
  }

  timeouts {
    create = "10m"
  }
}

resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name = "%s"
  image = "${openstack_images_image_v2.image_1.id}"
  coe = "kubernetes"
  docker_storage_driver = "devicemapper"
  docker_volume_size = %d
}
`, imageName, clusterTemplateName, dockerVolumeSize)
}

func testAccContainerInfraV1ClusterTemplateLabels(clusterTemplateName, imageName string) string {
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
  name   = "%s"
  image  = "${openstack_images_image_v2.image_1.id}"
  coe    = "kubernetes"
  labels = {
    kube_tag                         = "1.11.1"
    prometheus_monitoring            = "true"
    influx_grafana_dashboard_enabled = "true"
    kube_dashboard_enabled           = "true"
  }
}
`, imageName, clusterTemplateName)
}
