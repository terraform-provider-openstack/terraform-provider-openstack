package openstack

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceContainerInfraClusterTemplateV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerInfraClusterTemplateV1Create,
		Read:   resourceContainerInfraClusterTemplateV1Read,
		Update: resourceContainerInfraClusterTemplateV1Update,
		Delete: resourceContainerInfraClusterTemplateV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"apiserver_port": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validateClusterTemplateAPIServerPortV1,
			},
			"coe": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"cluster_distro": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"dns_nameserver": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"docker_storage_driver": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"docker_volume_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"external_network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"fixed_network": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"fixed_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"flavor": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				DefaultFunc: schema.EnvDefaultFunc("OS_MAGNUM_FLAVOR", nil),
			},
			"master_flavor": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				DefaultFunc: schema.EnvDefaultFunc("OS_MAGNUM_MASTER_FLAVOR", nil),
			},
			"floating_ip_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"http_proxy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"https_proxy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"image": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				DefaultFunc: schema.EnvDefaultFunc("OS_MAGNUM_IMAGE", nil),
			},
			"insecure_registry": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"keypair_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"links": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: false,
							Computed: true,
						},
						"rel": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: false,
							Computed: true,
						},
					},
				},
				ForceNew: false,
				Computed: true,
			},
			"master_lb_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"network_driver": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"no_proxy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"public": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"registry_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"server_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"tls_disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"volume_driver": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceContainerInfraClusterTemplateV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerInfraClient, err := config.containerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	// Get boolean parameters that will be passed by reference.
	floatingIPEnabled := d.Get("floating_ip_enabled").(bool)
	masterLBEnabled := d.Get("master_lb_enabled").(bool)
	public := d.Get("public").(bool)
	registryEnabled := d.Get("registry_enabled").(bool)
	tlsDisabled := d.Get("tls_disabled").(bool)

	createOpts := clustertemplates.CreateOpts{
		COE:                 d.Get("coe").(string),
		DNSNameServer:       d.Get("dns_nameserver").(string),
		DockerStorageDriver: d.Get("docker_storage_driver").(string),
		ExternalNetworkID:   d.Get("external_network_id").(string),
		FixedNetwork:        d.Get("fixed_network").(string),
		FixedSubnet:         d.Get("fixed_subnet").(string),
		FlavorID:            d.Get("flavor").(string),
		MasterFlavorID:      d.Get("master_flavor").(string),
		FloatingIPEnabled:   &floatingIPEnabled,
		HTTPProxy:           d.Get("http_proxy").(string),
		HTTPSProxy:          d.Get("https_proxy").(string),
		ImageID:             d.Get("image").(string),
		InsecureRegistry:    d.Get("insecure_registry").(string),
		KeyPairID:           d.Get("keypair_id").(string),
		Labels:              resourceClusterTemplateLabelsMapV1(d),
		MasterLBEnabled:     &masterLBEnabled,
		Name:                d.Get("name").(string),
		NetworkDriver:       d.Get("network_driver").(string),
		NoProxy:             d.Get("no_proxy").(string),
		Public:              &public,
		RegistryEnabled:     &registryEnabled,
		ServerType:          d.Get("server_type").(string),
		TLSDisabled:         &tlsDisabled,
		VolumeDriver:        d.Get("volume_driver").(string),
	}

	// Set int parameters that will be passed by reference.
	apiServerPort := d.Get("apiserver_port").(int)
	if apiServerPort > 0 {
		createOpts.APIServerPort = &apiServerPort
	}
	dockerVolumeSize := d.Get("docker_volume_size").(int)
	if dockerVolumeSize > 0 {
		createOpts.DockerVolumeSize = &dockerVolumeSize
	}

	s, err := clustertemplates.Create(containerInfraClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra Cluster template: %s", err)
	}

	d.SetId(s.UUID)

	log.Printf("[DEBUG] Created Cluster template %s: %#v", s.UUID, s)
	return resourceContainerInfraClusterTemplateV1Read(d, meta)
}

func resourceContainerInfraClusterTemplateV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerInfraClient, err := config.containerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	s, err := clustertemplates.Get(containerInfraClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "clustertemplate")
	}

	log.Printf("[DEBUG] Retrieved Clustertemplate %s: %#v", d.Id(), s)

	d.Set("apiserver_port", s.APIServerPort)
	d.Set("coe", s.COE)
	d.Set("cluster_distro", s.ClusterDistro)
	d.Set("dns_nameserver", s.DNSNameServer)
	d.Set("docker_storage_driver", s.DockerStorageDriver)
	d.Set("docker_volume_size", s.DockerVolumeSize)
	d.Set("external_network_id", s.ExternalNetworkID)
	d.Set("fixed_network", s.FixedNetwork)
	d.Set("fixed_subnet", s.FixedSubnet)
	d.Set("flavor_id", s.FlavorID)
	d.Set("master_flavor_id", s.MasterFlavorID)
	d.Set("floating_ip_enabled", s.FloatingIPEnabled)
	d.Set("http_proxy", s.HTTPProxy)
	d.Set("https_proxy", s.HTTPSProxy)
	d.Set("image_id", s.ImageID)
	d.Set("insecure_registry", s.InsecureRegistry)
	d.Set("keypair_id", s.KeyPairID)
	d.Set("labels", s.Labels)
	d.Set("links", resourceLinks(s.Links))
	d.Set("master_lb_enabled", s.MasterLBEnabled)
	d.Set("network_driver", s.NetworkDriver)
	d.Set("no_proxy", s.NoProxy)
	d.Set("public", s.Public)
	d.Set("registry_enabled", s.RegistryEnabled)
	d.Set("server_type", s.ServerType)
	d.Set("tls_disabled", s.TLSDisabled)
	d.Set("volume_driver", s.VolumeDriver)
	d.Set("region", GetRegion(d, config))
	d.Set("name", s.Name)
	d.Set("project_id", s.ProjectID)
	d.Set("user_id", s.UserID)
	d.Set("created_at", s.CreatedAt)
	d.Set("updated_at", s.UpdatedAt)

	return nil
}

