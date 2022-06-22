package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccContainerInfraV1NodeGroupDataSource_basic(t *testing.T) {
	resourceName := "openstack_containerinfra_nodegroup_v1.nodegroup_1"
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
				Config: testAccContainerInfraV1NodeGroupUpdate(keypairName, clusterTemplateName, clusterName, nodeGroupName, 1),
			},
			{
				Config: testAccContainerInfraV1NodeGroupDataSourceBasic(
					testAccContainerInfraV1NodeGroupUpdate(keypairName, clusterTemplateName, clusterName, nodeGroupName, 1),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1NodeGroupDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nodeGroupName),
					resource.TestCheckResourceAttr(resourceName, "node_count", "1"),
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

func testAccContainerInfraV1NodeGroupDataSourceBasic(clusterResource string) string {
	return fmt.Sprintf(`
%s

data "openstack_containerinfra_nodegroup_v1" "cluster_1" {
  cluster_id = "${openstack_containerinfra_cluster_v1.cluster_1.name}"
  name = "${openstack_containerinfra_nodegroup_v1.node_group_1.name}"
}
`, clusterResource)
}
