package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/nodegroups"
)

func dataSourceContainerInfraNodeGroupV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerInfraNodeGroupRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"docker_volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"min_node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"image": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceContainerInfraNodeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	containerInfraClient.Microversion = containerInfraV1NodeGroupMinMicroversion

	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	nodeGroup, err := nodegroups.Get(containerInfraClient, clusterID, name).Extract()
	if err != nil {
		return diag.Errorf("Error getting openstack_containerinfra_nodegroup_v1 %s: %s", name, err)
	}

	d.SetId(nodeGroup.UUID)

	d.Set("project_id", nodeGroup.ProjectID)
	d.Set("docker_volume_size", nodeGroup.DockerVolumeSize)
	d.Set("role", nodeGroup.Role)
	d.Set("node_count", nodeGroup.NodeCount)
	d.Set("min_node_count", nodeGroup.MinNodeCount)
	d.Set("max_node_count", nodeGroup.MaxNodeCount)
	d.Set("image", nodeGroup.ImageID)
	d.Set("flavor", nodeGroup.FlavorID)

	if err := d.Set("labels", nodeGroup.Labels); err != nil {
		log.Printf("[DEBUG] Unable to set labels for openstack_containerinfra_nodegroup_v1 %s: %s", nodeGroup.UUID, err)
	}
	if err := d.Set("created_at", nodeGroup.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set created_at for openstack_containerinfra_nodegroup_v1 %s: %s", nodeGroup.UUID, err)
	}
	if err := d.Set("updated_at", nodeGroup.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set updated_at for openstack_containerinfra_nodegroup_v1 %s: %s", nodeGroup.UUID, err)
	}

	d.Set("region", GetRegion(d, config))

	return nil
}
