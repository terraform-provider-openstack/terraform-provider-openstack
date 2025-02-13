package openstack

import (
	"context"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ZoneShare represents a shared DNS zone in OpenStack.
type ZoneShare struct {
	ID              string `json:"id"`
	TargetProjectID string `json:"target_project_id"`
}

// resourceDNSZoneShareV2 defines the resource to share a DNS zone.
func resourceDNSZoneShareV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSZoneShareV2Create,
		Read:   resourceDNSZoneShareV2Read,
		Delete: resourceDNSZoneShareV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the DNS zone to share",
			},
			"target_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target project ID with which to share the DNS zone",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The owner project ID required to authorize sharing",
			},
			"share_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the created share",
			},
		},
	}
}

func resourceDNSZoneShareV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.DNSV2Client(context.Background(), "")
	if err != nil {
		return fmt.Errorf("error creating OpenStack DNS client: %s", err)
	}

	zoneID := d.Get("zone_id").(string)
	projectID := d.Get("target_project_id").(string)
	sudoProjectID := d.Get("project_id").(string)

	shareOpts := zones.ShareZoneOpts{
		TargetProjectID: projectID,
	}

	client.Microversion = "2.0"
	client.MoreHeaders = map[string]string{
		"X-Auth-Sudo-Project-Id": sudoProjectID,
	}

	// Log API details for debugging
	fmt.Printf("DEBUG: Making API call to share DNS zone\n")
	fmt.Printf("DEBUG: URL: %s\n", client.ServiceURL("zones", zoneID, "shares"))
	fmt.Printf("DEBUG: Headers: X-Auth-Sudo-Project-Id: %s\n", sudoProjectID)
	fmt.Printf("DEBUG: Request Body: %+v\n", shareOpts)

	// Execute the API call
	if err := zones.Share(context.Background(), client, zoneID, shareOpts).ExtractErr(); err != nil {
		return fmt.Errorf("error sharing DNS zone %s with project %s: %s", zoneID, projectID, err)
	}

	// Fetch shares to get the share ID
	shares, err := listZoneShares(client, zoneID, sudoProjectID)
	if err != nil {
		return err
	}

	var shareID string
	for _, s := range shares {
		if s.TargetProjectID == projectID {
			shareID = s.ID
			break
		}
	}
	if shareID == "" {
		return fmt.Errorf("failed to locate share for DNS zone %s with project %s", zoneID, projectID)
	}

	// Store the resource ID
	id := fmt.Sprintf("%s/%s", zoneID, shareID)
	d.SetId(id)
	d.Set("share_id", shareID)
	return resourceDNSZoneShareV2Read(d, meta)
}

func resourceDNSZoneShareV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.DNSV2Client(context.Background(), "")
	if err != nil {
		return fmt.Errorf("error creating OpenStack DNS client: %s", err)
	}

	zoneID, shareID, err := parseDnsSharedZoneID(d.Id())
	if err != nil {
		return err
	}
	projectID := d.Get("target_project_id").(string)
	ownerProjectID := d.Get("project_id").(string)

	shares, err := listZoneShares(client, zoneID, ownerProjectID)
	if err != nil {
		return fmt.Errorf("error listing shares for DNS zone %s: %s", zoneID, err)
	}

	found := false
	for _, s := range shares {
		if s.ID == shareID && s.TargetProjectID == projectID {
			found = true
			break
		}
	}
	if !found {
		d.SetId("")
	}
	return nil
}

func resourceDNSZoneShareV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.DNSV2Client(context.Background(), "")
	if err != nil {
		return fmt.Errorf("error creating OpenStack DNS client: %s", err)
	}

	zoneID, shareID, err := parseDnsSharedZoneID(d.Id())
	if err != nil {
		return err
	}

	sudoProjectID := d.Get("project_id").(string)

	// Construct the URL for the delete request
	url := client.ServiceURL("zones", zoneID, "shares", shareID)

	// Create the request options with the necessary header
	reqOpts := &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{
			"X-Auth-Sudo-Project-Id": sudoProjectID,
		},
	}

	// Perform the delete request
	_, err = client.Delete(context.Background(), url, reqOpts)
	if err != nil {
		return fmt.Errorf("error unsharing DNS zone %s: %s", zoneID, err)
	}

	return nil
}

func parseDnsSharedZoneID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected ID format (%s), expected <zoneID>/<shareID>", id)
	}
	return parts[0], parts[1], nil
}

func listZoneShares(client *gophercloud.ServiceClient, zoneID string, ownerProjectID string) ([]ZoneShare, error) {
	url := client.ServiceURL("zones", zoneID, "shares")
	var result struct {
		Shares []ZoneShare `json:"shared_zones"`
	}

	reqOpts := &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{
			"X-Auth-Sudo-Project-Id": ownerProjectID,
		},
	}

	_, err := client.Get(context.Background(), url, &result, reqOpts)
	if err != nil {
		return nil, fmt.Errorf("error getting shares for zone %s: %s", zoneID, err)
	}

	return result.Shares, nil
}
