package openstack

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceContainerInfraClusterTemplateV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerInfraClusterTemplateV1Create,
		Read:   resourceContainerInfraClusterTemplateV1Read,
		// Update: resourceContainerInfraClusterTemplateV1Update,
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
				Computed:     true,
				ValidateFunc: validateClusterTemplateAPIServerPortV1,
			},
			"coe": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"cluster_distro": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},
			"dns_nameserver": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"docker_storage_driver": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"docker_volume_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"external_network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"fixed_network": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"fixed_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"flavor_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_FLAVOR_ID", nil),
			},
			"flavor_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_FLAVOR_NAME", nil),
			},
			"master_flavor_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_FLAVOR_ID", nil),
			},
			"master_flavor_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_FLAVOR_NAME", nil),
			},
			"floating_ip_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"http_proxy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"https_proxy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"image_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CONTAINERINFRA_IMAGE_ID", nil),
			},
			"image_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CONTAINERINFRA_IMAGE_NAME", nil),
			},
			"insecure_registry": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"keypair_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
				Computed: true,
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
				Computed: true,
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
				Computed: true,
			},
			"public": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"registry_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
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
				Computed: true,
			},
			"volume_driver": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
		},
	}
}

func resourceContainerInfraClusterTemplateV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}
	containerInfraClient, err := config.containerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	// Get imageID using following rules:
	//  - if an image_id was specified, use it
	//  - if an image_name was specified, look up the ID, report if error
	imageID, err := getContainerInfraImageID(computeClient, d)
	if err != nil {
		return err
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
		FloatingIPEnabled:   &floatingIPEnabled,
		HTTPProxy:           d.Get("http_proxy").(string),
		HTTPSProxy:          d.Get("https_proxy").(string),
		ImageID:             imageID,
		InsecureRegistry:    d.Get("insecure_registry").(string),
		KeyPairID:           d.Get("keypair_id").(string),
		Labels:              resourceClusterTemplateLabelsV1(d),
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

	// Get nodes and masters flavors using following rules:
	//  - if a flavor_id was specified, use it for regular nodes
	//  - if a flavor_name was specified, look up the ID, report if error
	//  - if a master_flavor_id was specified, use it for regular nodes
	//  - if a master_flavor_name was specified, look up the ID, report if error
	flavorID, err := getContainerInfraFlavorID(computeClient, d, "")
	if err != nil {
		return err
	}
	if flavorID != "" {
		createOpts.FlavorID = flavorID
	}
	masterFlavorID, err := getContainerInfraFlavorID(computeClient, d, "master")
	if err != nil {
		return err
	}
	if masterFlavorID != "" {
		createOpts.MasterFlavorID = masterFlavorID
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
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}
	containerInfraClient, err := config.containerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	s, err := clustertemplates.Get(containerInfraClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "clustertemplate")
	}

	log.Printf("[DEBUG] Retrieved Clustertemplate %s: %#v", d.Id(), s)

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

	// Set flavors names.
	if s.FlavorID != "" {
		flavor, err := flavors.Get(computeClient, s.FlavorID).Extract()
		if err != nil {
			return err
		}
		d.Set("flavor_name", flavor.Name)
	}
	if s.MasterFlavorID != "" {
		masterFlavor, err := flavors.Get(computeClient, s.MasterFlavorID).Extract()
		if err != nil {
			return err
		}
		d.Set("master_flavor_name", masterFlavor.Name)
	}

	// Set apiserver_port.
	if s.APIServerPort != "" {
		apiServerPort, err := strconv.Atoi(s.APIServerPort)
		if err != nil {
			return fmt.Errorf("Error setting Cluster template API server port: %v", s.APIServerPort)
		}
		d.Set("apiserver_port", apiServerPort)
	}

	return nil
}

// func resourceContainerInfraClusterTemplateV1Update(d *schema.ResourceData, meta interface{}) error {
//
// }

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
	if value < 1024 {
		err := fmt.Errorf("%s should be greater or equal to 1024", k)
		errors = append(errors, err)
	}
	return
}

func getContainerInfraImageID(computeClient *gophercloud.ServiceClient, d *schema.ResourceData) (string, error) {
	if imageID := d.Get("image_id").(string); imageID != "" {
		return imageID, nil
	}

	if imageName := d.Get("image_name").(string); imageName != "" {
		imageID, err := images.IDFromName(computeClient, imageName)
		if err != nil {
			return "", err
		}
		return imageID, nil
	}

	return "", fmt.Errorf("Neither an image_id or image_name were able to be determined.")
}

func getContainerInfraFlavorID(client *gophercloud.ServiceClient, d *schema.ResourceData, prefix string) (string, error) {
	attrID := "flavor_id"
	attrName := "flavor_name"
	if prefix != "" {
		attrID = strings.Join([]string{prefix, attrID}, "_")
		attrName = strings.Join([]string{prefix, attrName}, "_")
	}

	flavorID := d.Get("flavor_id").(string)
	if flavorID != "" {
		return flavorID, nil
	}

	flavorName := d.Get("flavor_name").(string)
	if flavorName == "" {
		return "", nil
	}

	return flavors.IDFromName(client, flavorName)
}

func resourceClusterTemplateLabelsV1(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
