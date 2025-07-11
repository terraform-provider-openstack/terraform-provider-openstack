package openstack

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/secgroups"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/tags"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	flavorsutils "github.com/gophercloud/utils/v2/openstack/compute/v2/flavors"
	imagesutils "github.com/gophercloud/utils/v2/openstack/image/v2/images"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-provider-openstack/utils/v2/hashcode"
)

func resourceComputeInstanceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeInstanceV2Create,
		ReadContext:   resourceComputeInstanceV2Read,
		UpdateContext: resourceComputeInstanceV2Update,
		DeleteContext: resourceComputeInstanceV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceOpenStackComputeInstanceV2ImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"hostname": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Computed:     true,
				ValidateFunc: validateHostname(),
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// just stash the hash for state & diff comparisons
				StateFunc: func(v any) string {
					switch v := v.(type) {
					case string:
						hash := sha1.Sum([]byte(v))

						return hex.EncodeToString(hash[:])
					default:
						return ""
					}
				},
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"availability_zone_hints": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"availability_zone"},
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ConflictsWith:    []string{"availability_zone_hints"},
				DiffSuppressFunc: suppressAvailabilityZoneDetailDiffs,
			},
			"network_mode": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      false,
				ConflictsWith: []string{"network"},
				ValidateFunc: validation.StringInSlice([]string{
					"auto", "none",
				}, true),
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"fixed_ip_v6": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_network": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"hypervisor_hostname": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"personality"},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"config_drive": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"admin_pass": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  false,
			},
			"access_ip_v4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_ip_v6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"block_device": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"destination_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"boot_index": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"guest_format": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"device_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_bus": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"multiattach": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
					},
				},
			},
			"scheduler_hints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"different_host": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"same_host": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"query": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"target_cell": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"different_cell": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"build_near_host_ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"additional_properties": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				Set: resourceComputeSchedulerHintsHash,
			},
			"personality": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:           resourceComputeInstancePersonalityHash,
				ConflictsWith: []string{"hypervisor_hostname"},
			},
			"stop_before_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"all_metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"power_state": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{
					"active", "shutoff", "shelved_offloaded", "paused",
				}, true),
				DiffSuppressFunc: suppressPowerStateDiffs,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vendor_options": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ignore_resize_confirmation": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"detach_ports_before_destroy": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
					},
				},
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CustomizeDiff: customdiff.All(
			// OpenStack cannot resize an instance, if its original flavor is deleted, that is why
			// we need to force recreation, if old flavor name or ID is reported as an empty string
			customdiff.ForceNewIfChange("flavor_id", func(_ context.Context, old, _, _ any) bool {
				return old.(string) == ""
			}),
			customdiff.ForceNewIfChange("flavor_name", func(_ context.Context, old, _, _ any) bool {
				return old.(string) == ""
			}),
			func(_ context.Context, d *schema.ResourceDiff, _ any) error {
				currentState, _ := d.GetChange("power_state")
				if currentState == "build" {
					// In "build" state, network and security groups are not yet available
					if err := d.Clear("network"); err != nil {
						return err
					}
					if err := d.Clear("security_groups"); err != nil {
						return err
					}
				}

				return nil
			},
		),
	}
}

