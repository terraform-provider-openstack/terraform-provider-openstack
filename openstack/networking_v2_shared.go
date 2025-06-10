package openstack

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func networkingV2ReadAttributesTags(d *schema.ResourceData, tags []string) {
	expandObjectReadTags(d, tags)
}

func networkingV2UpdateAttributesTags(d *schema.ResourceData) []string {
	return expandObjectUpdateTags(d)
}

func networkingV2AttributesTags(d *schema.ResourceData) []string {
	return expandObjectTags(d)
}

type neutronErrorWrap struct {
	NeutronError neutronError `json:"NeutronError"`
}

type neutronError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Detail  string `json:"detail"`
}

func retryOn409(err error) bool {
	var e gophercloud.ErrUnexpectedResponseCode

	ok := errors.As(err, &e)
	if !ok {
		return false
	}

	switch e.Actual {
	case http.StatusConflict: // 409
		neutronError, err := decodeNeutronError(e.Body)
		if err != nil {
			// retry, when error type cannot be detected
			log.Printf("[DEBUG] failed to decode a neutron error: %s", err)

			return true
		}

		if neutronError.Type == "IpAddressGenerationFailure" {
			return true
		}

		// don't retry on quota or other errors
		return false
	case http.StatusBadRequest: // 400
		neutronError, err := decodeNeutronError(e.Body)
		if err != nil {
			// retry, when error type cannot be detected
			log.Printf("[DEBUG] failed to decode a neutron error: %s", err)

			return true
		}

		if neutronError.Type == "ExternalIpAddressExhausted" {
			return true
		}

		// don't retry on quota or other errors
		return false
	case http.StatusNotFound: // this case is handled mostly for functional tests
		return true
	}

	return false
}

func decodeNeutronError(body []byte) (*neutronError, error) {
	e := &neutronErrorWrap{}
	if err := json.Unmarshal(body, e); err != nil {
		return nil, err
	}

	return &e.NeutronError, nil
}
