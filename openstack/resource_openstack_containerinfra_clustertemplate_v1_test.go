package openstack

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccContainerInfraV1ClusterTemplateBasic(t *testing.T) {
	var clusterTemplate clustertemplates.ClusterTemplate

	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	imageName := acctest.RandomWithPrefix("tf-acc-image")
	dockerVolumeSize := 5
	resourceName := fmt.Sprintf("openstack_containerinfra_clustertemplate_v1.%s", clusterTemplateName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerInfraV1ClusterTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerInfraV1ClusterTemplateBasic(clusterTemplateName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterTemplateExists(resourceName, &clusterTemplate),
					resource.TestCheckResourceAttr(resourceName, "name", clusterTemplateName),
					resource.TestCheckResourceAttr(resourceName, "coe", "kubernetes"),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", "127.0.0.1:8801"),
				),
			},
			resource.TestStep{
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

func TestAccContainerInfraV1ClusterTemplateLabels(t *testing.T) {

	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	imageName := acctest.RandomWithPrefix("tf-acc-image")
	resourceName := fmt.Sprintf("openstack_containerinfra_clustertemplate_v1.%s", clusterTemplateName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerInfraV1ClusterTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
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
		containerInfraClient, err := config.containerInfraV1Client(OS_REGION_NAME)
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
	containerInfraClient, err := config.containerInfraV1Client(OS_REGION_NAME)
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
resource "openstack_images_image_v2" "%s" {
	name   = "%s"
	image_source_url = "https://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
	container_format = "bare"
	disk_format = "raw"
  properties {
    os_distro = "fedora-atomic"
  }

	timeouts {
		create = "10m"
	}
}

resource "openstack_containerinfra_clustertemplate_v1" "%s" {
  name = "%s"
  image = "${openstack_images_image_v2.%s.id}"
	coe = "kubernetes"
	http_proxy = "127.0.0.1:8801"
}
`, imageName, imageName, clusterTemplateName, clusterTemplateName, imageName)
}

func testAccContainerInfraV1ClusterTemplateUpdate(clusterTemplateName, imageName string, dockerVolumeSize int) string {
	return fmt.Sprintf(`
resource "openstack_images_image_v2" "%s" {
  name   = "%s"
  image_source_url = "https://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
  container_format = "bare"
	disk_format = "raw"
  properties {
    os_distro = "fedora-atomic"
  }

  timeouts {
    create = "10m"
  }
}

resource "openstack_containerinfra_clustertemplate_v1" "%s" {
  name = "%s"
  image = "${openstack_images_image_v2.%s.id}"
  coe = "kubernetes"
  docker_storage_driver = "devicemapper"
  docker_volume_size = %d
}
`, imageName, imageName, clusterTemplateName, clusterTemplateName, imageName, dockerVolumeSize)
}

func testAccContainerInfraV1ClusterTemplateLabels(clusterTemplateName, imageName string) string {
	return fmt.Sprintf(`
resource "openstack_images_image_v2" "%s" {
  name   = "%s"
  image_source_url = "https://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
  container_format = "bare"
  disk_format = "raw"
  properties {
    os_distro = "fedora-atomic"
  }

  timeouts {
    create = "10m"
  }
}

resource "openstack_containerinfra_clustertemplate_v1" "%s" {
  name = "%s"
  image = "${openstack_images_image_v2.%s.id}"
  coe = "kubernetes"
  labels = {
		kube_tag = "1.11.1"
		prometheus_monitoring = "true"
		influx_grafana_dashboard_enabled = "true"
		kube_dashboard_enabled = "true"
	}
}
`, imageName, imageName, clusterTemplateName, clusterTemplateName, imageName)
}
