package openstack

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/objects"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	deleteAt = "2100-12-31T14:30:59+01:00"
)

func fooMD5() string {
	return fmt.Sprintf("%x", md5.Sum([]byte("foo")))
}

func barMD5() string {
	return fmt.Sprintf("%x", md5.Sum([]byte("bar")))
}

func foobarMD5() string {
	return fmt.Sprintf("%x", md5.Sum([]byte("foobar")))
}

func manifestMD5() string {
	return fmt.Sprintf("\"%x\"", md5.Sum([]byte(fmt.Sprintf("%s%s", fooMD5(), barMD5()))))
}

func TestAccObjectStorageV1Object_basic(t *testing.T) {
	var object objects.GetHeader

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSwift(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckObjectStorageV1ObjectDestroy(t.Context(), s, "terraform/test/myfile.txt")
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1ObjectBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectStorageV1ObjectExists(t.Context(),
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
						"openstack_objectstorage_object_v1.myfile", "etag", fooMD5()),
				),
			},
			{
				Config: testAccObjectStorageV1ObjectUpdateContentType(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "application/octet-stream"),
				),
			},
			{
				Config: testAccObjectStorageV1ObjectUpdateContent,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "application/octet-stream"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", foobarMD5()),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_length", "6"),
				),
			},
			{
				Config: testAccObjectStorageV1ObjectUpdateDeleteAfter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "delete_after", "3600"),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_basic_check_destroy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckSwift(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckObjectStorageV1ObjectDestroy(t.Context(), s, "terraform/test/myfile.txt")
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1ObjectBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_container_v1.container_1", "name", "tf_test_container_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "text/plain"),
				),
			},
			{
				Config:             testAccObjectStorageV1ObjectBasic(),
				ExpectNonEmptyPlan: true,
				Check:              testAccCheckObjectStorageV1DestroyContainer(t.Context(), "tf_test_container_1", "terraform/test/myfile.txt"),
			},
			{
				Config:                    testAccObjectStorageV1ObjectBasic(),
				Destroy:                   true,
				PreventPostDestroyRefresh: true,
			},
		},
	})
}

func TestAccObjectStorageV1Object_fromSource(t *testing.T) {
	content := []byte("foo")

	tmpfile, err := os.CreateTemp(t.TempDir(), "tf_test_objectstorage_object")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		os.Remove(tmpfile.Name())
		log.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		os.Remove(tmpfile.Name())
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSwift(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckObjectStorageV1ObjectDestroy(t.Context(), s, "terraform/test/myfile")
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccObjectStorageV1ObjectFromSource, tmpfile.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_container_v1.container_1", "name", "tf_test_container_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "text/plain"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_length", strconv.Itoa(len(content))),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", fooMD5()),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_detectContentType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSwift(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckObjectStorageV1ObjectDestroy(t.Context(), s, "terraform/test/myfile.csv")
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1ObjectDetectContentType(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_container_v1.container_1", "name", "tf_test_container_1"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "content_type", "text/csv"),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", fooMD5()),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_copyFrom(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSwift(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			if err := testAccCheckObjectStorageV1ObjectDestroy(t.Context(), s, "terraform/test/myfile.txt"); err != nil {
				return err
			}

			return testAccCheckObjectStorageV1ObjectDestroy(t.Context(), s, "terraform/test/myfilecopied.txt")
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1ObjectCopyFrom,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfilesource", "etag", fooMD5()),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfilecopied", "etag", fooMD5()),
				),
			},
		},
	})
}

func TestAccObjectStorageV1Object_objectManifest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSwift(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			if err := testAccCheckObjectStorageV1ObjectDestroy(t.Context(), s, "terraform/test.csv/part001"); err != nil {
				return err
			}
			if err := testAccCheckObjectStorageV1ObjectDestroy(t.Context(), s, "terraform/test.csv/part002"); err != nil {
				return err
			}
			if err := testAccCheckObjectStorageV1ObjectDestroy(t.Context(), s, "terraform/test.csv"); err != nil {
				return err
			}

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1ObjectManifest,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile_part1", "etag", fooMD5()),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile_part2", "etag", barMD5()),
					resource.TestCheckResourceAttr(
						"openstack_objectstorage_object_v1.myfile", "etag", manifestMD5()),
				),
			},
		},
	})
}

func testAccCheckObjectStorageV1ObjectDestroy(ctx context.Context, s *terraform.State, objectname string) error {
	config := testAccProvider.Meta().(*Config)

	objectStorageClient, err := config.ObjectStorageV1Client(ctx, osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack object storage client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_objectstorage_object_v1" {
			continue
		}

		_, err := objects.Get(ctx, objectStorageClient, "tf_test_container_1", objectname, &objects.GetOpts{}).Extract()
		if err == nil {
			return errors.New("Container object still exists")
		}
	}

	return nil
}

