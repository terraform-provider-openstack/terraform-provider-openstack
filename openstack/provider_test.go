package openstack

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/pathorcontents"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var (
	OS_DB_ENVIRONMENT               = os.Getenv("OS_DB_ENVIRONMENT")
	OS_DB_DATASTORE_VERSION         = os.Getenv("OS_DB_DATASTORE_VERSION")
	OS_DB_DATASTORE_TYPE            = os.Getenv("OS_DB_DATASTORE_TYPE")
	OS_DEPRECATED_ENVIRONMENT       = os.Getenv("OS_DEPRECATED_ENVIRONMENT")
	OS_DNS_ENVIRONMENT              = os.Getenv("OS_DNS_ENVIRONMENT")
	OS_EXTGW_ID                     = os.Getenv("OS_EXTGW_ID")
	OS_FLAVOR_ID                    = os.Getenv("OS_FLAVOR_ID")
	OS_FLAVOR_NAME                  = os.Getenv("OS_FLAVOR_NAME")
	OS_IMAGE_ID                     = os.Getenv("OS_IMAGE_ID")
	OS_IMAGE_NAME                   = os.Getenv("OS_IMAGE_NAME")
	OS_MAGNUM_FLAVOR                = os.Getenv("OS_MAGNUM_FLAVOR")
	OS_NETWORK_ID                   = os.Getenv("OS_NETWORK_ID")
	OS_POOL_NAME                    = os.Getenv("OS_POOL_NAME")
	OS_REGION_NAME                  = os.Getenv("OS_REGION_NAME")
	OS_SWIFT_ENVIRONMENT            = os.Getenv("OS_SWIFT_ENVIRONMENT")
	OS_LB_ENVIRONMENT               = os.Getenv("OS_LB_ENVIRONMENT")
	OS_FW_ENVIRONMENT               = os.Getenv("OS_FW_ENVIRONMENT")
	OS_VPN_ENVIRONMENT              = os.Getenv("OS_VPN_ENVIRONMENT")
	OS_USE_OCTAVIA                  = os.Getenv("OS_USE_OCTAVIA")
	OS_CONTAINER_INFRA_ENVIRONMENT  = os.Getenv("OS_CONTAINER_INFRA_ENVIRONMENT")
	OS_SFS_ENVIRONMENT              = os.Getenv("OS_SFS_ENVIRONMENT")
	OS_TRANSPARENT_VLAN_ENVIRONMENT = os.Getenv("OS_TRANSPARENT_VLAN_ENVIRONMENT")
	OS_KEYMANAGER_ENVIRONMENT       = os.Getenv("OS_KEYMANAGER_ENVIRONMENT")
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"openstack": testAccProvider,
	}
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	v := os.Getenv("OS_AUTH_URL")
	if v == "" {
		t.Fatal("OS_AUTH_URL must be set for acceptance tests")
	}

	if OS_IMAGE_ID == "" && OS_IMAGE_NAME == "" {
		t.Fatal("OS_IMAGE_ID or OS_IMAGE_NAME must be set for acceptance tests")
	}

	if OS_POOL_NAME == "" {
		t.Fatal("OS_POOL_NAME must be set for acceptance tests")
	}

	if OS_FLAVOR_ID == "" && OS_FLAVOR_NAME == "" {
		t.Fatal("OS_FLAVOR_ID or OS_FLAVOR_NAME must be set for acceptance tests")
	}

	if OS_NETWORK_ID == "" {
		t.Fatal("OS_NETWORK_ID must be set for acceptance tests")
	}

	if OS_EXTGW_ID == "" {
		t.Fatal("OS_EXTGW_ID must be set for acceptance tests")
	}
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	// Do not run the test if this is a deprecated testing environment.
	if OS_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}
}

func testAccPreCheckDeprecated(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}
}

func testAccPreCheckDNS(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_DNS_ENVIRONMENT == "" {
		t.Skip("This environment does not support DNS tests")
	}
}

func testAccPreCheckSwift(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_SWIFT_ENVIRONMENT == "" {
		t.Skip("This environment does not support Swift tests")
	}
}

func testAccPreCheckDatabase(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_DB_ENVIRONMENT == "" {
		t.Skip("This environment does not support Database tests")
	}
}

func testAccPreCheckLB(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_LB_ENVIRONMENT == "" {
		t.Skip("This environment does not support LB tests")
	}
}

