package openstack

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/pagination"
)

func TestAccComputeV2Instance_basic(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "availability_zone", "nova"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_initialStateActive(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceStateActive(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "active"),
					testAccCheckComputeV2InstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccComputeV2InstanceStateShutoff(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "shutoff"),
					testAccCheckComputeV2InstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: testAccComputeV2InstanceStateActive(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "active"),
					testAccCheckComputeV2InstanceState(&instance, "active"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_initialStateShutoff(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceStateShutoff(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "shutoff"),
					testAccCheckComputeV2InstanceState(&instance, "shutoff"),
				),
			},
			{
				Config: testAccComputeV2InstanceStateActive(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "active"),
					testAccCheckComputeV2InstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccComputeV2InstanceStateShutoff(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "shutoff"),
					testAccCheckComputeV2InstanceState(&instance, "shutoff"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_initialShelve(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceStateActive(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "active"),
					testAccCheckComputeV2InstanceState(&instance, "active"),
				),
			},
			{
				Config: testAccComputeV2InstanceStateShelve(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "shelved_offloaded"),
					testAccCheckComputeV2InstanceState(&instance, "shelved_offloaded"),
				),
			},
			{
				Config: testAccComputeV2InstanceStateActive(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "power_state", "active"),
					testAccCheckComputeV2InstanceState(&instance, "active"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_secgroupMulti(t *testing.T) {
	var instance1 servers.Server
	var secgroup1 secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceSecgroupMulti(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists(
						"openstack_compute_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance1),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_secgroupMultiUpdate(t *testing.T) {
	var instance1 servers.Server
	var secgroup1, secgroup2 secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceSecgroupMultiUpdate1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists(
						"openstack_compute_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckComputeV2SecGroupExists(
						"openstack_compute_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance1),
				),
			},
			{
				Config: testAccComputeV2InstanceSecgroupMultiUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists(
						"openstack_compute_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckComputeV2SecGroupExists(
						"openstack_compute_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance1),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeImage(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBootFromVolumeImage(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeVolume(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBootFromVolumeVolume(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeForceNew(t *testing.T) {
	var instance1 servers.Server
	var instance2 servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBootFromVolumeForceNew1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance1),
				),
			},
			{
				Config: testAccComputeV2InstanceBootFromVolumeForceNew2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance2),
					testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(&instance1, &instance2),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_blockDeviceNewVolume(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBlockDeviceNewVolume(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_blockDeviceNewVolumeTypeAndBus(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBlockDeviceNewVolumeTypeAndBus(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_blockDeviceExistingVolume(t *testing.T) {
	var instance servers.Server
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceBlockDeviceExistingVolume(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckBlockStorageV3VolumeExists(
						"openstack_blockstorage_volume_v3.volume_1", &volume),
				),
			},
		},
	})
}

// TODO: verify the personality really exists on the instance.
func TestAccComputeV2Instance_personality(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstancePersonality(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_multiEphemeral(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceMultiEphemeral(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_accessIPv4(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceAccessIPv4(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "access_ip_v4", "192.168.1.100"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_changeFixedIP(t *testing.T) {
	var instance1 servers.Server
	var instance2 servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceChangeFixedIP1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance1),
				),
			},
			{
				Config: testAccComputeV2InstanceChangeFixedIP2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"openstack_compute_instance_v2.instance_1", &instance2),
					testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(&instance1, &instance2),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_stopBeforeDestroy(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceStopBeforeDestroy(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_metadataRemove(t *testing.T) {
	var instance servers.Server

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceMetadataRemove1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceMetadata(&instance, "abc", "def"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "all_metadata.abc", "def"),
				),
			},
			{
				Config: testAccComputeV2InstanceMetadataRemove2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceMetadata(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceMetadata(&instance, "ghi", "jkl"),
					testAccCheckComputeV2InstanceNoMetadataKey(&instance, "abc"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "all_metadata.foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "all_metadata.ghi", "jkl"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_forceDelete(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceForceDelete(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_timeout(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceTimeout(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_networkModeAuto(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceNetworkModeAuto(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceNetworkExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_networkModeNone(t *testing.T) {
	var instance servers.Server
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceNetworkModeNone(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeV2InstanceNetworkDoesNotExist("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_networkNameToID(t *testing.T) {
	var instance servers.Server
	var network networks.Network
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceNetworkNameToID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.1.uuid", &network.ID),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_crazyNICs(t *testing.T) {
	var instance servers.Server
	var network1 networks.Network
	var network2 networks.Network
	var port1 ports.Port
	var port2 ports.Port
	var port3 ports.Port
	var port4 ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceCrazyNICs(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2NetworkExists(
						"openstack_networking_network_v2.network_1", &network1),
					testAccCheckNetworkingV2NetworkExists(
						"openstack_networking_network_v2.network_2", &network2),
					testAccCheckNetworkingV2PortExists(
						"openstack_networking_port_v2.port_1", &port1),
					testAccCheckNetworkingV2PortExists(
						"openstack_networking_port_v2.port_2", &port2),
					testAccCheckNetworkingV2PortExists(
						"openstack_networking_port_v2.port_3", &port3),
					testAccCheckNetworkingV2PortExists(
						"openstack_networking_port_v2.port_4", &port4),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.1.uuid", &network1.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.2.uuid", &network2.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.3.uuid", &network1.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.4.uuid", &network2.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.5.uuid", &network1.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.6.uuid", &network2.ID),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.1.name", "network_1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.2.name", "network_2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.3.name", "network_1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.4.name", "network_2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.5.name", "network_1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.6.name", "network_2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.7.name", "network_1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.8.name", "network_2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.1.fixed_ip_v4", "192.168.1.100"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.2.fixed_ip_v4", "192.168.2.100"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.3.fixed_ip_v4", "192.168.1.101"),
					resource.TestCheckResourceAttr(
						"openstack_compute_instance_v2.instance_1", "network.4.fixed_ip_v4", "192.168.2.101"),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.5.port", &port1.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.6.port", &port2.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.7.port", &port3.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.8.port", &port4.ID),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_tags(t *testing.T) {
	var instance servers.Server

	resourceName := "openstack_compute_instance_v2.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceTagsCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceTags(resourceName, []string{"tag1", "tag2", "tag3"}),
				),
			},
			{
				Config: testAccComputeV2InstanceTagsAdd(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceTags(resourceName, []string{"tag1", "tag2", "tag3", "tag4"}),
				),
			},
			{
				Config: testAccComputeV2InstanceTagsDelete(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceTags(resourceName, []string{"tag2", "tag3"}),
				),
			},
			{
				Config: testAccComputeV2InstanceTagsClear(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceTags(resourceName, nil),
				),
			},
		},
	})
}

func testAccCheckComputeV2InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_compute_instance_v2" {
			continue
		}

		server, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "SOFT_DELETED" && server.Status != "DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeV2InstanceExists(n string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		found, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccCheckComputeV2InstanceMetadata(
	instance *servers.Server, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Metadata == nil {
			return fmt.Errorf("No metadata")
		}

		for key, value := range instance.Metadata {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, value)
		}

		return fmt.Errorf("Metadata not found: %s", k)
	}
}

func testAccCheckComputeV2InstanceNoMetadataKey(
	instance *servers.Server, k string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance.Metadata == nil {
			return nil
		}

		for key := range instance.Metadata {
			if k == key {
				return fmt.Errorf("Metadata found: %s", k)
			}
		}

		return nil
	}
}

func testAccCheckComputeV2InstanceBootVolumeAttachment(
	instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var attachments []volumeattach.VolumeAttachment

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return err
		}

		err = volumeattach.List(computeClient, instance.ID).EachPage(
			func(page pagination.Page) (bool, error) {
				actual, err := volumeattach.ExtractVolumeAttachments(page)
				if err != nil {
					return false, fmt.Errorf("Unable to lookup attachment: %s", err)
				}

				attachments = actual
				return true, nil
			})
		if err != nil {
			return fmt.Errorf("Unable to list volume attachments")
		}

		if len(attachments) == 1 {
			return nil
		}

		return fmt.Errorf("No attached volume found")
	}
}

func testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(
	instance1, instance2 *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance1.ID == instance2.ID {
			return fmt.Errorf("Instance was not recreated")
		}

		return nil
	}
}

func testAccCheckComputeV2InstanceState(
	instance *servers.Server, state string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if strings.ToLower(instance.Status) != state {
			return fmt.Errorf("Instance state is not match")
		}

		return nil
	}
}

