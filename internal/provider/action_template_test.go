// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

func TestActionTemplate_Render(t *testing.T) {
	t.Parallel()

	template := `
Printf tffile: {{ printf "{{tffile %q}}" .ExampleFile }}
tffile: {{ tffile .ExampleFile }}
`
	expectedString := `
Printf tffile: {{tffile "action.tf"}}
tffile: terraform
action "scaffolding_example" "example1" {
  config {
    required_attr = "value-1"
  }
}

`

	tpl := actionTemplate(template)

	schema := tfjson.ActionSchema{
		Block: &tfjson.SchemaBlock{
			Attributes: map[string]*tfjson.SchemaAttribute{
				"required_attr": {
					AttributeType: cty.String,
					Description:   "Required attribute",
					Required:      true,
				},
				"optional_attr": {
					AttributeType: cty.String,
					Description:   "Optional attribute",
					Optional:      true,
				},
			},
		},
	}

	result, err := tpl.Render("testdata/test-action-dir", "testTemplate", "test-action", "test-action", "action", "action.tf", []string{"action.tf"}, &schema)
	if err != nil {
		t.Error(err)
	}

	cleanedResult := strings.ReplaceAll(result, "```", "")
	if !cmp.Equal(expectedString, cleanedResult) {
		t.Errorf("expected: %+v, got: %+v", expectedString, cleanedResult)
	}
}