func resourceComputeInstanceV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	imageClient, err := config.ImageV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	var availabilityZone string

	var networks any

	// Determines the Image ID using the following rules:
	// If a bootable block_device was specified, ignore the image altogether.
	// If an image_id was specified, use it.
	// If an image_name was specified, look up the image ID, report if error.
	imageID, err := getImageIDFromConfig(ctx, imageClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Determines the Flavor ID using the following rules:
	// If a flavor_id was specified, use it.
	// If a flavor_name was specified, lookup the flavor ID, report if error.
	flavorID, err := getFlavorID(ctx, computeClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// determine if block_device configuration is correct
	// this includes valid combinations and required attributes
	if err := checkBlockDeviceConfig(d); err != nil {
		return diag.FromErr(err)
	}

	if networkMode := d.Get("network_mode").(string); networkMode == "auto" || networkMode == "none" {
		// Use special string for network option
		computeClient.Microversion = computeV2InstanceCreateServerWithNetworkModeMicroversion
		networks = networkMode
		log.Printf("[DEBUG] Create with network options %s", networks)
	} else {
		log.Printf("[DEBUG] Create with specified network options")
		// Build a list of networks with the information given upon creation.
		// Error out if an invalid network configuration was used.
		allInstanceNetworks, err := getAllInstanceNetworks(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		// Build a []servers.Network to pass into the create options.
		networks = expandInstanceNetworks(allInstanceNetworks)
	}

	configDrive := d.Get("config_drive").(bool)

	// Retrieve tags and set microversion if they're provided.
	instanceTags := computeV2InstanceTags(d)
	if len(instanceTags) > 0 {
		computeClient.Microversion = computeV2InstanceCreateServerWithTagsMicroversion
	}

	var hypervisorHostname string
	if v, ok := getOkExists(d, "hypervisor_hostname"); ok {
		hypervisorHostname = v.(string)
		computeClient.Microversion = computeV2InstanceCreateServerWithHypervisorHostnameMicroversion
	}

	var hostname string
	if v, ok := getOkExists(d, "hostname"); ok {
		hostname = v.(string)
		if isValidHostname(hostname) {
			computeClient.Microversion = computeV2InstanceCreateServerWithHostnameMicroversion
		} else {
			computeClient.Microversion = computeV2InstanceCreateServerWithHostnameIsFqdnMicroversion
		}
	}

	if v, ok := getOkExists(d, "availability_zone"); ok {
		availabilityZone = v.(string)
	} else {
		availabilityZone = d.Get("availability_zone_hints").(string)
	}

	createOpts := &servers.CreateOpts{
		Name:               d.Get("name").(string),
		Hostname:           hostname,
		ImageRef:           imageID,
		FlavorRef:          flavorID,
		SecurityGroups:     resourceInstanceSecGroupsV2(d),
		AvailabilityZone:   availabilityZone,
		Networks:           networks,
		HypervisorHostname: hypervisorHostname,
		Metadata:           resourceInstanceMetadataV2(d),
		ConfigDrive:        &configDrive,
		AdminPass:          d.Get("admin_pass").(string),
		UserData:           []byte(d.Get("user_data").(string)),
		Personality:        resourceInstancePersonalityV2(d),
		Tags:               instanceTags,
	}

	if vL, ok := d.GetOk("block_device"); ok {
		blockDevices, err := resourceInstanceBlockDevicesV2(d, vL.([]any))
		if err != nil {
			return diag.FromErr(err)
		}

		// Check if Multiattach was set in any of the Block Devices.
		// If so, set the client's microversion appropriately.
		for _, bd := range d.Get("block_device").([]any) {
			if bd.(map[string]any)["multiattach"].(bool) {
				computeClient.Microversion = computeV2InstanceBlockDeviceMultiattachMicroversion
			}
		}

		// Check if VolumeType was set in any of the Block Devices.
		// If so, set the client's microversion appropriately.
		for _, bd := range blockDevices {
			if bd.VolumeType != "" {
				computeClient.Microversion = computeV2InstanceBlockDeviceVolumeTypeMicroversion
			}
		}

		createOpts.BlockDevice = blockDevices
	}

	var createOptsBuilder servers.CreateOptsBuilder = createOpts
	if keyName, ok := d.Get("key_pair").(string); ok && keyName != "" {
		createOptsBuilder = &keypairs.CreateOptsExt{
			CreateOptsBuilder: createOptsBuilder,
			KeyName:           keyName,
		}
	}

	var schedulerHints servers.SchedulerHintOpts

	schedulerHintsRaw := d.Get("scheduler_hints").(*schema.Set).List()
	if len(schedulerHintsRaw) > 0 {
		log.Printf("[DEBUG] schedulerhints: %+v", schedulerHintsRaw)
		schedulerHints = resourceInstanceSchedulerHintsV2(schedulerHintsRaw[0].(map[string]any))
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// If a block_device is used, use the bootfromvolume.Create function as it allows an empty ImageRef.
	// Otherwise, use the normal servers.Create function.
	server, err := servers.Create(ctx, computeClient, createOptsBuilder, schedulerHints).Extract()
	if err != nil {
		return diag.Errorf("Error creating OpenStack server: %s", err)
	}

	log.Printf("[INFO] Instance ID: %s", server.ID)

	// Store the ID now
	d.SetId(server.ID)

	// Wait for the instance to become running so we can get some attributes
	// that aren't available until later.
	log.Printf(
		"[DEBUG] Waiting for instance (%s) to become running",
		server.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, server.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	err = retry.RetryContext(ctx, stateConf.Timeout, func() *retry.RetryError {
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			log.Printf("[DEBUG] Retrying after error: %s", err)

			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			server.ID, err)
	}

	vmState := d.Get("power_state").(string)
	if strings.ToLower(vmState) == "shutoff" {
		err = servers.Stop(ctx, computeClient, d.Id()).ExtractErr()
		if err != nil {
			return diag.Errorf("Error stopping OpenStack instance: %s", err)
		}

		stopStateConf := &retry.StateChangeConf{
			// Pending:    []string{"ACTIVE"},
			Target:     []string{"SHUTOFF"},
			Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 3 * time.Second,
		}

		log.Printf("[DEBUG] Waiting for instance (%s) to stop", d.Id())

		_, err = stopStateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for instance (%s) to become inactive(shutoff): %s", d.Id(), err)
		}
	}

	return resourceComputeInstanceV2Read(ctx, d, meta)
}

func resourceComputeInstanceV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	imageClient, err := config.ImageV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack image client: %s", err)
	}

	server, err := servers.Get(ctx, computeClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "server"))
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	d.Set("name", server.Name)
	d.Set("created", server.Created.String())
	d.Set("updated", server.Updated.String())

	// Get the instance network and address information
	networks, err := flattenInstanceNetworks(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// Determine the best IPv4 and IPv6 addresses to access the instance with
	hostv4, hostv6 := getInstanceAccessAddresses(networks)

	// AccessIPv4/v6 isn't standard in OpenStack, but there have been reports
	// of them being used in some environments.
	if server.AccessIPv4 != "" && hostv4 == "" {
		hostv4 = server.AccessIPv4
	}

	if server.AccessIPv6 != "" && hostv6 == "" {
		hostv6 = server.AccessIPv6
	}

	d.Set("network", networks)
	d.Set("access_ip_v4", hostv4)
	d.Set("access_ip_v6", hostv6)

	// Determine the best IP address to use for SSH connectivity.
	// Prefer IPv4 over IPv6.
	var preferredSSHAddress string
	if hostv4 != "" {
		preferredSSHAddress = hostv4
	} else if hostv6 != "" {
		preferredSSHAddress = hostv6
	}

	if preferredSSHAddress != "" {
		// Initialize the connection info
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": preferredSSHAddress,
		})
	}

	d.Set("all_metadata", server.Metadata)

	secGrpNames := []string{}
	for _, sg := range server.SecurityGroups {
		secGrpNames = append(secGrpNames, sg["name"].(string))
	}

	d.Set("security_groups", secGrpNames)

	d.Set("key_pair", server.KeyName)

	flavorID, ok := server.Flavor["id"].(string)
	if !ok {
		return diag.Errorf("Error setting OpenStack server's flavor: %v", server.Flavor)
	}

	d.Set("flavor_id", flavorID)

	flavor, err := flavors.Get(ctx, computeClient, flavorID).Extract()
	if err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			// Original flavor was deleted, but it is possible that instance started
			// with this flavor is still running
			log.Printf("[DEBUG] Original instance flavor id %s could not be found", d.Id())
			d.Set("flavor_id", "")
			d.Set("flavor_name", "")
		} else {
			return diag.FromErr(err)
		}
	} else {
		d.Set("flavor_name", flavor.Name)
	}

	// Set the instance's image information appropriately
	if err := setImageInformation(ctx, imageClient, server, d); err != nil {
		return diag.FromErr(err)
	}

	// Set the availability zone
	d.Set("availability_zone", server.AvailabilityZone)

	// Set the region
	d.Set("region", GetRegion(d, config))

	// Set the current power_state
	currentStatus := strings.ToLower(server.Status)
	switch currentStatus {
	case "active", "shutoff", "error", "migrating", "shelved_offloaded", "shelved", "build", "paused":
		d.Set("power_state", currentStatus)
	default:
		return diag.Errorf("Invalid power_state for instance %s: %s", d.Id(), server.Status)
	}

	// Populate tags.
	computeClient.Microversion = computeV2TagsExtensionMicroversion

	instanceTags, err := tags.List(ctx, computeClient, server.ID).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get tags for openstack_compute_instance_v2: %s", err)
	} else {
		computeV2InstanceReadTags(d, instanceTags)
	}

	// Set the hypervisor hostname
	d.Set("hypervisor_hostname", server.HypervisorHostname)

	return nil
}

