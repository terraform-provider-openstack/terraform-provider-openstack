package openstack

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenStackObjectStorageTempurlV1_basic(t *testing.T) {
	objectName := "object/with/slashes"
	containerName := "container"
	ttl := 60

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSwift(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackObjectstorageTempurlV1ResourceBasic(containerName, objectName, "get", ttl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectstorageTempurlV1ResourceID("openstack_objectstorage_tempurl_v1.tempurl_1"),
					testAccCheckObjectstorageTempurlV1Get("openstack_objectstorage_tempurl_v1.tempurl_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_tempurl_v1.tempurl_1", "method", "get"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_tempurl_v1.tempurl_1", "container", containerName),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_tempurl_v1.tempurl_1", "object", objectName),
				),
			},
			{
				Config: testAccOpenStackObjectstorageTempurlV1ResourceBasic(containerName, objectName, "post", ttl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectstorageTempurlV1ResourceID("openstack_objectstorage_tempurl_v1.tempurl_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_tempurl_v1.tempurl_1", "method", "post"),
				),
			},
			/* TODO(flaper87): Find a good way to test the ttl expiration
			            resource.TestStep{
							Config: testAccOpenStackObjectstorageTempurlV1ResourceBasic(containerName, objectName, "get", ),
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

func testAccCheckObjectstorageTempurlV1Get(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find temp url resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Endpoint resource ID not set")
		}

		var url string
		if url, ok = rs.Primary.Attributes["url"]; !ok {
			return fmt.Errorf("Temp URL is not set")
		}

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Failed to retrieve tempurl: %s", url)
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Failed to read tempurl body: %s", url)
		}

		if v := string(data); v != "Hello, world!" {
			return fmt.Errorf("Tempurl body doesn't match the expected data: %s", v)
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

func testAccOpenStackObjectstorageTempurlV1ResourceBasic(container, object string, method string, ttl int) string {
	return fmt.Sprintf(`
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "%s"
  metadata = {
    Temp-URL-Key = "testkey"
  }
}

resource "openstack_objectstorage_object_v1" "object_1" {
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  name           = "%s"
  content        = "Hello, world!"
}

resource "openstack_objectstorage_tempurl_v1" "tempurl_1" {
  object = "${openstack_objectstorage_object_v1.object_1.name}"
  container = "${openstack_objectstorage_container_v1.container_1.name}"
  method = "%s"
  ttl = %d
}
`, container, object, method, ttl)
}
