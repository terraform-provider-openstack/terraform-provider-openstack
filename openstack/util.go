package openstack

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// BuildRequest takes an opts struct and builds a request body for
// Gophercloud to execute.
func BuildRequest(opts any, parent string) (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	b = AddValueSpecs(b)

	return map[string]any{parent: b}, nil
}

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, msg string) error {
	if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
		d.SetId("")

		return nil
	}

	return fmt.Errorf("%s %s: %w", msg, d.Id(), err)
}

// GetRegion returns the region that was specified in the resource. If a
// region was not set, the provider-level region is checked. The provider-level
// region can either be set by the region argument or by OS_REGION_NAME.
func GetRegion(d *schema.ResourceData, config *Config) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return config.Region
}

// AddValueSpecs expands the 'value_specs' object and removes 'value_specs'
// from the reqeust body.
func AddValueSpecs(body map[string]any) map[string]any {
	if body["value_specs"] != nil {
		for k, v := range body["value_specs"].(map[string]any) {
			// this hack allows to pass boolean values as strings
			if v == "true" || v == "false" {
				body[k] = v == "true"

				continue
			}

			body[k] = v
		}

		delete(body, "value_specs")
	}

	return body
}

// MapValueSpecs converts ResourceData into a map.
func MapValueSpecs(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("value_specs").(map[string]any) {
		m[key] = val.(string)
	}

	return m
}

func checkForRetryableError(err error) *retry.RetryError {
	var e gophercloud.ErrUnexpectedResponseCode

	ok := errors.As(err, &e)
	if !ok {
		return retry.NonRetryableError(err)
	}

	switch e.Actual {
	case http.StatusConflict, // 409
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return retry.RetryableError(err)
	}

	return retry.NonRetryableError(err)
}

func suppressEquivalentTimeDiffs(_, o, n string, _ *schema.ResourceData) bool {
	oldTime, err := time.Parse(time.RFC3339, o)
	if err != nil {
		return false
	}

	newTime, err := time.Parse(time.RFC3339, n)
	if err != nil {
		return false
	}

	return oldTime.Equal(newTime)
}

func resourceNetworkingAvailabilityZoneHintsV2(d *schema.ResourceData) []string {
	rawAZH := d.Get("availability_zone_hints").([]any)
	azh := make([]string, len(rawAZH))

	for i, raw := range rawAZH {
		azh[i] = raw.(string)
	}

	return azh
}

func expandVendorOptions(vendOptsRaw []any) map[string]any {
	vendorOptions := make(map[string]any)

	for _, option := range vendOptsRaw {
		for optKey, optValue := range option.(map[string]any) {
			vendorOptions[optKey] = optValue
		}
	}

	return vendorOptions
}

func expandObjectReadTags(d *schema.ResourceData, tags []string) {
	d.Set("all_tags", tags)

	allTags := d.Get("all_tags").(*schema.Set)
	desiredTags := d.Get("tags").(*schema.Set)
	actualTags := allTags.Intersection(desiredTags)

	if !actualTags.Equal(desiredTags) {
		d.Set("tags", expandToStringSlice(actualTags.List()))
	}
}

func expandObjectUpdateTags(d *schema.ResourceData) []string {
	allTags := d.Get("all_tags").(*schema.Set)
	oldTagsRaw, newTagsRaw := d.GetChange("tags")
	oldTags, newTags := oldTagsRaw.(*schema.Set), newTagsRaw.(*schema.Set)

	allTagsWithoutOld := allTags.Difference(oldTags)

	return expandToStringSlice(allTagsWithoutOld.Union(newTags).List())
}

func expandObjectTags(d *schema.ResourceData) []string {
	rawTags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, len(rawTags))

	for i, raw := range rawTags {
		tags[i] = raw.(string)
	}

	return tags
}

func expandToMapStringString(v map[string]any) map[string]string {
	m := make(map[string]string, len(v))

	for key, val := range v {
		if strVal, ok := val.(string); ok {
			m[key] = strVal
		}
	}

	return m
}

func expandToStringSlice(v []any) []string {
	s := make([]string, len(v))

	for i, val := range v {
		if strVal, ok := val.(string); ok {
			s[i] = strVal
		}
	}

	return s
}

