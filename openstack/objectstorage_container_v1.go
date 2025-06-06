package openstack

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers"
)

type containerCreateOpts struct {
	containers.CreateOpts
	// The storage class of the container. This feature is only available
	// in Ceph Rados Gateway Swift API implementation.
	StorageClass string `h:"X-Object-Storage-Class"`
}

func (opts containerCreateOpts) ToContainerCreateMap() (map[string]string, error) {
	h, err := gophercloud.BuildHeaders(opts.CreateOpts)
	if err != nil {
		return nil, err
	}
	// the BuildHeaders doesn't support nested struct, so we need to add the
	// StorageClass manually.
	if opts.StorageClass != "" {
		h["X-Object-Storage-Class"] = opts.StorageClass
	}

	for k, v := range opts.Metadata {
		h["X-Container-Meta-"+k] = v
	}

	return h, nil
}
