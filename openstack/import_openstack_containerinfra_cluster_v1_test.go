package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccContainerInfraV1ClusterImport_basic(t *testing.T) {
	resourceName := "openstack_containerinfra_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	imageName := acctest.RandomWithPrefix("tf-acc-image")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckContainerInfra(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerInfraV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterBasic(imageName, keypairName, clusterTemplateName, clusterName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keypair", "cluster_template_id"},
			},
		},
	})
}
