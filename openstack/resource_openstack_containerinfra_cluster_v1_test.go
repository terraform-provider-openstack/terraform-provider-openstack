package openstack

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
)

func TestAccContainerInfraV1Cluster_basic(t *testing.T) {
	var cluster clusters.Cluster

	resourceName := "openstack_containerinfra_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	imageName := acctest.RandomWithPrefix("tf-acc-image")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterBasic(imageName, keypairName, clusterTemplateName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "master_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "keypair", keypairName),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(5)),
				),
			},
			{
				Config: testAccContainerInfraV1ClusterUpdate(imageName, keypairName, clusterTemplateName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "master_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(resourceName, "keypair", keypairName),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(5)),
				),
			},
		},
	})
}

func testAccCheckContainerInfraV1ClusterExists(n string, cluster *clusters.Cluster) resource.TestCheckFunc {
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

		found, err := clusters.Get(containerInfraClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.UUID != rs.Primary.ID {
			return fmt.Errorf("Cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func testAccCheckContainerInfraV1ClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_containerinfra_cluster_v1" {
			continue
		}

		_, err := clusters.Get(containerInfraClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccContainerInfraV1ClusterBasic(imageName, keypairName, clusterTemplateName, clusterName string) string {
	return fmt.Sprintf(`
resource "openstack_images_image_v2" "image_1" {
  name             = "%s"
  image_source_url = "https://dl.fedoraproject.org/pub/fedora/linux/releases/27/CloudImages/x86_64/images/Fedora-Atomic-27-1.6.x86_64.qcow2"
  container_format = "bare"
  disk_format      = "qcow2"
  properties = {
    os_distro = "fedora-atomic"
  }
}

resource "openstack_compute_keypair_v2" "keypair_1" {
  name = "%s"
}

resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name                  = "%s"
  image                 = "${openstack_images_image_v2.image_1.name}"
  coe                   = "kubernetes"
  master_flavor         = "%s"
  flavor                = "%s"
  floating_ip_enabled   = true
  volume_driver         = "cinder"
  docker_storage_driver = "devicemapper"
  external_network_id   = "%s"
  network_driver        = "flannel"
  labels = {
    kubescheduler_options = "log-flush-frequency=1m"
  }
}

resource "openstack_containerinfra_cluster_v1" "cluster_1" {
  name                 = "%s"
  cluster_template_id  = "${openstack_containerinfra_clustertemplate_v1.clustertemplate_1.id}"
  master_count         = 1
  node_count           = 1
  keypair              = "${openstack_compute_keypair_v2.keypair_1.name}"
}
`, imageName, keypairName, clusterTemplateName, osMagnumFlavor, osMagnumFlavor, osExtGwID, clusterName)
}

func testAccContainerInfraV1ClusterUpdate(imageName, keypairName, clusterTemplateName, clusterName string) string {
	return fmt.Sprintf(`
resource "openstack_images_image_v2" "image_1" {
  name             = "%s"
  image_source_url = "https://dl.fedoraproject.org/pub/fedora/linux/releases/27/CloudImages/x86_64/images/Fedora-Atomic-27-1.6.x86_64.qcow2"
  container_format = "bare"
  disk_format      = "qcow2"
  properties = {
    os_distro = "fedora-atomic"
  }
}

resource "openstack_compute_keypair_v2" "keypair_1" {
  name = "%s"
}

resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name                  = "%s"
  image                 = "${openstack_images_image_v2.image_1.name}"
  coe                   = "kubernetes"
  master_flavor         = "%s"
  flavor                = "%s"
  floating_ip_enabled   = false
  external_network_id   = "%s"
  network_driver        = "flannel"
  labels = {
    kubescheduler_options = "log-flush-frequency=1m"
  }
}

resource "openstack_containerinfra_cluster_v1" "cluster_1" {
  name                 = "%s"
  cluster_template_id  = "${openstack_containerinfra_clustertemplate_v1.clustertemplate_1.id}"
  master_count         = 1
  node_count           = 2
  keypair              = "${openstack_compute_keypair_v2.keypair_1.name}"
}
`, imageName, keypairName, clusterTemplateName, osMagnumFlavor, osMagnumFlavor, osExtGwID, clusterName)
}
