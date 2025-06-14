package openstack

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	errs "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/errors"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shareaccessrules"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSharedFilesystemShareAccessV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFilesystemShareAccessV2Create,
		ReadContext:   resourceSharedFilesystemShareAccessV2Read,
		DeleteContext: resourceSharedFilesystemShareAccessV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSharedFilesystemShareAccessV2Import,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"access_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip", "user", "cert", "cephx",
				}, false),
			},

			"access_to": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"access_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"rw", "ro",
				}, false),
			},

			"access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSharedFilesystemShareAccessV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemV2MinMicroversion
	accessType := d.Get("access_type").(string)

	if accessType == "cephx" {
		sfsClient.Microversion = sharedFilesystemV2SharedAccessCephXMicroversion
	}

	shareID := d.Get("share_id").(string)

	grantOpts := shares.GrantAccessOpts{
		AccessType:  accessType,
		AccessTo:    d.Get("access_to").(string),
		AccessLevel: d.Get("access_level").(string),
	}

	log.Printf("[DEBUG] openstack_sharedfilesystem_share_access_v2 create options: %#v", grantOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	var access *shares.AccessRight

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		access, err = shares.GrantAccess(ctx, sfsClient, shareID, grantOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		detailedErr := errs.ErrorDetails{}

		e := errs.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Error creating openstack_sharedfilesystem_share_access_v2: %s: %s", err, e)
		}

		for k, msg := range detailedErr {
			return diag.Errorf("Error creating openstack_sharedfilesystem_share_access_v2: %s (%d): %s", k, msg.Code, msg.Message)
		}
	}

	log.Printf("[DEBUG] Waiting for openstack_sharedfilesystem_share_access_v2 %s to become available.", access.ID)
	stateConf := &retry.StateChangeConf{
		Target:     []string{"active"},
		Pending:    []string{"new", "queued_to_apply", "applying"},
		Refresh:    sharedFilesystemShareAccessV2StateRefreshFunc(ctx, sfsClient, shareID, access.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_sharedfilesystem_share_access_v2 %s to become available: %s", access.ID, err)
	}

	d.SetId(access.ID)

	return resourceSharedFilesystemShareAccessV2Read(ctx, d, meta)
}

func resourceSharedFilesystemShareAccessV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	shareID := d.Get("share_id").(string)

	access, _, err := sharedFilesystemShareAccessV2StateRefreshFunc(ctx, sfsClient, shareID, d.Id())()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Failed to retrieve openstack_sharedfilesystem_share_access_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_sharedfilesystem_share_access_v2 %s: %#v", d.Id(), access)

	switch access := access.(type) {
	case shareaccessrules.ShareAccess:
		d.Set("access_type", access.AccessType)
		d.Set("access_to", access.AccessTo)
		d.Set("access_level", access.AccessLevel)
		d.Set("region", GetRegion(d, config))
		d.Set("access_key", access.AccessKey)
		d.Set("state", access.State)

		return nil
	case shares.AccessRight:
		d.Set("access_type", access.AccessType)
		d.Set("access_to", access.AccessTo)
		d.Set("access_level", access.AccessLevel)
		d.Set("region", GetRegion(d, config))
		d.Set("state", access.State)

		// This will only be set if the Shared Filesystem environment supports
		// microversion 2.21.
		d.Set("access_key", access.AccessKey)

		return nil
	}

	return diag.Errorf("Unknown share access rules type: %T", access)
}

func resourceSharedFilesystemShareAccessV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemV2MinMicroversion

	shareID := d.Get("share_id").(string)

	revokeOpts := shares.RevokeAccessOpts{AccessID: d.Id()}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Attempting to delete openstack_sharedfilesystem_share_access_v2 %s", d.Id())

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = shares.RevokeAccess(ctx, sfsClient, shareID, revokeOpts).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		e := CheckDeleted(d, err, "Error deleting openstack_sharedfilesystem_share_access_v2")
		if e == nil {
			return nil
		}

		detailedErr := errs.ErrorDetails{}

		e = errs.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Error waiting for openstack_sharedfilesystem_share_access_v2 on %s to be removed: %s: %s", shareID, err, e)
		}

		for k, msg := range detailedErr {
			return diag.Errorf("Error waiting for openstack_sharedfilesystem_share_access_v2 on %s to be removed: %s (%d): %s", shareID, k, msg.Code, msg.Message)
		}
	}

	log.Printf("[DEBUG] Waiting for openstack_sharedfilesystem_share_access_v2 %s to become denied.", d.Id())
	stateConf := &retry.StateChangeConf{
		Target:     []string{"denied"},
		Pending:    []string{"active", "new", "queued_to_deny", "denying"},
		Refresh:    sharedFilesystemShareAccessV2StateRefreshFunc(ctx, sfsClient, shareID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return nil
		}

		return diag.Errorf("Error waiting for openstack_sharedfilesystem_share_access_v2 %s to become denied: %s", d.Id(), err)
	}

	return nil
}

func resourceSharedFilesystemShareAccessV2Import(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := errors.New("Invalid format specified for openstack_sharedfilesystem_share_access_v2. Format must be <share id>/<ACL id>")

		return nil, err
	}

	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack sharedfilesystem client: %w", err)
	}

	sfsClient.Microversion = sharedFilesystemV2MinMicroversion

	shareID := parts[0]
	accessID := parts[1]

	access, err := shares.ListAccessRights(ctx, sfsClient, shareID).Extract()
	if err != nil {
		return nil, fmt.Errorf("Unable to get %s openstack_sharedfilesystem_share_v2: %w", shareID, err)
	}

	for _, v := range access {
		if v.ID == accessID {
			log.Printf("[DEBUG] Retrieved openstack_sharedfilesystem_share_access_v2 %s: %#v", accessID, v)

			d.SetId(accessID)
			d.Set("share_id", shareID)

			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("[DEBUG] Unable to find openstack_sharedfilesystem_share_access_v2 %s", accessID)
}
