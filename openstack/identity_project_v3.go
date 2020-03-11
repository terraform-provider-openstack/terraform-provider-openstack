package openstack

func resourceIdentityProjectV3BuildTags(v []interface{}) *[]string {
	tags := make([]string, len(v))
	for i, tag := range v {
		tags[i] = tag.(string)
	}

	return &tags
}
