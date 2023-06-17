package openstack

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccContainerInfraV1ClusterDataSource_basic(t *testing.T) {
	resourceName := "data.openstack_containerinfra_cluster_v1.cluster_1"
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
				Config: testAccContainerInfraV1ClusterBasic(keypairName, clusterTemplateName, clusterName, 1),
			},
			{
				Config: testAccContainerInfraV1ClusterDataSourceBasic(
					testAccContainerInfraV1ClusterBasic(keypairName, clusterTemplateName, clusterName, 1),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "api_address"),
					resource.TestCheckResourceAttrSet(resourceName, "coe_version"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_template_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_version"),
					resource.TestCheckResourceAttrSet(resourceName, "create_timeout"),
					resource.TestCheckResourceAttrSet(resourceName, "discovery_url"),
					resource.TestCheckResourceAttrSet(resourceName, "docker_volume_size"),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "keypair", keypairName),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "master_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "master_addresses.#", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_addresses.#", strconv.Itoa(1)),
					resource.TestCheckResourceAttrSet(resourceName, "stack_id"),
					resource.TestCheckResourceAttrSet(resourceName, "floating_ip_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig"),
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
