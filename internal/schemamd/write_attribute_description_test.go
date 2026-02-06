// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schemamd_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-docs/internal/schemamd"
)

func TestWriteAttributeDescription(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		expected string
		att      *tfjson.SchemaAttribute
	}{
		// required
		{
			"(String, Required) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Required:      true,
				Description:   "This is an attribute.",
			},
		},
		{
			"(String, Required, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Required:      true,
				Description:   "This is an attribute.",
				WriteOnly:     true,
			},
		},
		{
			"(String, Required, Deprecated) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Required:      true,
				Description:   "This is an attribute.",
				Deprecated:    true,
			},
		},
		{
			"(String, Required, Sensitive, Deprecated) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Required:      true,
				Description:   "This is an attribute.",
				Deprecated:    true,
				Sensitive:     true,
			},
		},

		// optional
		{
			"(String, Optional) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Optional:      true,
				Description:   "This is an attribute.",
			},
		},
		{
			"(String, Optional) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Optional:      true,
				Computed:      true,
				Description:   "This is an attribute.",
			},
		},
		{
			"(String, Optional, Deprecated) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Optional:      true,
				Description:   "This is an attribute.",
				Deprecated:    true,
			},
		},
		{
			"(String, Optional, Deprecated) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Optional:      true,
				Computed:      true,
				Description:   "This is an attribute.",
				Deprecated:    true,
			},
		},
		{
			"(String, Optional, Sensitive, Deprecated) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Optional:      true,
				Computed:      true,
				Description:   "This is an attribute.",
				Deprecated:    true,
				Sensitive:     true,
			},
		},

		// computed
		{
			"(String, Read-only) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Computed:      true,
				Description:   "This is an attribute.",
			},
		},
		{
			"(String, Read-only, Deprecated) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Computed:      true,
				Description:   "This is an attribute.",
				Deprecated:    true,
			},
		},
		{
			"(String, Read-only, Sensitive, Deprecated) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Computed:      true,
				Description:   "This is an attribute.",
				Deprecated:    true,
				Sensitive:     true,
			},
		},

		// white space in descriptions
		{
			"(String, Required) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Required:      true,
				Description:   " This is an attribute.",
			},
		},
		{
			"(String, Required) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Required:      true,
				Description:   "This is an attribute. ",
			},
		},
		{
			"(String, Required) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Required:      true,
				Description:   "\n\t This is an attribute.\n\t ",
			},
		},
	} {
		t.Run(c.expected, func(t *testing.T) {
			t.Parallel()

			b := &strings.Builder{}
			err := schemamd.WriteAttributeDescription(b, c.att, true)
			if err != nil {
				t.Fatal(err)
			}
			actual := b.String()
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}
