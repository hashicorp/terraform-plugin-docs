// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package schemamd_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-docs/internal/schemamd"
)

func TestRender(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name          string
		inputFile     string
		expectedFile  string
		blocksSection bool
	}{
		{
			"aws_route_table_association",
			"testdata/aws_route_table_association.schema.json",
			"testdata/aws_route_table_association.md",
			false,
		},
		{
			"aws_acm_certificate",
			"testdata/aws_acm_certificate.schema.json",
			"testdata/aws_acm_certificate.md",
			false,
		},
		{
			"awscc_logs_log_group",
			"testdata/awscc_logs_log_group.schema.json",
			"testdata/awscc_logs_log_group.md",
			false,
		},
		{
			"awscc_acmpca_certificate",
			"testdata/awscc_acmpca_certificate.schema.json",
			"testdata/awscc_acmpca_certificate.md",
			false,
		},
		{
			"framework_types",
			"testdata/framework_types.schema.json",
			"testdata/framework_types.md",
			false,
		},
		{
			"framework_types_blocks_section",
			"testdata/framework_types.schema.json",
			"testdata/framework_types_blocks_section.md",
			true,
		},
		{
			// Reference: https://github.com/hashicorp/terraform-plugin-docs/issues/380
			"deep_nested_attributes",
			"testdata/deep_nested_attributes.schema.json",
			"testdata/deep_nested_attributes.md",
			false,
		},
		{
			// Reference: https://github.com/hashicorp/terraform-plugin-docs/issues/380
			"deep_nested_attributes_blocks_section",
			"testdata/deep_nested_attributes.schema.json",
			"testdata/deep_nested_attributes_blocks_section.md",
			true,
		},
		{
			"deep_nested_write_only_attributes",
			"testdata/deep_nested_write_only_attributes.schema.json",
			"testdata/deep_nested_write_only_attributes.md",
			false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			input, err := os.ReadFile(c.inputFile)
			if err != nil {
				t.Fatal(err)
			}

			expected, err := os.ReadFile(c.expectedFile)
			if err != nil {
				t.Fatal(err)
			}

			var schema tfjson.Schema

			err = json.Unmarshal(input, &schema)
			if err != nil {
				t.Fatal(err)
			}

			b := &strings.Builder{}
			err = schemamd.Render(&schema, b, c.blocksSection)
			if err != nil {
				t.Fatal(err)
			}

			// Remove \r characters so tests don't fail on windows
			expectedStr := strings.ReplaceAll(string(expected), "\r", "")

			// Remove trailing newlines before comparing (some text editors remove them).
			expectedStr = strings.TrimRight(expectedStr, "\n")
			actual := strings.TrimRight(b.String(), "\n")
			if diff := cmp.Diff(expectedStr, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}

func TestRenderIdentitySchema(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name         string
		inputFile    string
		expectedFile string
	}{
		{
			"identity",
			"testdata/identity.schema.json",
			"testdata/identity.md",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			input, err := os.ReadFile(c.inputFile)
			if err != nil {
				t.Fatal(err)
			}

			expected, err := os.ReadFile(c.expectedFile)
			if err != nil {
				t.Fatal(err)
			}

			var identitySchema tfjson.IdentitySchema

			err = json.Unmarshal(input, &identitySchema)
			if err != nil {
				t.Fatal(err)
			}

			b := &strings.Builder{}
			err = schemamd.RenderIdentitySchema(&identitySchema, b)
			if err != nil {
				t.Fatal(err)
			}

			// Remove \r characters so tests don't fail on windows
			expectedStr := strings.ReplaceAll(string(expected), "\r", "")

			// Remove trailing newlines before comparing (some text editors remove them).
			expectedStr = strings.TrimRight(expectedStr, "\n")
			actual := strings.TrimRight(b.String(), "\n")
			if diff := cmp.Diff(expectedStr, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}

func TestRenderAction(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name         string
		inputFile    string
		expectedFile string
	}{
		{
			"test_action",
			"testdata/actions/test_action.schema.json",
			"testdata/actions/test_action.md",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			input, err := os.ReadFile(c.inputFile)
			if err != nil {
				t.Fatal(err)
			}

			expected, err := os.ReadFile(c.expectedFile)
			if err != nil {
				t.Fatal(err)
			}

			var schema tfjson.ActionSchema

			err = json.Unmarshal(input, &schema)
			if err != nil {
				t.Fatal(err)
			}

			b := &strings.Builder{}
			err = schemamd.RenderAction(&schema, b)
			if err != nil {
				t.Fatal(err)
			}

			// Remove \r characters so tests don't fail on windows
			expectedStr := strings.ReplaceAll(string(expected), "\r", "")

			// Remove trailing newlines before comparing (some text editors remove them).
			expectedStr = strings.TrimRight(expectedStr, "\n")
			actual := strings.TrimRight(b.String(), "\n")
			if diff := cmp.Diff(expectedStr, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}
