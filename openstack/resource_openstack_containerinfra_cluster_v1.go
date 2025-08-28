package openstack

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clusters"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/nodegroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContainerInfraClusterV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerInfraClusterV1Create,
		ReadContext:   resourceContainerInfraClusterV1Read,
		UpdateContext: resourceContainerInfraClusterV1Update,
		DeleteContext: resourceContainerInfraClusterV1Delete,
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

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},

			"user_id": {
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

			"api_address": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"coe_version": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"cluster_template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				DefaultFunc: schema.EnvDefaultFunc("OS_MAGNUM_CLUSTER_TEMPLATE", nil),
			},

			"container_version": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"create_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"discovery_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"docker_volume_size": {
				Type:     schema.TypeInt,
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

			"master_flavor": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"keypair": {
				Type:     schema.TypeString,
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

			"merge_labels": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"master_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"master_lb_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  1,
			},

			"master_addresses": {
				Type:     schema.TypeList,
				ForceNew: false,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"node_addresses": {
				Type:     schema.TypeList,
				ForceNew: false,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"stack_id": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"fixed_network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"fixed_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"floating_ip_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"kubeconfig": {
				Type:      schema.TypeMap,
				Computed:  true,
				Sensitive: true,
				Elem:      &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceContainerInfraClusterV1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	containerInfraClient, err := config.ContainerInfraV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	// Get and check labels map.
	rawLabels := d.Get("labels").(map[string]any)

	labels, err := expandContainerInfraV1LabelsMap(rawLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := clusters.CreateOpts{
		ClusterTemplateID: d.Get("cluster_template_id").(string),
		DiscoveryURL:      d.Get("discovery_url").(string),
		FlavorID:          containerInfraClusterV1Flavor(d),
		Keypair:           d.Get("keypair").(string),
		Labels:            labels,
		MasterFlavorID:    containerInfraClusterV1MasterFlavor(d),
		Name:              d.Get("name").(string),
		FixedNetwork:      d.Get("fixed_network").(string),
		FixedSubnet:       d.Get("fixed_subnet").(string),
	}

	if v, ok := getOkExists(d, "floating_ip_enabled"); ok {
		v := v.(bool)
		createOpts.FloatingIPEnabled = &v
	}

	if v, ok := getOkExists(d, "create_timeout"); ok {
		v := v.(int)
		createOpts.CreateTimeout = &v
	}

	if v, ok := getOkExists(d, "docker_volume_size"); ok {
		v := v.(int)
		createOpts.DockerVolumeSize = &v
	}

	if v, ok := getOkExists(d, "master_count"); ok {
		v := v.(int)
		createOpts.MasterCount = &v
	}

	if v, ok := getOkExists(d, "node_count"); ok {
		v := v.(int)
		createOpts.NodeCount = &v

		if v == 0 {
			containerInfraClient.Microversion = containerInfraV1ZeroNodeCountMicroversion
		}
	}

	if v, ok := getOkExists(d, "merge_labels"); ok {
		v := v.(bool)
		createOpts.MergeLabels = &v
	}

	if v, ok := getOkExists(d, "master_lb_enabled"); ok {
		v := v.(bool)
		createOpts.MasterLBEnabled = &v
	}

	s, err := clusters.Create(ctx, containerInfraClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_containerinfra_cluster_v1: %s", err)
	}

	// Store the Cluster ID.
	d.SetId(s)

	stateConf := &retry.StateChangeConf{
		Pending:      []string{"CREATE_IN_PROGRESS"},
		Target:       []string{"CREATE_COMPLETE"},
		Refresh:      containerInfraClusterV1StateRefreshFunc(ctx, containerInfraClient, s),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        0,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_containerinfra_cluster_v1 %s to become ready: %s", s, err)
	}

	log.Printf("[DEBUG] Created openstack_containerinfra_cluster_v1 %s", s)

	return resourceContainerInfraClusterV1Read(ctx, d, meta)
}

func getDefaultNodegroupNodeCount(ctx context.Context, containerInfraClient *gophercloud.ServiceClient, clusterID string) (int, error) {
	containerInfraClient.Microversion = containerInfraV1NodeGroupMinMicroversion
	listOpts := nodegroups.ListOpts{}

	allPages, err := nodegroups.List(containerInfraClient, clusterID, listOpts).AllPages(ctx)
	if err != nil {
		return 0, err
	}

	ngs, err := nodegroups.ExtractNodeGroups(allPages)
	if err != nil {
		return 0, err
	}

	for _, ng := range ngs {
		if ng.IsDefault && ng.Role != "master" {
			return ng.NodeCount, nil
		}
	}

	return 0, errors.New("Default worker nodegroup not found")
}

func resourceContainerInfraClusterV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	containerInfraClient, err := config.ContainerInfraV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	s, err := clusters.Get(ctx, containerInfraClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_containerinfra_cluster_v1"))
	}

	log.Printf("[DEBUG] Retrieved openstack_containerinfra_cluster_v1 %s: %#v", d.Id(), s)

	labels := s.Labels

	if d.Get("merge_labels").(bool) {
		resourceDataLabels, err := expandContainerInfraV1LabelsMap(d.Get("labels").(map[string]any))
		if err != nil {
			return diag.FromErr(err)
		}

		labels = containerInfraV1GetLabelsMerged(s.LabelsAdded, s.LabelsSkipped, s.LabelsOverridden, s.Labels, resourceDataLabels)
	}

	if err := d.Set("labels", labels); err != nil {
		return diag.Errorf("Unable to set openstack_containerinfra_cluster_v1 labels: %s", err)
	}

	nodeCount, err := getDefaultNodegroupNodeCount(ctx, containerInfraClient, d.Id())
	if err != nil {
		log.Printf("[DEBUG] Can't retrieve node_count of the default worker node group %s: %s", d.Id(), err)

		nodeCount = s.NodeCount
	}

	d.Set("region", GetRegion(d, config))
	d.Set("name", s.Name)
	d.Set("project_id", s.ProjectID)
	d.Set("user_id", s.UserID)
	d.Set("api_address", s.APIAddress)
	d.Set("coe_version", s.COEVersion)
	d.Set("cluster_template_id", s.ClusterTemplateID)
	d.Set("container_version", s.ContainerVersion)
	d.Set("create_timeout", s.CreateTimeout)
	d.Set("discovery_url", s.DiscoveryURL)
	d.Set("docker_volume_size", s.DockerVolumeSize)
	d.Set("flavor", s.FlavorID)
	d.Set("master_flavor", s.MasterFlavorID)
	d.Set("keypair", s.KeyPair)
	d.Set("master_count", s.MasterCount)
	d.Set("master_lb_enabled", s.MasterLBEnabled)
	d.Set("node_count", nodeCount)
	d.Set("master_addresses", s.MasterAddresses)
	d.Set("node_addresses", s.NodeAddresses)
	d.Set("stack_id", s.StackID)
	d.Set("fixed_network", s.FixedNetwork)
	d.Set("fixed_subnet", s.FixedSubnet)
	d.Set("floating_ip_enabled", s.FloatingIPEnabled)

	kubeconfig, err := flattenContainerInfraV1Kubeconfig(ctx, d, containerInfraClient)
	if err != nil {
		return diag.Errorf("Error building kubeconfig for openstack_containerinfra_cluster_v1 %s: %s", d.Id(), err)
	}

	d.Set("kubeconfig", kubeconfig)

	if err := d.Set("created_at", s.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_containerinfra_cluster_v1 created_at: %s", err)
	}

	if err := d.Set("updated_at", s.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_containerinfra_cluster_v1 updated_at: %s", err)
	}

	return nil
}

func resourceContainerInfraClusterV1Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	containerInfraClient, err := config.ContainerInfraV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	if d.HasChange("cluster_template_id") {
		clusterTemplateID := d.Get("cluster_template_id").(string)
		upgradeOpts := clusters.UpgradeOpts{
			ClusterTemplate: clusterTemplateID,
		}

		containerInfraClient.Microversion = containerInfraV1ClusterUpgradeMinMicroversion

		_, err = clusters.Upgrade(ctx, containerInfraClient, d.Id(), upgradeOpts).Extract()
		if err != nil {
			return diag.Errorf("Error upgrading openstack_containerinfra_cluster_v1 %s: %s", d.Id(), err)
		}

		stateConf := &retry.StateChangeConf{
			Target:       []string{"UPDATE_COMPLETE"},
			Refresh:      containerInfraClusterV1StateRefreshFunc(ctx, containerInfraClient, d.Id()),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        0,
			PollInterval: 20 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"Error waiting for openstack_containerinfra_cluster_v1 %s to upgrade: %s", d.Id(), err)
		}
	}

	updateOpts := []clusters.UpdateOptsBuilder{}

	if d.HasChange("node_count") {
		nodeCount := d.Get("node_count").(int)

		if nodeCount == 0 {
			containerInfraClient.Microversion = containerInfraV1ZeroNodeCountMicroversion
		}

		updateOpts = append(updateOpts, clusters.UpdateOpts{
			Op:    clusters.ReplaceOp,
			Path:  strings.Join([]string{"/", "node_count"}, ""),
			Value: nodeCount,
		})
	}

	if len(updateOpts) > 0 {
		log.Printf(
			"[DEBUG] Updating openstack_containerinfra_cluster_v1 %s with options: %#v", d.Id(), updateOpts)

		_, err = clusters.Update(ctx, containerInfraClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_containerinfra_cluster_v1 %s: %s", d.Id(), err)
		}

		stateConf := &retry.StateChangeConf{
			Target:       []string{"UPDATE_COMPLETE"},
			Refresh:      containerInfraClusterV1StateRefreshFunc(ctx, containerInfraClient, d.Id()),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        0,
			PollInterval: 20 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"Error waiting for openstack_containerinfra_cluster_v1 %s to become updated: %s", d.Id(), err)
		}
	}

	return resourceContainerInfraClusterV1Read(ctx, d, meta)
}

func resourceContainerInfraClusterV1Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	containerInfraClient, err := config.ContainerInfraV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	if err := clusters.Delete(ctx, containerInfraClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_containerinfra_cluster_v1"))
	}

	stateConf := &retry.StateChangeConf{
		Target:       []string{"DELETE_COMPLETE"},
		Refresh:      containerInfraClusterV1StateRefreshFunc(ctx, containerInfraClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        0,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_containerinfra_cluster_v1 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}