func resourceComputeInstanceV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	var updateOpts servers.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("hostname") {
		hostname := d.Get("hostname").(string)
		updateOpts.Hostname = &hostname

		// Set the required microversion.
		if isValidHostname(*updateOpts.Hostname) {
			computeClient.Microversion = computeV2InstanceCreateServerWithHostnameMicroversion
		} else {
			computeClient.Microversion = computeV2InstanceCreateServerWithHostnameIsFqdnMicroversion
		}
	}

	if updateOpts != (servers.UpdateOpts{}) {
		_, err := servers.Update(ctx, computeClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating OpenStack server: %s", err)
		}
		// Reset microversion.
		computeClient.Microversion = ""
	}

	if d.HasChange("power_state") {
		powerStateOldRaw, powerStateNewRaw := d.GetChange("power_state")
		powerStateOld := powerStateOldRaw.(string)

		powerStateNew := powerStateNewRaw.(string)
		if strings.ToLower(powerStateNew) == "shelved_offloaded" {
			err = servers.Shelve(ctx, computeClient, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("Error shelve OpenStack instance: %s", err)
			}

			shelveStateConf := &retry.StateChangeConf{
				// Pending:    []string{"ACTIVE"},
				Target:     []string{"SHELVED_OFFLOADED"},
				Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      0,
				MinTimeout: 3 * time.Second,
			}

			log.Printf("[DEBUG] Waiting for instance (%s) to shelve", d.Id())

			_, err = shelveStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to become shelve: %s", d.Id(), err)
			}
		}

		if strings.ToLower(powerStateNew) == "paused" {
			err = servers.Pause(ctx, computeClient, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("Error pausing OpenStack instance: %s", err)
			}

			pauseStateConf := &retry.StateChangeConf{
				// Pending:    []string{"ACTIVE"},
				Target:     []string{"PAUSED"},
				Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      0,
				MinTimeout: 3 * time.Second,
			}

			log.Printf("[DEBUG] Waiting for instance (%s) to pause", d.Id())

			_, err = pauseStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to become paused: %s", d.Id(), err)
			}
		}

		if strings.ToLower(powerStateNew) == "shutoff" {
			err = servers.Stop(ctx, computeClient, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("Error stopping OpenStack instance: %s", err)
			}

			stopStateConf := &retry.StateChangeConf{
				// Pending:    []string{"ACTIVE"},
				Target:     []string{"SHUTOFF"},
				Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      0,
				MinTimeout: 3 * time.Second,
			}

			log.Printf("[DEBUG] Waiting for instance (%s) to stop", d.Id())

			_, err = stopStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to become inactive(shutoff): %s", d.Id(), err)
			}
		}

		if strings.ToLower(powerStateNew) == "active" {
			if strings.ToLower(powerStateOld) == "shelved" || strings.ToLower(powerStateOld) == "shelved_offloaded" {
				unshelveOpt := &servers.UnshelveOpts{
					AvailabilityZone: d.Get("availability_zone").(string),
				}

				err = servers.Unshelve(ctx, computeClient, d.Id(), unshelveOpt).ExtractErr()
				if err != nil {
					return diag.Errorf("Error unshelving OpenStack instance: %s", err)
				}
			} else if strings.ToLower(powerStateOld) == "paused" {
				err = servers.Unpause(ctx, computeClient, d.Id()).ExtractErr()
				if err != nil {
					return diag.Errorf("Error resuming OpenStack instance: %s", err)
				}
			} else if strings.ToLower(powerStateOld) != "build" {
				err = servers.Start(ctx, computeClient, d.Id()).ExtractErr()
				if err != nil {
					return diag.Errorf("Error starting OpenStack instance: %s", err)
				}
			}

			startStateConf := &retry.StateChangeConf{
				// Pending:    []string{"SHUTOFF"},
				Target:     []string{"ACTIVE"},
				Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      0,
				MinTimeout: 3 * time.Second,
			}

			log.Printf("[DEBUG] Waiting for instance (%s) to start/unshelve/resume", d.Id())

			_, err = startStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to become active: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("metadata") {
		oldMetadata, newMetadata := d.GetChange("metadata")

		var metadataToDelete []string

		// Determine if any metadata keys were removed from the configuration.
		// Then request those keys to be deleted.
		for oldKey := range oldMetadata.(map[string]any) {
			var found bool

			for newKey := range newMetadata.(map[string]any) {
				if oldKey == newKey {
					found = true
				}
			}

			if !found {
				metadataToDelete = append(metadataToDelete, oldKey)
			}
		}

		for _, key := range metadataToDelete {
			err := servers.DeleteMetadatum(ctx, computeClient, d.Id(), key).ExtractErr()
			if err != nil && CheckDeleted(d, err, "") != nil {
				return diag.Errorf("Error deleting metadata (%s) from server (%s): %s", key, d.Id(), err)
			}
		}

		// Update existing metadata and add any new metadata.
		metadataOpts := make(servers.MetadataOpts)
		for k, v := range newMetadata.(map[string]any) {
			metadataOpts[k] = v.(string)
		}

		_, err := servers.UpdateMetadata(ctx, computeClient, d.Id(), metadataOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating OpenStack server (%s) metadata: %s", d.Id(), err)
		}
	}

	if d.HasChange("security_groups") {
		oldSGRaw, newSGRaw := d.GetChange("security_groups")
		oldSGSet := oldSGRaw.(*schema.Set)
		newSGSet := newSGRaw.(*schema.Set)
		secgroupsToAdd := newSGSet.Difference(oldSGSet)
		secgroupsToRemove := oldSGSet.Difference(newSGSet)

		log.Printf("[DEBUG] Security groups to add: %v", secgroupsToAdd)

		log.Printf("[DEBUG] Security groups to remove: %v", secgroupsToRemove)

		for _, g := range secgroupsToRemove.List() {
			err := secgroups.RemoveServer(ctx, computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
					continue
				}

				return diag.Errorf("Error removing security group (%s) from OpenStack server (%s): %s", g, d.Id(), err)
			}

			log.Printf("[DEBUG] Removed security group (%s) from instance (%s)", g, d.Id())
		}

		for _, g := range secgroupsToAdd.List() {
			err := secgroups.AddServer(ctx, computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				return diag.Errorf("Error adding security group (%s) to OpenStack server (%s): %s", g, d.Id(), err)
			}

			log.Printf("[DEBUG] Added security group (%s) to instance (%s)", g, d.Id())
		}
	}

	if d.HasChange("admin_pass") {
		if newPwd, ok := d.Get("admin_pass").(string); ok {
			err := servers.ChangeAdminPassword(ctx, computeClient, d.Id(), newPwd).ExtractErr()
			if err != nil {
				return diag.Errorf("Error changing admin password of OpenStack server (%s): %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("flavor_id") || d.HasChange("flavor_name") {
		// Get vendor_options
		vendorOptionsRaw := d.Get("vendor_options").(*schema.Set)

		var ignoreResizeConfirmation bool

		if vendorOptionsRaw.Len() > 0 {
			vendorOptions := expandVendorOptions(vendorOptionsRaw.List())
			ignoreResizeConfirmation = vendorOptions["ignore_resize_confirmation"].(bool)
		}

		var newFlavorID string

		var err error

		if d.HasChange("flavor_id") {
			newFlavorID = d.Get("flavor_id").(string)
		} else {
			newFlavorName := d.Get("flavor_name").(string)

			newFlavorID, err = flavorsutils.IDFromName(ctx, computeClient, newFlavorName)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		resizeOpts := &servers.ResizeOpts{
			FlavorRef: newFlavorID,
		}
		log.Printf("[DEBUG] Resize configuration: %#v", resizeOpts)

		err = servers.Resize(ctx, computeClient, d.Id(), resizeOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("Error resizing OpenStack server: %s", err)
		}

		// Wait for the instance to finish resizing.
		log.Printf("[DEBUG] Waiting for instance (%s) to finish resizing", d.Id())

		// Resize instance without confirmation if specified by user.
		if ignoreResizeConfirmation {
			stateConf := &retry.StateChangeConf{
				Pending:    []string{"RESIZE", "VERIFY_RESIZE"},
				Target:     []string{"ACTIVE", "SHUTOFF"},
				Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      0,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to resize: %s", d.Id(), err)
			}
		} else {
			stateConf := &retry.StateChangeConf{
				Pending:    []string{"RESIZE"},
				Target:     []string{"VERIFY_RESIZE"},
				Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      0,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to resize: %s", d.Id(), err)
			}

			// Confirm resize.
			log.Printf("[DEBUG] Confirming resize")

			err = servers.ConfirmResize(ctx, computeClient, d.Id()).ExtractErr()
			if err != nil {
				return diag.Errorf("Error confirming resize of OpenStack server: %s", err)
			}

			stateConf = &retry.StateChangeConf{
				Pending:    []string{"VERIFY_RESIZE"},
				Target:     []string{"ACTIVE", "SHUTOFF"},
				Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      0,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for instance (%s) to confirm resize: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("image_id") || d.HasChange("image_name") || d.HasChange("personality") {
		var newImageID string

		imageClient, err := config.ImageV2Client(ctx, GetRegion(d, config))
		if err != nil {
			return diag.Errorf("Error creating OpenStack image client: %s", err)
		}

		if d.HasChange("image_id") {
			newImageID = d.Get("image_id").(string)
		} else if d.HasChange("image_name") {
			newImageName := d.Get("image_name").(string)

			newImageID, err = imagesutils.IDFromName(ctx, computeClient, newImageName)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			newImageID, err = getImageIDFromConfig(ctx, imageClient, d)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		var rebuildOpts servers.RebuildOptsBuilder = &servers.RebuildOpts{
			ImageRef:    newImageID,
			Personality: resourceInstancePersonalityV2(d),
		}

		log.Printf("[DEBUG] Rebuild configuration: %#v", rebuildOpts)

		_, err = servers.Rebuild(ctx, computeClient, d.Id(), rebuildOpts).Extract()
		if err != nil {
			return diag.Errorf("Error rebuilding OpenStack server: %s", err)
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{"REBUILD"},
			Target:     []string{"ACTIVE", "SHUTOFF"},
			Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      0,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for instance (%s) to rebuild: %s", d.Id(), err)
		}
	}

	// Perform any required updates to the tags.
	if d.HasChange("tags") {
		instanceTags := computeV2InstanceUpdateTags(d)
		instanceTagsOpts := tags.ReplaceAllOpts{Tags: instanceTags}
		computeClient.Microversion = computeV2TagsExtensionMicroversion

		instanceTags, err := tags.ReplaceAll(ctx, computeClient, d.Id(), instanceTagsOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_compute_instance_v2 %s: %s", d.Id(), err)
		}

		log.Printf("[DEBUG] Set tags %s on openstack_compute_instance_v2 %s", instanceTags, d.Id())
	}

	return resourceComputeInstanceV2Read(ctx, d, meta)
}

func resourceComputeInstanceV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	if d.Get("stop_before_destroy").(bool) {
		err = servers.Stop(ctx, computeClient, d.Id()).ExtractErr()
		if err != nil {
			log.Printf("[WARN] Error stopping openstack_compute_instance_v2: %s", err)
		} else {
			stopStateConf := &retry.StateChangeConf{
				Pending:    []string{"ACTIVE"},
				Target:     []string{"SHUTOFF"},
				Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutDelete),
				Delay:      0,
				MinTimeout: 3 * time.Second,
			}

			log.Printf("[DEBUG] Waiting for instance (%s) to stop", d.Id())

			_, err = stopStateConf.WaitForStateContext(ctx)
			if err != nil {
				log.Printf("[WARN] Error waiting for instance (%s) to stop: %s, proceeding to delete", d.Id(), err)
			}
		}
	}

	vendorOptionsRaw := d.Get("vendor_options").(*schema.Set)

	var detachPortBeforeDestroy bool

	if vendorOptionsRaw.Len() > 0 {
		vendorOptions := expandVendorOptions(vendorOptionsRaw.List())
		detachPortBeforeDestroy = vendorOptions["detach_ports_before_destroy"].(bool)
	}

	if detachPortBeforeDestroy {
		allInstanceNetworks, err := getAllInstanceNetworks(ctx, d, meta)
		if err != nil {
			log.Printf("[WARN] Unable to get openstack_compute_instance_v2 ports: %s", err)
		} else {
			for _, network := range allInstanceNetworks {
				if network.Port != "" {
					stateConf := &retry.StateChangeConf{
						Pending:    []string{""},
						Target:     []string{"DETACHED"},
						Refresh:    computeInterfaceAttachV2DetachFunc(ctx, computeClient, d.Id(), network.Port),
						Timeout:    d.Timeout(schema.TimeoutDelete),
						Delay:      0,
						MinTimeout: 5 * time.Second,
					}
					if _, err = stateConf.WaitForStateContext(ctx); err != nil {
						return diag.Errorf("Error detaching openstack_compute_instance_v2 %s: %s", d.Id(), err)
					}
				}
			}
		}
	}

	if d.Get("force_delete").(bool) {
		log.Printf("[DEBUG] Force deleting OpenStack Instance %s", d.Id())

		err = servers.ForceDelete(ctx, computeClient, d.Id()).ExtractErr()
		if err != nil {
			return diag.FromErr(CheckDeleted(d, err, "Error force deleting openstack_compute_instance_v2"))
		}
	} else {
		log.Printf("[DEBUG] Deleting OpenStack Instance %s", d.Id())

		err = servers.Delete(ctx, computeClient, d.Id()).ExtractErr()
		if err != nil {
			return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_compute_instance_v2"))
		}
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE", "SHUTOFF"},
		Target:     []string{"DELETED", "SOFT_DELETED"},
		Refresh:    ServerV2StateRefreshFunc(ctx, computeClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for instance (%s) to Delete:  %s",
			d.Id(), err)
	}

	return nil
}

func resourceOpenStackComputeInstanceV2ImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	var serverWithAttachments struct {
		VolumesAttached []map[string]any `json:"os-extended-volumes:volumes_attached"`
	}

	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack compute client: %w", err)
	}

	results := make([]*schema.ResourceData, 1)

	diagErr := resourceComputeInstanceV2Read(ctx, d, meta)
	if diagErr != nil {
		return nil, fmt.Errorf("Error reading openstack_compute_instance_v2 %s: %v", d.Id(), diagErr)
	}

	raw := servers.Get(ctx, computeClient, d.Id())
	if raw.Err != nil {
		return nil, CheckDeleted(d, raw.Err, "openstack_compute_instance_v2")
	}

	if err := raw.ExtractInto(&serverWithAttachments); err != nil {
		log.Printf("[DEBUG] unable to unmarshal raw struct to serverWithAttachments: %s", err)
	}

	log.Printf("[DEBUG] Retrieved openstack_compute_instance_v2 %s volume attachments: %#v",
		d.Id(), serverWithAttachments)

	bds := []map[string]any{}

	if len(serverWithAttachments.VolumesAttached) > 0 {
		blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
		if err == nil {
			volMetaData := struct {
				VolumeImageMetadata map[string]any `json:"volume_image_metadata"`
				ID                  string         `json:"id"`
				Size                int            `json:"size"`
				Bootable            string         `json:"bootable"`
			}{}

			for i, b := range serverWithAttachments.VolumesAttached {
				rawVolume := volumes.Get(ctx, blockStorageClient, b["id"].(string))
				if err := rawVolume.ExtractInto(&volMetaData); err != nil {
					log.Printf("[DEBUG] unable to unmarshal raw struct to volume metadata: %s", err)
				}

				log.Printf("[DEBUG] retrieved volume%+v", volMetaData)
				v := map[string]any{
					"delete_on_termination": true,
					"uuid":                  volMetaData.VolumeImageMetadata["image_id"],
					"boot_index":            i,
					"destination_type":      "volume",
					"source_type":           "image",
					"volume_size":           volMetaData.Size,
					"disk_bus":              "",
					"volume_type":           "",
					"device_type":           "",
				}

				if volMetaData.Bootable == "true" {
					bds = append(bds, v)
				}
			}
		} else {
			log.Print("[DEBUG] Could not create BlockStorageV3 client, trying BlockStorageV2")

			blockStorageClient, err := config.BlockStorageV2Client(ctx, GetRegion(d, config))
			if err != nil {
				return nil, fmt.Errorf("Error creating OpenStack volume V2 client: %w", err)
			}

			volMetaData := struct {
				VolumeImageMetadata map[string]any `json:"volume_image_metadata"`
				ID                  string         `json:"id"`
				Size                int            `json:"size"`
				Bootable            string         `json:"bootable"`
			}{}

			for i, b := range serverWithAttachments.VolumesAttached {
				rawVolume := volumes.Get(ctx, blockStorageClient, b["id"].(string))
				if err := rawVolume.ExtractInto(&volMetaData); err != nil {
					log.Printf("[DEBUG] unable to unmarshal raw struct to volume metadata: %s", err)
				}

				log.Printf("[DEBUG] retrieved volume%+v", volMetaData)
				v := map[string]any{
					"delete_on_termination": true,
					"uuid":                  volMetaData.VolumeImageMetadata["image_id"],
					"boot_index":            i,
					"destination_type":      "volume",
					"source_type":           "image",
					"volume_size":           volMetaData.Size,
					"disk_bus":              "",
					"volume_type":           "",
					"device_type":           "",
				}

				if volMetaData.Bootable == "true" {
					bds = append(bds, v)
				}
			}
		}

		d.Set("block_device", bds)
	}

	metadata, err := servers.Metadata(ctx, computeClient, d.Id()).Extract()
	if err != nil {
		return nil, fmt.Errorf("Unable to read metadata for openstack_compute_instance_v2 %s: %w", d.Id(), err)
	}

	d.Set("metadata", metadata)

	results[0] = d

	return results, nil
}

// ServerV2StateRefreshFunc returns a retry.StateRefreshFunc that is used to watch
// an OpenStack instance.
func ServerV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, instanceID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		s, err := servers.Get(ctx, client, instanceID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return s, "DELETED", nil
			}

			return nil, "", err
		}

		return s, s.Status, nil
	}
}

