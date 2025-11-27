// Copyright (c) HashiCorp, Inc.
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
	}, "terraform")

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

	result, err := tpl.Render("testdata/test-provider-dir", "testTemplate", "test-provider", "test-provider", "Resource", "provider.tf", []string{"provider.tf"}, "", "", "", &schema, nil, "terraform")
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

	result, err := tpl.Render("testdata/test-provider-dir", "testTemplate", "test-provider", "provider.tf", []string{"provider.tf"}, &schema, "terraform")
	if err != nil {
		t.Error(err)
	}

	cleanedResult := strings.ReplaceAll(result, "```", "")
	if !cmp.Equal(expectedString, cleanedResult) {
		t.Errorf("expected: %+v, got: %+v", expectedString, cleanedResult)
	}
}

func TestTffileSyntax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		template       string
		syntax         string
		expectedSyntax string
	}{
		{
			name:           "default syntax with terraform",
			template:       `{{ tffile .Code }}`,
			syntax:         "terraform",
			expectedSyntax: "terraform",
		},
		{
			name:           "default syntax with hcl",
			template:       `{{ tffile .Code }}`,
			syntax:         "hcl",
			expectedSyntax: "hcl",
		},
		{
			name:           "override with hcl",
			template:       `{{ tffile .Code "hcl" }}`,
			syntax:         "terraform",
			expectedSyntax: "hcl",
		},
		{
			name:           "override with terraform",
			template:       `{{ tffile .Code "terraform" }}`,
			syntax:         "hcl",
			expectedSyntax: "terraform",
		},
		{
			name:           "override with custom syntax",
			template:       `{{ tffile .Code "mySyntax" }}`,
			syntax:         "terraform",
			expectedSyntax: "mySyntax",
		},
		{
			name: "multiple tffile calls with different syntaxes",
			template: `Default: {{ tffile .Code }}
Override: {{ tffile .Code "hcl" }}
Custom: {{ tffile .Code "custom" }}`,
			syntax:         "terraform",
			expectedSyntax: "terraform", // We'll check for all three in the result
		},
		{
			name:           "override with variable format",
			template:       `{{ $format := .Format }}{{ tffile .Code $format }}`,
			syntax:         "terraform",
			expectedSyntax: "hcl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := struct {
				Code   string
				Format string
			}{
				Code:   "provider.tf",
				Format: "hcl", // Used for the variable format test case
			}

			result, err := renderStringTemplate("testdata/test-provider-dir", "testTemplate", tt.template, data, tt.syntax)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check that the expected syntax appears in the code block
			if !strings.Contains(result, "```"+tt.expectedSyntax) {
				t.Errorf("expected syntax %q in result, got: %s", tt.expectedSyntax, result)
			}

			// For the multiple calls test, verify all three syntaxes appear
			if tt.name == "multiple tffile calls with different syntaxes" {
				if !strings.Contains(result, "```terraform") {
					t.Errorf("expected terraform syntax in result, got: %s", result)
				}
				if !strings.Contains(result, "```hcl") {
					t.Errorf("expected hcl syntax in result, got: %s", result)
				}
				if !strings.Contains(result, "```custom") {
					t.Errorf("expected custom syntax in result, got: %s", result)
				}
			}
		})
	}
}
