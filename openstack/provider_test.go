package openstack

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/utils/v2/terraform/auth"
	"github.com/gophercloud/utils/v2/terraform/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/v3/openstack/internal/pathorcontents"
)

var (
	osBackupID                   = os.Getenv("OS_BACKUP_ID")
	osDBEnvironment              = os.Getenv("OS_DB_ENVIRONMENT")
	osDBDatastoreVersion         = os.Getenv("OS_DB_DATASTORE_VERSION")
	osDBDatastoreType            = os.Getenv("OS_DB_DATASTORE_TYPE")
	osDNSEnvironment             = os.Getenv("OS_DNS_ENVIRONMENT")
	osExtGwID                    = os.Getenv("OS_EXTGW_ID")
	osFlavorID                   = os.Getenv("OS_FLAVOR_ID")
	osFlavorName                 = os.Getenv("OS_FLAVOR_NAME")
	osImageID                    = os.Getenv("OS_IMAGE_ID")
	osImageName                  = os.Getenv("OS_IMAGE_NAME")
	osMagnumFlavor               = os.Getenv("OS_MAGNUM_FLAVOR")
	osMagnumImage                = os.Getenv("OS_MAGNUM_IMAGE")
	osNetworkID                  = os.Getenv("OS_NETWORK_ID")
	osPoolName                   = os.Getenv("OS_POOL_NAME")
	osRegionName                 = os.Getenv("OS_REGION_NAME")
	osSwiftEnvironment           = os.Getenv("OS_SWIFT_ENVIRONMENT")
	osLbEnvironment              = os.Getenv("OS_LB_ENVIRONMENT")
	osLbFlavorName               = os.Getenv("OS_LB_FLAVOR_NAME")
	osFwEnvironment              = os.Getenv("OS_FW_ENVIRONMENT")
	osVpnEnvironment             = os.Getenv("OS_VPN_ENVIRONMENT")
	osBgpVpnEnvironment          = os.Getenv("OS_BGPVPN_ENVIRONMENT")
	osContainerInfraEnvironment  = os.Getenv("OS_CONTAINER_INFRA_ENVIRONMENT")
	osSfsEnvironment             = os.Getenv("OS_SFS_ENVIRONMENT")
	osTransparentVlanEnvironment = os.Getenv("OS_TRANSPARENT_VLAN_ENVIRONMENT")
	osKeymanagerEnvironment      = os.Getenv("OS_KEYMANAGER_ENVIRONMENT")
	osGlanceimportEnvironment    = os.Getenv("OS_GLANCEIMPORT_ENVIRONMENT")
	osHypervisorEnvironment      = os.Getenv("OS_HYPERVISOR_HOSTNAME")
	osPortForwardingEnvironment  = os.Getenv("OS_PORT_FORWARDING_ENVIRONMENT")
	osWorkflowEnvironment        = os.Getenv("OS_WORKFLOW_ENVIRONMENT")
	osMagnumHTTPProxy            = os.Getenv("OS_MAGNUM_HTTP_PROXY")
	osMagnumHTTPSProxy           = os.Getenv("OS_MAGNUM_HTTPS_PROXY")
	osMagnumNoProxy              = os.Getenv("OS_MAGNUM_NO_PROXY")
	osMagnumLabels               = os.Getenv("OS_MAGNUM_LABELS")
)

var (
	testAccProviders map[string]func() (*schema.Provider, error)
	testAccProvider  *schema.Provider
)

func init() {
	// ensure that only "sharev2" service type endpoint is used for shared-file-system
	// see https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1847
	gophercloud.ServiceTypeAliases["shared-file-system"] = []string{"sharev2"}

	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"openstack": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	v := os.Getenv("OS_AUTH_URL")
	if v == "" {
		t.Fatal("OS_AUTH_URL must be set for acceptance tests")
	}

	if osImageID == "" && osImageName == "" {
		t.Fatal("OS_IMAGE_ID or OS_IMAGE_NAME must be set for acceptance tests")
	}

	if osPoolName == "" {
		t.Fatal("OS_POOL_NAME must be set for acceptance tests")
	}

	if osFlavorID == "" && osFlavorName == "" {
		t.Fatal("OS_FLAVOR_ID or OS_FLAVOR_NAME must be set for acceptance tests")
	}

	if osNetworkID == "" {
		t.Fatal("OS_NETWORK_ID must be set for acceptance tests")
	}

	if osExtGwID == "" {
		t.Fatal("OS_EXTGW_ID must be set for acceptance tests")
	}
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
}

func testAccPreCheckDNS(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osDNSEnvironment == "" {
		t.Skip("This environment does not support DNS tests")
	}
}