func resourceInstanceSecGroupsV2(d *schema.ResourceData) []string {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	res := make([]string, len(rawSecGroups))

	for i, raw := range rawSecGroups {
		res[i] = raw.(string)
	}

	return res
}

func resourceInstanceMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]any) {
		m[key] = val.(string)
	}

	return m
}

func resourceInstanceBlockDevicesV2(_ *schema.ResourceData, bds []any) ([]servers.BlockDevice, error) {
	blockDeviceOpts := make([]servers.BlockDevice, len(bds))

	for i, bd := range bds {
		bdM := bd.(map[string]any)
		blockDeviceOpts[i] = servers.BlockDevice{
			UUID:                bdM["uuid"].(string),
			VolumeSize:          bdM["volume_size"].(int),
			BootIndex:           bdM["boot_index"].(int),
			DeleteOnTermination: bdM["delete_on_termination"].(bool),
			GuestFormat:         bdM["guest_format"].(string),
			VolumeType:          bdM["volume_type"].(string),
			DeviceType:          bdM["device_type"].(string),
			DiskBus:             bdM["disk_bus"].(string),
		}

		sourceType := bdM["source_type"].(string)
		switch sourceType {
		case "blank":
			blockDeviceOpts[i].SourceType = servers.SourceBlank
		case "image":
			blockDeviceOpts[i].SourceType = servers.SourceImage
		case "snapshot":
			blockDeviceOpts[i].SourceType = servers.SourceSnapshot
		case "volume":
			blockDeviceOpts[i].SourceType = servers.SourceVolume
		default:
			return blockDeviceOpts, fmt.Errorf("unknown block device source type %s", sourceType)
		}

		destinationType := bdM["destination_type"].(string)
		switch destinationType {
		case "local":
			blockDeviceOpts[i].DestinationType = servers.DestinationLocal
		case "volume":
			blockDeviceOpts[i].DestinationType = servers.DestinationVolume
		default:
			return blockDeviceOpts, fmt.Errorf("unknown block device destination type %s", destinationType)
		}
	}

	log.Printf("[DEBUG] Block Device Options: %+v", blockDeviceOpts)

	return blockDeviceOpts, nil
}

