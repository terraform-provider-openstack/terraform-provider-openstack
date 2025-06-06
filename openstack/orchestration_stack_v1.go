package openstack

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stacks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func buildTE(t map[string]any) (*stacks.TE, error) {
	log.Printf("[DEBUG] Start to build TE structure")

	te := &stacks.TE{}

	if t["Bin"] != nil {
		if v, ok := t["Bin"].(string); ok {
			te.Bin = []byte(v)
		} else {
			return nil, errors.New("Bin value is expected to be a string")
		}
	}

	if t["URL"] != nil {
		if v, ok := t["URL"].(string); ok {
			te.URL = v
		} else {
			return nil, errors.New("URL value is expected to be a string")
		}
	}

	if t["Files"] != nil {
		if v, ok := t["Files"].(map[string]string); ok {
			te.Files = v
		} else {
			return nil, errors.New("URL value is expected to be a map of string")
		}
	}

	log.Printf("[DEBUG] TE structure builded")

	return te, nil
}

func buildTemplateOpts(d *schema.ResourceData) (*stacks.Template, error) {
	log.Printf("[DEBUG] Start building TemplateOpts")

	te, err := buildTE(d.Get("template_opts").(map[string]any))
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
		t := d.Get("environment_opts").(map[string]any)

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

func orchestrationStackV1StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, stackID string, isdelete bool) retry.StateRefreshFunc {
	return func() (any, string, error) {
		log.Printf("[DEBUG] Refresh Stack status %s", stackID)

		stack, err := stacks.Find(ctx, client, stackID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) && isdelete {
				return stack, "DELETE_COMPLETE", nil
			}

			return nil, "", err
		}

		if strings.Contains(stack.Status, "FAILED") {
			return stack, stack.Status, errors.New("The stack is in error status. " +
				"Please check with your cloud admin or check the orchestration " +
				"API logs to see why this error occurred.")
		}

		return stack, stack.Status, nil
	}
}
