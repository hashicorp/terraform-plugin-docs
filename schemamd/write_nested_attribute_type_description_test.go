package schemamd_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-docs/schemamd"
	"github.com/zclconf/go-cty/cty"
)

func TestWriteNestedAttributeTypeDescription(t *testing.T) {
	for _, c := range []struct {
		expected string
		att      *tfjson.SchemaAttribute
	}{
		{
			"(Attributes, Optional) This is an attribute.",
			&tfjson.SchemaAttribute{
				Description: "This is an attribute.",
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeSingle,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
							Description:   "This is a nested attribute.",
						},
					},
				},
				Optional: true,
			},
		},
		{
			"(Attributes List, Min: 2, Max: 3) This is an attribute.",
			&tfjson.SchemaAttribute{
				Description: "This is an attribute.",
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeList,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
							Description:   "This is a nested attribute.",
						},
					},
					MinItems: 2,
					MaxItems: 3,
				},
				Required: true,
			},
		},
		{
			"(Attributes Map) This is an attribute.",
			&tfjson.SchemaAttribute{
				Description: "This is an attribute.",
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeMap,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
							Description:   "This is a nested attribute.",
						},
					},
				},
			},
		},
		{
			"(Attributes Set, Min: 5) This is an attribute.",
			&tfjson.SchemaAttribute{
				Description: "This is an attribute.",
				AttributeNestedType: &tfjson.SchemaNestedAttributeType{
					NestingMode: tfjson.SchemaNestingModeSet,
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
							Description:   "This is a nested attribute.",
						},
					},
					MinItems: 5,
				},
			},
		},
	} {
		t.Run(c.expected, func(t *testing.T) {
			b := &strings.Builder{}
			err := schemamd.WriteNestedAttributeTypeDescription(b, c.att, true)
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
