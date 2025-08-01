// This set of code handles all functions required to configure networking
// on an openstack_compute_instance_v2 resource.
//
// This is a complicated task because it's not possible to obtain all
// information in a single API call. In fact, it even traverses multiple
// OpenStack services.
//
// The end result, from the user's point of view, is a structured set of
// understandable network information within the instance resource.
package openstack

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	computeV2InstanceCreateServerWithNetworkModeMicroversion        = "2.37"
	computeV2InstanceCreateServerWithTagsMicroversion               = "2.52"
	computeV2InstanceCreateServerWithHypervisorHostnameMicroversion = "2.74"
	computeV2InstanceCreateServerWithHostnameMicroversion           = "2.90"
	computeV2InstanceCreateServerWithHostnameIsFqdnMicroversion     = "2.94"
	computeV2TagsExtensionMicroversion                              = "2.26"
	computeV2InstanceBlockDeviceVolumeTypeMicroversion              = "2.67"
	computeV2InstanceBlockDeviceVolumeAttachTagsMicroversion        = "2.49"
	computeV2InstanceBlockDeviceMultiattachMicroversion             = "2.60"
)

// InstanceNIC is a structured representation of a Gophercloud servers.Server
// virtual NIC.
type InstanceNIC struct {
	FixedIPv4 string
	FixedIPv6 string
	MAC       string
}

// InstanceAddresses is a collection of InstanceNICs, grouped by the
// network name. An instance/server could have multiple NICs on the same
// network.
type InstanceAddresses struct {
	NetworkName  string
	InstanceNICs []InstanceNIC
}

// InstanceNetwork represents a collection of network information that a
// Terraform instance needs to satisfy all network information requirements.
type InstanceNetwork struct {
	UUID          string
	Name          string
	Port          string
	FixedIP       string
	AccessNetwork bool
}

// getAllInstanceNetworks loops through the networks defined in the Terraform
// configuration and structures that information into something standard that
// can be consumed by both OpenStack and Terraform.
//
// This would be simple, except we have ensure both the network name and
// network ID have been determined. This isn't just for the convenience of a
// user specifying a human-readable network name, but the network information
// returned by an OpenStack instance only has the network name set! So if a
// user specified a network ID, there's no way to correlate it to the instance
// unless we know both the name and ID.
//
// Not only that, but we have to account for two OpenStack network services
// running: nova-network (legacy) and Neutron (current).
//
// In addition, if a port was specified, not all of the port information
// will be displayed, such as multiple fixed and floating IPs. This resource
// isn't currently configured for that type of flexibility. It's better to
// reference the actual port resource itself.
//
// So, let's begin the journey.
func getAllInstanceNetworks(ctx context.Context, d *schema.ResourceData, meta any) ([]InstanceNetwork, error) {
	networks := d.Get("network").([]any)

	instanceNetworks := make([]InstanceNetwork, 0, len(networks))

	for _, v := range networks {
		network := v.(map[string]any)
		networkID := network["uuid"].(string)
		networkName := network["name"].(string)
		portID := network["port"].(string)

		if networkID == "" && networkName == "" && portID == "" {
			return nil, errors.New("At least one of network.uuid, network.name, or network.port must be set")
		}

		// If a user specified both an ID and name, that makes things easy
		// since both name and ID are already satisfied. No need to query
		// further.
		if networkID != "" && networkName != "" {
			v := InstanceNetwork{
				UUID:          networkID,
				Name:          networkName,
				Port:          portID,
				FixedIP:       network["fixed_ip_v4"].(string),
				AccessNetwork: network["access_network"].(bool),
			}
			instanceNetworks = append(instanceNetworks, v)

			continue
		}

		// But if at least one of name or ID was missing, we have to query
		// for that other piece.
		//
		// Priority is given to a port since a network ID or name usually isn't
		// specified when using a port.
		//
		// Next priority is given to the network ID since it's guaranteed to be
		// an exact match.
		queryType := "name"
		queryTerm := networkName

		if networkID != "" {
			queryType = "id"
			queryTerm = networkID
		}

		if portID != "" {
			queryType = "port"
			queryTerm = portID
		}

		networkInfo, err := getInstanceNetworkInfo(ctx, d, meta, queryType, queryTerm)
		if err != nil {
			return nil, err
		}

		v := InstanceNetwork{
			Port:          portID,
			FixedIP:       network["fixed_ip_v4"].(string),
			AccessNetwork: network["access_network"].(bool),
		}
		if networkInfo["uuid"] != nil {
			v.UUID = networkInfo["uuid"].(string)
		}

		if networkInfo["name"] != nil {
			v.Name = networkInfo["name"].(string)
		}

		instanceNetworks = append(instanceNetworks, v)
	}

	log.Printf("[DEBUG] getAllInstanceNetworks: %#v", instanceNetworks)

	return instanceNetworks, nil
}

