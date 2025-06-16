package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContainerInfraV1ClusterImport_basic(t *testing.T) {
	resourceName := "openstack_containerinfra_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Skip("Currently failing in GH-A: cant deploy cluster")
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterBasic(keypairName, clusterTemplateName, clusterName, 1),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kubeconfig"},
			},
		},
	})
}

func TestAccContainerInfraV1ClusterImport_mergeLabels(t *testing.T) {
	resourceName := "openstack_containerinfra_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Skip("Currently failing in GH-A: cant deploy cluster")
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterLabels(keypairName, clusterTemplateName, clusterName, 1, true),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kubeconfig", "merge_labels", "labels"},
			},
		},
	})
}

func TestAccContainerInfraV1ClusterImport_overrideLabels(t *testing.T) {
	resourceName := "openstack_containerinfra_cluster_v1.cluster_1"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster")
	keypairName := acctest.RandomWithPrefix("tf-acc-keypair")
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Skip("Currently failing in GH-A: cant deploy cluster")
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckContainerInfra(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterLabels(keypairName, clusterTemplateName, clusterName, 1, false),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kubeconfig", "merge_labels", "labels"},
			},
		},
	})
}
