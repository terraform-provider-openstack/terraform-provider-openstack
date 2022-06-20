package openstack

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/nodegroups"
)

func TestAccContainerInfraV1NodeGroup_basic(t *testing.T) {
	var nodeGroup nodegroups.NodeGroup

	resourceName := "openstack_containerinfra_nodegroup_v1.nodegroup_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	imageName := acctest.RandomWithPrefix("tf-acc-image")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	nodeGroupName := acctest.RandomWithPrefix("tf-acc-nodegroup")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1NodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1NodeGroupUpdate(imageName, keypairName, clusterTemplateName, clusterName, nodeGroupName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1NodeGroupExists(resourceName, &nodeGroup),
					resource.TestCheckResourceAttr(resourceName, "name", nodeGroupName),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(1)),
				),
			},
			{
				Config: testAccContainerInfraV1NodeGroupUpdate(imageName, keypairName, clusterTemplateName, clusterName, nodeGroupName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1NodeGroupExists(resourceName, &nodeGroup),
					resource.TestCheckResourceAttr(resourceName, "name", nodeGroupName),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(2)),
				),
			},
		},
	})
}

func testAccCheckContainerInfraV1NodeGroupExists(n string, nodeGroup *nodegroups.NodeGroup) resource.TestCheckFunc {
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

		clusterID, nodeGroupID, err := parseVolumeTypeAccessID(rs.Primary.ID)
		if err != nil {
			return err
		}
		found, err := nodegroups.Get(containerInfraClient, clusterID, nodeGroupID).Extract()
		if err != nil {
			return err
		}

		if found.UUID != rs.Primary.ID {
			return fmt.Errorf("Cluster template not found")
		}

		*nodeGroup = *found

		return nil
	}
}

func testAccCheckContainerInfraV1NodeGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_containerinfra_nodegroup_v1" {
			continue
		}
		clusterID, nodeGroupID, err := parseVolumeTypeAccessID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = nodegroups.Get(containerInfraClient, clusterID, nodeGroupID).Extract()
		if err == nil {
			return fmt.Errorf("node group still exists")
		}
	}

	return nil
}
func testAccContainerInfraV1NodeGroupUpdate(imageName, keypairName, clusterTemplateName, clusterName string, nodeGroupName string, nodeCount int) string {
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

resource "openstack_compute_keypair_v2" "keypair_1" {
  name = "%s"
}

resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name       = "%s"
  image      = "${openstack_images_image_v2.image_1.id}"
  coe        = "kubernetes"
  http_proxy = "127.0.0.1:8801"
}

resource "openstack_containerinfra_cluster_v1" "cluster_1" {
  name                 = "%s"
  cluster_template_id  = "${openstack_containerinfra_clustertemplate_v1.clustertemplate_1.id}"
  master_count         = 1
  node_count           = 1
  keypair              = "${openstack_compute_keypair_v2.keypair_1.name}"
}

resource "openstack_containerinfra_nodegroup_v1" "nodegroup_1" {
  name                 = "%s"
  cluster_id           = "${openstack_containerinfra_cluster_v1.cluster_1.id}"
  node_count           = %d
}
`, imageName, keypairName, clusterTemplateName, clusterName, nodeGroupName, nodeCount)
}
