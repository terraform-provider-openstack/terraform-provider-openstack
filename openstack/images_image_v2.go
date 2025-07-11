package openstack

import (
	"compress/bzip2"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/members"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/klauspost/compress/zstd"
	"github.com/terraform-provider-openstack/utils/v2/mutexkv"
	"github.com/ulikunitz/xz"
)

func resourceImagesImageV2MemberStatusFromString(v string) images.ImageMemberStatus {
	switch v {
	case string(images.ImageMemberStatusAccepted):
		return images.ImageMemberStatusAccepted
	case string(images.ImageMemberStatusPending):
		return images.ImageMemberStatusPending
	case string(images.ImageMemberStatusRejected):
		return images.ImageMemberStatusRejected
	case string(images.ImageMemberStatusAll):
		return images.ImageMemberStatusAll
	}

	return ""
}

func resourceImagesImageV2VisibilityFromString(v string) images.ImageVisibility {
	switch v {
	case string(images.ImageVisibilityPublic):
		return images.ImageVisibilityPublic
	case string(images.ImageVisibilityPrivate):
		return images.ImageVisibilityPrivate
	case string(images.ImageVisibilityShared):
		return images.ImageVisibilityShared
	case string(images.ImageVisibilityCommunity):
		return images.ImageVisibilityCommunity
	}

	return ""
}

func fileMD5Checksum(f *os.File) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func resourceImagesImageV2FileProps(filename string) (int64, string, error) {
	var filesize int64

	var filechecksum string

	file, err := os.Open(filename)
	if err != nil {
		return -1, "", fmt.Errorf("Error opening file for Image: %w", err)
	}
	defer file.Close()

	fstat, err := file.Stat()
	if err != nil {
		return -1, "", fmt.Errorf("Error reading image file %q: %w", file.Name(), err)
	}

	filesize = fstat.Size()

	filechecksum, err = fileMD5Checksum(file)
	if err != nil {
		return -1, "", fmt.Errorf("Error computing image file %q checksum: %w", file.Name(), err)
	}

	return filesize, filechecksum, nil
}

func resourceImagesImageV2File(ctx context.Context, client *gophercloud.ServiceClient, d *schema.ResourceData, mutexKV *mutexkv.MutexKV) (string, error) {
	if filename := d.Get("local_file_path").(string); filename != "" {
		return filename, nil
	}

	furl := d.Get("image_source_url").(string)
	if furl == "" {
		return "", errors.New("Error in config. no file specified")
	}

	dir := d.Get("image_cache_path").(string)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("unable to create dir %s: %w", dir, err)
	}

	// calculate the hashsum and create a lock to prevent simultaneous file access
	md5sum := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%t", furl, d.Get("decompress").(bool)))))
	mutexKV.Lock(md5sum)
	defer mutexKV.Unlock(md5sum)

	filename := filepath.Join(dir, md5sum+".img")
	// a cleanup func to delete a failed file
	delFile := func() {
		if err := os.Remove(filename); err != nil {
			log.Printf("[DEBUG] failed to cleanup the %q file: %s", filename, err)
		}
	}

	info, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("Error while trying to access file %q: %w", filename, err)
	}

	// check if the file size is zero
	// it could be a leftover from older provider versions
	if info != nil {
		if info.Size() != 0 {
			log.Printf("[DEBUG] File exists %s", filename)

			return filename, nil
		}
		// delete the zero size file
		delFile()
	}

	log.Printf("[DEBUG] File doens't exists %s. will download from %s", filename, furl)

	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("Error creating file %q: %w", filename, err)
	}

	defer file.Close()

	httpClient := &client.HTTPClient

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, furl, nil)
	if err != nil {
		delFile()

		return "", fmt.Errorf("Error creating a new request: %w", err)
	}

	username := d.Get("image_source_username").(string)
	password := d.Get("image_source_password").(string)

	if username != "" && password != "" {
		request.SetBasicAuth(username, password)
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		delFile()

		return "", fmt.Errorf("Error downloading image from %q: %w", furl, err)
	}

	// check for credential error among other errors
	if resp.StatusCode != http.StatusOK {
		delFile()

		return "", fmt.Errorf("Error downloading image from %q, statusCode is %d", furl, resp.StatusCode)
	}

	defer resp.Body.Close()
	reader := resp.Body

	decompress := d.Get("decompress").(bool)
	if decompress {
		reader, err = resourceImagesImageV2DetectCompression(resp)
		if err != nil {
			delFile()

			return "", err
		}
		defer reader.Close()
	}

	if _, err = io.Copy(file, reader); err != nil {
		delFile()

		return "", fmt.Errorf("Error downloading image %q to file %q: %w", furl, filename, err)
	}

	return filename, nil
}