func resourceInstanceSchedulerHintsV2(schedulerHintsRaw map[string]any) servers.SchedulerHintOpts {
	differentHost := []string{}

	if v, ok := schedulerHintsRaw["different_host"].([]any); ok {
		for _, dh := range v {
			differentHost = append(differentHost, dh.(string))
		}
	}

	sameHost := []string{}

	if v, ok := schedulerHintsRaw["same_host"].([]any); ok {
		for _, sh := range v {
			sameHost = append(sameHost, sh.(string))
		}
	}

	query := []any{}

	if v, ok := schedulerHintsRaw["query"].([]any); ok {
		for _, q := range v {
			query = append(query, q.(string))
		}
	}

	differentCell := []string{}

	if v, ok := schedulerHintsRaw["different_cell"].([]any); ok {
		for _, dh := range v {
			differentCell = append(differentCell, dh.(string))
		}
	}

	schedulerHints := servers.SchedulerHintOpts{
		Group:                schedulerHintsRaw["group"].(string),
		DifferentHost:        differentHost,
		SameHost:             sameHost,
		Query:                query,
		TargetCell:           schedulerHintsRaw["target_cell"].(string),
		DifferentCell:        differentCell,
		BuildNearHostIP:      schedulerHintsRaw["build_near_host_ip"].(string),
		AdditionalProperties: schedulerHintsRaw["additional_properties"].(map[string]any),
	}

	return schedulerHints
}