// getInstanceNetworkInfo will query for network information in order to make
// an accurate determination of a network's name and a network's ID.
func getInstanceNetworkInfo(ctx context.Context, d *schema.ResourceData, meta any, queryType, queryTerm string) (map[string]any, error) {
	config := meta.(*Config)

	networkClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack networking client: %w", err)
	}

	networkInfo, err := getInstanceNetworkInfoNeutron(ctx, networkClient, queryType, queryTerm)
	if err != nil {
		return nil, fmt.Errorf("Error trying to get network information from the Network API: %w", err)
	}

	return networkInfo, nil
}

// getInstanceNetworkInfoNeutron will query the neutron API for the network
// information.
func getInstanceNetworkInfoNeutron(ctx context.Context, client *gophercloud.ServiceClient, queryType, queryTerm string) (map[string]any, error) {
	// If a port was specified, use it to look up the network ID
	// and then query the network as if a network ID was originally used.
	if queryType == "port" {
		listOpts := ports.ListOpts{
			ID: queryTerm,
		}

		allPages, err := ports.List(client, listOpts).AllPages(ctx)
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve networks from the Network API: %w", err)
		}

		allPorts, err := ports.ExtractPorts(allPages)
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve networks from the Network API: %w", err)
		}

		var port ports.Port

		switch len(allPorts) {
		case 0:
			return nil, fmt.Errorf("Could not find any matching port for %s %s", queryType, queryTerm)
		case 1:
			port = allPorts[0]
		default:
			return nil, fmt.Errorf("More than one port found for %s %s", queryType, queryTerm)
		}

		queryType = "id"
		queryTerm = port.NetworkID
	}

	listOpts := networks.ListOpts{
		Status: "ACTIVE",
	}

	switch queryType {
	case "name":
		listOpts.Name = queryTerm
	default:
		listOpts.ID = queryTerm
	}

	allPages, err := networks.List(client, listOpts).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve networks from the Network API: %w", err)
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve networks from the Network API: %w", err)
	}

	var network networks.Network

	switch len(allNetworks) {
	case 0:
		return nil, fmt.Errorf("Could not find any matching network for %s %s", queryType, queryTerm)
	case 1:
		network = allNetworks[0]
	default:
		return nil, fmt.Errorf("More than one network found for %s %s", queryType, queryTerm)
	}

	v := map[string]any{
		"uuid": network.ID,
		"name": network.Name,
	}

	log.Printf("[DEBUG] getInstanceNetworkInfoNeutron: %#v", v)

	return v, nil
}

