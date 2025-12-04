// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemamd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-docs/internal/schemamd"
)

func TestWriteType(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		expected string
		ty       cty.Type
	}{
		{"Boolean", cty.Bool},
		{"Dynamic", cty.DynamicPseudoType},
		{"Number", cty.Number},
		{"String", cty.String},

		// Currently not supported
		// {"Number", cty.NilType},
		// {"Number", cty.Capsule("foo", reflect.TypeOf(1))},

		{"List of Booleans", cty.List(cty.Bool)},
		{"List of Dynamics", cty.List(cty.DynamicPseudoType)},

		{"Map of Booleans", cty.Map(cty.Bool)},

		{"Object", cty.EmptyObject},
		{"Object", cty.Object(map[string]cty.Type{
			"bool": cty.Bool,
		})},

		{"Set of Booleans", cty.Set(cty.Bool)},

		{"Tuple", cty.EmptyTuple},
		{"Tuple", cty.Tuple([]cty.Type{cty.Bool})},

		{"List of Map of Set of Objects", cty.List(cty.Map(cty.Set(cty.Object(map[string]cty.Type{
			"bool": cty.Bool,
		}))))},
	} {
		t.Run(fmt.Sprintf("%s %s", c.ty.FriendlyName(), c.expected), func(t *testing.T) {
			t.Parallel()

			b := &strings.Builder{}
			err := schemamd.WriteType(b, c.ty)
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
