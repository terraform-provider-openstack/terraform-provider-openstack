package openstack

import (
	"github.com/gophercloud/gophercloud/v2"
)

// TODO: implement this in gophercloud
// TODO: networks are not supported by the API
// TODO: local_as is not string, but int
// TODO: add TenantID to gophercloud.
type speakersCreateOpts struct {
	Name                          string   `json:"name,omitempty"`
	IPVersion                     int      `json:"ip_version,omitempty"`
	AdvertiseFloatingIPHostRoutes *bool    `json:"advertise_floating_ip_host_routes,omitempty"`
	AdvertiseTenantNetworks       *bool    `json:"advertise_tenant_networks,omitempty"`
	LocalAS                       int      `json:"local_as"`
	Networks                      []string `json:"networks,omitempty"`
	TenantID                      string   `json:"tenant_id,omitempty"`
}

func (opts speakersCreateOpts) ToSpeakerCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "bgp_speaker")
}

// TODO: implement this in gophercloud.
type speakersUpdateOpts struct {
	Name                          *string `json:"name,omitempty"`
	AdvertiseFloatingIPHostRoutes *bool   `json:"advertise_floating_ip_host_routes,omitempty"`
	AdvertiseTenantNetworks       *bool   `json:"advertise_tenant_networks,omitempty"`
}

func (opts speakersUpdateOpts) ToSpeakerUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "bgp_speaker")
}

// TODO: implement TenantID in gophercloud.
type peersCreateOpts struct {
	Name     string `json:"name,omitempty"`
	AuthType string `json:"auth_type,omitempty"`
	Password string `json:"password,omitempty"`
	PeerIP   string `json:"peer_ip,omitempty"`
	RemoteAS int    `json:"remote_as,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
}

func (opts peersCreateOpts) ToPeerCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "bgp_peer")
}

// TODO: implement this in gophercloud.
type peersUpdateOpts struct {
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (opts peersUpdateOpts) ToPeerUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "bgp_peer")
}
