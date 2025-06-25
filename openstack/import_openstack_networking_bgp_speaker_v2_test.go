package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkingV2BGPSpeaker_importBasic(t *testing.T) {
	resourceName := "openstack_networking_bgp_speaker_v2.speaker_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2BGPSpeakerDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2BGPSpeakerBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
