package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpenStackObjectStorageTempurlV1_basic(t *testing.T) {
	objectName := "object"
	containerName := "container"
	ttl := 60

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackObjectstorageTempurlV1Resource_basic(containerName, objectName, "get", ttl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectstorageTempurlV1ResourceID("openstack_objectstorage_tempurl_v1.tempurl_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_tempurl_v1.tempurl_1", "method", "get"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_tempurl_v1.tempurl_1", "container", containerName),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_tempurl_v1.tempurl_1", "object", objectName),
				),
			},
			{
				Config: testAccOpenStackObjectstorageTempurlV1Resource_basic(containerName, objectName, "post", ttl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectstorageTempurlV1ResourceID("openstack_objectstorage_tempurl_v1.tempurl_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_tempurl_v1.tempurl_1", "method", "post"),
				),
			},
			/* TODO(flaper87): Find a good way to test the ttl expiration
			            resource.TestStep{
							Config: testAccOpenStackObjectstorageTempurlV1Resource_basic(containerName, objectName, "get", ),
							Check: resource.ComposeTestCheckFunc(
								resource.TestCheckResourceAttr(
									"openstack_objectstorage_tempurl_v1.tempurl_1", "method", "get"),
								testAccCheckObjectstorageTempurlV1Expired("openstack_objectstorage_tempurl_v1.tempurl_1", 1),
							),
						},*/
		},
	})
}

func testAccCheckObjectstorageTempurlV1ResourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find temp url resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Endpoint resource ID not set")
		}

		return nil
	}
}

/*func testAccCheckObjectstorageTempurlV1Expired(n string, ttl int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(ttl))
		err := testAccCheckObjectstorageTempurlV1ResourceID(n)(s)
		if err == nil {
			return fmt.Errorf("The temp url didn't expire")
		}
		return nil
	}
}*/

func testAccOpenStackObjectstorageTempurlV1Resource_basic(container, object string, method string, ttl int) string {
	return fmt.Sprintf(`
	resource "openstack_objectstorage_tempurl_v1" "tempurl_1" {
      object = "%s"
      container = "%s"
      method = "%s"
      ttl = %d
	}
`, object, container, method, ttl)
}
