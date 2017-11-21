package clientconfig

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"

	"gopkg.in/yaml.v2"
)

// ClientOpts represents options to customize the way a client is
// configured.
type ClientOpts struct {
	// Cloud is the cloud entry in clouds.yaml to use.
	Cloud string

	// EnvPrefix allows a custom environment variable prefix to be used.
	EnvPrefix string
}

// LoadYAML will load a clouds.yaml file and return the full config.
func LoadYAML() (map[string]Cloud, error) {
	content, err := findAndReadYAML()
	if err != nil {
		return nil, err
	}

	var clouds Clouds
	err = yaml.Unmarshal(content, &clouds)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	return clouds.Clouds, nil
}

// GetCloudFromYAML will return a cloud entry from a clouds.yaml file.
func GetCloudFromYAML(opts *ClientOpts) (*Cloud, error) {
	clouds, err := LoadYAML()
	if err != nil {
		return nil, fmt.Errorf("unable to load clouds.yaml: %s", err)
	}

	// Determine which cloud to use.
	var cloudName string
	if opts != nil && opts.Cloud != "" {
		cloudName = opts.Cloud
	}

	if v := os.Getenv("OS_CLOUD"); v != "" {
		cloudName = v
	}

	var cloud *Cloud
	if cloudName != "" {
		v, ok := clouds[cloudName]
		if !ok {
			return nil, fmt.Errorf("cloud %s does not exist in clouds.yaml", cloudName)
		}
		cloud = &v
	}

	// If a cloud was not specified, and clouds only contains
	// a single entry, use that entry.
	if cloudName == "" && len(clouds) == 1 {
		for _, v := range clouds {
			cloud = &v
		}
	}

	if cloud == nil {
		return nil, fmt.Errorf("Unable to determine a valid entry in clouds.yaml")
	}

	return cloud, nil
}

// AuthOptions creates a gophercloud identity.AuthOptions structure with the
// settings found in a specific cloud entry of a clouds.yaml file.
// See http://docs.openstack.org/developer/os-client-config and
// https://github.com/openstack/os-client-config/blob/master/os_client_config/config.py.
func AuthOptions(opts *ClientOpts) (*gophercloud.AuthOptions, error) {
	// Get the requested cloud.
	cloud, err := GetCloudFromYAML(opts)
	if err != nil {
		return nil, err
	}

	auth := cloud.Auth

	// Create a gophercloud.AuthOptions struct based on the clouds.yaml entry.
	ao := &gophercloud.AuthOptions{
		IdentityEndpoint: auth.AuthURL,
		Username:         auth.Username,
		Password:         auth.Password,
		TenantID:         auth.ProjectID,
		TenantName:       auth.ProjectName,
		DomainID:         auth.DomainID,
		DomainName:       auth.DomainName,
	}

	// Domain scope overrides.
	if auth.ProjectDomainID != "" {
		ao.DomainID = auth.ProjectDomainID
	}

	if auth.ProjectDomainName != "" {
		ao.DomainName = auth.ProjectDomainName
	}

	if auth.UserDomainID != "" {
		ao.DomainID = auth.UserDomainID
	}

	if auth.UserDomainName != "" {
		ao.DomainName = auth.UserDomainName
	}

	// Environment variable overrides.
	envPrefix := "OS_"
	if opts != nil && opts.EnvPrefix != "" {
		envPrefix = opts.EnvPrefix
	}

	if v := os.Getenv(envPrefix + "AUTH_URL"); v != "" {
		ao.IdentityEndpoint = v
	}

	if v := os.Getenv(envPrefix + "USERNAME"); v != "" {
		ao.Username = v
	}

	if v := os.Getenv(envPrefix + "PASSWORD"); v != "" {
		ao.Password = v
	}

	if v := os.Getenv(envPrefix + "TENANT_ID"); v != "" {
		ao.TenantID = v
	}

	if v := os.Getenv(envPrefix + "PROJECT_ID"); v != "" {
		ao.TenantID = v
	}

	if v := os.Getenv(envPrefix + "TENANT_NAME"); v != "" {
		ao.TenantName = v
	}

	if v := os.Getenv(envPrefix + "PROJECT_NAME"); v != "" {
		ao.TenantName = v
	}

	if v := os.Getenv(envPrefix + "DOMAIN_ID"); v != "" {
		ao.DomainID = v
	}

	if v := os.Getenv(envPrefix + "PROJECT_DOMAIN_ID"); v != "" {
		ao.DomainID = v
	}

	if v := os.Getenv(envPrefix + "DOMAIN_NAME"); v != "" {
		ao.DomainName = v
	}

	if v := os.Getenv(envPrefix + "PROJECT_DOMAIN_NAME"); v != "" {
		ao.DomainName = v
	}

	// Check for absolute minimum requirements.
	if ao.IdentityEndpoint == "" {
		err := gophercloud.ErrMissingInput{Argument: "authURL"}
		return nil, err
	}

	if ao.Username == "" {
		err := gophercloud.ErrMissingInput{Argument: "username"}
		return nil, err
	}

	if ao.Password == "" {
		err := gophercloud.ErrMissingInput{Argument: "password"}
		return nil, err
	}

	return ao, nil
}

