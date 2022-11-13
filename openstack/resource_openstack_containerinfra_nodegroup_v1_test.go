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
				Config: testAccContainerInfraV1NodeGroupUpdate(keypairName, clusterTemplateName, clusterName, nodeGroupName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1NodeGroupExists(resourceName, &nodeGroup),
					resource.TestCheckResourceAttr(resourceName, "name", nodeGroupName),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(1)),
				),
			},
			{
				Config: testAccContainerInfraV1NodeGroupUpdate(keypairName, clusterTemplateName, clusterName, nodeGroupName, 2),
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

		containerInfraClient.Microversion = containerInfraV1NodeGroupMinMicroversion
		clusterID, nodeGroupID, err := parseNodeGroupID(rs.Primary.ID)
		if err != nil {
			return err
		}
		found, err := nodegroups.Get(containerInfraClient, clusterID, nodeGroupID).Extract()
		if err != nil {
			return err
		}

		if found.UUID != nodeGroupID {
			return fmt.Errorf("Nodegroup not found")
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

	containerInfraClient.Microversion = containerInfraV1NodeGroupMinMicroversion

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
func testAccContainerInfraV1NodeGroupUpdate(keypairName, clusterTemplateName, clusterName string, nodeGroupName string, nodeCount int) string {
	return fmt.Sprintf(`
resource "openstack_compute_keypair_v2" "keypair_1" {
  name = "%s"
}

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
    kubescheduler_options = "log-flush-frequency=1m",
  }
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
`, keypairName, clusterTemplateName, osMagnumImage, osMagnumFlavor, osMagnumFlavor, osExtGwID, osMagnumHttpProxy, osMagnumHttpsProxy, osMagnumNoProxy, clusterName, nodeGroupName, nodeCount)
}
