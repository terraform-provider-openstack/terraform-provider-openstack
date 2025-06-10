package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitDataSourceValidateImageSortFilter(t *testing.T) {
	var valid1 any = "name:asc"

	var valid2 any = "name:asc,status"

	var valid3 any = "name:desc,owner,created_at:asc"

	var invalid1 any = "hello"

	var invalid2 any = "hello:world"

	var invalid3 any = "name,hello:asc,owner:world"

	_, errs := dataSourceValidateImageSortFilter(valid1, "sort")
	assert.Empty(t, errs)

	_, errs = dataSourceValidateImageSortFilter(valid2, "sort")
	assert.Empty(t, errs)

	_, errs = dataSourceValidateImageSortFilter(valid3, "sort")
	assert.Empty(t, errs)

	_, errs = dataSourceValidateImageSortFilter(invalid1, "sort")
	assert.Len(t, errs, 1)
	require.Error(t, errs[0])

	_, errs = dataSourceValidateImageSortFilter(invalid2, "sort")
	assert.Len(t, errs, 2)
	require.Error(t, errs[0])
	require.Error(t, errs[1])

	_, errs = dataSourceValidateImageSortFilter(invalid3, "sort")
	assert.Len(t, errs, 2)
	require.Error(t, errs[0])
	require.Error(t, errs[1])
}
