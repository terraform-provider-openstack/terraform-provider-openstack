package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDNSQuotaV2_importBasic(t *testing.T) {
	resourceName := "openstack_dns_quota_v2.quota_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckDNS(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3ProjectDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSQuotaV2Basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