func getImageIDFromConfig(ctx context.Context, imageClient *gophercloud.ServiceClient, d *schema.ResourceData) (string, error) {
	// If block_device was used, an Image does not need to be specified, unless an image/local
	// combination was used. This emulates normal boot behavior. Otherwise, ignore the image altogether.
	if vL, ok := d.GetOk("block_device"); ok {
		needImage := false

		for _, v := range vL.([]any) {
			vM := v.(map[string]any)
			if vM["source_type"] == "image" && vM["destination_type"] == "local" {
				needImage = true
			}
		}

		if !needImage {
			return "", nil
		}
	}

	if imageID := d.Get("image_id").(string); imageID != "" {
		return imageID, nil
	}
	// try the OS_IMAGE_ID environment variable
	if v := os.Getenv("OS_IMAGE_ID"); v != "" {
		return v, nil
	}

	imageName := d.Get("image_name").(string)
	if imageName == "" {
		// try the OS_IMAGE_NAME environment variable
		if v := os.Getenv("OS_IMAGE_NAME"); v != "" {
			imageName = v
		}
	}

	if imageName != "" {
		imageID, err := imagesutils.IDFromName(ctx, imageClient, imageName)
		if err != nil {
			return "", err
		}

		return imageID, nil
	}

	return "", errors.New("Neither a boot device, image ID, or image name were able to be determined")
}