func resourceImagesImageV2DetectCompression(resp *http.Response) (io.ReadCloser, error) {
	ct := resp.Header.Get("Content-Type")

	// extract filename extension from "Content-Disposition" header
	cd := resp.Header.Get("Content-Disposition")
	_, p, _ := mime.ParseMediaType(cd)
	ext := strings.Trim(filepath.Ext(p["filename"]), ".")

	formats := []string{ct, ext}
	for _, format := range formats {
		switch format {
		case "gz", "gzip", "application/gzip", "application/x-gzip":
			reader, err := gzip.NewReader(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Error decompressing gzip image: %w", err)
			}

			return reader, nil
		case "bzip2", "application/bzip2", "application/x-bzip2":
			bz2Reader := bzip2.NewReader(resp.Body)

			return io.NopCloser(bz2Reader), nil
		case "xz", "application/xz", "application/x-xz":
			xzReader, err := xz.NewReader(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Error decompressing xz image: %w", err)
			}

			return io.NopCloser(xzReader), nil
		case "zst", "zstd", "application/zstd", "application/x-zstd":
			zstdReader, err := zstd.NewReader(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Error decompressing zstd image: %w", err)
			}

			return zstdReader.IOReadCloser(), nil
		case "application/octet-stream", "binary/octet-stream":
			// This is a fallback for cases where the server does not provide
			// a Content-Type header. In this case, we'll try to detect the
			// compression based on the Content-Disposition header.
		case "":
			// This triggers the error at the end of the function.
		default:
			return nil, fmt.Errorf("Error decompressing image, format %s is not supported", format)
		}
	}

	return nil, fmt.Errorf("Error decompressing image, detected %q formats are not supported", formats)
}

func resourceImagesImageV2RefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		img, err := images.Get(ctx, client, id).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenStack image status is: %s", img.Status)

		return img, string(img.Status), nil
	}
}

func resourceImagesImageV2BuildTags(v []any) []string {
	tags := make([]string, len(v))
	for i, tag := range v {
		tags[i] = tag.(string)
	}

	return tags
}

func resourceImagesImageV2ExpandProperties(v map[string]any) map[string]string {
	properties := map[string]string{}

	for key, value := range v {
		if v, ok := value.(string); ok {
			properties[key] = v
		}
	}

	return properties
}

func resourceImagesImageV2UpdateComputedAttributes(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	if diff.HasChange("properties") {
		// Only check if the image has been created.
		if diff.Id() != "" {
			// Try to reconcile the properties set by the server
			// with the properties set by the user.
			//
			// old = user properties + server properties
			// new = user properties only
			o, n := diff.GetChange("properties")

			newProperties := resourceImagesImageV2ExpandProperties(n.(map[string]any))

			for oldKey, oldValue := range o.(map[string]any) {
				// os_ keys are provided by the OpenStack Image service.
				if strings.HasPrefix(oldKey, "os_") {
					if v, ok := oldValue.(string); ok {
						newProperties[oldKey] = v
					}
				}

				// stores is provided by the OpenStack Image service.
				if oldKey == "stores" {
					if v, ok := oldValue.(string); ok {
						newProperties[oldKey] = v
					}
				}

				// direct_url is provided by some storage drivers.
				if oldKey == "direct_url" {
					if v, ok := oldValue.(string); ok {
						newProperties[oldKey] = v
					}
				}
			}

			// Set the diff to the newProperties, which includes the server-side
			// os_ properties.
			//
			// If the user has changed properties, they will be caught at this
			// point, too.
			if err := diff.SetNew("properties", newProperties); err != nil {
				log.Printf("[DEBUG] unable set diff for properties key: %s", err)
			}
		}
	}

	return nil
}