// getInstanceAddresses parses a Gophercloud server.Server's Address field into
// a structured InstanceAddresses struct.
func getInstanceAddresses(addresses map[string]any) []InstanceAddresses {
	allInstanceAddresses := make([]InstanceAddresses, 0, len(addresses))

	// Addresses includes a list of all IP addresses assigned to the server,
	// keyed by pool. This unfortunately causes problems because addresses are a
	// JSON object which means that they are unordered because gophercloud
	// uses a map[string]any to hold them and maps are unordered.
	// However OpenStack uses OrderedDict to return this data and this
	// provider uses a List type for the networks which is an ordered type,
	// this means that importing resources into terraform will result in
	// networks import out of order which causes terraform to say it has to
	// rebuild the instance. Unfortunately there is no way to ensure the
	// ordering is correct, lastly this issue occurs only when there's
	// "public" and "private" as the only networks so account for that
	// case specially. Create a list of the network names and if it's
	// the special case override it.
	networkNames := make([]string, len(addresses))
	i := 0

	for k := range addresses {
		networkNames[i] = k
		i++
	}

	if len(networkNames) == 2 {
		if networkNames[0] == "private" && networkNames[1] == "public" {
			networkNames[0] = "public"
			networkNames[1] = "private"
		}
	}

	for _, networkName := range networkNames {
		v := addresses[networkName]
		instanceAddresses := InstanceAddresses{
			NetworkName: networkName,
		}

		for _, v := range v.([]any) {
			instanceNIC := InstanceNIC{}

			var exists bool

			v := v.(map[string]any)
			if v, ok := v["OS-EXT-IPS-MAC:mac_addr"].(string); ok {
				instanceNIC.MAC = v
			}

			if v["OS-EXT-IPS:type"] == "fixed" || v["OS-EXT-IPS:type"] == nil {
				switch v["version"].(float64) {
				case 6:
					instanceNIC.FixedIPv6 = v["addr"].(string)
				default:
					instanceNIC.FixedIPv4 = v["addr"].(string)
				}
			}

			// To associate IPv4 and IPv6 on the right NIC,
			// key on the mac address and fill in the blanks.
			for i, v := range instanceAddresses.InstanceNICs {
				if v.MAC == instanceNIC.MAC {
					exists = true

					if instanceNIC.FixedIPv6 != "" {
						instanceAddresses.InstanceNICs[i].FixedIPv6 = instanceNIC.FixedIPv6
					}

					if instanceNIC.FixedIPv4 != "" {
						instanceAddresses.InstanceNICs[i].FixedIPv4 = instanceNIC.FixedIPv4
					}
				}
			}

			if !exists {
				instanceAddresses.InstanceNICs = append(instanceAddresses.InstanceNICs, instanceNIC)
			}
		}

		allInstanceAddresses = append(allInstanceAddresses, instanceAddresses)
	}

	log.Printf("[DEBUG] Addresses: %#v", addresses)
	log.Printf("[DEBUG] allInstanceAddresses: %#v", allInstanceAddresses)

	return allInstanceAddresses
}

// expandInstanceNetworks takes network information found in []InstanceNetwork
// and builds a Gophercloud []servers.Network for use in creating an Instance.
func expandInstanceNetworks(allInstanceNetworks []InstanceNetwork) []servers.Network {
	networks := make([]servers.Network, 0, len(allInstanceNetworks))

	for _, v := range allInstanceNetworks {
		n := servers.Network{
			UUID:    v.UUID,
			Port:    v.Port,
			FixedIP: v.FixedIP,
		}
		networks = append(networks, n)
	}

	return networks
}

