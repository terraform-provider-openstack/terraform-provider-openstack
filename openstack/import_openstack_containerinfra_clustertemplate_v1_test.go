package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContainerInfraV1ClusterTemplateImportBasic(t *testing.T) {
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	resourceName := "openstack_containerinfra_clustertemplate_v1.clustertemplate_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterTemplateDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterTemplateBasic(clusterTemplateName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
