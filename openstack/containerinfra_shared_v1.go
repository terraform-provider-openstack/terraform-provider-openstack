package openstack

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/certificates"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clusters"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clustertemplates"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/nodegroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
)

const (
	rsaPrivateKeyBlockType      = "RSA PRIVATE KEY"
	certificateRequestBlockType = "CERTIFICATE REQUEST"

	containerInfraV1ClusterUpgradeMinMicroversion = "1.8"
	containerInfraV1NodeGroupMinMicroversion      = "1.9"
	containerInfraV1ZeroNodeCountMicroversion     = "1.10"
)

func expandContainerInfraV1LabelsMap(v map[string]any) (map[string]string, error) {
	m := make(map[string]string)

	for key, val := range v {
		labelValue, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("label %s value should be string", key)
		}

		m[key] = labelValue
	}

	return m, nil
}

func expandContainerInfraV1LabelsString(v map[string]any) (string, error) {
	var formattedLabels string

	for key, val := range v {
		labelValue, ok := val.(string)
		if !ok {
			return "", fmt.Errorf("label %s value should be string", key)
		}

		formattedLabels = strings.Join([]string{
			formattedLabels,
			fmt.Sprintf("'%s':'%s'", key, labelValue),
		}, ",")
	}

	formattedLabels = strings.Trim(formattedLabels, ",")

	return fmt.Sprintf("{%s}", formattedLabels), nil
}

func containerInfraV1GetLabelsMerged(labelsAdded map[string]string, labelsSkipped map[string]string, labelsOverridden map[string]string, labels map[string]string, resourceDataLabels map[string]string) map[string]string {
	m := make(map[string]string)
	for key, val := range labelsAdded {
		m[key] = val
	}

	for key, val := range labelsSkipped {
		m[key] = val
	}

	for key := range labelsOverridden {
		// We have to get the actual value here, not the one overridden
		m[key] = labels[key]
	}
	// If defined resource's labels don't override label (are the same)
	for key, val := range resourceDataLabels {
		if _, exist := m[key]; !exist {
			m[key] = val
		}
	}

	return m
}

func containerInfraClusterTemplateV1AppendUpdateOpts(updateOpts []clustertemplates.UpdateOptsBuilder, attribute string, value any) []clustertemplates.UpdateOptsBuilder {
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

func containerInfraNodeGroupV1AppendUpdateOpts(updateOpts []nodegroups.UpdateOptsBuilder, attribute string, value int) []nodegroups.UpdateOptsBuilder {
	if value == 0 && attribute == "max_node_count" {
		updateOpts = append(updateOpts, nodegroups.UpdateOpts{
			Op:   nodegroups.RemoveOp,
			Path: strings.Join([]string{"/", attribute}, ""),
		})
	} else {
		updateOpts = append(updateOpts, nodegroups.UpdateOpts{
			Op:    nodegroups.ReplaceOp,
			Path:  strings.Join([]string{"/", attribute}, ""),
			Value: value,
		})
	}

	return updateOpts
}

// ContainerInfraClusterV1StateRefreshFunc returns a retry.StateRefreshFunc
// that is used to watch a container infra Cluster.
func containerInfraClusterV1StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, clusterID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		c, err := clusters.Get(ctx, client, clusterID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return c, "DELETE_COMPLETE", nil
			}

			return nil, "", err
		}

		errorStatuses := []string{
			"CREATE_FAILED",
			"UPDATE_FAILED",
			"DELETE_FAILED",
			"RESUME_FAILED",
			"ROLLBACK_FAILED",
		}
		for _, errorStatus := range errorStatuses {
			if c.Status == errorStatus {
				err = fmt.Errorf("openstack_containerinfra_cluster_v1 is in an error state: %s", c.StatusReason)

				return c, c.Status, err
			}
		}

		return c, c.Status, nil
	}
}

// ContainerInfraNodeGroupV1StateRefreshFunc returns a retry.StateRefreshFunc
// that is used to watch a container infra NodeGroup.
func containerInfraNodeGroupV1StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, clusterID string, nodeGroupID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		nodeGroup, err := nodegroups.Get(ctx, client, clusterID, nodeGroupID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return nodeGroup, "DELETE_COMPLETE", nil
			}

			return nil, "", err
		}

		errorStatuses := []string{
			"CREATE_FAILED",
			"UPDATE_FAILED",
			"DELETE_FAILED",
			"RESUME_FAILED",
			"ROLLBACK_FAILED",
		}
		for _, errorStatus := range errorStatuses {
			if nodeGroup.Status == errorStatus {
				err = fmt.Errorf("openstack_containerinfra_nodegroup_v1 is in an error state: %s", nodeGroup.StatusReason)

				return nodeGroup, nodeGroup.Status, err
			}
		}

		return nodeGroup, nodeGroup.Status, nil
	}
}

// containerInfraClusterV1Flavor will determine the flavor for a container infra
// cluster based on either what was set in the configuration or environment
// variable.
func containerInfraClusterV1Flavor(d *schema.ResourceData) string {
	if flavor := d.Get("flavor").(string); flavor != "" {
		return flavor
	}
	// Try the OS_MAGNUM_FLAVOR environment variable
	if v := os.Getenv("OS_MAGNUM_FLAVOR"); v != "" {
		return v
	}

	return ""
}