func testAccPreCheckSwift(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osSwiftEnvironment == "" {
		t.Skip("This environment does not support Swift tests")
	}
}

func testAccPreCheckDatabase(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osDBEnvironment == "" {
		t.Skip("This environment does not support Database tests")
	}
}

func testAccPreCheckLB(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osLbEnvironment == "" {
		t.Skip("This environment does not support LB tests")
	}

	if osLbFlavorName == "" {
		t.Skip("This environment does not support LB tests")
	}
}

func testAccPreCheckFW(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osFwEnvironment == "" {
		t.Skip("This environment does not support FW tests")
	}
}

func testAccPreCheckVPN(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osVpnEnvironment == "" {
		t.Skip("This environment does not support VPN tests")
	}
}

func testAccPreCheckBGPVPN(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osBgpVpnEnvironment == "" {
		t.Skip("This environment does not support BGP VPN tests")
	}
}

func testAccPreCheckKeyManager(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osKeymanagerEnvironment == "" {
		t.Skip("This environment does not support Barbican Keymanager tests")
	}
}

func testAccPreCheckContainerInfra(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osContainerInfraEnvironment == "" {
		t.Skip("This environment does not support Container Infra tests")
	}

	if osMagnumImage == "" {
		t.Fatal("OS_MAGNUM_IMAGE must be set for acceptance tests")
	}

	if osMagnumFlavor == "" {
		t.Fatal("OS_MAGNUM_FLAVOR must be set for acceptance tests")
	}
}

func testAccPreCheckSFS(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osSfsEnvironment == "" {
		t.Skip("This environment does not support Shared File Systems tests")
	}
}

func testAccPreOnlineResize(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	v := os.Getenv("OS_ONLINE_RESIZE")
	if v == "" {
		t.Skip("This environment does not support online blockstorage resize tests")
	}

	v = os.Getenv("OS_FLAVOR_NAME")
	if v == "" {
		t.Skip("OS_FLAVOR_NAME required to support online blockstorage resize tests")
	}
}

func testAccPreCheckTransparentVLAN(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osTransparentVlanEnvironment == "" {
		t.Skip("This environment does not support 'transparent-vlan' extension tests")
	}
}

func testAccPreCheckPortForwarding(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osPortForwardingEnvironment == "" {
		t.Skip("This environment does not support 'portforwarding' extension tests")
	}
}

func testAccPreCheckWorkflow(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if osWorkflowEnvironment == "" {
		t.Skip("This environment does not support 'workflow' extension tests")
	}
}

func testAccPreCheckAdminOnly(t *testing.T) {
	v := os.Getenv("OS_USERNAME")
	if v != "admin" {
		t.Skip("Skipping test because it requires the admin user")
	}
}

func testAccPreCheckNonAdminOnly(t *testing.T) {
	v := os.Getenv("OS_USERNAME")
	if v != "demo" {
		t.Skip("Skipping test because it requires the demo (non-admin) user")
	}
}

func testAccPreCheckGlanceImport(t *testing.T) {
	if osGlanceimportEnvironment == "" {
		t.Skip("This environment does not support Glance import tests")
	}
}

func testAccPreCheckHypervisor(t *testing.T) {
	if osHypervisorEnvironment == "" {
		t.Skip("This environment does not support Hypervisor data source tests")
	}
}

func TestUnitProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// Steps for configuring OpenStack with SSL validation are here:
// https://github.com/hashicorp/terraform/pull/6279#issuecomment-219020144
func TestAccProvider_caCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping OpenStack SSL test.")
	}

	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping OpenStack CA test.")
	}

	p := Provider()

	caFile, err := envVarFile("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(caFile)

	raw := map[string]any{
		"cacert_file": caFile,
	}

	diag := p.Configure(t.Context(), terraform.NewResourceConfigRaw(raw))
	if diag != nil {
		t.Fatalf("Unexpected err when specifying OpenStack CA by file: %v", diag)
	}
}

func TestAccProvider_caCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping OpenStack SSL test.")
	}

	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping OpenStack CA test.")
	}

	p := Provider()

	caContents, err := envVarContents("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]any{
		"cacert_file": caContents,
	}

	diag := p.Configure(t.Context(), terraform.NewResourceConfigRaw(raw))
	if diag != nil {
		t.Fatalf("Unexpected err when specifying OpenStack CA by string: %v", diag)
	}
}

