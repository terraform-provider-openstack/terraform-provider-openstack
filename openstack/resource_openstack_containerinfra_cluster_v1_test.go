package openstack

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
)

func TestAccContainerInfraV1Cluster_basic(t *testing.T) {
	var cluster clusters.Cluster

	resourceName := "openstack_containerinfra_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterBasic(keypairName, clusterTemplateName, clusterName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "api_address"),
					resource.TestCheckResourceAttrSet(resourceName, "coe_version"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_template_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_version"),
					resource.TestCheckResourceAttr(resourceName, "create_timeout", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "discovery_url"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "keypair", keypairName),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "master_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "master_addresses.#", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_addresses.#", strconv.Itoa(1)),
					resource.TestCheckResourceAttrSet(resourceName, "stack_id"),
					resource.TestCheckResourceAttr(resourceName, "floating_ip_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.raw_config"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_key"),
				),
			},
			{
				Config: testAccContainerInfraV1ClusterBasic(keypairName, clusterTemplateName, clusterName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "api_address"),
					resource.TestCheckResourceAttrSet(resourceName, "coe_version"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_template_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_version"),
					resource.TestCheckResourceAttr(resourceName, "create_timeout", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "discovery_url"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "keypair", keypairName),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "master_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(resourceName, "master_addresses.#", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_addresses.#", strconv.Itoa(2)),
					resource.TestCheckResourceAttrSet(resourceName, "stack_id"),
					resource.TestCheckResourceAttr(resourceName, "floating_ip_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.raw_config"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_key"),
				),
			},
		},
	})
}

func TestAccContainerInfraV1Cluster_mergeLabels(t *testing.T) {
	var cluster clusters.Cluster

	resourceName := "openstack_containerinfra_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterLabels(keypairName, clusterTemplateName, clusterName, 1, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "api_address"),
					resource.TestCheckResourceAttrSet(resourceName, "coe_version"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_template_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_version"),
					resource.TestCheckResourceAttr(resourceName, "create_timeout", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "discovery_url"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "keypair", keypairName),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "merge_labels", "true"),
					resource.TestCheckResourceAttr(resourceName, "master_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "master_addresses.#", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_addresses.#", strconv.Itoa(1)),
					resource.TestCheckResourceAttrSet(resourceName, "stack_id"),
					resource.TestCheckResourceAttr(resourceName, "floating_ip_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.raw_config"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_key"),
				),
			},
			{
				Config: testAccContainerInfraV1ClusterLabels(keypairName, clusterTemplateName, clusterName, 2, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "api_address"),
					resource.TestCheckResourceAttrSet(resourceName, "coe_version"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_template_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_version"),
					resource.TestCheckResourceAttr(resourceName, "create_timeout", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "discovery_url"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "keypair", keypairName),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "merge_labels", "true"),
					resource.TestCheckResourceAttr(resourceName, "master_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(resourceName, "master_addresses.#", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_addresses.#", strconv.Itoa(2)),
					resource.TestCheckResourceAttrSet(resourceName, "stack_id"),
					resource.TestCheckResourceAttr(resourceName, "floating_ip_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.raw_config"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_key"),
				),
			},
		},
	})
}

func TestAccContainerInfraV1Cluster_overrideLabels(t *testing.T) {
	var cluster clusters.Cluster

	resourceName := "openstack_containerinfra_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterLabels(keypairName, clusterTemplateName, clusterName, 1, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "api_address"),
					resource.TestCheckResourceAttrSet(resourceName, "coe_version"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_template_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_version"),
					resource.TestCheckResourceAttr(resourceName, "create_timeout", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "discovery_url"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "keypair", keypairName),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "merge_labels", "false"),
					resource.TestCheckResourceAttr(resourceName, "master_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "master_addresses.#", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_addresses.#", strconv.Itoa(1)),
					resource.TestCheckResourceAttrSet(resourceName, "stack_id"),
					resource.TestCheckResourceAttr(resourceName, "floating_ip_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.raw_config"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_key"),
				),
			},
			{
				Config: testAccContainerInfraV1ClusterLabels(keypairName, clusterTemplateName, clusterName, 2, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "api_address"),
					resource.TestCheckResourceAttrSet(resourceName, "coe_version"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_template_id"),
					resource.TestCheckResourceAttrSet(resourceName, "container_version"),
					resource.TestCheckResourceAttr(resourceName, "create_timeout", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "discovery_url"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", strconv.Itoa(10)),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "keypair", keypairName),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttr(resourceName, "merge_labels", "false"),
					resource.TestCheckResourceAttr(resourceName, "master_count", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_count", strconv.Itoa(2)),
					resource.TestCheckResourceAttr(resourceName, "master_addresses.#", strconv.Itoa(1)),
					resource.TestCheckResourceAttr(resourceName, "node_addresses.#", strconv.Itoa(2)),
					resource.TestCheckResourceAttrSet(resourceName, "stack_id"),
					resource.TestCheckResourceAttr(resourceName, "floating_ip_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.raw_config"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig.client_key"),
				),
			},
		},
	})
}

