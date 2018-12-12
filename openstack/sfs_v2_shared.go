package openstack

import (
	"fmt"
	"log"

	version "github.com/hashicorp/go-version"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/apiversions"
)

const (
	minManilaMicroversion   = "2.7"
	minOUManilaMicroversion = "2.44"
)

func setManilaMicroversion(sfsClient *gophercloud.ServiceClient) (bool, error) {
	apiVersions, err := apiversions.Get(sfsClient, "v2").Extract()
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] Minimum Manila API microversion: %s", apiVersions.MinVersion)
	log.Printf("[DEBUG] Current Manila API microversion: %s", apiVersions.Version)

	minReq, err := version.NewVersion(minManilaMicroversion)
	if err != nil {
		return false, err
	}
	req, err := version.NewVersion(minOUManilaMicroversion)
	if err != nil {
		return false, err
	}
	min, err := version.NewVersion(apiVersions.MinVersion)
	if err != nil {
		return false, err
	}
	curr, err := version.NewVersion(apiVersions.Version)
	if err != nil {
		return false, err
	}

	if (req.Equal(min) || req.GreaterThan(min)) && (req.Equal(curr) || req.LessThan(curr)) {
		sfsClient.Microversion = minOUManilaMicroversion
		return true, nil
	}

	if (minReq.Equal(min) || minReq.GreaterThan(min)) && (minReq.Equal(curr) || minReq.LessThan(curr)) {
		sfsClient.Microversion = minManilaMicroversion
		return true, fmt.Errorf("Organizational Unit field is not supported by %s Manila API microversion", apiVersions.Version)
	}

	return false, fmt.Errorf("Mimimum required %s Manila API microversion is not between %s and %s supported microversion range", minManilaMicroversion, apiVersions.MinVersion, apiVersions.Version)
}
