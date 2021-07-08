package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

func TestFlattenIdentityAuthScopeV3Roles(t *testing.T) {
	roles := []tokens.Role{
		{
			ID:   "1",
			Name: "foo",
		},
		{
			ID:   "2",
			Name: "bar",
		},
	}

	expected := []map[string]string{
		{
			"role_id":   "1",
			"role_name": "foo",
		},
		{
			"role_name": "bar",
			"role_id":   "2",
		},
	}

	actual := flattenIdentityAuthScopeV3Roles(roles)
	assert.Equal(t, expected, actual)
}

func TestFlattenIdentityAuthScopeV3ServiceCatalog(t *testing.T) {
	cinderEndpoints := []tokens.Endpoint{
		{
			ID:        "3",
			Interface: "public",
			Region:    "oakland",
			RegionID:  "oak",
			URL:       "http://oak.example.com/public_cinder",
		},
		{
			ID:        "4",
			Interface: "internal",
			Region:    "oakland",
			RegionID:  "oak",
			URL:       "http://oak.example.com/internal_cinder",
		},
	}
	keystoneEndpoints := []tokens.Endpoint{
		{
			ID:        "5",
			Interface: "public",
			Region:    "oakland",
			RegionID:  "oak",
			URL:       "http://oak.example.com/public_keystone",
		},
		{
			ID:        "6",
			Interface: "internal",
			Region:    "oakland",
			RegionID:  "oak",
			URL:       "http://oak.example.com/internal_keystone",
		},
	}

	catalogEntries := []tokens.CatalogEntry{
		{
			ID:        "1",
			Name:      "cinderv2",
			Type:      "volumev2",
			Endpoints: cinderEndpoints,
		},
		{
			ID:        "2",
			Name:      "keystone",
			Type:      "identity",
			Endpoints: keystoneEndpoints,
		},
	}
	catalog := tokens.ServiceCatalog{
		Entries: catalogEntries,
	}

	expected := []map[string]interface{}{
		{
			"id":   "1",
			"name": "cinderv2",
			"type": "volumev2",
			"endpoints": []map[string]string{
				{
					"id":        "3",
					"interface": "public",
					"region":    "oakland",
					"region_id": "oak",
					"url":       "http://oak.example.com/public_cinder",
				},
				{
					"id":        "4",
					"interface": "internal",
					"region":    "oakland",
					"region_id": "oak",
					"url":       "http://oak.example.com/internal_cinder",
				},
			},
		},
		{
			"id":   "2",
			"name": "keystone",
			"type": "identity",
			"endpoints": []map[string]string{
				{
					"id":        "5",
					"interface": "public",
					"region":    "oakland",
					"region_id": "oak",
					"url":       "http://oak.example.com/public_keystone",
				},
				{
					"id":        "6",
					"interface": "internal",
					"region":    "oakland",
					"region_id": "oak",
					"url":       "http://oak.example.com/internal_keystone",
				},
			},
		},
	}

	actual := flattenIdentityAuthScopeV3ServiceCatalog(&catalog)
	assert.Equal(t, expected, actual)
}
