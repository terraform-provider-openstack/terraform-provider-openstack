package openstack

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2PortSecgroupAssociate_import(t *testing.T) {
	resourceName := "openstack_networking_port_secgroup_associate_v2.port_1"

	if os.Getenv("TF_ACC") != "" {
		hiddenPort, err := testAccCheckNetworkingV2PortSecGroupCreatePort(t, "hidden_port", true)
		if err != nil {
			t.Fatal(err)
		}
		defer testAccCheckNetworkingV2PortSecGroupDeletePort(t, hiddenPort) //nolint:errcheck
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate0(),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_group_ids"},
			},
		},
	})
}
