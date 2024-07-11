package openstack

import "github.com/gophercloud/gophercloud/v2"

func identityEndpointAvailability(v string) gophercloud.Availability {
	availability := gophercloud.AvailabilityPublic

	switch v {
	case "internal":
		availability = gophercloud.AvailabilityInternal
	case "admin":
		availability = gophercloud.AvailabilityAdmin
	}

	return availability
}
