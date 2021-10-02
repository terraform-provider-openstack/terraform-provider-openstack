package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccContainerInfraV1ClusterDataSource_basic(t *testing.T) {
	resourceName := "openstack_containerinfra_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	imageName := acctest.RandomWithPrefix("tf-acc-image")
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
				Config: testAccContainerInfraV1ClusterBasic(imageName, keypairName, clusterTemplateName, clusterName),
			},
			{
				Config: testAccContainerInfraV1ClusterDataSourceBasic(
					testAccContainerInfraV1ClusterBasic(imageName, keypairName, clusterTemplateName, clusterName),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "master_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_count", "1"),
				),
			},
		},
	})
}

func testAccCheckContainerInfraV1ClusterDataSourceID(n string) resource.TestCheckFunc {
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

func testAccContainerInfraV1ClusterDataSourceBasic(clusterResource string) string {
	return fmt.Sprintf(`
%s

data "openstack_containerinfra_cluster_v1" "cluster_1" {
  name = "${openstack_containerinfra_cluster_v1.cluster_1.name}"
}
`, clusterResource)
}