// AuthenticatedClient is a convenience function to get a new provider client
// based on a clouds.yaml entry.
func AuthenticatedClient(opts *ClientOpts) (*gophercloud.ProviderClient, error) {
	ao, err := AuthOptions(opts)
	if err != nil {
		return nil, err
	}

	return openstack.AuthenticatedClient(*ao)
}

// NewServiceClient is a convenience function to get a new service client.
func NewServiceClient(service string, opts *ClientOpts) (*gophercloud.ServiceClient, error) {
	cloud, err := GetCloudFromYAML(opts)
	if err != nil {
		return nil, err
	}

	// Environment variable overrides.
	envPrefix := "OS_"
	if opts != nil && opts.EnvPrefix != "" {
		envPrefix = opts.EnvPrefix
	}

	// Get a Provider Client
	pClient, err := AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	// Determine the region to use.
	var region string
	if v := cloud.RegionName; v != "" {
		region = cloud.RegionName
	}

	if v := os.Getenv(envPrefix + "REGION_NAME"); v != "" {
		region = v
	}

	eo := gophercloud.EndpointOpts{
		Region: region,
	}

	switch service {
	case "compute":
		return openstack.NewComputeV2(pClient, eo)
	case "database":
		return openstack.NewDBV1(pClient, eo)
	case "dns":
		return openstack.NewDNSV2(pClient, eo)
	case "identity":
		identityVersion := "3"
		if v := cloud.IdentityAPIVersion; v != "" {
			identityVersion = v
		}

		switch identityVersion {
		case "v2", "2", "2.0":
			return openstack.NewIdentityV2(pClient, eo)
		case "v3", "3":
			return openstack.NewIdentityV3(pClient, eo)
		default:
			return nil, fmt.Errorf("invalid identity API version")
		}
	case "image":
		return openstack.NewImageServiceV2(pClient, eo)
	case "network":
		return openstack.NewNetworkV2(pClient, eo)
	case "object-store":
		return openstack.NewObjectStorageV1(pClient, eo)
	case "orchestration":
		return openstack.NewOrchestrationV1(pClient, eo)
	case "sharev2":
		return openstack.NewSharedFileSystemV2(pClient, eo)
	case "volume":
		volumeVersion := "2"
		if v := cloud.VolumeAPIVersion; v != "" {
			volumeVersion = v
		}

		switch volumeVersion {
		case "v1", "1":
			return openstack.NewBlockStorageV1(pClient, eo)
		case "v2", "2":
			return openstack.NewBlockStorageV2(pClient, eo)
		default:
			return nil, fmt.Errorf("invalid volume API version")
		}
	}

	return nil, fmt.Errorf("unable to create a service client for %s", service)
}