func setImageInformation(ctx context.Context, imageClient *gophercloud.ServiceClient, server *servers.Server, d *schema.ResourceData) error {
	// If block_device was used, an Image does not need to be specified, unless an image/local
	// combination was used. This emulates normal boot behavior. Otherwise, ignore the image altogether.
	if vL, ok := d.GetOk("block_device"); ok {
		needImage := false

		for _, v := range vL.([]any) {
			vM := v.(map[string]any)
			if vM["source_type"] == "image" && vM["destination_type"] == "local" {
				needImage = true
			}
		}

		if !needImage {
			d.Set("image_id", "Attempt to boot from volume - no image supplied")

			return nil
		}
	}

	if server.Image["id"] != nil {
		imageID := server.Image["id"].(string)
		if imageID != "" {
			d.Set("image_id", imageID)

			image, err := images.Get(ctx, imageClient, imageID).Extract()
			if err != nil {
				if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
					// If the image name can't be found, set the value to "Image not found".
					// The most likely scenario is that the image no longer exists in the Image Service
					// but the instance still has a record from when it existed.
					d.Set("image_name", "Image not found")

					return nil
				}

				return err
			}

			d.Set("image_name", image.Name)
		}
	}

	return nil
}

func getFlavorID(ctx context.Context, computeClient *gophercloud.ServiceClient, d *schema.ResourceData) (string, error) {
	if flavorID := d.Get("flavor_id").(string); flavorID != "" {
		return flavorID, nil
	}
	// Try the OS_FLAVOR_ID environment variable
	if v := os.Getenv("OS_FLAVOR_ID"); v != "" {
		return v, nil
	}

	flavorName := d.Get("flavor_name").(string)
	if flavorName == "" {
		// Try the OS_FLAVOR_NAME environment variable
		if v := os.Getenv("OS_FLAVOR_NAME"); v != "" {
			flavorName = v
		}
	}

	if flavorName != "" {
		flavorID, err := flavorsutils.IDFromName(ctx, computeClient, flavorName)
		if err != nil {
			return "", err
		}

		return flavorID, nil
	}

	return "", errors.New("Neither a flavor_id or flavor_name could be determined")
}

