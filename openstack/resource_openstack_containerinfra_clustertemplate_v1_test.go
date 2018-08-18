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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerInfraV1ClusterTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerInfraV1ClusterTemplateBasic(clusterTemplateName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterTemplateExists("openstack_containerinfra_clustertemplate_v1.clustertemplate_1", &clusterTemplate),
					resource.TestCheckResourceAttr(
						"openstack_containerinfra_clustertemplate_v1.clustertemplate_1", "name", clusterTemplateName),
				),
			},
			resource.TestStep{
				Config: testAccContainerInfraV1ClusterTemplateUpdate(clusterTemplateName, imageName, dockerVolumeSize),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_containerinfra_clustertemplate_v1.clustertemplate_1", "name", clusterTemplateName),
					resource.TestCheckResourceAttr(
						"openstack_containerinfra_clustertemplate_v1.clustertemplate_1", "coe", "kubernetes"),
					resource.TestCheckResourceAttr(
						"openstack_containerinfra_clustertemplate_v1.clustertemplate_1", "docker_storage_driver", "devicemapper"),
					resource.TestCheckResourceAttr(
						"openstack_containerinfra_clustertemplate_v1.clustertemplate_1", "docker_volume_size", strconv.Itoa(dockerVolumeSize)),
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

resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name = "%s"
  image = "${openstack_images_image_v2.%s.id}"
	coe = "kubernetes"
}
`, imageName, imageName, clusterTemplateName, imageName)
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

resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name = "%s"
  image = "${openstack_images_image_v2.%s.id}"
  coe = "kubernetes"
  docker_storage_driver = "devicemapper"
  docker_volume_size = %d
}
`, imageName, imageName, clusterTemplateName, imageName, dockerVolumeSize)
}
