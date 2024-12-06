// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemamd_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-docs/internal/schemamd"
)

func TestWriteBlockTypeDescription(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		expected string
		bt       *tfjson.SchemaBlockType
	}{
		// single
		{
			"(Block, Optional) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							Required: true,
						},
					},
				},
			},
		},
		{
			"(Block, Required) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				MinItems:    1,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block, Required, Deprecated) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				MinItems:    1,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Deprecated:  true,
				},
			},
		},

		// list
		{
			"(Block List) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeList,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block List, Min: 1) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeList,
				MinItems:    1,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block List, Max: 4) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeList,
				MaxItems:    4,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block List, Min: 1, Max: 4) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeList,
				MinItems:    1,
				MaxItems:    4,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block List, Min: 1, Max: 4, Deprecated) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeList,
				MinItems:    1,
				MaxItems:    4,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Deprecated:  true,
				},
			},
		},

		// set
		{
			"(Block Set) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSet,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block Set, Min: 1) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSet,
				MinItems:    1,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block Set, Max: 4) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSet,
				MaxItems:    4,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block Set, Min: 1, Max: 4) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSet,
				MinItems:    1,
				MaxItems:    4,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block Set, Min: 1, Max: 4, Deprecated) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSet,
				MinItems:    1,
				MaxItems:    4,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Deprecated:  true,
				},
			},
		},

		// map
		{
			"(Block Map) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeMap,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block Map, Min: 1) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeMap,
				MinItems:    1,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block Map, Max: 4) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeMap,
				MaxItems:    4,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block Map, Min: 1, Max: 4) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeMap,
				MinItems:    1,
				MaxItems:    4,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
				},
			},
		},
		{
			"(Block Map, Min: 1, Max: 4, Deprecated) This is a block.",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeMap,
				MinItems:    1,
				MaxItems:    4,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Deprecated:  true,
				},
			},
		},
	} {
		c := c
		t.Run(c.expected, func(t *testing.T) {
			t.Parallel()

			b := &strings.Builder{}
			err := schemamd.WriteBlockTypeDescription(b, c.bt)
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