// strSliceContains checks if a given string is contained in a slice
// When anybody asks why Go needs generics, here you go.
func strSliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}

	return false
}

func sliceUnion(a, b []string) []string {
	var res []string
	for _, i := range a {
		if !strSliceContains(res, i) {
			res = append(res, i)
		}
	}

	for _, k := range b {
		if !strSliceContains(res, k) {
			res = append(res, k)
		}
	}

	return res
}

// compatibleMicroversion will determine if an obtained microversion is
// compatible with a given microversion.
func compatibleMicroversion(direction, required, given string) (bool, error) {
	if direction != "min" && direction != "max" {
		return false, fmt.Errorf("Invalid microversion direction %s. Must be min or max", direction)
	}

	if required == "" || given == "" {
		return false, nil
	}

	requiredParts := strings.Split(required, ".")
	if len(requiredParts) != 2 {
		return false, fmt.Errorf("Not a valid microversion: %s", required)
	}

	givenParts := strings.Split(given, ".")
	if len(givenParts) != 2 {
		return false, fmt.Errorf("Not a valid microversion: %s", given)
	}

	requiredMajor, requiredMinor := requiredParts[0], requiredParts[1]
	givenMajor, givenMinor := givenParts[0], givenParts[1]

	requiredMajorInt, err := strconv.Atoi(requiredMajor)
	if err != nil {
		return false, fmt.Errorf("Unable to parse microversion: %s", required)
	}

	requiredMinorInt, err := strconv.Atoi(requiredMinor)
	if err != nil {
		return false, fmt.Errorf("Unable to parse microversion: %s", required)
	}

	givenMajorInt, err := strconv.Atoi(givenMajor)
	if err != nil {
		return false, fmt.Errorf("Unable to parse microversion: %s", given)
	}

	givenMinorInt, err := strconv.Atoi(givenMinor)
	if err != nil {
		return false, fmt.Errorf("Unable to parse microversion: %s", given)
	}

	switch direction {
	case "min":
		if requiredMajorInt == givenMajorInt {
			if requiredMinorInt <= givenMinorInt {
				return true, nil
			}
		}
	case "max":
		if requiredMajorInt == givenMajorInt {
			if requiredMinorInt >= givenMinorInt {
				return true, nil
			}
		}
	}

	return false, nil
}

func validateJSONObject(v any, k string) ([]string, []error) {
	if v == nil || v.(string) == "" {
		return nil, []error{fmt.Errorf("%q value must not be empty", k)}
	}

	var j map[string]any

	s := v.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return nil, []error{fmt.Errorf("%q must be a JSON object: %w", k, err)}
	}

	return nil, nil
}

func diffSuppressJSONObject(_, o, n string, _ *schema.ResourceData) bool {
	if strSliceContains([]string{"{}", ""}, o) &&
		strSliceContains([]string{"{}", ""}, n) {
		return true
	}

	return false
}

// Metadata in openstack are not fully replaced with a "set"
// operation, instead, it's only additive, and the existing
// metadata are only removed when set to `null` value in json.
func mapDiffWithNilValues(oldMap, newMap map[string]any) (output map[string]any) {
	output = make(map[string]any)

	for k, v := range newMap {
		output[k] = v
	}

	for key := range oldMap {
		_, ok := newMap[key]
		if !ok {
			output[key] = nil
		}
	}

	return
}

// parsePairedIDs is a helper function that parses a raw ID into two
// separate IDs. This is useful for resources that have a parent/child
// relationship.
func parsePairedIDs(id string, res string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Unable to determine %s ID from raw ID: %s", res, id)
	}

	return parts[0], parts[1], nil
}

// getOkExists is a helper function that replaces the deprecated GetOkExists
// schema method. It returns the value of the key if it exists in the
// configuration, along with a boolean indicating if the key exists.
func getOkExists(d *schema.ResourceData, key string) (any, bool) {
	v := d.GetRawConfig().GetAttr(key)
	if v.IsNull() {
		return nil, false
	}

	return d.Get(key), true
}