func resourceContainerInfraClusterTemplateV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerInfraClient, err := config.containerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	updateOpts := []clustertemplates.UpdateOptsBuilder{}

	if d.HasChange("name") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "name", d.Get("name").(string))
	}
	if d.HasChange("apiserver_port") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "apiserver_port", strconv.Itoa(d.Get("apiserver_port").(int)))
	}
	if d.HasChange("coe") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "coe", d.Get("coe").(string))
	}
	if d.HasChange("cluster_distro") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "cluster_distro", d.Get("cluster_distro").(string))
	}
	if d.HasChange("dns_nameserver") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "dns_nameserver", d.Get("dns_nameserver").(string))
	}
	if d.HasChange("docker_storage_driver") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "docker_storage_driver", d.Get("docker_storage_driver").(string))
	}
	if d.HasChange("docker_volume_size") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "docker_volume_size", strconv.Itoa(d.Get("docker_volume_size").(int)))
	}
	if d.HasChange("external_network_id") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "external_network_id", d.Get("external_network_id").(string))
	}
	if d.HasChange("fixed_network") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "fixed_network", d.Get("fixed_network").(string))
	}
	if d.HasChange("fixed_subnet") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "fixed_subnet", d.Get("fixed_subnet").(string))
	}
	if d.HasChange("flavor") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "flavor_id", d.Get("flavor").(string))
	}
	if d.HasChange("master_flavor") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "master_flavor_id", d.Get("master_flavor").(string))
	}
	if d.HasChange("floating_ip_enabled") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "floating_ip_enabled", strconv.FormatBool(d.Get("floating_ip_enabled").(bool)))
	}
	if d.HasChange("http_proxy") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "http_proxy", d.Get("http_proxy").(string))
	}
	if d.HasChange("https_proxy") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "https_proxy", d.Get("https_proxy").(string))
	}
	if d.HasChange("image") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "image_id", d.Get("image").(string))
	}
	if d.HasChange("insecure_registry") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "insecure_registry", d.Get("insecure_registry").(string))
	}
	if d.HasChange("keypair_id") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "keypair_id", d.Get("keypair_id").(string))
	}
	if d.HasChange("labels") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "labels", resourceClusterTemplateLabelsStringV1(d.Get("labels").(map[string]interface{})))
	}
	if d.HasChange("master_lb_enabled") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "master_lb_enabled", strconv.FormatBool(d.Get("master_lb_enabled").(bool)))
	}
	if d.HasChange("network_driver") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "network_driver", d.Get("network_driver").(string))
	}
	if d.HasChange("no_proxy") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "no_proxy", d.Get("no_proxy").(string))
	}
	if d.HasChange("public") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "public", strconv.FormatBool(d.Get("public").(bool)))
	}
	if d.HasChange("registry_enabled") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "registry_enabled", strconv.FormatBool(d.Get("registry_enabled").(bool)))
	}
	if d.HasChange("server_type") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "server_type", d.Get("server_type").(string))
	}
	if d.HasChange("tls_disabled") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "tls_disabled", strconv.FormatBool(d.Get("tls_disabled").(bool)))
	}
	if d.HasChange("volume_driver") {
		updateOpts = resourceClusterTemplateAppendUpdateOptsV1(updateOpts, "volume_driver", d.Get("volume_driver").(string))
	}

	log.Printf("[DEBUG] Updating Cluster template %s with options: %+v", d.Id(), updateOpts)

	_, err = clustertemplates.Update(containerInfraClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenStack container infra Cluster template: %s", err)
	}

	return resourceContainerInfraClusterTemplateV1Read(d, meta)
}

func resourceContainerInfraClusterTemplateV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerInfraClient, err := config.containerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	if err := clustertemplates.Delete(containerInfraClient, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf("Error deleting Cluster template: %v", err)
	}

	return nil
}

func validateClusterTemplateAPIServerPortV1(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1024 || value > 65535 {
		err := fmt.Errorf("%s should be between 1024 and 65535", k)
		errors = append(errors, err)
	}
	return
}

func resourceClusterTemplateLabelsMapV1(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceClusterTemplateLabelsStringV1(labels map[string]interface{}) string {
	var formattedLabels string
	for labelKey, labelValue := range labels {
		formattedLabels = strings.Join([]string{
			formattedLabels,
			fmt.Sprintf("%s=%s", labelKey, labelValue.(string)),
		}, ",")
	}
	formattedLabels = strings.Trim(formattedLabels, ",")

	return formattedLabels
}

func resourceClusterTemplateAppendUpdateOptsV1(updateOpts []clustertemplates.UpdateOptsBuilder, attribute string, value string) []clustertemplates.UpdateOptsBuilder {
	if value == "" {
		updateOpts = append(updateOpts, clustertemplates.UpdateOpts{
			Op:   clustertemplates.RemoveOp,
			Path: strings.Join([]string{"/", attribute}, ""),
		})
	} else {
		updateOpts = append(updateOpts, clustertemplates.UpdateOpts{
			Op:    clustertemplates.ReplaceOp,
			Path:  strings.Join([]string{"/", attribute}, ""),
			Value: value,
		})
	}
	return updateOpts
}