// containerInfraClusterV1Flavor will determine the master flavor for a
// container infra cluster based on either what was set in the configuration
// or environment variable.
func containerInfraClusterV1MasterFlavor(d *schema.ResourceData) string {
	if flavor := d.Get("master_flavor").(string); flavor != "" {
		return flavor
	}

	// Try the OS_MAGNUM_MASTER_FLAVOR environment variable
	if v := os.Getenv("OS_MAGNUM_MASTER_FLAVOR"); v != "" {
		return v
	}

	return ""
}

type kubernetesConfig struct {
	APIVersion     string                    `yaml:"apiVersion"`
	Kind           string                    `yaml:"kind"`
	Clusters       []kubernetesConfigCluster `yaml:"clusters"`
	Contexts       []kubernetesConfigContext `yaml:"contexts"`
	CurrentContext string                    `yaml:"current-context"`
	Users          []kubernetesConfigUser    `yaml:"users"`
}

type kubernetesConfigCluster struct {
	Cluster kubernetesConfigClusterData `yaml:"cluster"`
	Name    string                      `yaml:"name"`
}
type kubernetesConfigClusterData struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

type kubernetesConfigContext struct {
	Context kubernetesConfigContextData `yaml:"context"`
	Name    string                      `yaml:"name"`
}
type kubernetesConfigContextData struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type kubernetesConfigUser struct {
	Name string                   `yaml:"name"`
	User kubernetesConfigUserData `yaml:"user"`
}

type kubernetesConfigUserData struct {
	ClientKeyData         string `yaml:"client-key-data"`
	ClientCertificateData string `yaml:"client-certificate-data"`
}

func flattenContainerInfraV1Kubeconfig(ctx context.Context, d *schema.ResourceData, containerInfraClient *gophercloud.ServiceClient) (map[string]any, error) {
	clientSert, ok := d.Get("kubeconfig.client_certificate").(string)
	if ok && clientSert != "" {
		return d.Get("kubeconfig").(map[string]any), nil
	}

	certificateAuthority, err := certificates.Get(ctx, containerInfraClient, d.Id()).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error getting certificate authority: %w", err)
	}

	clientKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("Error generating client key: %w", err)
	}

	csrTemplate := x509.CertificateRequest{
		PublicKey:          clientKey.Public,
		SignatureAlgorithm: x509.SHA512WithRSA,
		Subject: pkix.Name{
			CommonName:         "admin",
			Organization:       []string{"system:masters"},
			OrganizationalUnit: []string{"terraform"},
		},
	}

	clientCsr, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, clientKey)
	if err != nil {
		return nil, fmt.Errorf("Error generating client CSR: %w", err)
	}

	pemClientKey := pem.EncodeToMemory(
		&pem.Block{
			Type:  rsaPrivateKeyBlockType,
			Bytes: x509.MarshalPKCS1PrivateKey(clientKey),
		},
	)

	pemClientCsr := pem.EncodeToMemory(
		&pem.Block{
			Type:  certificateRequestBlockType,
			Bytes: clientCsr,
		},
	)

	certificateCreateOpts := certificates.CreateOpts{
		ClusterUUID: d.Id(),
		CSR:         string(pemClientCsr),
	}

	clientCertificate, err := certificates.Create(ctx, containerInfraClient, certificateCreateOpts).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error requesting client certificate: %w", err)
	}

	name := d.Get("name").(string)
	host := d.Get("api_address").(string)

	rawKubeconfig, err := renderKubeconfig(name, host, []byte(certificateAuthority.PEM), []byte(clientCertificate.PEM), pemClientKey)
	if err != nil {
		return nil, fmt.Errorf("Error rendering kubeconfig: %w", err)
	}

	return map[string]any{
		"raw_config":             string(rawKubeconfig),
		"host":                   host,
		"cluster_ca_certificate": certificateAuthority.PEM,
		"client_certificate":     clientCertificate.PEM,
		"client_key":             string(pemClientKey),
	}, nil
}

func renderKubeconfig(name string, host string, clusterCaCertificate []byte, clientCertificate []byte, clientKey []byte) ([]byte, error) {
	userName := name + "-admin"

	config := kubernetesConfig{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: []kubernetesConfigCluster{
			{
				Name: name,
				Cluster: kubernetesConfigClusterData{
					CertificateAuthorityData: base64.StdEncoding.EncodeToString(clusterCaCertificate),
					Server:                   host,
				},
			},
		},
		Contexts: []kubernetesConfigContext{
			{
				Context: kubernetesConfigContextData{
					Cluster: name,
					User:    userName,
				},
				Name: name,
			},
		},
		CurrentContext: name,
		Users: []kubernetesConfigUser{
			{
				Name: userName,
				User: kubernetesConfigUserData{
					ClientCertificateData: base64.StdEncoding.EncodeToString(clientCertificate),
					ClientKeyData:         base64.StdEncoding.EncodeToString(clientKey),
				},
			},
		},
	}

	return yaml.Marshal(config)
}
