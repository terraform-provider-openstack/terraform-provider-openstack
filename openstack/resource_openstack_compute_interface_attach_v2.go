package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/attachinterfaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComputeInterfaceAttachV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeInterfaceAttachV2Create,
		ReadContext:   resourceComputeInterfaceAttachV2Read,
		DeleteContext: resourceComputeInterfaceAttachV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"port_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"network_id"},
			},

			"network_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"port_id"},
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fixed_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"port_id"},
			},
		},
	}
}

func resourceComputeInterfaceAttachV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	var portID string
	if v, ok := d.GetOk("port_id"); ok {
		portID = v.(string)
	}

	var networkID string
	if v, ok := d.GetOk("network_id"); ok {
		networkID = v.(string)
	}

	if networkID == "" && portID == "" {
		return diag.Errorf("Must set one of network_id and port_id")
	}

	// For some odd reason the API takes an array of IPs, but you can only have one element in the array.
	var fixedIPs []attachinterfaces.FixedIP
	if v, ok := d.GetOk("fixed_ip"); ok {
		fixedIPs = append(fixedIPs, attachinterfaces.FixedIP{IPAddress: v.(string)})
	}

	attachOpts := attachinterfaces.CreateOpts{
		PortID:    portID,
		NetworkID: networkID,
		FixedIPs:  fixedIPs,
	}

	log.Printf("[DEBUG] openstack_compute_interface_attach_v2 attach options: %#v", attachOpts)

	attachment, err := attachinterfaces.Create(ctx, computeClient, instanceID, attachOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ATTACHING"},
		Target:     []string{"ATTACHED"},
		Refresh:    computeInterfaceAttachV2AttachFunc(ctx, computeClient, instanceID, attachment.PortID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error creating openstack_compute_interface_attach_v2 %s: %s", instanceID, err)
	}

	// Use the instance ID and attachment ID as the resource ID.
	id := fmt.Sprintf("%s/%s", instanceID, attachment.PortID)

	log.Printf("[DEBUG] Created openstack_compute_interface_attach_v2 %s: %#v", id, attachment)

	d.SetId(id)

	return resourceComputeInterfaceAttachV2Read(ctx, d, meta)
}

func resourceComputeInterfaceAttachV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	instanceID, attachmentID, err := parsePairedIDs(d.Id(), "openstack_compute_interface_attach_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	attachment, err := attachinterfaces.Get(ctx, computeClient, instanceID, attachmentID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_compute_interface_attach_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_compute_interface_attach_v2 %s: %#v", d.Id(), attachment)

	if len(attachment.FixedIPs) > 0 {
		d.Set("fixed_ip", attachment.FixedIPs[0].IPAddress)
	}

	d.Set("instance_id", instanceID)
	d.Set("port_id", attachment.PortID)
	d.Set("network_id", attachment.NetID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeInterfaceAttachV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	instanceID, attachmentID, err := parsePairedIDs(d.Id(), "openstack_compute_interface_attach_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"DETACHED"},
		Refresh:    computeInterfaceAttachV2DetachFunc(ctx, computeClient, instanceID, attachmentID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error detaching openstack_compute_interface_attach_v2 %s: %s", d.Id(), err)
	}

	return nil
}
