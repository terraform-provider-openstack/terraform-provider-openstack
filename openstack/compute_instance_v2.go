package openstack

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/schedulerhints"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

// ServerV2StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an OpenStack instance.
func ServerV2StateRefreshFunc(client *gophercloud.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := servers.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return s, "DELETED", nil
			}
			return nil, "", err
		}

		return s, s.Status, nil
	}
}

func resourceInstanceSecGroupsV2(d *schema.ResourceData) []string {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secgroups := make([]string, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		secgroups[i] = raw.(string)
	}
	return secgroups
}

func resourceInstanceMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceInstanceBlockDevicesV2(d *schema.ResourceData, bds []interface{}) ([]bootfromvolume.BlockDevice, error) {
	blockDeviceOpts := make([]bootfromvolume.BlockDevice, len(bds))
	for i, bd := range bds {
		bdM := bd.(map[string]interface{})
		blockDeviceOpts[i] = bootfromvolume.BlockDevice{
			UUID:                bdM["uuid"].(string),
			VolumeSize:          bdM["volume_size"].(int),
			BootIndex:           bdM["boot_index"].(int),
			DeleteOnTermination: bdM["delete_on_termination"].(bool),
			GuestFormat:         bdM["guest_format"].(string),
			DeviceType:          bdM["device_type"].(string),
			DiskBus:             bdM["disk_bus"].(string),
		}

		sourceType := bdM["source_type"].(string)
		switch sourceType {
		case "blank":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceBlank
		case "image":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceImage
		case "snapshot":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceSnapshot
		case "volume":
			blockDeviceOpts[i].SourceType = bootfromvolume.SourceVolume
		default:
			return blockDeviceOpts, fmt.Errorf("unknown block device source type %s", sourceType)
		}

		destinationType := bdM["destination_type"].(string)
		switch destinationType {
		case "local":
			blockDeviceOpts[i].DestinationType = bootfromvolume.DestinationLocal
		case "volume":
			blockDeviceOpts[i].DestinationType = bootfromvolume.DestinationVolume
		default:
			return blockDeviceOpts, fmt.Errorf("unknown block device destination type %s", destinationType)
		}
	}

	log.Printf("[DEBUG] Block Device Options: %+v", blockDeviceOpts)
	return blockDeviceOpts, nil
}

func resourceInstanceSchedulerHintsV2(d *schema.ResourceData, schedulerHintsRaw map[string]interface{}) schedulerhints.SchedulerHints {
	differentHost := []string{}
	if len(schedulerHintsRaw["different_host"].([]interface{})) > 0 {
		for _, dh := range schedulerHintsRaw["different_host"].([]interface{}) {
			differentHost = append(differentHost, dh.(string))
		}
	}

	sameHost := []string{}
	if len(schedulerHintsRaw["same_host"].([]interface{})) > 0 {
		for _, sh := range schedulerHintsRaw["same_host"].([]interface{}) {
			sameHost = append(sameHost, sh.(string))
		}
	}

	query := make([]interface{}, len(schedulerHintsRaw["query"].([]interface{})))
	if len(schedulerHintsRaw["query"].([]interface{})) > 0 {
		for _, q := range schedulerHintsRaw["query"].([]interface{}) {
			query = append(query, q.(string))
		}
	}

	schedulerHints := schedulerhints.SchedulerHints{
		Group:                schedulerHintsRaw["group"].(string),
		DifferentHost:        differentHost,
		SameHost:             sameHost,
		Query:                query,
		TargetCell:           schedulerHintsRaw["target_cell"].(string),
		BuildNearHostIP:      schedulerHintsRaw["build_near_host_ip"].(string),
		AdditionalProperties: schedulerHintsRaw["additional_properties"].(map[string]interface{}),
	}

	return schedulerHints
}

func getImageIDFromConfig(computeClient *gophercloud.ServiceClient, d *schema.ResourceData) (string, error) {
	// If block_device was used, an Image does not need to be specified, unless an image/local
	// combination was used. This emulates normal boot behavior. Otherwise, ignore the image altogether.
	if vL, ok := d.GetOk("block_device"); ok {
		needImage := false
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})
			if vM["source_type"] == "image" && vM["destination_type"] == "local" {
				needImage = true
			}
		}
		if !needImage {
			return "", nil
		}
	}

	if imageId := d.Get("image_id").(string); imageId != "" {
		return imageId, nil
	} else {
		// try the OS_IMAGE_ID environment variable
		if v := os.Getenv("OS_IMAGE_ID"); v != "" {
			return v, nil
		}
	}

	imageName := d.Get("image_name").(string)
	if imageName == "" {
		// try the OS_IMAGE_NAME environment variable
		if v := os.Getenv("OS_IMAGE_NAME"); v != "" {
			imageName = v
		}
	}

	if imageName != "" {
		imageId, err := images.IDFromName(computeClient, imageName)
		if err != nil {
			return "", err
		}
		return imageId, nil
	}

	return "", fmt.Errorf("Neither a boot device, image ID, or image name were able to be determined.")
}

