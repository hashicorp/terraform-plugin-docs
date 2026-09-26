// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestRenderStringTemplate(t *testing.T) {
	t.Parallel()

	template := `
Plainmarkdown: {{ plainmarkdown .Text }}
Split: {{ $arr := split .Text " "}}{{ index $arr 3 }}
Trimspace: {{ trimspace .Text }}
Lower: {{ upper .Text }}
Upper: {{ lower .Text }}
Title: {{ title .Text }}
Prefixlines:
{{ prefixlines "  " .MultiLineTest }}
Printf tffile: {{ printf "{{tffile %q}}" .Code }}
tffile: {{ tffile .Code }}
`

	expectedString := `
Plainmarkdown: my Odly cAsed striNg
Split: striNg
Trimspace: my Odly cAsed striNg
Lower: MY ODLY CASED STRING
Upper: my odly cased string
Title: My Odly Cased String
Prefixlines:
  This text used
  multiple lines
Printf tffile: {{tffile "provider.tf"}}
tffile: terraform
provider "scaffolding" {
  # example configuration here
}

`
	result, err := renderStringTemplate("testdata/test-provider-dir", "testTemplate", template, struct {
		Text          string
		MultiLineTest string
		Code          string
	}{
		Text: "my Odly cAsed striNg",
		MultiLineTest: `This text used
multiple lines`,
		Code: "provider.tf",
	})

	if err != nil {
		t.Error(err)
	}
	cleanedResult := strings.ReplaceAll(result, "```", "")
	if !cmp.Equal(expectedString, cleanedResult) {
		t.Errorf("expected: %+v, got: %+v", expectedString, cleanedResult)
	}
}

func TestResourceTemplate_Render(t *testing.T) {
	t.Parallel()

	template := `
Printf tffile: {{ printf "{{tffile %q}}" .ExampleFile }}
tffile: {{ tffile .ExampleFile }}
`
	expectedString := `
Printf tffile: {{tffile "provider.tf"}}
tffile: terraform
provider "scaffolding" {
  # example configuration here
}

`

	tpl := resourceTemplate(template)

	schema := tfjson.Schema{
		Version: 3,
		Block: &tfjson.SchemaBlock{
			Attributes:      nil,
			NestedBlocks:    nil,
			Description:     "",
			DescriptionKind: "",
			Deprecated:      false,
		},
	}

	result, err := tpl.Render("testdata/test-provider-dir", "testTemplate", "test-provider", "test-provider", "Resource", "provider.tf", []string{"provider.tf"}, "", "", "", &schema, nil, false)
	if err != nil {
		t.Error(err)
	}

	cleanedResult := strings.ReplaceAll(result, "```", "")
	if !cmp.Equal(expectedString, cleanedResult) {
		t.Errorf("expected: %+v, got: %+v", expectedString, cleanedResult)
	}
}

func TestProviderTemplate_Render(t *testing.T) {
	t.Parallel()

	template := `
Printf tffile: {{ printf "{{tffile %q}}" .ExampleFile }}
tffile: {{ tffile .ExampleFile }}
`
	expectedString := `
Printf tffile: {{tffile "provider.tf"}}
tffile: terraform
provider "scaffolding" {
  # example configuration here
}

`

	tpl := providerTemplate(template)

	schema := tfjson.Schema{
		Version: 3,
		Block: &tfjson.SchemaBlock{
			Attributes:      nil,
			NestedBlocks:    nil,
			Description:     "",
			DescriptionKind: "",
			Deprecated:      false,
		},
	}

	result, err := tpl.Render("testdata/test-provider-dir", "testTemplate", "test-provider", "provider.tf", []string{"provider.tf"}, &schema, false)
	if err != nil {
		t.Error(err)
	}

	cleanedResult := strings.ReplaceAll(result, "```", "")
	if !cmp.Equal(expectedString, cleanedResult) {
		t.Errorf("expected: %+v, got: %+v", expectedString, cleanedResult)
	}
}
