package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitDataSourceValidateImageSortFilter(t *testing.T) {
	var valid1 interface{} = "name:asc"
	var valid2 interface{} = "name:asc,status"
	var valid3 interface{} = "name:desc,owner,created_at:asc"

	var invalid1 interface{} = "hello"
	var invalid2 interface{} = "hello:world"
	var invalid3 interface{} = "name,hello:asc,owner:world"

	_, errs := dataSourceValidateImageSortFilter(valid1, "sort")
	assert.Len(t, errs, 0)

	_, errs = dataSourceValidateImageSortFilter(valid2, "sort")
	assert.Len(t, errs, 0)

	_, errs = dataSourceValidateImageSortFilter(valid3, "sort")
	assert.Len(t, errs, 0)

	_, errs = dataSourceValidateImageSortFilter(invalid1, "sort")
	assert.Len(t, errs, 1)
	assert.Error(t, errs[0])

	_, errs = dataSourceValidateImageSortFilter(invalid2, "sort")
	assert.Len(t, errs, 2)
	assert.Error(t, errs[0])
	assert.Error(t, errs[1])

	_, errs = dataSourceValidateImageSortFilter(invalid3, "sort")
	assert.Len(t, errs, 2)
	assert.Error(t, errs[0])
	assert.Error(t, errs[1])
}
