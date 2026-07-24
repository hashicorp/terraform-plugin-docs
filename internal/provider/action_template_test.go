// Copyright IBM Corp. 2020, 2026
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

	result, err := tpl.Render("testdata/test-action-dir", "testTemplate", "test-action", "test-action", "action", "action.tf", []string{"action.tf"}, "", &schema)
	if err != nil {
		t.Error(err)
	}

	cleanedResult := strings.ReplaceAll(result, "```", "")
	if !cmp.Equal(expectedString, cleanedResult) {
		t.Errorf("expected: %+v, got: %+v", expectedString, cleanedResult)
	}
}

func TestActionTemplate_Render_WithInvocation(t *testing.T) {
	t.Parallel()

	template := `
Printf codefile: {{ printf "{{codefile %q %q}}" "shell" .InvocationFile }}
codefile: {{ codefile "shell" .InvocationFile }}
`
	expectedString := `
Printf codefile: {{codefile "shell" "invoke.sh"}}
codefile: ` + "```" + `shell
terraform run scaffolding_example.example1
` + "```" + `
`

	tpl := actionTemplate(template)

	schema := tfjson.ActionSchema{
		Block: &tfjson.SchemaBlock{},
	}

	result, err := tpl.Render("testdata/test-action-dir", "testTemplate", "test-action", "test-action", "action", "", nil, "invoke.sh", &schema)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(expectedString, result) {
		t.Errorf("expected: %+v, got: %+v", expectedString, result)
	}
}
