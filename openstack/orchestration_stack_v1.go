package openstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacks"
)

func buildTE(t map[string]interface{}) (*stacks.TE, error) {
	log.Printf("[DEBUG] Start to build TE structure")
	te := &stacks.TE{}
	if t["Bin"] != nil {
		if v, ok := t["Bin"].(string); ok {
			te.Bin = []byte(v)
		} else {
			return nil, fmt.Errorf("Bin value is expected to be a string")
		}
	}
	if t["URL"] != nil {
		if v, ok := t["URL"].(string); ok {
			te.URL = v
		} else {
			return nil, fmt.Errorf("URL value is expected to be a string")
		}
	}
	if t["Files"] != nil {
		if v, ok := t["Files"].(map[string]string); ok {
			te.Files = v
		} else {
			return nil, fmt.Errorf("URL value is expected to be a map of string")
		}
	}
	log.Printf("[DEBUG] TE structure builded")
	return te, nil
}

func buildTemplateOpts(d *schema.ResourceData) (*stacks.Template, error) {
	log.Printf("[DEBUG] Start building TemplateOpts")
	te, err := buildTE(d.Get("template_opts").(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Return TemplateOpts")
	return &stacks.Template{
		TE: *te,
	}, nil
}

func buildEnvironmentOpts(d *schema.ResourceData) (*stacks.Environment, error) {
	log.Printf("[DEBUG] Start building EnvironmentOpts")
	if d.Get("environment_opts") != nil {
		t := d.Get("environment_opts").(map[string]interface{})
		te, err := buildTE(t)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] Return EnvironmentOpts")
		return &stacks.Environment{
			TE: *te,
		}, nil
	}
	return nil, nil
}

func orchestrationStackV1StateRefreshFunc(client *gophercloud.ServiceClient, stackID string, isdelete bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Refresh Stack status %s", stackID)
		stack, err := stacks.Find(client, stackID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok && isdelete {
				return stack, "DELETE_COMPLETE", nil
			}

			return nil, "", err
		}

		if strings.Contains(stack.Status, "FAILED") {
			return stack, stack.Status, fmt.Errorf("The stack is in error status. " +
				"Please check with your cloud admin or check the orchestration " +
				"API logs to see why this error occurred.")
		}

		return stack, stack.Status, nil
	}
}
