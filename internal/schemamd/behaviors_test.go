// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package schemamd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

func TestChildAttributeIsRequired(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name     string
		att      *tfjson.SchemaAttribute
		expected bool
	}{
		{
			"required",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Required:      true,
				Description:   "This is an attribute.",
			},
			true,
		},
		{
			"not required",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Optional:      true,
				Description:   "This is an attribute.",
			},
			false,
		},
	} {

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := childAttributeIsRequired(c.att)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}

func TestChildAttributeIsOptional(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name     string
		att      *tfjson.SchemaAttribute
		expected bool
	}{
		{
			"not optional",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Required:      true,
				Description:   "This is an attribute.",
			},
			false,
		},
		{
			"optional",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Optional:      true,
				Description:   "This is an attribute.",
			},
			true,
		},
	} {

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := childAttributeIsOptional(c.att)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}

func TestChildAttributeIsReadOnly(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name     string
		att      *tfjson.SchemaAttribute
		expected bool
	}{
		{
			"required, not compted",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Required:      true,
				Description:   "This is an attribute.",
			},
			false,
		},
		{
			"optional, not computed",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Optional:      true,
				Description:   "This is an attribute.",
			},
			false,
		},
		{
			"optional, computed",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Optional:      true,
				Description:   "This is an attribute.",
			},
			false,
		},
		{
			"computed",
			&tfjson.SchemaAttribute{
				AttributeType: cty.String,
				Computed:      true,
				Description:   "This is an attribute.",
			},
			true,
		},
	} {

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := childAttributeIsReadOnly(c.att)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}

func TestChildBlockIsRequired(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name     string
		block    *tfjson.SchemaBlockType
		expected bool
	}{
		{
			"required",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
							Description:   "This is an attribute.",
						},
					},
				},
				MinItems: 1,
			},
			true,
		},
		{
			"not required",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
							Description:   "This is an attribute.",
						},
					},
				},
			},
			false,
		},
	} {

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := childBlockIsRequired(c.block)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}

func TestChildBlockIsOptional(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name     string
		block    *tfjson.SchemaBlockType
		expected bool
	}{
		{
			"min items 1",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
							Description:   "This is an attribute.",
						},
					},
				},
				MinItems: 1,
			},
			false,
		},
		{
			"required child attribute",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
							Description:   "This is an attribute.",
						},
					},
				},
			},
			true,
		},
		{
			"optional child attribute",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Optional:      true,
							Description:   "This is an attribute.",
						},
					},
				},
			},
			true,
		},
		{
			"computed child attribute",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Computed:      true,
							Description:   "This is an attribute.",
						},
					},
				},
			},
			false,
		},
		{
			"child block min items 1",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					NestedBlocks: map[string]*tfjson.SchemaBlockType{
						"foo": {
							Block: &tfjson.SchemaBlock{
								Attributes: map[string]*tfjson.SchemaAttribute{
									"foo": {
										AttributeType: cty.String,
										Required:      true,
										Description:   "This is an attribute.",
									},
								},
							},
							MinItems: 1,
						},
					},
				},
			},
			true,
		},
		{
			"optional child block attribute",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					NestedBlocks: map[string]*tfjson.SchemaBlockType{
						"foo": {
							Block: &tfjson.SchemaBlock{
								Attributes: map[string]*tfjson.SchemaAttribute{
									"foo": {
										AttributeType: cty.String,
										Optional:      true,
										Description:   "This is an attribute.",
									},
								},
							},
						},
					},
				},
			},
			true,
		},
		{
			"computed child block attribute",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					NestedBlocks: map[string]*tfjson.SchemaBlockType{
						"foo": {
							Block: &tfjson.SchemaBlock{
								Attributes: map[string]*tfjson.SchemaAttribute{
									"foo": {
										AttributeType: cty.String,
										Computed:      true,
										Description:   "This is an attribute.",
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"empty block",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is an empty block.",
				},
			},
			true,
		},
	} {

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := childBlockIsOptional(c.block)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}

func TestChildBlockIsReadOnly(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name     string
		block    *tfjson.SchemaBlockType
		expected bool
	}{
		{
			"max items 1",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
							Description:   "This is an attribute.",
						},
					},
				},
				MaxItems: 1,
			},
			false,
		},
		{
			"required child attribute",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Required:      true,
							Description:   "This is an attribute.",
						},
					},
				},
			},
			false,
		},
		{
			"optional child attribute",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Optional:      true,
							Description:   "This is an attribute.",
						},
					},
				},
			},
			false,
		},
		{
			"computed child attribute",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					Attributes: map[string]*tfjson.SchemaAttribute{
						"foo": {
							AttributeType: cty.String,
							Computed:      true,
							Description:   "This is an attribute.",
						},
					},
				},
			},
			true,
		},
		{
			"child block min items 1",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					NestedBlocks: map[string]*tfjson.SchemaBlockType{
						"foo": {
							Block: &tfjson.SchemaBlock{
								Attributes: map[string]*tfjson.SchemaAttribute{
									"foo": {
										AttributeType: cty.String,
										Required:      true,
										Description:   "This is an attribute.",
									},
								},
							},
							MinItems: 1,
						},
					},
				},
			},
			false,
		},
		{
			"optional child block attribute",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					NestedBlocks: map[string]*tfjson.SchemaBlockType{
						"foo": {
							Block: &tfjson.SchemaBlock{
								Attributes: map[string]*tfjson.SchemaAttribute{
									"foo": {
										AttributeType: cty.String,
										Optional:      true,
										Description:   "This is an attribute.",
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"computed child block attribute",
			&tfjson.SchemaBlockType{
				NestingMode: tfjson.SchemaNestingModeSingle,
				Block: &tfjson.SchemaBlock{
					Description: "This is a block.",
					NestedBlocks: map[string]*tfjson.SchemaBlockType{
						"foo": {
							Block: &tfjson.SchemaBlock{
								Attributes: map[string]*tfjson.SchemaAttribute{
									"foo": {
										AttributeType: cty.String,
										Computed:      true,
										Description:   "This is an attribute.",
									},
								},
							},
						},
					},
				},
			},
			true,
		},
	} {

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			actual := childBlockIsReadOnly(c.block)
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}
