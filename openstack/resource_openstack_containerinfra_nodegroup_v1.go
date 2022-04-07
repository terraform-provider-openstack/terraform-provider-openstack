package openstack

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/nodegroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContainerInfraNodeGroupV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerInfraNodeGroupV1Create,
		ReadContext:   resourceContainerInfraNodeGroupV1Read,
		UpdateContext: resourceContainerInfraNodeGroupV1Update,
		DeleteContext: resourceContainerInfraNodeGroupV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"docker_volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"role": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"min_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"max_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"image": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"flavor": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceContainerInfraNodeGroupV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	containerInfraClient.Microversion = containerInfraV1NodeGroupMinMicroversion

	// Get and check labels map.
	rawLabels := d.Get("labels").(map[string]interface{})
	labels, err := expandContainerInfraV1LabelsMap(rawLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := nodegroups.CreateOpts{
		Name:         d.Get("name").(string),
		Labels:       labels,
		MinNodeCount: d.Get("min_node_count").(int),
		Role:         d.Get("role").(string),
		ImageID:      d.Get("image").(string),
		FlavorID:     d.Get("flavor").(string),
	}

	// Set int parameters that will be passed by reference.
	dockerVolumeSize := d.Get("docker_volume_size").(int)
	if dockerVolumeSize > 0 {
		createOpts.DockerVolumeSize = &dockerVolumeSize
	}
	nodeCount := d.Get("node_count").(int)
	if nodeCount > 0 {
		createOpts.NodeCount = &nodeCount
	}
	maxNodeCount := d.Get("max_node_count").(int)
	if maxNodeCount > 0 {
		createOpts.MaxNodeCount = &maxNodeCount
	}

	log.Printf("[DEBUG] openstack_containerinfra_nodegroup_v1 create options: %#v", createOpts)

	clusterId := d.Get("cluster_id").(string)
	nodeGroup, err := nodegroups.Create(containerInfraClient, clusterId, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_containerinfra_nodegroup_v1: %s", err)
	}

	id := fmt.Sprintf("%s/%s", clusterId, nodeGroup.UUID)
	d.SetId(id)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATE_IN_PROGRESS"},
		Target:       []string{"CREATE_COMPLETE"},
		Refresh:      containerInfraNodeGroupV1StateRefreshFunc(containerInfraClient, clusterId, nodeGroup.UUID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Minute,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_containerinfra_nodegroup_v1 %s to become ready: %s", nodeGroup.UUID, err)
	}

	log.Printf("[DEBUG] Created openstack_containerinfra_nodegroup_v1 %s", nodeGroup.UUID)

	return resourceContainerInfraNodeGroupV1Read(ctx, d, meta)
}

func resourceContainerInfraNodeGroupV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	containerInfraClient.Microversion = containerInfraV1NodeGroupMinMicroversion

	clusterId, nodeGroupId, err := parseNodeGroupID(d.Id())
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error parsing ID of openstack_containerinfra_nodegroup_v1"))
	}

	nodeGroup, err := nodegroups.Get(containerInfraClient, clusterId, nodeGroupId).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_containerinfra_nodegroup_v1"))
	}

	log.Printf("[DEBUG] Retrieved openstack_containerinfra_nodegroup_v1 %s: %#v", d.Id(), nodeGroup)

	if err := d.Set("labels", nodeGroup.Labels); err != nil {
		return diag.Errorf("Unable to set labels: %s", err)
	}

	d.Set("cluster_id", clusterId)
	d.Set("region", GetRegion(d, config))
	d.Set("name", nodeGroup.Name)
	d.Set("project_id", nodeGroup.ProjectID)
	d.Set("role", nodeGroup.Role)
	d.Set("node_count", nodeGroup.NodeCount)
	d.Set("min_node_count", nodeGroup.NodeCount)
	d.Set("max_node_count", nodeGroup.NodeCount)
	d.Set("image_id", nodeGroup.ImageID)
	d.Set("flavor_id", nodeGroup.FlavorID)

	if err := d.Set("created_at", nodeGroup.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_containerinfra_nodegroup_v1 created_at: %s", err)
	}
	if err := d.Set("updated_at", nodeGroup.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_containerinfra_nodegroup_v1 updated_at: %s", err)
	}

	return nil
}

func resourceContainerInfraNodeGroupV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	containerInfraClient.Microversion = containerInfraV1NodeGroupMinMicroversion

	clusterId, nodeGroupId, err := parseNodeGroupID(d.Id())
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error parsing ID of openstack_containerinfra_nodegroup_v1"))
	}

	updateOpts := []nodegroups.UpdateOptsBuilder{}

	if d.HasChange("min_node_count") {
		v := d.Get("min_node_count").(int)
		minNodeCount := strconv.Itoa(v)
		updateOpts = containerInfraNodeGroupV1AppendUpdateOpts(
			updateOpts, "min_node_count", minNodeCount)
	}

	if d.HasChange("max_node_count") {
		v := d.Get("max_node_count").(int)
		maxNodeCount := strconv.Itoa(v)
		updateOpts = containerInfraNodeGroupV1AppendUpdateOpts(
			updateOpts, "max_node_count", maxNodeCount)
	}

	if len(updateOpts) > 0 {
		log.Printf(
			"[DEBUG] Updating openstack_containerinfra_nodegroup_v1 %s with options: %#v", d.Id(), updateOpts)

		_, err = nodegroups.Update(containerInfraClient, clusterId, nodeGroupId, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_containerinfra_nodegroup_v1 %s: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"UPDATE_IN_PROGRESS"},
			Target:       []string{"UPDATE_COMPLETE"},
			Refresh:      containerInfraNodeGroupV1StateRefreshFunc(containerInfraClient, clusterId, nodeGroupId),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        1 * time.Minute,
			PollInterval: 20 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"Error waiting for openstack_containerinfra_node_group_v1 %s to become updated: %s", d.Id(), err)
		}
	}
	return resourceContainerInfraNodeGroupV1Read(ctx, d, meta)
}

func resourceContainerInfraNodeGroupV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	containerInfraClient.Microversion = containerInfraV1NodeGroupMinMicroversion

	clusterId, nodeGroupId, err := parseNodeGroupID(d.Id())
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error parsing ID of openstack_containerinfra_nodegroup_v1"))
	}

	if err := nodegroups.Delete(containerInfraClient, clusterId, nodeGroupId).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_containerinfra_nodegroup_v1"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETE_IN_PROGRESS"},
		Target:       []string{"DELETE_COMPLETE"},
		Refresh:      containerInfraNodeGroupV1StateRefreshFunc(containerInfraClient, clusterId, nodeGroupId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_containerinfra_nodegroup_v1 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}

func parseNodeGroupID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("Unable to determine nodegroup ID %s", id)
	}

	return idParts[0], idParts[1], nil
}
