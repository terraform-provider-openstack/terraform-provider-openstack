package clientconfig

// Clouds represents a collection of Cloud entries in a clouds.yaml file.
// The format of clouds.yaml is documented at
// https://docs.openstack.org/os-client-config/latest/user/configuration.html.
type Clouds struct {
	Clouds map[string]Cloud `yaml:"clouds"`
}

// Cloud represents an entry in a clouds.yaml file.
type Cloud struct {
	Auth       *CloudAuth    `yaml:"auth"`
	RegionName string        `yaml:"region_name"`
	Regions    []interface{} `yaml:"regions"`

	// API Version overrides.
	IdentityAPIVersion string `yaml:"identity_api_version"`
	VolumeAPIVersion   string `yaml:"volume_api_version"`
}

// CloudAuth represents the auth section of a cloud entry.
type CloudAuth struct {
	AuthURL           string `yaml:"auth_url"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	ProjectName       string `yaml:"project_name"`
	ProjectID         string `yaml:"project_id"`
	DomainName        string `yaml:"domain_name"`
	DomainID          string `yaml:"domain_id"`
	UserDomainName    string `yaml:"user_domain_name"`
	UserDomainID      string `yaml:"user_domain_id"`
	ProjectDomainName string `yaml:"project_domain_name"`
	ProjectDomainID   string `yaml:"project_domain_id"`
}