func setImageInformation(computeClient *gophercloud.ServiceClient, server *servers.Server, d *schema.ResourceData) error {
	// If block_device was used, an Image does not need to be specified, unless an image/local
	// combination was used. This emulates normal boot behavior. Otherwise, ignore the image altogether.
	if vL, ok := d.GetOk("block_device"); ok {
		needImage := false
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})
			if vM["source_type"] == "image" && vM["destination_type"] == "local" {
				needImage = true
			}
		}
		if !needImage {
			d.Set("image_id", "Attempt to boot from volume - no image supplied")
			return nil
		}
	}

	imageId := server.Image["id"].(string)
	if imageId != "" {
		d.Set("image_id", imageId)
		if image, err := images.Get(computeClient, imageId).Extract(); err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				// If the image name can't be found, set the value to "Image not found".
				// The most likely scenario is that the image no longer exists in the Image Service
				// but the instance still has a record from when it existed.
				d.Set("image_name", "Image not found")
				return nil
			}
			return err
		} else {
			d.Set("image_name", image.Name)
		}
	}

	return nil
}

func getFlavorID(computeClient *gophercloud.ServiceClient, d *schema.ResourceData) (string, error) {
	if flavorId := d.Get("flavor_id").(string); flavorId != "" {
		return flavorId, nil
	} else {
		// Try the OS_FLAVOR_ID environment variable
		if v := os.Getenv("OS_FLAVOR_ID"); v != "" {
			return v, nil
		}
	}

	flavorName := d.Get("flavor_name").(string)
	if flavorName == "" {
		// Try the OS_FLAVOR_NAME environment variable
		if v := os.Getenv("OS_FLAVOR_NAME"); v != "" {
			flavorName = v
		}
	}

	if flavorName != "" {
		flavorId, err := flavors.IDFromName(computeClient, flavorName)
		if err != nil {
			return "", err
		}
		return flavorId, nil
	}

	return "", fmt.Errorf("Neither a flavor_id or flavor_name could be determined.")
}

func resourceComputeSchedulerHintsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["group"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["group"].(string)))
	}

	if m["target_cell"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["target_cell"].(string)))
	}

	if m["build_host_near_ip"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["build_host_near_ip"].(string)))
	}

	if m["additional_properties"] != nil {
		for _, v := range m["additional_properties"].(map[string]interface{}) {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}

	buf.WriteString(fmt.Sprintf("%s-", m["different_host"].([]interface{})))
	buf.WriteString(fmt.Sprintf("%s-", m["same_host"].([]interface{})))
	buf.WriteString(fmt.Sprintf("%s-", m["query"].([]interface{})))

	return hashcode.String(buf.String())
}

func checkBlockDeviceConfig(d *schema.ResourceData) error {
	if vL, ok := d.GetOk("block_device"); ok {
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})

			if vM["source_type"] != "blank" && vM["uuid"] == "" {
				return fmt.Errorf("You must specify a uuid for %s block device types", vM["source_type"])
			}

			if vM["source_type"] == "image" && vM["destination_type"] == "volume" {
				if vM["volume_size"] == 0 {
					return fmt.Errorf("You must specify a volume_size when creating a volume from an image")
				}
			}

			if vM["source_type"] == "blank" && vM["destination_type"] == "local" {
				if vM["volume_size"] == 0 {
					return fmt.Errorf("You must specify a volume_size when creating a blank block device")
				}
			}
		}
	}

	return nil
}

func resourceComputeInstancePersonalityHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["file"].(string)))

	return hashcode.String(buf.String())
}

func resourceInstancePersonalityV2(d *schema.ResourceData) servers.Personality {
	var personalities servers.Personality

	if v := d.Get("personality"); v != nil {
		personalityList := v.(*schema.Set).List()
		if len(personalityList) > 0 {
			for _, p := range personalityList {
				rawPersonality := p.(map[string]interface{})
				file := servers.File{
					Path:     rawPersonality["file"].(string),
					Contents: []byte(rawPersonality["content"].(string)),
				}

				log.Printf("[DEBUG] OpenStack Compute Instance Personality: %+v", file)

				personalities = append(personalities, &file)
			}
		}
	}

	return personalities
}

// suppressAvailabilityZoneDetailDiffs will suppress diffs when a user specifies an
// availability zone in the format of `az:host:node` and Nova/Compute responds with
// only `az`.
func suppressAvailabilityZoneDetailDiffs(k, old, new string, d *schema.ResourceData) bool {
	newParts := strings.Split(new, ":")
	newAZ := newParts[0]
	oldParts := strings.Split(old, ":")
	oldAZ := oldParts[0]

	if oldAZ == newAZ {
		return true
	}

	if len(oldParts) > 1 && oldAZ == "" {
		return true
	}

	return false
}

// suppressPowerStateDiffs will allow a state of "error" or "migrating" even though we don't
// allow them as a user input.
func suppressPowerStateDiffs(k, old, new string, d *schema.ResourceData) bool {
	if old == "error" || old == "migrating" {
		return true
	}

	return false
}
