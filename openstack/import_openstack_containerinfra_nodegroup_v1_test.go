package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccContainerInfraV1NodeGroupImport_basic(t *testing.T) {
	resourceName := "openstack_containerinfra_nodegroup_v1.nodegroup_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	imageName := acctest.RandomWithPrefix("tf-acc-image")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	nodeGroupName := acctest.RandomWithPrefix("tf-acc-cluster")

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
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id"},
			},
		},
	})
}
