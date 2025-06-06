package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceComputeKeypairV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeKeypairV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// computed-only
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeKeypairV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	computeClient.Microversion = computeKeyPairV2UserIDMicroversion

	opts := keypairs.GetOpts{}

	// Check if searching for the keypair of another user
	userID := d.Get("user_id").(string)
	if userID != "" {
		opts.UserID = userID
	}

	name := d.Get("name").(string)
	kp, err := keypairs.Get(ctx, computeClient, name, opts).Extract()
	if err != nil {
		return diag.Errorf("Error retrieving openstack_compute_keypair_v2 %s: %s", name, err)
	}

	d.SetId(name)

	log.Printf("[DEBUG] Retrieved openstack_compute_keypair_v2 %s: %#v", d.Id(), kp)

	d.Set("fingerprint", kp.Fingerprint)
	d.Set("public_key", kp.PublicKey)
	d.Set("region", GetRegion(d, config))
	d.Set("user_id", kp.UserID)

	return nil
}