func testAccCheckComputeV2InstanceTags(name string, tags []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		if _, ok := rs.Primary.Attributes["tags.#"]; !ok {
			return fmt.Errorf("resource tags not found: %s.tags", name)
		}

		var rtags []string
		for key, val := range rs.Primary.Attributes {
			if !strings.HasPrefix(key, "tags.") {
				continue
			}

			if key == "tags.#" {
				continue
			}

			rtags = append(rtags, val)
		}

		sort.Strings(rtags)
		sort.Strings(tags)
		if !reflect.DeepEqual(rtags, tags) {
			return fmt.Errorf(
				"%s.tags: expected: %#v, got %#v", name, tags, rtags)
		}
		return nil
	}
}

func testAccCheckComputeV2InstanceNetworkExists(n string, _ *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		networkCount, ok := rs.Primary.Attributes["network.#"]

		if !ok {
			return fmt.Errorf("network attributes not found: %s", n)
		}

		if networkCount != "1" {
			return fmt.Errorf("network should be exists when network mode 'auto': %s", n)
		}

		return nil
	}
}

func testAccCheckComputeV2InstanceNetworkDoesNotExist(n string, _ *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		networkCount, ok := rs.Primary.Attributes["network.#"]

		if !ok {
			return fmt.Errorf("network attributes not found: %s", n)
		}

		if networkCount != "0" {
			return fmt.Errorf("network should not exists when network mode 'none': %s", n)
		}

		return nil
	}
}

