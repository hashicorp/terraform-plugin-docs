package schemamd_test

import (
	"strings"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-docs/schemamd"
)

func TestWriteAttributeDescription(t *testing.T) {
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
			"(String, Required, Deprecated) This is an attribute.",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Required:      true,
				Description:   "This is an attribute.",
				Deprecated:    true,
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
			b := &strings.Builder{}
			err := schemamd.WriteAttributeDescription(b, c.att)
			if err != nil {
				t.Fatal(err)
			}
			actual := b.String()
			if c.expected != actual {
				t.Fatalf("expected %q, got %q", c.expected, actual)
			}
		})
	}
}