func resourceImagesImageAccessV2DetectMemberID(ctx context.Context, client *gophercloud.ServiceClient, imageID string) (string, error) {
	allPages, err := members.List(client, imageID).AllPages(ctx)
	if err != nil {
		return "", fmt.Errorf("Unable to list image members: %w", err)
	}

	allMembers, err := members.ExtractMembers(allPages)
	if err != nil {
		return "", fmt.Errorf("Unable to extract image members: %w", err)
	}

	if len(allMembers) == 0 {
		return "", fmt.Errorf("No members found for the %q image", imageID)
	}

	if len(allMembers) > 1 {
		return "", fmt.Errorf("Too many members found for the %q image, please specify the member_id explicitly", imageID)
	}

	return allMembers[0].MemberID, nil
}

func imagesFilterByRegex(imageArr []images.Image, nameRegex string) []images.Image {
	var result []images.Image

	r := regexp.MustCompile(nameRegex)

	for _, image := range imageArr {
		// Check for a very rare case where the response would include no
		// image name. No name means nothing to attempt a match against,
		// therefore we are skipping such image.
		if image.Name == "" {
			log.Printf("[WARN] Unable to find image name to match against "+
				"for image ID %q owned by %q, nothing to do.",
				image.ID, image.Owner)

			continue
		}

		if r.MatchString(image.Name) {
			result = append(result, image)
		}
	}

	return result
}

// v - slice of images to filter
// p - field "properties" of schema.Resource from dataSourceImagesImageIDsV2
//
//	or dataSourceImagesImageV2. If p is empty no filtering applies and the
//	function returns the v.
func imagesFilterByProperties(v []images.Image, p map[string]string) []images.Image {
	var result []images.Image

	if len(p) > 0 {
		for _, image := range v {
			if len(image.Properties) > 0 {
				match := true

				for searchKey, searchValue := range p {
					imageValue, ok := image.Properties[searchKey]
					if !ok {
						match = false

						break
					}

					if searchValue != imageValue {
						match = false

						break
					}
				}

				if match {
					result = append(result, image)
				}
			}
		}
	} else {
		result = v
	}

	return result
}

// dataSourceValidateImageSortFilter checks that the sorting filter passed
// by a user follows the glance API definition of "sort_key:sort_dir,sort_key:sort_dir:.."
// where sort_dir is optional. For more details, check the glance api docs
// https://docs.openstack.org/api-ref/image/v2/index.html#list-images
// at the `sorting` section.
func dataSourceValidateImageSortFilter(v any, _ string) (ws []string, errors []error) {
	validSortKeys := []string{
		"name",
		"owner",
		"protected",
		"status",
		"tag",
		"visibility",
		"container_format",
		"disk_format",
		"os_hidden",
		"member_status",
		"created_at",
		"updated_at",
	}
	validSortDirections := []string{
		"asc",
		"desc",
	}

	sortPairs := strings.Split(v.(string), ",")
	for _, sortPair := range sortPairs {
		parts := strings.Split(sortPair, ":")
		if !strSliceContains(validSortKeys, parts[0]) {
			errors = append(errors, fmt.Errorf("Invalid sorting key: %s. Has to be one of: %+q", parts[0], validSortKeys))
		}

		if len(parts) > 1 && !strSliceContains(validSortDirections, parts[1]) {
			errors = append(errors, fmt.Errorf("Invalid sorting direction: %s. Has to be one of: %+q", parts[0], validSortDirections))
		}
	}

	return ws, errors
}