func testAccComputeV2InstanceBasic() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceSecgroupMulti() string {
	return fmt.Sprintf(`
resource "openstack_compute_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "a security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default", "${openstack_compute_secgroup_v2.secgroup_1.name}"]
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceSecgroupMultiUpdate1() string {
	return fmt.Sprintf(`
resource "openstack_compute_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "a security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "another security group"
  rule {
    from_port = 80
    to_port = 80
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceSecgroupMultiUpdate2() string {
	return fmt.Sprintf(`
resource "openstack_compute_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "a security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "another security group"
  rule {
    from_port = 80
    to_port = 80
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default", "${openstack_compute_secgroup_v2.secgroup_1.name}", "${openstack_compute_secgroup_v2.secgroup_2.name}"]
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceBootFromVolumeImage() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 5
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeV2InstanceBootFromVolumeVolume() string {
	return fmt.Sprintf(`
resource "openstack_blockstorage_volume_v3" "vol_1" {
  name = "vol_1"
  size = 5
  image_id = "%s"
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "${openstack_blockstorage_volume_v3.vol_1.id}"
    source_type = "volume"
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeV2InstanceBootFromVolumeForceNew1() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 5
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeV2InstanceBootFromVolumeForceNew2() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 4
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeV2InstanceBlockDeviceNewVolume() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    destination_type = "local"
    boot_index = 0
    delete_on_termination = true
  }
  block_device {
    source_type = "blank"
    destination_type = "volume"
    volume_size = 1
    boot_index = 1
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeV2InstanceBlockDeviceNewVolumeTypeAndBus() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    destination_type = "local"
    boot_index = 0
		delete_on_termination = true
		device_type = "disk"
		disk_bus = "virtio"
  }
  block_device {
    source_type = "blank"
    destination_type = "volume"
    volume_size = 1
    boot_index = 1
		delete_on_termination = true
		device_type = "disk"
		disk_bus = "virtio"
  }
  network {
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeV2InstanceBlockDeviceExistingVolume() string {
	return fmt.Sprintf(`
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  block_device {
    uuid = "%s"
    source_type = "image"
    destination_type = "local"
    boot_index = 0
    delete_on_termination = true
  }
  block_device {
    uuid = "${openstack_blockstorage_volume_v3.volume_1.id}"
    source_type = "volume"
    destination_type = "volume"
    boot_index = 1
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeV2InstancePersonality() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  personality {
    file = "/tmp/foobar.txt"
    content = "happy"
  }
  personality {
    file = "/tmp/barfoo.txt"
    content = "angry"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceMultiEphemeral() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "terraform-test"
  security_groups = ["default"]
  block_device {
    boot_index = 0
    delete_on_termination = true
    destination_type = "local"
    source_type = "image"
    uuid = "%s"
  }
  block_device {
    boot_index = -1
    delete_on_termination = true
    destination_type = "local"
    source_type = "blank"
    volume_size = 1
  }
  block_device {
    boot_index = -1
    delete_on_termination = true
    destination_type = "local"
    source_type = "blank"
    volume_size = 1
  }
  network {
    uuid = "%s"
  }
}
`, osImageID, osNetworkID)
}

func testAccComputeV2InstanceAccessIPv4() string {
	return fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "openstack_compute_instance_v2" "instance_1" {
  depends_on = ["openstack_networking_subnet_v2.subnet_1"]

  name = "instance_1"
  security_groups = ["default"]

  network {
    uuid = "%s"
  }

  network {
    uuid = "${openstack_networking_network_v2.network_1.id}"
    fixed_ip_v4 = "192.168.1.100"
    access_network = true
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceChangeFixedIP1() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
    fixed_ip_v4 = "10.0.0.24"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceChangeFixedIP2() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
    fixed_ip_v4 = "10.0.0.25"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceStopBeforeDestroy() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  stop_before_destroy = true
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceDetachPortsBeforeDestroy() string {
	return fmt.Sprintf(`

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = "%s"
  admin_state_up = "true"
}


resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  vendor_options {
    detach_ports_before_destroy = true
  }
  network {
    port = "${openstack_networking_port_v2.port_1.id}"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceMetadataRemove1() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  metadata = {
    foo = "bar"
    abc = "def"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceMetadataRemove2() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  metadata = {
    foo = "bar"
    ghi = "jkl"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceForceDelete() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  force_delete = true
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceTimeout() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]

  timeouts {
    create = "10m"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceNetworkModeAuto() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"

  network_mode = "auto"
}
`)
}

func testAccComputeV2InstanceNetworkModeNone() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "test-instance-1"

  network_mode = "none"
}
`)
}

func testAccComputeV2InstanceNetworkNameToID() string {
	return fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "openstack_compute_instance_v2" "instance_1" {
  depends_on = ["openstack_networking_subnet_v2.subnet_1"]

  name = "instance_1"
  security_groups = ["default"]

  network {
    uuid = "%s"
  }

  network {
    name = "${openstack_networking_network_v2.network_1.name}"
  }

}
`, osNetworkID)
}

func testAccComputeV2InstanceCrazyNICs() string {
	return fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "openstack_networking_network_v2" "network_2" {
  name = "network_2"
}

resource "openstack_networking_subnet_v2" "subnet_2" {
  name = "subnet_2"
  network_id = "${openstack_networking_network_v2.network_2.id}"
  cidr = "192.168.2.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.1.103"
  }
}

resource "openstack_networking_port_v2" "port_2" {
  name = "port_2"
  network_id = "${openstack_networking_network_v2.network_2.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_2.id}"
    ip_address = "192.168.2.103"
  }
}

resource "openstack_networking_port_v2" "port_3" {
  name = "port_3"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.1.104"
  }
}

resource "openstack_networking_port_v2" "port_4" {
  name = "port_4"
  network_id = "${openstack_networking_network_v2.network_2.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_2.id}"
    ip_address = "192.168.2.104"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  depends_on = [
    "openstack_networking_subnet_v2.subnet_1",
    "openstack_networking_subnet_v2.subnet_2",
    "openstack_networking_port_v2.port_1",
    "openstack_networking_port_v2.port_2",
  ]

  name = "instance_1"
  security_groups = ["default"]

  network {
    uuid = "%s"
  }

  network {
    uuid = "${openstack_networking_network_v2.network_1.id}"
    fixed_ip_v4 = "192.168.1.100"
  }

  network {
    uuid = "${openstack_networking_network_v2.network_2.id}"
    fixed_ip_v4 = "192.168.2.100"
  }

  network {
    uuid = "${openstack_networking_network_v2.network_1.id}"
    fixed_ip_v4 = "192.168.1.101"
  }

  network {
    uuid = "${openstack_networking_network_v2.network_2.id}"
    fixed_ip_v4 = "192.168.2.101"
  }

  network {
    port = "${openstack_networking_port_v2.port_1.id}"
  }

  network {
    port = "${openstack_networking_port_v2.port_2.id}"
  }

  network {
    port = "${openstack_networking_port_v2.port_3.id}"
  }

  network {
    port = "${openstack_networking_port_v2.port_4.id}"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceStateActive() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  power_state = "active"
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceStateShutoff() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  power_state = "shutoff"
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceStateShelve() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  power_state = "shelved_offloaded"
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceTagsCreate() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  tags = ["tag1", "tag2", "tag3"]
}
`, osNetworkID)
}

func testAccComputeV2InstanceTagsAdd() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  tags = ["tag1", "tag2", "tag3", "tag4"]
}
`, osNetworkID)
}

func testAccComputeV2InstanceTagsDelete() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
  tags = ["tag2", "tag3"]
}
`, osNetworkID)
}

func testAccComputeV2InstanceTagsClear() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}