func testAccCheckContainerInfraV1ClusterExists(n string, cluster *clusters.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		containerInfraClient, err := config.ContainerInfraV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
		}

		found, err := clusters.Get(containerInfraClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.UUID != rs.Primary.ID {
			return fmt.Errorf("Cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func testAccCheckContainerInfraV1ClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_containerinfra_cluster_v1" {
			continue
		}

		_, err := clusters.Get(containerInfraClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccContainerInfraV1ClusterBasic(keypairName, clusterTemplateName, clusterName string, nodeCount int) string {
	return fmt.Sprintf(`
resource "openstack_compute_keypair_v2" "keypair_1" {
  name = "%s"
}

resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name                  = "%s"
  image                 = "%s"
  coe                   = "kubernetes"
  floating_ip_enabled   = false
  volume_driver         = "cinder"
  docker_storage_driver = "overlay2"
  docker_volume_size    = 5
  external_network_id   = "%s"
  network_driver        = "flannel"
  http_proxy            = "%s"
  https_proxy           = "%s"
  no_proxy              = "%s"
  labels = {
    kubescheduler_options = "log-flush-frequency=1m",
	kube_dashboard_enabled = "true",
	%s
  }
}

resource "openstack_containerinfra_cluster_v1" "cluster_1" {
  region               = "%s"
  name                 = "%s"
  cluster_template_id  = "${openstack_containerinfra_clustertemplate_v1.clustertemplate_1.id}"
  create_timeout       = "40"
  docker_volume_size   = "10"
  flavor                = "%s"
  master_flavor         = "%s"
  keypair              = "${openstack_compute_keypair_v2.keypair_1.name}"
  master_count         = 1
  node_count           = %d
  floating_ip_enabled  = true
}
`, keypairName, clusterTemplateName, osMagnumImage, osExtGwID, osMagnumHTTPProxy, osMagnumHTTPSProxy, osMagnumNoProxy, osMagnumLabels, osRegionName, clusterName, osMagnumFlavor, osMagnumFlavor, nodeCount)
}

func testAccContainerInfraV1ClusterLabels(keypairName, clusterTemplateName, clusterName string, nodeCount int, mergeLabels bool) string {
	return fmt.Sprintf(`
resource "openstack_compute_keypair_v2" "keypair_1" {
  name = "%s"
}

resource "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name                  = "%s"
  image                 = "%s"
  coe                   = "kubernetes"
  floating_ip_enabled   = false
  volume_driver         = "cinder"
  docker_storage_driver = "overlay2"
  docker_volume_size    = 5
  external_network_id   = "%s"
  network_driver        = "flannel"
  http_proxy            = "%s"
  https_proxy           = "%s"
  no_proxy              = "%s"
  labels = {
    kubescheduler_options = "log-flush-frequency=1m",
	kube_dashboard_enabled = "true",
  }
}

resource "openstack_containerinfra_cluster_v1" "cluster_1" {
  region               = "%s"
  name                 = "%s"
  cluster_template_id  = "${openstack_containerinfra_clustertemplate_v1.clustertemplate_1.id}"
  create_timeout       = "40"
  docker_volume_size   = "10"
  flavor                = "%s"
  master_flavor         = "%s"
  keypair              = "${openstack_compute_keypair_v2.keypair_1.name}"
  master_count         = 1
  node_count           = %d
  floating_ip_enabled  = true
  merge_labels         = %t
  labels = {
	kube_dashboard_enabled = "false",
	%s
  }
}
`, keypairName, clusterTemplateName, osMagnumImage, osExtGwID, osMagnumHTTPProxy, osMagnumHTTPSProxy, osMagnumNoProxy, osRegionName, clusterName, osMagnumFlavor, osMagnumFlavor, nodeCount, mergeLabels, osMagnumLabels)
}
