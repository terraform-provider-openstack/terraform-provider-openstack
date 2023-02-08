package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccContainerInfraV1ClusterTemplateDataSource_basic(t *testing.T) {
	resourceName := "data.openstack_containerinfra_clustertemplate_v1.clustertemplate_1"
	clusterTemplateName := acctest.RandomWithPrefix("tf-acc-clustertemplate")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckContainerInfra(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckContainerInfraV1ClusterTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerInfraV1ClusterTemplateBasic(clusterTemplateName),
			},
			{
				Config: testAccContainerInfraV1ClusterTemplateDataSourceBasic(
					testAccContainerInfraV1ClusterTemplateBasic(clusterTemplateName),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerInfraV1ClusterTemplateDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "region", osRegionName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterTemplateName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "apiserver_port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "coe", "kubernetes"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_distro"),
					resource.TestCheckResourceAttr(resourceName, "dns_nameserver", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "docker_storage_driver", "overlay2"),
					resource.TestCheckResourceAttr(resourceName, "docker_volume_size", "5"),
					resource.TestCheckResourceAttr(resourceName, "external_network_id", osExtGwID),
					resource.TestCheckResourceAttr(resourceName, "fixed_network", "cluster-network"),
					resource.TestCheckResourceAttr(resourceName, "fixed_subnet", "cluster-network-subnet"),
					resource.TestCheckResourceAttr(resourceName, "flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_flavor", osMagnumFlavor),
					resource.TestCheckResourceAttr(resourceName, "floating_ip_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "http_proxy", osMagnumHTTPProxy),
					resource.TestCheckResourceAttr(resourceName, "https_proxy", osMagnumHTTPSProxy),
					resource.TestCheckResourceAttr(resourceName, "image", osMagnumImage),
					resource.TestCheckResourceAttr(resourceName, "insecure_registry", "registry.example.com"),
					resource.TestCheckResourceAttr(resourceName, "keypair_id", "my-keypair"),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_tag", "1.11.1"),
					resource.TestCheckResourceAttr(resourceName, "labels.prometheus_monitoring", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.influx_grafana_dashboard_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "labels.kube_dashboard_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "master_lb_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "network_driver", "flannel"),
					resource.TestCheckResourceAttr(resourceName, "no_proxy", osMagnumNoProxy),
					resource.TestCheckResourceAttrSet(resourceName, "public"),
					resource.TestCheckResourceAttr(resourceName, "registry_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "server_type", "vm"),
					resource.TestCheckResourceAttr(resourceName, "tls_disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "volume_driver", "cinder"),
					resource.TestCheckResourceAttr(resourceName, "hidden", "false"),
				),
			},
		},
	})
}

func testAccCheckContainerInfraV1ClusterTemplateDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ct, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cluster template data source: %s", n)
		}

		if ct.Primary.ID == "" {
			return fmt.Errorf("Cluster template data source ID not set")
		}

		return nil
	}
}

func testAccContainerInfraV1ClusterTemplateDataSourceBasic(clusterTemplateResource string) string {
	return fmt.Sprintf(`
%s

data "openstack_containerinfra_clustertemplate_v1" "clustertemplate_1" {
  name = "${openstack_containerinfra_clustertemplate_v1.clustertemplate_1.name}"
}
`, clusterTemplateResource)
}
