package openstack

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/objects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceObjectstorageTempurlV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObjectstorageTempurlV1Create,
		ReadContext:   resourceObjectstorageTempurlV1Read,
		Delete:        schema.RemoveFromState,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"container": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"object": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"method": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "get",
				ValidateFunc: func(v any, _ string) (ws []string, errs []error) {
					value := v.(string)
					if value != "get" && value != "post" {
						errs = append(errs, errors.New("Only 'get', and 'post' are supported values for 'method'"))
					}

					return
				},
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"split": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"digest": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"sha1", "sha256", "sha512"}, false),
			},

			"regenerate": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"url": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

// resourceObjectstorageTempurlV1Create performs the image lookup.
func resourceObjectstorageTempurlV1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	objectStorageClient, err := config.ObjectStorageV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	method := objects.GET

	switch d.Get("method") {
	case "post":
		method = objects.POST
		// gophercloud doesn't have support for PUT yet,
		// although it's a valid method for swift
		// case "put":
		//	method = objects.PUT
	}

	turlOptions := objects.CreateTempURLOpts{
		Method:     method,
		TTL:        d.Get("ttl").(int),
		Split:      d.Get("split").(string),
		TempURLKey: d.Get("key").(string),
		Digest:     d.Get("digest").(string),
	}

	containerName := d.Get("container").(string)
	objectName := d.Get("object").(string)

	log.Printf("[DEBUG] Create temporary url Options: %#v", turlOptions)

	url, err := objects.CreateTempURL(ctx, objectStorageClient, containerName, objectName, turlOptions)
	if err != nil {
		return diag.Errorf("Unable to generate a temporary url for the object %s in container %s: %s",
			objectName, containerName, err)
	}

	log.Printf("[DEBUG] URL Generated: %s", url)

	// Set the URL and Id fields.
	hasher := md5.New()
	hasher.Write([]byte(url))
	d.SetId(hex.EncodeToString(hasher.Sum(nil)))
	d.Set("url", url)
	d.Set("region", GetRegion(d, config))

	return nil
}

// resourceObjectstorageTempurlV1Read performs the image lookup.
func resourceObjectstorageTempurlV1Read(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	d.Set("region", GetRegion(d, config))

	turl := d.Get("url").(string)
	u, err := url.Parse(turl)
	if err != nil {
		return diag.Errorf("Failed to read the temporary url %s: %s", turl, err)
	}

	qp, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return diag.Errorf("Failed to parse the temporary url %s query string: %s", turl, err)
	}

	tempURLExpires := qp.Get("temp_url_expires")
	expiry, err := strconv.ParseInt(tempURLExpires, 10, 64)
	if err != nil {
		return diag.Errorf(
			"Failed to parse the temporary url %s expiration time %s: %s",
			turl, tempURLExpires, err)
	}

	// Regenerate the URL if it has expired and if the user requested it to be.
	regen := d.Get("regenerate").(bool)
	now := time.Now().Unix()

	if expiry < now && regen {
		log.Printf("[DEBUG] temporary url %s expired, generating a new one", turl)
		d.SetId("")
	}

	return nil
}
