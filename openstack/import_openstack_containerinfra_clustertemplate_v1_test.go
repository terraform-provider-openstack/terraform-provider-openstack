package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccContainerInfraV1ClusterTemplateImportBasic(t *testing.T) {
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")
	imageName := acctest.RandomWithPrefix("tf-acc-image")
	resourceName := fmt.Sprintf("openstack_containerinfra_clustertemplate_v1.%s", clusterTemplateName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterTemplateBasic(clusterTemplateName, imageName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image"},
			},
		},
	})
}
