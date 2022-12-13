package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccContainerInfraV1NodeGroupDataSource_basic(t *testing.T) {
	resourceName := "data.openstack_containerinfra_nodegroup_v1.nodegroup_1"
	nodeGroupName := acctest.RandomWithPrefix("tf-acc-nodegroup")
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1NodeGroupBasic(keypairName, clusterTemplateName, clusterName, nodeGroupName, 1),
			},
			{
				Config: testAccContainerInfraV1NodeGroupDataSourceBasic(
					testAccContainerInfraV1NodeGroupBasic(keypairName, clusterTemplateName, clusterName, nodeGroupName, 1),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1NodeGroupDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "name", nodeGroupName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "docker_volume_size"),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttrSet(resourceName, "role"),
					resource.TestCheckResourceAttr(resourceName, "node_count", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "min_node_count"),
					resource.TestCheckResourceAttrSet(resourceName, "max_node_count"),
					resource.TestCheckResourceAttr(resourceName, "image", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
				),
			},
		},
	})
}

func testAccCheckContainerInfraV1NodeGroupDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ct, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cluster data source: %s", n)
		}

		if ct.Primary.ID == "" {
			return fmt.Errorf("Cluster data source ID is not set")
		}

		return nil
	}
}

func testAccContainerInfraV1NodeGroupDataSourceBasic(nodeGroupResource string) string {
	return fmt.Sprintf(`
%s

data "openstack_containerinfra_nodegroup_v1" "nodegroup_1" {
  cluster_id = "${openstack_containerinfra_cluster_v1.cluster_1.name}"
  name = "${openstack_containerinfra_nodegroup_v1.nodegroup_1.name}"
}
`, nodeGroupResource)
}