func testAccPreCheckFW(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_FW_ENVIRONMENT == "" {
		t.Skip("This environment does not support FW tests")
	}
}

func testAccPreCheckVPN(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_VPN_ENVIRONMENT == "" {
		t.Skip("This environment does not support VPN tests")
	}
}

func testAccPreCheckKeyManager(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_KEYMANAGER_ENVIRONMENT == "" {
		t.Skip("This environment does not support Barbican Keymanager tests")
	}
}

func testAccPreCheckContainerInfra(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_CONTAINER_INFRA_ENVIRONMENT == "" {
		t.Skip("This environment does not support Container Infra tests")
	}
}

func testAccPreCheckSFS(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	/* TODO: enable when ready in OpenLab
	if OS_SFS_ENVIRONMENT == "" {
		t.Skip("This environment does not support Shared File Systems tests")
	}
	*/
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

	if OS_TRANSPARENT_VLAN_ENVIRONMENT == "" {
		t.Skip("This environment does not support 'transparent-vlan' extension tests")
	}
}

func testAccPreCheckAdminOnly(t *testing.T) {
	v := os.Getenv("OS_USERNAME")
	if v != "admin" {
		t.Skip("Skipping test because it requires the admin user")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
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

	raw := map[string]interface{}{
		"cacert_file": caFile,
	}
	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("Unexpected err when specifying OpenStack CA by file: %s", err)
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
	raw := map[string]interface{}{
		"cacert_file": caContents,
	}
	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("Unexpected err when specifying OpenStack CA by string: %s", err)
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

	raw := map[string]interface{}{
		"cert": certFile,
		"key":  keyFile,
	}
	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("Unexpected err when specifying OpenStack Client keypair by file: %s", err)
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

	raw := map[string]interface{}{
		"cert": certContents,
		"key":  keyContents,
	}
	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("Unexpected err when specifying OpenStack Client keypair by contents: %s", err)
	}
}

func envVarContents(varName string) (string, error) {
	contents, _, err := pathorcontents.Read(os.Getenv(varName))
	if err != nil {
		return "", fmt.Errorf("Error reading %s: %s", varName, err)
	}
	return contents, nil
}

func envVarFile(varName string) (string, error) {
	contents, err := envVarContents(varName)
	if err != nil {
		return "", err
	}

	tmpFile, err := ioutil.TempFile("", varName)
	if err != nil {
		return "", fmt.Errorf("Error creating temp file: %s", err)
	}
	if _, err := tmpFile.Write([]byte(contents)); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error writing temp file: %s", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error closing temp file: %s", err)
	}
	return tmpFile.Name(), nil
}

func testAccAuthFromEnv() (*Config, error) {
	tenantID := os.Getenv("OS_TENANT_ID")
	if tenantID == "" {
		tenantID = os.Getenv("OS_PROJECT_ID")
	}

	tenantName := os.Getenv("OS_TENANT_NAME")
	if tenantName == "" {
		tenantName = os.Getenv("OS_PROJECT_NAME")
	}

	config := Config{
		CACertFile:        os.Getenv("OS_CACERT"),
		ClientCertFile:    os.Getenv("OS_CERT"),
		ClientKeyFile:     os.Getenv("OS_KEY"),
		Cloud:             os.Getenv("OS_CLOUD"),
		DefaultDomain:     os.Getenv("OS_DEFAULT_DOMAIN"),
		DomainID:          os.Getenv("OS_DOMAIN_ID"),
		DomainName:        os.Getenv("OS_DOMAIN_NAME"),
		EndpointType:      os.Getenv("OS_ENDPOINT_TYPE"),
		IdentityEndpoint:  os.Getenv("OS_AUTH_URL"),
		Password:          os.Getenv("OS_PASSWORD"),
		ProjectDomainID:   os.Getenv("OS_PROJECT_DOMAIN_ID"),
		ProjectDomainName: os.Getenv("OS_PROJECT_DOMAIN_NAME"),
		Region:            os.Getenv("OS_REGION"),
		Token:             os.Getenv("OS_TOKEN"),
		TenantID:          tenantID,
		TenantName:        tenantName,
		UserDomainID:      os.Getenv("OS_USER_DOMAIN_ID"),
		UserDomainName:    os.Getenv("OS_USER_DOMAIN_NAME"),
		Username:          os.Getenv("OS_USERNAME"),
		UserID:            os.Getenv("OS_USER_ID"),
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}
