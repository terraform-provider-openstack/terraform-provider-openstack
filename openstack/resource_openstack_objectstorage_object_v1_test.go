package openstack

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	deleteAt = "2100-12-31T14:30:59+01:00"
)

var (
	fooMD5      = fmt.Sprintf("%x", md5.Sum([]byte("foo")))
	barMD5      = fmt.Sprintf("%x", md5.Sum([]byte("bar")))
	foobarMD5   = fmt.Sprintf("%x", md5.Sum([]byte("foobar")))
	manifestMD5 = fmt.Sprintf("\"%x\"", md5.Sum([]byte(fmt.Sprintf("%s%s", fooMD5, barMD5))))
)

func TestAccObjectStorageV1Object_basic(t *testing.T) {
	var object objects.GetHeader

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckSwift(t) },
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test/myfile.txt")
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1Object_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectStorageV1ObjectExists(
						"openstack_objectstorage_object_v1.myfile", &object),
					testAccCheckObjectStorageV1ObjectDeleteAtMatches(deleteAt, &object),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_container_v1.container_1", "name", "tf_test_container_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "text/plain"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_length", "3"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_disposition", "foo"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_encoding", "utf8"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", fooMD5),
				),
			},
			{
				Config: testAccObjectStorageV1Object_updateContentType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "application/octet-stream"),
				),
			},
			{
				Config: testAccObjectStorageV1Object_updateContent,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "application/octet-stream"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", foobarMD5),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_length", "6"),
				),
			},
			{
				Config: testAccObjectStorageV1Object_updateDeleteAfter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "delete_after", "3600"),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_fromSource(t *testing.T) {
	content := []byte("foo")
	tmpfile, err := ioutil.TempFile("", "tf_test_objectstorage_object")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckSwift(t) },
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test/myfile")
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccObjectStorageV1Object_fromSource, tmpfile.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_container_v1.container_1", "name", "tf_test_container_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "text/plain"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_length", fmt.Sprintf("%v", len(content))),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", fooMD5),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_detectContentType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckSwift(t) },
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test/myfile.csv")
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1Object_detectContentType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_container_v1.container_1", "name", "tf_test_container_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "text/csv"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", fooMD5),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_copyFrom(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckSwift(t) },
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			if err := testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test/myfile.txt"); err != nil {
				return err
			}
			return testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test/myfilecopied.txt")
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1Object_copyFrom,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfilesource", "etag", fooMD5),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfilecopied", "etag", fooMD5),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_objectManifest(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckSwift(t) },
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			if err := testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test.csv/part001"); err != nil {
				return err
			}
			if err := testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test.csv/part002"); err != nil {
				return err
			}
			if err := testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test.csv"); err != nil {
				return err
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1Object_objectManifest,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile_part1", "etag", fooMD5),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile_part2", "etag", barMD5),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", manifestMD5),
				),
			},
		},
	})
}

func testAccCheckObjectStorageV1ObjectDestroy(s *terraform.State, objectname string) error {
	config := testAccProvider.Meta().(*Config)
	objectStorageClient, err := config.objectStorageV1Client(OS_REGION_NAME)
	if err != nil {

		return fmt.Errorf("Error creating OpenStack object storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_objectstorage_object_v1" {
			continue
		}

		_, err := objects.Get(objectStorageClient, "tf_test_container_1", objectname, &objects.GetOpts{}).Extract()
		if err == nil {
			return fmt.Errorf("Container object still exists")
		}
	}

	return nil
}

func testAccCheckObjectStorageV1ObjectExists(n string, object *objects.GetHeader) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		objectStorageClient, err := config.objectStorageV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack object storage client: %s", err)
		}

		parts := strings.SplitN(rs.Primary.ID, "/", 2)
		if len(parts) < 2 {
			return fmt.Errorf("Malformed object name: %s", rs.Primary.ID)
		}

		found, err := objects.Get(objectStorageClient, parts[0], parts[1], nil).Extract()
		if err != nil {
			return err
		}

		*object = *found

		return nil
	}
}

func testAccCheckObjectStorageV1ObjectDeleteAtMatches(expected string, object *objects.GetHeader) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		expectedTime, err := time.Parse(time.RFC3339, expected)
		if err != nil {
			return err
		}

		if !expectedTime.Equal(object.DeleteAt) {
			return fmt.Errorf("%s and %s do not match", expected, object.DeleteAt)
		}

		return nil
	}
}

var testAccObjectStorageV1Object_basic = fmt.Sprintf(`
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
  content_type = "text/plain"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.txt"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  content = "foo"

  content_disposition = "foo"
  content_encoding = "utf8"
  delete_at = "%s"
}
`, deleteAt)

var testAccObjectStorageV1Object_detectContentType = fmt.Sprintf(`
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
  content_type = "text/plain"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.csv"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  detect_content_type = true
  content = "foo"
  content_disposition = "foo"
  content_encoding = "utf8"
  delete_at = "%s"
}
`, deleteAt)

var testAccObjectStorageV1Object_updateContentType = fmt.Sprintf(`
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
  content_type = "text/plain"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.txt"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  content_type = "application/octet-stream"
  content = "foo"
  content_disposition = "foo"
  content_encoding = "utf8"
  delete_at = "%s"
}
`, deleteAt)

const testAccObjectStorageV1Object_updateDeleteAfter = `
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
  content_type = "text/plain"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.txt"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  content_type = "application/octet-stream"
  content = "foo"
  content_encoding = "utf8"
  delete_after = "3600"
}
`

const testAccObjectStorageV1Object_updateContent = `
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
  content_type = "text/plain"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.txt"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  content_type = "application/octet-stream"
  content = "foobar"

}
`

const testAccObjectStorageV1Object_fromSource = `
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.txt"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  detect_content_type = true
  source = "%s"
}
`

const testAccObjectStorageV1Object_copyFrom = `
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
}

resource "openstack_objectstorage_object_v1" "myfilesource" {
  name = "terraform/test/myfile.txt"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  content = "foo"
}

resource "openstack_objectstorage_object_v1" "myfilecopied" {
  name = "terraform/test/myfilecopied.txt"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  copy_from = "${openstack_objectstorage_container_v1.container_1.name}/${openstack_objectstorage_object_v1.myfilesource.name}"
}
`

const testAccObjectStorageV1Object_objectManifest = `
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
}

resource "openstack_objectstorage_object_v1" "myfile_part1" {
  name = "terraform/test.csv/part001"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  content = "foo"
}
resource "openstack_objectstorage_object_v1" "myfile_part2" {
  name = "terraform/test.csv/part002"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  content = "bar"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test.csv"
  container_name = "${openstack_objectstorage_container_v1.container_1.name}"
  object_manifest = "${format("%s/terraform/test.csv/part",openstack_objectstorage_container_v1.container_1.name)}"

  metadata = {
    race = "${openstack_objectstorage_object_v1.myfile_part1.id}"
    condition = "${openstack_objectstorage_object_v1.myfile_part2.id}"
  }
}
`
