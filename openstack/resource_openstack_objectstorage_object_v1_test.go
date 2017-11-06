package openstack

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccObjectStorageV1Object_basic(t *testing.T) {
	foomd5 := fmt.Sprintf("%x", md5.Sum([]byte("foo")))
	foobarmd5 := fmt.Sprintf("%x", md5.Sum([]byte("foobar")))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSwift(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test/myfile.txt")
		},
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccObjectStorageV1Object_basic,
				Check: resource.ComposeTestCheckFunc(
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
						"openstack_objectstorage_object_v1.myfile", "delete_at", "2100-12-31T14:30:59+01:00"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", foomd5),
				),
			},
			resource.TestStep{
				Config: testAccObjectStorageV1Object_update_content_type,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "application/octet-stream"),
				),
			},
			resource.TestStep{
				Config: testAccObjectStorageV1Object_update_content,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "application/octet-stream"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", foobarmd5),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_length", "6"),
				),
			},
			resource.TestStep{
				Config: testAccObjectStorageV1Object_update_delete,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "delete_after", "3600"),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_fromsource(t *testing.T) {
	content := []byte("foo")
	foomd5 := fmt.Sprintf("%x", md5.Sum(content))
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
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSwift(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test/myfile")
		},
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccObjectStorageV1Object_fromsource, tmpfile.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_container_v1.container_1", "name", "tf_test_container_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "text/plain"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_length", fmt.Sprintf("%v", len(content))),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", foomd5),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_detectcontenttype(t *testing.T) {
	foomd5 := fmt.Sprintf("%x", md5.Sum([]byte("foo")))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSwift(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test/myfile.csv")
		},
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccObjectStorageV1Object_detectcontenttype,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_container_v1.container_1", "name", "tf_test_container_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "text/csv"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", foomd5),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_copyfrom(t *testing.T) {
	foomd5 := fmt.Sprintf("%x", md5.Sum([]byte("foo")))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSwift(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			if err := testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test/myfile.txt"); err != nil {
				return err
			}
			return testAccCheckObjectStorageV1ObjectDestroy(s, "terraform/test/myfilecopied.txt")
		},
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccObjectStorageV1Object_copyfrom,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfilesource", "etag", foomd5),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfilecopied", "etag", foomd5),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_objectmanifest(t *testing.T) {
	foomd5 := fmt.Sprintf("%x", md5.Sum([]byte("foo")))
	barmd5 := fmt.Sprintf("%x", md5.Sum([]byte("bar")))
	manifestmd5 := fmt.Sprintf("\"%x\"", md5.Sum([]byte(fmt.Sprintf("%s%s", foomd5, barmd5))))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSwift(t)
		},
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
			resource.TestStep{
				Config: testAccObjectStorageV1Object_objectmanifest,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile_part1", "etag", foomd5),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile_part2", "etag", barmd5),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", manifestmd5),
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

const testAccObjectStorageV1Object_basic = `
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
  delete_at = "2100-12-31T14:30:59+01:00"
}
`
const testAccObjectStorageV1Object_detectcontenttype = `
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
  delete_at = "2100-12-31T14:30:59+01:00"
}
`

const testAccObjectStorageV1Object_update_content_type = `
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
  delete_at = "2100-12-31T14:30:59+01:00"

}
`

const testAccObjectStorageV1Object_update_delete = `
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

const testAccObjectStorageV1Object_update_content = `
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

const testAccObjectStorageV1Object_fromsource = `
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

const testAccObjectStorageV1Object_copyfrom = `
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

const testAccObjectStorageV1Object_objectmanifest = `
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
}
`