func resourceComputeSchedulerHintsHash(v any) int {
	var buf bytes.Buffer

	m, ok := v.(map[string]any)
	if !ok {
		return hashcode.String(buf.String())
	}

	if m == nil {
		return hashcode.String(buf.String())
	}

	if m["group"] != nil {
		buf.WriteString(m["group"].(string) + "-")
	}

	if m["target_cell"] != nil {
		buf.WriteString(m["target_cell"].(string) + "-")
	}

	if m["build_host_near_ip"] != nil {
		buf.WriteString(m["build_host_near_ip"].(string) + "-")
	}

	if m["additional_properties"] != nil {
		for _, v := range m["additional_properties"].(map[string]any) {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}

	buf.WriteString(fmt.Sprintf("%s-", m["different_host"].([]any)))
	buf.WriteString(fmt.Sprintf("%s-", m["same_host"].([]any)))
	buf.WriteString(fmt.Sprintf("%s-", m["query"].([]any)))
	buf.WriteString(fmt.Sprintf("%s-", m["different_cell"].([]any)))

	return hashcode.String(buf.String())
}

func checkBlockDeviceConfig(d *schema.ResourceData) error {
	if vL, ok := d.GetOk("block_device"); ok {
		for _, v := range vL.([]any) {
			vM := v.(map[string]any)

			if vM["source_type"] != "blank" && vM["uuid"] == "" {
				return fmt.Errorf("You must specify a uuid for %s block device types", vM["source_type"])
			}

			if vM["source_type"] == "image" && vM["destination_type"] == "volume" {
				if vM["volume_size"] == 0 {
					return errors.New("You must specify a volume_size when creating a volume from an image")
				}
			}

			if vM["source_type"] == "blank" && vM["destination_type"] == "local" {
				if vM["volume_size"] == 0 {
					return errors.New("You must specify a volume_size when creating a blank block device")
				}
			}
		}
	}

	return nil
}

func resourceComputeInstancePersonalityHash(v any) int {
	var buf bytes.Buffer

	m := v.(map[string]any)
	buf.WriteString(m["file"].(string) + "-")

	return hashcode.String(buf.String())
}

func resourceInstancePersonalityV2(d *schema.ResourceData) servers.Personality {
	var personalities servers.Personality

	if v := d.Get("personality"); v != nil {
		personalityList := v.(*schema.Set).List()
		if len(personalityList) > 0 {
			for _, p := range personalityList {
				rawPersonality := p.(map[string]any)
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
func suppressAvailabilityZoneDetailDiffs(_, o, n string, _ *schema.ResourceData) bool {
	if strings.Contains(n, ":") {
		parts := strings.Split(n, ":")
		az := parts[0]

		if az == o {
			return true
		}
	}

	return false
}

// suppressPowerStateDiffs will allow a state of "error" or "migrating" even though we don't
// allow them as a user input.
func suppressPowerStateDiffs(_, old, _ string, _ *schema.ResourceData) bool {
	if old == "error" || old == "migrating" {
		return true
	}

	return false
}

// validateHostname retruns a validation function which checks if the supplied hostname is a
// valid FQDN or hostname. While underscores are not allowed in RFC952, nova accepts them.
// https://github.com/openstack/nova/blob/0d586ccca88ae90b9634ee00b8f7f86a78b09cd0/nova/api/validation/parameter_types.py#L269-L279
func validateHostname() schema.SchemaValidateFunc {
	r := regexp.MustCompile(`^[a-zA-Z0-9-\._]{1,255}$`)

	return validation.StringMatch(r, "Invalid hostname. only alphanumeric, . (dot), - (dash) and _ (underscore) are allowed characters in the hostname.")
}

// isValidHostname checks if the supplied hostname matches the regexp defined in the nova API.
// https://github.com/openstack/nova/blob/0d586ccca88ae90b9634ee00b8f7f86a78b09cd0/nova/api/validation/parameter_types.py#L262-L266
func isValidHostname(hostname string) bool {
	if len(hostname) < 2 || len(hostname) > 63 {
		return false
	}

	return regexp.MustCompile(`^[a-zA-Z0-9]+[a-zA-Z0-9-]*[a-zA-Z0-9]+$`).MatchString(hostname)
}