// flattenInstanceNetworks collects instance network information from different
// sources and aggregates it all together into a map array.
func flattenInstanceNetworks(ctx context.Context, d *schema.ResourceData, meta any) ([]map[string]any, error) {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack compute client: %w", err)
	}

	server, err := servers.Get(ctx, computeClient, d.Id()).Extract()
	if err != nil {
		return nil, CheckDeleted(d, err, "server")
	}

	allInstanceAddresses := getInstanceAddresses(server.Addresses)

	allInstanceNetworks, err := getAllInstanceNetworks(ctx, d, meta)
	if err != nil {
		return nil, err
	}

	networks := []map[string]any{}

	// If there were no instance networks returned, this means that there
	// was not a network specified in the Terraform configuration. When this
	// happens, the instance will be launched on a "default" network, if one
	// is available. If there isn't, the instance will fail to launch, so
	// this is a safe assumption at this point.
	if len(allInstanceNetworks) == 0 {
		for _, instanceAddresses := range allInstanceAddresses {
			for _, instanceNIC := range instanceAddresses.InstanceNICs {
				v := map[string]any{
					"name":        instanceAddresses.NetworkName,
					"fixed_ip_v4": instanceNIC.FixedIPv4,
					"fixed_ip_v6": instanceNIC.FixedIPv6,
					"mac":         instanceNIC.MAC,
				}

				// Use the same method as getAllInstanceNetworks to get the network uuid
				networkInfo, err := getInstanceNetworkInfo(ctx, d, meta, "name", instanceAddresses.NetworkName)
				if err != nil {
					log.Printf("[WARN] Error getting default network uuid: %s", err)
				} else {
					if v["uuid"] != nil {
						v["uuid"] = networkInfo["uuid"].(string)
					} else {
						log.Printf("[WARN] Could not get default network uuid")
					}
				}

				networks = append(networks, v)
			}
		}

		log.Printf("[DEBUG] flattenInstanceNetworks: %#v", networks)

		return networks, nil
	}

	// Loop through all networks and addresses, merge relevant address details.
	for _, instanceNetwork := range allInstanceNetworks {
		for _, instanceAddresses := range allInstanceAddresses {
			// Skip if instanceAddresses has no NICs
			if len(instanceAddresses.InstanceNICs) == 0 {
				continue
			}

			if instanceNetwork.Name == instanceAddresses.NetworkName {
				// Only use one NIC since it's possible the user defined another NIC
				// on this same network in another Terraform network block.
				instanceNIC := instanceAddresses.InstanceNICs[0]
				copy(instanceAddresses.InstanceNICs, instanceAddresses.InstanceNICs[1:])
				v := map[string]any{
					"name":           instanceAddresses.NetworkName,
					"fixed_ip_v4":    instanceNIC.FixedIPv4,
					"fixed_ip_v6":    instanceNIC.FixedIPv6,
					"mac":            instanceNIC.MAC,
					"uuid":           instanceNetwork.UUID,
					"port":           instanceNetwork.Port,
					"access_network": instanceNetwork.AccessNetwork,
				}
				networks = append(networks, v)
			}
		}
	}

	log.Printf("[DEBUG] flattenInstanceNetworks: %#v", networks)

	return networks, nil
}

// getInstanceAccessAddresses determines the best IP address to communicate
// with the instance. It does this by looping through all networks and looking
// for a valid IP address. Priority is given to a network that was flagged as
// an access_network.
func getInstanceAccessAddresses(networks []map[string]any) (string, string) {
	var hostv4, hostv6 string

	// Loop through all networks
	// If the network has a valid fixed v4 or fixed v6 address
	// and hostv4 or hostv6 is not set, set hostv4/hostv6.
	// If the network is an "access_network" overwrite hostv4/hostv6.
	for _, n := range networks {
		var accessNetwork bool

		if an, ok := n["access_network"].(bool); ok && an {
			accessNetwork = true
		}

		if fixedIPv4, ok := n["fixed_ip_v4"].(string); ok && fixedIPv4 != "" {
			if hostv4 == "" || accessNetwork {
				hostv4 = fixedIPv4
			}
		}

		if fixedIPv6, ok := n["fixed_ip_v6"].(string); ok && fixedIPv6 != "" {
			if hostv6 == "" || accessNetwork {
				hostv6 = fixedIPv6
			}
		}
	}

	log.Printf("[DEBUG] OpenStack Instance Network Access Addresses: %s, %s", hostv4, hostv6)

	return hostv4, hostv6
}

func computeV2InstanceReadTags(d *schema.ResourceData, tags []string) {
	expandObjectReadTags(d, tags)
}

func computeV2InstanceUpdateTags(d *schema.ResourceData) []string {
	return expandObjectUpdateTags(d)
}

func computeV2InstanceTags(d *schema.ResourceData) []string {
	return expandObjectTags(d)
}