func testAccCheckObjectStorageV1ObjectExists(ctx context.Context, n string, object *objects.GetHeader) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		objectStorageClient, err := config.ObjectStorageV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack object storage client: %w", err)
		}

		container, obj, err := parsePairedIDs(rs.Primary.ID, "openstack_objectstorage_object_v1")
		if err != nil {
			return err
		}

		found, err := objects.Get(ctx, objectStorageClient, container, obj, nil).Extract()
		if err != nil {
			return err
		}

		*object = *found

		return nil
	}
}

func testAccCheckObjectStorageV1ObjectDeleteAtMatches(expected string, object *objects.GetHeader) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
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

func testAccCheckObjectStorageV1DestroyContainer(ctx context.Context, container, object string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		objectStorageClient, err := config.ObjectStorageV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack object storage client: %w", err)
		}

		_, err = objects.Delete(ctx, objectStorageClient, container, object, nil).Extract()
		if err != nil {
			return err
		}

		_, err = containers.Delete(ctx, objectStorageClient, container).Extract()
		if err != nil {
			return err
		}

		return err
	}
}

func testAccObjectStorageV1ObjectBasic() string {
	return fmt.Sprintf(`
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
  content_type = "text/plain"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.txt"
  container_name = openstack_objectstorage_container_v1.container_1.name
  content = "foo"

  content_disposition = "foo"
  content_encoding = "utf8"
  delete_at = "%s"
}
`, deleteAt)
}

func testAccObjectStorageV1ObjectDetectContentType() string {
	return fmt.Sprintf(`
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
  content_type = "text/plain"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.csv"
  container_name = openstack_objectstorage_container_v1.container_1.name
  detect_content_type = true
  content = "foo"
  content_disposition = "foo"
  content_encoding = "utf8"
  delete_at = "%s"
}
`, deleteAt)
}

func testAccObjectStorageV1ObjectUpdateContentType() string {
	return fmt.Sprintf(`
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
  content_type = "text/plain"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.txt"
  container_name = openstack_objectstorage_container_v1.container_1.name
  content_type = "application/octet-stream"
  content = "foo"
  content_disposition = "foo"
  content_encoding = "utf8"
  delete_at = "%s"
}
`, deleteAt)
}

const testAccObjectStorageV1ObjectUpdateDeleteAfter = `
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
  content_type = "text/plain"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.txt"
  container_name = openstack_objectstorage_container_v1.container_1.name
  content_type = "application/octet-stream"
  content = "foo"
  content_encoding = "utf8"
  delete_after = "3600"
}
`

const testAccObjectStorageV1ObjectUpdateContent = `
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
  content_type = "text/plain"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.txt"
  container_name = openstack_objectstorage_container_v1.container_1.name
  content_type = "application/octet-stream"
  content = "foobar"

}
`

const testAccObjectStorageV1ObjectFromSource = `
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test/myfile.txt"
  container_name = openstack_objectstorage_container_v1.container_1.name
  detect_content_type = true
  source = "%s"
}
`

const testAccObjectStorageV1ObjectCopyFrom = `
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
}

resource "openstack_objectstorage_object_v1" "myfilesource" {
  name = "terraform/test/myfile.txt"
  container_name = openstack_objectstorage_container_v1.container_1.name
  content = "foo"
}

resource "openstack_objectstorage_object_v1" "myfilecopied" {
  name = "terraform/test/myfilecopied.txt"
  container_name = openstack_objectstorage_container_v1.container_1.name
  copy_from = "${openstack_objectstorage_container_v1.container_1.name}/${openstack_objectstorage_object_v1.myfilesource.name}"
}
`

const testAccObjectStorageV1ObjectManifest = `
resource "openstack_objectstorage_container_v1" "container_1" {
  name = "tf_test_container_1"
}

resource "openstack_objectstorage_object_v1" "myfile_part1" {
  name = "terraform/test.csv/part001"
  container_name = openstack_objectstorage_container_v1.container_1.name
  content = "foo"
}
resource "openstack_objectstorage_object_v1" "myfile_part2" {
  name = "terraform/test.csv/part002"
  container_name = openstack_objectstorage_container_v1.container_1.name
  content = "bar"
}

resource "openstack_objectstorage_object_v1" "myfile" {
  name = "terraform/test.csv"
  container_name = openstack_objectstorage_container_v1.container_1.name
  object_manifest = format("%s/terraform/test.csv/part",openstack_objectstorage_container_v1.container_1.name)

  metadata = {
    race = openstack_objectstorage_object_v1.myfile_part1.id
    condition = openstack_objectstorage_object_v1.myfile_part2.id
  }
}
`
