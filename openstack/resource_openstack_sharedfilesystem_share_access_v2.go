package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

func resourceSharedFilesystemShareAccessV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceSharedFilesystemShareAccessV2Grant,
		Read:   resourceSharedFilesystemShareAccessV2Read,
		Delete: resourceSharedFilesystemShareAccessV2Revoke,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"share_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"access_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip", "user", "cert",
				}, true),
			},

			"access_to": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"access_level": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"rw", "ro",
				}, true),
			},
		},
	}
}

func resourceSharedFilesystemShareAccessV2Grant(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = minManilaMicroversion

	shareId := d.Get("share_id").(string)

	grantOpts := shares.GrantAccessOpts{
		AccessType:  d.Get("access_type").(string),
		AccessTo:    d.Get("access_to").(string),
		AccessLevel: d.Get("access_level").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", grantOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	log.Printf("[DEBUG] Attempting to grant access")
	var access *shares.AccessRight
	err = resource.Retry(timeout, func() *resource.RetryError {
		access, err = shares.GrantAccess(sfsClient, shareId, grantOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error granting access: %s", err)
	}

	d.SetId(access.ID)

	pending := []string{"new", "queued_to_apply", "applying"}
	// Wait for access to become active before continuing
	err = waitForSFV2Access(sfsClient, shareId, access.ID, "active", pending, timeout)
	if err != nil {
		return err
	}

	return resourceSharedFilesystemShareAccessV2Read(d, meta)
}

func resourceSharedFilesystemShareAccessV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = minManilaMicroversion

	shareId := d.Get("share_id").(string)

	access, err := shares.ListAccessRights(sfsClient, shareId).Extract()
	if err != nil {
		return CheckDeleted(d, err, "share")
	}

	for _, v := range access {
		if v.ID == d.Id() {
			log.Printf("[DEBUG] Retrieved access %s: %#v", d.Id(), access)

			d.Set("access_type", v.AccessType)
			d.Set("access_to", v.AccessTo)
			d.Set("access_level", v.AccessLevel)

			return nil
		}
	}

	d.Set("access_type", "")
	d.Set("access_to", "")
	d.Set("access_level", "")

	return nil
}

func resourceSharedFilesystemShareAccessV2Revoke(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = minManilaMicroversion

	shareId := d.Get("share_id").(string)

	revokeOpts := shares.RevokeAccessOpts{AccessID: d.Id()}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Attempting to revoke access %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = shares.RevokeAccess(sfsClient, shareId, revokeOpts).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Wait for access to become deleted before continuing
	pending := []string{"new", "queued_to_deny", "denying"}
	err = waitForSFV2Access(sfsClient, shareId, d.Id(), "denied", pending, timeout)
	if err != nil {
		return err
	}

	return nil
}

// Full list of the share access statuses: https://developer.openstack.org/api-ref/shared-file-system/?expanded=list-services-detail,list-access-rules-detail#list-access-rules
func waitForSFV2Access(sfsClient *gophercloud.ServiceClient, shareId string, id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for access %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceSFV2AccessRefreshFunc(sfsClient, shareId, id),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "denied":
				return nil
			default:
				return fmt.Errorf("Error: access %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("Error waiting for access %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceSFV2AccessRefreshFunc(sfsClient *gophercloud.ServiceClient, shareId string, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		access, err := shares.ListAccessRights(sfsClient, shareId).Extract()
		if err != nil {
			return nil, "", err
		}
		for _, v := range access {
			if v.ID == id {
				return v, v.State, nil
			}
		}
		return nil, "", gophercloud.ErrDefault404{}
	}
}