func TestAccProvider_clientCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping OpenStack SSL test.")
	}

	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping OpenStack client SSL auth test.")
	}

	p := Provider()

	certFile, err := envVarFile("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(certFile)

	keyFile, err := envVarFile("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(keyFile)

	raw := map[string]any{
		"cert": certFile,
		"key":  keyFile,
	}

	diag := p.Configure(t.Context(), terraform.NewResourceConfigRaw(raw))
	if diag != nil {
		t.Fatalf("Unexpected err when specifying OpenStack Client keypair by file: %v", diag)
	}
}

func TestAccProvider_clientCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping OpenStack SSL test.")
	}

	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping OpenStack client SSL auth test.")
	}

	p := Provider()

	certContents, err := envVarContents("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}

	keyContents, err := envVarContents("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]any{
		"cert": certContents,
		"key":  keyContents,
	}

	diag := p.Configure(t.Context(), terraform.NewResourceConfigRaw(raw))
	if diag != nil {
		t.Fatalf("Unexpected err when specifying OpenStack Client keypair by contents: %v", diag)
	}
}

func envVarContents(varName string) (string, error) {
	contents, _, err := pathorcontents.Read(os.Getenv(varName))
	if err != nil {
		return "", fmt.Errorf("Error reading %s: %w", varName, err)
	}

	return contents, nil
}

func envVarFile(varName string) (string, error) {
	contents, err := envVarContents(varName)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", varName)
	if err != nil {
		return "", fmt.Errorf("Error creating temp file: %w", err)
	}

	if _, err := tmpFile.WriteString(contents); err != nil {
		_ = os.Remove(tmpFile.Name())

		return "", fmt.Errorf("Error writing temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())

		return "", fmt.Errorf("Error closing temp file: %w", err)
	}

	return tmpFile.Name(), nil
}

func testAccAuthFromEnv(ctx context.Context) (*Config, error) {
	// If the provider has already been configured, return the existing config.
	if v, ok := testAccProvider.Meta().(*Config); ok {
		return v, nil
	}

	tenantID := os.Getenv("OS_TENANT_ID")
	if tenantID == "" {
		tenantID = os.Getenv("OS_PROJECT_ID")
	}

	tenantName := os.Getenv("OS_TENANT_NAME")
	if tenantName == "" {
		tenantName = os.Getenv("OS_PROJECT_NAME")
	}

	authOpts := &gophercloud.AuthOptions{
		Scope: &gophercloud.AuthScope{System: testGetenvBool("OS_SYSTEM_SCOPE")},
	}

	config := Config{
		auth.Config{
			CACertFile:                  os.Getenv("OS_CACERT"),
			ClientCertFile:              os.Getenv("OS_CERT"),
			ClientKeyFile:               os.Getenv("OS_KEY"),
			Cloud:                       os.Getenv("OS_CLOUD"),
			DefaultDomain:               os.Getenv("OS_DEFAULT_DOMAIN"),
			DomainID:                    os.Getenv("OS_DOMAIN_ID"),
			DomainName:                  os.Getenv("OS_DOMAIN_NAME"),
			EndpointType:                os.Getenv("OS_ENDPOINT_TYPE"),
			IdentityEndpoint:            os.Getenv("OS_AUTH_URL"),
			Password:                    os.Getenv("OS_PASSWORD"),
			ProjectDomainID:             os.Getenv("OS_PROJECT_DOMAIN_ID"),
			ProjectDomainName:           os.Getenv("OS_PROJECT_DOMAIN_NAME"),
			Region:                      os.Getenv("OS_REGION"),
			Token:                       os.Getenv("OS_TOKEN"),
			TenantID:                    tenantID,
			TenantName:                  tenantName,
			UserDomainID:                os.Getenv("OS_USER_DOMAIN_ID"),
			UserDomainName:              os.Getenv("OS_USER_DOMAIN_NAME"),
			Username:                    os.Getenv("OS_USERNAME"),
			UserID:                      os.Getenv("OS_USER_ID"),
			UseOctavia:                  true,
			ApplicationCredentialID:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
			ApplicationCredentialName:   os.Getenv("OS_APPLICATION_CREDENTIAL_NAME"),
			ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
			DelayedAuth:                 testGetenvBool("OS_DELAYED_AUTH"),
			AllowReauth:                 testGetenvBool("OS_ALLOW_REAUTH"),
			AuthOpts:                    authOpts,
			MutexKV:                     mutexkv.NewMutexKV(),
		},
	}

	if err := config.LoadAndValidate(ctx); err != nil {
		return nil, err
	}

	return &config, nil
}

func testGetenvBool(env string) bool {
	ret, _ := strconv.ParseBool(os.Getenv(env))

	return ret
}
